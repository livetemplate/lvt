package testing

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

//go:embed livetemplate-client.browser.js
var clientLibraryJS []byte

const (
	dockerImage           = "chromedp/headless-shell:latest"
	chromeContainerPrefix = "chrome-e2e-test-"
	staleContainerGrace   = 10 * time.Minute
)

func removeContainersByFilter(filter string, shouldRemove func(string) (bool, error)) ([]string, error) {
	listCmd := exec.Command("docker", "ps", "-aq", "-f", fmt.Sprintf("name=%s", filter))
	output, err := listCmd.Output()
	if err != nil {
		return nil, err
	}
	ids := strings.Fields(strings.TrimSpace(string(output)))
	if shouldRemove != nil {
		filtered := make([]string, 0, len(ids))
		for _, id := range ids {
			ok, err := shouldRemove(id)
			if err != nil {
				return nil, err
			}
			if ok {
				filtered = append(filtered, id)
			}
		}
		ids = filtered
	}
	if len(ids) == 0 {
		return nil, nil
	}
	args := append([]string{"rm", "-f"}, ids...)
	rmCmd := exec.Command("docker", args...)
	if rmOutput, err := rmCmd.CombinedOutput(); err != nil {
		return ids, fmt.Errorf("docker rm failed: %w (%s)", err, strings.TrimSpace(string(rmOutput)))
	}
	return ids, nil
}

func cleanupContainerByName(tb testing.TB, name string) {
	if removed, err := removeContainersByFilter(name, nil); err != nil {
		if tb != nil {
			tb.Logf("warning: failed to clean Chrome container %s: %v", name, err)
		} else {
			fmt.Fprintf(os.Stderr, "warning: failed to clean Chrome container %s: %v\n", name, err)
		}
	} else if len(removed) > 0 {
		msg := fmt.Sprintf("Cleaned up leftover Chrome container(s): %s", strings.Join(removed, ", "))
		if tb != nil {
			tb.Log(msg)
		} else {
			fmt.Fprintln(os.Stderr, msg)
		}
	}
}

// CleanupChromeContainers removes any lingering Chrome containers created by the test helpers.
func CleanupChromeContainers() {
	shouldRemove := func(id string) (bool, error) {
		inspectCmd := exec.Command("docker", "inspect", "--format", "{{.State.Running}} {{.State.StartedAt}}", id)
		inspectOutput, err := inspectCmd.Output()
		if err != nil {
			return true, err
		}
		fields := strings.Fields(strings.TrimSpace(string(inspectOutput)))
		if len(fields) < 2 {
			return true, nil
		}
		running := fields[0] == "true"
		startedAt, err := time.Parse(time.RFC3339Nano, fields[1])
		if err != nil {
			startedAt, err = time.Parse(time.RFC3339, fields[1])
			if err != nil {
				return true, nil
			}
		}
		if !running {
			return true, nil
		}
		return time.Since(startedAt) > staleContainerGrace, nil
	}
	if removed, err := removeContainersByFilter(chromeContainerPrefix, shouldRemove); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to clean Chrome containers: %v\n", err)
	} else if len(removed) > 0 {
		fmt.Fprintf(os.Stderr, "Cleaned up %d lingering Chrome container(s): %s\n", len(removed), strings.Join(removed, ", "))
	}
}

// GetFreePort asks the kernel for a free open port that is ready to use
func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

// GetChromeTestURL returns the URL for Chrome (in Docker) to access the test server
// Chrome container uses host.docker.internal to reach the host on all platforms
func GetChromeTestURL(port int) string {
	portStr := fmt.Sprintf("%d", port)
	return "http://host.docker.internal:" + portStr
}

