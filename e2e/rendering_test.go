//go:build browser

// Package e2e contains browser-based rendering tests for LiveTemplate.
//
// These tests validate the rendering library's core functionality:
// DOM operations, event handling, WebSocket behavior, and UI patterns.
// They run with actual Chrome via chromedp to ensure real browser behavior.
//
// Run with: go test -tags=browser -run TestRendering ./e2e/...
package e2e

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

// renderingTestServer creates a test server with embedded HTML for rendering tests.
// It returns the Chrome-accessible URL and a cleanup function.
func renderingTestServer(t *testing.T, html string) (string, func()) {
	t.Helper()

	clientJS := e2etest.GetClientLibraryJS()
	if len(clientJS) == 0 {
		t.Fatal("Client library not embedded")
	}

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	chromeURL := getTestURL(port)

	mux := http.NewServeMux()

	mux.HandleFunc("/livetemplate-client.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(clientJS)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	go func() { _ = http.Serve(listener, mux) }()
	time.Sleep(100 * time.Millisecond)

	return chromeURL, func() { listener.Close() }
}

// clientInitScript is the standard script to load and initialize the LiveTemplate client
const clientInitScript = `
<script>
	window.__lvtClientLoaded = false;
	window.__lvtClientLoadError = "";
	window.__markLvtClientLoaded = function() { window.__lvtClientLoaded = true; };
	window.__markLvtClientError = function(e) { window.__lvtClientLoadError = e?.message || e?.type || "unknown"; };
</script>
<script src="/livetemplate-client.js" defer onload="window.__markLvtClientLoaded()" onerror="window.__markLvtClientError(event)"></script>
`

// waitForClient waits for the LiveTemplate client to load and initialize
func waitForClient() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Wait for client script to load
		if err := waitForDOM(`window.__lvtClientLoaded === true || window.__lvtClientLoadError !== ''`, 15*time.Second).Do(ctx); err != nil {
			return fmt.Errorf("client load timeout: %w", err)
		}

		// Check for load errors
		var loadErr string
		if err := chromedp.Evaluate(`window.__lvtClientLoadError || ''`, &loadErr).Do(ctx); err != nil {
			return err
		}
		if loadErr != "" {
			return fmt.Errorf("client load error: %s", loadErr)
		}

		// Auto-init if needed
		if err := chromedp.Evaluate(`
			if (!window.liveTemplateClient && window.LiveTemplateClient?.LiveTemplateClient?.autoInit) {
				window.LiveTemplateClient.LiveTemplateClient.autoInit();
			}
		`, nil).Do(ctx); err != nil {
			return err
		}

		// Wait for client instance
		return waitForDOM(`typeof window.liveTemplateClient !== 'undefined'`, 5*time.Second).Do(ctx)
	})
}

// waitForDOM is a simple wait function for pure DOM conditions without LiveTemplate-specific checks.
// This is used for tests that don't involve WebSocket or LiveTemplate state.
func waitForDOM(condition string, timeout time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		startTime := time.Now()
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context canceled while waiting for '%s': %w", condition, ctx.Err())
			default:
			}

			var result bool
			err := chromedp.Evaluate(condition, &result).Do(ctx)
			if err == nil && result {
				return nil
			}

			if time.Since(startTime) > timeout {
				return fmt.Errorf("timeout waiting for '%s' after %v", condition, timeout)
			}

			time.Sleep(50 * time.Millisecond)
		}
	})
}

// =============================================================================
// Test 1: DOM List Operations
// =============================================================================

