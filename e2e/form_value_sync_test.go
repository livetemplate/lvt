//go:build browser

package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	cdpnetwork "github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// TestFormValueSync checks that the sort <select> and the search input show
// the server's value after each patch, and keep the user's value while they
// are focused. The kits shipped a data-expected-value script to force the
// select back in line; this test is what decides whether the client still
// needs it.
func TestFormValueSync(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := createTestApp(t, tmpDir, "syncapp", &AppOptions{Kit: "multi"})
	setupLocalClientLibrary(t, appDir)

	if err := runLvtCommand(t, appDir, "gen", "resource", "items", "name:string"); err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}
	if err := runLvtCommand(t, appDir, "migration", "up"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Newest first gives Cherry, Banana, Apple; name A-Z gives Apple first.
	type row = struct {
		SQL  string
		Args []interface{}
	}
	var rows []row
	for i, name := range []string{"Apple", "Banana", "Cherry"} {
		rows = append(rows, row{
			SQL:  "INSERT INTO items (id, name, created_at) VALUES (?, ?, datetime('now', ?))",
			Args: []interface{}{fmt.Sprintf("item-%d", i), name, fmt.Sprintf("-%d minutes", 10-i)},
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
	rec := recordBrowser(bctx, appDir)

	if err := chromedp.Run(bctx,
		cdpnetwork.Enable(),
		chromedp.Navigate(getTestURL(serverPort)+"/items"),
		waitForWebSocketReady(5*time.Second),
		waitFor(firstRowIs("Cherry"), 5*time.Second),
	); err != nil {
		t.Fatalf("list didn't load newest first: %v%s", err, rec.dump())
	}

	const sel = `document.querySelector('select[name="sort_by"]')`
	const search = `document.querySelector('input[name="query"]')`
	selectValue := func(want string) string { return sel + `.value === ` + strconv.Quote(want) }
	searchValue := func(want string) string { return search + `.value === ` + strconv.Quote(want) }
	step := func(label string, actions ...chromedp.Action) {
		t.Helper()
		if err := chromedp.Run(bctx, actions...); err != nil {
			t.Fatalf("%s: %v%s", label, err, rec.dump())
		}
	}

	// 1. Pick a sort the way a user does. The server re-sorts and the select
	//    shows the choice.
	step("sort by name A-Z",
		chromedp.Evaluate(sel+`.value = "name_asc"; `+sel+`.dispatchEvent(new Event("change", {bubbles: true}))`, nil),
		waitFor(firstRowIs("Apple"), 5*time.Second),
		waitFor(selectValue("name_asc"), 2*time.Second),
	)

	// 2. A search patches the page; the select keeps the server's sort and
	//    the input keeps what was typed.
	step("search keeps the sort",
		chromedp.Focus(`input[name="query"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="query"]`, "an", chromedp.ByQuery),
		waitFor(firstRowIs("Banana"), 5*time.Second),
		waitFor(selectValue("name_asc")+` && `+searchValue("an"), 2*time.Second),
	)

	// 3. Clearing the search patches again; the input is empty and the
	//    select still shows the sort.
	step("clearing the search keeps the sort",
		chromedp.Evaluate(`document.querySelector('button[name="search"]').click()`, nil),
		waitFor(firstRowIs("Apple"), 5*time.Second),
		waitFor(selectValue("name_asc")+` && `+searchValue(""), 2*time.Second),
	)

	// 4. A patch that arrives while the select is focused, with a choice the
	//    user hasn't committed, leaves that choice alone.
	step("a patch keeps a focused select's choice",
		chromedp.Focus(`select[name="sort_by"]`, chromedp.ByQuery),
		chromedp.Evaluate(sel+`.value = "oldest_first"`, nil),
		chromedp.Evaluate(`window.liveTemplateClient.send({action: "search", data: {query: "e"}})`, nil),
		waitFor(rowCountIs(2), 5*time.Second),
		chromedp.Sleep(500*time.Millisecond),
		waitFor(selectValue("oldest_first"), 2*time.Second),
	)

	// 5. A patch that arrives while the user is typing in the search input
	//    leaves the typed text alone.
	step("a patch keeps text being typed",
		chromedp.Evaluate(sel+`.value = "name_asc"`, nil),
		chromedp.Focus(`input[name="query"]`, chromedp.ByQuery),
		chromedp.Evaluate(search+`.value = "Che"`, nil),
		chromedp.Evaluate(`window.liveTemplateClient.send({action: "sort", data: {sort_by: "name_desc"}})`, nil),
		waitFor(firstRowIs("Cherry"), 5*time.Second),
		chromedp.Sleep(500*time.Millisecond),
		waitFor(searchValue("Che"), 2*time.Second),
	)

	// 6. The server changed the sort to name_desc in step 5. Once the select
	//    loses focus, a patch shows the server's value.
	step("an unfocused select shows the server's sort",
		chromedp.Evaluate(`document.activeElement.blur()`, nil),
		chromedp.Evaluate(`window.liveTemplateClient.send({action: "search", data: {query: ""}})`, nil),
		waitFor(rowCountIs(3), 5*time.Second),
		waitFor(selectValue("name_desc"), 3*time.Second),
	)
}

func firstRowIs(name string) string {
	return `(() => { const r = document.querySelector('tbody tr[data-key]'); return !!r && r.textContent.includes(` + strconv.Quote(name) + `); })()`
}

func rowCountIs(n int) string {
	return `document.querySelectorAll('tbody tr[data-key]').length === ` + strconv.Itoa(n)
}
