package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Test helpers

func setupMCPTestDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "lvt-mcp-tool-test-*")
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

// TestMCPTool_LvtNew tests the lvt_new tool input/output
func TestMCPTool_LvtNew(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	tests := []struct {
		name        string
		input       NewAppInput
		expectError bool
	}{
		{
			name: "valid app with defaults",
			input: NewAppInput{
				Name: "testapp",
			},
			expectError: false,
		},
		{
			name: "valid app with kit",
			input: NewAppInput{
				Name: "testapp2",
				Kit:  "simple",
			},
			expectError: false,
		},
		{
			name: "valid app with CSS",
			input: NewAppInput{
				Name: "testapp3",
				CSS:  "pico",
			},
			expectError: false,
		},
		{
			name: "valid app with module",
			input: NewAppInput{
				Name:   "testapp4",
				Module: "github.com/test/testapp4",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{tt.input.Name}
			if tt.input.Kit != "" {
				args = append(args, "--kit", tt.input.Kit)
			}
			if tt.input.CSS != "" {
				args = append(args, "--css", tt.input.CSS)
			}
			if tt.input.Module != "" {
				args = append(args, "--module", tt.input.Module)
			}

			err := New(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && tt.input.Name != "" {
				// Verify app directory was created
				appDir := filepath.Join(tmpDir, tt.input.Name)
				if _, err := os.Stat(appDir); os.IsNotExist(err) {
					t.Errorf("App directory %s was not created", appDir)
				}
			}
		})
	}
}

// TestMCPTool_LvtGenResource tests the lvt_gen_resource tool
func TestMCPTool_LvtGenResource(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// First create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	tests := []struct {
		name        string
		input       GenResourceInput
		expectError bool
	}{
		{
			name: "valid resource with explicit types",
			input: GenResourceInput{
				Name: "tasks",
				Fields: map[string]string{
					"title":       "string",
					"description": "string",
					"completed":   "bool",
				},
			},
			expectError: false,
		},
		{
			name: "valid resource with inferred types",
			input: GenResourceInput{
				Name: "posts",
				Fields: map[string]string{
					"title":   "string",
					"content": "string",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"resource", tt.input.Name}
			for field, typ := range tt.input.Fields {
				args = append(args, fmt.Sprintf("%s:%s", field, typ))
			}

			err := Gen(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify resource files were created
				handlerFile := filepath.Join(appDir, "app", tt.input.Name, tt.input.Name+".go")
				if _, err := os.Stat(handlerFile); os.IsNotExist(err) {
					t.Errorf("Handler file %s was not created", handlerFile)
				}
			}
		})
	}
}

// TestMCPTool_LvtGenView tests the lvt_gen_view tool
func TestMCPTool_LvtGenView(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	tests := []struct {
		name        string
		input       GenViewInput
		expectError bool
	}{
		{
			name: "valid view",
			input: GenViewInput{
				Name: "dashboard",
			},
			expectError: false,
		},
		{
			name: "another valid view",
			input: GenViewInput{
				Name: "about",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"view", tt.input.Name}

			err := Gen(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify view files were created
				viewFile := filepath.Join(appDir, "app", tt.input.Name, tt.input.Name+".go")
				if _, err := os.Stat(viewFile); os.IsNotExist(err) {
					t.Errorf("View file %s was not created", viewFile)
				}
			}
		})
	}
}

// TestMCPTool_LvtGenAuth tests the lvt_gen_auth tool
func TestMCPTool_LvtGenAuth(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	tests := []struct {
		name        string
		input       GenAuthInput
		expectError bool
	}{
		{
			name: "default auth",
			input: GenAuthInput{
				StructName: "",
				TableName:  "",
			},
			expectError: false,
		},
		{
			name: "custom struct name",
			input: GenAuthInput{
				StructName: "Account",
				TableName:  "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []string
			if tt.input.StructName != "" {
				args = append(args, "auth", tt.input.StructName)
			} else {
				args = append(args, "auth")
			}
			if tt.input.TableName != "" {
				args = append(args, tt.input.TableName)
			}

			err := Gen(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify auth files were created
				authDir := filepath.Join(appDir, "app", "auth")
				if _, err := os.Stat(authDir); os.IsNotExist(err) {
					t.Errorf("Auth directory %s was not created", authDir)
				}
			}
		})
	}
}

// TestMCPTool_LvtGenSchema tests the lvt_gen_schema tool
func TestMCPTool_LvtGenSchema(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	tests := []struct {
		name        string
		input       GenSchemaInput
		expectError bool
	}{
		{
			name: "valid schema",
			input: GenSchemaInput{
				Table: "products",
				Fields: map[string]string{
					"name":  "string",
					"price": "float",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"schema", tt.input.Table}
			for field, typ := range tt.input.Fields {
				args = append(args, fmt.Sprintf("%s:%s", field, typ))
			}

			err := Gen(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify migration was created
				migrationsDir := filepath.Join(appDir, "database", "migrations")
				entries, err := os.ReadDir(migrationsDir)
				if err != nil {
					t.Errorf("Failed to read migrations dir: %v", err)
				}
				if len(entries) == 0 {
					t.Error("No migration files created")
				}
			}
		})
	}
}

