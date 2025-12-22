//go:build browser

package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

// Shared test resources that persist across all tests
var (
	sharedChromePort  int = 9222
	sharedTestApp     string
	sharedTestAppDir  string
	sharedSetupOnce   sync.Once
	sharedCleanupOnce sync.Once
	sharedSetupError  error
)

// TestMain sets up shared resources before running tests and cleans up after
func TestMain(m *testing.M) {
	// Cleanup any leftover containers from previous runs
	// This is safe to run even if Docker is not available - it will just log a warning
	cleanupChromeContainers()

	// Run tests
	code := m.Run()

	// Cleanup Chrome pool if it was initialized
	chromePoolMu.Lock()
	if chromePool != nil {
		chromePool.Cleanup()
	}
	chromePoolMu.Unlock()

	// Final cleanup of any remaining containers
	cleanupChromeContainers()

	os.Exit(code)
}

// setupSharedResources initializes shared Chrome container and pre-built test app
func setupSharedResources() error {
	var setupErr error
	sharedSetupOnce.Do(func() {
		log.Println("🚀 Setting up shared test resources...")

		// 1. Start shared Docker Chrome container
		log.Printf("Starting shared Docker Chrome on port %d...", sharedChromePort)
		if err := e2etest.StartDockerChrome(&testing.T{}, sharedChromePort); err != nil {
			setupErr = fmt.Errorf("failed to start shared Chrome container: %w", err)
			return
		}

		// Wait a bit for Chrome to be fully ready
		time.Sleep(2 * time.Second)
		log.Println("✅ Shared Docker Chrome started")

		// 2. Pre-compile a test application for reuse
		log.Println("Pre-compiling test application...")
		tmpDir, err := os.MkdirTemp("", "lvt-e2e-shared-*")
		if err != nil {
			setupErr = fmt.Errorf("failed to create temp dir: %w", err)
			return
		}
		sharedTestAppDir = tmpDir

		// Create a standard test app
		appDir := filepath.Join(tmpDir, "sharedapp")

		// Use createTestApp if we can, but do it manually for now
		if err := runLvtCommand(&testing.T{}, tmpDir, "new", "sharedapp", "--dev"); err != nil {
			setupErr = fmt.Errorf("failed to create shared app: %w", err)
			return
		}

		// Generate a posts resource
		if err := runLvtCommand(&testing.T{}, appDir, "gen", "resource", "posts", "title", "content:text", "published:bool"); err != nil {
			setupErr = fmt.Errorf("failed to generate posts resource: %w", err)
			return
		}

		// Run migrations
		if err := runLvtCommand(&testing.T{}, appDir, "migration", "up"); err != nil {
			setupErr = fmt.Errorf("failed to run migrations: %w", err)
			return
		}

		// Run go mod tidy and vendor in a Docker container to isolate background processes
		// This prevents background processes from affecting the test process
		// Using the published version of livetemplate (no replace directive)
		log.Println("Running go mod tidy and vendor in Docker container...")
		dockerCmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s:/app", appDir),
			"-w", "/app",
			"-e", "GOWORK=off",
			"golang:1.25",
			"sh", "-c", "go mod tidy && go mod vendor")

		dockerOutput, err := dockerCmd.CombinedOutput()
		if err != nil {
			setupErr = fmt.Errorf("failed to run go mod tidy/vendor in Docker: %w\nOutput: %s", err, string(dockerOutput))
			return
		}
		log.Println("✅ go mod tidy and vendor completed in Docker")

		// Run sqlc generate in Docker to isolate background processes
		log.Println("Running sqlc generate...")
		sqlcDockerCmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s:/app", appDir),
			"-w", "/app",
			"-e", "GOWORK=off",
			"golang:1.25",
			"go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", "database/sqlc.yaml")

		sqlcOutput, err := sqlcDockerCmd.CombinedOutput()
		if err != nil {
			setupErr = fmt.Errorf("failed to run sqlc in Docker: %w\nOutput: %s", err, string(sqlcOutput))
			return
		}
		log.Println("✅ sqlc generate completed")

		// Build the app using vendored dependencies
		appBinary := filepath.Join(appDir, "sharedapp")
		buildCmd := exec.Command("go", "build", "-mod=vendor", "-o", appBinary, "./cmd/sharedapp")
		buildCmd.Dir = appDir
		if output, err := buildCmd.CombinedOutput(); err != nil {
			setupErr = fmt.Errorf("failed to build shared app: %w\n%s", err, output)
			return
		}

		sharedTestApp = appBinary
		log.Printf("✅ Pre-compiled test application: %s", sharedTestApp)

		sharedSetupError = nil
	})

	return setupErr
}

// cleanupSharedResources tears down shared resources
func cleanupSharedResources() {
	sharedCleanupOnce.Do(func() {
		log.Println("🧹 Cleaning up shared test resources...")

		// Note: go mod tidy spawns background processes that may cause I/O wait warnings
		// These are harmless and don't affect test results

		// Stop shared Chrome container
		log.Println("Stopping shared Docker Chrome...")
		e2etest.StopDockerChrome(nil, sharedChromePort)
		log.Println("✅ Shared Docker Chrome stopped")

		// Clean up shared test app directory
		if sharedTestAppDir != "" {
			log.Printf("Cleaning up shared test app directory: %s", sharedTestAppDir)
			os.RemoveAll(sharedTestAppDir)
		}

		log.Println("✅ Shared test resources cleaned up")
	})
}

// getIsolatedChromeContext creates a dedicated Chrome container per test
// This enables true parallel execution without contention on shared Chrome instance
func getIsolatedChromeContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	// Allocate unique port for this test
	port, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("Failed to allocate port: %v", err)
	}

	// Start dedicated Chrome container for this test
	if err := e2etest.StartDockerChrome(t, port); err != nil {
		t.Fatalf("Failed to start isolated Chrome container: %v", err)
	}

	// Create allocator context for isolated Chrome
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(),
		fmt.Sprintf("http://localhost:%d", port))

	// Create browser context
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Return combined cancel function that cleans up Chrome container
	combinedCancel := func() {
		cancel()
		allocCancel()
		e2etest.StopDockerChrome(t, port)
	}

	return ctx, combinedCancel
}


// Unused: Kept for potential future use
// getSharedTestApp returns the path to the pre-compiled test application
// func getSharedTestApp(t *testing.T) string {
// 	t.Helper()
//
// 	// Ensure shared resources are setup
// 	if err := setupSharedResources(); err != nil {
// 		t.Fatalf("Failed to setup shared resources: %v", err)
// 	}
//
// 	if sharedSetupError != nil {
// 		t.Fatalf("Shared setup failed: %v", sharedSetupError)
// 	}
//
// 	return sharedTestApp
// }

// Unused: Kept for potential future use
// getSharedTestAppDir returns the directory containing the pre-built test app
// func getSharedTestAppDir(t *testing.T) string {
// 	t.Helper()
//
// 	if err := setupSharedResources(); err != nil {
// 		t.Fatalf("Failed to setup shared resources: %v", err)
// 	}
//
// 	if sharedSetupError != nil {
// 		t.Fatalf("Shared setup failed: %v", sharedSetupError)
// 	}
//
// 	return sharedTestAppDir
// }
