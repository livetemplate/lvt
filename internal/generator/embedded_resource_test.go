package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/internal/parser"
)

// TestEmbeddedResourceGeneration tests the end-to-end flow:
// 1. Generate posts as standalone resource
// 2. Generate comments with --parent posts
// 3. Verify child files exist
// 4. Verify parent files are modified
// 5. Verify no separate route for comments
func TestEmbeddedResourceGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up a minimal project structure
	setupMinimalProject(t, tmpDir)

	// Step 1: Generate posts as a standalone resource
	postFields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "content", Type: "text", GoType: "string", SQLType: "TEXT", IsTextarea: true},
	}

	err := GenerateResource(tmpDir, "testapp", "posts", postFields, "multi", "tailwind", "tailwind", "infinite", 20, "page", "", false, false)
	if err != nil {
		t.Fatalf("failed to generate posts: %v", err)
	}

	// Verify posts files exist
	assertFileExists(t, filepath.Join(tmpDir, "app", "posts", "posts.go"))
	assertFileExists(t, filepath.Join(tmpDir, "app", "posts", "posts.tmpl"))

	// Step 2: Generate comments with --parent posts
	commentFields := []parser.Field{
		{Name: "post_id", Type: "references:posts", GoType: "string", SQLType: "TEXT", IsReference: true, ReferencedTable: "posts"},
		{Name: "author", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "text", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	err = GenerateResource(tmpDir, "testapp", "comments", commentFields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "posts", false, false)
	if err != nil {
		t.Fatalf("failed to generate embedded comments: %v", err)
	}

	// Step 3: Verify child files exist
	assertFileExists(t, filepath.Join(tmpDir, "app", "comments", "comments.go"))
	assertFileExists(t, filepath.Join(tmpDir, "app", "comments", "comments.tmpl"))

	// Verify child handler uses EmbeddedState/EmbeddedController
	childHandler, err := os.ReadFile(filepath.Join(tmpDir, "app", "comments", "comments.go"))
	if err != nil {
		t.Fatal(err)
	}
	childSrc := string(childHandler)

	if !strings.Contains(childSrc, "EmbeddedState") {
		t.Error("child handler should define EmbeddedState")
	}
	if !strings.Contains(childSrc, "EmbeddedController") {
		t.Error("child handler should define EmbeddedController")
	}
	if !strings.Contains(childSrc, "NewEmbeddedController") {
		t.Error("child handler should define NewEmbeddedController")
	}
	// Child handler should NOT have a Handler() func
	if strings.Contains(childSrc, "func Handler(") {
		t.Error("child handler should NOT have a standalone Handler func")
	}

	// Step 4: Verify parent files are modified
	parentHandler, err := os.ReadFile(filepath.Join(tmpDir, "app", "posts", "posts.go"))
	if err != nil {
		t.Fatal(err)
	}
	parentSrc := string(parentHandler)

	if !strings.Contains(parentSrc, `"testapp/app/comments"`) {
		t.Error("parent handler should import comments package")
	}
	if !strings.Contains(parentSrc, `Comments *comments.EmbeddedState`) {
		t.Error("parent state should have Comments field")
	}
	if !strings.Contains(parentSrc, `CommentsCtrl *comments.EmbeddedController`) {
		t.Error("parent controller should have CommentsCtrl field")
	}
	if !strings.Contains(parentSrc, "CommentAdd") {
		t.Error("parent should have CommentAdd forwarding method")
	}
	if !strings.Contains(parentSrc, `"app/comments/comments.tmpl"`) {
		t.Error("parent should parse child template")
	}

	// Verify parent template includes child section
	parentTmpl, err := os.ReadFile(filepath.Join(tmpDir, "app", "posts", "posts.tmpl"))
	if err != nil {
		t.Fatal(err)
	}
	tmplSrc := string(parentTmpl)

	if !strings.Contains(tmplSrc, `{{template "comments:section" .Comments}}`) {
		t.Error("parent template should include comments section")
	}

	// Step 5: Verify no separate route for comments
	mainGo, err := os.ReadFile(filepath.Join(tmpDir, "cmd", "testapp", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	mainSrc := string(mainGo)

	if strings.Contains(mainSrc, `"/comments"`) {
		t.Error("should NOT inject a separate route for embedded child resource")
	}
	// But posts route should exist
	if !strings.Contains(mainSrc, `"/posts"`) {
		t.Error("posts route should be injected")
	}

	// Step 6: Verify queries include filtered-by-parent query
	queriesSQL, err := os.ReadFile(filepath.Join(tmpDir, "database", "queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	queriesSrc := string(queriesSQL)

	if !strings.Contains(queriesSrc, "GetCommentsByPostID") {
		t.Error("queries should include GetCommentsByPostID")
	}
	// Update query should NOT include the parent FK field
	if strings.Contains(queriesSrc, "SET post_id") {
		t.Error("Update query should NOT include parent FK field (post_id)")
	}
	// But it should include the other fields
	if !strings.Contains(queriesSrc, "author = ?") {
		t.Error("Update query should include author field")
	}

	// Step 7: Verify child template has correct action names
	childTmpl, err := os.ReadFile(filepath.Join(tmpDir, "app", "comments", "comments.tmpl"))
	if err != nil {
		t.Fatal(err)
	}
	childTmplSrc := string(childTmpl)

	if !strings.Contains(childTmplSrc, `lvt-submit="comment_add"`) {
		t.Error("child template should use comment_add action")
	}
	if !strings.Contains(childTmplSrc, `lvt-click="comment_edit"`) {
		t.Error("child template should use comment_edit action")
	}
	if !strings.Contains(childTmplSrc, `{{define "comments:section"}}`) {
		t.Error("child template should define comments:section block")
	}
}

func TestEmbeddedResourceGeneration_MissingParent(t *testing.T) {
	tmpDir := t.TempDir()
	setupMinimalProject(t, tmpDir)

	fields := []parser.Field{
		{Name: "post_id", Type: "references:posts", GoType: "string", SQLType: "TEXT", IsReference: true, ReferencedTable: "posts"},
		{Name: "text", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	// Should fail because posts resource doesn't exist
	err := GenerateResource(tmpDir, "testapp", "comments", fields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "posts", false, false)
	if err == nil {
		t.Error("expected error when parent resource doesn't exist")
	}
}

func TestEmbeddedResourceGeneration_MissingReferenceField(t *testing.T) {
	tmpDir := t.TempDir()
	setupMinimalProject(t, tmpDir)

	// Generate parent first
	postFields := []parser.Field{
		{Name: "title", Type: "string", GoType: "string", SQLType: "TEXT"},
	}
	if err := GenerateResource(tmpDir, "testapp", "posts", postFields, "multi", "tailwind", "tailwind", "infinite", 20, "page", "", false, false); err != nil {
		t.Fatalf("failed to generate posts: %v", err)
	}

	// Try to generate child WITHOUT a reference field to the parent
	commentFields := []parser.Field{
		{Name: "author", Type: "string", GoType: "string", SQLType: "TEXT"},
		{Name: "text", Type: "string", GoType: "string", SQLType: "TEXT"},
	}

	err := GenerateResource(tmpDir, "testapp", "comments", commentFields, "multi", "tailwind", "tailwind", "infinite", 20, "modal", "posts", false, false)
	if err == nil {
		t.Error("expected error when child has no reference field for parent")
	}
}

// setupMinimalProject creates a minimal project structure for testing
func setupMinimalProject(t *testing.T, dir string) {
	t.Helper()

	// go.mod
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module testapp\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// cmd/testapp/main.go
	cmdDir := filepath.Join(dir, "cmd", "testapp")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatal(err)
	}
	mainGo := `package main

import (
	"net/http"
	"testapp/database"
)

func main() {
	_, err := database.InitDB("app.db")
	if err != nil {
		panic(err)
	}

	// TODO: Add routes here
	// Example: http.Handle("/path", handler.Handler(queries))

	http.ListenAndServe(":8080", nil)
}
`
	if err := os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}

	// database directory
	dbDir := filepath.Join(dir, "database")
	if err := os.MkdirAll(filepath.Join(dbDir, "migrations"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dbDir, "schema.sql"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dbDir, "queries.sql"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// .lvtrc config
	lvtrc := `{"kit": "multi", "styles": "tailwind"}`
	if err := os.WriteFile(filepath.Join(dir, ".lvtrc"), []byte(lvtrc), 0644); err != nil {
		t.Fatal(err)
	}

	// .lvtresources
	if err := os.WriteFile(filepath.Join(dir, ".lvtresources"), []byte("[]"), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}