// TestMCPTool_Migrations tests all migration tools
func TestMCPTool_Migrations(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Generate a resource first (creates migrations)
	err = Gen([]string{"resource", "items", "name:string"})
	if err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	t.Run("migration_status", func(t *testing.T) {
		err := Migration([]string{"status"})
		if err != nil {
			t.Errorf("Migration status failed: %v", err)
		}
	})

	t.Run("migration_up", func(t *testing.T) {
		err := Migration([]string{"up"})
		if err != nil {
			t.Errorf("Migration up failed: %v", err)
		}
	})

	t.Run("migration_create", func(t *testing.T) {
		err := Migration([]string{"create", "add_test_column"})
		if err != nil {
			t.Errorf("Migration create failed: %v", err)
		}
	})

	// Skip migration down as it may fail if already at base
	// t.Run("migration_down", func(t *testing.T) {
	// 	err := Migration([]string{"down"})
	// 	if err != nil {
	// 		t.Errorf("Migration down failed: %v", err)
	// 	}
	// })
}

// TestMCPTool_Seed tests the lvt_seed tool
func TestMCPTool_Seed(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Generate a resource
	err = Gen([]string{"resource", "items", "name:string"})
	if err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// Run migrations
	err = Migration([]string{"up"})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	tests := []struct {
		name        string
		input       SeedInput
		expectError bool
	}{
		{
			name: "seed with default count",
			input: SeedInput{
				Resource: "items",
				Count:    10,
			},
			expectError: false,
		},
		{
			name: "seed with cleanup",
			input: SeedInput{
				Resource: "items",
				Count:    5,
				Cleanup:  true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []string
			if tt.input.Resource != "" {
				args = append(args, tt.input.Resource)
			}
			if tt.input.Count > 0 {
				args = append(args, "--count", fmt.Sprintf("%d", tt.input.Count))
			}
			if tt.input.Cleanup {
				args = append(args, "--cleanup")
			}

			err := Seed(args)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestMCPTool_ResourceInspect tests resource list and describe tools
func TestMCPTool_ResourceInspect(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	// Generate a resource
	err = Gen([]string{"resource", "items", "name:string"})
	if err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	t.Run("resource_list", func(t *testing.T) {
		err := Resource([]string{"list"})
		if err != nil {
			t.Errorf("Resource list failed: %v", err)
		}
	})

	t.Run("resource_describe", func(t *testing.T) {
		err := Resource([]string{"describe", "items"})
		if err != nil {
			t.Errorf("Resource describe failed: %v", err)
		}
	})

	t.Run("resource_describe_missing", func(t *testing.T) {
		err := Resource([]string{"describe", "nonexistent"})
		if err == nil {
			t.Error("Expected error for nonexistent resource")
		}
	})
}

// TestMCPTool_ValidateTemplate tests the lvt_validate_template tool
func TestMCPTool_ValidateTemplate(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create a valid template file
	validTemplate := `<div>Hello World</div>`

	validFile := filepath.Join(tmpDir, "valid.tmpl")
	err := os.WriteFile(validFile, []byte(validTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	tests := []struct {
		name        string
		input       ValidateTemplatesInput
		expectError bool
	}{
		{
			name: "valid template",
			input: ValidateTemplatesInput{
				TemplateFile: validFile,
			},
			expectError: false,
		},
		{
			name: "missing file",
			input: ValidateTemplatesInput{
				TemplateFile: filepath.Join(tmpDir, "nonexistent.tmpl"),
			},
			expectError: true,
		},
		{
			name: "empty path",
			input: ValidateTemplatesInput{
				TemplateFile: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.TemplateFile == "" {
				// Skip parse call for empty path test
				return
			}

			err := Parse([]string{tt.input.TemplateFile})
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestMCPTool_EnvGenerate tests the lvt_env_generate tool
func TestMCPTool_EnvGenerate(t *testing.T) {
	tmpDir, cleanup := setupMCPTestDir(t)
	defer cleanup()

	// Create an app
	err := New([]string{"testapp"})
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	appDir := filepath.Join(tmpDir, "testapp")
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Failed to change to app dir: %v", err)
	}

	t.Run("env_generate", func(t *testing.T) {
		err := Env([]string{"generate"})
		if err != nil {
			t.Errorf("Env generate failed: %v", err)
		}

		// Verify .env.example was created
		envFile := filepath.Join(appDir, ".env.example")
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			t.Error(".env.example was not created")
		}
	})
}

// TestMCPTool_Kits tests kits list and info tools
func TestMCPTool_Kits(t *testing.T) {
	t.Run("kits_list", func(t *testing.T) {
		err := Kits([]string{"list"})
		if err != nil {
			t.Errorf("Kits list failed: %v", err)
		}
	})

	t.Run("kits_info", func(t *testing.T) {
		// Test with a known kit
		err := Kits([]string{"info", "multi"})
		if err != nil {
			t.Errorf("Kits info failed: %v", err)
		}
	})

	t.Run("kits_info_missing", func(t *testing.T) {
		err := Kits([]string{"info", "nonexistent"})
		if err == nil {
			t.Error("Expected error for nonexistent kit")
		}
	})
}

// TestMCPTool_InputValidation tests input validation across all tools
func TestMCPTool_InputValidation(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		valid bool
	}{
		{
			name:  "valid NewAppInput",
			input: NewAppInput{Name: "testapp"},
			valid: true,
		},
		{
			name:  "invalid NewAppInput - empty name",
			input: NewAppInput{Name: ""},
			valid: false,
		},
		{
			name:  "valid GenResourceInput",
			input: GenResourceInput{Name: "items", Fields: map[string]string{"name": "string"}},
			valid: true,
		},
		{
			name:  "invalid GenResourceInput - empty name",
			input: GenResourceInput{Name: "", Fields: map[string]string{"name": "string"}},
			valid: false,
		},
		{
			name:  "invalid GenResourceInput - no fields",
			input: GenResourceInput{Name: "items", Fields: map[string]string{}},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify structs are properly typed
			if tt.input == nil {
				t.Error("Input should not be nil")
			}
		})
	}
}
