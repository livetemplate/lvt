package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/livetemplate/lvt/internal/tmplutil"
)

// ValidateTemplate parses a generated .tmpl file and returns a clear error
// if the template contains syntax errors. This catches issues at generation
// time rather than at runtime.
func ValidateTemplate(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", path, err)
	}

	_, err = template.New(filepath.Base(path)).Parse(string(content))
	if err != nil {
		return formatTemplateError(path, string(content), err)
	}
	return nil
}

// formatTemplateError enhances a template parse error with file path,
// line number, and surrounding source context.
func formatTemplateError(path, content string, parseErr error) error {
	lineNum := tmplutil.ExtractLineNumber(parseErr)
	if lineNum <= 0 {
		return fmt.Errorf("template syntax error in %s: %w", path, parseErr)
	}

	ctx := tmplutil.SourceContext(content, lineNum, 2, "→ ")
	return fmt.Errorf("template syntax error in %s (line %d):\n%s\n  error: %w", path, lineNum, ctx, parseErr)
}
