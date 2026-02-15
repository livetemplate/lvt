//go:build browser

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/livetemplate/livetemplate"
	e2etest "github.com/livetemplate/lvt/testing"
)

// --- Controller + State for wire format test ---

type WireItem struct {
	ID   string
	Name string
	Done bool
}

type WireFormatState struct {
	Title   string
	Count   int
	Visible bool
	Items   []WireItem
}

type WireFormatController struct{}

func (c *WireFormatController) Increment(state WireFormatState, _ *livetemplate.Context) (WireFormatState, error) {
	state.Count++
	return state, nil
}

func (c *WireFormatController) Toggle(state WireFormatState, _ *livetemplate.Context) (WireFormatState, error) {
	state.Visible = !state.Visible
	return state, nil
}

func (c *WireFormatController) AddItem(state WireFormatState, ctx *livetemplate.Context) (WireFormatState, error) {
	item := WireItem{
		ID:   ctx.GetString("id"),
		Name: ctx.GetString("name"),
	}
	state.Items = append(state.Items, item)
	return state, nil
}

func (c *WireFormatController) RemoveItem(state WireFormatState, ctx *livetemplate.Context) (WireFormatState, error) {
	id := ctx.GetString("id")
	items := make([]WireItem, 0, len(state.Items))
	for _, item := range state.Items {
		if item.ID != id {
			items = append(items, item)
		}
	}
	state.Items = items
	return state, nil
}

func (c *WireFormatController) UpdateItem(state WireFormatState, ctx *livetemplate.Context) (WireFormatState, error) {
	id := ctx.GetString("id")
	name := ctx.GetString("name")
	for i, item := range state.Items {
		if item.ID == id {
			state.Items[i].Name = name
			break
		}
	}
	return state, nil
}

func (c *WireFormatController) ReorderItems(state WireFormatState, _ *livetemplate.Context) (WireFormatState, error) {
	// Reverse the items
	n := len(state.Items)
	reversed := make([]WireItem, n)
	for i, item := range state.Items {
		reversed[n-1-i] = item
	}
	state.Items = reversed
	return state, nil
}

// --- Template ---

const wireFormatTemplate = `<!DOCTYPE html>
<html>
<head><title>Wire Format Test</title></head>
<body>
	<h1 id="title">{{.Title}}</h1>
	<span id="count">{{.Count}}</span>
	{{if .Visible}}<div id="visible-section">Visible</div>{{end}}
	<ul id="item-list">
		{{range .Items}}
		<li data-key="{{.ID}}" class="item">{{.Name}}{{if .Done}} ✓{{end}}</li>
		{{end}}
	</ul>
	<button id="btn-increment" lvt-click="increment">+</button>
	<button id="btn-toggle" lvt-click="toggle">Toggle</button>
	<button id="btn-add" lvt-click="add_item" lvt-data-id="item-4" lvt-data-name="Delta">Add</button>
	<button id="btn-remove" lvt-click="remove_item" lvt-data-id="item-2">Remove</button>
	<button id="btn-update" lvt-click="update_item" lvt-data-id="item-1" lvt-data-name="Alpha Updated">Update</button>
	<button id="btn-reorder" lvt-click="reorder_items">Reorder</button>
	<script src="/client.js"></script>
</body>
</html>`

// --- Helpers ---

// hasStaticsAnywhere recursively checks if a tree contains any "s" key.
func hasStaticsAnywhere(v interface{}) bool {
	switch val := v.(type) {
	case map[string]interface{}:
		if _, ok := val["s"]; ok {
			return true
		}
		for _, child := range val {
			if hasStaticsAnywhere(child) {
				return true
			}
		}
	case []interface{}:
		for _, child := range val {
			if hasStaticsAnywhere(child) {
				return true
			}
		}
	}
	return false
}

// findRangeOps searches the tree for range operation arrays.
// Range ops are arrays like ["a", ...], ["i", ...], ["r", ...], ["u", ...], ["o", ...], ["p", ...].
func findRangeOps(v interface{}) [][]interface{} {
	var ops [][]interface{}

	switch val := v.(type) {
	case map[string]interface{}:
		for _, child := range val {
			ops = append(ops, findRangeOps(child)...)
		}
	case []interface{}:
		// Check if this is a range operation array
		if len(val) > 0 {
			if opStr, ok := val[0].(string); ok {
				switch opStr {
				case "a", "i", "p", "r", "u", "o":
					ops = append(ops, val)
					return ops
				}
			}
			// Could be an array of operations
			for _, child := range val {
				ops = append(ops, findRangeOps(child)...)
			}
		}
	}
	return ops
}

