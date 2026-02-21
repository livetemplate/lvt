package validation

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/livetemplate/lvt/internal/validator"
)

// tmplLinePattern matches template parse errors like "template: name:5:" or "template: name:5:22:".
var tmplLinePattern = regexp.MustCompile(`template:.*?:(\d+)`)

// TemplateCheck validates all .tmpl files in an app directory.
type TemplateCheck struct{}

func (c *TemplateCheck) Name() string { return "templates" }

func (c *TemplateCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	var found bool

	walkErr := filepath.WalkDir(appPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
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
		lineNum := extractLineNumber(parseErr)
		hint := ""
		if lineNum > 0 {
			hint = sourceContext(src, lineNum, 2)
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

// extractLineNumber pulls the line number from a template parse error.
func extractLineNumber(err error) int {
	m := tmplLinePattern.FindStringSubmatch(err.Error())
	if len(m) < 2 {
		return 0
	}
	n, convErr := strconv.Atoi(m[1])
	if convErr != nil {
		return 0
	}
	return n
}

// sourceContext returns lines around lineNum with an arrow on the target line.
func sourceContext(content string, lineNum, surrounding int) string {
	lines := strings.Split(content, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}

	start := max(lineNum-surrounding-1, 0)
	end := min(lineNum+surrounding, len(lines))

	var b strings.Builder
	for i := start; i < end; i++ {
		marker := "  "
		if i+1 == lineNum {
			marker = "> "
		}
		fmt.Fprintf(&b, "  %s%4d | %s\n", marker, i+1, lines[i])
	}
	return strings.TrimRight(b.String(), "\n")
}
