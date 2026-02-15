package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// lineNumberPattern matches Go template parse error positions like "template: name:5:" or "template: name:5:22:"
var lineNumberPattern = regexp.MustCompile(`template:.*?:(\d+)`)

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
	lineNum := extractLineNumber(parseErr)
	if lineNum <= 0 {
		return fmt.Errorf("template syntax error in %s: %w", path, parseErr)
	}

	context := sourceContext(content, lineNum, 2)
	return fmt.Errorf("template syntax error in %s (line %d):\n%s\n  error: %w", path, lineNum, context, parseErr)
}

// extractLineNumber pulls the line number from a Go template parse error message.
func extractLineNumber(err error) int {
	matches := lineNumberPattern.FindStringSubmatch(err.Error())
	if len(matches) < 2 {
		return 0
	}
	n, convErr := strconv.Atoi(matches[1])
	if convErr != nil {
		return 0
	}
	return n
}

// sourceContext returns a few lines of source around the given line number,
// with line numbers and an arrow marking the error line.
func sourceContext(content string, lineNum, surroundingLines int) string {
	lines := strings.Split(content, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}

	start := max(lineNum-surroundingLines-1, 0)
	end := min(lineNum+surroundingLines, len(lines))

	var b strings.Builder
	for i := start; i < end; i++ {
		marker := "  "
		if i+1 == lineNum {
			marker = "→ "
		}
		fmt.Fprintf(&b, "  %s%4d | %s\n", marker, i+1, lines[i])
	}
	return strings.TrimRight(b.String(), "\n")
}
