//go:build browser

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cdpnetwork "github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

// startItemsApp generates a multi-kit app with an items resource (name:string,
// plus genArgs), migrates it, seeds rows, and runs it natively. It returns the
// app directory and port.
func startItemsApp(t *testing.T, name string, rows []struct {
	SQL  string
	Args []interface{}
}, genArgs ...string) (string, int) {
	t.Helper()
	appDir := createTestApp(t, t.TempDir(), name, &AppOptions{Kit: "multi"})
	setupLocalClientLibrary(t, appDir)

	if err := runLvtCommand(t, appDir, append([]string{"gen", "resource", "items", "name:string"}, genArgs...)...); err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	if err := seedTestData(filepath.Join(appDir, "app.db"), rows); err != nil {
		t.Fatalf("Failed to seed items: %v", err)
	}

	_ = injectFrameworkForTest(t, appDir)
	port := allocateTestPort()
	_ = buildAndRunNative(t, appDir, port)
	waitForServer(t, fmt.Sprintf("http://localhost:%d/items", port), 10*time.Second)
	return appDir, port
}

// browserDiag is a pooled-Chrome tab that records its console (uncaught
// exceptions included) and WebSocket frames, so a failing test can print them
// with the app's server log and the rendered HTML. It enables the network
// domain, so frames arrive from the first navigation.
type browserDiag struct {
	ctx     context.Context
	appDir  string
	console *e2etest.ConsoleLogger
	ws      *e2etest.WSMessageLogger
}

func newRecordedTab(t *testing.T, appDir string) *browserDiag {
	t.Helper()
	pool, _, release := GetPooledChrome(t)
	t.Cleanup(release)
	ctx, cancel := chromedp.NewContext(pool, chromedp.WithLogf(t.Logf))
	t.Cleanup(cancel)
	ctx, timeoutCancel := context.WithTimeout(ctx, getBrowserTimeout())
	t.Cleanup(timeoutCancel)

	if err := chromedp.Run(ctx, cdpnetwork.Enable()); err != nil {
		t.Fatalf("enable the network domain: %v", err)
	}
	console := e2etest.NewConsoleLogger()
	console.Start(ctx)
	return &browserDiag{ctx: ctx, appDir: appDir, console: console, ws: e2etest.RecordWSFrames(ctx)}
}

// sentAction reports whether the page sent a frame with this action.
func (d *browserDiag) sentAction(action string) bool {
	for _, m := range d.ws.GetSent() {
		if m.Parsed["action"] == action {
			return true
		}
	}
	return false
}

// dump returns the server log (buildAndRunNative's server.log), console, WS
// frames and rendered HTML, ready to append to a failure message.
func (d *browserDiag) dump() string {
	var html string
	_ = chromedp.Run(d.ctx, chromedp.OuterHTML(`body`, &html, chromedp.ByQuery))
	serverLog, _ := os.ReadFile(filepath.Join(d.appDir, "server.log"))

	var console []string
	for _, l := range d.console.GetLogs() {
		console = append(console, l.Type+" "+l.Message)
	}
	var frames []string
	for _, m := range d.ws.GetMessages() {
		data := m.Data
		if len(data) > 300 {
			data = data[:300] + "…"
		}
		frames = append(frames, m.Direction+" "+data)
	}
	return fmt.Sprintf("\n--- server ---\n%s\n--- console ---\n%s\n--- ws ---\n%s\n--- html ---\n%s",
		serverLog, strings.Join(console, "\n"), strings.Join(frames, "\n"), html)
}