// findMetadata recursively searches a tree for range metadata (any "m" key).
func findMetadata(v interface{}) map[string]interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		if m, ok := val["m"]; ok {
			if mMap, ok := m.(map[string]interface{}); ok {
				return mMap
			}
		}
		for _, child := range val {
			if found := findMetadata(child); found != nil {
				return found
			}
		}
	case []interface{}:
		for _, child := range val {
			if found := findMetadata(child); found != nil {
				return found
			}
		}
	}
	return nil
}

// getReceivedMessagesWithTree returns all received WS messages that have a "tree" field.
func getReceivedMessagesWithTree(wsLogger *e2etest.WSMessageLogger) []e2etest.WSMessage {
	msgs := wsLogger.GetReceived()
	var result []e2etest.WSMessage
	for _, msg := range msgs {
		if msg.Parsed != nil {
			if _, ok := msg.Parsed["tree"]; ok {
				result = append(result, msg)
			}
		}
	}
	return result
}

// waitForActionResponse waits for a received WS message matching a specific action name.
// startFrom specifies the message index to start scanning from, so callers should
// snapshot len(wsLogger.GetReceived()) BEFORE triggering the action to avoid matching
// stale messages from previous actions.
func waitForActionResponse(wsLogger *e2etest.WSMessageLogger, action string, startFrom int, timeout time.Duration) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msgs := wsLogger.GetReceived()
		for i := startFrom; i < len(msgs); i++ {
			msg := msgs[i]
			if msg.Parsed == nil {
				continue
			}
			meta, ok := msg.Parsed["meta"]
			if !ok {
				continue
			}
			metaMap, ok := meta.(map[string]interface{})
			if !ok {
				continue
			}
			if metaMap["action"] == action {
				return msg.Parsed, nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, fmt.Errorf("timeout waiting for action response %q", action)
}

// --- Test ---

func TestWireFormat(t *testing.T) {
	controller := &WireFormatController{}
	state := &WireFormatState{
		Title:   "Wire Format Test",
		Count:   0,
		Visible: true,
		Items: []WireItem{
			{ID: "item-1", Name: "Alpha"},
			{ID: "item-2", Name: "Beta"},
			{ID: "item-3", Name: "Gamma"},
		},
	}

	tmpl := livetemplate.Must(livetemplate.New("wire-format"))
	if _, err := tmpl.Parse(wireFormatTemplate); err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tmpl.Handle(controller, livetemplate.AsState(state)))
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	port, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", port), 10*time.Second)

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			t.Logf("Server shutdown warning: %v", err)
		}
	}()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 60*time.Second)
	defer cleanup()

	ctx := chromeCtx.Context

	// Set up loggers for diagnostics
	wsLogger := e2etest.NewWSMessageLogger()
	consoleLogger := e2etest.NewConsoleLogger()
	wsLogger.Start(ctx)
	consoleLogger.Start(ctx)

	// Enable network domain to capture WebSocket frames
	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		t.Fatalf("Failed to enable network domain: %v", err)
	}

	url := e2etest.GetChromeTestURL(port)

	// Navigate and wait for initial render
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#title`, chromedp.ByID),
		e2etest.WaitForWebSocketReady(10*time.Second),
	)
	if err != nil {
		t.Fatalf("Navigation failed: %v", err)
	}

	// Dump diagnostics on failure
	defer func() {
		if t.Failed() {
			wsLogger.Print()
			consoleLogger.Print()
			var html string
			_ = chromedp.Run(ctx, chromedp.OuterHTML("body", &html, chromedp.ByQuery))
			t.Logf("Rendered HTML:\n%s", html)
		}
	}()

	// --- Subtest 1: Initial render has statics ---
	t.Run("1_Initial_Render_Has_Statics", func(t *testing.T) {
		treeMsgs := getReceivedMessagesWithTree(wsLogger)
		if len(treeMsgs) == 0 {
			t.Fatal("No tree messages received after initial render")
		}

		initialMsg := treeMsgs[0]
		tree, ok := initialMsg.Parsed["tree"]
		if !ok {
			t.Fatal("Initial message has no 'tree' field")
		}

		if !hasStaticsAnywhere(tree) {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Initial render tree should contain statics ('s' keys), got:\n%s", string(prettyTree))
		}

		// Also verify DOM
		var title string
		var count string
		var itemCount int
		err := chromedp.Run(ctx,
			chromedp.Text(`#title`, &title, chromedp.ByID),
			chromedp.Text(`#count`, &count, chromedp.ByID),
			chromedp.Evaluate(`document.querySelectorAll('.item').length`, &itemCount),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}

		if title != "Wire Format Test" {
			t.Errorf("Title = %q, want %q", title, "Wire Format Test")
		}
		if count != "0" {
			t.Errorf("Count = %q, want %q", count, "0")
		}
		if itemCount != 3 {
			t.Errorf("Item count = %d, want 3", itemCount)
		}

		t.Log("Initial render has statics and DOM is correct")
	})

	// --- Subtest 2: Update omits statics ---
	t.Run("2_Update_Omits_Statics", func(t *testing.T) {
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-increment`, chromedp.ByID),
			e2etest.WaitFor(`document.getElementById('count').textContent === '1'`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "increment", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No update message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Update message has no 'tree' field")
		}

		if hasStaticsAnywhere(tree) {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Update tree should NOT contain statics ('s' keys), got:\n%s", string(prettyTree))
		}

		// Verify DOM
		var count string
		if err := chromedp.Run(ctx, chromedp.Text(`#count`, &count, chromedp.ByID)); err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if count != "1" {
			t.Errorf("Count = %q, want %q", count, "1")
		}

		t.Log("Update correctly omits statics")
	})

	// --- Subtest 3: Range insert ---
	t.Run("3_Range_Insert", func(t *testing.T) {
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-add`, chromedp.ByID),
			e2etest.WaitFor(`document.querySelectorAll('.item').length === 4`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "add_item", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No insert message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		ops := findRangeOps(tree)
		foundAppend := false
		for _, op := range ops {
			if len(op) > 0 {
				if opStr, ok := op[0].(string); ok && (opStr == "a" || opStr == "i" || opStr == "p") {
					foundAppend = true
					break
				}
			}
		}

		if !foundAppend {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Expected append/insert operation ['a'/'i'/'p', ...] in tree, got:\n%s", string(prettyTree))
		}

		// Verify DOM: 4 items, last one is Delta
		var lastItemText string
		err = chromedp.Run(ctx,
			chromedp.Text(`.item:last-child`, &lastItemText, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if !strings.Contains(lastItemText, "Delta") {
			t.Errorf("Last item text = %q, want to contain 'Delta'", lastItemText)
		}

		t.Log("Range append/insert operation found and DOM updated")
	})

	// --- Subtest 4: Range remove ---
	t.Run("4_Range_Remove", func(t *testing.T) {
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-remove`, chromedp.ByID),
			e2etest.WaitFor(`document.querySelectorAll('.item').length === 3`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "remove_item", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No remove message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		ops := findRangeOps(tree)
		foundRemove := false
		for _, op := range ops {
			if len(op) >= 2 {
				if opStr, ok := op[0].(string); ok && opStr == "r" {
					foundRemove = true
					break
				}
			}
		}

		if !foundRemove {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Expected remove operation ['r', ...] in tree, got:\n%s", string(prettyTree))
		}

		// Verify DOM: item-2 (Beta) should be gone
		var hasItem2 bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('[data-key="item-2"]') !== null`, &hasItem2),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if hasItem2 {
			t.Error("item-2 (Beta) should have been removed from DOM")
		}

		t.Log("Range remove operation found and DOM updated")
	})

	// --- Subtest 5: Range update ---
	t.Run("5_Range_Update", func(t *testing.T) {
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-update`, chromedp.ByID),
			e2etest.WaitFor(`
				(() => {
					const el = document.querySelector('[data-key="item-1"]');
					return el && el.textContent.includes('Alpha Updated');
				})()
			`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "update_item", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No update message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		// The update_item response may contain the ["u", ...] operation,
		// OR the tree may be empty {} if the diff engine already batched
		// the update into a prior response (e.g., remove_item).
		// Either way is valid — check both the wire format and the DOM.
		ops := findRangeOps(tree)
		foundUpdate := false
		for _, op := range ops {
			if len(op) >= 2 {
				if opStr, ok := op[0].(string); ok && opStr == "u" {
					foundUpdate = true
					break
				}
			}
		}

		// Also check if the remove_item response contained the update
		// (diff engine may batch operations). Search all historical messages
		// since the remove_item response arrived before this subtest.
		if !foundUpdate {
			for _, msg := range wsLogger.GetReceived() {
				if msg.Parsed == nil {
					continue
				}
				meta, _ := msg.Parsed["meta"].(map[string]interface{})
				if meta == nil || meta["action"] != "remove_item" {
					continue
				}
				if removeTree, ok := msg.Parsed["tree"]; ok {
					for _, op := range findRangeOps(removeTree) {
						if len(op) >= 2 {
							if opStr, ok := op[0].(string); ok && opStr == "u" {
								foundUpdate = true
								break
							}
						}
					}
				}
			}
		}

		if !foundUpdate {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Expected update operation ['u', ...] in either update_item or remove_item response, got update_item tree:\n%s", string(prettyTree))
		}

		// Verify DOM — this is the ground truth
		var itemText string
		err = chromedp.Run(ctx,
			chromedp.Text(`[data-key="item-1"]`, &itemText, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if !strings.Contains(itemText, "Alpha Updated") {
			t.Errorf("Item text = %q, want to contain 'Alpha Updated'", itemText)
		}

		t.Log("Range update operation found and DOM updated")
	})

	// --- Subtest 6: Range reorder ---
	t.Run("6_Range_Reorder", func(t *testing.T) {
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-reorder`, chromedp.ByID),
			// After reorder: Delta, Gamma, Alpha Updated (reversed)
			e2etest.WaitFor(`
				(() => {
					const items = document.querySelectorAll('.item');
					return items.length === 3 && items[0].textContent.includes('Delta');
				})()
			`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "reorder_items", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No reorder message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		ops := findRangeOps(tree)
		foundReorder := false
		for _, op := range ops {
			if len(op) >= 2 {
				if opStr, ok := op[0].(string); ok && opStr == "o" {
					foundReorder = true
					break
				}
			}
		}

		if !foundReorder {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Expected reorder operation ['o', [...]] in tree, got:\n%s", string(prettyTree))
		}

		// Verify DOM order: first item should be Delta, last should be Alpha Updated
		var firstItemText, lastItemText string
		err = chromedp.Run(ctx,
			chromedp.Text(`.item:first-child`, &firstItemText, chromedp.ByQuery),
			chromedp.Text(`.item:last-child`, &lastItemText, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if !strings.Contains(firstItemText, "Delta") {
			t.Errorf("First item = %q, want to contain 'Delta'", firstItemText)
		}
		if !strings.Contains(lastItemText, "Alpha Updated") {
			t.Errorf("Last item = %q, want to contain 'Alpha Updated'", lastItemText)
		}

		t.Log("Range reorder operation found and DOM updated")
	})

	// --- Subtest 7: Conditional toggle ---
	t.Run("7_Conditional_Toggle", func(t *testing.T) {
		// Verify visible-section exists before toggle
		var visibleBefore bool
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`document.getElementById('visible-section') !== null`, &visibleBefore),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if !visibleBefore {
			t.Fatal("visible-section should exist before toggle")
		}

		msgCount := len(wsLogger.GetReceived())
		err = chromedp.Run(ctx,
			chromedp.Click(`#btn-toggle`, chromedp.ByID),
			e2etest.WaitFor(`document.getElementById('visible-section') === null`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "toggle", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No toggle message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		// The toggle message should not have statics (it's an update)
		if hasStaticsAnywhere(tree) {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Toggle update should not contain statics, got:\n%s", string(prettyTree))
		}

		// Verify DOM: visible-section should be gone
		var visibleAfter bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.getElementById('visible-section') !== null`, &visibleAfter),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if visibleAfter {
			t.Error("visible-section should not exist after toggle")
		}

		t.Log("Conditional toggle correctly hides section")
	})

	// --- Subtest 8: Structure change resends statics ---
	t.Run("8_Structure_Change_Resends_Statics", func(t *testing.T) {
		// After subtest 7, Visible is false (conditional hidden).
		// Toggling back to true changes the tree structure: the conditional
		// position goes from "" to a TreeNode with statics. The server detects
		// a fingerprint mismatch and must resend statics for the new subtree.
		msgCount := len(wsLogger.GetReceived())
		err := chromedp.Run(ctx,
			chromedp.Click(`#btn-toggle`, chromedp.ByID),
			e2etest.WaitFor(`document.getElementById('visible-section') !== null`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Click/wait failed: %v", err)
		}

		msg, err := waitForActionResponse(wsLogger, "toggle", msgCount, 5*time.Second)
		if err != nil {
			t.Fatalf("No toggle message received: %v", err)
		}

		tree, ok := msg["tree"]
		if !ok {
			t.Fatal("Message has no 'tree' field")
		}

		if !hasStaticsAnywhere(tree) {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Structure change (conditional reappearing) should resend statics, got:\n%s", string(prettyTree))
		}

		var visibleAfter bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.getElementById('visible-section') !== null`, &visibleAfter),
		)
		if err != nil {
			t.Fatalf("DOM query failed: %v", err)
		}
		if !visibleAfter {
			t.Error("visible-section should exist after toggling back to visible")
		}

		t.Log("Structure change correctly resends statics")
	})

	// --- Subtest 9: Range metadata idKey ---
	t.Run("9_Range_Metadata_IDKey", func(t *testing.T) {
		// The initial render tree must contain range metadata with idKey.
		// Without idKey, the client cannot match items for differential
		// operations (update, remove, reorder).
		treeMsgs := getReceivedMessagesWithTree(wsLogger)
		if len(treeMsgs) == 0 {
			t.Fatal("No tree messages received")
		}

		initialMsg := treeMsgs[0]
		tree, ok := initialMsg.Parsed["tree"]
		if !ok {
			t.Fatal("Initial message has no 'tree' field")
		}

		metadata := findMetadata(tree)
		if metadata == nil {
			prettyTree, _ := json.MarshalIndent(tree, "", "  ")
			t.Fatalf("Initial render should contain range metadata ('m' key), got:\n%s", string(prettyTree))
		}

		idKey, ok := metadata["idKey"]
		if !ok {
			t.Fatalf("Range metadata should contain 'idKey', got: %v", metadata)
		}

		idKeyStr, ok := idKey.(string)
		if !ok || idKeyStr == "" {
			t.Fatalf("idKey should be a non-empty string, got: %v", idKey)
		}

		t.Logf("Range metadata found with idKey=%q", idKeyStr)
	})

	// --- Subtest 10: Envelope schema validation ---
	t.Run("10_Envelope_Schema_Validation", func(t *testing.T) {
		// Validate UpdateResponse envelope structure across all captured messages.
		// Every message with a "tree" field must have a "meta" object with
		// a boolean "success" field and optional string "action" / map "errors".
		treeMsgs := getReceivedMessagesWithTree(wsLogger)
		if len(treeMsgs) == 0 {
			t.Fatal("No tree messages to validate")
		}

		for i, msg := range treeMsgs {
			meta, ok := msg.Parsed["meta"]
			if !ok {
				t.Errorf("Message %d: missing 'meta' field", i)
				continue
			}

			metaMap, ok := meta.(map[string]interface{})
			if !ok {
				t.Errorf("Message %d: 'meta' is not a map, got %T", i, meta)
				continue
			}

			success, ok := metaMap["success"]
			if !ok {
				t.Errorf("Message %d: meta missing 'success' field", i)
			} else if _, ok := success.(bool); !ok {
				t.Errorf("Message %d: meta.success is %T, want bool", i, success)
			}

			if action, ok := metaMap["action"]; ok {
				if _, ok := action.(string); !ok {
					t.Errorf("Message %d: meta.action is %T, want string", i, action)
				}
			}

			if errors, ok := metaMap["errors"]; ok && errors != nil {
				if _, ok := errors.(map[string]interface{}); !ok {
					t.Errorf("Message %d: meta.errors is %T, want map or nil", i, errors)
				}
			}
		}

		t.Logf("Validated envelope schema for %d messages", len(treeMsgs))
	})
}
