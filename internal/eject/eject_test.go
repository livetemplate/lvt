package eject

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAvailableComponents(t *testing.T) {
	components := AvailableComponents()
	if len(components) == 0 {
		t.Error("AvailableComponents() returned empty list")
	}

	// Verify each component has required fields
	for _, c := range components {
		if c.Name == "" {
			t.Error("Component with empty name")
		}
		if c.Package == "" {
			t.Errorf("Component %s has empty package", c.Name)
		}
		if c.Description == "" {
			t.Errorf("Component %s has empty description", c.Name)
		}
		if len(c.Templates) == 0 {
			t.Errorf("Component %s has no templates", c.Name)
		}
	}
}

func TestFindComponent(t *testing.T) {
	tests := []struct {
		name      string
		component string
		wantFound bool
	}{
		{"existing component", "dropdown", true},
		{"another existing", "modal", true},
		{"non-existent", "nonexistent", false},
		{"empty name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindComponent(tt.component)
			if tt.wantFound && result == nil {
				t.Errorf("FindComponent(%q) returned nil, expected component", tt.component)
			}
			if !tt.wantFound && result != nil {
				t.Errorf("FindComponent(%q) returned component, expected nil", tt.component)
			}
			if tt.wantFound && result != nil && result.Name != tt.component {
				t.Errorf("FindComponent(%q) returned %q", tt.component, result.Name)
			}
		})
	}
}

func TestEjectOptions_DefaultDestDir(t *testing.T) {
	// Test that EjectComponent uses default destination
	opts := EjectOptions{
		ComponentName: "nonexistent",
		DestDir:       "", // empty means use default
	}

	err := EjectComponent(opts)
	if err == nil {
		t.Error("Expected error for nonexistent component")
	}

	// Error should mention the component name
	if err != nil && !containsString(err.Error(), "nonexistent") {
		t.Errorf("Error should mention component name, got: %s", err.Error())
	}
}

func TestEjectComponent_UnknownComponent(t *testing.T) {
	opts := EjectOptions{
		ComponentName: "unknown-component",
	}

	err := EjectComponent(opts)
	if err == nil {
		t.Error("Expected error for unknown component")
	}

	errMsg := err.Error()
	if !containsString(errMsg, "unknown component") {
		t.Errorf("Error should mention 'unknown component', got: %s", errMsg)
	}
}

func TestEjectComponent_DestinationExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-eject-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create destination that already exists
	destDir := filepath.Join(tmpDir, "internal", "components", "dropdown")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}

	// Try to eject without force
	opts := EjectOptions{
		ComponentName: "dropdown",
		DestDir:       destDir,
		Force:         false,
	}

	err = EjectComponent(opts)
	if err == nil {
		t.Error("Expected error when destination exists and force=false")
	}

	if err != nil && !containsString(err.Error(), "already exists") {
		t.Errorf("Error should mention 'already exists', got: %s", err.Error())
	}
}

func TestEjectTemplateOptions_DefaultDestDir(t *testing.T) {
	// Test that EjectTemplate uses default destination
	opts := EjectTemplateOptions{
		ComponentName: "nonexistent",
		TemplateName:  "default",
		DestDir:       "",
	}

	err := EjectTemplate(opts)
	if err == nil {
		t.Error("Expected error for nonexistent component")
	}
}

func TestEjectTemplate_UnknownComponent(t *testing.T) {
	opts := EjectTemplateOptions{
		ComponentName: "unknown-component",
		TemplateName:  "default",
	}

	err := EjectTemplate(opts)
	if err == nil {
		t.Error("Expected error for unknown component")
	}

	errMsg := err.Error()
	if !containsString(errMsg, "unknown component") {
		t.Errorf("Error should mention 'unknown component', got: %s", errMsg)
	}
}

func TestEjectTemplate_UnknownTemplate(t *testing.T) {
	opts := EjectTemplateOptions{
		ComponentName: "dropdown",
		TemplateName:  "unknown-template",
	}

	err := EjectTemplate(opts)
	if err == nil {
		t.Error("Expected error for unknown template")
	}

	errMsg := err.Error()
	if !containsString(errMsg, "unknown template") {
		t.Errorf("Error should mention 'unknown template', got: %s", errMsg)
	}
}

