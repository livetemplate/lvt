//go:build browser

package e2e

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

// Wrapper functions to use cmd/lvt/testing utilities with shorter names in tests

// Unused: Kept for potential future use
// func startDockerChrome(t *testing.T, debugPort int) *exec.Cmd {
// 	return e2etest.StartDockerChrome(t, debugPort)
// }

// Unused: Kept for potential future use
// func stopDockerChrome(t *testing.T, cmd *exec.Cmd, debugPort int) {
// 	e2etest.StopDockerChrome(t, cmd, debugPort)
// }

func getTestURL(port int) string {
	return e2etest.GetChromeTestURL(port)
}

func waitFor(condition string, timeout time.Duration) chromedp.Action {
	return e2etest.WaitFor(condition, timeout)
}

func waitForWebSocketReady(timeout time.Duration) chromedp.Action {
	// Use optimized timeout: 10s local, 30s CI (unless explicitly overridden)
	optimizedTimeout := getTimeout("WEBSOCKET_TIMEOUT", 10*time.Second, 30*time.Second)
	// If caller passes a custom timeout, respect it
	if timeout > 0 && timeout != 30*time.Second {
		optimizedTimeout = timeout
	}
	return e2etest.WaitForWebSocketReady(optimizedTimeout)
}

func validateNoTemplateExpressions(selector string) chromedp.Action {
	return e2etest.ValidateNoTemplateExpressions(selector)
}

// getBrowserTimeout returns optimized browser operation timeout
// Local: 20s (faster feedback), CI: 60s for stable operation
func getBrowserTimeout() time.Duration {
	return getTimeout("BROWSER_TIMEOUT", 20*time.Second, 120*time.Second)
}

// clickUntilModalOpens retries clicking a button until a target element becomes
// visible. This handles the race where the LiveTemplate client library's DOM
// event handlers haven't attached yet when the click fires.
func clickUntilModalOpens(buttonSelector, targetSelector string, timeout time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			_ = chromedp.Click(buttonSelector, chromedp.ByQuery).Do(ctx)
			var visible bool
			_ = chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const el = document.querySelector(%q);
					return el !== null && el.offsetParent !== null;
				})()
			`, targetSelector), &visible).Do(ctx)
			if visible {
				return nil
			}
			time.Sleep(200 * time.Millisecond)
		}
		return fmt.Errorf("target %q did not become visible after clicking %q for %v", targetSelector, buttonSelector, timeout)
	})
}
