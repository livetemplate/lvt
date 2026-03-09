//go:build browser

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestTemplateEdgeCases tests generation combinations:
// 1. App with no resources
// 2. App with only views
// 3. App with resources but no auth
// 4. App with auth but no resources
func TestTemplateEdgeCases(t *testing.T) {
	t.Run("NoResources", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		appDir := createTestApp(t, tmpDir, "bareapp", &AppOptions{Kit: "multi"})
		verifyAppBuildsAndRuns(t, appDir, false)
	})

	t.Run("OnlyViews", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		appDir := createTestApp(t, tmpDir, "viewapp", &AppOptions{Kit: "multi"})

		// Generate a view (no database)
		if err := runLvtCommand(t, appDir, "gen", "view", "dashboard"); err != nil {
			t.Fatalf("Failed to generate view: %v", err)
		}

		verifyAppBuildsAndRuns(t, appDir, false)
	})

	t.Run("ResourceNoAuth", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		appDir := createTestApp(t, tmpDir, "resapp", &AppOptions{Kit: "multi"})

		// Generate resource
		if err := runLvtCommand(t, appDir, "gen", "resource", "posts", "title", "content:text"); err != nil {
			t.Fatalf("Failed to generate resource: %v", err)
		}

		verifyAppBuildsAndRuns(t, appDir, true)
	})

	t.Run("AuthNoResources", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		appDir := createTestApp(t, tmpDir, "authapp", &AppOptions{Kit: "multi"})

		// Generate auth (with no resources to protect)
		if err := runLvtCommand(t, appDir, "gen", "auth"); err != nil {
			t.Fatalf("Failed to generate auth: %v", err)
		}

		verifyAppBuildsAndRuns(t, appDir, true)
	})
}

// verifyAppBuildsAndRuns builds and starts a generated app, verifying it responds to HTTP.
// If hasQueries is true, runs sqlc generate first.
func verifyAppBuildsAndRuns(t *testing.T, appDir string, hasQueries bool) {
	t.Helper()

	// Inject components
	injectComponentsForTest(t, appDir)

	// Write embedded client library
	writeEmbeddedClientLibrary(t, appDir)

	// Run sqlc if needed
	if hasQueries {
		sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
		if _, err := os.Stat(sqlcPath); err == nil {
			sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
			sqlcCmd.Dir = appDir
			sqlcCmd.Env = append(os.Environ(), "GOWORK=off")
			if output, err := sqlcCmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed to run sqlc generate: %v\nOutput: %s", err, output)
			}
		}
	}

	// Run migrations if available
	_ = runLvtCommand(t, appDir, "migration", "up")

	// Build the app
	binaryPath := filepath.Join(appDir, "server")
	var buildCmd *exec.Cmd
	if _, err := os.Stat(filepath.Join(appDir, "main.go")); err == nil {
		buildCmd = exec.Command("go", "build", "-o", binaryPath, ".")
	} else {
		buildCmd = exec.Command("go", "build", "-o", binaryPath, "./cmd/...")
	}
	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "GOWORK=off", "CGO_ENABLED=1")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build app: %v\nOutput: %s", err, output)
	}
	t.Log("Build succeeded")

	// Start and verify server responds
	port := allocateTestPort()
	serverCmd := exec.Command(binaryPath)
	serverCmd.Dir = appDir
	serverCmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		"TEST_MODE=1",
	)
	serverCmd.Stdout = os.Stderr
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	t.Cleanup(func() {
		if serverCmd.Process != nil {
			_ = serverCmd.Process.Kill()
			_ = serverCmd.Wait()
		}
	})

	serverURL := fmt.Sprintf("http://localhost:%d", port)
	waitForServer(t, serverURL+"/health/live", 15*time.Second)
	t.Logf("Server running at %s", serverURL)
}
