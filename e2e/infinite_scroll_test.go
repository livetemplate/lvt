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

	cdpnetwork "github.com/chromedp/cdproto/network"
	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestInfiniteScrollLoadsMore checks a generated app's infinite pagination in a
// browser: scrolling the [lvt-scroll-sentinel] into view sends load_more and
// renders the next page. The client's observer does this; the kit carries no
// scroll script of its own.
func TestInfiniteScrollLoadsMore(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "scrollapp", &AppOptions{Kit: "multi"})
	setupLocalClientLibrary(t, appDir)

	if err := runLvtCommand(t, appDir, "gen", "resource", "items", "name:string", "--page-size", "10"); err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	var rows []struct {
		SQL  string
		Args []interface{}
	}
	for i := 1; i <= 25; i++ {
		rows = append(rows, struct {
			SQL  string
			Args []interface{}
		}{
			SQL:  "INSERT INTO items (id, name, created_at) VALUES (?, ?, datetime('now', ?))",
			Args: []interface{}{fmt.Sprintf("item-%02d", i), fmt.Sprintf("Item %02d", i), fmt.Sprintf("-%d minutes", i)},
		})
	}
	if err := seedTestData(filepath.Join(appDir, "app.db"), rows); err != nil {
		t.Fatalf("Failed to seed items: %v", err)
	}

	_ = injectFrameworkForTest(t, appDir)
	serverPort := allocateTestPort()
	_ = buildAndRunNative(t, appDir, serverPort)
	waitForServer(t, fmt.Sprintf("http://localhost:%d/items", serverPort), 10*time.Second)

	ctx, _, cleanup := GetPooledChrome(t)
	defer cleanup()
	bctx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(t.Logf))
	defer cancel()
	bctx, timeoutCancel := context.WithTimeout(bctx, getBrowserTimeout())
	defer timeoutCancel()

	var mu sync.Mutex
	var consoleLines, wsFrames []string
	chromedp.ListenTarget(bctx, func(ev interface{}) {
		mu.Lock()
		defer mu.Unlock()
		switch e := ev.(type) {
		case *cdpruntime.EventConsoleAPICalled:
			parts := []string{string(e.Type)}
			for _, a := range e.Args {
				parts = append(parts, string(a.Value))
			}
			consoleLines = append(consoleLines, strings.Join(parts, " "))
		case *cdpnetwork.EventWebSocketFrameSent:
			wsFrames = append(wsFrames, "sent "+e.Response.PayloadData)
		case *cdpnetwork.EventWebSocketFrameReceived:
			frame := e.Response.PayloadData
			if len(frame) > 300 {
				frame = frame[:300] + "…"
			}
			wsFrames = append(wsFrames, "recv "+frame)
		}
	})
	diag := func() string {
		var html string
		_ = chromedp.Run(bctx, chromedp.OuterHTML(`body`, &html, chromedp.ByQuery))
		serverLog, _ := os.ReadFile(filepath.Join(appDir, "server.log"))
		mu.Lock()
		defer mu.Unlock()
		return fmt.Sprintf("\n--- server ---\n%s\n--- console ---\n%s\n--- ws ---\n%s\n--- html ---\n%s",
			serverLog, strings.Join(consoleLines, "\n"), strings.Join(wsFrames, "\n"), html)
	}

	const rowCount = `document.querySelectorAll('tbody tr[data-key]').length`
	if err := chromedp.Run(bctx,
		cdpnetwork.Enable(),
		// Short enough that the first page's sentinel starts below the fold.
		chromedp.EmulateViewport(1200, 400),
		chromedp.Navigate(getTestURL(serverPort)+"/items"),
		waitForWebSocketReady(5*time.Second),
		waitFor(rowCount+` === 10`, 5*time.Second),
		chromedp.Sleep(time.Second),
		waitFor(rowCount+` === 10`, time.Second),
	); err != nil {
		t.Fatalf("first page didn't show 10 rows, or loaded more without a scroll: %v%s", err, diag())
	}

	if err := chromedp.Run(bctx,
		chromedp.Evaluate(`document.querySelector('[lvt-scroll-sentinel]').scrollIntoView()`, nil),
		waitFor(rowCount+` === 20`, 10*time.Second),
	); err != nil {
		t.Fatalf("scrolling to the sentinel didn't load the second page: %v%s", err, diag())
	}

	mu.Lock()
	sent := strings.Join(wsFrames, "\n")
	mu.Unlock()
	if !strings.Contains(sent, `"load_more"`) {
		t.Errorf("no load_more action went over the socket%s", diag())
	}
}