// StartDockerChrome starts the chromedp headless-shell Docker container
// Returns an error if the container fails to start or Chrome fails to become ready
func StartDockerChrome(t *testing.T, debugPort int) error {
	t.Helper()

	// Check if Docker is available
	if _, err := exec.Command("docker", "version").CombinedOutput(); err != nil {
		t.Skip("Docker not available, skipping E2E test")
	}

	containerName := fmt.Sprintf("%s%d", chromeContainerPrefix, debugPort)
	cleanupContainerByName(t, containerName)

	// Check if image exists, if not try to pull it (with timeout)
	checkCmd := exec.Command("docker", "image", "inspect", dockerImage)
	if _, err := checkCmd.CombinedOutput(); err != nil {
		// Image doesn't exist, try to pull with timeout
		t.Log("Pulling chromedp/headless-shell Docker image...")

		// Use a context with timeout for the pull operation
		pullCtx, pullCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer pullCancel()

		pullCmd := exec.CommandContext(pullCtx, "docker", "pull", dockerImage)
		// Use Output() to properly close pipes and avoid I/O wait
		if output, err := pullCmd.CombinedOutput(); err != nil {
			if pullCtx.Err() == context.DeadlineExceeded {
				t.Fatal("Docker pull timed out after 60 seconds")
			}
			t.Fatalf("Failed to pull Docker image: %v\nOutput: %s", err, output)
		}
		t.Log("✅ Docker image pulled successfully")
	} else {
		t.Log("✅ Docker image already exists, skipping pull")
	}

	// Start the container in detached mode to avoid I/O wait issues
	t.Log("Starting Chrome headless Docker container...")
	portMapping := fmt.Sprintf("%d:9222", debugPort)

	// Run in detached mode (-d) to avoid process I/O issues during cleanup
	// Don't use --rm here; we'll clean up manually in StopDockerChrome
	cmd := exec.Command("docker", "run", "-d",
		"-p", portMapping,
		"--name", containerName,
		"--add-host", "host.docker.internal:host-gateway",
		dockerImage,
	)

	// Use Output() instead of Run() to properly close pipes and avoid I/O wait
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("failed to start Chrome Docker container: %w", err)
	}

	// Container runs independently in detached mode

	// Wait for Chrome to be ready (increased timeout for slower systems)
	t.Log("Waiting for Chrome to be ready...")
	chromeURL := fmt.Sprintf("http://localhost:%d/json/version", debugPort)
	ready := false
	var lastErr error
	for i := 0; i < 120; i++ { // 120 iterations × 500ms = 60 seconds
		resp, err := http.Get(chromeURL)
		if err == nil {
			resp.Body.Close()
			ready = true
			t.Logf("✅ Chrome ready after %d attempts (%.1fs)", i+1, float64(i+1)*0.5)
			break
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		t.Logf("Chrome failed to start within 60 seconds. Last error: %v", lastErr)

		// Try to get container logs for debugging
		logsCmd := exec.Command("docker", "logs", "--tail", "50", containerName)
		if output, err := logsCmd.CombinedOutput(); err == nil && len(output) > 0 {
			t.Logf("Chrome container logs:\n%s", string(output))
		}

		// Clean up the container since Chrome didn't start properly
		_, _ = exec.Command("docker", "rm", "-f", containerName).CombinedOutput()
		return fmt.Errorf("Chrome failed to start within 60 seconds: %w", lastErr)
	}

	t.Log("✅ Chrome headless Docker container ready")
	return nil
}

// StopDockerChrome stops and removes the Chrome Docker container
func StopDockerChrome(t *testing.T, debugPort int) {
	t.Helper()
	t.Log("Stopping Chrome Docker container...")

	containerName := fmt.Sprintf("chrome-e2e-test-%d", debugPort)

	// docker rm -f will stop and remove the container in one command
	// The -f flag forces removal even if the container is running
	rmCmd := exec.Command("docker", "rm", "-f", containerName)
	// Use Output() to properly close pipes and avoid I/O wait
	if output, err := rmCmd.CombinedOutput(); err != nil {
		// Only log if it's not a "no such container" error (which is fine)
		errMsg := string(output)
		if !strings.Contains(errMsg, "No such container") && !strings.Contains(err.Error(), "No such container") {
			t.Logf("Warning: Failed to remove Docker container: %v (output: %s)", err, errMsg)
		}
	}
}