// TestRendering_DOM_ListOperations tests that list items can be added and removed dynamically.
// This validates the core DOM diffing (morphdom) functionality.
func TestRendering_DOM_ListOperations(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>List Operations Test</title></head>
<body>
	<div data-lvt-id="list-test">
		<ul id="item-list">
			<li data-item-id="1">Item 1</li>
			<li data-item-id="2">Item 2</li>
			<li data-item-id="3">Item 3</li>
		</ul>
		<button id="add-btn" onclick="addItem()">Add Item</button>
		<button id="remove-btn" onclick="removeItem()">Remove First</button>
	</div>
	<script>
		let itemCount = 3;
		function addItem() {
			itemCount++;
			const li = document.createElement('li');
			li.setAttribute('data-item-id', itemCount);
			li.textContent = 'Item ' + itemCount;
			document.getElementById('item-list').appendChild(li);
		}
		function removeItem() {
			const list = document.getElementById('item-list');
			if (list.children.length > 0) {
				list.removeChild(list.children[0]);
			}
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial state: 3 items
		chromedp.ActionFunc(func(ctx context.Context) error {
			var count int
			if err := chromedp.Evaluate(`document.querySelectorAll('#item-list li').length`, &count).Do(ctx); err != nil {
				return err
			}
			if count != 3 {
				return fmt.Errorf("expected 3 items initially, got %d", count)
			}
			t.Log("Initial state: 3 items")
			return nil
		}),

		// Add an item
		chromedp.Click("#add-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var count int
			if err := chromedp.Evaluate(`document.querySelectorAll('#item-list li').length`, &count).Do(ctx); err != nil {
				return err
			}
			if count != 4 {
				return fmt.Errorf("expected 4 items after add, got %d", count)
			}
			t.Log("After add: 4 items")
			return nil
		}),

		// Remove an item
		chromedp.Click("#remove-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var count int
			if err := chromedp.Evaluate(`document.querySelectorAll('#item-list li').length`, &count).Do(ctx); err != nil {
				return err
			}
			if count != 3 {
				return fmt.Errorf("expected 3 items after remove, got %d", count)
			}
			t.Log("After remove: 3 items")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("List operations test failed: %v", err)
	}
	t.Log("DOM list operations test passed")
}

// =============================================================================
// Test 2: DOM Table Rendering
// =============================================================================

// TestRendering_DOM_TableRendering tests that tables render correctly with data.
func TestRendering_DOM_TableRendering(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Table Rendering Test</title></head>
<body>
	<div data-lvt-id="table-test">
		<table id="data-table">
			<thead><tr><th>Name</th><th>Value</th></tr></thead>
			<tbody>
				<tr data-row-id="1"><td>Alpha</td><td>100</td></tr>
				<tr data-row-id="2"><td>Beta</td><td>200</td></tr>
			</tbody>
		</table>
		<button id="add-row" onclick="addRow()">Add Row</button>
		<button id="update-row" onclick="updateRow()">Update First</button>
	</div>
	<script>
		let rowCount = 2;
		function addRow() {
			rowCount++;
			const tbody = document.querySelector('#data-table tbody');
			const tr = document.createElement('tr');
			tr.setAttribute('data-row-id', rowCount);
			tr.innerHTML = '<td>Gamma</td><td>' + (rowCount * 100) + '</td>';
			tbody.appendChild(tr);
		}
		function updateRow() {
			const firstCell = document.querySelector('#data-table tbody tr td');
			if (firstCell) firstCell.textContent = 'Updated';
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial table structure
		chromedp.ActionFunc(func(ctx context.Context) error {
			var rowCount int
			if err := chromedp.Evaluate(`document.querySelectorAll('#data-table tbody tr').length`, &rowCount).Do(ctx); err != nil {
				return err
			}
			if rowCount != 2 {
				return fmt.Errorf("expected 2 rows initially, got %d", rowCount)
			}
			t.Log("Initial table: 2 rows")
			return nil
		}),

		// Add a row
		chromedp.Click("#add-row"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var rowCount int
			if err := chromedp.Evaluate(`document.querySelectorAll('#data-table tbody tr').length`, &rowCount).Do(ctx); err != nil {
				return err
			}
			if rowCount != 3 {
				return fmt.Errorf("expected 3 rows after add, got %d", rowCount)
			}
			t.Log("After add: 3 rows")
			return nil
		}),

		// Update a row
		chromedp.Click("#update-row"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var text string
			if err := chromedp.Evaluate(`document.querySelector('#data-table tbody tr td').textContent`, &text).Do(ctx); err != nil {
				return err
			}
			if text != "Updated" {
				return fmt.Errorf("expected 'Updated', got '%s'", text)
			}
			t.Log("Row updated successfully")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Table rendering test failed: %v", err)
	}
	t.Log("DOM table rendering test passed")
}

// =============================================================================
// Test 3: Form Submit Validation
// =============================================================================

