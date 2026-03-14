//go:build browser

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestProductionTemplate verifies production-ready features in generated apps:
// health endpoints, security headers, structured logging, graceful shutdown,
// .env.example generation, and .gitignore generation.
func TestProductionTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Step 1: Create app
	t.Log("Step 1: Creating test app...")
	appDir := createTestApp(t, tmpDir, "prodtest", &AppOptions{Kit: "multi"})

	// Step 2: Verify .env.example and .gitignore were generated
	t.Run("EnvFiles", func(t *testing.T) {
		envExample := filepath.Join(appDir, ".env.example")
		content, err := os.ReadFile(envExample)
		if err != nil {
			t.Fatalf(".env.example not found: %v", err)
		}
		envStr := string(content)

		for _, expected := range []string{"PORT=", "APP_ENV=", "LOG_LEVEL=", "DATABASE_PATH=", "CLIENT_LIB_PATH"} {
			if !strings.Contains(envStr, expected) {
				t.Errorf(".env.example missing %q", expected)
			}
		}

		gitignore := filepath.Join(appDir, ".gitignore")
		giContent, err := os.ReadFile(gitignore)
		if err != nil {
			t.Fatalf(".gitignore not found: %v", err)
		}
		if !strings.Contains(string(giContent), ".env") {
			t.Error(".gitignore missing .env entry")
		}
	})

	// Step 3: Build and run the app natively
	t.Log("Step 3: Building and running app...")
	port := allocateTestPort()
	serverCmd, serverLogPath := buildAndRunProdApp(t, appDir, port)
	serverURL := fmt.Sprintf("http://localhost:%d", port)

	// Step 4: Test health endpoints
	t.Run("HealthLive", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/health/live")
		if err != nil {
			t.Fatalf("GET /health/live failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]string
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if result["status"] != "healthy" {
			t.Errorf("expected status=healthy, got %q", result["status"])
		}
	})

	t.Run("HealthReady", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/health/ready")
		if err != nil {
			t.Fatalf("GET /health/ready failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]string
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if result["status"] != "healthy" {
			t.Errorf("expected status=healthy, got %q", result["status"])
		}
	})

	// Step 5: Test security headers
	t.Run("SecurityHeaders", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()

		expectedHeaders := map[string]string{
			"X-Content-Type-Options": "nosniff",
			"X-Xss-Protection":       "1; mode=block",
			"X-Frame-Options":        "DENY",
		}

		for header, expected := range expectedHeaders {
			got := resp.Header.Get(header)
			if got != expected {
				t.Errorf("header %s: expected %q, got %q", header, expected, got)
			}
		}

		csp := resp.Header.Get("Content-Security-Policy")
		if csp == "" {
			t.Error("missing Content-Security-Policy header")
		}
		if !strings.Contains(csp, "default-src 'self'") {
			t.Errorf("CSP missing default-src 'self': %q", csp)
		}
	})

	// Step 6: Test home page renders via chromedp
	t.Run("HomePageRenders", func(t *testing.T) {
		ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(t.Logf))
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var html string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(serverURL),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.OuterHTML("html", &html),
		); err != nil {
			t.Fatalf("chromedp failed: %v", err)
		}

		if !strings.Contains(html, "prodtest") {
			t.Error("home page does not contain app name")
		}
	})

	// Step 7: Test graceful shutdown
	t.Run("GracefulShutdown", func(t *testing.T) {
		// Send SIGTERM and verify server exits cleanly
		if err := serverCmd.Process.Signal(syscall.SIGTERM); err != nil {
			t.Fatalf("failed to send SIGTERM: %v", err)
		}

		// Wait for process to exit (with timeout)
		done := make(chan error, 1)
		go func() { done <- serverCmd.Wait() }()

		select {
		case err := <-done:
			if err != nil {
				// On SIGTERM, exit code may be non-zero depending on OS
				t.Logf("Server exited: %v (expected on SIGTERM)", err)
			} else {
				t.Log("Server exited cleanly")
			}
		case <-time.After(10 * time.Second):
			_ = serverCmd.Process.Kill()
			t.Fatal("server did not exit within 10 seconds after SIGTERM")
		}

		// Verify server logs show clean shutdown
		logContent, err := os.ReadFile(serverLogPath)
		if err != nil {
			t.Fatalf("failed to read server log: %v", err)
		}
		logStr := string(logContent)
		if !strings.Contains(logStr, "Shutting down") && !strings.Contains(logStr, "shutting down") {
			t.Error("server log missing shutdown message")
		}
	})
}

// buildAndRunProdApp builds and runs a generated app natively for production template testing.
// Returns the server command and log file path.
func buildAndRunProdApp(t *testing.T, appDir string, port int) (*exec.Cmd, string) {
	t.Helper()

	// Inject components for test
	injectComponentsForTest(t, appDir)

	// Write embedded client library
	writeEmbeddedClientLibrary(t, appDir)

	// Run sqlc generate (skip failures for empty query files in bare apps)
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err == nil {
		sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
		sqlcCmd.Dir = appDir
		sqlcCmd.Env = append(os.Environ(), "GOWORK=off")
		if output, err := sqlcCmd.CombinedOutput(); err != nil {
			t.Logf("sqlc generate (non-fatal, no queries yet): %v\n%s", err, output)
		}
	}

	// Run migrations
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Logf("Warning: migration up failed (may be empty): %v", err)
	}

	// Build the app binary
	binaryPath := filepath.Join(appDir, "server")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/...")
	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "GOWORK=off", "CGO_ENABLED=1")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build app: %v\nOutput: %s", err, output)
	}
	t.Log("App built successfully")

	// Start the server
	portStr := fmt.Sprintf("%d", port)
	serverCmd := exec.Command(binaryPath)
	serverCmd.Dir = appDir
	serverCmd.Env = append(os.Environ(),
		"PORT="+portStr,
		"TEST_MODE=1",
	)

	serverLogPath := filepath.Join(appDir, "server.log")
	serverLogFile, err := os.Create(serverLogPath)
	if err != nil {
		t.Fatalf("Failed to create server log: %v", err)
	}
	serverCmd.Stdout = serverLogFile
	serverCmd.Stderr = serverLogFile

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	waitForServer(t, serverURL+"/health/live", 15*time.Second)

	t.Cleanup(func() {
		if serverCmd.Process != nil {
			_ = serverCmd.Process.Kill()
			_ = serverCmd.Wait()
		}
		if serverLogFile != nil {
			serverLogFile.Close()
		}
	})

	return serverCmd, serverLogPath
}
