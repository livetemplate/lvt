//go:build browser

// NOTE: this file declares `package e2e_test` (external test package) rather
// than `package e2e` because the `e2e` package already contains a different
// `waitForCondition` helper (in `helpers.go`) and a `waitForCondition` helper (in
// `rendering_test.go`) with chromedp.Action signatures. Moving this file in
// would require renaming the fatal-test-helper polling loop, and the directory
// already mixes both package forms (e.g., `livetemplate_core_test.go` is
// `package e2e_test` too).
package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/websocket"
	"github.com/livetemplate/livetemplate"
	e2etest "github.com/livetemplate/lvt/testing"
)

// =============================================================================
// E2E coverage for the browser-observable parts of the lifecycle-ergonomics
// batch: #339 (ErrSessionDisconnected), #340 (IsInitialMount / IsReconnect),
// #345 (ClearAllFlash). Issue #341 (validateLifecycleSignatures warn-at-boot)
// is server-startup-time, not browser-observable, and is covered by the
// upstream livetemplate unit tests.
//
// Verifies behavior end-to-end against a real Chromium and a real livetemplate
// http.Handler. Per CLAUDE.md (global): chromedp e2e tests surface (1) browser
// console, (2) server logs, (3) websocket messages, (4) rendered HTML.
//
// Important semantic note observed during test development:
//
//   In a normal browser flow — HTTP GET → server-renders HTML → WS connects —
//   the framework persists state at the end of the HTTP Mount. When the WS
//   then connects, `restorePersistedState` returns ok, which is the signal
//   that flips IsReconnect() to true. So even on a "fresh" tab load,
//   `ctx.IsReconnect()` is true inside the WS-path Mount/OnConnect.
//
//   IsInitialMount(), by contrast, is only true on the actual HTTP GET. By
//   the time the WS has connected, OnConnect has overwritten state with
//   IsInitialMount=false. To assert IsInitialMount=true you must inspect the
//   server-rendered HTML before the WS arrives — these tests use http.Get
//   for that, not the chromedp DOM.
// =============================================================================

type lifecycleE2EState struct {
	// Count is persisted so a WS reconnect restores prior state — that
	// restoration is precisely what IsReconnect() reads.
	Count int `lvt:"persist"`

	// Tags carry the connect-kind helpers' values out to the DOM/HTML so the
	// test can read them. Persist them too so the *most recent* lifecycle
	// call wins regardless of which path (HTTP vs WS) ran last.
	InitialMountTag string `lvt:"persist"`
	ReconnectTag    string `lvt:"persist"`
}

type lifecycleE2EController struct{}

func (c *lifecycleE2EController) Mount(state lifecycleE2EState, ctx *livetemplate.Context) (lifecycleE2EState, error) {
	state.InitialMountTag = boolTag(ctx.IsInitialMount())
	state.ReconnectTag = boolTag(ctx.IsReconnect())
	return state, nil
}

func (c *lifecycleE2EController) OnConnect(state lifecycleE2EState, ctx *livetemplate.Context) (lifecycleE2EState, error) {
	// OnConnect overwrites the tags so the DOM reflects the WS-path classification.
	state.InitialMountTag = boolTag(ctx.IsInitialMount())
	state.ReconnectTag = boolTag(ctx.IsReconnect())
	state.Count++
	return state, nil
}

func (c *lifecycleE2EController) SetFlashes(state lifecycleE2EState, ctx *livetemplate.Context) (lifecycleE2EState, error) {
	ctx.SetFlash("success", "Saved!")
	ctx.SetFlash("info", "Heads up")
	return state, nil
}

func (c *lifecycleE2EController) ClearFlashes(state lifecycleE2EState, ctx *livetemplate.Context) (lifecycleE2EState, error) {
	ctx.ClearAllFlash()
	return state, nil
}

func boolTag(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

const lifecycleE2ETemplate = `<!DOCTYPE html>
<html>
<head><title>Lifecycle E2E</title></head>
<body>
  <div id="im">im={{.InitialMountTag}}</div>
  <div id="rc">rc={{.ReconnectTag}}</div>
  <div id="count">count={{.Count}}</div>
  <div id="flash-success">{{.lvt.Flash "success"}}</div>
  <div id="flash-info">{{.lvt.Flash "info"}}</div>
  <button id="set-flashes" lvt-on:click="SetFlashes">Set flashes</button>
  <button id="clear-flashes" lvt-on:click="ClearFlashes">Clear all</button>
  <script src="/client.js"></script>
</body>
</html>`

// startLifecycleE2EServer boots an in-process http.Server with the test
// handler. Mirrors TestLoadingIndicator's setup pattern.
func startLifecycleE2EServer(t *testing.T) (baseURL string, port int, shutdown func()) {
	t.Helper()

	tmpl := livetemplate.Must(livetemplate.New("lifecycle-e2e"))
	if _, err := tmpl.Parse(lifecycleE2ETemplate); err != nil {
		t.Fatalf("parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tmpl.Handle(&lifecycleE2EController{}, livetemplate.AsState(&lifecycleE2EState{})))
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	p, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server: %v", err)
		}
	}()

	url := fmt.Sprintf("http://localhost:%d", p)
	e2etest.WaitForServer(t, url, 10*time.Second)

	return url, p, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
}

// installConsoleLogger pipes browser console.* calls into t.Logf so failures
// are debuggable without a re-run.
func installConsoleLogger(t *testing.T, ctx context.Context) {
	t.Helper()
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if call, ok := ev.(*runtime.EventConsoleAPICalled); ok {
			parts := make([]string, 0, len(call.Args))
			for _, a := range call.Args {
				parts = append(parts, string(a.Value))
			}
			t.Logf("[console.%s] %s", call.Type, strings.Join(parts, " "))
		}
	})
}