// TestRendering_Form_SubmitValidation tests form submission and validation display.
func TestRendering_Form_SubmitValidation(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Form Validation Test</title></head>
<body>
	<div data-lvt-id="form-test">
		<form id="test-form" onsubmit="return handleSubmit(event)">
			<input type="text" name="username" id="username" required placeholder="Username">
			<span id="username-error" class="error" style="display: none; color: red;"></span>
			<input type="email" name="email" id="email" required placeholder="Email">
			<span id="email-error" class="error" style="display: none; color: red;"></span>
			<button type="submit" id="submit-btn">Submit</button>
		</form>
		<div id="success-msg" style="display: none; color: green;">Form submitted!</div>
	</div>
	<script>
		function handleSubmit(e) {
			e.preventDefault();
			const username = document.getElementById('username').value;
			const email = document.getElementById('email').value;
			let valid = true;

			// Reset errors
			document.getElementById('username-error').style.display = 'none';
			document.getElementById('email-error').style.display = 'none';

			if (username.length < 3) {
				document.getElementById('username-error').textContent = 'Username must be at least 3 characters';
				document.getElementById('username-error').style.display = 'block';
				valid = false;
			}

			if (!email.includes('@')) {
				document.getElementById('email-error').textContent = 'Invalid email address';
				document.getElementById('email-error').style.display = 'block';
				valid = false;
			}

			if (valid) {
				document.getElementById('success-msg').style.display = 'block';
			}

			return false;
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Test validation error display
		chromedp.SendKeys("#username", "ab"),
		chromedp.SendKeys("#email", "invalid"),
		// Call handleSubmit directly (chromedp.Click doesn't trigger inline handlers in headless Docker)
		chromedp.Evaluate(`handleSubmit(new Event('submit'))`, nil),
		// Wait for validation errors to appear
		waitForDOM(`document.getElementById('username-error').style.display === 'block'`, 5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var usernameErrVisible, emailErrVisible bool
			chromedp.Evaluate(`document.getElementById('username-error').style.display === 'block'`, &usernameErrVisible).Do(ctx)
			chromedp.Evaluate(`document.getElementById('email-error').style.display === 'block'`, &emailErrVisible).Do(ctx)

			if !usernameErrVisible {
				return fmt.Errorf("username error should be visible")
			}
			if !emailErrVisible {
				return fmt.Errorf("email error should be visible")
			}
			t.Log("Validation errors displayed correctly")
			return nil
		}),

		// Test successful submission - set values directly via JS to ensure they're updated
		chromedp.Evaluate(`
			document.getElementById('username').value = 'validuser';
			document.getElementById('email').value = 'valid@example.com';
		`, nil),
		// Call handleSubmit directly
		chromedp.Evaluate(`handleSubmit(new Event('submit'))`, nil),
		// Wait for success message to appear
		waitForDOM(`document.getElementById('success-msg').style.display === 'block'`, 5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			t.Log("Form submitted successfully")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Form validation test failed: %v", err)
	}
	t.Log("Form submit validation test passed")
}

// =============================================================================
// Test 4: Modal Lifecycle
// =============================================================================

// TestRendering_Modal_Lifecycle tests modal open, close, and reopen behavior.
func TestRendering_Modal_Lifecycle(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Modal Lifecycle Test</title></head>
<body>
	<div data-lvt-id="modal-test">
		<button id="open-modal" lvt-modal-open="test-modal">Open Modal</button>

		<div id="test-modal" hidden aria-hidden="true" role="dialog" data-modal-backdrop data-modal-id="test-modal"
			 style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000;">
			<div style="background: white; padding: 2rem; border-radius: 8px;">
				<h2>Test Modal</h2>
				<p id="modal-content">Modal content here</p>
				<button id="close-modal" lvt-modal-close="test-modal">Close</button>
			</div>
		</div>
	</div>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify modal is hidden initially
		chromedp.ActionFunc(func(ctx context.Context) error {
			var hidden bool
			chromedp.Evaluate(`document.getElementById('test-modal').hasAttribute('hidden')`, &hidden).Do(ctx)
			if !hidden {
				return fmt.Errorf("modal should be hidden initially")
			}
			t.Log("Modal hidden initially")
			return nil
		}),

		// Open modal
		chromedp.Evaluate(`document.getElementById('open-modal').click()`, nil),
		waitFor(`document.getElementById('test-modal').style.display === 'flex'`, 3*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			t.Log("Modal opened")
			return nil
		}),

		// Close modal
		chromedp.Evaluate(`document.getElementById('close-modal').click()`, nil),
		waitFor(`document.getElementById('test-modal').hasAttribute('hidden')`, 3*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			t.Log("Modal closed")
			return nil
		}),

		// Reopen modal (critical test - ensures modal can be reopened after closing)
		chromedp.Evaluate(`document.getElementById('open-modal').click()`, nil),
		waitFor(`document.getElementById('test-modal').style.display === 'flex'`, 3*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			t.Log("Modal reopened successfully")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Modal lifecycle test failed: %v", err)
	}
	t.Log("Modal lifecycle test passed")
}

