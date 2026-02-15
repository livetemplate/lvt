package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateTemplate_ValidTemplate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.tmpl")

	content := `{{define "layout"}}
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
  {{block "content" .}}{{end}}
</body>
</html>
{{end}}`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateTemplate(path); err != nil {
		t.Errorf("expected valid template, got error: %v", err)
	}
}

func TestValidateTemplate_InvalidSyntax(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.tmpl")

	content := `{{define "layout"}}
<html>
<body>
  {{if .Show}
  <p>Hello</p>
  {{end}}
</body>
</html>
{{end}}`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTemplate(path)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}

	errMsg := err.Error()

	// Should include file path
	if !strings.Contains(errMsg, path) {
		t.Errorf("error should include file path %q, got: %s", path, errMsg)
	}

	// Should include line number
	if !strings.Contains(errMsg, "line") {
		t.Errorf("error should include line number, got: %s", errMsg)
	}

	// Should include source context with arrow marker
	if !strings.Contains(errMsg, "→") {
		t.Errorf("error should include source context arrow, got: %s", errMsg)
	}
}

func TestValidateTemplate_FileNotFound(t *testing.T) {
	err := ValidateTemplate("/nonexistent/path/template.tmpl")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read template") {
		t.Errorf("expected read error, got: %v", err)
	}
}

func TestValidateTemplate_EmptyTemplate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.tmpl")

	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateTemplate(path); err != nil {
		t.Errorf("empty template should be valid, got error: %v", err)
	}
}

func TestValidateTemplate_UnclosedAction(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "unclosed.tmpl")

	content := `<html>
<body>
  {{range .Items
</body>
</html>`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTemplate(path)
	if err == nil {
		t.Fatal("expected error for unclosed action")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error should include file path, got: %v", err)
	}
}

func TestExtractLineNumber(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected int
	}{
		{"standard error", "template: test.tmpl:5: unexpected", 5},
		{"with column", "template: test.tmpl:12:22: function not defined", 12},
		{"no line number", "some other error", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLineNumber(fmt.Errorf("%s", tt.errMsg))
			if got != tt.expected {
				t.Errorf("extractLineNumber(%q) = %d, want %d", tt.errMsg, got, tt.expected)
			}
		})
	}
}

func TestSourceContext(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7"

	ctx := sourceContext(content, 4, 2)

	// Should include lines 2-6 (4 ± 2)
	if !strings.Contains(ctx, "line2") {
		t.Error("context should include line 2")
	}
	if !strings.Contains(ctx, "line6") {
		t.Error("context should include line 6")
	}
	// Line 4 should be marked with arrow
	if !strings.Contains(ctx, "→") {
		t.Error("context should mark error line with →")
	}
	// Should show line numbers
	if !strings.Contains(ctx, "4 |") {
		t.Error("context should show line numbers")
	}
}

func TestSourceContext_FirstLine(t *testing.T) {
	content := "line1\nline2\nline3"

	ctx := sourceContext(content, 1, 2)

	// Should include lines 1-3
	if !strings.Contains(ctx, "line1") {
		t.Error("context should include line 1")
	}
	if !strings.Contains(ctx, "→") {
		t.Error("context should mark error line")
	}
}

func TestSourceContext_OutOfRange(t *testing.T) {
	content := "line1\nline2"

	ctx := sourceContext(content, 10, 2)
	if ctx != "" {
		t.Errorf("out of range should return empty, got: %s", ctx)
	}

	ctx = sourceContext(content, 0, 2)
	if ctx != "" {
		t.Errorf("zero line should return empty, got: %s", ctx)
	}
}
