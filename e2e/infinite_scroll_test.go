//go:build browser

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestInfiniteScrollLoadsMore checks a generated app's infinite pagination in a
// browser: scrolling the [lvt-scroll-sentinel] into view sends load_more and
// renders the next page. The client's observer does this; the kit carries no
// scroll script of its own.
func TestInfiniteScrollLoadsMore(t *testing.T) {
	var rows []struct {
		SQL  string
		Args []interface{}
	}
	for i := 1; i <= 25; i++ {
		rows = append(rows, struct {
			SQL  string
			Args []interface{}
		}{
			SQL:  "INSERT INTO items (id, name, created_at) VALUES (?, ?, datetime('now', ?))",
			Args: []interface{}{fmt.Sprintf("item-%02d", i), fmt.Sprintf("Item %02d", i), fmt.Sprintf("-%d minutes", i)},
		})
	}
	appDir, serverPort := startItemsApp(t, "scrollapp", rows, "--page-size", "10")
	rec := newRecordedTab(t, appDir)
	bctx := rec.ctx

	const rowCount = `document.querySelectorAll('tbody tr[data-key]').length`
	if err := chromedp.Run(bctx,
		// Short enough that the first page's sentinel starts below the fold.
		chromedp.EmulateViewport(1200, 400),
		chromedp.Navigate(getTestURL(serverPort)+"/items"),
		waitForWebSocketReady(5*time.Second),
		waitFor(rowCount+` === 10`, 5*time.Second),
	); err != nil {
		t.Fatalf("first page didn't show 10 rows: %v%s", err, rec.dump())
	}
	// The client observes the sentinel with a 200px rootMargin, so it must
	// start further below the fold than that for the scroll to be the trigger.
	var belowFold bool
	if err := chromedp.Run(bctx, chromedp.Evaluate(
		`document.querySelector('[lvt-scroll-sentinel]').getBoundingClientRect().top > window.innerHeight + 200`, &belowFold),
	); err != nil || !belowFold {
		t.Fatalf("the sentinel starts within the observer's margin, so the test can't tell a scroll from a page load (err=%v)%s", err, rec.dump())
	}

	if err := chromedp.Run(bctx,
		chromedp.Evaluate(`document.querySelector('[lvt-scroll-sentinel]').scrollIntoView()`, nil),
		waitFor(rowCount+` === 20`, 10*time.Second),
	); err != nil {
		t.Fatalf("scrolling to the sentinel didn't load the second page: %v%s", err, rec.dump())
	}

	// The last page holds 5 rows, and with nothing left the sentinel goes.
	if err := chromedp.Run(bctx,
		chromedp.Evaluate(`document.querySelector('[lvt-scroll-sentinel]').scrollIntoView()`, nil),
		waitFor(rowCount+` === 25 && !document.querySelector('[lvt-scroll-sentinel]')`, 10*time.Second),
	); err != nil {
		t.Fatalf("the last page didn't load, or the sentinel stayed after it: %v%s", err, rec.dump())
	}

	if !rec.sentAction("load_more") {
		t.Errorf("no load_more action went over the socket%s", rec.dump())
	}
}
