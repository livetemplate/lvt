package testing

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	stdtesting "testing"
	"time"

	"github.com/livetemplate/lvt/testing/providers"
)

// DeploymentTest manages the lifecycle of a deployed test application
type DeploymentTest struct {
	T *stdtesting.T

	// App metadata
	Provider Provider
	AppName  string
	AppDir   string
	AppURL   string
	Region   string

	// Provider-specific clients (for debugging and inspection)
	DockerClient *providers.DockerClient

	// Cleanup tracking
	cleanupFuncs []func() error
	cleanupMu    sync.Mutex
	cleaned      bool
}

// DeploymentOptions configures deployment test setup
type DeploymentOptions struct {
	Provider Provider
	AppName  string // If empty, generates unique name
	AppDir   string // If empty, creates test app
	Region   string // Cloud region (defaults vary by provider)
	Kit      string // App kit: multi, single, simple (default: multi)

	// Feature flags
	WithAuth       bool
	WithLitestream bool
	WithS3Backup   bool

	// Resource options
	Resources []string // Resources to add (e.g., "posts title content")
}

// SetupDeployment creates a test app and prepares it for deployment
func SetupDeployment(t *stdtesting.T, opts *DeploymentOptions) *DeploymentTest {
	t.Helper()

	if opts == nil {
		opts = &DeploymentOptions{}
	}

	// Set defaults
	if opts.Provider == "" {
		opts.Provider = ProviderFly
	}
	if opts.Kit == "" {
		opts.Kit = "multi"
	}
	if opts.Region == "" {
		switch opts.Provider {
		case ProviderFly:
			opts.Region = "sjc" // San Jose
		case ProviderDigitalOcean:
			opts.Region = "nyc"
		default:
			opts.Region = "us-east-1"
		}
	}

	// Generate unique app name if not provided
	if opts.AppName == "" {
		opts.AppName = GenerateTestAppName("lvt-test")
	}

	// Validate app name
	if err := ValidateAppName(opts.AppName); err != nil {
		t.Fatalf("Invalid app name: %v", err)
	}

	dt := &DeploymentTest{
		T:        t,
		Provider: opts.Provider,
		AppName:  opts.AppName,
		Region:   opts.Region,
	}

	// Create test app if appDir not provided
	if opts.AppDir == "" {
		// Use working directory instead of system temp to avoid macOS permission dialogs
		// and to keep test artifacts accessible for debugging
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get working directory: %v", err)
		}

		testDir := filepath.Join(wd, ".test-deployments", t.Name())
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Register cleanup to remove test directory
		t.Cleanup(func() {
			if err := os.RemoveAll(testDir); err != nil {
				t.Logf("Warning: failed to clean up test directory %s: %v", testDir, err)
			}
		})

		// Create the app using lvt new (it will create the directory)
		t.Logf("Creating test app: %s (kit: %s)", opts.AppName, opts.Kit)
		if err := runLvtNew(testDir, opts.AppName, opts.Kit); err != nil {
			t.Fatalf("Failed to create app: %v", err)
		}

		// Set the app directory path
		dt.AppDir = filepath.Join(testDir, opts.AppName)
	} else {
		dt.AppDir = opts.AppDir
	}

	// Add resources if specified
	for _, resource := range opts.Resources {
		t.Logf("Adding resource: %s", resource)
		if err := runLvtGenResource(dt.AppDir, resource); err != nil {
			t.Fatalf("Failed to add resource %s: %v", resource, err)
		}
	}

	// Add auth if requested
	if opts.WithAuth {
		t.Logf("Adding authentication")
		if err := runLvtGenAuth(dt.AppDir); err != nil {
			t.Fatalf("Failed to add auth: %v", err)
		}
	}

	// Register cleanup
	t.Cleanup(func() {
		if err := dt.Cleanup(); err != nil {
			t.Errorf("Cleanup failed: %v", err)
		}
	})

	return dt
}

