//go:build browser

package e2e

import (
	"context"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TestAttributeRegistry_PublicRegistration covers the client's lvt-* attribute
// handler registry (livetemplate/livetemplate#473, Phase 1) in a real browser.
//
// What this test is FOR, given the rest of e2e already exercises the built-in
// attributes: the registry is a refactor whose failure mode is SILENCE. A
// handler that never runs, or one that runs but warns on every page load,
// breaks nothing an existing assertion looks at. So this test asserts the two
// things only a browser can show:
//
//  1. Loading a page that uses lvt-* attributes produces no registry
//     diagnostics — the registry must not warn about the framework's own
//     handlers.
//  2. A handler registered from a plain <script> AFTER the client tag takes
//     effect immediately, against the DOM that is already on screen. That is
//     the "late registration catch-up" path, and it is the ONLY path a
//     third-party bundle can use: under the documented `defer` load pattern the
//     core bundle has already auto-initialized and connected by the time any
//     second script evaluates, so there is no pre-init window to register in.
//
// It also pins the global's spelling. The browser build runs
// `--global-name=LiveTemplateClient` over a module that exports a class of the
// same name, so `window.LiveTemplateClient` is the module NAMESPACE and
// `LiveTemplateClient.registerAttribute(...)` — what every doc example writes —
// resolves to a module-level named export, not to the class static. Those two
// are indistinguishable in a unit test and trivially confused in a refactor.
func TestAttributeRegistry_PublicRegistration(t *testing.T) {
	t.Parallel()

	html := `<!DOCTYPE html>
<html>
<head><title>Attribute Registry</title></head>
<body>
	<div data-lvt-id="registry-test">
		<!-- Built-in attributes, present so the registry has real work to do
		     and so any warning about the framework's own handlers surfaces. -->
		<div id="fx-scroll" lvt-fx:scroll="into-view">scroll target</div>
		<div id="fx-highlight" lvt-fx:highlight="42">highlight target</div>
		<div id="scroll-away" lvt-scroll-away="up">scroll away</div>
		<div id="sentinel" lvt-scroll-sentinel></div>

		<!-- Claimed by a handler registered below, after the client loads. -->
		<button id="custom" lvt-x:e2e-copy="hello-from-template">Copy</button>
	</div>
	` + clientInitScript + `
</body>
</html>`

	chromeURL, cleanup := renderingTestServer(t, html)
	defer cleanup()

	ctx, _, cleanupChrome := GetPooledChrome(t)
	defer cleanupChrome()

	ctx, cancel := context.WithTimeout(ctx, getBrowserTimeout())
	defer cancel()

	// Console capture. Collected for the whole run so a failure below can be
	// read against what the page actually said, rather than guessed at.
	var (
		consoleMu  = make(chan struct{}, 1)
		consoleLog []string
	)
	consoleMu <- struct{}{}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		e, ok := ev.(*runtime.EventConsoleAPICalled)
		if !ok {
			return
		}
		var b strings.Builder
		b.WriteString(string(e.Type))
		b.WriteString(": ")
		for _, arg := range e.Args {
			if arg.Value != nil {
				b.Write(arg.Value)
				b.WriteByte(' ')
			}
		}
		<-consoleMu
		consoleLog = append(consoleLog, b.String())
		consoleMu <- struct{}{}
	})

	var (
		hasNamedExport bool
		hasClassStatic bool
		copiedValue    string
		addedCallCount int
		renderedHTML   string
	)

	err := chromedp.Run(ctx,
		chromedp.Navigate(chromeURL),
		chromedp.WaitReady("body"),
		waitForClient(),

		// The spelling every doc example uses must resolve.
		chromedp.Evaluate(`typeof window.LiveTemplateClient?.registerAttribute === "function"`, &hasNamedExport),
		chromedp.Evaluate(`typeof window.LiveTemplateClient?.LiveTemplateClient?.registerAttribute === "function"`, &hasClassStatic),

		// Register AFTER the client has loaded, connected and rendered — the
		// only sequence available to a second bundle. No render is triggered
		// afterwards on purpose: catch-up means the handler reaches the DOM
		// that is already on screen.
		chromedp.Evaluate(`(() => {
			window.__lvtAdded = 0;
			window.__lvtCopied = "";
			window.LiveTemplateClient.registerAttribute({
				attribute: "lvt-x:e2e-copy",
				onElementAdded(el, ctx) {
					window.__lvtAdded++;
					// ctx.value is a live accessor, not a captured string.
					el.addEventListener("click", () => { window.__lvtCopied = ctx.value; });
				},
			});
			return true;
		})()`, nil),

		chromedp.Click("#custom", chromedp.ByID),
		chromedp.Evaluate(`window.__lvtCopied`, &copiedValue),
		chromedp.Evaluate(`window.__lvtAdded`, &addedCallCount),
		chromedp.OuterHTML("html", &renderedHTML),
	)
	if err != nil {
		t.Fatalf("browser run failed: %v\nconsole:\n%s", err, strings.Join(consoleLog, "\n"))
	}

	if !hasNamedExport {
		t.Error("LiveTemplateClient.registerAttribute is not a function — the IIFE global is the module namespace, so registerAttribute must be a module-level named export")
	}
	if !hasClassStatic {
		t.Error("LiveTemplateClient.LiveTemplateClient.registerAttribute is not a function — the class static mirror is missing")
	}
	if addedCallCount != 1 {
		t.Errorf("onElementAdded fired %d times, want exactly 1 (late-registration catch-up should reach the live DOM once)", addedCallCount)
	}
	if copiedValue != "hello-from-template" {
		t.Errorf("ctx.value = %q, want %q — the handler did not read the attribute from the already-rendered element", copiedValue, "hello-from-template")
	}

	// Registry diagnostics must be silent on a page using only built-ins plus
	// one well-formed custom handler.
	//
	// Deliberately NOT "zero console output": this fixture serves no WebSocket
	// endpoint, so the client's connect() legitimately reports a transport
	// failure, and asserting on total silence would make the test a hostage to
	// unrelated noise. The assertion is scoped to what THIS feature emits.
	registryNoise := []string{
		"AttributeRegistry",
		"already claims this name",
		"never both",
		"has an empty value",
	}
	<-consoleMu
	defer func() { consoleMu <- struct{}{} }()
	for _, line := range consoleLog {
		for _, needle := range registryNoise {
			if strings.Contains(line, needle) {
				t.Errorf("registry emitted a diagnostic on a well-formed page: %s", line)
			}
		}
	}

	if !strings.Contains(renderedHTML, `lvt-x:e2e-copy`) {
		t.Error("custom attribute missing from the rendered DOM — fixture did not load as expected")
	}
	t.Logf("console (%d lines):\n%s", len(consoleLog), strings.Join(consoleLog, "\n"))
}
