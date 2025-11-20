package e2e

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestDeleteWithMultiplePosts specifically tests delete functionality when database has multiple posts
// This reproduces the issue seen in TestCompleteWorkflow_BlogApp where Delete_Post fails after Create_and_Edit_Post
func TestDeleteWithMultiplePosts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := t.TempDir()

	// Create blog app
	t.Log("Creating blog app...")
	appDir := createTestApp(t, tmpDir, "blog", &AppOptions{
		Kit: "multi",
	})
	setupLocalClientLibrary(t, appDir)

	// Generate posts resource
	t.Log("Generating posts resource...")
	runLvtCommand(t, appDir, "gen", "resource", "posts", "title", "content:text", "published:bool")

	// Run migrations
	t.Log("Running migrations...")
	runLvtCommand(t, appDir, "migration", "up")

	// Build Docker image
	t.Log("Building Docker image...")
	serverPort := allocateTestPort()
	imageName := "lvt-test-delete-multi:latest"
	buildDockerImage(t, appDir, imageName)
	_ = runDockerContainer(t, imageName, serverPort)

	serverURL := fmt.Sprintf("http://localhost:%d", serverPort)
	waitForServer(t, serverURL+"/posts", 10*time.Second)
	t.Log("✅ Server running")

	// Get Chrome from pool
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	testURL := getTestURL(serverPort)

	// Console logs collection
	var consoleLogs []string
	consoleLogsMutex := &sync.Mutex{}

	createBrowserContext := func() (context.Context, context.CancelFunc) {
		subCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(t.Logf))

		chromedp.ListenTarget(subCtx, func(ev interface{}) {
			if consoleEvent, ok := ev.(*runtime.EventConsoleAPICalled); ok {
				for _, arg := range consoleEvent.Args {
					if arg.Type == runtime.TypeString {
						logMsg := string(arg.Value)
						consoleLogsMutex.Lock()
						consoleLogs = append(consoleLogs, logMsg)
						consoleLogsMutex.Unlock()
						if strings.Contains(logMsg, "WebSocket") || strings.Contains(logMsg, "delete") || strings.Contains(logMsg, "DELETE") {
							t.Logf("Browser console: %s", logMsg)
						}
					}
				}
			}
		})

		return subCtx, cancel
	}

	// Step 1: Create first post
	t.Run("Create_First_Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			// Create first post
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="title"]`, "First Post", chromedp.ByQuery),
			chromedp.SendKeys(`textarea[name="content"]`, "Content of first post", chromedp.ByQuery),
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
			waitFor(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'First Post';
					});
				})()
			`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to create first post: %v", err)
		}
		t.Log("✅ First post created")
	})

	// Step 2: Create second post
	t.Run("Create_Second_Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			// Create second post
			chromedp.WaitVisible(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`input[name="title"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="title"]`, "Second Post", chromedp.ByQuery),
			chromedp.SendKeys(`textarea[name="content"]`, "Content of second post", chromedp.ByQuery),
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
			waitFor(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'Second Post';
					});
				})()
			`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to create second post: %v", err)
		}
		t.Log("✅ Second post created")
	})

	// Step 3: Edit the first post (to mimic Create_and_Edit_Post test behavior)
	t.Run("Edit_First_Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			// Find "First Post" and click Edit
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'First Post';
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
			`, nil),

			// Wait for edit form
			waitFor(`document.querySelector('form[lvt-submit="update"] input[name="title"]') !== null`, 5*time.Second),

			// Update title
			chromedp.Clear(`form[lvt-submit="update"] input[name="title"]`),
			chromedp.SendKeys(`form[lvt-submit="update"] input[name="title"]`, "First Post - Edited", chromedp.ByQuery),

			// Submit
			chromedp.Click(`form[lvt-submit="update"] button[type="submit"]`, chromedp.ByQuery),

			// Wait for update in table
			waitFor(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'First Post - Edited';
					});
				})()
			`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to edit first post: %v", err)
		}
		t.Log("✅ First post edited")
	})

	// Step 4: Verify both posts exist
	t.Run("Verify_Both_Posts_Exist", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		var postCount int
		var tableHTML string

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return 0;
					return table.querySelectorAll('tbody tr').length;
				})()
			`, &postCount),

			chromedp.Evaluate(`document.querySelector('table')?.outerHTML || 'NO TABLE'`, &tableHTML),
		)
		if err != nil {
			t.Fatalf("Failed to check posts: %v", err)
		}

		t.Logf("Post count: %d", postCount)
		t.Logf("Table HTML (first 500 chars): %s", tableHTML[:min(500, len(tableHTML))])

		if postCount != 2 {
			t.Fatalf("Expected 2 posts, got %d", postCount)
		}
		t.Log("✅ Both posts exist")
	})

	// Step 5: Delete the second post (the one we want to delete)
	t.Run("Delete_Second_Post", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		var clickedEdit bool
		var foundDeleteButton bool
		var clickedDelete bool

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			// Find "Second Post" and click Edit
			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					console.log('[DELETE TEST] Total rows:', rows.length);

					const targetRow = rows.find(row => {
						const cells = row.querySelectorAll('td');
						const titleText = cells.length > 0 ? cells[0].textContent.trim() : '';
						console.log('[DELETE TEST] Checking row with title:', titleText);
						return titleText === 'Second Post';
					});

					if (targetRow) {
						console.log('[DELETE TEST] Found target row for Second Post');
						const editButton = targetRow.querySelector('button[lvt-click="edit"]');
						if (editButton) {
							console.log('[DELETE TEST] Found edit button, clicking...');
							editButton.click();
							return true;
						} else {
							console.log('[DELETE TEST] Edit button not found in row');
						}
					} else {
						console.log('[DELETE TEST] Target row not found');
					}
					return false;
				})()
			`, &clickedEdit),
		)
		if err != nil {
			t.Fatalf("Failed during edit click: %v", err)
		}
		if !clickedEdit {
			t.Fatal("❌ Failed to click edit button for Second Post")
		}
		t.Log("✅ Clicked edit button")

		// Wait for edit modal and delete button
		err = chromedp.Run(ctx,
			waitFor(`document.querySelector('button[lvt-click="delete"]') !== null`, 5*time.Second),

			// Check if delete button exists
			chromedp.Evaluate(`
				(() => {
					const deleteBtn = document.querySelector('button[lvt-click="delete"]');
					if (deleteBtn) {
						console.log('[DELETE TEST] Delete button found');
						console.log('[DELETE TEST] Delete button data-id:', deleteBtn.getAttribute('lvt-data-id'));
						console.log('[DELETE TEST] Delete button outerHTML:', deleteBtn.outerHTML.substring(0, 300));
						return true;
					}
					console.log('[DELETE TEST] Delete button NOT found');
					return false;
				})()
			`, &foundDeleteButton),
		)
		if err != nil || !foundDeleteButton {
			t.Fatalf("Delete button not found: %v", err)
		}
		t.Log("✅ Delete button found")

		// Override confirm and click delete
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`window.confirm = () => { console.log('[DELETE TEST] Confirm called'); return true; };`, nil),

			chromedp.Evaluate(`
				(() => {
					const deleteButton = document.querySelector('button[lvt-click="delete"]');
					if (deleteButton) {
						console.log('[DELETE TEST] Clicking delete button...');
						deleteButton.click();
						return true;
					}
					return false;
				})()
			`, &clickedDelete),
		)
		if err != nil || !clickedDelete {
			t.Fatalf("Failed to click delete button: %v", err)
		}
		t.Log("✅ Clicked delete button")

		// Wait for post to disappear
		err = chromedp.Run(ctx,
			waitFor(`
				(() => {
					const table = document.querySelector('table tbody');
					if (!table) return true;
					const rows = Array.from(table.querySelectorAll('tr'));
					console.log('[DELETE TEST] Waiting for deletion, current row count:', rows.length);
					const stillExists = rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.includes('Second Post');
					});
					console.log('[DELETE TEST] Second Post still exists:', stillExists);
					return !stillExists;
				})()
			`, 10*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to delete second post: %v", err)
		}
		t.Log("✅ Second post deleted successfully")
	})

	// Step 6: Verify only first post remains
	t.Run("Verify_Only_First_Post_Remains", func(t *testing.T) {
		ctx, cancel := createBrowserContext()
		defer cancel()
		ctx, timeoutCancel := context.WithTimeout(ctx, 30*time.Second)
		defer timeoutCancel()

		var postCount int
		var hasFirstPost bool
		var hasSecondPost bool

		err := chromedp.Run(ctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),

			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return 0;
					return table.querySelectorAll('tbody tr').length;
				})()
			`, &postCount),

			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'First Post - Edited';
					});
				})()
			`, &hasFirstPost),

			chromedp.Evaluate(`
				(() => {
					const table = document.querySelector('table');
					if (!table) return false;
					const rows = Array.from(table.querySelectorAll('tbody tr'));
					return rows.some(row => {
						const cells = row.querySelectorAll('td');
						return cells.length > 0 && cells[0].textContent.trim() === 'Second Post';
					});
				})()
			`, &hasSecondPost),
		)
		if err != nil {
			t.Fatalf("Failed to verify posts: %v", err)
		}

		t.Logf("Post count after deletion: %d", postCount)
		t.Logf("Has First Post: %v", hasFirstPost)
		t.Logf("Has Second Post: %v", hasSecondPost)

		if postCount != 1 {
			t.Errorf("Expected 1 post after deletion, got %d", postCount)
		}
		if !hasFirstPost {
			t.Error("First Post - Edited should still exist")
		}
		if hasSecondPost {
			t.Error("Second Post should have been deleted")
		}

		if postCount == 1 && hasFirstPost && !hasSecondPost {
			t.Log("✅ Only first post remains")
		}
	})
}
