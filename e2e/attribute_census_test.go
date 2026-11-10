//go:build browser

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/livetemplate/livetemplate"
	"github.com/livetemplate/lvt/e2e"
	e2etest "github.com/livetemplate/lvt/testing"
)

// TestE2E_AttributeLifecycle covers Phase 3 of livetemplate/livetemplate#473 in
// a real browser against a real livetemplate server, because three of its four
// claims are about load order and only a browser loading a page orders scripts:
//
//  1. The window.lvtAttributes queue. The handlers are pushed from an inline
//     <script> in <head>, BEFORE the deferred client tag. An inline script
//     ignores `defer`, so it runs first, when no LiveTemplateClient global
//     exists — calling registerAttribute there would throw.
//  2. The connect pass. Those handlers register before any client exists, and
//     the page load only hydrates the server-rendered markup, so no render
//     changes anything. They must still arm the SSR'd element, exactly once,
//     and stay at once across the renders an action causes.
//  3. The attribute census. The server sends meta.attributes in the authored
//     spelling (lvt-el:addClass); the client must claim every core name in it,
//     case-folded, and warn exactly once for the one unhandled name.
//  4. delegatedEvents. A custom element's event reaches a Controller action
//     through lvt-on:my-widget-change once a handler declares the type.
//
// Set LVT_LOCAL_CLIENT_JS to a built client bundle to run this against an
// unreleased client (see e2etest.GetClientLibraryJS); otherwise it uses the
// published one.
func TestE2E_AttributeLifecycle(t *testing.T) {
	serverLogger := e2etest.NewServerLogger()
	serverLogger.Start()
	defer serverLogger.Stop()
	// slog.SetDefault is process-global, so this test must not run in parallel.
	prevSlog := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(os.Stderr, serverLogger.Writer()),
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))
	t.Cleanup(func() { slog.SetDefault(prevSlog) })

	port, shutdown := startAttributeLifecycleServer(t)
	defer shutdown()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 60*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context

	consoleLog := e2etest.NewConsoleLogger()
	consoleLog.Start(ctx)
	// The client logs its scope prefix as the first argument, so a census line
	// is only recognisable with every argument joined.
	consoleLines := func() []string {
		var lines []string
		for _, l := range consoleLog.GetLogs() {
			lines = append(lines, l.Type+": "+strings.Join(l.Args, " "))
		}
		return lines
	}
	wsLog := e2etest.RecordWSFrames(ctx)

	dump := func() {
		t.Logf("──────── CONSOLE ────────\n%s", strings.Join(consoleLines(), "\n"))
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
		chromedp.WaitVisible("#count", chromedp.ByID),
	); err != nil {
		dump()
		t.Fatalf("navigate: %v", err)
	}

	// (1)+(2): the queued declarative handler armed the SSR'd element once.
	// The count is checked first so && short-circuits before __ssrValue, which
	// only exists once the handler has run.
	e2e.PollUntil(t, ctx, `window.__ssrAdded === 1 && window.__ssrValue() === "armed-0"`,
		10*time.Second, "queued handler never armed the server-rendered element", dump)

	// (3): the census runs after window load, so wait for its warning itself
	// rather than for load, which only schedules it. The count is re-checked
	// once more below, after the renders the action causes.
	// The matched strings are the client's own wording: attribute-census.ts
	// for the census, attribute-registry.ts for the registry diagnostics.
	const censusNeedle = "no attribute handler claims it"
	censusLines := func() []string {
		var out []string
		for _, line := range consoleLines() {
			if strings.Contains(line, censusNeedle) {
				out = append(out, line)
			}
		}
		return out
	}
	for deadline := time.Now().Add(10 * time.Second); len(censusLines()) == 0; {
		if time.Now().After(deadline) {
			dump()
			t.Fatal("the census never warned about lvt-x:nope")
		}
		time.Sleep(50 * time.Millisecond)
	}

	// (4): the custom element's event dispatches Bump through lvt-on:.
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.getElementById("widget").dispatchEvent(new CustomEvent("my-widget-change", { bubbles: true }))`, nil),
	); err != nil {
		dump()
		t.Fatalf("dispatch custom event: %v", err)
	}
	e2e.PollUntil(t, ctx, `document.getElementById("count").textContent.trim() === "1"`,
		10*time.Second, "lvt-on:my-widget-change never reached the Bump action", dump)

	// A render later the handler has still fired once, and its live ctx.value
	// re-reads the re-rendered attribute rather than a captured one.
	e2e.PollUntil(t, ctx, `window.__ssrAdded === 1 && window.__ssrValue() === "armed-1"`,
		10*time.Second, "handler re-fired on re-render, or ctx.value was captured", dump)

	// Give a duplicate warning time to show up after the action's renders
	// before counting: "exactly once" is the claim under test.
	time.Sleep(500 * time.Millisecond)
	census := censusLines()
	var registry []string
	for _, line := range consoleLines() {
		for _, needle := range []string{"AttributeRegistry", "already claims this name", "can never match"} {
			if strings.Contains(line, needle) {
				registry = append(registry, line)
			}
		}
	}
	if len(census) != 1 || !strings.Contains(census[0], "lvt-x:nope") {
		dump()
		t.Errorf("want exactly one census warning, for lvt-x:nope; got %d:\n%s", len(census), strings.Join(census, "\n"))
	}
	if len(registry) != 0 {
		dump()
		t.Errorf("registry emitted diagnostics on a well-formed page:\n%s", strings.Join(registry, "\n"))
	}

	// The wire carried the census in the authored spelling — what makes (3)'s
	// case folding a real test rather than a lowercase-only one.
	var sawAuthored bool
	for _, m := range wsLog.GetMessages() {
		var frame struct {
			Meta struct {
				Attributes []string `json:"attributes"`
			} `json:"meta"`
		}
		if m.Direction != "received" || json.Unmarshal([]byte(m.Data), &frame) != nil {
			continue
		}
		if slices.Contains(frame.Meta.Attributes, "lvt-el:addClass") {
			sawAuthored = true
		}
	}
	if !sawAuthored {
		dump()
		t.Error(`no WS frame carried meta.attributes with "lvt-el:addClass" — the census never reached the client`)
	}
}

type attrLifecycleState struct{ Count int }

type attrLifecycleController struct{}

func (c *attrLifecycleController) Bump(s attrLifecycleState, _ *livetemplate.Context) (attrLifecycleState, error) {
	s.Count++
	return s, nil
}

// Every core namespace the census must claim appears on the page, with one
// mixed-case lvt-el: name, plus one name nothing handles (lvt-x:nope).
const attrLifecycleTemplate = `<!DOCTYPE html>
<html>
<head>
<title>Attribute lifecycle</title>
<script>
  window.__ssrAdded = 0;
  (window.lvtAttributes ||= []).push({
    attribute: "lvt-x:e2e-ssr",
    onElementAdded(el, ctx) {
      // Recorded on window, not on the element: morphdom strips attributes
      // the server did not render, so a data-* marker would not survive.
      window.__ssrAdded++;
      window.__ssrValue = () => ctx.value;
    },
  });
  (window.lvtAttributes ||= []).push({
    attribute: "lvt-x:widget",
    delegatedEvents: ["my-widget-change"],
    onElementAdded() {},
  });