func TestEjectTemplate_DestinationExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-eject-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create destination file that already exists
	destDir := filepath.Join(tmpDir, "internal", "templates")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}

	destFile := filepath.Join(destDir, "dropdown-default.tmpl")
	if err := os.WriteFile(destFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create dest file: %v", err)
	}

	// Try to eject without force
	opts := EjectTemplateOptions{
		ComponentName: "dropdown",
		TemplateName:  "default",
		DestDir:       destDir,
		Force:         false,
	}

	err = EjectTemplate(opts)
	if err == nil {
		t.Error("Expected error when destination exists and force=false")
	}

	if err != nil && !containsString(err.Error(), "already exists") {
		t.Errorf("Error should mention 'already exists', got: %s", err.Error())
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := "test content\nwith multiple lines"
	if err := os.WriteFile(srcFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	dstFile := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content
	copied, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copied) != content {
		t.Errorf("Copied content mismatch: got %q, want %q", string(copied), content)
	}
}

func TestCopyDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-copydir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory structure
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(filepath.Join(srcDir, "templates"), 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}

	// Create some files
	files := map[string]string{
		"main.go":                "package main",
		"options.go":             "package main\nvar opts = 1",
		"templates/default.tmpl": "{{define \"test\"}}hello{{end}}",
	}

	for name, content := range files {
		path := filepath.Join(srcDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create %s: %v", name, err)
		}
	}

	// Also create a test file that should be skipped
	testFile := filepath.Join(srcDir, "main_test.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Copy directory
	dstDir := filepath.Join(tmpDir, "dst")
	copied, err := copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Verify files were copied (test file should be skipped)
	if len(copied) != len(files) {
		t.Errorf("Expected %d copied files, got %d", len(files), len(copied))
	}

	// Verify test file was not copied
	testDst := filepath.Join(dstDir, "main_test.go")
	if _, err := os.Stat(testDst); !os.IsNotExist(err) {
		t.Error("Test file should not have been copied")
	}

	// Verify regular files were copied
	for name, expectedContent := range files {
		path := filepath.Join(dstDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read copied file %s: %v", name, err)
			continue
		}
		if string(content) != expectedContent {
			t.Errorf("Content mismatch for %s: got %q, want %q", name, string(content), expectedContent)
		}
	}
}

func TestUpdateImports(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-imports-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a Go file with imports
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package test

import (
	"github.com/livetemplate/components/dropdown"
	"github.com/livetemplate/components/dropdown/internal"
)

func main() {
	dropdown.New("test")
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	// Update imports
	oldPkg := "github.com/livetemplate/components/dropdown"
	newPkg := "mymodule/internal/components/dropdown"
	if err := updateImports(goFile, oldPkg, newPkg); err != nil {
		t.Fatalf("updateImports failed: %v", err)
	}

	// Verify imports were updated
	updated, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	if containsString(string(updated), oldPkg) {
		t.Error("Old package import still present")
	}
	if !containsString(string(updated), newPkg) {
		t.Error("New package import not found")
	}
}

func TestComponentInfo_Templates(t *testing.T) {
	// Verify known components have expected templates
	testCases := []struct {
		component string
		templates []string
	}{
		{"dropdown", []string{"default", "searchable", "multi"}},
		{"modal", []string{"default", "confirm", "sheet"}},
		{"tabs", []string{"horizontal", "vertical", "pills"}},
	}

	for _, tc := range testCases {
		t.Run(tc.component, func(t *testing.T) {
			comp := FindComponent(tc.component)
			if comp == nil {
				t.Fatalf("Component %s not found", tc.component)
			}

			for _, expected := range tc.templates {
				found := false
				for _, actual := range comp.Templates {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Template %q not found in %s templates: %v", expected, tc.component, comp.Templates)
				}
			}
		})
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
