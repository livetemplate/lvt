package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/internal/generator"
	"github.com/livetemplate/lvt/internal/parser"
)

// runHandlerGoldenTest generates a resource and compares the handler output against a golden file.
func runHandlerGoldenTest(t *testing.T, resourceName string, fields []parser.Field, goldenPath, handlerSubpath string, withAuthz ...bool) {
	t.Helper()
	authz := len(withAuthz) > 0 && withAuthz[0]
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", resourceName, fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", authz); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	handlerPath := filepath.Join(tmpDir, handlerSubpath)
	generated, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read generated handler: %v", err)
	}

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata/golden", 0755); err != nil {
			t.Fatalf("Failed to create golden directory: %v", err)
		}
		if err := os.WriteFile(goldenPath, generated, 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Log("Updated golden file")
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("Golden file doesn't exist. Run with UPDATE_GOLDEN=1 to create it.")
		}
		t.Fatalf("Failed to read golden file: %v", err)
	}

	if string(golden) != string(generated) {
		t.Errorf("Generated code differs from golden file.\n"+
			"Run with UPDATE_GOLDEN=1 to update.\n"+
			"Golden: %d bytes, Generated: %d bytes",
			len(golden), len(generated))

		goldenLines := strings.Split(string(golden), "\n")
		genLines := strings.Split(string(generated), "\n")

		for i := 0; i < len(goldenLines) && i < len(genLines); i++ {
			if goldenLines[i] != genLines[i] {
				t.Logf("First difference at line %d:", i+1)
				t.Logf("  Golden:    %s", goldenLines[i])
				t.Logf("  Generated: %s", genLines[i])
				break
			}
		}
	}
}

// TestResourceHandlerGolden validates handler generation against golden file
func TestResourceHandlerGolden(t *testing.T) {
	fields := []parser.Field{
		{Name: "name", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "age", Type: "int", GoType: "int64", SQLType: "INTEGER", Metadata: parser.GetFieldMetadata("int")},
	}
	runHandlerGoldenTest(t, "User", fields,
		"testdata/golden/resource_handler.go.golden",
		"app/user/user.go")
}

// TestFileUploadResourceHandlerGolden validates handler generation for file/image fields against golden file
func TestFileUploadResourceHandlerGolden(t *testing.T) {
	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "photo", Type: "image", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: true, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
		{Name: "doc", Type: "file", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: false, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
	}
	runHandlerGoldenTest(t, "Gallery", fields,
		"testdata/golden/file_upload_handler.go.golden",
		"app/gallery/gallery.go")
}

// TestAuthzResourceHandlerGolden validates handler generation with --with-authz
func TestAuthzResourceHandlerGolden(t *testing.T) {
	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "content", Type: "text", GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: parser.GetFieldMetadata("text")},
	}
	runHandlerGoldenTest(t, "Post", fields,
		"testdata/golden/resource_handler_authz.go.golden",
		"app/post/post.go", true)
}

// TestResourceHandlerUnstyledImport verifies that styles="unstyled" generates the unstyled import
func TestResourceHandlerUnstyledImport(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "name", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Item", fields, "multi", "tailwind", "unstyled", "infinite", 20, "modal", "", false); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	handlerPath := filepath.Join(tmpDir, "app", "item", "item.go")
	generated, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read generated handler: %v", err)
	}

	content := string(generated)
	if !strings.Contains(content, `styles/unstyled"`) {
		t.Error("Expected styles/unstyled import in generated handler, not found")
	}
	if strings.Contains(content, `styles/tailwind"`) {
		t.Error("Expected no styles/tailwind import in generated handler for unstyled mode, but found one")
	}
}

// TestResourceHandlerInvalidStyles verifies that an invalid styles value is rejected
func TestResourceHandlerInvalidStyles(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "name", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
	}

	err := generator.GenerateResource(tmpDir, "testmodule", "Item", fields, "multi", "tailwind", "bootstrap", "infinite", 20, "modal", "", false)
	if err == nil {
		t.Fatal("Expected error for invalid styles adapter, got nil")
	}
	if !strings.Contains(err.Error(), "invalid styles adapter") {
		t.Errorf("Expected 'invalid styles adapter' error, got: %v", err)
	}
}

