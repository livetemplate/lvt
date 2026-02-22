package validation

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/livetemplate/lvt/internal/validator"
)

// RuntimeCheck builds the app binary, starts it, and probes HTTP routes.
// It implements the Check interface and is opt-in (not in DefaultEngine)
// because it is expensive: it compiles a binary, starts a subprocess,
// and makes HTTP requests.
type RuntimeCheck struct {
	// StartupTimeout is how long to wait for the app to start listening.
	// Defaults to 15s if zero.
	StartupTimeout time.Duration
	// Port to run the app on. 0 means auto-allocate a free port.
	Port int
	// Routes to probe after the app starts. Defaults to ["/"] if empty.
	Routes []string
}

func (c *RuntimeCheck) Name() string { return "runtime" }

func (c *RuntimeCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()

	timeout := c.StartupTimeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	routes := c.Routes
	if len(routes) == 0 {
		routes = []string{"/"}
	}

	// 1. Build the binary
	binaryPath := filepath.Join(appPath, "lvt-runtime-check")
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = appPath
	buildCmd.Env = envWithGOWORKOff()
	if output, err := buildCmd.CombinedOutput(); err != nil {
		result.AddError(fmt.Sprintf("runtime check: build failed: %s", trimOutput(output)), "", 0)
		return result
	}
	defer os.Remove(binaryPath)

	// 2. Allocate a free port
	port := c.Port
	if port == 0 {
		var err error
		port, err = getFreePort()
		if err != nil {
			result.AddError(fmt.Sprintf("runtime check: failed to allocate port: %v", err), "", 0)
			return result
		}
	}

	// 3. Start the binary
	appCtx, appCancel := context.WithCancel(ctx)
	defer appCancel()

	cmd := exec.CommandContext(appCtx, binaryPath)
	cmd.Dir = appPath
	cmd.Env = append(envWithGOWORKOff(), fmt.Sprintf("PORT=%d", port))
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		result.AddError(fmt.Sprintf("runtime check: failed to start app: %v", err), "", 0)
		return result
	}

	// Ensure process cleanup — kill + wait to avoid zombies.
	defer func() {
		appCancel()
		// Process may already be gone; ignore errors.
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	// 4. Wait for the app to be ready
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if !waitForReady(ctx, baseURL, timeout) {
		result.AddError(fmt.Sprintf("runtime check: app did not start within %s", timeout), "", 0)
		return result
	}

	// 5. Probe routes
	client := &http.Client{Timeout: 5 * time.Second}
	for _, route := range routes {
		probeURL := baseURL + route
		resp, err := client.Get(probeURL)
		if err != nil {
			result.AddError(fmt.Sprintf("runtime check: failed to reach %s: %v", route, err), "", 0)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 500 {
			result.AddError(fmt.Sprintf("runtime check: %s returned %d", route, resp.StatusCode), "", 0)
		} else if resp.StatusCode >= 400 {
			result.AddWarning(fmt.Sprintf("runtime check: %s returned %d", route, resp.StatusCode), "", 0)
		}
	}

	return result
}

// getFreePort asks the OS for an available TCP port.
func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port, nil
}

// waitForReady polls the base URL until it gets a response or times out.
func waitForReady(ctx context.Context, baseURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 1 * time.Second}

	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			return false
		}
		resp, err := client.Get(baseURL)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// trimOutput trims and shortens build output for error messages.
func trimOutput(output []byte) string {
	s := string(output)
	if len(s) > 500 {
		s = s[:500] + "..."
	}
	return s
}