// waitForCondition polls a JS expression until it returns true or the timeout
// elapses. Used for "DOM should reflect X after WS round-trip" assertions.
//
// A cancelled chromedp context (e.g. browser crash, Docker shutdown) is fatal
// immediately — without that early return the loop would silently treat
// context.Canceled as "condition not yet true" and report a misleading timeout.
func waitForCondition(t *testing.T, ctx context.Context, jsExpr string, timeout time.Duration, why string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var ok bool
		err := chromedp.Run(ctx, chromedp.Evaluate(jsExpr, &ok))
		if err == nil && ok {
			return
		}
		if ctx.Err() != nil {
			t.Fatalf("chromedp context cancelled waiting for %s: %v", why, ctx.Err())
		}
		time.Sleep(50 * time.Millisecond)
	}
	var html string
	_ = chromedp.Run(ctx, chromedp.OuterHTML("body", &html, chromedp.ByQuery))
	t.Fatalf("timed out waiting for %s. Final DOM:\n%s", why, html)
}

// =============================================================================
// Test 1 — Issue #340: IsInitialMount fires on the initial HTTP GET.
// The server-rendered HTML (before the WS even connects) must reflect
// IsInitialMount()=true. After the WS connects and OnConnect runs, the DOM
// flips to im=NO because OnConnect's context is from the WS path.
// =============================================================================

func TestE2E_IsInitialMount_RendersInInitialHTML(t *testing.T) {
	baseURL, _, shutdown := startLifecycleE2EServer(t)
	defer shutdown()

	// 1. The server-rendered HTML — captured BEFORE any WS connects — must
	// reflect IsInitialMount()=true. http.Get takes no WS round-trip, so
	// this isolates the HTTP-path Mount classification.
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	html := string(body)

	if !strings.Contains(html, ">im=YES<") {
		t.Errorf("server-rendered HTML missing im=YES (IsInitialMount wiring broken on HTTP path).\nHTML:\n%s", html)
	}
	if !strings.Contains(html, ">rc=NO<") {
		t.Errorf("server-rendered HTML missing rc=NO (initial GET should not be a reconnect).\nHTML:\n%s", html)
	}
}

// =============================================================================
// Test 2 — Issue #340: IsReconnect fires after a real reconnect (page reload
// with persisted state). On the second visit the WS-path Mount/OnConnect see
// state restored from the prior visit, so IsReconnect()=true and the count
// (which OnConnect increments) accumulates across reloads.
// =============================================================================

func TestE2E_IsReconnect_ReflectsStateRestoration(t *testing.T) {
	_, port, shutdown := startLifecycleE2EServer(t)
	defer shutdown()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 45*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context
	installConsoleLogger(t, ctx)
	wsLog := e2etest.RecordWSFrames(ctx)

	chromeURL := e2etest.GetChromeTestURL(port)

	// First visit: HTTP GET + WS connect. State is persisted on Mount, then
	// restored when the WS connects, so the WS-path OnConnect sees IsReconnect
	// = true. The DOM settles with im=NO (WS overwrites HTTP's im=YES) and
	// rc=YES.
	if err := chromedp.Run(ctx, chromedp.Navigate(chromeURL)); err != nil {
		t.Fatalf("first navigate: %v", err)
	}
	waitForCondition(t, ctx,
		`document.getElementById('rc').innerText.includes('rc=YES')
		 && document.getElementById('im').innerText.includes('im=NO')`,
		5*time.Second,
		"first WS connect to settle (rc=YES, im=NO)")

	receivedBefore := wsLog.CountByDirection("received")

	// Reload: a fresh HTTP GET re-renders im=YES (initial mount) and rc=YES
	// (state restored from prior visit). The WS then reconnects and OnConnect
	// runs again, flipping the DOM back to im=NO rc=YES. We assert the WS
	// frame counter advanced (a direct, unambiguous "re-dialed" signal) rather
	// than reasoning about post-reload count values, which race with the
	// in-flight WS connect.
	if err := chromedp.Run(ctx, chromedp.Reload()); err != nil {
		t.Fatalf("reload: %v", err)
	}
	waitForCondition(t, ctx,
		`document.getElementById('rc').innerText.includes('rc=YES')
		 && document.getElementById('im').innerText.includes('im=NO')`,
		10*time.Second,
		"reload + WS reconnect to settle (rc=YES, im=NO)")

	if got := wsLog.CountByDirection("received") - receivedBefore; got == 0 {
		wsLog.PrintLast(5)
		t.Errorf("no new WS frames after reload — client may not have re-dialed")
	}
}

