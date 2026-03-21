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

	if err := generator.GenerateResource(tmpDir, "testmodule", "User", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
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

	if err := generator.GenerateApp("testapp", "testapp", "multi", "tailwind", false); err != nil { // false = production mode
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

	if err := generator.GenerateResource(appDir, "testapp", "Post", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
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

	if err := generator.GenerateResource(tmpDir, "testmodule", "Post", parentFields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
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

	if err := generator.GenerateResource(tmpDir, "testmodule", "Comment", childFields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
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
	if err := generator.GenerateApp(appName, appName, "multi", "tailwind", false); err != nil {
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
	if err := generator.GenerateResource(appDir, appName, "Post", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}
	t.Log("✅ Resource generated")

	// Step 3: Add replace directives BEFORE generating auth
	// GenerateAuth runs `go get github.com/livetemplate/lvt@latest` which
	// transitively requires lvt/components — a local-only module not published
	// to any module proxy. The replace directives must be in place first.
	t.Log("Step 3: Adding replace directives for local packages...")
	goModPath := filepath.Join(appDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	replaceDirective := fmt.Sprintf("\nreplace github.com/livetemplate/lvt => %s\nreplace github.com/livetemplate/lvt/components => %s/components\n", origDir, origDir)
	if err := os.WriteFile(goModPath, append(goModContent, []byte(replaceDirective)...), 0644); err != nil {
		t.Fatalf("Failed to update go.mod: %v", err)
	}
	t.Log("✅ Replace directives added")

	// Step 4: Generate auth (with delay to avoid migration timestamp collision)
	t.Log("Step 4: Generating auth...")
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

	// Step 5: Run go mod tidy
	t.Log("Step 5: Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = appDir
	// Don't set GOWORK=off - we want to use the go.work file
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ go mod tidy completed")

	// Step 6: Generate sqlc code
	t.Log("Step 6: Generating sqlc code...")
	cmd = exec.Command("sqlc", "generate")
	cmd.Dir = filepath.Join(appDir, "database")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sqlc generate failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ sqlc generate completed")

	// Step 7: Build the app (uses go.work to find local lvt packages)
	t.Log("Step 7: Building app...")
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = appDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ Build successful")

	// Step 8: Run short tests (skip E2E which requires Docker/Chrome)
	t.Log("Step 8: Running generated tests (short mode)...")
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

// TestFileUploadResourceGeneration validates that generating a resource with
// file/image fields produces correct handler, SQL, and template output.
func TestFileUploadResourceGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "photo", Type: "image", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: true, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
		{Name: "doc", Type: "file", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: false, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Gallery", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// --- Verify handler ---
	handlerPath := filepath.Join(tmpDir, "app", "gallery", "gallery.go")
	handler, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read handler: %v", err)
	}
	handlerContent := string(handler)

	handlerChecks := map[string]string{
		"storage import":       `"github.com/livetemplate/lvt/pkg/storage"`,
		"os import":            `"os"`,
		"Store field":          `Store   storage.Store`,
		"Handler with store":   `func Handler(queries *models.Queries, store storage.Store) http.Handler`,
		"WithUpload photo":     `livetemplate.WithUpload("photo"`,
		"WithUpload doc":       `livetemplate.WithUpload("doc"`,
		"image accept filter":  `"image/*"`,
		"GetCompletedUploads":  `ctx.GetCompletedUploads("photo")`,
		"Store.Save":           `c.Store.Save(`,
		"Store.URL":            `c.Store.URL(`,
		"AddInput struct":      "type AddInput struct",
		"NonFile title in Add": "Title string",
		"photo upload var":     "var photoVal, photoFilename, photoContentType string",
		"doc upload var":       "var docVal, docFilename, docContentType string",
		"PhotoFilename param":  "PhotoFilename:",
		"DocContentType param": "DocContentType:",
		"Store.Delete":         "c.Store.Delete(",
	}

	for desc, substr := range handlerChecks {
		if !strings.Contains(handlerContent, substr) {
			t.Errorf("Handler missing %s: expected %q", desc, substr)
		}
	}

	// Verify file fields are NOT in AddInput struct (they come via upload protocol)
	// AddInput should have Title but not Photo or Doc fields
	addInputIdx := strings.Index(handlerContent, "type AddInput struct")
	updateInputIdx := strings.Index(handlerContent, "type UpdateInput struct")
	if addInputIdx < 0 || updateInputIdx < 0 {
		t.Fatal("Could not find AddInput or UpdateInput structs")
	}
	addInputSection := handlerContent[addInputIdx:updateInputIdx]
	if strings.Contains(addInputSection, `json:"photo"`) {
		t.Error("AddInput should not contain photo field (file data comes via upload protocol)")
	}
	if strings.Contains(addInputSection, `json:"doc"`) {
		t.Error("AddInput should not contain doc field (file data comes via upload protocol)")
	}
	if !strings.Contains(addInputSection, `json:"title"`) {
		t.Error("AddInput should contain title field")
	}

	// Verify search skips file fields
	if !strings.Contains(handlerContent, "Search across all text fields (skip file/image fields)") {
		t.Error("Handler search comment should mention skipping file/image fields")
	}

	// --- Verify handler has valid Go syntax ---
	cmd := exec.Command("go", "tool", "compile", "-o", "/dev/null", handlerPath)
	output, _ := cmd.CombinedOutput()
	if strings.Contains(string(output), "syntax error") {
		t.Errorf("Handler has syntax errors:\n%s", output)
	}

	// --- Verify SQL schema ---
	schemaPath := filepath.Join(tmpDir, "database", "schema.sql")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}
	schema := string(schemaData)

	sqlChecks := map[string]string{
		"photo column":       "photo TEXT NOT NULL DEFAULT ''",
		"photo_filename":     "photo_filename TEXT NOT NULL DEFAULT ''",
		"photo_content_type": "photo_content_type TEXT NOT NULL DEFAULT ''",
		"photo_size":         "photo_size INTEGER NOT NULL DEFAULT 0",
		"doc column":         "doc TEXT NOT NULL DEFAULT ''",
		"doc_filename":       "doc_filename TEXT NOT NULL DEFAULT ''",
		"doc_content_type":   "doc_content_type TEXT NOT NULL DEFAULT ''",
		"doc_size":           "doc_size INTEGER NOT NULL DEFAULT 0",
		"title column":       "title TEXT NOT NULL",
	}

	for desc, substr := range sqlChecks {
		if !strings.Contains(schema, substr) {
			t.Errorf("Schema missing %s: expected %q\nSchema:\n%s", desc, substr, schema)
		}
	}

	// --- Verify queries ---
	queriesPath := filepath.Join(tmpDir, "database", "queries.sql")
	queriesData, err := os.ReadFile(queriesPath)
	if err != nil {
		t.Fatalf("Failed to read queries.sql: %v", err)
	}
	queries := string(queriesData)

	// INSERT should have expanded columns
	queryChecks := map[string]string{
		"INSERT photo columns": "photo, photo_filename, photo_content_type, photo_size",
		"INSERT doc columns":   "doc, doc_filename, doc_content_type, doc_size",
		"UPDATE photo columns": "photo = ?, photo_filename = ?, photo_content_type = ?, photo_size = ?",
		"UPDATE doc columns":   "doc = ?, doc_filename = ?, doc_content_type = ?, doc_size = ?",
	}

	for desc, substr := range queryChecks {
		if !strings.Contains(queries, substr) {
			t.Errorf("Queries missing %s: expected %q\nQueries:\n%s", desc, substr, queries)
		}
	}

	// --- Verify template (form) ---
	tmplPath := filepath.Join(tmpDir, "app", "gallery", "gallery.tmpl")
	tmplData, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template: %v", err)
	}
	tmpl := string(tmplData)

	tmplChecks := map[string]string{
		"photo upload input":      `lvt-upload="photo"`,
		"photo accept":            `accept="image/*"`,
		"doc upload input":        `lvt-upload="doc"`,
		"upload progress":         `.lvt.Uploads "photo"`,
		"upload error check":      `.lvt.HasUploadError "photo"`,
		"edit current file":       `Leave empty to keep current file`,
		"image preview in edit":   `<img src="`,
		"file download in detail": `target="_blank"`,
	}

	for desc, substr := range tmplChecks {
		if !strings.Contains(tmpl, substr) {
			t.Errorf("Template missing %s: expected %q", desc, substr)
		}
	}

	// doc should NOT have accept="image/*"
	// Find the doc upload input and verify it doesn't have image accept
	docInputIdx := strings.Index(tmpl, `lvt-upload="doc"`)
	if docInputIdx < 0 {
		t.Error("Template should have doc upload input")
	} else {
		// Check the line containing the doc input
		lineStart := strings.LastIndex(tmpl[:docInputIdx], "\n")
		lineEnd := strings.Index(tmpl[docInputIdx:], "\n")
		if lineEnd < 0 {
			lineEnd = len(tmpl) - docInputIdx
		}
		docLine := tmpl[lineStart : docInputIdx+lineEnd]
		if strings.Contains(docLine, `accept="image/*"`) {
			t.Error("Doc file input should NOT have accept=\"image/*\" (only image fields should)")
		}
	}

	// --- Verify migration ---
	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations: %v", err)
	}

	found := false
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "galleries") {
			data, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read migration: %v", err)
			}
			migration := string(data)
			if !strings.Contains(migration, "photo_filename TEXT NOT NULL DEFAULT ''") {
				t.Error("Migration missing photo_filename column")
			}
			if !strings.Contains(migration, "doc_size INTEGER NOT NULL DEFAULT 0") {
				t.Error("Migration missing doc_size column")
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Migration file for galleries not found")
	}

	t.Log("✅ File upload resource generation test passed")
}

// TestFileUploadFullFlow generates a complete app with file upload fields,
// runs sqlc, and verifies the generated code compiles.
func TestFileUploadFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full flow test in short mode")
	}

	if _, err := exec.LookPath("sqlc"); err != nil {
		t.Fatal("sqlc not installed - run: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
	}

	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	appName := "uploadapp"
	appDir := filepath.Join(tmpDir, appName)

	// Step 1: Generate app
	t.Log("Step 1: Generating app...")
	if err := generator.GenerateApp(appName, appName, "multi", "tailwind", false); err != nil {
		t.Fatalf("Failed to generate app: %v", err)
	}

	// Step 2: Generate resource with file/image fields
	t.Log("Step 2: Generating resource with file/image fields...")
	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "photo", Type: "image", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: true, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
		{Name: "doc", Type: "file", GoType: "string", SQLType: "TEXT", IsFile: true, IsImage: false, Metadata: parser.FieldMetadata{HTMLInputType: "file"}},
		{Name: "views", Type: "int", GoType: "int64", SQLType: "INTEGER", Metadata: parser.GetFieldMetadata("int")},
	}
	if err := generator.GenerateResource(appDir, appName, "Gallery", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", false); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}
	t.Log("✅ Resource with file/image fields generated")

	// Step 3: Add replace directives
	t.Log("Step 3: Adding replace directives...")
	goModPath := filepath.Join(appDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	replaceDirective := fmt.Sprintf("\nreplace github.com/livetemplate/lvt => %s\nreplace github.com/livetemplate/lvt/components => %s/components\n", origDir, origDir)
	if err := os.WriteFile(goModPath, append(goModContent, []byte(replaceDirective)...), 0644); err != nil {
		t.Fatalf("Failed to update go.mod: %v", err)
	}

	// Step 4: Wire up route manually (file upload resources skip auto-injection)
	t.Log("Step 4: Wiring up file upload route in main.go...")
	mainGoPath := filepath.Join(appDir, "cmd", appName, "main.go")
	mainGoContent, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	mainGoStr := string(mainGoContent)
	// Add storage import
	mainGoStr = strings.Replace(mainGoStr,
		`"net/http"`,
		"\"net/http\"\n\t\"github.com/livetemplate/lvt/pkg/storage\"\n\t\""+appName+"/app/gallery\"",
		1)
	// Enable the queries variable (normally done by route injector)
	mainGoStr = strings.Replace(mainGoStr,
		`_, err := database.InitDB(dbPath)`,
		`queries, err := database.InitDB(dbPath)`,
		1)
	// Add store + route at the TODO comment
	mainGoStr = strings.Replace(mainGoStr,
		`// TODO: Add routes here (added automatically by `+"`lvt gen`"+`)`,
		"// File upload store\n\tstore := storage.NewLocalStore(\"uploads\", \"/uploads\")\n\thttp.Handle(\"/uploads/\", http.StripPrefix(\"/uploads/\", store.FileServer()))\n\thttp.Handle(\"/gallery\", gallery.Handler(queries, store))",
		1)
	if err := os.WriteFile(mainGoPath, []byte(mainGoStr), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}
	t.Log("✅ Route wired up")

	// Step 5: go mod tidy
	t.Log("Step 5: Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = appDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, output)
	}

	// Step 6: sqlc generate
	t.Log("Step 6: Generating sqlc code...")
	cmd = exec.Command("sqlc", "generate")
	cmd.Dir = filepath.Join(appDir, "database")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sqlc generate failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ sqlc generate completed")

	// Step 7: Verify sqlc generated the expanded columns in the model
	modelsDir := filepath.Join(appDir, "database", "models")
	modelsEntries, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("Failed to read models directory: %v", err)
	}
	foundModel := false
	for _, entry := range modelsEntries {
		if strings.HasSuffix(entry.Name(), ".go") && entry.Name() != "db.go" && entry.Name() != "models.go" {
			data, err := os.ReadFile(filepath.Join(modelsDir, entry.Name()))
			if err != nil {
				continue
			}
			content := string(data)
			if strings.Contains(content, "Gallery") {
				// Verify expanded file columns are present in the model
				for _, field := range []string{"PhotoFilename", "PhotoContentType", "PhotoSize", "DocFilename", "DocContentType", "DocSize"} {
					if !strings.Contains(content, field) {
						t.Errorf("sqlc model missing field %s", field)
					}
				}
				foundModel = true
				break
			}
		}
	}
	if !foundModel {
		t.Error("Could not find sqlc-generated Gallery model")
	}

	// Step 8: Build
	t.Log("Step 8: Building app...")
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = appDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\nOutput: %s", err, output)
	}
	t.Log("✅ Build successful — file upload code compiles")

	t.Log("✅ File upload full flow test passed!")
}

// TestAuthzResourceGeneration validates that generating a resource with --with-authz
// produces correct handler, SQL, and template output.
func TestAuthzResourceGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	fields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT", Metadata: parser.GetFieldMetadata("string")},
		{Name: "content", Type: "text", GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: parser.GetFieldMetadata("text")},
	}

	if err := generator.GenerateResource(tmpDir, "testmodule", "Post", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "", true); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// --- Verify handler ---
	handlerPath := filepath.Join(tmpDir, "app", "post", "post.go")
	handler, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read handler: %v", err)
	}
	handlerContent := string(handler)

	handlerChecks := map[string]string{
		"authz import":           `"github.com/livetemplate/lvt/pkg/authz"`,
		"policy registration":    `authz.Register("posts"`,
		"authz Can update":       `authz.Can(user, authz.ActionUpdate`,
		"authz Can delete":       `authz.Can(user, authz.ActionDelete`,
		"OwnedBy check":          `authz.OwnedBy(item.CreatedBy)`,
		"CreatedBy in Create":    `CreatedBy: ctx.UserID()`,
		"getUserRole method":     `func (c *PostController) getUserRole`,
		"CookieAuthenticator":    `authz.NewCookieAuthenticator`,
		"WithAuthenticator":      `livetemplate.WithAuthenticator`,
	}

	for desc, substr := range handlerChecks {
		if !strings.Contains(handlerContent, substr) {
			t.Errorf("Handler missing %s: expected %q", desc, substr)
		}
	}

	// Verify handler has valid Go syntax
	cmd := exec.Command("go", "tool", "compile", "-o", "/dev/null", handlerPath)
	output, _ := cmd.CombinedOutput()
	if strings.Contains(string(output), "syntax error") {
		t.Errorf("Handler has syntax errors:\n%s", output)
	}

	// --- Verify SQL schema ---
	schemaPath := filepath.Join(tmpDir, "database", "schema.sql")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}
	schema := string(schemaData)

	sqlChecks := map[string]string{
		"created_by column": "created_by TEXT NOT NULL REFERENCES users(id)",
		"created_by index":  "idx_posts_created_by",
	}
	for desc, substr := range sqlChecks {
		if !strings.Contains(schema, substr) {
			t.Errorf("Schema missing %s: expected %q\nSchema:\n%s", desc, substr, schema)
		}
	}

	// --- Verify queries ---
	queriesPath := filepath.Join(tmpDir, "database", "queries.sql")
	queriesData, err := os.ReadFile(queriesPath)
	if err != nil {
		t.Fatalf("Failed to read queries.sql: %v", err)
	}
	queries := string(queriesData)

	if !strings.Contains(queries, "created_by") {
		t.Errorf("Queries should include created_by in INSERT\nQueries:\n%s", queries)
	}

	// --- Verify migration ---
	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations: %v", err)
	}
	found := false
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "posts") {
			data, _ := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			migration := string(data)
			if !strings.Contains(migration, "created_by TEXT NOT NULL REFERENCES users(id)") {
				t.Error("Migration missing created_by column")
			}
			if !strings.Contains(migration, "idx_posts_created_by") {
				t.Error("Migration missing created_by index")
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Migration file for posts not found")
	}

	t.Log("✅ Authz resource generation test passed")
}
