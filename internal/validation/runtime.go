package validation

import (
	"bytes"
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

	// Early exit if context is already cancelled — avoids a confusing
	// "build failed: " message when the build never actually ran.
	if ctx.Err() != nil {
		result.AddError("runtime check: cancelled: "+ctx.Err().Error(), "", 0)
		return result
	}

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

	// 3. Start the binary, capturing stdout/stderr for diagnostics.
	appCtx, appCancel := context.WithCancel(ctx)
	defer appCancel()

	var appOut bytes.Buffer
	cmd := exec.CommandContext(appCtx, binaryPath)
	cmd.Dir = appPath
	cmd.Env = append(envWithGOWORKOff(), fmt.Sprintf("PORT=%d", port))
	cmd.Stdout = &appOut
	cmd.Stderr = &appOut

	if err := cmd.Start(); err != nil {
		result.AddError(fmt.Sprintf("runtime check: failed to start app: %v", err), "", 0)
		return result
	}

	// Monitor for early exit so we can report the real error instead
	// of waiting until the full startup timeout.
	procDone := make(chan error, 1)
	go func() { procDone <- cmd.Wait() }()

	// Track whether procDone was already consumed (by waitForReadyOrExit
	// detecting early exit). Without this, the deferred cleanup would
	// deadlock trying to receive from an already-drained channel.
	procConsumed := false

	// Ensure process cleanup — kill + wait to avoid zombies.
	defer func() {
		appCancel()
		_ = cmd.Process.Kill()
		if !procConsumed {
			<-procDone
		}
	}()

	// 4. Wait for the app to be ready, or detect early process exit.
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	ready, earlyExit := waitForReadyOrExit(ctx, baseURL, timeout, procDone)
	if earlyExit {
		procConsumed = true
	}
	if !ready {
		var msg string
		switch {
		case ctx.Err() != nil:
			msg = "runtime check: cancelled: " + ctx.Err().Error()
		case earlyExit:
			msg = "runtime check: app exited before becoming ready"
		default:
			msg = fmt.Sprintf("runtime check: app did not start within %s", timeout)
		}
		if output := trimOutput(appOut.Bytes()); output != "" {
			msg += "\nOutput: " + output
		}
		result.AddError(msg, "", 0)
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
// Note: there is an inherent TOCTOU race between closing the listener and
// the subprocess binding the port. In practice this is rare, but it can
// cause flaky failures in heavily parallel CI environments.
func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port, nil
}

// waitForReadyOrExit polls the base URL until it gets a response, the process
// exits early, or the timeout elapses. Returns (ready, earlyExit).
func waitForReadyOrExit(ctx context.Context, baseURL string, timeout time.Duration, procDone <-chan error) (ready bool, earlyExit bool) {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 1 * time.Second}

	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			return false, false
		}
		// Check if process has exited.
		select {
		case <-procDone:
			return false, true
		default:
		}
		resp, err := client.Get(baseURL)
		if err == nil {
			resp.Body.Close()
			return true, false
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false, false
}

// trimOutput trims and shortens build output for error messages.
// Uses rune-aware slicing to avoid splitting multi-byte UTF-8 characters.
func trimOutput(output []byte) string {
	s := string(output)
	runes := []rune(s)
	if len(runes) > 500 {
		s = string(runes[:500]) + "..."
	}
	return s
}
