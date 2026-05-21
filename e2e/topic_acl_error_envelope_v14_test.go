//go:build browser

// V14 (Tier-2, user-visible leg) — broadcast-action-redesign #415, Phase 5.
//
// An ACL-denied ctx.Subscribe in the WS-connect Mount must, end-to-end in a
// real browser: (a) surface as an `lvt:error` CustomEvent { code, topic } on
// the wrapper, and (b) leave the WebSocket OPEN AND FUNCTIONAL (the keep-open
// finalization shipped in livetemplate Phase 4). The {code,topic} asserted
// here is byte-for-byte the server-emitted contract
// (livetemplate topic_runtime.go topicErrorEnvelope) and the client logic-leg
// jest contract (../client tests/topic-error-envelope.test.ts).
//
// Phase 5 cleanup vs the original Phase-4 draft (closed lvt#327, superseded):
//   - go.mod now pins livetemplate v0.10.0 directly (no `replace`).
//   - The client bundle is served via the canonical e2etest.ServeClientLibrary
//     (no `serveLocalPhase4ClientBundle`) now that @livetemplate/client@0.9.2
//     is on npm with the `lvt:error` branch and there's no cross-repo gating.
package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/livetemplate/livetemplate"
	"github.com/livetemplate/lvt/e2e"
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
// pre-existing condition unchanged by Phase 4/5). If we didn't skip the
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
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	p, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: mux}
	go func() {
		// slog (not log.Printf) so a startup error lands in the
		// serverLogger-captured stream the test surfaces on failure,
		// not in the default process stderr.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("v14 server ListenAndServe failed", slog.Any("error", err))
		}
	}()

	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", p), 10*time.Second)
	return p, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Log shutdown failures so flaky teardown on slow CI is diagnosable
		// (currently captured via serverLogger / surfaced on test failure).
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("v14 server Shutdown failed", slog.Any("error", err))
		}
	}
}

func TestE2E_V14_TopicACLDeniedEmitsLvtErrorAndKeepsWSOpen(t *testing.T) {
	// Not parallel: this test mutates the process-global slog.Default. Other
	// tests in this package that also touch slog.Default (or that observe
	// livetemplate's internal logging via the same global) would race if
	// this ran via t.Parallel(). Sibling browser e2es (e.g.
	// lifecycle_ergonomics_test.go) follow the same non-parallel convention.

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
	e2e.PollUntil(t, ctx,
		`Array.isArray(window.__lvtErrors) && window.__lvtErrors.length === 1
		 && window.__lvtErrors[0].code === 'topic_forbidden'
		 && window.__lvtErrors[0].topic === 'private/admin'`,
		10*time.Second,
		"lvt:error CustomEvent {code:topic_forbidden, topic:private/admin}",
		dump,
	)

	// (b) Verify the server actually took the Phase-4 keep-open code path —
	// not just that the client happened to stay connected. The server-side
	// WARN log is the load-bearing signal that mount.go's *TopicForbiddenError
	// branch (Option B fall-through) ran instead of the pre-Phase-4 return.
	// Guards against silent regressions where the WARN is removed without the
	// corresponding behavior change.
	//
	// Coupling: the structured slog attribute `event=topic_acl_denied_keep_open`
	// emitted by livetemplate's mount.go (grep anchor: `slog.String("event",
	// "topic_acl_denied_keep_open")` near the `*TopicForbiddenError` branch).
	// Structured-key assertion is the v0.10.1+ hardening of the previous
	// substring-of-message coupling — robust against prose rewordings of the
	// WARN message itself.
	if !serverLogger.HasLog("event=topic_acl_denied_keep_open") {
		dump()
		t.Fatal("expected server to log the keep-open WARN with structured " +
			"attribute event=topic_acl_denied_keep_open after ACL denial — " +
			"the attribute is the proof the server took the keep-open code " +
			"path, not just that the client happened to stay connected")
	}

	// (c) The WS is OPEN AND FUNCTIONAL: a real action round-trips and
	// re-renders. If the server had closed the socket (pre-Phase-4 behavior),
	// the client's auto-reconnect would re-Mount → re-deny → loop, and #pong
	// would never become PONG over a live WS.
	if err := chromedp.Run(ctx, chromedp.Click("#ping", chromedp.ByID)); err != nil {
		dump()
		t.Fatalf("click #ping: %v", err)
	}
	e2e.PollUntil(t, ctx,
		`document.getElementById('pong') && document.getElementById('pong').textContent === 'PONG'`,
		10*time.Second,
		"#pong === 'PONG' (WS action round-trip after the error envelope)",
		dump,
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
