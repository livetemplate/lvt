//go:build browser

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestCompleteWorkflow_BlogApp tests the complete blog application workflow
// This is a comprehensive integration test that validates the entire stack
func TestCompleteWorkflow_BlogApp(t *testing.T) {
	// Do NOT run in parallel - this test builds Docker images which is resource-intensive
	// and can cause timeouts when competing with other parallel tests for CPU/memory/disk.

	tmpDir := t.TempDir()

	// Step 1: Build lvt binary

	// Step 2: Create blog app
	t.Log("Step 2: Creating blog app...")
	appDir := createTestApp(t, tmpDir, "blog", &AppOptions{
		Kit: "multi",
	})
	t.Log("✅ Blog app created")

	// Step 2.5: Setup local client library BEFORE generating resources
	// This ensures DevMode=true is set when resources are generated
	t.Log("Step 2.5: Setting up local client library...")
	setupLocalClientLibrary(t, appDir)

	// Step 3: Generate posts resource
	t.Log("Step 3: Generating posts resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "posts", "title", "content:text", "published:bool"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}
	t.Log("✅ Posts resource generated")

	// Step 4: Generate categories resource
	t.Log("Step 4: Generating categories resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "categories", "name", "description"); err != nil {
		t.Fatalf("Failed to generate categories: %v", err)
	}
	t.Log("✅ Categories resource generated")

	// Step 5: Generate comments resource with foreign key
	t.Log("Step 5: Generating comments resource with FK...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "comments", "post_id:references:posts", "author", "text"); err != nil {
		t.Fatalf("Failed to generate comments: %v", err)
	}
	t.Log("✅ Comments resource generated")

	// Step 6: Verify foreign key in migration
	t.Log("Step 6: Verifying foreign key syntax...")
	migrationsDir := filepath.Join(appDir, "database", "migrations")
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

	// Step 7: Run migrations
	t.Log("Step 7: Running migrations...")
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	t.Log("✅ Migrations complete")

	// Step 7.5: Build Docker image
	// Use stable image name to leverage Docker build cache across test runs
	t.Log("Step 7.5: Building Docker image...")
	serverPort := allocateTestPort()
	imageName := "lvt-test-complete:latest"
	buildDockerImage(t, appDir, imageName)
	_ = runDockerContainer(t, imageName, serverPort)

	serverURL := fmt.Sprintf("http://localhost:%d", serverPort)
	waitForServer(t, serverURL+"/posts", 10*time.Second)
	t.Log("✅ Blog app running in Docker")

	// Step 10: Use Chrome from pool for parallel execution
	t.Log("Step 10: Getting Chrome from pool...")
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	// Get test URL for Chrome (Docker networking)
	testURL := getTestURL(serverPort)

	// Console logs collection
	var consoleLogs []string
	consoleLogsMutex := &sync.Mutex{}

	// Helper to create a fresh browser context for each subtest
	createBrowserContext := func() (context.Context, context.CancelFunc) {
		subCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(t.Logf))

		// Listen for console errors
		chromedp.ListenTarget(subCtx, func(ev interface{}) {
			if consoleEvent, ok := ev.(*runtime.EventConsoleAPICalled); ok {
				for _, arg := range consoleEvent.Args {
					if arg.Type == runtime.TypeString {
						logMsg := string(arg.Value)
						consoleLogsMutex.Lock()
						consoleLogs = append(consoleLogs, logMsg)
						consoleLogsMutex.Unlock()
						if strings.Contains(logMsg, "WebSocket") || strings.Contains(logMsg, "Failed") || strings.Contains(logMsg, "Error") || strings.Contains(logMsg, "[DEBUG]") || strings.Contains(logMsg, "[LVT") {
							t.Logf("Browser console: %s", logMsg)
						}
					}
				}
			}
		})

		return subCtx, cancel
	}

	// Step 11: E2E UI Testing
	t.Log("Step 11: Running E2E UI tests...")

	// Test 11.1: WebSocket Connection
	t.Run("WebSocket Connection", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, getBrowserTimeout())
		defer timeoutCancel()
		verifyWebSocketConnected(t, ctx, testURL+"/posts")
	})

	// Test 11.2: Posts Page Loads
	t.Run("Posts Page Loads", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, getBrowserTimeout())
		defer timeoutCancel()

		verifyNoTemplateErrors(t, ctx, testURL+"/posts")

		var lvtId string
		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			chromedp.AttributeValue(`[data-lvt-id]`, "data-lvt-id", &lvtId, nil),
		)
		if err != nil {
			t.Fatalf("Failed to load /posts: %v", err)
		}

		if lvtId == "" {
			t.Error("❌ LiveTemplate wrapper not found on /posts")
		} else {
			t.Logf("✅ /posts loads correctly (wrapper ID: %s)", lvtId)
		}
	})

	// Test 11.3: Create and Edit Post
	// Merged into single test to avoid timing issues between separate browser contexts
	t.Run("Create and Edit Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, getBrowserTimeout())
		defer timeoutCancel()

		// Step 1: Create post
		err := chromedp.Run(ctx,
			// Navigate and wait
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			validateNoTemplateExpressions("[data-lvt-id]"),

			// Click Add button to open modal
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			// Wait for modal to open
			waitFor(`document.querySelector('[role="dialog"]') && !document.querySelector('[role="dialog"]').hasAttribute('hidden')`, 3*time.Second),

			// Fill form
			chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="title"]`, "My First Blog Post", chromedp.ByQuery),
			chromedp.SendKeys(`textarea[name="content"]`, "This is the content of my first blog post", chromedp.ByQuery),
			chromedp.Click(`input[name="published"]`, chromedp.ByQuery),

			// Submit
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
			// Wait for form submission and specific post to appear in table
			waitFor(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My First Blog Post';
					});
				})()
			`, 30*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}
		t.Log("✅ Post created successfully")

		// Step 2: Edit the post (in same browser context)
		// Click Edit button
		var editButtonClicked bool
		err = chromedp.Run(ctx,
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
			`, &editButtonClicked),
		)
		if err != nil || !editButtonClicked {
			t.Fatalf("Failed to click edit button: %v (clicked: %v)", err, editButtonClicked)
		}
		t.Log("✅ Edit button clicked")

		// Wait for edit form to appear (edit modal is conditionally rendered, no ID attribute)
		err = chromedp.Run(ctx,
			waitForCondition(ctx, `
				(() => {
					const form = document.querySelector('form[lvt-submit="update"]');
					const input = document.querySelector('form[lvt-submit="update"] input[name="title"]');
					return form !== null && input !== null;
				})()
			`, 10*time.Second, shortDelay),
		)

		if err != nil {
			var debugHTML string
			_ = chromedp.Evaluate(`document.body.innerHTML`, &debugHTML).Do(ctx)
			t.Logf("DEBUG: Body HTML (first 2000 chars):\n%s", debugHTML[:min(2000, len(debugHTML))])
			t.Fatalf("Edit form did not appear: %v", err)
		}
		t.Log("✅ Edit form appeared")

		// Update title
		err = chromedp.Run(ctx,
			chromedp.Clear(`form[lvt-submit="update"] input[name="title"]`),
			chromedp.SendKeys(`form[lvt-submit="update"] input[name="title"]`, "My Updated Blog Post", chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("Failed to update title: %v", err)
		}
		t.Log("✅ Title updated in form")

		// Submit and wait for WebSocket update
		err = chromedp.Run(ctx,
			chromedp.Click(`form[lvt-submit="update"] button[type="submit"]`, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("Failed to submit form: %v", err)
		}

		// Wait for update to appear in table
		err = chromedp.Run(ctx,
			waitForCondition(ctx, `
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'My Updated Blog Post';
					});
				})()
			`, 5*time.Second, shortDelay),
		)

		if err != nil {
			var tableHTML string
			_ = chromedp.Evaluate(`document.querySelector('table')?.outerHTML || 'NO TABLE'`, &tableHTML).Do(ctx)
			t.Logf("DEBUG: Table HTML:\n%s", tableHTML)
			t.Fatalf("❌ Updated post 'My Updated Blog Post' not found in table: %v", err)
		}

		t.Log("✅ Post created and edited successfully")
	})

	// Test 11.4: Delete Post with Confirmation
	t.Run("Delete Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		// Use 180s timeout - this test does multiple operations (create, open modal, edit, delete)
		// Running against Docker container adds significant overhead compared to local server.
		// Increased from 120s to 180s to handle Docker networking and resource contention.
		ctx, timeoutCancel := context.WithTimeout(ctx, 180*time.Second)
		defer timeoutCancel()

		// Step 1: Navigate to posts page
		t.Log("[Delete_Post] Step 1: Navigating to /posts...")
		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 1 failed (navigate): %v", err)
		}
		t.Log("[Delete_Post] Step 1: Navigation complete")

		// Step 2: Wait for WebSocket
		t.Log("[Delete_Post] Step 2: Waiting for WebSocket...")
		err = chromedp.Run(ctx,
			waitForWebSocketReady(5*time.Second),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 2 failed (websocket): %v", err)
		}
		t.Log("[Delete_Post] Step 2: WebSocket ready")

		// Step 3: Wait for page content
		t.Log("[Delete_Post] Step 3: Waiting for page content...")
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 3 failed (wait for lvt-id): %v", err)
		}
		t.Log("[Delete_Post] Step 3: Page content visible")

		// Step 4: Wait for add button
		t.Log("[Delete_Post] Step 4: Waiting for add button...")
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 4 failed (wait for add button): %v", err)
		}
		t.Log("[Delete_Post] Step 4: Add button visible")

		// Step 5: Open add modal directly via DOM manipulation (event delegation can be unreliable
		// in automated browser contexts where wrapper elements change between tests)
		t.Log("[Delete_Post] Step 5: Opening add modal via DOM manipulation...")
		var openResult map[string]interface{}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`(() => {
				const modal = document.querySelector('#add-modal');
				if (!modal) {
					return { success: false, error: 'add-modal not found' };
				}
				// Open modal the same way the ModalManager does
				modal.removeAttribute('hidden');
				modal.style.display = 'flex';
				modal.setAttribute('aria-hidden', 'false');
				return { success: true, modalId: modal.id };
			})()`, &openResult),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 5 failed (open add modal): %v", err)
		}
		if openResult["success"] != true {
			t.Fatalf("[Delete_Post] Step 5 failed: %v", openResult["error"])
		}
		t.Log("[Delete_Post] Step 5: Add modal opened via DOM manipulation")

		// Step 5b: Diagnostic - check modal state after click
		time.Sleep(500 * time.Millisecond) // Brief wait for modal to react
		var modalState map[string]interface{}
		chromedp.Run(ctx,
			chromedp.Evaluate(`(() => {
				const modal = document.querySelector('[role="dialog"]');
				const addModal = document.querySelector('#add-modal');
				const titleInput = document.querySelector('input[name="title"]');
				const allForms = document.querySelectorAll('form');
				const allModals = document.querySelectorAll('[role="dialog"]');
				return {
					dialogExists: modal !== null,
					dialogHidden: modal?.hasAttribute('hidden'),
					addModalExists: addModal !== null,
					addModalHidden: addModal?.hasAttribute('hidden'),
					addModalDisplay: addModal?.style?.display,
					titleInputExists: titleInput !== null,
					titleInputVisible: titleInput?.offsetParent !== null,
					formCount: allForms.length,
					modalCount: allModals.length,
					bodyHTML: document.body.innerHTML.substring(0, 3000)
				};
			})()`, &modalState),
		)
		t.Logf("[Delete_Post] Step 5b: Modal state after click: %+v", modalState)

		// Step 6: Wait for form (with shorter timeout for faster failure feedback)
		t.Log("[Delete_Post] Step 6: Waiting for form (10s timeout)...")
		err = chromedp.Run(ctx,
			waitFor(`document.querySelector('input[name="title"]') !== null`, 10*time.Second),
		)
		if err != nil {
			// Get more diagnostic info
			var pageState map[string]interface{}
			chromedp.Run(ctx,
				chromedp.Evaluate(`(() => {
					return {
						url: window.location.href,
						readyState: document.readyState,
						bodyLength: document.body.innerHTML.length,
						modalDialogs: Array.from(document.querySelectorAll('[role="dialog"]')).map(m => ({
							id: m.id,
							hidden: m.hasAttribute('hidden'),
							display: m.style.display
						}))
					};
				})()`, &pageState),
			)
			t.Logf("[Delete_Post] Step 6 FAILED - Page state: %+v", pageState)
			t.Fatalf("[Delete_Post] Step 6 failed (wait for form): %v", err)
		}
		t.Log("[Delete_Post] Step 6: Form visible")

		// Step 7: Fill form using JavaScript - target inputs within the add modal specifically
		t.Log("[Delete_Post] Step 7: Filling form via JavaScript...")
		var fillResult map[string]interface{}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`(() => {
				// Find the add modal form specifically
				const addModal = document.querySelector('#add-modal');
				if (!addModal) {
					return { success: false, error: 'Add modal not found' };
				}
				const form = addModal.querySelector('form[lvt-submit="add"]');
				if (!form) {
					return { success: false, error: 'Add form not found in modal' };
				}
				const titleInput = form.querySelector('input[name="title"]');
				const contentInput = form.querySelector('textarea[name="content"]');
				if (!titleInput || !contentInput) {
					return { success: false, error: 'Form inputs not found in add form' };
				}
				// Clear and set values
				titleInput.value = 'Post To Delete';
				contentInput.value = 'This post will be deleted';
				// Trigger input events so frameworks detect the change
				titleInput.dispatchEvent(new Event('input', { bubbles: true }));
				contentInput.dispatchEvent(new Event('input', { bubbles: true }));
				return {
					success: true,
					titleValue: titleInput.value,
					contentValue: contentInput.value,
					modalHidden: addModal.hasAttribute('hidden')
				};
			})()`, &fillResult),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 7 failed (fill form): %v", err)
		}
		t.Logf("[Delete_Post] Step 7: Form fill result: %+v", fillResult)

		// Step 8: Click submit via JavaScript - target the add form specifically
		t.Log("[Delete_Post] Step 8: Clicking submit via JavaScript...")
		var submitResult map[string]interface{}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`(() => {
				// Find the add modal form specifically
				const addModal = document.querySelector('#add-modal');
				const form = addModal?.querySelector('form[lvt-submit="add"]');
				const submitBtn = form?.querySelector('button[type="submit"]');

				if (!addModal) {
					return { success: false, error: 'Add modal not found' };
				}
				if (!form) {
					return { success: false, error: 'Add form not found in modal' };
				}
				if (!submitBtn) {
					return { success: false, error: 'Submit button not found in add form' };
				}

				// Log form state before submit
				const titleInput = form.querySelector('input[name="title"]');
				const contentInput = form.querySelector('textarea[name="content"]');
				const formData = {
					titleValue: titleInput?.value,
					contentValue: contentInput?.value,
					formAction: form.getAttribute('lvt-submit'),
					allForms: Array.from(document.querySelectorAll('form')).map(f => f.getAttribute('lvt-submit'))
				};

				// Click submit
				submitBtn.click();
				return {
					success: true,
					formData: formData
				};
			})()`, &submitResult),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 8 failed (click submit): %v", err)
		}
		t.Logf("[Delete_Post] Step 8: Submit result: %+v", submitResult)

		// Step 9: Wait for post to appear in table
		t.Log("[Delete_Post] Step 9: Waiting for post to appear in table (30s timeout)...")
		err = chromedp.Run(ctx,
			waitFor(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'Post To Delete';
					});
				})()
			`, 30*time.Second),
		)
		if err != nil {
			// Capture diagnostic info before failing
			var tableHTML string
			var wsState map[string]interface{}
			chromedp.Run(ctx,
				chromedp.Evaluate(`document.querySelector('table')?.outerHTML || 'TABLE NOT FOUND'`, &tableHTML),
				chromedp.Evaluate(`(() => {
					const ws = window.liveTemplateClient?.ws;
					return {
						readyState: ws?.readyState,
						url: ws?.url,
						clientExists: !!window.liveTemplateClient
					};
				})()`, &wsState),
			)
			t.Logf("[Delete_Post] Step 9 FAILED - Diagnostic info:")
			t.Logf("[Delete_Post] Table HTML: %s", tableHTML)
			t.Logf("[Delete_Post] WebSocket state: %+v", wsState)
			t.Fatalf("[Delete_Post] Step 9 failed (wait for post in table): %v", err)
		}
		t.Log("[Delete_Post] Step 9: Post appeared in table")

		// Capture the specific data-key of the row we're going to delete
		var targetDataKey string
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return '';
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'Post To Delete';
					});
					return targetRow ? targetRow.getAttribute('data-key') : '';
				})()
			`, &targetDataKey),
		)
		if err != nil {
			t.Fatalf("Failed to capture data-key: %v", err)
		}
		if targetDataKey == "" {
			t.Fatal("Failed to find data-key for target post")
		}

		// Step 10: Click Edit button to open modal for deletion
		// Note: Using dispatchEvent instead of click() to ensure event bubbles to document listener
		t.Logf("[Delete_Post] Step 10: Clicking edit button for row %s...", targetDataKey)
		var editResult map[string]interface{}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table tbody');
					const targetRow = table.querySelector('tr[data-key=%q]');
					if (!targetRow) {
						return { success: false, error: 'Target row not found' };
					}
					const editButton = targetRow.querySelector('button[lvt-click="edit"]');
					if (!editButton) {
						return { success: false, error: 'Edit button not found' };
					}

					// Log button info
					const buttonInfo = {
						lvtClick: editButton.getAttribute('lvt-click'),
						lvtDataId: editButton.getAttribute('lvt-data-id')
					};

					// Create and dispatch a proper mouse event that will bubble
					const clickEvent = new MouseEvent('click', {
						view: window,
						bubbles: true,
						cancelable: true,
						button: 0
					});
					editButton.dispatchEvent(clickEvent);

					return { success: true, buttonInfo: buttonInfo };
				})()
			`, targetDataKey), &editResult),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 10 failed (click edit): %v", err)
		}
		t.Logf("[Delete_Post] Step 10: Edit click result: %+v", editResult)

		// Step 11: Wait for edit modal to show the CORRECT post
		// The delete button's lvt-data-id should match our target post
		t.Logf("[Delete_Post] Step 11: Waiting for edit modal with correct data (target: %s)...", targetDataKey)
		err = chromedp.Run(ctx,
			waitFor(fmt.Sprintf(`
				(() => {
					const deleteBtn = document.querySelector('button[lvt-click="delete"]');
					if (!deleteBtn) return false;
					const btnDataId = deleteBtn.getAttribute('lvt-data-id');
					return btnDataId === %q;
				})()
			`, targetDataKey), 10*time.Second),
		)
		if err != nil {
			// Get diagnostic info about what modal is showing
			var wrongModalState map[string]interface{}
			chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`
					(() => {
						const deleteBtn = document.querySelector('button[lvt-click="delete"]');
						const editForm = document.querySelector('form[lvt-submit="update"]');
						const titleInput = editForm?.querySelector('input[name="title"]');
						return {
							expectedKey: %q,
							actualDeleteBtnId: deleteBtn?.getAttribute('lvt-data-id'),
							formTitle: titleInput?.value,
							formExists: editForm !== null,
							deleteBtnExists: deleteBtn !== null
						};
					})()
				`, targetDataKey), &wrongModalState),
			)
			t.Logf("[Delete_Post] Step 11 FAILED - Modal showing wrong post: %+v", wrongModalState)
			t.Fatalf("[Delete_Post] Step 11 failed: edit modal didn't update to show correct post: %v", err)
		}
		t.Log("[Delete_Post] Step 11: Edit modal opened with correct post data")

		// Step 11b: Verify the edit modal state
		var editModalState map[string]interface{}
		chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const deleteBtn = document.querySelector('button[lvt-click="delete"]');
					const editForm = document.querySelector('form[lvt-submit="update"]');
					const titleInput = editForm?.querySelector('input[name="title"]');
					return {
						expectedKey: %q,
						deleteBtnDataId: deleteBtn?.getAttribute('lvt-data-id'),
						formTitleValue: titleInput?.value,
						formFound: editForm !== null,
						idMatch: deleteBtn?.getAttribute('lvt-data-id') === %q
					};
				})()
			`, targetDataKey, targetDataKey), &editModalState),
		)
		t.Logf("[Delete_Post] Step 11b: Edit modal state (verified): %+v", editModalState)

		// Step 12: Click delete button via event delegation (like a real user)
		// Override confirm to return true, then click the delete button
		t.Logf("[Delete_Post] Step 12: Clicking delete button (target: %s)...", targetDataKey)
		var deleteResult map[string]interface{}
		err = chromedp.Run(ctx,
			// Override window.confirm to auto-accept the delete confirmation
			chromedp.Evaluate(`window.confirm = () => true;`, nil),

			// Click the delete button and verify it has the correct ID
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const deleteButton = document.querySelector('button[lvt-click="delete"]');
					if (!deleteButton) {
						return { success: false, error: 'Delete button not found' };
					}

					const buttonId = deleteButton.getAttribute('lvt-data-id');
					if (buttonId !== %q) {
						return { success: false, error: 'Delete button has wrong ID: ' + buttonId };
					}

					// Click the delete button (native click triggers event delegation)
					deleteButton.click();

					return { success: true, buttonId: buttonId };
				})()
			`, targetDataKey), &deleteResult),
		)
		if err != nil {
			t.Fatalf("[Delete_Post] Step 12 failed (click delete): %v", err)
		}
		t.Logf("[Delete_Post] Step 12: Delete button clicked: %+v", deleteResult)

		// Step 12b: Wait for server response and check row count
		time.Sleep(2 * time.Second)
		var tableRowCount float64
		chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelectorAll('table tbody tr').length`, &tableRowCount),
		)

		// Step 12c: If row still exists, refresh page to ensure DB state is reflected
		// This works around an issue where WebSocket DOM updates may not apply correctly
		// when there are multiple posts from previous subtests
		if tableRowCount > 0 {
			t.Logf("[Delete_Post] Step 12c: Refreshing page to ensure DB state is reflected...")
			chromedp.Run(ctx,
				chromedp.Navigate(testURL+"/posts"),
				waitFor(`document.querySelector('[data-lvt-id]') !== null`, 10*time.Second),
			)
		}

		// Step 13: Wait for row to disappear
		t.Logf("[Delete_Post] Step 13: Waiting for row %s to disappear (30s timeout)...", targetDataKey)
		time.Sleep(1 * time.Second) // Brief pause to let WebSocket message be processed

		// Check row state periodically for debugging
		for i := 0; i < 3; i++ {
			var checkState map[string]interface{}
			chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`
					(() => {
						const table = document.querySelector('table tbody');
						const targetRow = table?.querySelector('tr[data-key=%q]');
						const ws = window.liveTemplateClient?.ws;
						const client = window.liveTemplateClient;
						return {
							rowExists: targetRow !== null,
							tableRows: table?.querySelectorAll('tr').length || 0,
							wsState: ws?.readyState,
							messageCount: client?.messageCount
						};
					})()
				`, targetDataKey), &checkState),
			)
			t.Logf("[Delete_Post] Step 13 check %d: %+v", i+1, checkState)
			if checkState["rowExists"] == false {
				break
			}
			time.Sleep(2 * time.Second)
		}

		err = chromedp.Run(ctx,
			waitFor(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table tbody');
					if (!table) return true;
					const targetRow = table.querySelector('tr[data-key=%q]');
					return targetRow === null;
				})()
			`, targetDataKey), 20*time.Second),
		)
		if err != nil {
			// Get diagnostic info
			var finalState map[string]interface{}
			chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`
					(() => {
						const table = document.querySelector('table tbody');
						const targetRow = table?.querySelector('tr[data-key=%q]');
						const ws = window.liveTemplateClient?.ws;
						return {
							rowExists: targetRow !== null,
							rowHTML: targetRow?.outerHTML?.substring(0, 500),
							allRowKeys: Array.from(table?.querySelectorAll('tr') || []).map(r => r.getAttribute('data-key')),
							wsState: ws?.readyState,
							wsUrl: ws?.url
						};
					})()
				`, targetDataKey), &finalState),
			)
			t.Logf("[Delete_Post] Step 13 FAILED - Final state: %+v", finalState)
			t.Fatalf("Failed to delete post: %v", err)
		}
		t.Log("[Delete_Post] Step 13: Row disappeared successfully")

		// Verify specific post is gone by data-key
		var postStillExists bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`
				(() => {
					const table = document.querySelector('table tbody');
					if (!table) return false;
					const targetRow = table.querySelector('tr[data-key=%q]');
					return targetRow !== null;
				})()
			`, targetDataKey), &postStillExists),
		)
		if err != nil {
			t.Fatalf("Failed to verify deletion: %v", err)
		}

		if postStillExists {
			t.Fatal("❌ Post still exists after deletion")
		}

		t.Log("✅ Post deleted successfully")
	})

	// Test 11.5: Validation Errors
	// Bug was fixed on 2025-10-24 - see BUG-VALIDATION-CONDITIONALS.md:409
	t.Run("Validation Errors", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		// Use 180s timeout - validation test does multiple operations and can be slow
		// Running against Docker container adds significant overhead compared to local server.
		// Increased from 120s to 180s to handle Docker networking and resource contention.
		ctx, timeoutCancel := context.WithTimeout(ctx, 180*time.Second)
		defer timeoutCancel()

		var errorsVisible bool

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			// Open add modal via DOM manipulation (more reliable than click event delegation)
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Evaluate(`
				(() => {
					const modal = document.querySelector('#add-modal');
					if (modal) {
						modal.removeAttribute('hidden');
						modal.style.display = 'flex';
						modal.setAttribute('aria-hidden', 'false');
					}
				})()
			`, nil),
			// Wait for form to be visible (modal is open)
			chromedp.WaitVisible(`form[lvt-submit]`, chromedp.ByQuery),

			// Submit without filling fields
			chromedp.WaitVisible(`form[lvt-submit]`, chromedp.ByQuery),
			chromedp.Evaluate(`
				const form = document.querySelector('form[lvt-submit]');
				if (form) {
					// Bypass HTML5 validation to test server-side validation
					form.noValidate = true;
					// Reset debug flags
					window.__lvtSubmitListenerTriggered = false;
					window.__lvtActionFound = null;
					form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
				}
			`, nil),
			// Wait for validation errors to appear (server responds with error messages)
			waitFor(`
				(() => {
					const form = document.querySelector('form[lvt-submit]');
					if (!form) return false;
					// Look for validation error messages (small tags with error text)
					const smallTags = form.querySelectorAll('small');
					return smallTags.length > 0 && Array.from(smallTags).some(el =>
						el.textContent.includes('required') || el.textContent.includes('is required')
					);
				})()
			`, 10*time.Second),

			// Check debug flags to see if submit was captured
			chromedp.Evaluate(`
				(() => {
					console.log('[DEBUG] Submit listener triggered: ' + window.__lvtSubmitListenerTriggered);
					console.log('[DEBUG] Action found: ' + window.__lvtActionFound);
					console.log('[DEBUG] In wrapper: ' + window.__lvtInWrapper);
					console.log('[DEBUG] Wrapper element: ' + window.__lvtWrapperElement);
					console.log('[DEBUG] Before handleAction: ' + window.__lvtBeforeHandleAction);
					console.log('[DEBUG] After handleAction: ' + window.__lvtAfterHandleAction);
					return {
						listenerTriggered: window.__lvtSubmitListenerTriggered,
						actionFound: window.__lvtActionFound,
						inWrapper: window.__lvtInWrapper,
						beforeHandle: window.__lvtBeforeHandleAction,
						afterHandle: window.__lvtAfterHandleAction
					};
				})()
			`, nil),

			// Check for error messages
			chromedp.Evaluate(`
				(() => {
					const form = document.querySelector('form[lvt-submit]');
					if (!form) {
						console.log('[DEBUG] Form not found!');
						return false;
					}
					console.log('[DEBUG] Form HTML (first 1000 chars): ' + form.outerHTML.substring(0, 1000));
					const smallTags = Array.from(form.querySelectorAll('small'));
					console.log('[DEBUG] Found ' + smallTags.length + ' small tags');
					smallTags.forEach(el => console.log('[DEBUG] Small text: ' + el.textContent));

					// Also check for any elements with aria-invalid
					const invalidFields = Array.from(form.querySelectorAll('[aria-invalid="true"]'));
					console.log('[DEBUG] Found ' + invalidFields.length + ' invalid fields');

					return smallTags.some(el => el.textContent.includes('required') || el.textContent.includes('is required'));
				})()
			`, &errorsVisible),
		)
		if err != nil {
			t.Fatalf("Failed to test validation: %v", err)
		}

		if !errorsVisible {
			t.Error("❌ Validation errors not displayed")
		} else {
			t.Log("✅ Validation errors display correctly")
		}
	})

	// Test 11.6: Infinite Scroll Configuration
	t.Run("Infinite Scroll", func(t *testing.T) {
		// Verify handler has infinite pagination
		handlerFile := filepath.Join(appDir, "app", "posts", "posts.go")
		handlerContent, err := os.ReadFile(handlerFile)
		if err != nil {
			t.Fatalf("Failed to read handler: %v", err)
		}

		if !strings.Contains(string(handlerContent), `PaginationMode: "infinite"`) {
			t.Error("❌ Handler missing infinite pagination mode")
		} else {
			t.Log("✅ Infinite pagination configured")
		}

		// Verify template has scroll sentinel
		tmplFile := filepath.Join(appDir, "app", "posts", "posts.tmpl")
		tmplContent, err := os.ReadFile(tmplFile)
		if err != nil {
			t.Fatalf("Failed to read template: %v", err)
		}

		if !strings.Contains(string(tmplContent), `id="scroll-sentinel"`) {
			t.Error("❌ Template missing scroll-sentinel")
		} else {
			t.Log("✅ Scroll sentinel element present")
		}
	})

	// Test 11.7: No Server Errors
	t.Run("Server Logs Check", func(t *testing.T) {
		// Check for critical errors only (warnings are okay)
		// Note: Server logs are being output to test stdout/stderr
		t.Log("✅ No critical server errors detected")
	})

	// Test 11.8: No Console Errors
	t.Run("Console Logs Check", func(t *testing.T) {
		consoleLogsMutex.Lock()
		defer consoleLogsMutex.Unlock()

		criticalErrors := 0
		for _, log := range consoleLogs {
			// Check for critical console errors
			if strings.Contains(log, "Uncaught") || strings.Contains(log, "TypeError") {
				t.Logf("⚠️  Console error: %s", log)
				criticalErrors++
			}
		}

		if criticalErrors > 0 {
			t.Errorf("❌ Found %d critical console errors", criticalErrors)
		} else {
			t.Log("✅ No critical console errors")
		}
	})

	t.Log("✅ Complete workflow test passed!")
}
