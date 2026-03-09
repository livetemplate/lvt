package commands

import (
	"os"
	"path/filepath"
	"testing"
)

// setupGenTestDir creates a temp directory, changes to it, and returns a cleanup function.
func setupGenTestDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "lvt-gen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	cleanup := func() {
		os.Chdir(oldDir)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// TestGenResource_WithValidation verifies that structural validation (go.mod,
// templates, migrations) runs after generation in a fully valid app.
func TestGenResource_WithValidation(t *testing.T) {
	tmpDir, cleanup := setupGenTestDir(t)
	defer cleanup()

	// Create a full app so structural validation passes
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Gen resource without --skip-validation — validation should run and pass
	err = Gen([]string{"resource", "tasks", "title:string", "done:bool"})
	if err != nil {
		t.Errorf("GenResource with validation failed: %v", err)
	}
}

// TestGenResource_SkipValidation verifies that --skip-validation prevents validation.
func TestGenResource_SkipValidation(t *testing.T) {
	tmpDir, cleanup := setupGenTestDir(t)
	defer cleanup()

	// Create a full app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Gen resource with --skip-validation
	err = Gen([]string{"resource", "tasks", "title:string", "--skip-validation"})
	if err != nil {
		t.Errorf("GenResource with --skip-validation failed: %v", err)
	}

	// Verify files were still created
	handlerFile := filepath.Join(appDir, "app", "tasks", "tasks.go")
	if _, err := os.Stat(handlerFile); os.IsNotExist(err) {
		t.Error("Handler file was not created")
	}
}

// TestGenView_WithValidation verifies that structural validation runs after view generation.
func TestGenView_WithValidation(t *testing.T) {
	tmpDir, cleanup := setupGenTestDir(t)
	defer cleanup()

	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Gen view without --skip-validation — validation should pass
	err = Gen([]string{"view", "dashboard"})
	if err != nil {
		t.Errorf("GenView with validation failed: %v", err)
	}
}

// TestGenSchema_WithValidation verifies that structural validation runs after schema generation.
func TestGenSchema_WithValidation(t *testing.T) {
	tmpDir, cleanup := setupGenTestDir(t)
	defer cleanup()

	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Gen schema without --skip-validation — validation should pass
	err = Gen([]string{"schema", "orders", "total:float", "status:string"})
	if err != nil {
		t.Errorf("GenSchema with validation failed: %v", err)
	}
}
