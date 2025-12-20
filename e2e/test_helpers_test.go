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

	"github.com/chromedp/chromedp"
	"github.com/livetemplate/lvt/commands"
	"github.com/livetemplate/lvt/internal/serve"
	e2etest "github.com/livetemplate/lvt/testing"
)

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

// DockerContainerHandle provides control over a running Docker container
type DockerContainerHandle struct {
	containerID string
	port        int
}

// Stop stops and removes the Docker container
func (h *DockerContainerHandle) Stop(t *testing.T) {
	t.Helper()
	if h.containerID == "" {
		return
	}

	t.Logf("Stopping Docker container %s...", h.containerID)

	// Stop container
	stopCmd := exec.Command("docker", "stop", h.containerID)
	if output, err := stopCmd.CombinedOutput(); err != nil {
		t.Logf("Warning: Failed to stop container: %v\nOutput: %s", err, output)
	}

	// Remove container
	rmCmd := exec.Command("docker", "rm", h.containerID)
	if output, err := rmCmd.CombinedOutput(); err != nil {
		t.Logf("Warning: Failed to remove container: %v\nOutput: %s", err, output)
	} else {
		t.Logf("✅ Container %s stopped and removed", h.containerID)
	}
}

// enableDevMode enables development mode for the test app by writing .lvtrc config
// In DevMode, the app serves the local client library instead of using CDN
func enableDevMode(t *testing.T, appDir string) {
	t.Helper()
	lvtrcPath := filepath.Join(appDir, ".lvtrc")
	lvtrcContent := "dev_mode=true\n"
	if err := os.WriteFile(lvtrcPath, []byte(lvtrcContent), 0644); err != nil {
		t.Fatalf("Failed to write .lvtrc: %v", err)
	}
	t.Log("✅ Enabled DevMode for test app")
}

// writeEmbeddedClientLibrary writes the embedded client library to the app directory
// This allows Docker-based e2e tests to serve it locally instead of using CDN
func writeEmbeddedClientLibrary(t *testing.T, appDir string) {
	t.Helper()
	clientPath := filepath.Join(appDir, "livetemplate-client.js")
	if err := os.WriteFile(clientPath, e2etest.GetClientLibraryJS(), 0644); err != nil {
		t.Fatalf("Failed to write client library: %v", err)
	}
	t.Logf("✅ Wrote embedded client library to %s (%d bytes)", clientPath, len(e2etest.GetClientLibraryJS()))
}

// setupLocalClientLibrary configures the test app to use the embedded local client library
// Call this before building Docker images for Docker-based e2e tests
func setupLocalClientLibrary(t *testing.T, appDir string) {
	t.Helper()
	enableDevMode(t, appDir)
	writeEmbeddedClientLibrary(t, appDir)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// buildDockerImage builds a Docker image from the app directory
func buildDockerImage(t *testing.T, appDir, imageName string) {
	t.Helper()
	t.Logf("Building Docker image: %s", imageName)

	// Ensure base image exists
	buildBaseImage(t)

	// Create Dockerfile that builds on base
	// The base image has common dependencies cached, so go mod tidy will be fast
	dockerfile := `FROM lvt-base:latest

# Copy app-specific code
COPY . .

# Tidy and download dependencies using cache mount to avoid re-downloading across builds
# This shares the Go module cache across all parallel Docker builds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy && go mod download

# Generate database code if sqlc.yaml exists
RUN if [ -f database/sqlc.yaml ]; then \
      echo "Running sqlc generate..." && \
      sqlc generate -f database/sqlc.yaml; \
    fi

# Build the app
# Auto-detect if main.go is in root (simple kit) or cmd/ (multi kit)
RUN if [ -f main.go ]; then \
      CGO_ENABLED=1 go build -o server .; \
    else \
      CGO_ENABLED=1 go build -o server ./cmd/*; \
    fi

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite-libs
WORKDIR /app
COPY --from=0 /app/server /app/server
# Copy directories that might exist (use shell to handle missing dirs)
COPY --from=0 /app /app/
# Clean up build artifacts we don't need at runtime
RUN rm -rf /app/cmd /app/go.mod /app/go.sum /app/README.md /app/.git* 2>/dev/null || true
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./server"]
`

	dockerfilePath := filepath.Join(appDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		t.Fatalf("Failed to write Dockerfile: %v", err)
	}

	// Build only the app layer (fast, ~5-10 seconds)
	// Enable BuildKit for cache mount support
	buildCmd := exec.Command("docker", "build", "-t", imageName, ".")
	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Docker build failed: %v\nOutput: %s", err, output)
	}

	t.Log("✅ Docker image built successfully")
}

