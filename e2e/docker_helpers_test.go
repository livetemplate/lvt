//go:build deployment || browser

package e2e

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	e2etest "github.com/livetemplate/lvt/testing"
)

var (
	baseImageBuilt bool
	baseImageMutex sync.Mutex
)

// buildBaseImage builds the shared base Docker image once
func buildBaseImage(t *testing.T) {
	baseImageMutex.Lock()
	defer baseImageMutex.Unlock()

	if baseImageBuilt {
		return
	}

	// Check if Docker is available
	if _, err := exec.Command("docker", "version").CombinedOutput(); err != nil {
		t.Skip("Docker not available, skipping test that requires Docker image build")
	}

	t.Log("Building base Docker image (one-time setup)...")

	// Get the directory where this test file lives
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	e2eDir := filepath.Dir(filename)

	cmd := exec.Command("docker", "build",
		"-f", "Dockerfile.base",
		"-t", "lvt-base:latest",
		".",
	)
	cmd.Dir = e2eDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build base image: %v\nOutput: %s", err, output)
	}

	t.Log("Base image built successfully")
	baseImageBuilt = true
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
		t.Logf("Container %s stopped and removed", h.containerID)
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
	t.Log("Enabled DevMode for test app")
}

// writeEmbeddedClientLibrary writes the embedded client library to the app directory
// This allows Docker-based e2e tests to serve it locally instead of using CDN
func writeEmbeddedClientLibrary(t *testing.T, appDir string) {
	t.Helper()
	clientPath := filepath.Join(appDir, "livetemplate-client.js")
	if err := os.WriteFile(clientPath, e2etest.GetClientLibraryJS(), 0644); err != nil {
		t.Fatalf("Failed to write client library: %v", err)
	}
	t.Logf("Wrote embedded client library to %s (%d bytes)", clientPath, len(e2etest.GetClientLibraryJS()))
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

	t.Log("Docker image built successfully")
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

	t.Logf("Container started: %s", containerID)
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
FROM golang:1.26-alpine AS builder

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

	t.Log("Dockerfile generated")
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
		t.Log("sqlc generate completed")
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
	t.Log("App built successfully")

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
	t.Logf("Server logs will be written to: %s", serverLogPath)

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
		t.Fatalf("Server failed to respond within 10 seconds. Last error: %v", lastErr)
	}

	t.Logf("App running on http://localhost:%d", port)

	// Register cleanup
	t.Cleanup(func() {
		if serverCmd.Process != nil {
			t.Logf("Stopping native server (PID: %d)...", serverCmd.Process.Pid)
			if err := serverCmd.Process.Kill(); err != nil {
				t.Logf("Warning: Failed to kill server: %v", err)
			} else {
				t.Log("Native server stopped")
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
					t.Log("DEBUG LOGS ([TIMING], [PUMP], [SEND]):")
					for _, line := range debugLines {
						t.Log(line)
					}
				}
			}
		}
	})

	return serverCmd
}
