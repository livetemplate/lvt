package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidComponentName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		// Valid names
		{"simple name", "button", true},
		{"with hyphen", "date-picker", true},
		{"with numbers", "card2", true},
		{"multiple hyphens", "multi-select-dropdown", true},

		// Invalid names
		{"empty", "", false},
		{"starts with hyphen", "-button", false},
		{"ends with hyphen", "button-", false},
		{"uppercase", "Button", false},
		{"mixed case", "datePicker", false},
		{"underscore", "date_picker", false},
		{"space", "date picker", false},
		{"special chars", "button@1", false},
		{"double hyphen", "date--picker", true}, // This is technically valid but ugly
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidComponentName(tt.input)
			if result != tt.valid {
				t.Errorf("isValidComponentName(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"button", "Button"},
		{"date-picker", "DatePicker"},
		{"multi-select-dropdown", "MultiSelectDropdown"},
		{"a", "A"},
		{"a-b-c", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"button", "button"},
		{"date-picker", "datePicker"},
		{"multi-select-dropdown", "multiSelectDropdown"},
		{"a", "a"},
		{"a-b-c", "aBC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPackageName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"button", "button"},
		{"date-picker", "datepicker"},
		{"multi-select-dropdown", "multiselectdropdown"},
		{"a-b-c", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPackageName(tt.input)
			if result != tt.expected {
				t.Errorf("toPackageName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewComponent_Creates_Files(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-new-component-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run NewComponent
	err = NewComponent([]string{"test-widget"})
	if err != nil {
		t.Fatalf("NewComponent failed: %v", err)
	}

	// Verify directory structure
	expectedDirs := []string{
		"components/test-widget",
		"components/test-widget/templates",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tmpDir, dir)
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			t.Errorf("Expected directory missing: %s", dir)
		} else if !info.IsDir() {
			t.Errorf("Expected directory is not a directory: %s", dir)
		}
	}

	// Verify files
	expectedFiles := []string{
		"components/test-widget/test-widget.go",
		"components/test-widget/options.go",
		"components/test-widget/templates.go",
		"components/test-widget/test-widget_test.go",
		"components/test-widget/templates/default.tmpl",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file missing: %s", file)
		}
	}
}

func TestNewComponent_ValidGoCode(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-new-component-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run NewComponent with hyphenated name
	err = NewComponent([]string{"date-picker"})
	if err != nil {
		t.Fatalf("NewComponent failed: %v", err)
	}

	// Read the main file
	mainFile := filepath.Join(tmpDir, "components/date-picker/date-picker.go")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatalf("Failed to read main file: %v", err)
	}

	contentStr := string(content)

	// Verify package name doesn't have hyphen (Go syntax error)
	if strings.Contains(contentStr, "package date-picker") {
		t.Error("Package name contains hyphen which is invalid Go syntax")
	}

	// Verify package name is correct
	if !strings.Contains(contentStr, "package datepicker") {
		t.Error("Package name should be 'datepicker' (no hyphens)")
	}

	// Verify struct name is PascalCase
	if !strings.Contains(contentStr, "type DatePicker struct") {
		t.Error("Struct should be named DatePicker")
	}

	// Verify namespace uses original hyphenated name
	if !strings.Contains(contentStr, `"date-picker"`) {
		t.Error("Namespace should use original hyphenated name")
	}
}

func TestNewComponent_CustomDest(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-new-component-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run NewComponent with custom destination
	customDest := filepath.Join(tmpDir, "internal", "ui", "button")
	err = NewComponent([]string{"button", "--dest", customDest})
	if err != nil {
		t.Fatalf("NewComponent failed: %v", err)
	}

	// Verify files are in custom location
	expectedFiles := []string{
		"button.go",
		"options.go",
		"templates.go",
		"button_test.go",
		"templates/default.tmpl",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(customDest, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file missing in custom dest: %s", file)
		}
	}
}

func TestNewComponent_DestinationExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-new-component-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Create destination first
	destDir := filepath.Join(tmpDir, "components", "button")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}

	// Run NewComponent (should fail)
	err = NewComponent([]string{"button"})
	if err == nil {
		t.Error("Expected error when destination exists")
	}

	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Error should mention 'already exists', got: %s", err.Error())
	}
}

func TestNewComponent_InvalidName(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"empty args", []string{}, false}, // shows help, no error
		{"uppercase", []string{"Button"}, true},
		{"starts with hyphen", []string{"-button"}, true},
		{"ends with hyphen", []string{"button-"}, true},
		{"special chars", []string{"button@1"}, true},
		{"flag-like", []string{"--help"}, false}, // shows help, no error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for each test
			tmpDir, err := os.MkdirTemp("", "lvt-new-component-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get working dir: %v", err)
			}
			defer os.Chdir(oldDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change dir: %v", err)
			}

			err = NewComponent(tt.args)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for args %v", tt.args)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for args %v: %v", tt.args, err)
			}
		})
	}
}

func TestNewComponent_Help(t *testing.T) {
	// --help should not return an error
	err := NewComponent([]string{"--help"})
	if err != nil {
		t.Errorf("--help should not return error, got: %v", err)
	}

	// -h should not return an error
	err = NewComponent([]string{"-h"})
	if err != nil {
		t.Errorf("-h should not return error, got: %v", err)
	}
}

func TestComponentData(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPackage string
		wantPascal  string
		wantCamel   string
	}{
		{"simple", "button", "button", "Button", "button"},
		{"hyphenated", "date-picker", "datepicker", "DatePicker", "datePicker"},
		{"multi-hyphen", "multi-select-dropdown", "multiselectdropdown", "MultiSelectDropdown", "multiSelectDropdown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := componentData{
				Name:        tt.input,
				PackageName: toPackageName(tt.input),
				NamePascal:  toPascalCase(tt.input),
				NameCamel:   toCamelCase(tt.input),
			}

			if data.PackageName != tt.wantPackage {
				t.Errorf("PackageName = %q, want %q", data.PackageName, tt.wantPackage)
			}
			if data.NamePascal != tt.wantPascal {
				t.Errorf("NamePascal = %q, want %q", data.NamePascal, tt.wantPascal)
			}
			if data.NameCamel != tt.wantCamel {
				t.Errorf("NameCamel = %q, want %q", data.NameCamel, tt.wantCamel)
			}
		})
	}
}
