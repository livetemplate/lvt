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
	"sync/atomic"
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
	h := startBeforeConnectHarness(t)

	e2e.PollUntil(t, h.ctx, `window.__rateArmed === 1`, wsUpgradeDelay,
		"the queued handler to arm the server-rendered button before the WebSocket opened", h.dump)
	h.clickBeforeConnect("#rate")
	h.expectOneWSAction()
}

// TestE2E_LvtOnClickBeforeConnect: an lvt-on:click made before the
// WebSocket opens is queued and reaches the Controller exactly once after the
// socket opens, the same as a custom attribute handler's send.
func TestE2E_LvtOnClickBeforeConnect(t *testing.T) {
	h := startBeforeConnectHarness(t)

	h.clickBeforeConnect("#inc")
	h.expectOneWSAction()
}

// TestE2E_FormSubmitBeforeConnectIsNative: with the loading indicator off, a
// form submitted before the WebSocket opens is left to the browser, so it
// reaches the server as one native POST even if the socket never opens. The
// client must not also queue it and send it again over the WebSocket. With the
// indicator on (the default), the client disables form controls until the
// first render, so there is no pre-connect submit to test.
func TestE2E_FormSubmitBeforeConnectIsNative(t *testing.T) {
	h := startBeforeConnectHarness(t, livetemplate.WithLoadingDisabled())

	h.clickBeforeConnect("#submit")
	if !waitFor(func() bool { return h.posts.Load() >= 1 }, 10*time.Second) {
		h.dump()
		t.Fatal("a form submitted before the WebSocket opened did not POST natively")
	}
	// The POST loads a new page, whose socket also waits wsUpgradeDelay. Wait
	// for it to open, so a queued copy of the submit would have been sent.
	if !waitFor(func() bool { return h.wsLog.CountByDirection("received") >= 1 }, wsUpgradeDelay+10*time.Second) {
		h.dump()
		t.Fatal("the page loaded by the POST never opened its WebSocket")
	}
	time.Sleep(500 * time.Millisecond)
	if n := h.wsLog.CountByDirection("sent"); n != 0 {
		h.dump()
		t.Errorf("%d WS frames sent, want 0: the submit already went out as a native POST", n)
	}
	if n := h.posts.Load(); n != 1 {
		h.dump()
		t.Errorf("%d HTTP POSTs, want 1", n)
	}
}

type beforeConnectHarness struct {
	t     *testing.T
	ctx   context.Context
	posts *atomic.Int32
	wsLog *e2etest.WSMessageLogger
	dump  func()
}

// startBeforeConnectHarness starts the app with a delayed WebSocket upgrade
// and loads it in Chrome, recording console logs, server logs and WS frames.
// slog.SetDefault is process-global, so these tests must not run in parallel.
func startBeforeConnectHarness(t *testing.T, opts ...livetemplate.Option) *beforeConnectHarness {
	t.Helper()
	serverLogger := e2etest.NewServerLogger()
	serverLogger.Start()
	t.Cleanup(serverLogger.Stop)
	prevSlog := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(os.Stderr, serverLogger.Writer()),
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))
	t.Cleanup(func() { slog.SetDefault(prevSlog) })

	port, posts, shutdown := startBeforeConnectServer(t, opts...)
	t.Cleanup(shutdown)
	chromeCtx, chromeCleanup := e2etest.SetupDockerChrome(t, 60*time.Second)
	t.Cleanup(chromeCleanup)
	ctx := chromeCtx.Context

	consoleLog := e2etest.NewConsoleLogger()
	consoleLog.Start(ctx)
	wsLog := e2etest.RecordWSFrames(ctx)

	h := &beforeConnectHarness{t: t, ctx: ctx, posts: posts, wsLog: wsLog}
	h.dump = func() {
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
		h.dump()
		t.Fatalf("navigate: %v", err)
	}
	return h
}

// clickBeforeConnect clicks sel and fails if the WebSocket had already
// carried a frame, since the click would then not test the pre-connect window.
func (h *beforeConnectHarness) clickBeforeConnect(sel string) {
	h.t.Helper()
	if err := chromedp.Run(h.ctx, chromedp.Click(sel, chromedp.ByQuery)); err != nil {
		h.dump()
		h.t.Fatalf("click %s: %v", sel, err)
	}
	if n := h.wsLog.Count(); n != 0 {
		h.dump()
		h.t.Fatalf("the click was not inside the pre-connect window: %d WS frames before it (a slow runner, not the client)", n)
	}
}

// expectOneWSAction waits for the pre-connect click to reach the Controller
// and checks it went out over the WebSocket only.
func (h *beforeConnectHarness) expectOneWSAction() {
	h.t.Helper()
	e2e.PollUntil(h.t, h.ctx, `document.getElementById("count").textContent.trim() === "1"`,
		wsUpgradeDelay+10*time.Second, "a click made before the WebSocket opened to reach the Controller", h.dump)

	// Exactly once: a queued send must not also go out over HTTP or twice over
	// the WebSocket. A duplicate would land within a round trip, so a short
	// wait catches it, but a later one would slip past.
	time.Sleep(500 * time.Millisecond)
	if n := h.posts.Load(); n != 0 {
		h.dump()
		h.t.Errorf("%d HTTP POSTs, want 0: the queued send must go out over the WebSocket only", n)
	}
	if n := h.sentActions(); n != 1 {
		h.dump()
		h.t.Errorf("%d WS action frames, want 1", n)
	}
	var count string
	if err := chromedp.Run(h.ctx, chromedp.Text("#count", &count, chromedp.ByID)); err != nil {
		h.dump()
		h.t.Fatalf("read count: %v", err)
	}
	if strings.TrimSpace(count) != "1" {
		h.dump()
		h.t.Errorf("count = %q after one pre-connect click, want 1", count)
	}
}

// sentActions counts the WS frames the page sent that carry the rate action.
func (h *beforeConnectHarness) sentActions() int {
	n := 0
	for _, m := range h.wsLog.GetMessages() {
		if m.Direction == "sent" && strings.Contains(strings.ToLower(m.Data), `"action":"rate"`) {
			n++
		}
	}
	return n
}

func waitFor(cond func() bool, timeout time.Duration) bool {
	for deadline := time.Now().Add(timeout); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		if cond() {
			return true
		}
	}
	return cond()
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
  <button id="inc" type="button" lvt-on:click="Rate">inc</button>
  <form method="POST"><button id="submit" name="rate">submit</button></form>
  <script src="/client.js" defer></script>
</body>
</html>`

// startBeforeConnectServer returns the port, a count of POSTs the app
// received, and a shutdown func.
func startBeforeConnectServer(t *testing.T, opts ...livetemplate.Option) (int, *atomic.Int32, func()) {
	t.Helper()

	tmpl := livetemplate.Must(livetemplate.New("before-connect-e2e", opts...))
	if _, err := tmpl.Parse(beforeConnectTemplate); err != nil {
		t.Fatalf("parse template: %v", err)
	}
	app := tmpl.Handle(&beforeConnectController{}, livetemplate.AsState(&beforeConnectState{}))

	mux := http.NewServeMux()
	var posts atomic.Int32
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			posts.Add(1)
		}
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
	return p, &posts, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("before-connect server shutdown failed", slog.Any("error", err))
		}
	}
}