// TestViewHandlerGolden validates view generation against golden file
func TestViewHandlerGolden(t *testing.T) {
	tmpDir := t.TempDir()

	if err := generator.GenerateView(tmpDir, "testmodule", "Counter", "multi", "tailwind"); err != nil {
		t.Fatalf("Failed to generate view: %v", err)
	}

	handlerPath := filepath.Join(tmpDir, "app", "counter", "counter.go")
	generated, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read generated handler: %v", err)
	}

	goldenPath := "testdata/golden/view_handler.go.golden"

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata/golden", 0755); err != nil {
			t.Fatalf("Failed to create golden directory: %v", err)
		}
		if err := os.WriteFile(goldenPath, generated, 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Log("Updated golden file")
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("Golden file doesn't exist. Run with UPDATE_GOLDEN=1 to create it.")
		}
		t.Fatalf("Failed to read golden file: %v", err)
	}

	if string(golden) != string(generated) {
		t.Errorf("Generated code differs from golden file.\n"+
			"Run with UPDATE_GOLDEN=1 to update.\n"+
			"Golden: %d bytes, Generated: %d bytes",
			len(golden), len(generated))
	}
}

// TestResourceTemplateGolden validates template generation
func TestResourceTemplateGolden(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "published", Type: "bool", GoType: "bool", SQLType: "BOOLEAN", Metadata: parser.GetFieldMetadata("bool")},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Post", fields, "multi", "tailwind", "tailwind", "prev-next", 10, "modal", "", false); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	tmplPath := filepath.Join(tmpDir, "app", "post", "post.tmpl")
	generated, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read generated template: %v", err)
	}

	goldenPath := "testdata/golden/resource_template.tmpl.golden"

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata/golden", 0755); err != nil {
			t.Fatalf("Failed to create golden directory: %v", err)
		}
		if err := os.WriteFile(goldenPath, generated, 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Log("Updated golden file")
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("Golden file doesn't exist. Run with UPDATE_GOLDEN=1 to create it.")
		}
		t.Fatalf("Failed to read golden file: %v", err)
	}

	if string(golden) != string(generated) {
		t.Errorf("Generated template differs from golden file.\n" +
			"Run with UPDATE_GOLDEN=1 to update.")
	}
}

// TestAPIHandlerGolden validates API handler generation against golden file
func TestAPIHandlerGolden(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "content", Type: "text", GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: parser.GetFieldMetadata("text")},
		{Name: "published", Type: "bool", GoType: "bool", SQLType: "BOOLEAN", Metadata: parser.GetFieldMetadata("bool")},
	}

	if err := generator.GenerateAPI(tmpDir, "testmodule", "Post", fields, "multi"); err != nil {
		t.Fatalf("Failed to generate API: %v", err)
	}

	handlerPath := filepath.Join(tmpDir, "app", "api", "post.go")
	generated, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read generated handler: %v", err)
	}

	goldenPath := "testdata/golden/api_handler.go.golden"

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata/golden", 0755); err != nil {
			t.Fatalf("Failed to create golden directory: %v", err)
		}
		if err := os.WriteFile(goldenPath, generated, 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Log("Updated API golden file")
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("Golden file doesn't exist. Run with UPDATE_GOLDEN=1 to create it.")
		}
		t.Fatalf("Failed to read golden file: %v", err)
	}

	if string(golden) != string(generated) {
		t.Errorf("Generated code differs from golden file.\n"+
			"Run with UPDATE_GOLDEN=1 to update.\n"+
			"Golden: %d bytes, Generated: %d bytes",
			len(golden), len(generated))

		goldenLines := strings.Split(string(golden), "\n")
		genLines := strings.Split(string(generated), "\n")
		for i := 0; i < len(goldenLines) && i < len(genLines); i++ {
			if goldenLines[i] != genLines[i] {
				t.Logf("First difference at line %d:", i+1)
				t.Logf("  Golden:    %s", goldenLines[i])
				t.Logf("  Generated: %s", genLines[i])
				break
			}
		}
	}
}