// =============================================================================
// Test 3 — Issue #345: ClearAllFlash from an action wipes every flash entry
// while leaving the rest of the DOM untouched.
// =============================================================================

func TestE2E_ClearAllFlash_RemovesAllFlashFromDOM(t *testing.T) {
	_, port, shutdown := startLifecycleE2EServer(t)
	defer shutdown()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 30*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context
	installConsoleLogger(t, ctx)
	wsLog := e2etest.RecordWSFrames(ctx)

	chromeURL := e2etest.GetChromeTestURL(port)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitVisible("#set-flashes", chromedp.ByID),
	); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	// 1. Click "Set flashes" — server sets two flashes. DOM updates over WS.
	if err := chromedp.Run(ctx, chromedp.Click("#set-flashes", chromedp.ByID)); err != nil {
		t.Fatalf("click set-flashes: %v", err)
	}
	waitForCondition(t, ctx,
		`document.getElementById('flash-success').innerText.includes('Saved!')
		 && document.getElementById('flash-info').innerText.includes('Heads up')`,
		5*time.Second,
		"both flashes to land in DOM after SetFlashes")

	// 2. Click "Clear all" — invokes ctx.ClearAllFlash. Both flashes vanish.
	if err := chromedp.Run(ctx, chromedp.Click("#clear-flashes", chromedp.ByID)); err != nil {
		t.Fatalf("click clear-flashes: %v", err)
	}
	waitForCondition(t, ctx,
		`document.getElementById('flash-success').innerText.trim() === ''
		 && document.getElementById('flash-info').innerText.trim() === ''`,
		5*time.Second,
		"both flashes to clear from DOM after ClearAllFlash")

	// Sanity: ClearAllFlash must NOT touch non-flash state. OnConnect's
	// `state.Count++` has run at least once by now, so count must be ≥ 1 —
	// asserting "count=0" would catch a regression where ClearAllFlash zeroed
	// numeric state instead of just flash keys.
	var count string
	if err := chromedp.Run(ctx, chromedp.Text("#count", &count, chromedp.NodeVisible, chromedp.ByID)); err != nil {
		t.Fatalf("read count: %v", err)
	}
	if !strings.Contains(count, "count=") || count == "count=0" {
		t.Errorf("count tag missing or zeroed after ClearAllFlash — non-flash state was wiped (got %q)", count)
	}

	// WS log: at least one frame round-tripped each way during the click sequence.
	// Sent: one frame per click (SetFlashes, ClearFlashes). Received: server's
	// state-update response for each action.
	if sent := wsLog.CountByDirection("sent"); sent < 2 {
		wsLog.PrintLast(10)
		t.Errorf("expected at least 2 WS frames sent (one per click), got %d", sent)
	}
	if received := wsLog.CountByDirection("received"); received < 2 {
		wsLog.PrintLast(10)
		t.Errorf("expected at least 2 WS frames received (one response per click), got %d", received)
	}
}

// =============================================================================
// Test 4 — Issue #339: regression check that wrapping the "no connected
// sessions" error with the ErrSessionDisconnected sentinel didn't break the
// WS happy path. The sentinel itself (and its errors.Is wiring) is asserted
// in livetemplate's TestLocalSession_TriggerActionDisconnectedReturnsError;
// this test only verifies that a fresh WS connect + action round-trip still
// works post-change.
// =============================================================================

func TestE2E_WS_HappyPathRegressionAfterSentinelWrap(t *testing.T) {
	baseURL, _, shutdown := startLifecycleE2EServer(t)
	defer shutdown()

	wsURL := "ws" + strings.TrimPrefix(baseURL, "http") + "/"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	if _, _, err := conn.ReadMessage(); err != nil {
		t.Fatalf("initial frame read: %v", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("set write deadline: %v", err)
	}
	msg, _ := json.Marshal(map[string]interface{}{"action": "SetFlashes", "data": map[string]interface{}{}})
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatalf("ws write: %v", err)
	}
	if _, _, err := conn.ReadMessage(); err != nil {
		t.Fatalf("ws read after action: %v", err)
	}
}
