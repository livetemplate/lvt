//go:build http

package e2e

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestKitRuntime_AllKits tests that all kits can successfully run and serve content.
// This is a critical regression test to catch issues like template name mismatches,
// routing problems, or handler configuration issues in fast tests (~5s per kit)
// instead of slow deployment tests (3-4 minutes).
func TestKitRuntime_AllKits(t *testing.T) {
	tests := []struct {
		name string
		kit  string
		port string
	}{
		{
			name: "simple kit",
			kit:  "simple",
			port: "9880",
		},
		{
			name: "single kit",
			kit:  "single",
			port: "9881",
		},
		{
			name: "multi kit",
			kit:  "multi",
			port: "9882",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create app with specific kit
			opts := &AppOptions{
				Kit: tt.kit,
			}
			appDir := createTestApp(t, tmpDir, "testapp", opts)

			t.Logf("Testing %s runtime...", tt.kit)

			// Start server
			handle, err := startServeInBackground(t, appDir, "--no-browser", "--port", tt.port)
			if err != nil {
				t.Fatalf("Failed to start server for %s: %v", tt.kit, err)
			}
			defer func() { _ = handle.Shutdown() }()

			// Wait for server to be ready
			url := fmt.Sprintf("http://localhost:%s", tt.port)
			waitForServer(t, url, 10*time.Second)

			// Verify HTTP 200
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Server not responding for %s: %v", tt.kit, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Fatalf("Expected status 200 for %s, got %d", tt.kit, resp.StatusCode)
			}

			// Verify response is HTML
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body for %s: %v", tt.kit, err)
			}

			bodyStr := string(body)
			if !strings.Contains(bodyStr, "<html") && !strings.Contains(bodyStr, "<!DOCTYPE") {
				t.Fatalf("Response for %s doesn't appear to be HTML. Got: %s", tt.kit, bodyStr[:min(200, len(bodyStr))])
			}

			t.Logf("✅ %s: HTTP 200 + HTML rendering verified", tt.kit)
		})
	}
}

// TestKitRuntime_TemplateRendering is a focused regression test for template name matching.
// This test specifically verifies that templates are correctly discovered and rendered,
// preventing bugs like the one where main.go used livetemplate.Must(livetemplate.New("[[.AppName]]", ...))
// but the actual template file was named index.tmpl.
func TestKitRuntime_TemplateRendering(t *testing.T) {
	tests := []struct {
		name           string
		kit            string
		port           string
		expectedInHTML []string // Content that should appear in rendered HTML
	}{
		{
			name: "simple kit template",
			kit:  "simple",
			port: "9890",
			expectedInHTML: []string{
				"<html",
				"<body",
				// Simple kit should have counter functionality
				"counter",
			},
		},
		{
			name: "single kit template",
			kit:  "single",
			port: "9891",
			expectedInHTML: []string{
				"<html",
				"<body",
				// Single kit has posts functionality
			},
		},
		{
			name: "multi kit template",
			kit:  "multi",
			port: "9892",
			expectedInHTML: []string{
				"<html",
				"<body",
				// Multi kit has posts functionality
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			opts := &AppOptions{
				Kit: tt.kit,
			}
			appDir := createTestApp(t, tmpDir, "testapp", opts)

			t.Logf("Testing template rendering for %s...", tt.kit)

			handle, err := startServeInBackground(t, appDir, "--no-browser", "--port", tt.port)
			if err != nil {
				t.Fatalf("Failed to start server: %v", err)
			}
			defer func() { _ = handle.Shutdown() }()

			// Wait for server to be ready
			url := fmt.Sprintf("http://localhost:%s", tt.port)
			waitForServer(t, url, 10*time.Second)

			// Fetch root page
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Failed to fetch page: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Fatalf("Expected status 200, got %d (template may not be found)", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			bodyStr := strings.ToLower(string(body))

			// Verify expected content is present
			for _, expected := range tt.expectedInHTML {
				if !strings.Contains(bodyStr, strings.ToLower(expected)) {
					t.Errorf("Expected to find '%s' in rendered HTML for %s kit", expected, tt.kit)
					t.Logf("Response preview: %s", bodyStr[:min(500, len(bodyStr))])
				}
			}

			t.Logf("✅ %s: Template rendered correctly with expected content", tt.kit)
		})
	}
}
