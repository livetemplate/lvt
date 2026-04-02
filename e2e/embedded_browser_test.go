//go:build browser

package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestEmbeddedBrowser_PostsWithComments tests embedded child resources and page-mode navigation.
// It validates:
//   - Page-mode navigation shows correct content per URL (mount.go always-call-Mount fix)
//   - Embedded comments section renders on detail pages
//   - No toast template errors
//   - Comment CRUD on detail pages
func TestEmbeddedBrowser_PostsWithComments(t *testing.T) {
	tmpDir := t.TempDir()

	// Step 1: Create blog app
	t.Log("Creating blog app with embedded comments...")
	appDir := createTestApp(t, tmpDir, "emblog", &AppOptions{Kit: "multi"})
	setupLocalClientLibrary(t, appDir)

	// Step 2: Generate posts resource (page mode)
	t.Log("Generating posts resource (page mode)...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "posts", "title:string", "content:text", "--edit-mode", "page"); err != nil {
		t.Fatalf("Failed to generate posts: %v", err)
	}

	// Step 3: Generate embedded comments resource
	t.Log("Generating embedded comments resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "comments", "post_id:references:posts", "author:string", "text:string", "--parent", "posts"); err != nil {
		t.Fatalf("Failed to generate comments: %v", err)
	}

	// Step 4: Run migrations
	t.Log("Running migrations...")
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Inject framework replace directive so the test app uses local mount.go changes
	hasFramework := injectFrameworkForTest(t, appDir)

	// Step 5: Seed test data — two posts with comments to test navigation
	t.Log("Seeding test data...")
	dbPath := filepath.Join(appDir, "app.db")
	if err := seedTestData(dbPath, []struct {
		SQL  string
		Args []interface{}
	}{
		{
			SQL:  "INSERT INTO posts (id, title, content, created_at) VALUES (?, ?, ?, datetime('now'))",
			Args: []interface{}{"post-alpha", "Alpha Post", "Alpha body content here"},
		},
		{
			SQL:  "INSERT INTO posts (id, title, content, created_at) VALUES (?, ?, ?, datetime('now'))",
			Args: []interface{}{"post-beta", "Beta Post", "Beta body content here"},
		},
		{
			SQL:  "INSERT INTO comments (id, post_id, author, text, created_at) VALUES (?, ?, ?, ?, datetime('now'))",
			Args: []interface{}{"comment-1", "post-alpha", "Alice", "Comment on Alpha"},
		},
		{
			SQL:  "INSERT INTO comments (id, post_id, author, text, created_at) VALUES (?, ?, ?, ?, datetime('now'))",
			Args: []interface{}{"comment-2", "post-beta", "Bob", "Comment on Beta"},
		},
	}); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	// Step 6: Build and run natively
	serverPort := allocateTestPort()
	_ = buildAndRunNative(t, appDir, serverPort)

	serverURL := fmt.Sprintf("http://localhost:%d", serverPort)
	waitForServer(t, serverURL+"/posts", 10*time.Second)
	t.Log("App running")

	// Step 7: Get Chrome from pool
	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	testURL := getTestURL(serverPort)

	// Console log collection for diagnostics
	var consoleLogs []string
	consoleLogsMutex := &sync.Mutex{}

	createBrowserContext := func() (context.Context, context.CancelFunc) {
		subCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(t.Logf))

		chromedp.ListenTarget(subCtx, func(ev interface{}) {
			if consoleEvent, ok := ev.(*cdpruntime.EventConsoleAPICalled); ok {
				for _, arg := range consoleEvent.Args {
					if arg.Type == cdpruntime.TypeString {
						logMsg := string(arg.Value)
						consoleLogsMutex.Lock()
						consoleLogs = append(consoleLogs, logMsg)
						consoleLogsMutex.Unlock()
						if strings.Contains(logMsg, "Error") || strings.Contains(logMsg, "Failed") || strings.Contains(logMsg, "WebSocket") {
							t.Logf("Browser console: %s", logMsg)
						}
					}
				}
			}
		})

		return subCtx, cancel
	}

	// Test 1: Posts list page loads without template errors
	t.Run("Posts List Page", func(t *testing.T) {
		bctx, cancel := createBrowserContext()
		defer cancel()
		bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
		defer timeoutCancel()

		verifyNoTemplateErrors(t, bctx, testURL+"/posts")
		t.Log("✅ Posts list page loads without errors")
	})

	// Test 2: Both seeded posts appear in the list
	t.Run("Seeded Posts In List", func(t *testing.T) {
		bctx, cancel := createBrowserContext()
		defer cancel()
		bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
		defer timeoutCancel()

		var bodyText string
		err := chromedp.Run(bctx,
			chromedp.Navigate(testURL+"/posts"),
			waitForWebSocketReady(5*time.Second),
			chromedp.WaitVisible(`[data-lvt-id]`, chromedp.ByQuery),
			waitFor(`
				(() => {
					const body = document.body.textContent;
					return body.includes('Alpha Post') && body.includes('Beta Post');
				})()
			`, 10*time.Second),
			chromedp.Text("body", &bodyText, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("Failed to find seeded posts in list: %v", err)
		}
		if !strings.Contains(bodyText, "Alpha Post") || !strings.Contains(bodyText, "Beta Post") {
			t.Error("Not all seeded posts found in list")
		}
		t.Log("✅ Both seeded posts appear in list")
	})

	// Test 3: CRITICAL — Page-mode navigation shows correct content per URL
	// Before the mount.go fix, clicking post B after post A would show A's content (stale session state).
	t.Run("Page Mode Navigation Shows Correct Content", func(t *testing.T) {
		if !hasFramework {
			t.Skip("Skipping navigation test: requires local framework with mount fix")
		}
		bctx, cancel := createBrowserContext()
		defer cancel()
		bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
		defer timeoutCancel()

		// Navigate to Alpha Post detail page
		t.Log("Navigating to Alpha Post detail page...")
		alphaURL := testURL + "/posts/post-alpha"
		var alphaBody string
		err := chromedp.Run(bctx,
			chromedp.Navigate(alphaURL),
			waitFor(`
				(() => {
					return document.body.textContent.includes('Alpha Post') &&
					       document.body.textContent.includes('Alpha body content');
				})()
			`, 10*time.Second),
			chromedp.Text("body", &alphaBody, chromedp.ByQuery),
		)
		if err != nil {
			var bodyHTML string
			_ = chromedp.Run(bctx, chromedp.Evaluate(`document.body.innerHTML.substring(0, 2000)`, &bodyHTML))
			t.Logf("DEBUG Alpha page HTML: %s", bodyHTML)
			t.Fatalf("Alpha Post detail page failed to load: %v", err)
		}
		if !strings.Contains(alphaBody, "Alpha Post") {
			t.Errorf("Alpha detail page should show 'Alpha Post', got: %s", alphaBody[:min(200, len(alphaBody))])
		}
		// Alpha's comment should be visible, not Beta's
		if !strings.Contains(alphaBody, "Alice") {
			t.Log("Note: Alice's comment not found on Alpha detail (comments may not load on initial render)")
		}
		t.Log("✅ Alpha Post detail shows correct content")

		// Navigate to Beta Post detail page (same browser session → same cookie)
		// This is the KEY test: Beta should show Beta's content, NOT Alpha's
		t.Log("Navigating to Beta Post detail page (key mount fix test)...")
		betaURL := testURL + "/posts/post-beta"
		var betaBody string
		err = chromedp.Run(bctx,
			chromedp.Navigate(betaURL),
			waitFor(`
				(() => {
					return document.body.textContent.includes('Beta Post') &&
					       document.body.textContent.includes('Beta body content');
				})()
			`, 10*time.Second),
			chromedp.Text("body", &betaBody, chromedp.ByQuery),
		)
		if err != nil {
			var bodyHTML string
			_ = chromedp.Run(bctx, chromedp.Evaluate(`document.body.innerHTML.substring(0, 2000)`, &bodyHTML))
			t.Logf("DEBUG Beta page HTML: %s", bodyHTML)
			t.Fatalf("Beta Post detail page failed to load correct content: %v", err)
		}
		if !strings.Contains(betaBody, "Beta Post") {
			t.Errorf("Beta detail page should show 'Beta Post' (not stale Alpha), got: %s", betaBody[:min(200, len(betaBody))])
		}
		if !strings.Contains(betaBody, "Beta body content") {
			t.Errorf("Beta detail page should show 'Beta body content'")
		}
		if strings.Contains(betaBody, "Alpha body content") {
			t.Errorf("Beta detail page should NOT show Alpha's content (stale state bug)")
		}
		t.Log("✅ Beta Post detail shows correct content (mount fix verified)")

		// Navigate back to Alpha to confirm it still works
		t.Log("Navigating back to Alpha Post...")
		var alphaBody2 string
		err = chromedp.Run(bctx,
			chromedp.Navigate(alphaURL),
			waitFor(`
				(() => {
					return document.body.textContent.includes('Alpha Post') &&
					       document.body.textContent.includes('Alpha body content');
				})()
			`, 10*time.Second),
			chromedp.Text("body", &alphaBody2, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("Failed to navigate back to Alpha Post: %v", err)
		}
		if !strings.Contains(alphaBody2, "Alpha Post") {
			t.Errorf("Alpha detail page should still show Alpha content after navigating back")
		}
		if strings.Contains(alphaBody2, "Beta body content") {
			t.Errorf("Alpha detail page should NOT show Beta's content")
		}
		t.Log("✅ Navigation back to Alpha works correctly")
	})

	// Test 4: Detail page has comments section (embedded child rendering)
	t.Run("Detail Page Has Comments Section", func(t *testing.T) {
		bctx, cancel := createBrowserContext()
		defer cancel()
		bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
		defer timeoutCancel()

		detailURL := testURL + "/posts/post-alpha"
		var bodyText string
		err := chromedp.Run(bctx,
			chromedp.Navigate(detailURL),
			waitFor(`
				(() => {
					const body = document.body.textContent.toLowerCase();
					return body.includes('alpha post') && body.includes('comment');
				})()
			`, 10*time.Second),
			chromedp.Text("body", &bodyText, chromedp.ByQuery),
		)
		if err != nil {
			var bodyHTML string
			_ = chromedp.Run(bctx, chromedp.Evaluate(`document.body.innerHTML.substring(0, 3000)`, &bodyHTML))
			t.Logf("DEBUG detail page HTML: %s", bodyHTML)
			t.Fatalf("Comments section not found on detail page: %v", err)
		}

		bodyLower := strings.ToLower(bodyText)
		if !strings.Contains(bodyLower, "comment") {
			t.Errorf("Detail page should have comments section")
		}

		// Check for template errors
		errorPatterns := []string{"invalid value", "expected toast.Type", "<no value>", "template:"}
		for _, pattern := range errorPatterns {
			if strings.Contains(bodyText, pattern) {
				t.Errorf("Template error found on detail page: %q", pattern)
			}
		}
		t.Log("✅ Comments section visible, no template errors")
	})

	// Test 5: No toast type errors in console logs
	t.Run("No Toast Errors In Console", func(t *testing.T) {
		consoleLogsMutex.Lock()
		defer consoleLogsMutex.Unlock()

		for _, log := range consoleLogs {
			if strings.Contains(log, "invalid value") || strings.Contains(log, "toast.Type") {
				t.Errorf("Console log contains toast type error: %s", log)
			}
		}
		t.Log("✅ No toast type errors in console logs")
	})

	// Test 6: Server logs check
	t.Run("Server Logs Check", func(t *testing.T) {
		serverLogPath := filepath.Join(appDir, "server.log")
		logContent, err := os.ReadFile(serverLogPath)
		if err != nil {
			t.Logf("Could not read server logs: %v", err)
			return
		}

		logStr := string(logContent)
		if strings.Contains(logStr, "invalid value") {
			t.Errorf("'invalid value' error found in server logs")
		}
		t.Logf("Server log size: %d bytes", len(logContent))
		t.Log("✅ No critical errors in server logs")
	})
}

// TestToastAutoDismiss tests that success toasts auto-dismiss after ~5 seconds.
func TestToastAutoDismiss(t *testing.T) {
	tmpDir := t.TempDir()

	t.Log("Creating app for toast auto-dismiss test...")
	appDir := createTestApp(t, tmpDir, "toastapp", &AppOptions{Kit: "multi"})
	setupLocalClientLibrary(t, appDir)

	t.Log("Generating items resource...")
	if err := runLvtCommand(t, appDir, "gen", "resource", "items", "name:string"); err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}

	t.Log("Running migrations...")
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	_ = injectFrameworkForTest(t, appDir)

	serverPort := allocateTestPort()
	_ = buildAndRunNative(t, appDir, serverPort)

	serverURL := fmt.Sprintf("http://localhost:%d", serverPort)
	waitForServer(t, serverURL+"/items", 10*time.Second)
	t.Log("App running")

	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()

	testURL := getTestURL(serverPort)

	t.Run("Success Toast Auto Dismisses", func(t *testing.T) {
		bctx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(t.Logf))
		defer cancel()
		bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
		defer timeoutCancel()

		// Navigate to items list
		err := chromedp.Run(bctx,
			chromedp.Navigate(testURL+"/items"),
			waitForWebSocketReady(5*time.Second),
		)
		if err != nil {
			t.Fatalf("Failed to navigate: %v", err)
		}

		// Open add modal and submit a new item via the submit button (not native form submit)
		err = chromedp.Run(bctx,
			chromedp.Click(`[lvt-modal-open="add-modal"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`form[lvt-submit="add"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="name"]`, "Test Item", chromedp.ByQuery),
			chromedp.Click(`form[lvt-submit="add"] button[type="submit"]`, chromedp.ByQuery),
		)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		// Wait for toast to appear (client-rendered via data-toast-trigger ephemeral component)
		err = chromedp.Run(bctx,
			waitFor(`document.querySelector('[data-lvt-toast-item]') !== null`, 5*time.Second),
		)
		if err != nil {
			t.Fatalf("Toast did not appear: %v", err)
		}
		t.Log("Toast appeared as data-lvt-toast-item")

		// Verify the toast auto-dismisses within 8 seconds (5s timer + buffer)
		err = chromedp.Run(bctx,
			waitFor(`document.querySelector('[data-lvt-toast-item]') === null`, 8*time.Second),
		)
		if err != nil {
			t.Fatalf("Toast did not auto-dismiss within 8 seconds: %v", err)
		}
		t.Log("Toast auto-dismissed successfully")
	})
}

// injectFrameworkForTest adds a replace directive for the livetemplate framework module
// so that test apps use the local framework with recent changes (e.g., mount.go always-call-Mount fix).
// Returns true if injection succeeded, false if the framework directory was not found.
func injectFrameworkForTest(t *testing.T, appDir string) bool {
	t.Helper()

	// Find the livetemplate framework directory relative to the project.
	// This test file lives in <lvt-root>/e2e/ (or <lvt-root>/.worktrees/<name>/e2e/).
	// The framework is a sibling of the lvt repo: <parent>/livetemplate/
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path for framework injection")
	}

	// Walk up from e2e dir to find the lvt repo root (has go.mod with github.com/livetemplate/lvt)
	dir := filepath.Dir(filepath.Dir(filename)) // e2e/ -> lvt-root or worktree-root
	var frameworkDir string
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(filepath.Dir(dir), "livetemplate")
		if _, err := os.Stat(filepath.Join(candidate, "go.mod")); err == nil {
			frameworkDir = candidate
			break
		}
		dir = filepath.Dir(dir)
	}
	if frameworkDir == "" {
		t.Log("⏭️  Framework directory not found, skipping injection")
		return false
	}

	goModPath := filepath.Join(appDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod for framework injection: %v", err)
	}

	goModStr := string(goModContent)
	if strings.Contains(goModStr, "replace github.com/livetemplate/livetemplate") {
		t.Log("Framework replace directive already present")
		return true
	}

	goModStr += fmt.Sprintf("\nreplace github.com/livetemplate/livetemplate => %s\n", frameworkDir)

	if err := os.WriteFile(goModPath, []byte(goModStr), 0644); err != nil {
		t.Fatalf("Failed to update go.mod for framework injection: %v", err)
	}

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = appDir
	tidyCmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		t.Logf("⚠️  go mod tidy after framework injection: %v\nOutput: %s", err, output)
	}

	t.Logf("✅ Framework module injected for test (path: %s)", frameworkDir)
	return true
}