// Deploy executes the deployment to the configured provider
func (dt *DeploymentTest) Deploy() error {
	dt.T.Helper()
	dt.T.Logf("Deploying %s to %s (region: %s)", dt.AppName, dt.Provider, dt.Region)

	startTime := time.Now()

	// Deploy based on provider (will be implemented by provider-specific files)
	var err error
	switch dt.Provider {
	case ProviderFly:
		err = dt.deployToFly()
	case ProviderDocker:
		err = dt.deployToDocker()
	default:
		return fmt.Errorf("unsupported provider: %s", dt.Provider)
	}

	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	duration := time.Since(startTime)
	dt.T.Logf("Deployment completed in %v", duration)

	return nil
}

// VerifyHealth checks if the deployed app is responding
func (dt *DeploymentTest) VerifyHealth() error {
	dt.T.Helper()

	if dt.AppURL == "" {
		return fmt.Errorf("app URL not set")
	}

	dt.T.Logf("Verifying health at: %s", dt.AppURL)

	// TODO: Implement HTTP health check
	// For now, just check URL is set
	return nil
}

// VerifyWebSocket checks if WebSocket connection can be established
func (dt *DeploymentTest) VerifyWebSocket() error {
	dt.T.Helper()

	if dt.AppURL == "" {
		return fmt.Errorf("app URL not set")
	}

	dt.T.Logf("Verifying WebSocket connection")

	// TODO: Implement WebSocket connection test
	return nil
}

// AddCleanup registers a cleanup function to be called when test ends
func (dt *DeploymentTest) AddCleanup(fn func() error) {
	dt.cleanupMu.Lock()
	defer dt.cleanupMu.Unlock()
	dt.cleanupFuncs = append(dt.cleanupFuncs, fn)
}

// Cleanup destroys all resources created during the test
func (dt *DeploymentTest) Cleanup() error {
	dt.cleanupMu.Lock()
	defer dt.cleanupMu.Unlock()

	if dt.cleaned {
		return nil
	}
	dt.cleaned = true

	dt.T.Logf("Cleaning up deployment: %s", dt.AppName)

	var errors []error

	// Run cleanup functions in reverse order
	for i := len(dt.cleanupFuncs) - 1; i >= 0; i-- {
		if err := dt.cleanupFuncs[i](); err != nil {
			errors = append(errors, err)
			dt.T.Logf("Cleanup function %d failed: %v", i, err)
		}
	}

	// Provider-specific cleanup
	switch dt.Provider {
	case ProviderFly:
		if err := dt.cleanupFly(); err != nil {
			errors = append(errors, err)
		}
	case ProviderDocker:
		if err := dt.cleanupDocker(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup had %d error(s): %v", len(errors), errors)
	}

	dt.T.Logf("Cleanup completed successfully")
	return nil
}

// Helper functions for creating apps (will use e2e/test_helpers.go patterns)

func runLvtNew(parentDir, appName, kit string) error {
	// Build arguments for lvt new command
	args := []string{appName}
	if kit != "" {
		args = append(args, "--kit", kit)
	}

	// Call commands.New directly
	return runLvtCommandInDir(parentDir, "new", args...)
}

func runLvtGenResource(appDir, resourceSpec string) error {
	// Parse resource spec (e.g., "posts title:string content:text")
	parts := strings.Fields(resourceSpec)
	if len(parts) == 0 {
		return fmt.Errorf("empty resource spec")
	}

	// Prepend "resource" subcommand: lvt gen resource posts title:string ...
	args := append([]string{"resource"}, parts...)
	return runLvtCommandInDir(appDir, "gen", args...)
}

func runLvtGenAuth(appDir string) error {
	// Call commands.Gen with auth arguments
	return runLvtCommandInDir(appDir, "gen", "auth")
}

// runLvtCommandInDir executes an lvt command in a specific directory
// This is imported from e2e/test_helpers.go pattern - calls commands directly
func runLvtCommandInDir(workDir, command string, args ...string) error {
	// Save and restore working directory
	origDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	defer os.Chdir(origDir)

	if workDir != "" {
		if err := os.Chdir(workDir); err != nil {
			return fmt.Errorf("failed to change directory to %s: %w", workDir, err)
		}
	}

	// Import commands package and call directly
	// This avoids the complexity and slowness of shelling out
	// Note: We need to import github.com/livetemplate/lvt/commands
	// For now, use exec to build and run from the project root
	// Find the project root (where go.mod is) - start from current package location
	projectRoot, err := findProjectRoot(origDir)
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Build lvt binary once and cache it
	lvtBin, err := buildLvtBinary(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to build lvt binary: %w", err)
	}

	// Run the command
	cmdArgs := append([]string{command}, args...)
	cmd := exec.Command(lvtBin, cmdArgs...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: lvt %s: %w\nOutput: %s", strings.Join(cmdArgs, " "), err, string(output))
	}

	return nil
}

// findProjectRoot finds the directory containing go.mod, starting from startDir
func findProjectRoot(startDir string) (string, error) {
	dir := startDir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	// Clean the path to get absolute path
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (go.mod)")
		}
		dir = parent
	}
}

