//go:build browser

package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	_ "github.com/mattn/go-sqlite3"
)

// E2E test timing constants
// These constants define wait times for various browser operations to make tests
// more maintainable and easier to tune for different environments
const (
	// shortDelay is used for brief pauses between operations (e.g., after clicking buttons)
	shortDelay = 500 * time.Millisecond
)

// waitForCondition polls a JavaScript condition until it returns true or times out
// This is more reliable than manual retry loops with fixed delays
func waitForCondition(ctx context.Context, jsCondition string, timeout time.Duration, pollInterval time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// Check if parent context already has a deadline
		deadline, hasDeadline := ctx.Deadline()

		// Calculate effective timeout (minimum of requested and remaining parent timeout)
		effectiveTimeout := timeout
		if hasDeadline {
			remaining := time.Until(deadline)
			if remaining < timeout {
				effectiveTimeout = remaining
			}
		}

		// Only create new timeout context if we have time remaining
		if effectiveTimeout <= 0 {
			return fmt.Errorf("parent context already expired while waiting for condition: %s", jsCondition)
		}

		ctx, cancel := context.WithTimeout(ctx, effectiveTimeout)
		defer cancel()

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Provide better error message with actual timeout used
				return fmt.Errorf("timeout (%.1fs) waiting for condition: %s (error: %v)", effectiveTimeout.Seconds(), jsCondition, ctx.Err())
			case <-ticker.C:
				var result bool
				if err := chromedp.Evaluate(jsCondition, &result).Do(ctx); err != nil {
					// Continue polling even on evaluation errors (DOM might not be ready)
					continue
				}
				if result {
					return nil
				}
			}
		}
	}
}

// PollUntil polls a JS expression against the chromedp ctx until it returns
// true or the timeout elapses. On failure (timeout or chromedp context
// cancellation), it invokes the optional onTimeout callback for diagnostic
// artifact dumping (server logs, WS frames, rendered HTML, etc.) and then
// calls t.Fatalf with a descriptive message.
//
// This is the canonical "dump-aware poll" helper for browser e2e tests. Pass
// nil for onTimeout if no extra diagnostic dumping is needed beyond the body
// OuterHTML printed in the Fatalf message.
//
// A cancelled chromedp context (browser crash, Docker shutdown, parent test
// timeout) is fatal immediately — without that early return the loop would
// silently treat context.Canceled as "condition not yet true" and report a
// misleading timeout.
//
// 100ms poll interval matches lvt/testing's WaitFor convention — chosen there
// for stability and to avoid CPU thrashing on slow CI runners.
func PollUntil(t *testing.T, ctx context.Context, jsExpr string, timeout time.Duration, why string, onTimeout func()) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var ok bool
		err := chromedp.Run(ctx, chromedp.Evaluate(jsExpr, &ok))
		if err == nil && ok {
			return
		}
		// Context cancellation is fatal immediately — never wait for the
		// timeout in that case.
		if ctx.Err() != nil {
			if onTimeout != nil {
				onTimeout()
			}
			t.Fatalf("chromedp context cancelled waiting for %s: %v", why, ctx.Err())
		}
		// Non-cancellation chromedp errors (e.g. detached frame, JS syntax
		// error, evaluation panic) are surfaced as warnings rather than
		// swallowed — without this a real error gets misreported as a
		// generic timeout. We still keep polling: transient errors (DOM not
		// ready yet, navigation in progress) are expected.
		if err != nil {
			t.Logf("PollUntil[%s]: chromedp error (will retry): %v", why, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if onTimeout != nil {
		onTimeout()
	}
	var html string
	_ = chromedp.Run(ctx, chromedp.OuterHTML("body", &html, chromedp.ByQuery))
	t.Fatalf("timed out waiting for %s. Final DOM body:\n%s", why, html)
}

// seedTestData seeds test data into SQLite database using parameterized queries
// This is safer than string concatenation and prevents SQL injection
func seedTestData(dbPath string, queries []struct {
	SQL  string
	Args []interface{}
}) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Begin transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, q := range queries {
		if _, err := tx.Exec(q.SQL, q.Args...); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