// StartTestServer starts a Go server on the specified port
// mainPath should be the path to main.go (e.g., "main.go" or "../../examples/counter/main.go")
func StartTestServer(t *testing.T, mainPath string, port int) *exec.Cmd {
	t.Helper()

	portStr := fmt.Sprintf("%d", port)
	serverURL := fmt.Sprintf("http://localhost:%d", port)

	t.Logf("Starting test server on port %s", portStr)
	cmd := exec.Command("go", "run", mainPath)
	cmd.Env = append([]string{
		"PORT=" + portStr,
		"LVT_DEV_MODE=true", // Use local client library in tests
	}, cmd.Environ()...)

	// Start the server
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to be ready
	ready := false
	for i := 0; i < 50; i++ { // 5 seconds
		resp, err := http.Get(serverURL)
		if err == nil {
			resp.Body.Close()
			ready = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !ready {
		_ = cmd.Process.Kill()
		t.Fatal("Server failed to start within 5 seconds")
	}

	// Register cleanup handler to kill server process on test completion/failure
	t.Cleanup(func() {
		if cmd.Process != nil {
			t.Logf("Killing test server process (PID: %d)...", cmd.Process.Pid)
			if err := cmd.Process.Kill(); err != nil {
				t.Logf("Warning: Failed to kill test server process (PID: %d): %v", cmd.Process.Pid, err)
			} else {
				t.Logf("✅ Test server process killed (PID: %d)", cmd.Process.Pid)
			}
			// Wait for the process to exit to clean up zombie processes and avoid I/O wait
			_ = cmd.Wait()
		}
	})

	t.Logf("✅ Test server ready at %s", serverURL)
	return cmd
}

// ServeClientLibrary serves the LiveTemplate client browser bundle from embedded bytes.
// This is for development/testing purposes only. In production, serve from CDN.
func ServeClientLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(clientLibraryJS)
}

// WaitFor polls a JavaScript condition until it returns true or timeout is reached.
// This is a generic condition-based wait utility that eliminates arbitrary sleeps.
//
// The condition must be a JavaScript expression that evaluates to a boolean.
//
// Examples:
//   - WaitFor("document.getElementById('modal').style.display === 'flex'", 5*time.Second)
//   - WaitFor("document.querySelector('.item').textContent === 'Hello'", 3*time.Second)
//   - WaitFor("document.querySelectorAll('.item').length === 5", 5*time.Second)
//   - WaitFor("!document.getElementById('modal').hasAttribute('hidden')", 2*time.Second)
func WaitFor(condition string, timeout time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		startTime := time.Now()

		for {
			// Check context first to fail fast if parent canceled
			select {
			case <-ctx.Done():
				return fmt.Errorf("context canceled while waiting for condition '%s': %w", condition, ctx.Err())
			default:
			}

			var result bool
			err := chromedp.Evaluate(condition, &result).Do(ctx)

			if err != nil {
				// Don't fail immediately on evaluation errors - the DOM might not be ready yet
				// Just log and continue polling
				if time.Since(startTime) > timeout {
					return fmt.Errorf("timeout waiting for condition '%s' (last error: %v)", condition, err)
				}
			} else if result {
				// Condition met
				return nil
			}

			if time.Since(startTime) > timeout {
				// Get debug info on timeout
				var debugInfo string
				_ = chromedp.Evaluate(`document.readyState`, &debugInfo).Do(ctx)
				return fmt.Errorf("timeout waiting for condition '%s' after %v (readyState: %s)", condition, timeout, debugInfo)
			}

			// Poll every 100ms (increased from 10ms for stability)
			// This reduces CPU thrashing and makes checks more reliable
			time.Sleep(100 * time.Millisecond)
		}
	})
}

// WaitForText waits for an element's text content to include the specified text.
// This is a convenience wrapper around WaitFor for common text-matching scenarios.
//
// The selector is a CSS selector, and text is the substring to match.
// Returns an error if the condition is not met within the timeout.
//
// Examples:
//
//	WaitForText("section", "No results found", 5*time.Second)
//	WaitForText("tbody", "First Todo Item", 10*time.Second)
//	WaitForText(".status", "Connected", 3*time.Second)
func WaitForText(selector, text string, timeout time.Duration) chromedp.Action {
	condition := fmt.Sprintf(`
		(() => {
			const el = document.querySelector('%s');
			return el && el.textContent && el.textContent.includes(%q);
		})()
	`, selector, text)
	return WaitFor(condition, timeout)
}

// WaitForCount waits for a specific number of elements to match the selector.
// This is a convenience wrapper around WaitFor for counting elements.
//
// The selector is a CSS selector, and expectedCount is the exact number of elements expected.
// Returns an error if the condition is not met within the timeout.
//
// Examples:
//
//	WaitForCount("tbody tr", 3, 5*time.Second)
//	WaitForCount(".todo-item", 10, 10*time.Second)
//	WaitForCount("button[disabled]", 0, 3*time.Second)
func WaitForCount(selector string, expectedCount int, timeout time.Duration) chromedp.Action {
	condition := fmt.Sprintf(`document.querySelectorAll('%s').length === %d`, selector, expectedCount)
	return WaitFor(condition, timeout)
}

// WaitForWebSocketReady waits for the first WebSocket update to be applied
// by polling for the removal of data-lvt-loading attribute (condition-based waiting).
// This ensures E2E tests run after the WebSocket connection is established and
// initial state is synchronized.
func WaitForWebSocketReady(timeout time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Wait for wrapper element to exist
		err := chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery).Do(ctx)
		if err != nil {
			return fmt.Errorf("wrapper element not found: %w", err)
		}

		// Give client time to initialize before we start polling (reduces race conditions)
		time.Sleep(50 * time.Millisecond)

		// Use improved WaitFor with slower polling (100ms instead of 10ms)
		// This reduces CPU thrashing and makes the check more stable
		return WaitFor(`
			(() => {
				const wrapper = document.querySelector('[data-lvt-id]');
				// Also verify that the client is actually initialized
				const client = window.liveTemplateClient;
				const clientInitialized = client !== undefined;
				const clientReady = typeof client?.isReady === 'function' ? client.isReady() : false;
				return wrapper && clientInitialized && clientReady;
			})()
		`, timeout).Do(ctx)
	})
}

