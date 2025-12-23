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
	"sync"
	"testing"
	"time"

	"github.com/livetemplate/lvt/commands"
	"github.com/livetemplate/lvt/internal/serve"
)

// Port allocation for parallel tests
var (
	portMutex    sync.Mutex
	nextTestPort int = 8800
)

// getTimeout returns local or CI timeout based on environment
func getTimeout(envVar string, localDefault, ciDefault time.Duration) time.Duration {
	if os.Getenv("CI") == "true" {
		return ciDefault
	}
	if val := os.Getenv(envVar); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return localDefault
}

// allocateTestPort returns a unique port for a test
func allocateTestPort() int {
	portMutex.Lock()
	defer portMutex.Unlock()

	port := nextTestPort
	nextTestPort++
	return port
}

// runSqlc runs sqlc generate in the app directory
func runSqlc(t *testing.T, appDir string) {
	t.Helper()
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err == nil {
		t.Log("Running sqlc generate...")
		sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
		sqlcCmd.Dir = appDir
		sqlcCmd.Env = append(os.Environ(), "GOWORK=off")
		if output, err := sqlcCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to run sqlc generate: %v\nOutput: %s", err, output)
		}
		t.Log("sqlc generate completed")
	}
}

// chdirMutex protects os.Chdir calls in parallel tests
// os.Chdir affects the entire process, so we need to serialize these operations
var chdirMutex sync.Mutex

// AppOptions contains options for creating a test app
type AppOptions struct {
	Kit           string // Kit name (multi, single, simple)
	Module        string // Go module name
	SkipGoModTidy bool   // Skip go mod tidy (for Docker-based tests that run it inside Docker)
}

// runLvtCommand executes an lvt command directly by calling the command functions
// This is much faster than shelling out and avoids working directory issues
func runLvtCommand(t *testing.T, workDir string, args ...string) error {
	t.Helper()
	_, err := runLvtCommandWithOutput(t, workDir, args...)
	return err
}

// runLvtCommandWithOutput executes an lvt command and returns its output
func runLvtCommandWithOutput(t *testing.T, workDir string, args ...string) (string, error) {
	t.Helper()
	t.Logf("Running: lvt %s", strings.Join(args, " "))

	if len(args) == 0 {
		return "", fmt.Errorf("no command specified")
	}

	// Lock to serialize directory changes across parallel tests
	// os.Chdir affects the entire process, not just the current goroutine
	chdirMutex.Lock()
	defer chdirMutex.Unlock()

	// Save and restore working directory
	origDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			return "", fmt.Errorf("failed to change directory to %s: %w", workDir, err)
		}
	}

	// Save and restore stdout/stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Capture output to buffer
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	var outputBuf strings.Builder
	outputDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(&outputBuf, r)
		close(outputDone)
	}()

	// Execute the command directly
	command := args[0]
	cmdArgs := args[1:]

	var cmdErr error
	switch command {
	case "new":
		cmdErr = commands.New(cmdArgs)
	case "gen":
		cmdErr = commands.Gen(cmdArgs)
	case "migration":
		cmdErr = commands.Migration(cmdArgs)
	case "kits", "kit":
		cmdErr = commands.Kits(cmdArgs)
	case "serve":
		cmdErr = commands.Serve(cmdArgs)
	case "resource", "res":
		cmdErr = commands.Resource(cmdArgs)
	case "seed":
		cmdErr = commands.Seed(cmdArgs)
	case "parse":
		cmdErr = commands.Parse(cmdArgs)
	default:
		cmdErr = fmt.Errorf("unknown command: %s", command)
	}

	// Restore stdout/stderr FIRST to prevent inheritance by any late spawned processes
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Then close the write end of the pipe
	w.Close()

	// Close the read end explicitly after reading is done to prevent lingering
	<-outputDone
	r.Close()

	output := outputBuf.String()

	if cmdErr != nil {
		return output, fmt.Errorf("command failed: lvt %s: %w", strings.Join(args, " "), cmdErr)
	}

	return output, nil
}

// ServerHandle provides control over a running test server
type ServerHandle struct {
	server  interface{ Shutdown() error }
	cancel  context.CancelFunc
	errChan chan error
}

// Shutdown stops the server gracefully
func (h *ServerHandle) Shutdown() error {
	if h.cancel != nil {
		h.cancel()
	}
	if h.server != nil {
		return h.server.Shutdown()
	}
	return nil
}

// Wait waits for the server to finish and returns any error
func (h *ServerHandle) Wait() error {
	return <-h.errChan
}

