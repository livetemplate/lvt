//go:build browser

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestTutorialE2E tests the complete blog tutorial workflow
func TestTutorialE2E(t *testing.T) {
	t.Parallel() // Can run concurrently with Chrome pool

	// Create temp directory for test blog
	tmpDir := t.TempDir()
	blogDir := filepath.Join(tmpDir, "testblog")

	// Step 1: lvt new testblog (production mode for Docker compatibility)
	t.Log("Step 1: Creating new blog app...")
	if err := runLvtCommand(t, tmpDir, "new", "testblog"); err != nil {
		t.Fatalf("Failed to create new app: %v", err)
	}
	t.Log("✅ Blog app created")

	// Enable DevMode BEFORE generating resources so DevMode=true gets baked into handler code
	enableDevMode(t, blogDir)

	// Step 2: Generate posts resource
	t.Log("Step 2: Generating posts resource...")
	if err := runLvtCommand(t, blogDir, "gen", "resource", "posts", "title", "content", "published:bool"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}
	t.Log("✅ Posts resource generated")

	// Step 3: Generate categories resource
	t.Log("Step 3: Generating categories resource...")
	if err := runLvtCommand(t, blogDir, "gen", "resource", "categories", "name", "description"); err != nil {
		t.Fatalf("Failed to generate categories: %v", err)
	}
	t.Log("✅ Categories resource generated")

	// Step 4: Generate comments resource with foreign key
	t.Log("Step 4: Generating comments resource with FK...")
	if err := runLvtCommand(t, blogDir, "gen", "resource", "comments", "post_id:references:posts", "author", "text"); err != nil {
		t.Fatalf("Failed to generate comments: %v", err)
	}
	t.Log("✅ Comments resource generated with foreign key")

	// Step 5: Run migrations
	t.Log("Step 5: Running migrations...")
	if err := runLvtCommand(t, blogDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	t.Log("✅ Migrations complete")

	// Verify foreign key in migration file
	t.Log("Verifying foreign key syntax...")
	migrationsDir := filepath.Join(blogDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations dir: %v", err)
	}

	var commentsMigration string
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "comments") {
			data, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read migration: %v", err)
			}
			commentsMigration = string(data)
			break
		}
	}

	// Verify inline FOREIGN KEY (not ALTER TABLE)
	if strings.Contains(commentsMigration, "ALTER TABLE") && strings.Contains(commentsMigration, "ADD CONSTRAINT") {
		t.Error("❌ Migration uses ALTER TABLE ADD CONSTRAINT (should use inline FOREIGN KEY)")
	} else if strings.Contains(commentsMigration, "FOREIGN KEY (post_id) REFERENCES posts(id)") {
		t.Log("✅ Foreign key uses correct inline syntax")
	} else {
		t.Error("❌ Foreign key definition not found in migration")
	}

	// Step 6 & 7: Build and run app natively (much faster than Docker)
	// This test focuses on UI functionality, not deployment - Docker testing is in deployment_docker_test.go
	serverPort := allocateTestPort() // Use unique port for parallel testing
	_ = buildAndRunNative(t, blogDir, serverPort)

	serverURL := fmt.Sprintf("http://localhost:%d", serverPort)

	// Step 9: E2E UI Testing with Chrome
	t.Log("Step 9: Testing UI with Chrome...")

	// Use Chrome from pool for parallel execution
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	// Determine URL for Chrome to access (Docker networking)
	testURL := getTestURL(serverPort)

	// Listen for console errors (especially WebSocket errors)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if consoleEvent, ok := ev.(*runtime.EventConsoleAPICalled); ok {
			for _, arg := range consoleEvent.Args {
				if arg.Type == runtime.TypeString {
					logMsg := string(arg.Value)
					if strings.Contains(logMsg, "WebSocket") || strings.Contains(logMsg, "Failed") || strings.Contains(logMsg, "Error") {
						t.Logf("Browser console: %s", logMsg)
					}
				}
			}
		}
	})

	// Test WebSocket Connection
	t.Run("WebSocket Connection", func(t *testing.T) {
		// Create fresh browser context for this subtest
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		// Note: Apps use CDN client from unpkg.com, not local /livetemplate-client.js

		var wsConnected bool
		var wsURL string
		var wsReadyState int
		var pathname string
		var liveUrl string
		err = chromedp.Run(testCtx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(30*time.Second), // Wait for WebSocket init and first update
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"), // Validate no raw template expressions
			chromedp.Evaluate(`window.location.pathname`, &pathname),
			chromedp.Evaluate(`window.liveTemplateClient ? window.liveTemplateClient.options.liveUrl : null`, &liveUrl),
			chromedp.Evaluate(`(() => {
				// Get WebSocket URL being used
				return window.liveTemplateClient && window.liveTemplateClient.ws ? window.liveTemplateClient.ws.url : null;
			})()`, &wsURL),
			chromedp.Evaluate(`(() => {
				// Get WebSocket readyState
				return window.liveTemplateClient && window.liveTemplateClient.ws ? window.liveTemplateClient.ws.readyState : -1;
			})()`, &wsReadyState),
			chromedp.Evaluate(`(() => {
				// Check if WebSocket connection exists
				return window.liveTemplateClient &&
				       window.liveTemplateClient.ws &&
				       window.liveTemplateClient.ws.readyState === WebSocket.OPEN;
			})()`, &wsConnected),
		)
		if err != nil {
			t.Fatalf("Failed to check WebSocket connection: %v", err)
		}

		t.Logf("window.location.pathname: %s", pathname)
		t.Logf("client.options.liveUrl: %s", liveUrl)
		t.Logf("WebSocket URL: %s, ReadyState: %d (0=CONNECTING, 1=OPEN, 2=CLOSING, 3=CLOSED)", wsURL, wsReadyState)

		if !wsConnected {
			t.Error("❌ WebSocket did not connect to /posts endpoint")
		} else {
			t.Log("✅ WebSocket connected successfully to " + wsURL)
		}
	})

	// Test /posts Endpoint Serves Content
	t.Run("Posts Page", func(t *testing.T) {
		// Create fresh browser context for this subtest
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		var lvtId string
		err := chromedp.Run(testCtx,
			chromedp.Navigate(testURL+"/posts"),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			chromedp.AttributeValue(`[data-lvt-id]`, "data-lvt-id", &lvtId, nil),
		)
		if err != nil {
			t.Fatalf("Failed to test /posts endpoint: %v", err)
		}

		if lvtId == "" {
			t.Error("❌ LiveTemplate wrapper not found on /posts endpoint")
		} else {
			t.Logf("✅ /posts endpoint serves LiveTemplate content (wrapper ID: %s)", lvtId)
		}
	})

	// Test Add Post
	t.Run("Add Post", func(t *testing.T) {
		// Create per-subtest context with individual timeout
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		var titleValueBeforeSubmit string
		var contentValueBeforeSubmit string
		var activeElementAfterTyping string
		var lastFocusedName string
		err := chromedp.Run(testCtx,
			// Navigate to /posts and wait for it to load
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(30*time.Second), // Wait for WebSocket init and first update
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"), // Validate no raw template expressions

			// Click the "+ Add Posts" button in toolbar to open modal
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			// Wait for modal to appear
			waitFor(`document.querySelector('[role="dialog"]') && !document.querySelector('[role="dialog"]').hasAttribute('hidden')`, 10*time.Second),

			// Fill in the form in the modal
			chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="title"]`, "My First Blog Post", chromedp.ByQuery),
			chromedp.SendKeys(`textarea[name="content"]`, "This is the content of my first blog post", chromedp.ByQuery),
			chromedp.Evaluate(`document.querySelector('input[name="title"]').value`, &titleValueBeforeSubmit),
			chromedp.Evaluate(`document.querySelector('textarea[name="content"]').value`, &contentValueBeforeSubmit),
			chromedp.Evaluate(`document.activeElement ? document.activeElement.getAttribute('name') || document.activeElement.tagName : 'none'`, &activeElementAfterTyping),
			chromedp.Click(`input[name="published"]`, chromedp.ByQuery),

			// Click the submit button in the modal
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),

			// Wait for the post to appear in the table
			waitFor(`
				(() => {
					const table = document.querySelector('table tbody');
					return table && table.querySelectorAll('tr').length > 0;
				})()
			`, 30*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			// Wait for page to load
			waitFor(`document.readyState === 'complete'`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to add post: %v", err)
		}

		// Debug: Capture the tree state from the client
		var treeState string
		var lastWSMessage string
		var lastFormData string
		_ = chromedp.Run(testCtx,
			chromedp.Evaluate(`JSON.stringify(window.liveTemplateClient?.treeState || {}, null, 2)`, &treeState),
			chromedp.Evaluate(`window.__lastWSMessage || 'No WS message captured'`, &lastWSMessage),
			chromedp.Evaluate(`JSON.stringify(window.__lvtLastFormData || JSON.parse(window.sessionStorage.getItem("__lvtLastFormData") || "null"))`, &lastFormData),
			chromedp.Evaluate(`window.__lvtLastFocusedName || 'unknown'`, &lastFocusedName),
		)
		t.Logf("=== CLIENT TREE STATE ===\n%s\n", treeState)
		t.Logf("=== LAST WS MESSAGE ===\n%s\n", lastWSMessage)
		t.Logf("=== LAST FORM DATA ===\n%s\n", lastFormData)
		t.Logf("=== VALUES BEFORE SUBMIT === title=%q content=%q activeElement=%q lastFocused=%q\n", titleValueBeforeSubmit, contentValueBeforeSubmit, activeElementAfterTyping, lastFocusedName)

		// Verify the post appears in the table
		var postInTable bool
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My First Blog Post';
					});
				})()
			`, &postInTable),
		)
		if err != nil {
			t.Fatalf("Failed to check table: %v", err)
		}

		if !postInTable {
			var tableSummary string
			var firstCellHTML string
			_ = chromedp.Run(testCtx,
				chromedp.Evaluate(`
					(() => {
						const table = document.querySelector('table');
						if (!table) {
							const wrapper = document.querySelector('[data-lvt-id]');
							return 'No table found. Wrapper exists: ' + !!wrapper + '. Body text: ' + document.body.textContent.substring(0, 200);
						}
						const rows = Array.from(table.querySelectorAll('tbody tr'));
						return 'Table has ' + rows.length + ' rows. Titles: ' + rows.map(r => {
							const cells = r.querySelectorAll('td');
							return cells.length > 0 ? cells[0].textContent.trim() : '';
						}).join(', ');
					})()
				`, &tableSummary),
				chromedp.Evaluate(`
					(() => {
						const table = document.querySelector('table');
						if (!table) return 'No table';
						const firstRow = table.querySelector('tbody tr');
						if (!firstRow) return 'No rows';
						const firstCell = firstRow.querySelector('td');
						if (!firstCell) return 'No cell';
						return firstCell.innerHTML;
					})()
				`, &firstCellHTML),
			)
			t.Fatalf("❌ Post not found in table.\nTable summary: %s\nFirst cell HTML: %s", tableSummary, firstCellHTML)
		}

		t.Log("✅ Post 'My First Blog Post' added successfully and appears in table")
	})

	// Test Modal Delete with Confirmation
	t.Run("Modal Delete with Confirmation", func(t *testing.T) {
		// Create per-subtest context with individual timeout
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		// First, ensure a post exists for this test (make test independent)
		t.Log("Creating post fixture for modal delete test")
		if err := ensureTutorialPostExists(testCtx, testURL); err != nil {
			t.Fatalf("Failed to create post fixture: %v", err)
		}

		// Verify the post exists
		var postExists bool
		err := chromedp.Run(testCtx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(30*time.Second), // Wait for WebSocket init and first update
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"), // Validate no raw template expressions
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My First Blog Post';
					});
				})()
			`, &postExists),
		)
		if err != nil {
			t.Fatalf("Failed to check for post: %v", err)
		}

		if !postExists {
			t.Fatal("❌ Post 'My First Blog Post' not found - cannot test deletion")
		}

		// Verify there's NO delete button in table rows (modal mode)
		var deleteButtonInRow bool
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My First Blog Post';
					});
					if (targetRow) {
						const deleteButton = targetRow.querySelector('button[lvt-click="delete"]');
						return !!deleteButton;
					}
					return false;
				})()
			`, &deleteButtonInRow),
		)
		if err != nil {
			t.Fatalf("Failed to check for delete button in row: %v", err)
		}

		if deleteButtonInRow {
			t.Error("❌ Delete button should NOT be in table rows in modal mode")
		} else {
			t.Log("✅ No delete button in table rows (modal mode)")
		}

		// Click Edit button to open modal
		var editButtonFound bool
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My First Blog Post';
					});
					if (targetRow) {
						const editButton = targetRow.querySelector('button[lvt-click="edit"]');
						if (editButton) {
							editButton.click();
							return true;
						}
					}
					return false;
				})()
			`, &editButtonFound),
			// Wait for EDIT modal to open (has lvt-submit="update", not "add")
			waitFor(`document.querySelector('form[lvt-submit="update"]') !== null`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to click edit button: %v", err)
		}

		if !editButtonFound {
			t.Fatal("❌ Edit button not found in table row")
		}
		t.Log("✅ Edit button clicked, modal should be open")

		// Verify delete button exists in modal with lvt-confirm attribute
		var deleteButtonInModal bool
		var hasConfirmAttr bool
		var modalHTML string
		err = chromedp.Run(testCtx,
			// Capture the edit form HTML for debugging (not add modal)
			chromedp.Evaluate(`
				(() => {
					const editForm = document.querySelector('form[lvt-submit="update"]');
					return editForm ? editForm.outerHTML : 'Edit form not found';
				})()
			`, &modalHTML),
			chromedp.Evaluate(`
				(() => {
					const deleteButton = document.querySelector('button[lvt-click="delete"]');
					return !!deleteButton;
				})()
			`, &deleteButtonInModal),
			chromedp.Evaluate(`
				(() => {
					const deleteButton = document.querySelector('button[lvt-click="delete"]');
					return deleteButton ? deleteButton.hasAttribute('lvt-confirm') : false;
				})()
			`, &hasConfirmAttr),
		)
		if err != nil {
			t.Fatalf("Failed to check delete button in modal: %v", err)
		}

		if !deleteButtonInModal {
			t.Logf("❌ Modal HTML:\n%s", modalHTML)
			t.Fatal("❌ Delete button not found in modal")
		}
		t.Log("✅ Delete button found in modal")

		if !hasConfirmAttr {
			t.Error("❌ Delete button missing lvt-confirm attribute")
		} else {
			t.Log("✅ Delete button has lvt-confirm attribute")
		}

		// We've already verified the key requirements:
		// 1. No delete button in table rows ✅
		// 2. Delete button exists in modal ✅
		// 3. Delete button has lvt-confirm attribute ✅
		// The confirmation dialog functionality is client-side JavaScript (window.confirm)
		// which is tested implicitly through the next test case

		t.Log("✅ Modal delete confirmation setup verified")
	})

	// Test Delete Post (Accept Confirmation)
	t.Run("Delete Post with Accepted Confirmation", func(t *testing.T) {
		// Create per-subtest context with individual timeout
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()

		// Add console log listener for this specific test
		chromedp.ListenTarget(testCtx, func(ev interface{}) {
			if consoleEvent, ok := ev.(*runtime.EventConsoleAPICalled); ok {
				for _, arg := range consoleEvent.Args {
					if arg.Type == runtime.TypeString {
						logMsg := string(arg.Value)
						// Log ALL console messages for debugging
						t.Logf("🔍 Browser console: %s", logMsg)
					}
				}
			}
		})

		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		const (
			defaultPostTitle   = "My First Blog Post"
			findExistingPostJS = `
		(() => {
			const titles = ['My Updated Blog Post', 'My First Blog Post'];
			const table = document.querySelector('table');
			if (!table) return '';
			const rows = Array.from(table.querySelectorAll('tbody tr'));
			for (const title of titles) {
				if (rows.some(row => {
					const cells = row.querySelectorAll('td');
					return cells.length > 0 && cells[0].textContent.trim() === title;
				})) {
					return title;
				}
			}
			return '';
		})()
		`
		)

		var targetTitle string
		err := chromedp.Run(testCtx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(30*time.Second), // Wait for WebSocket init and first update
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"), // Validate no raw template expressions
			chromedp.Evaluate(findExistingPostJS, &targetTitle),
		)
		if err != nil {
			t.Fatalf("Failed to check for existing post: %v", err)
		}

		if targetTitle == "" {
			t.Log("ℹ️ Post not found, creating fixture for deletion test")
			if err := ensureTutorialPostExists(testCtx, testURL); err != nil {
				t.Fatalf("Failed to create post fixture: %v", err)
			}
			targetTitle = defaultPostTitle
			postVisibleJS := fmt.Sprintf(`
		(() => {
			const table = document.querySelector('table');
			if (!table) return false;
			const rows = Array.from(table.querySelectorAll('tbody tr'));
			return rows.some(row => {
				const cells = row.querySelectorAll('td');
				return cells.length > 0 && cells[0].textContent.trim() === %q;
			});
		})()
		`, targetTitle)
			var postVisible bool
			if err := chromedp.Run(testCtx,
				waitFor(postVisibleJS, 30*time.Second),
				chromedp.Evaluate(postVisibleJS, &postVisible),
			); err != nil {
				t.Fatalf("Failed to verify post fixture: %v", err)
			}
			if !postVisible {
				t.Fatal("❌ Unable to create post fixture for deletion test")
			}
		}

		if targetTitle == "" {
			t.Fatal("❌ No post available for deletion test")
		}

		checkPostExistsJS := fmt.Sprintf(`
		(() => {
			const table = document.querySelector('table');
			if (!table) return false;
			const rows = Array.from(table.querySelectorAll('tbody tr'));
			return rows.some(row => {
				const cells = row.querySelectorAll('td');
				return cells.length > 0 && cells[0].textContent.trim() === %q;
			});
		})()
		`, targetTitle)

		var postExists bool
		if err := chromedp.Run(testCtx,
			chromedp.Evaluate(checkPostExistsJS, &postExists),
		); err != nil {
			t.Fatalf("Failed to verify post before deletion: %v", err)
		}
		if !postExists {
			t.Fatalf("❌ Target post %q not found before deletion", targetTitle)
		}

		// Find the specific data-key of the row we're going to delete
		var targetDataKey string
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return '';
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === %q;
					});
					return targetRow ? targetRow.getAttribute('data-key') : '';
				})()
			`, targetTitle), &targetDataKey),
		)
		if err != nil {
			t.Fatalf("Failed to find target row data-key: %v", err)
		}
		if targetDataKey == "" {
			t.Fatalf("❌ Could not find data-key for post %q", targetTitle)
		}
		t.Logf("🎯 Target post data-key: %s", targetDataKey)

		// Click Edit button to open modal
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
					return cells.length > 0 && cells[0].textContent.trim() === %q;
					});
					if (targetRow) {
						const editButton = targetRow.querySelector('button[lvt-click="edit"]');
						if (editButton) {
							editButton.click();
							return true;
						}
					}
					return false;
			})()
			`, targetTitle), &postExists),
			// Wait for modal controls to be ready before continuing
			waitFor(`document.querySelector('button[lvt-click="delete"]') !== null`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to open edit modal: %v", err)
		}

		// CAPTURE STATE BEFORE CLICKING DELETE (so we can see the starting state)
		var screenshotBefore []byte
		var htmlBefore string
		err = chromedp.Run(testCtx,
			chromedp.FullScreenshot(&screenshotBefore, 100),
			chromedp.OuterHTML("html", &htmlBefore),
		)
		if err == nil {
			beforePath := fmt.Sprintf("/tmp/delete_before_%d.png", time.Now().Unix())
			beforeHTMLPath := fmt.Sprintf("/tmp/delete_before_%d.html", time.Now().Unix())
			os.WriteFile(beforePath, screenshotBefore, 0644)
			os.WriteFile(beforeHTMLPath, []byte(htmlBefore), 0644)
			t.Logf("💾 BEFORE DELETE - screenshot: %s", beforePath)
			t.Logf("💾 BEFORE DELETE - HTML: %s", beforeHTMLPath)
		}

		// Install WebSocket message interceptor to log all messages
		// Use a more robust approach that waits for the client to be ready
		var interceptorResult string
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`
				(() => {
					if (window.wsMessagesLogged) {
						return 'Already installed, count=' + window.wsMessageCount;
					}

					window.wsMessageCount = 0;
					window.wsMessages = [];

					// Try multiple ways to find the WebSocket
					let ws = null;

					// Method 1: Via window.liveTemplateClient (exposed by client)
					const client = window.liveTemplateClient;
					if (client && client.webSocketManager) {
						ws = client.webSocketManager.getSocket();
						if (ws) {
							console.log('[WS-LOG] Found WebSocket via window.liveTemplateClient');
						}
					}

					// Method 2: Intercept WebSocket constructor (for future connections)
					if (!ws) {
						console.log('[WS-LOG] Client not ready, will intercept next WebSocket creation');
						const OriginalWebSocket = window.WebSocket;
						window.WebSocket = function(url, protocols) {
							const ws = new OriginalWebSocket(url, protocols);
							console.log('[WS-LOG] New WebSocket created:', url);
							window.currentWebSocket = ws;

							// Install interceptor on this new WebSocket
							const originalOnMessage = ws.onmessage;
							ws.onmessage = function(event) {
								window.wsMessageCount++;
								const msg = JSON.parse(event.data);
								window.wsMessages.push({
									count: window.wsMessageCount,
									time: new Date().toISOString(),
									action: msg.meta?.action,
									success: msg.meta?.success,
									size: event.data.length
								});
								console.log('[WS-RECV #' + window.wsMessageCount + ']',
									'action=' + msg.meta?.action,
									'success=' + msg.meta?.success,
									'size=' + event.data.length + 'B');

								if (originalOnMessage) {
									originalOnMessage.call(this, event);
								}
							};

							return ws;
						};
						window.wsMessagesLogged = true;
						return 'Installed WebSocket constructor interceptor';
					}

					// Install interceptor on existing WebSocket
					console.log('[WS-LOG] Installing interceptor on existing WebSocket, readyState:', ws.readyState);

					const originalOnMessage = ws.onmessage;
					console.log('[WS-LOG] Original onmessage handler exists:', !!originalOnMessage);
					ws.onmessage = function(event) {
						window.wsMessageCount++;
						const msg = JSON.parse(event.data);
						window.wsMessages.push({
							count: window.wsMessageCount,
							time: new Date().toISOString(),
							action: msg.meta?.action,
							success: msg.meta?.success,
							size: event.data.length
						});
						console.log('[WS-RECV #' + window.wsMessageCount + ']',
							'action=' + msg.meta?.action,
							'success=' + msg.meta?.success,
							'size=' + event.data.length + 'B');

						if (originalOnMessage) {
							console.log('[WS-LOG] Calling original handler');
							originalOnMessage.call(this, event);
							console.log('[WS-LOG] Original handler returned');
						} else {
							console.log('[WS-LOG] WARNING: No original handler to call!');
						}
					};

					window.wsMessagesLogged = true;
					return 'Interceptor installed on existing WebSocket (readyState=' + ws.readyState + ')';
				})()
			`, &interceptorResult),
		)
		if err != nil {
			t.Logf("⚠️ Failed to install WebSocket interceptor: %v", err)
		} else {
			t.Logf("📡 WebSocket interceptor: %s", interceptorResult)
		}

		// Override window.confirm to return true (accept)
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`window.confirm = () => true;`, nil),
			chromedp.Evaluate(`
				(() => {
					const deleteButton = document.querySelector('button[lvt-click="delete"]');
					if (deleteButton) {
						deleteButton.click();
						return true;
					}
					return false;
				})()
			`, &postExists),
		)
		if err != nil {
			t.Fatalf("Failed to click delete button: %v", err)
		}

		// Wait a moment for the delete to be sent
		time.Sleep(500 * time.Millisecond)

		// Check WebSocket messages received so far
		var wsMessageInfo string
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(`
				(() => {
					if (!window.wsMessages) return 'No WebSocket interceptor installed';
					return 'Received ' + window.wsMessageCount + ' messages: ' +
						JSON.stringify(window.wsMessages.slice(-3)); // Last 3 messages
				})()
			`, &wsMessageInfo),
		)
		if err == nil {
			t.Logf("📡 WebSocket messages after delete click: %s", wsMessageInfo)
		}

		// CAPTURE STATE AFTER CLICKING DELETE (before waiting for UI update)
		var screenshotAfter []byte
		var htmlAfter string
		err = chromedp.Run(testCtx,
			chromedp.FullScreenshot(&screenshotAfter, 100),
			chromedp.OuterHTML("html", &htmlAfter),
		)
		if err == nil {
			afterPath := fmt.Sprintf("/tmp/delete_after_click_%d.png", time.Now().Unix())
			afterHTMLPath := fmt.Sprintf("/tmp/delete_after_click_%d.html", time.Now().Unix())
			os.WriteFile(afterPath, screenshotAfter, 0644)
			os.WriteFile(afterHTMLPath, []byte(htmlAfter), 0644)
			t.Logf("💾 AFTER DELETE CLICK - screenshot: %s", afterPath)
			t.Logf("💾 AFTER DELETE CLICK - HTML: %s", afterHTMLPath)
		}

		// Now wait for deletion to complete in the UI
		// Wait for the specific row with this data-key to be removed
		err = chromedp.Run(testCtx,
			// Wait for deletion to process - check that the specific data-key is gone
			waitFor(fmt.Sprintf(`
			(() => {
				const table = document.querySelector('table tbody');
				if (!table) return true;
				const targetRow = table.querySelector('tr[data-key=%q]');
				return targetRow === null; // Row with this data-key should be gone
			})()
			`, targetDataKey), 30*time.Second),

			waitForWebSocketReady(30*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			// Wait for page to load
			waitFor(`document.readyState === 'complete'`, 10*time.Second),
		)
		if err != nil {
			t.Logf("⚠️ Delete wait timed out: %v", err)

			// Check final WebSocket message count
			var finalWSInfo string
			chromedp.Run(testCtx,
				chromedp.Evaluate(`
					(() => {
						if (!window.wsMessages) return 'No WebSocket interceptor';
						return 'Total messages: ' + window.wsMessageCount +
							', Last 5: ' + JSON.stringify(window.wsMessages.slice(-5));
					})()
				`, &finalWSInfo),
			)
			t.Logf("📡 Final WebSocket state: %s", finalWSInfo)

			t.Logf("💾 Check BEFORE DELETE files for modal state")
			t.Logf("💾 Check AFTER DELETE CLICK files for post-delete state")
			t.Fatalf("Failed to delete post: %v", err)
		}

		// Verify the specific row with this data-key is no longer in the table
		var postStillExists bool
		err = chromedp.Run(testCtx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table tbody');
					if (!table) return false;
					const targetRow = table.querySelector('tr[data-key=%q]');
					return targetRow !== null; // Returns true if row still exists
				})()
			`, targetDataKey), &postStillExists),
		)
		if err != nil {
			t.Fatalf("Failed to check if post was deleted: %v", err)
		}

		if postStillExists {
			t.Fatalf("❌ Post with data-key %q still exists after deletion", targetDataKey)
		}

		t.Logf("✅ Post %q deleted successfully after confirming", targetTitle)
	})

	// Test Validation Errors
	t.Run("Validation Errors", func(t *testing.T) {
		// Create per-subtest context with individual timeout
		testCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		testCtx, timeoutCancel := context.WithTimeout(testCtx, getBrowserTimeout())
		defer timeoutCancel()

		var (
			errorsVisible    bool
			titleErrorText   string
			contentErrorText string
			formHTML         string
		)

		err := chromedp.Run(testCtx,
			// Navigate to /posts
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(30*time.Second), // Wait for WebSocket init and first update
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"), // Validate no raw template expressions

			// Click the "+ Add Posts" button in toolbar to open modal
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			// Wait for modal to appear
			chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),

			// Submit form WITHOUT filling required fields
			chromedp.WaitVisible(`form[lvt-submit]`, chromedp.ByQuery),
			chromedp.Evaluate(`
				const form = document.querySelector('form[lvt-submit]');
				if (form) {
					// Bypass HTML5 validation to test server-side validation
					form.noValidate = true;
					form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
				}
			`, nil),

			// Wait for validation response - form should still be visible
			waitFor(`document.querySelector('form[lvt-submit]') !== null`, 10*time.Second),

			// Wait a bit for error messages to be injected by client library
			chromedp.Sleep(2*time.Second),

			// Debug: Capture the form HTML
			chromedp.Evaluate(`document.querySelector('form[lvt-submit]')?.outerHTML || 'Form not found'`, &formHTML),

			// Check if error messages are visible in the UI (rendered server-side)
			chromedp.Evaluate(`
				(() => {
					// Look for error messages in <small> tags (server-side rendered via .lvt.HasError)
					const form = document.querySelector('form[lvt-submit]');
					if (!form) return false;
					const smallTags = Array.from(form.querySelectorAll('small'));
					return smallTags.some(el => el.textContent.includes('required') || el.textContent.includes('is required'));
				})()
			`, &errorsVisible),

			// Get specific error texts (server-side rendered)
			chromedp.Evaluate(`
				(() => {
					const form = document.querySelector('form[lvt-submit]');
					if (!form) return '';
					// Find the small tag near the title input
					const titleDiv = Array.from(form.querySelectorAll('div')).find(div => {
						const label = div.querySelector('label');
						return label && label.textContent.includes('Title');
					});
					return titleDiv ? (titleDiv.querySelector('small')?.textContent || '') : '';
				})()
			`, &titleErrorText),
			chromedp.Evaluate(`
				(() => {
					const form = document.querySelector('form[lvt-submit]');
					if (!form) return '';
					// Find the small tag near the content input
					const contentDiv = Array.from(form.querySelectorAll('div')).find(div => {
						const label = div.querySelector('label');
						return label && label.textContent.includes('Content');
					});
					return contentDiv ? (contentDiv.querySelector('small')?.textContent || '') : '';
				})()
			`, &contentErrorText),
		)
		if err != nil {
			t.Fatalf("Failed to test validation: %v", err)
		}

		// Debug: Check what the client has
		var lastWSMessage, clientErrors, activeFormStatus, handleResponseCalled, renderCalled, responseMeta, allWSMessages, errorElementsCount string
		chromedp.Run(testCtx,
			chromedp.Evaluate(`window.__lastWSMessage || 'No WS message'`, &lastWSMessage),
			chromedp.Evaluate(`JSON.stringify(window.liveTemplateClient?.errors || {})`, &clientErrors),
			chromedp.Evaluate(`window.liveTemplateClient?.formLifecycleManager?.activeForm ? 'active' : 'not-active'`, &activeFormStatus),
			chromedp.Evaluate(`window.__lvtHandleResponseCalled ? 'yes' : 'no'`, &handleResponseCalled),
			chromedp.Evaluate(`window.__lvtRenderFieldErrorsCalled ? 'yes' : 'no'`, &renderCalled),
			chromedp.Evaluate(`JSON.stringify(window.__lvtResponseMeta || {})`, &responseMeta),
			chromedp.Evaluate(`JSON.stringify(window.__wsMessages?.slice(-5) || [])`, &allWSMessages),
			chromedp.Evaluate(`document.querySelectorAll('small[data-lvt-error]').length.toString()`, &errorElementsCount),
		)

		t.Logf("Last WS message: %s", lastWSMessage)
		t.Logf("All WS messages (last 5): %s", allWSMessages)
		t.Logf("Client errors state: %s", clientErrors)
		t.Logf("Active form status: %s", activeFormStatus)
		t.Logf("HandleResponse called: %s", handleResponseCalled)
		t.Logf("RenderFieldErrors called: %s", renderCalled)
		t.Logf("Response meta: %s", responseMeta)
		t.Logf("Error elements count: %s", errorElementsCount)
		t.Logf("Form HTML (first 500 chars): %s", formHTML[:min(500, len(formHTML))])

		// Verify errors are displayed in the UI (server-side rendered)
		if !errorsVisible {
			t.Error("❌ Error messages are not visible in the UI")
		}
		t.Log("✅ Error messages are visible in the UI")

		// Verify specific field errors
		if titleErrorText == "" {
			t.Error("❌ Title field error not displayed")
		} else {
			t.Logf("✅ Title error: %s", titleErrorText)
		}

		if contentErrorText == "" {
			t.Error("❌ Content field error not displayed")
		} else {
			t.Logf("✅ Content error: %s", contentErrorText)
		}
	})

	// Test Infinite Scroll Sentinel
	t.Run("Infinite Scroll Sentinel", func(t *testing.T) {
		// The sentinel only appears when HasMore is true (more items to load)
		// Since we added 1 post and default page size is 20, HasMore will be false
		// So we check that the template is configured for infinite scroll by:
		// 1. Checking the generated handler has PaginationMode: "infinite"
		// 2. Verifying template contains infiniteScroll define

		// Read handler file to verify pagination mode
		handlerFile := filepath.Join(blogDir, "app", "posts", "posts.go")
		handlerContent, err := os.ReadFile(handlerFile)
		if err != nil {
			t.Fatalf("Failed to read posts handler: %v", err)
		}

		if !strings.Contains(string(handlerContent), `PaginationMode: "infinite"`) {
			t.Error("❌ Posts handler does not have PaginationMode: \"infinite\"")
		} else {
			t.Log("✅ Posts handler configured with infinite pagination mode")
		}

		// Read template file to verify infiniteScroll block exists
		tmplFile := filepath.Join(blogDir, "app", "posts", "posts.tmpl")
		tmplContent, err := os.ReadFile(tmplFile)
		if err != nil {
			t.Fatalf("Failed to read posts template: %v", err)
		}

		tmplStr := string(tmplContent)
		if !strings.Contains(tmplStr, `id="scroll-sentinel"`) {
			t.Error("❌ Template does not contain scroll-sentinel element")
		} else {
			t.Log("✅ Template contains scroll-sentinel element for infinite scroll")
		}

		// Verify the sentinel appears in actual rendered HTML when there are no template errors
		// (The sentinel won't be visible with only 1 post, but we've verified the configuration)
		t.Log("✅ Infinite scroll pagination configured correctly")
	})

	// Note: Server logs check removed since we're using Docker containers
	// Docker container logs can be checked with: docker logs <container_id>
	t.Run("Server Health Check", func(t *testing.T) {
		// Verify server is still responding
		resp, err := http.Get(serverURL + "/posts")
		if err != nil {
			t.Fatalf("Server not responding: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("Server returned status %d, expected 200", resp.StatusCode)
		} else {
			t.Log("✅ Server is healthy and responding")
		}
	})
}

