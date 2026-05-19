//go:build browser

// V14 (Tier-2, user-visible leg) — broadcast-action-redesign #415, Phase 4.
//
// An ACL-denied ctx.Subscribe in the WS-connect Mount must, end-to-end in a
// real browser: (a) surface as an `lvt:error` CustomEvent { code, topic } on
// the wrapper, and (b) leave the WebSocket OPEN AND FUNCTIONAL (Phase 4's
// keep-open finalization of the decision Phase 1 deferred). The {code,topic}
// asserted here is byte-for-byte the server-emitted contract
// (livetemplate topic_runtime.go topicErrorEnvelope) and the client logic-leg
// jest contract (../client tests/topic-error-envelope.test.ts).
//
// CROSS-REPO / RELEASE ORDER: this file is built against the unreleased
// Phase-0..4 livetemplate via the committed `replace` in go.mod (pointing at
// the livetemplate Phase-4 worktree). Until livetemplate ships and Phase 5
// resolves that replace into a real version pin, the lvt Phase-4 branch is
// intentionally NOT independently mergeable / CI-runnable — this is the
// proposal's documented "lvt e2e gates last (consumes both)" release order,
// recorded in docs learnings/phase-4.md, NOT a skipped test. It is run green
// locally against the worktree build + locally-built client bundle. Phase 5/6
// swaps serveLocalPhase4ClientBundle back to e2etest.ServeClientLibrary once
// @livetemplate/client is published.
package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/livetemplate/livetemplate"
	e2etest "github.com/livetemplate/lvt/testing"
)

type v14ACLState struct{ Pong string }

type v14ACLController struct{}

// Mount subscribes a developer topic the ACL denies and propagates the error
// (the canonical `return s, err` real apps use), exercising the WS-connect
// Mount-failure path that emits the topic_forbidden envelope.
//
// IsInitialMount guard: Mount runs on every HTTP request AND WebSocket
// connect (livetemplate CLAUDE.md "Mount() guard pattern"). On the initial
// HTTP GET (page render), a denied ctx.Subscribe surfaces as HTTP 500
// (mount.go HTTP-GET path; phase-1.md "HTTP-GET ACL-denial → HTTP 500" — a
// pre-existing condition unchanged by Phase 4). If we didn't skip the
// Subscribe on the GET, the page would 500, client.js would never load, and
// there'd be no WS for V14 to exercise. V14's pinned scenario is exactly the
// WS-connect Mount denied-Subscribe (envelope + keep-open), so the guard
// scopes the denial to that path — which is also V14's exact spec wording.
func (c *v14ACLController) Mount(s v14ACLState, ctx *livetemplate.Context) (v14ACLState, error) {
	if ctx.IsInitialMount() {
		return s, nil
	}
	if err := ctx.Subscribe("private/admin"); err != nil {
		return s, err
	}
	return s, nil
}

// Ping is the post-envelope usability probe: a successful action round-trip
// re-renders #pong, proving the socket stayed FUNCTIONAL — not merely
// "not closed in the first N ms".
func (c *v14ACLController) Ping(s v14ACLState, _ *livetemplate.Context) (v14ACLState, error) {
	s.Pong = "PONG"
	return s, nil
}

// The <head> script installs a CAPTURE-phase document listener BEFORE
// client.js loads and before the WS connects. The client dispatches
// `lvt:error` on the wrapper as a non-bubbling CustomEvent; a capture-phase
// document listener still observes it (capture traverses window→target
// regardless of `bubbles`), which removes the listener-vs-first-frame race
// without instrumenting the client.
const v14ACLTemplate = `<!DOCTYPE html>
<html>
<head>
<title>V14 topic ACL error envelope</title>
<script>
window.__lvtErrors = [];
document.addEventListener('lvt:error', function (e) { window.__lvtErrors.push(e.detail); }, true);
</script>
</head>
<body>
  <div id="pong">{{.Pong}}</div>
  <button id="ping" lvt-on:click="Ping">Ping</button>
  <script src="/client.js"></script>
</body>
</html>`

