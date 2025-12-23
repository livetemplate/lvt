//go:build integration

// Package e2e contains end-to-end tests for LiveTemplate.
// NOTE: These tests are tagged 'integration' because they create full apps,
// which is slow and requires all dependencies (sqlc, go, etc).
// Run with: go test -tags=integration ./e2e/...
//
// HTTP Tests (this file):
// Tests tagged with "http" run without a browser, using direct HTTP requests.
// They validate server-side rendering, form submission, and template correctness.
// Run with: go test -tags=http ./e2e/...
//
// Browser Tests:
// Tests tagged with "browser" require Chrome/chromedp for full browser testing.
// They validate JavaScript execution, WebSocket behavior, and DOM interactions.
// Run with: go test -tags=browser ./e2e/...
package e2e

import (
	"path/filepath"
	"testing"
	"time"

	lvttest "github.com/livetemplate/lvt/testing"
)

// TestHTTP_AppHealthCheck demonstrates HTTP testing pattern for health checks.
// This is a simple example showing how to use the HTTP testing framework.
func TestHTTP_AppHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	// Create a test app using existing helpers
	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "httptest", nil)

	// Generate a resource to have something to test
	if err := runLvtCommand(t, appDir, "gen", "resource", "Post", "title:string", "body:text"); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}

	// Run sqlc to generate database code
	runSqlc(t, appDir)

	// Setup HTTP test - starts the app server
	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
		AppPath: filepath.Join(appDir, "cmd", "httptest", "main.go"),
		AppDir:  appDir,
		Timeout: 30 * time.Second,
	})

	// Wait for server to be ready
	if err := test.WaitForServer(10 * time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test health endpoint
	resp := test.Get("/health")
	assert := lvttest.NewHTTPAssert(resp)
	assert.StatusOK(t)
	assert.Contains(t, `"status":"ok"`)
}

// TestHTTP_HomePageRendering tests that the home page renders correctly.
func TestHTTP_HomePageRendering(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "httptest2", nil)
	runSqlc(t, appDir)

	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
		AppPath: filepath.Join(appDir, "cmd", "httptest2", "main.go"),
		AppDir:  appDir,
		Timeout: 30 * time.Second,
	})

	if err := test.WaitForServer(10 * time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test home page
	resp := test.Get("/")
	assert := lvttest.NewHTTPAssert(resp)
	assert.StatusOK(t)
	assert.ContentTypeHTML(t)
	assert.Contains(t, "httptest2") // App name should appear
	assert.NoTemplateErrors(t)      // No unflattened {{.Field}} expressions
}

// TestHTTP_ResourceListPage tests that a resource list page renders correctly.
func TestHTTP_ResourceListPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "httptest3", nil)

	if err := runLvtCommand(t, appDir, "gen", "resource", "Task", "title:string", "done:bool"); err != nil {
		t.Fatalf("Failed to generate resource: %v", err)
	}
	runSqlc(t, appDir)

	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
		AppPath: filepath.Join(appDir, "cmd", "httptest3", "main.go"),
		AppDir:  appDir,
		Timeout: 30 * time.Second,
	})

	if err := test.WaitForServer(10 * time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test resource list page
	resp := test.Get("/task")
	assert := lvttest.NewHTTPAssert(resp)
	assert.StatusOK(t)
	assert.ContentTypeHTML(t)
	assert.Contains(t, "Task")
	assert.NoTemplateErrors(t)
	assert.HasElement(t, "table") // Should have a table for listing
}

// TestHTTP_SecurityHeaders tests that security headers are set correctly.
func TestHTTP_SecurityHeaders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "httptest4", nil)
	runSqlc(t, appDir)

	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
		AppPath: filepath.Join(appDir, "cmd", "httptest4", "main.go"),
		AppDir:  appDir,
		Timeout: 30 * time.Second,
	})

	if err := test.WaitForServer(10 * time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	resp := test.Get("/")
	assert := lvttest.NewHTTPAssert(resp)
	assert.StatusOK(t)

	// Check security headers
	assert.HasHeader(t, "X-Content-Type-Options")
	assert.HasHeader(t, "X-Frame-Options")
	assert.HasHeader(t, "X-Xss-Protection")
}

// TestHTTP_TemplateExpressionValidation tests that unflattened templates are caught.
func TestHTTP_TemplateExpressionValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTTP test in short mode")
	}

	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "httptest5", nil)

	// Generate multiple resources to ensure all templates render correctly
	resources := []struct {
		name   string
		fields string
	}{
		{"Article", "title:string,content:text,published:bool"},
		{"Comment", "author:string,body:text"},
	}

	for _, r := range resources {
		if err := runLvtCommand(t, appDir, "gen", "resource", r.name, r.fields); err != nil {
			t.Fatalf("Failed to generate resource %s: %v", r.name, err)
		}
	}
	runSqlc(t, appDir)

	test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
		AppPath: filepath.Join(appDir, "cmd", "httptest5", "main.go"),
		AppDir:  appDir,
		Timeout: 30 * time.Second,
	})

	if err := test.WaitForServer(10 * time.Second); err != nil {
		t.Fatalf("Server not ready: %v", err)
	}

	// Test each resource page for template errors
	pages := []string{"/", "/article", "/comment"}
	for _, page := range pages {
		t.Run(page, func(t *testing.T) {
			resp := test.Get(page)
			assert := lvttest.NewHTTPAssert(resp)
			assert.StatusOK(t)
			assert.NoTemplateErrors(t)
		})
	}
}