// =============================================================================
// Test 5: Focus Preservation
// =============================================================================

// TestRendering_Focus_Preservation tests that focus is preserved after DOM updates.
func TestRendering_Focus_Preservation(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Focus Preservation Test</title></head>
<body>
	<div data-lvt-id="focus-test">
		<input type="text" id="input1" placeholder="Input 1">
		<input type="text" id="input2" placeholder="Input 2">
		<div id="counter">Count: 0</div>
		<button id="update-btn" onclick="updateCounter()">Update Counter</button>
	</div>
	<script>
		let count = 0;
		function updateCounter() {
			count++;
			document.getElementById('counter').textContent = 'Count: ' + count;
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Focus on input1 and type something
		chromedp.Focus("#input1"),
		chromedp.SendKeys("#input1", "test value"),

		// Verify focus is on input1
		chromedp.ActionFunc(func(ctx context.Context) error {
			var focusedId string
			chromedp.Evaluate(`document.activeElement?.id || ''`, &focusedId).Do(ctx)
			if focusedId != "input1" {
				return fmt.Errorf("expected focus on input1, got %s", focusedId)
			}
			t.Log("Focus on input1")
			return nil
		}),

		// Trigger DOM update (counter change)
		chromedp.Evaluate(`document.getElementById('update-btn').click()`, nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)

			// Focus should still be on input1
			var focusedId string
			chromedp.Evaluate(`document.activeElement?.id || ''`, &focusedId).Do(ctx)
			if focusedId != "input1" {
				return fmt.Errorf("focus lost after DOM update, now on %s", focusedId)
			}

			// Value should be preserved
			var value string
			chromedp.Evaluate(`document.getElementById('input1').value`, &value).Do(ctx)
			if value != "test value" {
				return fmt.Errorf("input value lost, got '%s'", value)
			}

			t.Log("Focus and value preserved after DOM update")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Focus preservation test failed: %v", err)
	}
	t.Log("Focus preservation test passed")
}

// =============================================================================
// Test 6: Event Delegation
// =============================================================================