// serveLocalPhase4ClientBundle serves the locally-built Phase-4 client bundle
// (NOT e2etest.ServeClientLibrary, which fetches the published client from a
// CDN and would not contain the unreleased lvt:error branch — and whose 1h
// disk cache can shadow LVT_CLIENT_CDN_URL anyway). Phase 5/6 reverts this to
// ServeClientLibrary after @livetemplate/client publishes.
//
// Bundle resolution order — both honor the cross-repo Phase-4 setup the
// proposal's release order assumes:
//   1. LVT_CLIENT_BUNDLE_PATH env var (explicit override)
//   2. Repo-relative convention path (assumes both worktrees live at
//      <repo>/.worktrees/broadcast-redesign-phase-4 — the standing rule)
//
// If neither resolves, t.Skip — this test is gated on the Phase-4 cross-repo
// setup. NOT a t.Fatalf: failing the whole -tags=browser run for a
// contributor not set up for V14 / Phase 4 would be a hostile default, and
// the proposal's release order explicitly says this e2e is intentionally
// not CI-runnable until Phase 5's pin bump. The skip message points at the
// learnings doc for setup.
func serveLocalPhase4ClientBundle(t *testing.T) http.HandlerFunc {
	t.Helper()
	const repoRelativeDefault = "../../../../client/.worktrees/broadcast-redesign-phase-4/dist/livetemplate-client.browser.js"
	path := os.Getenv("LVT_CLIENT_BUNDLE_PATH")
	if path == "" {
		path = repoRelativeDefault
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("V14 e2e gated on the Phase-4 client bundle (cross-repo, not CI-runnable "+
			"until livetemplate Phase 5's pin bump — see docs/proposals/broadcast-action-"+
			"redesign-proposal/learnings/phase-4.md). Set LVT_CLIENT_BUNDLE_PATH or build "+
			"the bundle at the repo-relative default %q: (cd <umbrella>/client/.worktrees/"+
			"broadcast-redesign-phase-4 && npm run build). read %q: %v",
			repoRelativeDefault, path, err)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write(data)
	}
}

func startV14ACLServer(t *testing.T) (port int, shutdown func()) {
	t.Helper()

	tmpl := livetemplate.Must(livetemplate.New("v14-acl-e2e",
		livetemplate.WithTopicACL(func(topic, _ string, _ *http.Request) (bool, error) {
			return topic != "private/admin", nil
		}),
	))
	if _, err := tmpl.Parse(v14ACLTemplate); err != nil {
		t.Fatalf("parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tmpl.Handle(&v14ACLController{}, livetemplate.AsState(&v14ACLState{})))
	mux.HandleFunc("/client.js", serveLocalPhase4ClientBundle(t))

	p, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("v14 server: %v", err)
		}
	}()

	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", p), 10*time.Second)
	return p, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
}

