package evolution

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livetemplate/lvt/internal/validator"
)

func TestTestFix_ValidFix(t *testing.T) {
	// Create a source directory with a file to fix
	srcDir := t.TempDir()
	handlerDir := filepath.Join(srcDir, "app", "posts")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := `package posts

type State struct {
	IsAdding bool
	EditingID string
}
`
	if err := os.WriteFile(filepath.Join(handlerDir, "handler.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a passing validator
	passValidator := func(ctx context.Context, appPath string) *validator.ValidationResult {
		return validator.NewValidationResult()
	}

	tester := NewTester(passValidator)
	fix := Fix{
		ID:          "fix-1",
		TargetFile:  "*/handler.go",
		FindPattern: "IsAdding bool",
		Replace:     `IsAdding bool ` + "`" + `lvt:"transient"` + "`",
		IsRegex:     false,
	}

	result, err := tester.TestFix(context.Background(), fix, srcDir)
	if err != nil {
		t.Fatalf("test fix: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got error: %s", result.Error)
	}
	if result.Validation == nil {
		t.Error("expected validation result")
	}
}

func TestTestFix_BadFix(t *testing.T) {
	srcDir := t.TempDir()
	handlerDir := filepath.Join(srcDir, "app", "posts")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := `package posts

type State struct {
	IsAdding bool
}
`
	if err := os.WriteFile(filepath.Join(handlerDir, "handler.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a failing validator
	failValidator := func(ctx context.Context, appPath string) *validator.ValidationResult {
		r := validator.NewValidationResult()
		r.AddError("compilation failed", "handler.go", 1)
		return r
	}

	tester := NewTester(failValidator)
	fix := Fix{
		ID:          "fix-2",
		TargetFile:  "*/handler.go",
		FindPattern: "IsAdding bool",
		Replace:     "IsAdding BROKEN_TYPE", // introduces error
		IsRegex:     false,
	}

	result, err := tester.TestFix(context.Background(), fix, srcDir)
	if err != nil {
		t.Fatalf("test fix: %v", err)
	}
	if result.Success {
		t.Error("expected failure for bad fix")
	}
}

func TestTestFix_NoMatch(t *testing.T) {
	srcDir := t.TempDir()
	handlerDir := filepath.Join(srcDir, "app", "posts")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(handlerDir, "handler.go"), []byte("package posts\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tester := NewTester(nil)
	fix := Fix{
		ID:          "fix-3",
		TargetFile:  "*/handler.go",
		FindPattern: "NON_EXISTENT_TEXT",
		Replace:     "replacement",
		IsRegex:     false,
	}

	result, err := tester.TestFix(context.Background(), fix, srcDir)
	if err != nil {
		t.Fatalf("test fix: %v", err)
	}
	if result.Success {
		t.Error("expected failure when fix target not found")
	}
	if result.Error == "" {
		t.Error("expected error message")
	}
}

func TestTestFix_Cleanup(t *testing.T) {
	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	tester := NewTester(nil)
	fix := Fix{
		ID:          "fix-cleanup",
		TargetFile:  "test.txt",
		FindPattern: "hello",
		Replace:     "world",
		IsRegex:     false,
	}

	_, err := tester.TestFix(context.Background(), fix, srcDir)
	if err != nil {
		t.Fatalf("test fix: %v", err)
	}

	// Verify temp dir was cleaned up by checking that source is untouched
	data, err := os.ReadFile(filepath.Join(srcDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("source should be untouched, got %q", string(data))
	}
}
