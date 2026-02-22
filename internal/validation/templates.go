package validation

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/livetemplate/lvt/internal/tmplutil"
	"github.com/livetemplate/lvt/internal/validator"
)

// TemplateCheck validates all .tmpl files in an app directory using
// html/template (matching the existing generator/validate.go convention).
type TemplateCheck struct{}

func (c *TemplateCheck) Name() string { return "templates" }

func (c *TemplateCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	var found bool

	walkErr := filepath.WalkDir(appPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			relPath, _ := filepath.Rel(appPath, path)
			if relPath == "" {
				relPath = path
			}
			result.AddWarning("skipping "+relPath+": "+err.Error(), relPath, 0)
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if ctx.Err() != nil {
			return filepath.SkipAll
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		found = true
		c.validateFile(path, appPath, result)
		return nil
	})

	if walkErr != nil {
		result.AddWarning("template walk incomplete: "+walkErr.Error(), appPath, 0)
	}

	// If the context was cancelled during the walk, record it so the
	// caller knows the result is partial.
	if ctx.Err() != nil {
		result.AddError("validation cancelled: "+ctx.Err().Error(), "", 0)
	}

	if !found {
		result.AddInfo("no .tmpl files found", "", 0)
	}

	return result
}

func (c *TemplateCheck) validateFile(path, appPath string, result *validator.ValidationResult) {
	relPath, _ := filepath.Rel(appPath, path)
	if relPath == "" {
		relPath = path
	}

	content, err := os.ReadFile(path)
	if err != nil {
		result.AddError("failed to read template: "+err.Error(), relPath, 0)
		return
	}

	src := string(content)

	// Parse check.
	_, parseErr := template.New(filepath.Base(path)).Parse(src)
	if parseErr != nil {
		lineNum := tmplutil.ExtractLineNumber(parseErr)
		hint := ""
		if lineNum > 0 {
			hint = tmplutil.SourceContext(src, lineNum, 2, "> ")
		}
		result.AddErrorWithHint(parseErr.Error(), relPath, lineNum, hint)
		return
	}

	// Structural check: mismatched delimiters. Only warn when the parser
	// succeeded — if the parser found them balanced, the count mismatch
	// is likely harmless (e.g. {{ in string literals or comments).
	if opens, closes := strings.Count(src, "{{"), strings.Count(src, "}}"); opens != closes {
		result.AddWarning(
			fmt.Sprintf("mismatched delimiters: %d opening {{ vs %d closing }}", opens, closes),
			relPath, 0,
		)
	}
}
