//go:build browser

package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/livetemplate/livetemplate"
	"github.com/livetemplate/lvt/e2e"
	e2etest "github.com/livetemplate/lvt/testing"
)

// wsUpgradeDelay holds every WebSocket handshake, so the page stays clickable
// with no open socket for this long, as it does on a slow remote site.
const wsUpgradeDelay = 5 * time.Second

// TestE2E_AttributeClickBeforeConnect: a handler registered through the queue
// arms on the server-rendered markup before the WebSocket opens, and a click
// in that window reaches the Controller exactly once after the socket opens.
// Without that, the click is lost: on a remote site the page is interactive
// well before the handshake finishes.
//
// Set LVT_LOCAL_CLIENT_JS to a built client bundle to run this against an
// unreleased client (see e2etest.GetClientLibraryJS).
func TestE2E_AttributeClickBeforeConnect(t *testing.T) {
	serverLogger := e2etest.NewServerLogger()
	serverLogger.Start()
	defer serverLogger.Stop()
	// slog.SetDefault is process-global, so this test must not run in parallel.
	prevSlog := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(os.Stderr, serverLogger.Writer()),
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))
	defer slog.SetDefault(prevSlog)

	port, shutdown := startBeforeConnectServer(t)
	defer shutdown()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 60*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context

	consoleLog := e2etest.NewConsoleLogger()
	consoleLog.Start(ctx)
	wsLog := e2etest.RecordWSFrames(ctx)

	dump := func() {
		var lines []string
		for _, l := range consoleLog.GetLogs() {
			lines = append(lines, l.Type+": "+strings.Join(l.Args, " "))
		}
		t.Logf("──────── CONSOLE ────────\n%s", strings.Join(lines, "\n"))
		t.Log("──────── SERVER LOGS ────────")
		for _, l := range serverLogger.GetLogs() {
			t.Log(l)
		}
		t.Log("──────── WS FRAMES (last 20) ────────")
		wsLog.PrintLast(20)
		var html string
		if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &html, chromedp.ByQuery)); err == nil {
			t.Logf("──────── RENDERED HTML ────────\n%s", html)
		}
	}

	if err := chromedp.Run(ctx,
		chromedp.Navigate(e2etest.GetChromeTestURL(port)),
		chromedp.WaitVisible("#rate", chromedp.ByID),
	); err != nil {
		dump()
		t.Fatalf("navigate: %v", err)
	}
	e2e.PollUntil(t, ctx, `window.__rateArmed === 1`, wsUpgradeDelay/2,
		"the queued handler to arm the server-rendered button before the WebSocket opened", dump)

	if err := chromedp.Run(ctx, chromedp.Click("#rate", chromedp.ByID)); err != nil {
		dump()
		t.Fatalf("click: %v", err)
	}
	if n := wsLog.Count(); n != 0 {
		dump()
		t.Fatalf("the click was not inside the pre-connect window: %d WS frames before it", n)
	}

	e2e.PollUntil(t, ctx, `document.getElementById("count").textContent.trim() === "1"`,
		wsUpgradeDelay+10*time.Second, "a click made before the WebSocket opened to reach the Controller", dump)

	// Exactly once: a queued send must not also go out over HTTP. A duplicate
	// would land within a round trip, so a short wait catches it, but a later
	// one would slip past.
	time.Sleep(500 * time.Millisecond)
	var count string
	if err := chromedp.Run(ctx, chromedp.Text("#count", &count, chromedp.ByID)); err != nil {
		dump()
		t.Fatalf("read count: %v", err)
	}
	if strings.TrimSpace(count) != "1" {
		dump()
		t.Errorf("count = %q after one pre-connect click, want 1", count)
	}
}

type beforeConnectState struct{ Count int }

type beforeConnectController struct{}

func (c *beforeConnectController) Rate(s beforeConnectState, _ *livetemplate.Context) (beforeConnectState, error) {
	s.Count++
	return s, nil
}

const beforeConnectTemplate = `<!DOCTYPE html>
<html>
<head>
<title>Click before connect</title>
<script>
  window.__rateArmed = 0;
  (window.lvtAttributes ||= []).push({
    attribute: "lvt-x:e2e-rate",
    needsServerChannel: true,
    onElementAdded(el, ctx) {
      window.__rateArmed++;
      el.addEventListener("click", () => ctx.send?.({ action: ctx.value }), { signal: ctx.signal });
    },
  });
</script>
</head>
<body>
  <div id="count">{{.Count}}</div>
  <button id="rate" type="button" lvt-x:e2e-rate="Rate">rate</button>
  <script src="/client.js" defer></script>
</body>
</html>`

func startBeforeConnectServer(t *testing.T) (int, func()) {
	t.Helper()

	tmpl := livetemplate.Must(livetemplate.New("before-connect-e2e"))
	if _, err := tmpl.Parse(beforeConnectTemplate); err != nil {
		t.Fatalf("parse template: %v", err)
	}
	app := tmpl.Handle(&beforeConnectController{}, livetemplate.AsState(&beforeConnectState{}))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			slog.Info("delaying WebSocket upgrade", slog.Duration("delay", wsUpgradeDelay))
			select {
			case <-time.After(wsUpgradeDelay):
			case <-r.Context().Done():
				return
			}
		}
		app.ServeHTTP(w, r)
	}))
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	p, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("before-connect server failed", slog.Any("error", err))
		}
	}()
	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", p), 10*time.Second)
	return p, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("before-connect server shutdown failed", slog.Any("error", err))
		}
	}
}
