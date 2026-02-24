package evolution

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/livetemplate/lvt/internal/validator"
)

// ValidateFunc validates an app directory and returns the result.
type ValidateFunc func(ctx context.Context, appPath string) *validator.ValidationResult

// Tester validates proposed fixes by applying them in an isolated temp directory.
type Tester struct {
	validate ValidateFunc
}

// NewTester creates a Tester using the given validation function.
// Pass nil to create a tester that skips validation (for unit tests).
func NewTester(validate ValidateFunc) *Tester {
	return &Tester{validate: validate}
}

// TestFix applies a fix in a temp directory and validates the result.
func (t *Tester) TestFix(ctx context.Context, fix Fix, sourceDir string) (*TestResult, error) {
	// Create isolated temp directory
	tempDir, err := os.MkdirTemp("", "lvt-fix-test-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy source files to temp directory
	if err := copyDir(sourceDir, tempDir); err != nil {
		return nil, fmt.Errorf("copy source: %w", err)
	}

	// Apply the fix
	applied, err := applyFix(tempDir, fix)
	if err != nil {
		return &TestResult{
			FixID:   fix.ID,
			Success: false,
			Error:   fmt.Sprintf("apply fix: %v", err),
		}, nil
	}
	if !applied {
		return &TestResult{
			FixID:   fix.ID,
			Success: false,
			Error:   "fix target not found: no files matched pattern or find text not found",
		}, nil
	}

	// Run validation if available
	if t.validate != nil {
		result := t.validate(ctx, tempDir)
		return &TestResult{
			FixID:      fix.ID,
			Success:    !result.HasErrors(),
			Validation: result,
		}, nil
	}

	// No validation function — consider it a pass (fix applied successfully)
	return &TestResult{
		FixID:   fix.ID,
		Success: true,
	}, nil
}

// applyFix finds files matching the fix target glob and performs find/replace.
// Returns true if at least one file was modified.
func applyFix(dir string, fix Fix) (bool, error) {
	// Find matching files
	matches, err := filepath.Glob(filepath.Join(dir, fix.TargetFile))
	if err != nil {
		return false, fmt.Errorf("glob %q: %w", fix.TargetFile, err)
	}

	// Also try recursive walk when the glob has a wildcard directory prefix (e.g. "*/handler.go.tmpl")
	if len(matches) == 0 && strings.HasPrefix(fix.TargetFile, "*/") {
		suffix := strings.TrimPrefix(fix.TargetFile, "*/")
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				rel, relErr := filepath.Rel(dir, path)
				if relErr != nil {
					return nil
				}
				matched, _ := filepath.Match(fix.TargetFile, rel)
				if !matched {
					// Also try matching just the filename for simple patterns
					matched, _ = filepath.Match(suffix, info.Name())
				}
				if matched {
					matches = append(matches, path)
				}
			}
			return nil
		})
		if err != nil {
			return false, err
		}
	}

	applied := false
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)

		var newContent string
		if fix.IsRegex {
			re, err := regexp.Compile(fix.FindPattern)
			if err != nil {
				return false, fmt.Errorf("compile fix regex: %w", err)
			}
			newContent = re.ReplaceAllString(content, fix.Replace)
		} else {
			newContent = strings.ReplaceAll(content, fix.FindPattern, fix.Replace)
		}

		if newContent != content {
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return false, fmt.Errorf("write %s: %w", path, err)
			}
			applied = true
		}
	}

	return applied, nil
}

// copyDir recursively copies src to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