// runDockerContainer starts a Docker container and returns a handle
func runDockerContainer(t *testing.T, imageName string, port int) *DockerContainerHandle {
	t.Helper()
	t.Logf("Starting Docker container from %s on port %d", imageName, port)

	containerID := fmt.Sprintf("lvt-test-%d-%d", time.Now().Unix(), port)

	runCmd := exec.Command("docker", "run", "-d",
		"--name", containerID,
		"-p", fmt.Sprintf("%d:8080", port),
		imageName)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Docker run failed: %v\nOutput: %s", err, output)
	}

	handle := &DockerContainerHandle{
		containerID: containerID,
		port:        port,
	}

	// Register cleanup
	t.Cleanup(func() {
		handle.Stop(t)
	})

	t.Logf("✅ Container started: %s", containerID)
	return handle
}

// ensureDockerfile creates a Dockerfile if it doesn't exist
func ensureDockerfile(t *testing.T, appDir string) {
	t.Helper()

	dockerfilePath := filepath.Join(appDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err == nil {
		return // Already exists
	}

	t.Log("Generating Dockerfile...")

	// Use the multi-stage Dockerfile pattern from testing/deployment.go
	dockerfile := `# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev curl

# Install sqlc for database code generation
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ]; then SQLC_ARCH="arm64"; else SQLC_ARCH="amd64"; fi && \
    curl -L https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_${SQLC_ARCH}.tar.gz | tar -xz -C /usr/local/bin

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy after copying source (in case source files affect dependencies)
RUN go mod tidy

# Generate sqlc models if sqlc.yaml exists (multi kit with database)
RUN if [ -f database/sqlc.yaml ]; then \
      echo "Running sqlc generate..." && \
      sqlc generate -f database/sqlc.yaml; \
    fi

# Build binary with CGO enabled for SQLite
# Auto-detect if main.go is in root (simple kit) or cmd/ (multi kit)
RUN if [ -f main.go ]; then \
      CGO_ENABLED=1 GOOS=linux go build -o main .; \
    else \
      CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/*; \
    fi

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy all source files needed at runtime
COPY --from=builder /app .

# Clean up build artifacts we don't need at runtime
RUN rm -rf /app/cmd /app/go.mod /app/go.sum /app/README.md /app/.git* 2>/dev/null || true

# Create data directory for SQLite
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./main"]
`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		t.Fatalf("Failed to write Dockerfile: %v", err)
	}

	t.Log("✅ Dockerfile generated")
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

// verifyNoTemplateErrors checks that the page has no template errors
func verifyNoTemplateErrors(t *testing.T, ctx context.Context, url string) {
	t.Helper()

	var bodyText string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
		chromedp.Text("body", &bodyText, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	// Check for common template error patterns
	errorPatterns := []string{
		"template:",
		"<no value>",
		"{{.",
		"executing template",
		"parse error",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(bodyText, pattern) {
			t.Errorf("❌ Template error found on page: contains %q", pattern)
		}
	}
}

// verifyWebSocketConnected checks that WebSocket connection is established
func verifyWebSocketConnected(t *testing.T, ctx context.Context, url string) {
	t.Helper()

	var wsConnected bool
	var wsURL string
	var wsReadyState int

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		e2etest.WaitForWebSocketReady(30*time.Second), // Increased for CDN loading + WebSocket init
		chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
		chromedp.Evaluate(`window.liveTemplateClient && window.liveTemplateClient.ws ? window.liveTemplateClient.ws.url : null`, &wsURL),
		chromedp.Evaluate(`window.liveTemplateClient && window.liveTemplateClient.ws ? window.liveTemplateClient.ws.readyState : -1`, &wsReadyState),
		chromedp.Evaluate(`(() => {
			return window.liveTemplateClient &&
			       window.liveTemplateClient.ws &&
			       window.liveTemplateClient.ws.readyState === WebSocket.OPEN;
		})()`, &wsConnected),
	)
	if err != nil {
		t.Fatalf("Failed to check WebSocket: %v", err)
	}

	t.Logf("WebSocket URL: %s, ReadyState: %d (1=OPEN)", wsURL, wsReadyState)

	if !wsConnected {
		t.Errorf("❌ WebSocket not connected (readyState: %d)", wsReadyState)
	} else {
		t.Log("✅ WebSocket connected")
	}
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

// buildAndRunNative builds the app natively and starts it on the specified port
// This is much faster than Docker build (~5s vs ~245s)
// Returns the server process command
func buildAndRunNative(t *testing.T, appDir string, port int) *exec.Cmd {
	t.Helper()

	t.Log("Step 6: Building app natively (fast path)...")

	// Write embedded client library (DevMode should already be enabled)
	writeEmbeddedClientLibrary(t, appDir)

	// Run sqlc generate if sqlc.yaml exists
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err == nil {
		t.Log("Running sqlc generate...")
		sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
		sqlcCmd.Dir = appDir
		sqlcCmd.Env = append(os.Environ(), "GOWORK=off")
		if output, err := sqlcCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to run sqlc generate: %v\nOutput: %s", err, output)
		}
		t.Log("✅ sqlc generate completed")
	}

	// Build the app
	// Check if simple kit (main.go in root) or multi kit (main.go in cmd/)
	binaryPath := filepath.Join(appDir, "server")
	t.Log("Building binary...")

	var buildCmd *exec.Cmd
	if _, err := os.Stat(filepath.Join(appDir, "main.go")); err == nil {
		// Simple kit - main.go in root
		buildCmd = exec.Command("go", "build", "-o", binaryPath, ".")
	} else {
		// Multi kit - main.go in cmd/
		buildCmd = exec.Command("go", "build", "-o", binaryPath, "./cmd/...")
	}

	buildCmd.Dir = appDir
	buildCmd.Env = append(os.Environ(), "GOWORK=off", "CGO_ENABLED=1")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build app: %v\nOutput: %s", err, output)
	}
	t.Log("✅ App built successfully")

	// Step 7: Start the app
	t.Log("Step 7: Starting app natively...")
	portStr := fmt.Sprintf("%d", port)
	serverCmd := exec.Command(binaryPath)
	serverCmd.Dir = appDir
	serverCmd.Env = append(os.Environ(),
		"PORT="+portStr,
		"LVT_DEV_MODE=true",
	)

	// Redirect output to file for debugging
	serverLogPath := filepath.Join(appDir, "server.log")
	serverLogFile, err := os.Create(serverLogPath)
	if err != nil {
		t.Fatalf("Failed to create server log file: %v", err)
	}
	serverCmd.Stdout = serverLogFile
	serverCmd.Stderr = serverLogFile
	t.Logf("📝 Server logs will be written to: %s", serverLogPath)

	// Start the server
	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	serverURL := fmt.Sprintf("http://localhost:%d", port)
	ready := false
	var lastErr error
	consecutiveSuccesses := 0
	const requiredSuccesses = 2

	for i := 0; i < 50; i++ {
		resp, err := http.Get(serverURL)
		if err == nil {
			if resp.StatusCode == 200 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				bodyStr := string(body)
				if strings.Contains(bodyStr, "<!DOCTYPE html>") || strings.Contains(bodyStr, "<html") {
					consecutiveSuccesses++
					if consecutiveSuccesses >= requiredSuccesses {
						ready = true
						break
					}
				}
			} else {
				resp.Body.Close()
				consecutiveSuccesses = 0
			}
		} else {
			lastErr = err
			consecutiveSuccesses = 0
		}
		time.Sleep(200 * time.Millisecond)
	}

	if !ready {
		_ = serverCmd.Process.Kill()
		t.Fatalf("❌ Server failed to respond within 10 seconds. Last error: %v", lastErr)
	}

	t.Logf("✅ App running on http://localhost:%d", port)

	// Register cleanup
	t.Cleanup(func() {
		if serverCmd.Process != nil {
			t.Logf("Stopping native server (PID: %d)...", serverCmd.Process.Pid)
			if err := serverCmd.Process.Kill(); err != nil {
				t.Logf("Warning: Failed to kill server: %v", err)
			} else {
				t.Log("✅ Native server stopped")
			}
			_ = serverCmd.Wait()
		}

		// Close log file and print timing logs
		if serverLogFile != nil {
			serverLogFile.Close()

			// Read and print debug logs ([TIMING], [PUMP], [SEND])
			if content, err := os.ReadFile(serverLogPath); err == nil {
				lines := strings.Split(string(content), "\n")
				debugLines := []string{}
				for _, line := range lines {
					if strings.Contains(line, "[TIMING]") || strings.Contains(line, "[PUMP]") || strings.Contains(line, "[SEND]") {
						debugLines = append(debugLines, line)
					}
				}
				if len(debugLines) > 0 {
					t.Log("📊 DEBUG LOGS ([TIMING], [PUMP], [SEND]):")
					for _, line := range debugLines {
						t.Log(line)
					}
				}
			}
		}
	})

	return serverCmd
}