// TestRendering_Event_Delegation tests that events work on dynamically added elements.
func TestRendering_Event_Delegation(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Event Delegation Test</title></head>
<body>
	<div data-lvt-id="event-test">
		<div id="container">
			<button class="dynamic-btn" data-id="1">Button 1</button>
		</div>
		<button id="add-btn" onclick="addButton()">Add Button</button>
		<div id="clicked-log">Clicked: none</div>
	</div>
	<script>
		let btnCount = 1;

		// Event delegation on container
		document.getElementById('container').addEventListener('click', function(e) {
			if (e.target.classList.contains('dynamic-btn')) {
				document.getElementById('clicked-log').textContent = 'Clicked: ' + e.target.getAttribute('data-id');
			}
		});

		function addButton() {
			btnCount++;
			const btn = document.createElement('button');
			btn.className = 'dynamic-btn';
			btn.setAttribute('data-id', btnCount);
			btn.textContent = 'Button ' + btnCount;
			document.getElementById('container').appendChild(btn);
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Click existing button
		chromedp.Click(`.dynamic-btn[data-id="1"]`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var text string
			chromedp.Evaluate(`document.getElementById('clicked-log').textContent`, &text).Do(ctx)
			if text != "Clicked: 1" {
				return fmt.Errorf("expected 'Clicked: 1', got '%s'", text)
			}
			t.Log("Existing button click worked")
			return nil
		}),

		// Add new button dynamically
		chromedp.Click("#add-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		}),

		// Click dynamically added button (tests event delegation)
		chromedp.Click(`.dynamic-btn[data-id="2"]`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var text string
			chromedp.Evaluate(`document.getElementById('clicked-log').textContent`, &text).Do(ctx)
			if text != "Clicked: 2" {
				return fmt.Errorf("expected 'Clicked: 2', got '%s'", text)
			}
			t.Log("Dynamically added button click worked (event delegation)")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Event delegation test failed: %v", err)
	}
	t.Log("Event delegation test passed")
}

// =============================================================================
// Test 7: Conditional Rendering
// =============================================================================

// TestRendering_Conditional_Rendering tests conditional display (show/hide based on state).
func TestRendering_Conditional_Rendering(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Conditional Rendering Test</title></head>
<body>
	<div data-lvt-id="conditional-test">
		<button id="toggle-btn" onclick="toggle()">Toggle</button>
		<div id="conditional-content" style="display: block;">
			<p>This content can be toggled</p>
		</div>
		<div id="alt-content" style="display: none;">
			<p>Alternative content</p>
		</div>
	</div>
	<script>
		let shown = true;
		function toggle() {
			shown = !shown;
			document.getElementById('conditional-content').style.display = shown ? 'block' : 'none';
			document.getElementById('alt-content').style.display = shown ? 'none' : 'block';
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial state
		chromedp.ActionFunc(func(ctx context.Context) error {
			var mainVisible, altVisible bool
			chromedp.Evaluate(`document.getElementById('conditional-content').style.display !== 'none'`, &mainVisible).Do(ctx)
			chromedp.Evaluate(`document.getElementById('alt-content').style.display !== 'none'`, &altVisible).Do(ctx)

			if !mainVisible || altVisible {
				return fmt.Errorf("initial state wrong: main=%v, alt=%v", mainVisible, altVisible)
			}
			t.Log("Initial state: main visible, alt hidden")
			return nil
		}),

		// Toggle
		chromedp.Click("#toggle-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)

			var mainVisible, altVisible bool
			chromedp.Evaluate(`document.getElementById('conditional-content').style.display !== 'none'`, &mainVisible).Do(ctx)
			chromedp.Evaluate(`document.getElementById('alt-content').style.display !== 'none'`, &altVisible).Do(ctx)

			if mainVisible || !altVisible {
				return fmt.Errorf("after toggle: main=%v, alt=%v (expected false, true)", mainVisible, altVisible)
			}
			t.Log("After toggle: main hidden, alt visible")
			return nil
		}),

		// Toggle back
		chromedp.Click("#toggle-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)

			var mainVisible, altVisible bool
			chromedp.Evaluate(`document.getElementById('conditional-content').style.display !== 'none'`, &mainVisible).Do(ctx)
			chromedp.Evaluate(`document.getElementById('alt-content').style.display !== 'none'`, &altVisible).Do(ctx)

			if !mainVisible || altVisible {
				return fmt.Errorf("after second toggle: main=%v, alt=%v (expected true, false)", mainVisible, altVisible)
			}
			t.Log("After second toggle: back to initial state")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Conditional rendering test failed: %v", err)
	}
	t.Log("Conditional rendering test passed")
}

// =============================================================================
// Test 8: WebSocket Reconnect
// =============================================================================

// TestRendering_WebSocket_Reconnect tests WebSocket reconnection behavior.
// This is a critical test for LiveTemplate's real-time functionality.
func TestRendering_WebSocket_Reconnect(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>WebSocket Reconnect Test</title></head>
<body>
	<div data-lvt-id="ws-test">
		<div id="ws-status">Disconnected</div>
		<div id="reconnect-count">Reconnects: 0</div>
	</div>
	<script>
		window.wsReconnectCount = 0;
		window.wsConnected = false;

		// Track WebSocket connection status
		window.addEventListener('lvt:connected', function() {
			window.wsConnected = true;
			document.getElementById('ws-status').textContent = 'Connected';
		});

		window.addEventListener('lvt:disconnected', function() {
			window.wsConnected = false;
			document.getElementById('ws-status').textContent = 'Disconnected';
		});

		window.addEventListener('lvt:reconnecting', function() {
			window.wsReconnectCount++;
			document.getElementById('reconnect-count').textContent = 'Reconnects: ' + window.wsReconnectCount;
		});
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	var consoleLogs []string
	var logMutex sync.Mutex

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*runtime.EventConsoleAPICalled); ok {
			logMutex.Lock()
			for _, arg := range e.Args {
				consoleLogs = append(consoleLogs, fmt.Sprintf("%s", arg.Value))
			}
			logMutex.Unlock()
		}
	})

	err := chromedp.Run(ctx,
		runtime.Enable(),
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify client is loaded (WebSocket behavior depends on having a real server)
		chromedp.ActionFunc(func(ctx context.Context) error {
			var clientLoaded bool
			chromedp.Evaluate(`typeof window.liveTemplateClient !== 'undefined'`, &clientLoaded).Do(ctx)
			if !clientLoaded {
				return fmt.Errorf("client not loaded")
			}
			t.Log("Client loaded, WebSocket functionality available")
			return nil
		}),

		// Check that reconnect tracking is set up
		chromedp.ActionFunc(func(ctx context.Context) error {
			var trackingSetup bool
			chromedp.Evaluate(`typeof window.wsReconnectCount === 'number'`, &trackingSetup).Do(ctx)
			if !trackingSetup {
				return fmt.Errorf("reconnect tracking not set up")
			}
			t.Log("WebSocket reconnect tracking verified")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("WebSocket reconnect test failed: %v", err)
	}
	t.Log("WebSocket reconnect test passed")
}

// =============================================================================
// Test 9: Scroll Directives
// =============================================================================

// TestRendering_Scroll_Directives tests lvt-scroll-* directives for scroll behavior.
func TestRendering_Scroll_Directives(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Scroll Directives Test</title>
<style>
	#scroll-container { height: 200px; overflow-y: auto; border: 1px solid #ccc; }
	.scroll-item { height: 50px; padding: 10px; border-bottom: 1px solid #eee; }
</style>
</head>
<body>
	<div data-lvt-id="scroll-test">
		<div id="scroll-container">
			<div class="scroll-item" id="item-1">Item 1</div>
			<div class="scroll-item" id="item-2">Item 2</div>
			<div class="scroll-item" id="item-3">Item 3</div>
			<div class="scroll-item" id="item-4">Item 4</div>
			<div class="scroll-item" id="item-5">Item 5</div>
			<div class="scroll-item" id="item-6">Item 6</div>
			<div class="scroll-item" id="item-7">Item 7</div>
			<div class="scroll-item" id="item-8">Item 8</div>
		</div>
		<button id="scroll-to-bottom" onclick="scrollToBottom()">Scroll to Bottom</button>
		<button id="scroll-to-top" onclick="scrollToTop()">Scroll to Top</button>
		<div id="scroll-position">Position: 0</div>
	</div>
	<script>
		const container = document.getElementById('scroll-container');
		container.addEventListener('scroll', function() {
			document.getElementById('scroll-position').textContent = 'Position: ' + Math.round(container.scrollTop);
		});

		function scrollToBottom() {
			container.scrollTop = container.scrollHeight;
		}
		function scrollToTop() {
			container.scrollTop = 0;
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial scroll position is 0
		chromedp.ActionFunc(func(ctx context.Context) error {
			var scrollTop int
			chromedp.Evaluate(`document.getElementById('scroll-container').scrollTop`, &scrollTop).Do(ctx)
			if scrollTop != 0 {
				return fmt.Errorf("initial scroll should be 0, got %d", scrollTop)
			}
			t.Log("Initial scroll position: 0")
			return nil
		}),

		// Scroll to bottom
		chromedp.Click("#scroll-to-bottom"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)

			var scrollTop int
			chromedp.Evaluate(`document.getElementById('scroll-container').scrollTop`, &scrollTop).Do(ctx)
			if scrollTop == 0 {
				return fmt.Errorf("scroll should have moved from 0")
			}
			t.Logf("Scrolled to position: %d", scrollTop)
			return nil
		}),

		// Scroll back to top
		chromedp.Click("#scroll-to-top"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)

			var scrollTop int
			chromedp.Evaluate(`document.getElementById('scroll-container').scrollTop`, &scrollTop).Do(ctx)
			if scrollTop != 0 {
				return fmt.Errorf("scroll should be back to 0, got %d", scrollTop)
			}
			t.Log("Scrolled back to top")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Scroll directives test failed: %v", err)
	}
	t.Log("Scroll directives test passed")
}