func ensureTutorialPostExists(ctx context.Context, baseURL string) error {
	const (
		title   = "My First Blog Post"
		content = "This is the content of my first blog post"
	)

	existenceCheck := fmt.Sprintf(`
		(() => {
			const table = document.querySelector('table');
			if (!table) return false;
			const rows = Array.from(table.querySelectorAll('tbody tr'));
			return rows.some(row => {
				const cells = row.querySelectorAll('td');
				return cells.length > 0 && cells[0].textContent.trim() === %q;
			});
		})()
		`, title)

	return chromedp.Run(ctx,
		chromedp.Navigate(baseURL+"/posts"),
		waitForWebSocketReady(30*time.Second),
		chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
		chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
		chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
		waitFor(`document.querySelector('[role="dialog"]') && !document.querySelector('[role="dialog"]').hasAttribute('hidden')`, 10*time.Second),
		chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector('input[name="title"]').value = ''`, nil),
		chromedp.Evaluate(`document.querySelector('textarea[name="content"]').value = ''`, nil),
		chromedp.SendKeys(`input[name="title"]`, title, chromedp.ByQuery),
		chromedp.SendKeys(`textarea[name="content"]`, content, chromedp.ByQuery),
		chromedp.Click(`input[name="published"]`, chromedp.ByQuery),
		chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
		waitFor(existenceCheck, 30*time.Second),
		chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
		waitFor(`document.readyState === 'complete'`, 10*time.Second),
	)
}
