package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livetemplate/lvt/internal/generator"
	"github.com/livetemplate/lvt/internal/parser"
)

// TestGeneratedCodeSyntax validates that generated code has valid Go syntax
func TestGeneratedCodeSyntax(t *testing.T) {
	tmpDir := t.TempDir()

	// Create database directory structure
	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	// Generate CRUD resource
	fields := []parser.Field{
		{Name: "name", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "email", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "User", fields, "multi", "tailwind", "infinite", 20, "modal"); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// Generate view
	if err := generator.GenerateView(tmpDir, "testmodule", "Counter", "multi", "tailwind"); err != nil {
		t.Fatalf("Failed to generate view: %v", err)
	}

	// Check generated Go files for syntax errors
	goFiles := []string{
		filepath.Join(tmpDir, "app", "user", "user.go"),
		filepath.Join(tmpDir, "app", "counter", "counter.go"),
	}

	for _, file := range goFiles {
		// Use go/parser to check syntax
		cmd := exec.Command("go", "tool", "compile", "-o", "/dev/null", file)
		// We expect this to fail due to unresolved imports, but syntax should be valid
		output, _ := cmd.CombinedOutput()

		// Check for syntax errors (not import errors)
		if strings.Contains(string(output), "syntax error") {
			t.Errorf("Syntax error in %s:\n%s", file, output)
		}
	}

	t.Log("✅ Generated code has valid Go syntax")
}

// TestGeneratedFilesExist validates that all expected files are generated
func TestGeneratedFilesExist(t *testing.T) {
	tmpDir := t.TempDir()

	// Save original working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	// Restore working directory after test
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Logf("Warning: failed to restore directory: %v", err)
		}
	})

	// Generate app
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if err := generator.GenerateApp("testapp", "testapp", "multi", false); err != nil { // false = production mode
		t.Fatalf("Failed to generate app: %v", err)
	}

	appDir := "testapp"

	// Check app files
	expectedAppFiles := []string{
		"go.mod",
		"README.md",
		"cmd/testapp/main.go",
		"database/db.go",
		"database/schema.sql",
		"database/queries.sql",
		"database/sqlc.yaml",
	}

	for _, file := range expectedAppFiles {
		path := filepath.Join(appDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file not found: %s", path)
		}
	}

	// Generate resource
	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	if err := generator.GenerateResource(appDir, "testapp", "Post", fields, "multi", "tailwind", "infinite", 20, "modal"); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// Check resource files
	expectedResourceFiles := []string{
		"app/post/post.go",
		"app/post/post.tmpl",
		"app/post/post_test.go",
	}

	for _, file := range expectedResourceFiles {
		path := filepath.Join(appDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected resource file not found: %s", path)
		}
	}

	// Generate view
	if err := generator.GenerateView(appDir, "testapp", "Dashboard", "multi", "tailwind"); err != nil {
		t.Fatalf("Failed to generate view: %v", err)
	}

	// Check view files
	expectedViewFiles := []string{
		"app/dashboard/dashboard.go",
		"app/dashboard/dashboard.tmpl",
		"app/dashboard/dashboard_test.go",
	}

	for _, file := range expectedViewFiles {
		path := filepath.Join(appDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected view file not found: %s", path)
		}
	}

	t.Log("✅ All expected files generated")
}

// TestForeignKeyGeneration validates that foreign keys are properly generated
func TestForeignKeyGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create database directory structure
	dbDir := filepath.Join(tmpDir, "database", "migrations")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	// Generate parent resource (posts)
	parentFields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "content", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Post", parentFields, "multi", "tailwind", "infinite", 20, "modal"); err != nil {
		t.Fatalf("Failed to generate parent resource: %v", err)
	}

	// Generate child resource with foreign key
	childFields := []parser.Field{
		{
			Name:            "post_id",
			Type:            "references:posts",
			GoType:          "string",
			SQLType:         "TEXT",
			IsReference:     true,
			ReferencedTable: "posts",
			OnDelete:        "CASCADE",
		},
		{Name: "author", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "text", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Comment", childFields, "multi", "tailwind", "infinite", 20, "modal"); err != nil {
		t.Fatalf("Failed to generate child resource: %v", err)
	}

	// Read the generated migration file for comments
	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	var commentsMigration string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "comments") {
			data, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read migration file: %v", err)
			}
			commentsMigration = string(data)
			break
		}
	}

	if commentsMigration == "" {
		t.Fatal("Comments migration file not found")
	}

	// Validate that foreign key is defined inline (not using ALTER TABLE)
	if strings.Contains(commentsMigration, "ALTER TABLE") && strings.Contains(commentsMigration, "ADD CONSTRAINT") {
		t.Errorf("Migration uses ALTER TABLE ADD CONSTRAINT which is not supported in SQLite")
	}

	// Validate that inline FOREIGN KEY definition exists
	if !strings.Contains(commentsMigration, "FOREIGN KEY") {
		t.Errorf("Migration missing FOREIGN KEY definition")
	}

	// Validate that FOREIGN KEY references the correct table
	if !strings.Contains(commentsMigration, "REFERENCES posts(id)") {
		t.Errorf("FOREIGN KEY does not reference posts(id)")
	}

	// Validate ON DELETE CASCADE
	if !strings.Contains(commentsMigration, "ON DELETE CASCADE") {
		t.Errorf("FOREIGN KEY missing ON DELETE CASCADE")
	}

	// Validate that index is created for foreign key column
	if !strings.Contains(commentsMigration, "idx_comments_post_id") {
		t.Errorf("Index not created for foreign key column post_id")
	}

	// Check schema.sql as well
	schemaPath := filepath.Join(tmpDir, "database", "schema.sql")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}
	schema := string(schemaData)

	// Validate schema also has inline FOREIGN KEY
	if !strings.Contains(schema, "FOREIGN KEY (post_id) REFERENCES posts(id)") {
		t.Errorf("schema.sql missing inline FOREIGN KEY definition")
	}

	t.Log("✅ Foreign key generation test passed")
	t.Logf("Migration content:\n%s", commentsMigration)
}