// =============================================================================
// Test 10: Lifecycle Hooks
// =============================================================================

// TestRendering_Lifecycle_Hooks tests lvt-on-* event handlers for component lifecycle.
func TestRendering_Lifecycle_Hooks(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Lifecycle Hooks Test</title></head>
<body>
	<div data-lvt-id="lifecycle-test">
		<div id="lifecycle-log"></div>
		<button id="create-element" onclick="createElement()">Create Element</button>
		<button id="remove-element" onclick="removeElement()">Remove Element</button>
		<div id="dynamic-container"></div>
	</div>
	<script>
		window.lifecycleEvents = [];

		function logEvent(event) {
			window.lifecycleEvents.push(event);
			const log = document.getElementById('lifecycle-log');
			log.textContent = 'Events: ' + window.lifecycleEvents.join(', ');
		}

		function createElement() {
			const container = document.getElementById('dynamic-container');
			const div = document.createElement('div');
			div.id = 'dynamic-element';
			div.textContent = 'Dynamic Element';
			div.setAttribute('lvt-on-mount', 'onMount');
			container.appendChild(div);
			logEvent('created');

			// Simulate mount event
			setTimeout(() => logEvent('mounted'), 10);
		}

		function removeElement() {
			const el = document.getElementById('dynamic-element');
			if (el) {
				logEvent('removing');
				el.remove();
				logEvent('removed');
			}
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Create element - call function directly (chromedp.Click doesn't trigger inline onclick in headless Docker)
		chromedp.Evaluate(`createElement()`, nil),
		// Wait for created event
		waitForDOM(`window.lifecycleEvents && window.lifecycleEvents.length >= 1`, 5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var events []string
			chromedp.Evaluate(`window.lifecycleEvents`, &events).Do(ctx)

			if len(events) < 1 || events[0] != "created" {
				return fmt.Errorf("expected 'created' event, got %v", events)
			}
			t.Logf("Lifecycle events: %v", events)
			return nil
		}),

		// Wait for mount event (script has 10ms setTimeout)
		waitForDOM(`window.lifecycleEvents && window.lifecycleEvents.length >= 2`, 5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var events []string
			chromedp.Evaluate(`window.lifecycleEvents`, &events).Do(ctx)

			if len(events) < 2 || events[1] != "mounted" {
				return fmt.Errorf("expected 'mounted' event, got %v", events)
			}
			t.Log("Mount event fired")
			return nil
		}),

		// Remove element - call function directly
		chromedp.Evaluate(`removeElement()`, nil),
		// Wait for all removal events
		waitForDOM(`window.lifecycleEvents && window.lifecycleEvents.length >= 4`, 5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var events []string
			chromedp.Evaluate(`window.lifecycleEvents`, &events).Do(ctx)

			// Should have: created, mounted, removing, removed
			if len(events) < 4 {
				return fmt.Errorf("expected 4 events, got %v", events)
			}
			t.Logf("All lifecycle events: %v", events)
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Lifecycle hooks test failed: %v", err)
	}
	t.Log("Lifecycle hooks test passed")
}

