//go:build browser

package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestSkillPipeline validates the complete LLM skill pipeline:
// lvt new → gen resource → gen resource (with references) → gen auth → seed → build → run → verify
func TestSkillPipeline(t *testing.T) {
	tmpDir := t.TempDir()

	// Step 1: Create blog app
	t.Log("Step 1: Creating blog app with multi kit and tailwind...")
	appDir := createTestApp(t, tmpDir, "blogapp", &AppOptions{Kit: "multi"})

	// Step 2: Generate posts resource
	t.Log("Step 2: Generating posts resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "posts", "title", "content:text", "published:bool"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}

	// Step 3: Generate comments resource (with reference to posts)
	t.Log("Step 3: Generating comments resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "comments", "post_id:references:posts", "author", "text"); err != nil {
		t.Fatalf("Failed to generate comments: %v", err)
	}

	// Step 4: Generate auth
	t.Log("Step 4: Generating auth...")
	if err := runLvtCommand(t, appDir, "gen", "auth"); err != nil {
		t.Fatalf("Failed to generate auth: %v", err)
	}

	// Step 5: Run migrations
	t.Log("Step 5: Running migrations...")
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Step 6: Seed posts
	t.Log("Step 6: Seeding posts...")
	if err := runLvtCommand(t, appDir, "seed", "posts", "--count", "5"); err != nil {
		t.Logf("Warning: seed failed (may not be critical): %v", err)
	}

	// Step 7: Generate env
	t.Log("Step 7: Verifying .env.example exists...")
	if _, err := os.Stat(filepath.Join(appDir, ".env.example")); err != nil {
		t.Fatalf(".env.example not found: %v", err)
	}

	// Step 8: Inject components and build
	t.Log("Step 8: Building app...")
	injectComponentsForTest(t, appDir)
	writeEmbeddedClientLibrary(t, appDir)

	// Run sqlc
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	sqlcCmd := exec.Command("go", "run", sqlcPackage, "generate", "-f", sqlcPath)
	sqlcCmd.Dir = appDir
	sqlcCmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := sqlcCmd.CombinedOutput(); err != nil {
		t.Fatalf("sqlc generate failed: %v\nOutput: %s", err, output)
	}

	binaryPath := filepath.Join(appDir, "server")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/...")
	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "GOWORK=off", "CGO_ENABLED=1")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\nOutput: %s", err, output)
	}
	t.Log("Build succeeded")

	// Step 9: Start server
	t.Log("Step 9: Starting server...")
	port := allocateTestPort()
	serverCmd := exec.Command(binaryPath)
	serverCmd.Dir = appDir
	serverCmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
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
	t.Cleanup(func() {
		if serverCmd.Process != nil {
			_ = serverCmd.Process.Kill()
			_ = serverCmd.Wait()
		}
		serverLogFile.Close()
	})

	serverURL := fmt.Sprintf("http://localhost:%d", port)
	waitForServer(t, serverURL+"/health/live", 15*time.Second)

	// Step 10: Verify via HTTP
	t.Run("HealthEndpoints", func(t *testing.T) {
		for _, path := range []string{"/health/live", "/health/ready"} {
			resp, err := http.Get(serverURL + path)
			if err != nil {
				t.Fatalf("GET %s failed: %v", path, err)
			}
			if resp.StatusCode != 200 {
				t.Errorf("%s returned %d", path, resp.StatusCode)
			}
			resp.Body.Close()
		}
	})

	t.Run("SecurityHeaders", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
			t.Error("missing X-Content-Type-Options: nosniff")
		}
	})

	t.Run("HomePageLoads", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "<!DOCTYPE html>") && !strings.Contains(string(body), "<html") {
			t.Error("home page does not return HTML")
		}
	})

	// Step 11: Verify via chromedp (using pooled Docker Chrome)
	t.Run("BrowserRenders", func(t *testing.T) {
		ctx, _, cleanup := GetPooledChrome(t)
		defer cleanup()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		testURL := getTestURL(port)
		var html string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(testURL),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.OuterHTML("html", &html),
		); err != nil {
			t.Fatalf("chromedp failed: %v", err)
		}

		if !strings.Contains(html, "blogapp") {
			t.Error("page does not contain app name")
		}
	})

	t.Run("AuthPageLoads", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/auth")
		if err != nil {
			t.Fatalf("GET /auth failed: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		// Auth page should return HTML (either directly or via redirect/WebSocket)
		if resp.StatusCode != 200 && resp.StatusCode != 101 {
			t.Errorf("/auth returned unexpected status: %d\nBody: %s", resp.StatusCode, string(body[:min(len(body), 500)]))
		}
	})

	t.Run("ResourceRoutesExist", func(t *testing.T) {
		for _, path := range []string{"/posts", "/comments"} {
			resp, err := http.Get(serverURL + path)
			if err != nil {
				t.Fatalf("GET %s failed: %v", path, err)
			}
			// Protected routes may redirect (302) or require auth (401/403)
			// or return 200 if no auth middleware wraps them yet
			if resp.StatusCode >= 500 {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("%s returned server error %d: %s", path, resp.StatusCode, string(body[:min(len(body), 500)]))
			}
			resp.Body.Close()
		}
	})

	// Step 12: Check server logs for errors
	t.Run("NoServerErrors", func(t *testing.T) {
		logContent, err := os.ReadFile(serverLogPath)
		if err != nil {
			t.Fatalf("Failed to read server log: %v", err)
		}
		logStr := string(logContent)
		if strings.Contains(logStr, "\"level\":\"ERROR\"") {
			// Filter out expected errors (like WebSocket upgrade on non-WebSocket requests)
			lines := strings.Split(logStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "\"level\":\"ERROR\"") && !strings.Contains(line, "websocket") {
					t.Errorf("server error in logs: %s", line)
				}
			}
		}
	})
}
