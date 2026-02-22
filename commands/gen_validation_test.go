package commands

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestGenResource_WithValidation verifies that structural validation (go.mod,
// templates, migrations) runs after generation in a fully valid app.
func TestGenResource_WithValidation(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
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
	tmpDir, cleanup := setupMCPTestDir(t)
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
	tmpDir, cleanup := setupMCPTestDir(t)
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

// TestValidationOutput_JSONMarshal verifies ValidationOutput serializes correctly.
func TestValidationOutput_JSONMarshal(t *testing.T) {
	vo := &ValidationOutput{
		Valid:      false,
		ErrorCount: 2,
		WarnCount:  1,
		Issues: []ValidationIssueOutput{
			{Level: "error", File: "main.go", Line: 10, Message: "undefined: foo"},
			{Level: "error", File: "main.go", Line: 20, Message: "undefined: bar"},
			{Level: "warning", Message: "no templates found"},
		},
	}

	data, err := json.Marshal(vo)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationOutput: %v", err)
	}

	// Verify it round-trips
	var decoded ValidationOutput
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ValidationOutput: %v", err)
	}

	if decoded.Valid {
		t.Error("expected Valid=false")
	}
	if decoded.ErrorCount != 2 {
		t.Errorf("expected ErrorCount=2, got %d", decoded.ErrorCount)
	}
	if decoded.WarnCount != 1 {
		t.Errorf("expected WarnCount=1, got %d", decoded.WarnCount)
	}
	if len(decoded.Issues) != 3 {
		t.Errorf("expected 3 issues, got %d", len(decoded.Issues))
	}
}

// TestMCPGenResource_IncludesValidation verifies MCP gen_resource response has validation.
func TestMCPGenResource_IncludesValidation(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
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

	// Simulate MCP handler: call Gen with --skip-validation, then run validation
	err = Gen([]string{"resource", "items", "name:string", "--skip-validation"})
	if err != nil {
		t.Fatalf("Gen resource failed: %v", err)
	}

	// Run MCP validation
	result := runMCPValidation(context.Background(), appDir)
	if result == nil {
		t.Fatal("expected non-nil validation result")
	}

	// In a valid app, validation should pass
	if !result.Valid {
		t.Errorf("expected valid=true, got issues: %+v", result.Issues)
	}
	if result.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", result.ErrorCount)
	}

	// Verify it can be placed into GenResourceOutput.
	// Success reflects generation success (always true here), not validation state.
	output := GenResourceOutput{
		Success:    true,
		Message:    "test",
		Validation: result,
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal GenResourceOutput with validation: %v", err)
	}

	// Verify "validation" key exists in JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if _, ok := raw["validation"]; !ok {
		t.Error("expected 'validation' key in GenResourceOutput JSON")
	}
}
