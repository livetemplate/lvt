//go:build browser

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cdpnetwork "github.com/chromedp/cdproto/network"
	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// browserDiag records a tab's console and WebSocket frames so a failing test
// can print them with the app's server log and the rendered HTML. Run
// cdpnetwork.Enable() on the tab before navigating, or no frames arrive.
type browserDiag struct {
	ctx    context.Context
	appDir string

	mu       sync.Mutex
	console  []string
	wsFrames []string
}

func recordBrowser(ctx context.Context, appDir string) *browserDiag {
	d := &browserDiag{ctx: ctx, appDir: appDir}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		d.mu.Lock()
		defer d.mu.Unlock()
		switch e := ev.(type) {
		case *cdpruntime.EventConsoleAPICalled:
			parts := []string{string(e.Type)}
			for _, a := range e.Args {
				parts = append(parts, string(a.Value))
			}
			d.console = append(d.console, strings.Join(parts, " "))
		case *cdpnetwork.EventWebSocketFrameSent:
			d.wsFrames = append(d.wsFrames, "sent "+e.Response.PayloadData)
		case *cdpnetwork.EventWebSocketFrameReceived:
			frame := e.Response.PayloadData
			if len(frame) > 300 {
				frame = frame[:300] + "…"
			}
			d.wsFrames = append(d.wsFrames, "recv "+frame)
		}
	})
	return d
}

// sent joins the frames the page sent, for asserting on actions.
func (d *browserDiag) sent() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []string
	for _, f := range d.wsFrames {
		if strings.HasPrefix(f, "sent ") {
			out = append(out, f)
		}
	}
	return strings.Join(out, "\n")
}

// dump returns the server log (buildAndRunNative's server.log), console, WS
// frames and rendered HTML, ready to append to a failure message.
func (d *browserDiag) dump() string {
	var html string
	_ = chromedp.Run(d.ctx, chromedp.OuterHTML(`body`, &html, chromedp.ByQuery))
	serverLog, _ := os.ReadFile(filepath.Join(d.appDir, "server.log"))
	d.mu.Lock()
	defer d.mu.Unlock()
	return fmt.Sprintf("\n--- server ---\n%s\n--- console ---\n%s\n--- ws ---\n%s\n--- html ---\n%s",
		serverLog, strings.Join(d.console, "\n"), strings.Join(d.wsFrames, "\n"), html)
}
