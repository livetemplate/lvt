package base

import (
	"embed"
	"html/template"
	"testing"
)

//go:embed testdata/*.tmpl
var testTemplateFS embed.FS

func TestNewTemplateSet(t *testing.T) {
	ts := NewTemplateSet(testTemplateFS, "testdata/*.tmpl", "test")

	if ts.Pattern != "testdata/*.tmpl" {
		t.Errorf("expected Pattern 'testdata/*.tmpl', got '%s'", ts.Pattern)
	}

	if ts.Namespace != "test" {
		t.Errorf("expected Namespace 'test', got '%s'", ts.Namespace)
	}

	if ts.Funcs != nil {
		t.Error("expected Funcs to be nil by default")
	}
}

func TestTemplateSet_WithFuncs(t *testing.T) {
	ts := NewTemplateSet(testTemplateFS, "testdata/*.tmpl", "test")

	funcs := template.FuncMap{
		"testFunc": func() string { return "test" },
	}

	ts2 := ts.WithFuncs(funcs)

	// Original should be unchanged
	if ts.Funcs != nil {
		t.Error("WithFuncs should not modify the original")
	}

	// New set should have funcs
	if ts2.Funcs == nil {
		t.Error("expected Funcs to be set")
	}

	// Should have the function
	if _, ok := ts2.Funcs["testFunc"]; !ok {
		t.Error("expected 'testFunc' to be in Funcs")
	}

	// Other fields should be preserved
	if ts2.Pattern != ts.Pattern {
		t.Error("WithFuncs should preserve Pattern")
	}

	if ts2.Namespace != ts.Namespace {
		t.Error("WithFuncs should preserve Namespace")
	}
}

func TestTemplateSet_FSReadable(t *testing.T) {
	ts := NewTemplateSet(testTemplateFS, "testdata/*.tmpl", "test")

	// Should be able to read from FS
	entries, err := ts.FS.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata dir: %v", err)
	}

	// Should have at least one template file
	if len(entries) == 0 {
		t.Error("expected at least one file in testdata")
	}
}