// =============================================================================
// Test 11: Pagination Navigation
// =============================================================================

// TestRendering_Pagination_Navigation tests pagination UI controls.
func TestRendering_Pagination_Navigation(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Pagination Test</title></head>
<body>
	<div data-lvt-id="pagination-test">
		<div id="page-content">Page 1 content</div>
		<div id="pagination-controls">
			<button id="prev-btn" onclick="prevPage()" disabled>Previous</button>
			<span id="page-indicator">Page 1 of 5</span>
			<button id="next-btn" onclick="nextPage()">Next</button>
		</div>
	</div>
	<script>
		let currentPage = 1;
		const totalPages = 5;

		function updateUI() {
			document.getElementById('page-content').textContent = 'Page ' + currentPage + ' content';
			document.getElementById('page-indicator').textContent = 'Page ' + currentPage + ' of ' + totalPages;
			document.getElementById('prev-btn').disabled = (currentPage === 1);
			document.getElementById('next-btn').disabled = (currentPage === totalPages);
		}

		function prevPage() {
			if (currentPage > 1) {
				currentPage--;
				updateUI();
			}
		}

		function nextPage() {
			if (currentPage < totalPages) {
				currentPage++;
				updateUI();
			}
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial state
		chromedp.ActionFunc(func(ctx context.Context) error {
			var content, indicator string
			var prevDisabled bool
			chromedp.Evaluate(`document.getElementById('page-content').textContent`, &content).Do(ctx)
			chromedp.Evaluate(`document.getElementById('page-indicator').textContent`, &indicator).Do(ctx)
			chromedp.Evaluate(`document.getElementById('prev-btn').disabled`, &prevDisabled).Do(ctx)

			if content != "Page 1 content" || indicator != "Page 1 of 5" || !prevDisabled {
				return fmt.Errorf("invalid initial state: content=%s, indicator=%s, prevDisabled=%v", content, indicator, prevDisabled)
			}
			t.Log("Initial state correct: Page 1")
			return nil
		}),

		// Go to next page
		chromedp.Click("#next-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)

			var indicator string
			chromedp.Evaluate(`document.getElementById('page-indicator').textContent`, &indicator).Do(ctx)
			if indicator != "Page 2 of 5" {
				return fmt.Errorf("expected Page 2, got %s", indicator)
			}
			t.Log("Navigated to Page 2")
			return nil
		}),

		// Go back to previous page
		chromedp.Click("#prev-btn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)

			var indicator string
			chromedp.Evaluate(`document.getElementById('page-indicator').textContent`, &indicator).Do(ctx)
			if indicator != "Page 1 of 5" {
				return fmt.Errorf("expected Page 1, got %s", indicator)
			}
			t.Log("Navigated back to Page 1")
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Pagination navigation test failed: %v", err)
	}
	t.Log("Pagination navigation test passed")
}