// TestGeneratedAppFullFlow tests the complete workflow: generate app, resource, auth,
// then build and run the generated tests. This ensures the generated app works in a single shot.
func TestGeneratedAppFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full flow test in short mode")
	}

	// sqlc is required - fail if not installed
	if _, err := exec.LookPath("sqlc"); err != nil {
		t.Fatal("sqlc not installed - run: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
	}

	tmpDir := t.TempDir()

	// Save original working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Logf("Warning: failed to restore directory: %v", err)
		}
	})

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	appName := "testblog"
	appDir := filepath.Join(tmpDir, appName)

	// Step 1: Generate app
	t.Log("Step 1: Generating app...")
	if err := generator.GenerateApp(appName, appName, "multi", false); err != nil {
		t.Fatalf("Failed to generate app: %v", err)
	}
	t.Log("✅ App generated")

	// Step 2: Generate resource
	t.Log("Step 2: Generating resource...")
	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "content", Type: "text", GoType: "string", SQLType: "TEXT"},
		{Name: "published", Type: "bool", GoType: "bool", SQLType: "BOOLEAN"},
	}
	if err := generator.GenerateResource(appDir, appName, "Post", fields, "multi", "tailwind", "infinite", 20, "modal"); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}
	t.Log("✅ Resource generated")

	// Step 3: Generate auth (with delay to avoid migration timestamp collision)
	t.Log("Step 3: Generating auth...")
	time.Sleep(2 * time.Second)
	authConfig := &generator.AuthConfig{
		ModuleName:         appName,
		TableName:          "users",
		EnablePassword:     true,
		EnableMagicLink:    true,
		EnableEmailConfirm: true,
	}
	if err := generator.GenerateAuth(appDir, authConfig); err != nil {
		t.Fatalf("Failed to generate auth: %v", err)
	}
	t.Log("✅ Auth generated")

	// Step 3.5: Add replace directive to use local lvt packages
	// This is needed because the new pkg/* packages are not yet published
	t.Log("Step 3.5: Adding replace directive for local packages...")
	goModPath := filepath.Join(appDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	// Add replace directive at the end
	replaceDirective := fmt.Sprintf("\nreplace github.com/livetemplate/lvt => %s\n", origDir)
	if err := os.WriteFile(goModPath, append(goModContent, []byte(replaceDirective)...), 0644); err != nil {
		t.Fatalf("Failed to update go.mod: %v", err)
	}
	t.Log("✅ Replace directive added")

	// Step 4: Run go mod tidy
	t.Log("Step 4: Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = appDir
	// Don't set GOWORK=off - we want to use the go.work file
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ go mod tidy completed")

	// Step 5: Generate sqlc code
	t.Log("Step 5: Generating sqlc code...")
	cmd = exec.Command("sqlc", "generate")
	cmd.Dir = filepath.Join(appDir, "database")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sqlc generate failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ sqlc generate completed")

	// Step 6: Build the app (uses go.work to find local lvt packages)
	t.Log("Step 6: Building app...")
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = appDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ Build successful")

	// Step 7: Run short tests (skip E2E which requires Docker/Chrome)
	t.Log("Step 7: Running generated tests (short mode)...")
	cmd = exec.Command("go", "test", "./...", "-short", "-v")
	cmd.Dir = appDir
	output, err := cmd.CombinedOutput()
	t.Logf("Test output:\n%s", output)
	if err != nil {
		t.Fatalf("go test failed: %v", err)
	}
	t.Log("✅ Short tests passed")

	t.Log("✅ Full flow test passed - generated app works in a single shot!")
}
