//go:build browser

package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

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
			t.Errorf("Template error found on page: contains %q", pattern)
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
		t.Errorf("WebSocket not connected (readyState: %d)", wsReadyState)
	} else {
		t.Log("WebSocket connected")
	}
}