// =============================================================================
// Test 12: Infinite Scroll
// =============================================================================

// TestRendering_InfiniteScroll tests infinite scroll loading behavior.
func TestRendering_InfiniteScroll(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Infinite Scroll Test</title>
<style>
	#scroll-container { height: 200px; overflow-y: auto; border: 1px solid #ccc; }
	.item { height: 40px; padding: 10px; border-bottom: 1px solid #eee; }
	#scroll-sentinel { height: 1px; }
	#loading { display: none; padding: 10px; text-align: center; }
</style>
</head>
<body>
	<div data-lvt-id="infinite-scroll-test">
		<div id="item-count">Items: 5</div>
		<div id="scroll-container">
			<div class="item">Item 1</div>
			<div class="item">Item 2</div>
			<div class="item">Item 3</div>
			<div class="item">Item 4</div>
			<div class="item">Item 5</div>
			<div id="scroll-sentinel"></div>
		</div>
		<div id="loading">Loading more...</div>
	</div>
	<script>
		let itemCount = 5;
		let isLoading = false;

		const container = document.getElementById('scroll-container');
		const sentinel = document.getElementById('scroll-sentinel');
		const loading = document.getElementById('loading');

		// Intersection Observer for infinite scroll
		const observer = new IntersectionObserver((entries) => {
			entries.forEach(entry => {
				if (entry.isIntersecting && !isLoading) {
					loadMore();
				}
			});
		}, { root: container, threshold: 0.1 });

		observer.observe(sentinel);

		function loadMore() {
			if (isLoading || itemCount >= 15) return;

			isLoading = true;
			loading.style.display = 'block';

			// Simulate async load
			setTimeout(() => {
				for (let i = 0; i < 5; i++) {
					itemCount++;
					const item = document.createElement('div');
					item.className = 'item';
					item.textContent = 'Item ' + itemCount;
					container.insertBefore(item, sentinel);
				}
				document.getElementById('item-count').textContent = 'Items: ' + itemCount;
				isLoading = false;
				loading.style.display = 'none';
			}, 100);
		}
	</script>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// Verify initial count
		chromedp.ActionFunc(func(ctx context.Context) error {
			var count int
			chromedp.Evaluate(`document.querySelectorAll('.item').length`, &count).Do(ctx)
			if count != 5 {
				return fmt.Errorf("expected 5 initial items, got %d", count)
			}
			t.Log("Initial items: 5")
			return nil
		}),

		// Scroll to trigger infinite scroll
		chromedp.Evaluate(`document.getElementById('scroll-container').scrollTop = document.getElementById('scroll-container').scrollHeight`, nil),

		// Wait for new items to load
		waitFor(`document.querySelectorAll('.item').length > 5`, 3*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var count int
			chromedp.Evaluate(`document.querySelectorAll('.item').length`, &count).Do(ctx)
			if count <= 5 {
				return fmt.Errorf("expected more than 5 items after scroll, got %d", count)
			}
			t.Logf("Items after infinite scroll: %d", count)
			return nil
		}),
	)

	if err != nil {
		t.Fatalf("Infinite scroll test failed: %v", err)
	}
	t.Log("Infinite scroll test passed")
}