// startServeInBackground starts a development server in the background using a context
// Returns a handle that can be used to shut down the server via context cancellation
func startServeInBackground(t *testing.T, workDir string, args ...string) (*ServerHandle, error) {
	t.Helper()

	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Change to work directory if specified
	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			return nil, fmt.Errorf("failed to change directory to %s: %w", workDir, err)
		}
		defer func() { _ = os.Chdir(origDir) }()
	}

	// Parse serve arguments to create config
	config := serve.DefaultConfig()
	config.OpenBrowser = false // Never open browser in tests
	config.LiveReload = false  // Disable live reload by default in tests

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port", "-p":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--port requires a value")
			}
			var port int
			if _, err := fmt.Sscanf(args[i+1], "%d", &port); err != nil {
				return nil, fmt.Errorf("invalid port: %s", args[i+1])
			}
			config.Port = port
			i++
		case "--mode", "-m":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--mode requires a value")
			}
			config.Mode = serve.ServeMode(args[i+1])
			config.AutoDetect = false
			i++
		case "--no-browser":
			config.OpenBrowser = false
		case "--no-reload":
			config.LiveReload = false
		}
	}

	// Create server
	server, err := serve.NewServer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	// Create context with cancellation for server control
	ctx, cancel := context.WithCancel(context.Background())

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		// Change directory for the goroutine
		if workDir != "" {
			_ = os.Chdir(workDir)
		}

		// Start returns when server shuts down or errors
		if err := server.Start(ctx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	handle := &ServerHandle{
		server:  server,
		cancel:  cancel,
		errChan: errChan,
	}

	// Register cleanup handler to shut down server on test completion/failure
	t.Cleanup(func() {
		if err := handle.Shutdown(); err != nil {
			t.Logf("Warning: Failed to shutdown server: %v", err)
		} else {
			t.Log("✅ Server shutdown complete")
		}
	})

	return handle, nil
}

// createTestApp creates a new test application and sets it up for testing
func createTestApp(t *testing.T, tmpDir, appName string, opts *AppOptions) string {
	t.Helper()
	t.Logf("Creating test app: %s", appName)

	// Set defaults
	if opts == nil {
		opts = &AppOptions{
			Kit: "multi",
		}
	}

	// Build lvt new command
	args := []string{"new", appName}

	if opts.Kit != "" && opts.Kit != "multi" {
		args = append(args, "--kit", opts.Kit)
	}

	if opts.Module != "" {
		args = append(args, "--module", opts.Module)
	}

	// Create app
	if err := runLvtCommand(t, tmpDir, args...); err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	appDir := filepath.Join(tmpDir, appName)

	// Skip go mod tidy if requested (e.g., for Docker-based tests that run it inside Docker)
	// For non-Docker tests (lvt serve), go mod tidy is required to work properly
	if !opts.SkipGoModTidy {
		t.Log("Running go mod tidy...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = appDir
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to run go mod tidy: %v\nOutput: %s", err, output)
		}
		t.Log("✅ go mod tidy completed")
	} else {
		t.Log("⏭️  Skipping go mod tidy (will be run in Docker)")
	}

	// Register cleanup handler to remove app directory on test completion/failure
	// This is registered at the END to ensure it only runs after successful setup
	t.Cleanup(func() {
		if err := os.RemoveAll(appDir); err != nil {
			t.Logf("Warning: Failed to cleanup app directory %s: %v", appDir, err)
		} else {
			t.Logf("✅ Cleaned up app directory: %s", appDir)
		}
	})

	t.Log("✅ Test app created")
	return appDir
}

// waitForServer waits for the server to be ready and responding
func waitForServer(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	t.Logf("Waiting for server at %s...", url)

	deadline := time.Now().Add(timeout)
	var lastErr error
	consecutiveSuccesses := 0
	const requiredSuccesses = 2 // Require 2 consecutive successful responses for stability

	// Use exponential backoff for faster server detection
	retryDelay := 10 * time.Millisecond
	const maxRetryDelay = 100 * time.Millisecond

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			if resp.StatusCode == 200 {
				resp.Body.Close()
				consecutiveSuccesses++
				if consecutiveSuccesses >= requiredSuccesses {
					// Give server a bit more time to fully initialize WebSocket handlers
					time.Sleep(50 * time.Millisecond)
					t.Logf("✅ Server ready (verified with %d consecutive successful requests)", requiredSuccesses)
					return
				}
				// Reset delay on success
				retryDelay = 10 * time.Millisecond
			} else {
				resp.Body.Close()
				lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
				consecutiveSuccesses = 0
			}
		} else {
			lastErr = err
			consecutiveSuccesses = 0
		}

		time.Sleep(retryDelay)
		// Exponential backoff up to max
		retryDelay = retryDelay * 2
		if retryDelay > maxRetryDelay {
			retryDelay = maxRetryDelay
		}
	}

	t.Fatalf("❌ Server failed to respond within %v. Last error: %v", timeout, lastErr)
}

// readLvtrc reads and parses the .lvtrc file
func readLvtrc(t *testing.T, appDir string) (kit string) {
	t.Helper()

	lvtrcPath := filepath.Join(appDir, ".lvtrc")
	content, err := os.ReadFile(lvtrcPath)
	if err != nil {
		t.Fatalf("Failed to read .lvtrc: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "kit=") {
			kit = strings.TrimPrefix(line, "kit=")
		}
	}

	return kit
}
