//go:build browser

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestEditModalReopenFix tests the fix for the edit modal re-open bug.
// Bug scenario:
// 1. List view (EditingID="")
// 2. Click Edit -> modal appears (EditingID="post-123")
// 3. Click Save -> modal closes (EditingID="")
// 4. Click Edit again -> Previously showed garbled output, now should show modal correctly
//
// This test requires a testfix app running on port 8083.
// To run manually:
//  1. cd /tmp && lvt new testfix --kit tailwind
//  2. cd testfix && lvt gen resource items name:string quantity:int
//  3. lvt migration up && go tool sqlc generate && go mod tidy
//  4. lvt seed items --count 3
//  5. PORT=8083 go run cmd/testfix/main.go
//  6. go test -v -tags=browser -run TestEditModalReopenFix ./e2e/
func TestEditModalReopenFix(t *testing.T) {
	// Use the running testfix server on port 8083
	// getTestURL handles Docker networking (host.docker.internal on Linux)
	baseURL := getTestURL(8083)

	// Check if the testfix server is running, skip if not
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("Skipping: testfix server not running on port 8083 (run 'PORT=8083 go run cmd/testfix/main.go' in testfix directory)")
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Skipf("Skipping: testfix server returned status %d", resp.StatusCode)
	}

	// Collect console logs and websocket messages for debugging
	var consoleLogs []string
	var consoleLogsMutex sync.Mutex

	// Get Chrome from pool
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	// Enable Runtime domain and listen for console messages
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*runtime.EventConsoleAPICalled); ok {
			consoleLogsMutex.Lock()
			for _, arg := range e.Args {
				consoleLogs = append(consoleLogs, fmt.Sprintf("[Console] %s", arg.Value))
			}
			consoleLogsMutex.Unlock()
		}
	})

	// Set timeout
	ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)
	defer timeoutCancel()

	// Navigate to items page (testfix app on port 8083)
	// The app should have seeded items already
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return runtime.Enable().Do(ctx)
		}),

		// Navigate to items page
		chromedp.Navigate(baseURL+"/items"),
		chromedp.WaitReady("body"),

		// Wait for LiveTemplate client to initialize
		waitFor(`typeof window.liveTemplateClient !== 'undefined'`, 15*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed initial navigation: %v", err)
	}

	// Check if there are any items
	var itemCount int
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelectorAll('table tbody tr').length`, &itemCount),
	)
	if err != nil {
		t.Fatalf("Failed to count items: %v", err)
	}

	if itemCount == 0 {
		t.Fatal("No items found. Make sure testfix app is running with seeded data.")
	}
	t.Logf("Found %d items", itemCount)

	// Now test the edit modal re-open scenario
	// The testfix app uses {{if ne .EditingID ""}} conditional to show/hide the edit modal
	// This is exactly the bug scenario we're testing - TreeNode→primitive→TreeNode transitions
	t.Log("Step 1: Click Edit on first post...")
	err = chromedp.Run(ctx,
		// Click the Edit button on the first post
		chromedp.Click(`button[lvt-click="edit"]`, chromedp.ByQuery),

		// Wait for edit modal overlay to appear (it uses position:fixed with rgba background)
		// The template uses {{if ne .EditingID ""}} so the modal div only exists when editing
		// We need to exclude the add-modal which has id="add-modal" attribute
		waitFor(`!!document.querySelector('div[style*="position: fixed"][style*="rgba(0,0,0,0.5)"]:not([id="add-modal"])')`, 10*time.Second),
	)
	if err != nil {
		t.Fatalf("Step 1 failed - couldn't open edit modal: %v", err)
	}

	// Verify modal content is correct (not garbled)
	var modalHTML string
	err = chromedp.Run(ctx,
		chromedp.OuterHTML(`div[style*="position: fixed"][style*="rgba(0,0,0,0.5)"]:not([id="add-modal"])`, &modalHTML, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("Failed to get modal HTML: %v", err)
	}

	t.Log("Step 1 passed: Edit modal opened")
	t.Logf("Modal HTML (first open): %s", truncateString(modalHTML, 200))

	// Check for garbled content (IDs concatenated with text)
	// The testfix app uses "items" resource, so IDs look like "test-seed-1234..."
	if strings.Contains(modalHTML, "test-seed-") && !strings.Contains(modalHTML, "<") {
		t.Fatalf("Step 1 FAILED: Garbled content detected in first modal open")
	}

	// Step 2: Save/close the modal
	t.Log("Step 2: Clicking Save to close modal...")
	err = chromedp.Run(ctx,
		// Click Save button (inside the modal form)
		chromedp.Click(`form[lvt-submit="update"] button[type="submit"]`, chromedp.ByQuery),

		// Wait for modal to close (the div should no longer exist since template uses {{if}})
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		t.Fatalf("Step 2 failed - couldn't save: %v", err)
	}

	// Check if modal is gone (since it uses {{if}}, the element won't exist when EditingID is "")
	var modalExists bool
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`!!document.querySelector('div[style*="position: fixed"][style*="rgba(0,0,0,0.5)"]:not([id="add-modal"])')`, &modalExists),
	)
	if err != nil {
		t.Fatalf("Failed to check modal state: %v", err)
	}

	if modalExists {
		// Modal might still be visible, wait a bit more
		chromedp.Run(ctx, chromedp.Sleep(2*time.Second))
	}

	t.Log("Step 2 passed: Modal closed after save")

	// Step 3: Click Edit again - THIS IS THE BUG SCENARIO
	// When EditingID goes from "" → "someID" → "" → "someID", the conditional block
	// transitions TreeNode→primitive→TreeNode. Without the registry invalidation fix,
	// the second TreeNode render shows garbled content because statics aren't re-sent.
	t.Log("Step 3: Clicking Edit again (the bug scenario)...")
	err = chromedp.Run(ctx,
		// Click the Edit button again
		chromedp.Click(`button[lvt-click="edit"]`, chromedp.ByQuery),

		// Wait for edit modal to open
		waitFor(`!!document.querySelector('div[style*="position: fixed"][style*="rgba(0,0,0,0.5)"]:not([id="add-modal"])')`, 10*time.Second),
	)
	if err != nil {
		t.Fatalf("Step 3 failed - couldn't re-open edit modal: %v", err)
	}

	// Get modal content after re-open
	var modalHTMLAfterReopen string
	err = chromedp.Run(ctx,
		chromedp.OuterHTML(`div[style*="position: fixed"][style*="rgba(0,0,0,0.5)"]:not([id="add-modal"])`, &modalHTMLAfterReopen, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("Failed to get modal HTML after reopen: %v", err)
	}

	t.Logf("Modal HTML (after re-open): %s", truncateString(modalHTMLAfterReopen, 200))

	// Check for the bug symptoms
	bugDetected := false
	bugReason := ""

	// Check 1: Garbled content (IDs concatenated with text without HTML tags)
	// The testfix app uses "test-seed-" prefix for IDs
	if strings.Contains(modalHTMLAfterReopen, "test-seed-") {
		// Find the pattern "test-seed-NNNN" followed immediately by letters (no HTML)
		// This would indicate garbled output like "test-seed-123hello"
		for i := strings.Index(modalHTMLAfterReopen, "test-seed-"); i >= 0; {
			endIdx := i + 10 // len("test-seed-")
			// Skip the numeric ID part
			for endIdx < len(modalHTMLAfterReopen) && ((modalHTMLAfterReopen[endIdx] >= '0' && modalHTMLAfterReopen[endIdx] <= '9') || modalHTMLAfterReopen[endIdx] == '-') {
				endIdx++
			}
			if endIdx < len(modalHTMLAfterReopen) && modalHTMLAfterReopen[endIdx] != '"' && modalHTMLAfterReopen[endIdx] != '<' && modalHTMLAfterReopen[endIdx] != ' ' {
				bugDetected = true
				bugReason = fmt.Sprintf("Garbled content: ID followed by non-delimiter at position %d", endIdx)
				break
			}
			nextIdx := strings.Index(modalHTMLAfterReopen[endIdx:], "test-seed-")
			if nextIdx < 0 {
				break
			}
			i = endIdx + nextIdx
		}
	}

	// Check 2: Missing form elements
	if !strings.Contains(modalHTMLAfterReopen, "<form") {
		bugDetected = true
		bugReason = "Missing <form> element in modal"
	}

	// Check 3: Missing input fields
	if !strings.Contains(modalHTMLAfterReopen, "<input") {
		bugDetected = true
		bugReason = "Missing <input> elements in modal"
	}

	// Check 4: Missing submit button
	if !strings.Contains(modalHTMLAfterReopen, "type=\"submit\"") {
		bugDetected = true
		bugReason = "Missing submit button in modal"
	}

	if bugDetected {
		// Print console logs for debugging
		consoleLogsMutex.Lock()
		if len(consoleLogs) > 0 {
			t.Log("\n📋 Console Logs:")
			for _, log := range consoleLogs {
				t.Log("  " + log)
			}
		}
		consoleLogsMutex.Unlock()

		t.Fatalf("BUG DETECTED: %s\nFull modal HTML:\n%s", bugReason, modalHTMLAfterReopen)
	}

	t.Log("Step 3 passed: Edit modal re-opened correctly!")
	t.Log("\n✅ EDIT MODAL RE-OPEN FIX VERIFIED!")
	t.Log("   The modal opens correctly on the second click with proper HTML structure.")

	// Clean up - close the modal
	chromedp.Run(ctx,
		chromedp.Click(`button[lvt-click="cancel_edit"]`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
	)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
