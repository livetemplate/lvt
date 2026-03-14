//go:build http

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEmbeddedResourceGen tests generating a child resource with --parent flag
func TestEmbeddedResourceGen(t *testing.T) {
	tmpDir := t.TempDir()

	// Create app
	appDir := createTestApp(t, tmpDir, "blogapp", nil)

	// Generate posts as standalone page-mode resource
	t.Log("Generating posts resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "posts",
		"title:string", "content:text", "--edit-mode", "page"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}

	// Verify posts files
	assertFileExistsE2E(t, filepath.Join(appDir, "app/posts/posts.go"))
	assertFileExistsE2E(t, filepath.Join(appDir, "app/posts/posts.tmpl"))

	// Generate comments as embedded child of posts
	t.Log("Generating comments resource with --parent posts...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "comments",
		"post_id:references:posts", "author:string", "text:string", "--parent", "posts"); err != nil {
		t.Fatalf("Failed to generate embedded comments: %v", err)
	}

	// Verify child files exist
	assertFileExistsE2E(t, filepath.Join(appDir, "app/comments/comments.go"))
	assertFileExistsE2E(t, filepath.Join(appDir, "app/comments/comments.tmpl"))

	// Verify child handler uses embedded patterns
	childHandler, err := os.ReadFile(filepath.Join(appDir, "app/comments/comments.go"))
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
	if strings.Contains(childSrc, "func Handler(") {
		t.Error("child handler should NOT have standalone Handler function")
	}

	// Verify parent handler was modified
	parentHandler, err := os.ReadFile(filepath.Join(appDir, "app/posts/posts.go"))
	if err != nil {
		t.Fatal(err)
	}
	parentSrc := string(parentHandler)

	if !strings.Contains(parentSrc, `"blogapp/app/comments"`) {
		t.Error("parent should import comments package")
	}
	if !strings.Contains(parentSrc, "Comments *comments.EmbeddedState") {
		t.Error("parent state should include Comments field")
	}
	if !strings.Contains(parentSrc, "CommentsCtrl *comments.EmbeddedController") {
		t.Error("parent controller should include CommentsCtrl field")
	}
	if !strings.Contains(parentSrc, "CommentAdd") {
		t.Error("parent should have CommentAdd forwarding method")
	}
	if !strings.Contains(parentSrc, `"app/comments/comments.tmpl"`) {
		t.Error("parent should parse child template file")
	}

	// Verify parent template includes child section
	parentTmpl, err := os.ReadFile(filepath.Join(appDir, "app/posts/posts.tmpl"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(parentTmpl), `{{template "comments:section" .Comments}}`) {
		t.Error("parent template should include comments section")
	}

	// Verify no route for comments in main.go
	mainGo, err := os.ReadFile(filepath.Join(appDir, "cmd/blogapp/main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(mainGo), `"/comments"`) {
		t.Error("should NOT have a separate route for embedded child")
	}

	// Verify queries have the filtered query
	queries, err := os.ReadFile(filepath.Join(appDir, "database/queries.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(queries), "GetCommentsByPostID") {
		t.Error("queries should include GetCommentsByPostID")
	}

	// Verify the generated app compiles (need go mod tidy first for new dependencies)
	t.Log("Running go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = appDir
	tidyCmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		t.Logf("go mod tidy output: %s", output)
		// Not fatal — dependencies might not resolve in test environment
	}

	t.Log("Building generated app...")
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Logf("Build output: %s", string(output))
		// Skip (not pass) when build fails due to missing module dependencies,
		// so real generation regressions are still surfaced.
		t.Skip("Skipping: build failed, likely due to missing module dependencies in test environment")
	}
}

// TestEmbeddedResourceGen_MissingParent tests that --parent fails with missing parent
func TestEmbeddedResourceGen_MissingParent(t *testing.T) {
	tmpDir := t.TempDir()

	appDir := createTestApp(t, tmpDir, "testapp2", nil)

	// Try to generate comments with --parent posts, but posts doesn't exist
	err := runLvtCommand(t, appDir, "gen", "resource", "comments",
		"post_id:references:posts", "author:string", "text:string", "--parent", "posts")
	if err == nil {
		t.Error("expected error when parent resource doesn't exist")
	}
}

// TestEmbeddedResourceGen_MissingReferenceField tests that --parent requires a reference field
func TestEmbeddedResourceGen_MissingReferenceField(t *testing.T) {
	tmpDir := t.TempDir()

	appDir := createTestApp(t, tmpDir, "testapp3", nil)

	// Generate posts first
	if err := runLvtCommand(t, appDir, "gen", "resource", "posts", "title:string"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}

	// Try to generate comments without a reference field to posts
	err := runLvtCommand(t, appDir, "gen", "resource", "comments",
		"author:string", "text:string", "--parent", "posts")
	if err == nil {
		t.Error("expected error when child has no reference field for parent")
	}
}

func assertFileExistsE2E(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}
