package generator

import (
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