func TestE2E_V14_TopicACLDeniedEmitsLvtErrorAndKeepsWSOpen(t *testing.T) {
	// Artifact 2/4 — server logs. livetemplate logs via the global slog
	// default (incl. the new "Mount Subscribe denied … connection kept open"
	// WARN). Tee into a captured buffer; restore the default in Cleanup
	// (slog.SetDefault is process-global and the default persists across
	// tests in this package).
	serverLogger := e2etest.NewServerLogger()
	serverLogger.Start()
	// Teardown ordering matters: `defer serverLogger.Stop()` runs at function
	// return (LIFO with other defers), while `t.Cleanup` runs *after* all
	// defers. So slog.SetDefault is restored *after* the logger pipe closes —
	// no late writes from a shut-down logger into the global slog.
	defer serverLogger.Stop()
	prevSlog := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(os.Stderr, serverLogger.Writer()),
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))
	t.Cleanup(func() { slog.SetDefault(prevSlog) })

	port, shutdown := startV14ACLServer(t)
	defer shutdown()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 60*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context

	installConsoleLogger(t, ctx)         // artifact 1/4 — browser console (streamed live)
	wsLog := e2etest.RecordWSFrames(ctx) // artifact 3/4 — WS frames

	// dump surfaces all four artifacts on failure while the chrome ctx is
	// still alive (defer cleanup() cancels it before t.Cleanup runs, so HTML
	// must be captured here, not in Cleanup). Console is already streamed.
	dump := func() {
		t.Log("──────── SERVER LOGS ────────")
		for _, l := range serverLogger.GetLogs() {
			t.Log(l)
		}
		t.Log("──────── WS FRAMES (last 50) ────────")
		wsLog.PrintLast(50)
		var html string
		if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &html, chromedp.ByQuery)); err == nil {
			t.Logf("──────── RENDERED HTML ────────\n%s", html)
		} else {
			t.Logf("(rendered HTML capture failed: %v)", err)
		}
	}
	poll := func(jsExpr string, timeout time.Duration, why string) {
		t.Helper()
		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			var ok bool
			if err := chromedp.Run(ctx, chromedp.Evaluate(jsExpr, &ok)); err == nil && ok {
				return
			}
			if ctx.Err() != nil {
				dump()
				t.Fatalf("chrome ctx cancelled waiting for %s: %v", why, ctx.Err())
			}
			time.Sleep(100 * time.Millisecond)
		}
		dump()
		t.Fatalf("timed out waiting for %s", why)
	}

	chromeURL := e2etest.GetChromeTestURL(port)
	if err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitVisible("#ping", chromedp.ByID),
	); err != nil {
		dump()
		t.Fatalf("navigate: %v", err)
	}

	// (a) The denied Subscribe surfaced as an lvt:error CustomEvent on the
	// wrapper with the exact {code,topic} the server emitted.
	poll(
		`Array.isArray(window.__lvtErrors) && window.__lvtErrors.length === 1
		 && window.__lvtErrors[0].code === 'topic_forbidden'
		 && window.__lvtErrors[0].topic === 'private/admin'`,
		10*time.Second,
		"lvt:error CustomEvent {code:topic_forbidden, topic:private/admin}",
	)

	// (b) Verify the server actually took the Phase-4 keep-open code path —
	// not just that the client happened to stay connected. The server-side
	// WARN log is the load-bearing signal that mount.go's *TopicForbiddenError
	// branch (Option B fall-through) ran instead of the pre-Phase-4 return.
	// Guards against silent regressions where the log message changes or the
	// WARN is removed without the corresponding behavior change.
	if !serverLogger.HasLog("connection kept open") {
		dump()
		t.Fatal("expected server to log the Phase-4 keep-open WARN " +
			"('Mount Subscribe denied by topic ACL; surfaced to client, " +
			"connection kept open') after ACL denial — the WARN is the proof " +
			"the server took the keep-open code path, not just that the " +
			"client happened to stay connected")
	}

	// (c) The WS is OPEN AND FUNCTIONAL: a real action round-trips and
	// re-renders. If the server had closed the socket (pre-Phase-4 behavior),
	// the client's auto-reconnect would re-Mount → re-deny → loop, and #pong
	// would never become PONG over a live WS.
	if err := chromedp.Run(ctx, chromedp.Click("#ping", chromedp.ByID)); err != nil {
		dump()
		t.Fatalf("click #ping: %v", err)
	}
	poll(
		`document.getElementById('pong') && document.getElementById('pong').textContent === 'PONG'`,
		10*time.Second,
		"#pong === 'PONG' (WS action round-trip after the error envelope)",
	)

	// Sanity: the envelope did NOT also drive a spurious second lvt:error or
	// leak into the diff path.
	var errCount int
	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.__lvtErrors.length`, &errCount)); err != nil {
		dump()
		t.Fatalf("read __lvtErrors.length: %v", err)
	}
	if errCount != 1 {
		var errsJSON string
		_ = chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify(window.__lvtErrors)`, &errsJSON))
		dump()
		t.Fatalf("expected exactly one lvt:error, got %d; window.__lvtErrors = %s", errCount, errsJSON)
	}
}