var lvtBinaryCache string
var lvtBinaryMutex sync.Mutex

// buildLvtBinary builds the lvt binary and caches it
// If LVT_BINARY environment variable is set, uses that instead
func buildLvtBinary(projectRoot string) (string, error) {
	lvtBinaryMutex.Lock()
	defer lvtBinaryMutex.Unlock()

	// Check for pre-built binary via environment variable
	if envBinary := os.Getenv("LVT_BINARY"); envBinary != "" {
		if _, err := os.Stat(envBinary); err == nil {
			lvtBinaryCache = envBinary
			return envBinary, nil
		}
		return "", fmt.Errorf("LVT_BINARY set to %s but file does not exist", envBinary)
	}

	// Return cached binary if it exists
	if lvtBinaryCache != "" {
		if _, err := os.Stat(lvtBinaryCache); err == nil {
			return lvtBinaryCache, nil
		}
	}

	// Build the binary
	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, "lvt-test-binary")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to build lvt: %w\nOutput: %s", err, string(output))
	}

	lvtBinaryCache = binaryPath
	return binaryPath, nil
}

// Provider-specific deployment methods (stubs for now, will be implemented in providers/)

func (dt *DeploymentTest) deployToFly() error {
	// Load credentials
	creds, err := LoadTestCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds.FlyAPIToken == "" {
		return fmt.Errorf("FLY_API_TOKEN not set")
	}

	// Create Fly.io client
	client := providers.NewFlyClient(creds.FlyAPIToken, "personal")

	// Launch app
	dt.T.Logf("Launching Fly.io app: %s", dt.AppName)
	if err := client.Launch(dt.AppName, dt.Region); err != nil {
		return fmt.Errorf("failed to launch app: %w", err)
	}

	// Register cleanup for the app
	dt.AddCleanup(func() error {
		dt.T.Logf("Destroying Fly.io app: %s", dt.AppName)
		return client.Destroy(dt.AppName)
	})

	// Create volume if needed
	dt.T.Logf("Creating volume for app: %s", dt.AppName)
	volumeID, err := client.CreateVolume(dt.AppName, dt.Region, 1) // 1GB volume
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}
	dt.T.Logf("Created volume: %s", volumeID)

	// Deploy the app
	dt.T.Logf("Deploying app from: %s", dt.AppDir)
	if err := client.Deploy(dt.AppName, dt.AppDir, dt.Region); err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	// Wait for app to be ready
	dt.T.Logf("Waiting for app to be ready...")
	if err := client.WaitForAppReady(dt.AppName, 5*time.Minute); err != nil {
		return fmt.Errorf("app failed to become ready: %w", err)
	}

	// Get app URL
	appURL, err := client.GetAppURL(dt.AppName)
	if err != nil {
		return fmt.Errorf("failed to get app URL: %w", err)
	}

	dt.AppURL = appURL
	dt.T.Logf("App deployed successfully at: %s", dt.AppURL)

	return nil
}