</script>
</head>
<body>
  <div id="count">{{.Count}}</div>
  <button id="ssr" lvt-x:e2e-ssr="armed-{{.Count}}">ssr</button>
  <my-widget id="widget" lvt-x:widget="on" lvt-on:my-widget-change="Bump"></my-widget>
  <button id="bump" lvt-on:click="Bump" lvt-mod:debounce="10" lvt-el:addClass:on:click="pressed">bump</button>
  <input id="k" lvt-on:keydown="Bump" lvt-key="Enter">
  <div lvt-ignore>ignored</div>
  <div lvt-fx:highlight="{{.Count}}">hl</div>
  <span lvt-x:nope="1">nobody handles this</span>
  <script src="/client.js" defer></script>
</body>
</html>`

func startAttributeLifecycleServer(t *testing.T) (int, func()) {
	t.Helper()

	tmpl := livetemplate.Must(livetemplate.New("attr-lifecycle-e2e"))
	if _, err := tmpl.Parse(attrLifecycleTemplate); err != nil {
		t.Fatalf("parse template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tmpl.Handle(&attrLifecycleController{}, livetemplate.AsState(&attrLifecycleState{})))
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	p, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("attribute lifecycle server failed", slog.Any("error", err))
		}
	}()
	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", p), 10*time.Second)
	return p, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("attribute lifecycle server shutdown failed", slog.Any("error", err))
		}
	}
}