// ValidateNoTemplateExpressions checks that the specified element does not contain
// raw Go template expressions like {{if}}, {{range}}, {{define}}, etc.
// This catches the bug where unflattened templates are used in WebSocket tree generation.
func ValidateNoTemplateExpressions(selector string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var innerHTML string
		if err := chromedp.InnerHTML(selector, &innerHTML, chromedp.ByQuery).Do(ctx); err != nil {
			return fmt.Errorf("failed to get innerHTML of %s: %w", selector, err)
		}

		// Check for common template expressions
		templateExpressions := []string{
			"{{if",
			"{{range",
			"{{define",
			"{{template",
			"{{with",
			"{{block",
			"{{else",
			"{{end}}",
		}

		for _, expr := range templateExpressions {
			if strings.Contains(innerHTML, expr) {
				// Find context around the expression for better error messages
				idx := strings.Index(innerHTML, expr)
				start := idx - 50
				if start < 0 {
					start = 0
				}
				end := idx + 100
				if end > len(innerHTML) {
					end = len(innerHTML)
				}
				context := innerHTML[start:end]

				return fmt.Errorf("raw template expression '%s' found in HTML. Context: ...%s...", expr, context)
			}
		}

		return nil
	})
}

// WaitForMessageCount waits for the WebSocket message counter to reach the expected value.
// This is a deterministic way to wait for updates without relying on arbitrary timeouts.
// The client increments window.__wsMessageCount after each DOM update completes.
//
// Example:
//
//	var initialCount int
//	chromedp.Evaluate(`window.__wsMessageCount || 0`, &initialCount)
//	// ... trigger action ...
//	WaitForMessageCount(initialCount+1, 5*time.Second)  // Wait for exactly 1 new message
func WaitForMessageCount(expectedCount int, timeout time.Duration) chromedp.Action {
	condition := fmt.Sprintf(`(window.__wsMessageCount || 0) >= %d`, expectedCount)
	return WaitFor(condition, timeout)
}

// WaitForActionResponse waits for a WebSocket message with the specified action name in metadata.
// This is deterministic - it waits for the exact action to complete, not arbitrary time.
//
// The test should clear window.__wsMessages before triggering the action:
//
//	chromedp.Evaluate(`window.__wsMessages = [];`, nil)
//	// ... trigger action ...
//	WaitForActionResponse("search", 5*time.Second)
//
// Example:
//
//	chromedp.Evaluate(`window.__wsMessages = [];`, nil),
//	chromedp.SendKeys(`input[name="query"]`, "test", chromedp.ByQuery),
//	WaitForActionResponse("search", 5*time.Second),
func WaitForActionResponse(actionName string, timeout time.Duration) chromedp.Action {
	condition := fmt.Sprintf(`
		(() => {
			const msgs = window.__wsMessages || [];
			return msgs.some(m => m.meta && m.meta.action === %q);
		})()
	`, actionName)
	return WaitFor(condition, timeout)
}

// SetupUpdateEventListener sets up a non-blocking event listener that captures
// 'lvt:updated' events in window.__capturedEvents array.
// This must be called BEFORE the action that triggers the event.
// Then use WaitForUpdateEvent to poll for the captured event.
//
// Example:
//
//	chromedp.Run(ctx,
//	    SetupUpdateEventListener(),           // Setup listener (non-blocking)
//	    chromedp.SendKeys(...),                // Trigger action
//	    WaitForUpdateEvent("search", 5*time.Second),  // Poll for captured event
//	)
func SetupUpdateEventListener() chromedp.Action {
	return chromedp.Evaluate(`
		(() => {
			window.__capturedEvents = [];
			const wrapper = document.querySelector('[data-lvt-id]');
			if (wrapper) {
				// Remove existing listener if any
				if (window.__capturedEventsHandler) {
					wrapper.removeEventListener('lvt:updated', window.__capturedEventsHandler);
				}
				// Create and store new handler
				window.__capturedEventsHandler = (e) => {
					window.__capturedEvents.push({
						action: e.detail.action,
						messageCount: e.detail.messageCount,
						success: e.detail.success,
						timestamp: Date.now()
					});
				};
				wrapper.addEventListener('lvt:updated', window.__capturedEventsHandler);
			}
		})();
	`, nil)
}

// WaitForUpdateEvent polls for a captured 'lvt:updated' event.
// Must be used after SetupUpdateEventListener().
// Optionally filters by action name if provided.
//
// This is deterministic - it waits for the actual event to fire, not arbitrary timeouts.
func WaitForUpdateEvent(actionName string, timeout time.Duration) chromedp.Action {
	condition := `
		(() => {
			const events = window.__capturedEvents || [];
			return events.length > 0;
		})()
	`

	if actionName != "" {
		condition = fmt.Sprintf(`
			(() => {
				const events = window.__capturedEvents || [];
				return events.some(e => e.action === %q);
			})()
		`, actionName)
	}

	return WaitFor(condition, timeout)
}