func (dt *DeploymentTest) deployToDocker() error {
	// Create Docker client
	// Use a unique port for this test (8000 + last 3 digits of timestamp)
	port := 8000 + (int(time.Now().Unix()) % 1000)
	client := providers.NewDockerClient(dt.AppName, port)

	// Store Docker client for debugging access
	dt.DockerClient = client

	// Ensure Dockerfile and dependencies are ready
	if err := dt.ensureDockerfile(); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	// Note: go mod tidy is now handled in the Dockerfile during docker build
	// This ensures builds use the latest tagged version without local modifications

	// Build Docker image
	dt.T.Logf("Building Docker image for: %s", dt.AppName)
	if err := client.Build(dt.AppDir); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	// Register cleanup for image and container
	dt.AddCleanup(func() error {
		// Skip cleanup if KEEP_DOCKER_CONTAINER is set (for debugging)
		if os.Getenv("KEEP_DOCKER_CONTAINER") == "true" {
			dt.T.Logf("KEEP_DOCKER_CONTAINER=true: Skipping container cleanup for debugging")
			dt.T.Logf("Container: %s, Port: %d, URL: %s", dt.AppName, port, client.GetContainerURL())
			return nil
		}
		dt.T.Logf("Destroying Docker container and image: %s", dt.AppName)
		return client.Destroy()
	})

	// Run container
	dt.T.Logf("Starting Docker container: %s on port %d", dt.AppName, port)
	if err := client.Run(); err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	// Wait for container to be ready
	dt.T.Logf("Waiting for container to be ready...")
	if err := client.WaitForReady(2 * time.Minute); err != nil {
		return fmt.Errorf("container failed to become ready: %w", err)
	}

	// Get container URL
	dt.AppURL = client.GetContainerURL()
	dt.T.Logf("Container deployed successfully at: %s", dt.AppURL)

	return nil
}

func (dt *DeploymentTest) cleanupFly() error {
	// Cleanup is handled by the cleanup functions registered during deployment
	// The Destroy() call is already registered in deployToFly()
	dt.T.Logf("Fly.io cleanup completed for: %s", dt.AppName)
	return nil
}

func (dt *DeploymentTest) cleanupDocker() error {
	// Cleanup is handled by the cleanup functions registered during deployment
	// The Destroy() call is already registered in deployToDocker()
	dt.T.Logf("Docker cleanup completed for: %s", dt.AppName)
	return nil
}

// ensureDockerfile creates a minimal Dockerfile if one doesn't exist
func (dt *DeploymentTest) ensureDockerfile() error {
	dockerfilePath := filepath.Join(dt.AppDir, "Dockerfile")

	// Check if Dockerfile already exists
	if _, err := os.Stat(dockerfilePath); err == nil {
		dt.T.Logf("Dockerfile already exists")
		return nil
	}

	dt.T.Logf("Creating minimal Dockerfile for testing")

	// Create a minimal Dockerfile for testing
	// Supports both simple kit (main.go in root) and multi kit (main.go in cmd/<appname>/)
	dockerfile := `# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev curl

# Install sqlc (for multi kit database code generation)
# Detect architecture and install appropriate binary
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ]; then SQLC_ARCH="arm64"; else SQLC_ARCH="amd64"; fi && \
    curl -L https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_${SQLC_ARCH}.tar.gz | tar -xz -C /usr/local/bin

# Copy go mod file
COPY go.mod ./

# Copy go.sum if it exists
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
# Simple kit: only needs index.tmpl in root
# Multi kit: needs internal/ directory with templates and .lvtresources
COPY --from=builder /app .

# Clean up build artifacts we don't need at runtime
RUN rm -rf /app/cmd /app/go.mod /app/go.sum /app/README.md /app/.git* 2>/dev/null || true

# Create data directory for SQLite
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./main"]
`

	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	dt.T.Logf("Created Dockerfile at: %s", dockerfilePath)
	return nil
}
