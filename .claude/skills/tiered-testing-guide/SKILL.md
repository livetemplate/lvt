---
name: tiered-testing-guide
description: AUTOMATICALLY USE when writing tests, adding test files, creating e2e tests, or implementing test coverage in the lvt codebase. Guides test tier selection (HTTP vs browser vs Jest). Activates for any testing work in this project.
---

# Tiered Testing Guide for lvt Development

**This skill activates AUTOMATICALLY** when Claude detects testing-related work in the lvt codebase.

## AUTOMATIC ACTIVATION

Claude MUST use this skill when:
- Writing any new test file
- Adding test functions
- Implementing test coverage
- Creating e2e tests
- Discussing test strategy
- Fixing failing tests
- Reviewing test code

**NO user prompt required** - skill activates based on context.

### Detection Signals
- File paths containing `_test.go`
- Mentions of: test, testing, e2e, coverage, assertion
- Creating files in `e2e/` or `*_test.go`
- Build tags: http, browser, deployment
- Import of `testing` package

---

## Brainstorming: Which Test Type?

When adding a new test, ask these questions:

### Question 1: Does it require JavaScript execution?

**NO** → Use **HTTP test** (Tier 1)
- Form submissions, CRUD, page loads, redirects, template rendering

**YES** → Continue to Question 2

### Question 2: What JavaScript behavior?

| Behavior | Test Type |
|----------|-----------|
| Client-side form validation logic | **Jest** (Tier 0) |
| DOM manipulation (add/remove elements) | **Browser** (Tier 2) |
| Focus management | **Browser** (Tier 2) |
| Modal/dialog lifecycle | **Browser** (Tier 2) |
| WebSocket connection handling | **Browser** (Tier 2) |
| CSS animations/transitions | **Browser** (Tier 2) |
| Event bubbling/delegation | **Browser** (Tier 2) |

### Question 3: Is it infrastructure/deployment?

**YES** → Use `//go:build deployment` (Tier 3)
- Docker builds, Fly.io deployment, K8s configs

### Decision Flowchart

```
┌─────────────────────────────────────┐
│ What are you testing?               │
└─────────────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────┐
│ Does it need JavaScript execution?  │
└─────────────────────────────────────┘
        │                    │
       NO                   YES
        │                    │
        ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│ HTTP Test       │  │ Is it pure JS   │
│ (Tier 1)        │  │ logic? No DOM?  │
│ //go:build http │  └─────────────────┘
└─────────────────┘          │
                        YES  │  NO
                        │    │
                        ▼    ▼
              ┌──────────┐ ┌────────────────┐
              │ Jest     │ │ Browser Test   │
              │ (Tier 0) │ │ (Tier 2)       │
              └──────────┘ │ //go:build     │
                           │ browser        │
                           └────────────────┘
```

### Example Brainstorming

**"I need to test that clicking 'Add Item' adds a row to the table"**

1. Does it need JS? → YES (click handler, DOM update)
2. Is it pure JS logic? → NO (requires DOM)
3. **→ Browser test (Tier 2)**

**"I need to test that submitting the form creates a database record"**

1. Does it need JS? → NO (server-side processing)
2. **→ HTTP test (Tier 1)** - POST form, check database

**"I need to test that invalid email shows error message"**

1. Does it need JS? → DEPENDS
   - If server-side validation: **HTTP test** (check error in response HTML)
   - If client-side validation: **Jest** (test validation function) or **Browser** (test UI)

---

## Quick Decision Guide

| What You're Testing | Tier | Build Tag | Why |
|---------------------|------|-----------|-----|
| Form submission works | **Tier 1 (HTTP)** | `//go:build http` | No JS needed |
| Database CRUD | **Tier 1 (HTTP)** | `//go:build http` | Server logic only |
| Template renders correctly | **Tier 1 (HTTP)** | `//go:build http` | HTML output check |
| CSRF protection | **Tier 1 (HTTP)** | `//go:build http` | Token extraction |
| Modal opens/closes | **Tier 2 (Browser)** | `//go:build browser` | Requires DOM |
| Focus preservation | **Tier 2 (Browser)** | `//go:build browser` | Browser tracks focus |
| WebSocket reconnection | **Tier 2 (Browser)** | `//go:build browser` | Real WS API |
| Animation plays | **Tier 2 (Browser)** | `//go:build browser` | CSS in browser |
| Docker deployment | **Tier 3** | `//go:build deployment` | Infrastructure |
| Client-side validation JS | **Tier 0** | Jest | Pure JS logic |

---

## The 4 Tiers

### Tier 0: Client JS (Jest + jsdom)
- **Speed**: ~5s
- **When**: Pure JavaScript logic, no server needed
- **Run**: `npm test` (separate from Go tests)

### Tier 1: HTTP Tests (Go httptest)
- **Speed**: ~10s
- **When**: Server rendering, CRUD, forms - NO browser needed
- **Build Tag**: `//go:build http`
- **Run**: `make test-http` or `go test -tags http ./...`

### Tier 2: Browser Tests (chromedp)
- **Speed**: ~45s
- **When**: ONLY for JS execution, DOM manipulation, focus, animations
- **Build Tag**: `//go:build browser`
- **Run**: `make test-browser` (weekly/pre-release)
- **Limit**: Only 12 focused test types needed

### Tier 3: Application/Deployment Tests
- **When**: Full workflows, deployment infrastructure
- **Build Tag**: `//go:build deployment` (for Docker/Fly tests)
- **Run**: `make test-all`

---

## CRITICAL: When to Use Browser Tests

**Browser tests are EXPENSIVE** (CPU, memory, time). Only use for:

1. DOM list operations (add/remove/reorder items)
2. Table rendering
3. Form client-side validation (JS-based)
4. Modal lifecycle (open/close/reopen)
5. Pagination navigation
6. Infinite scroll
7. Focus preservation after actions
8. Scroll directives
9. Lifecycle hooks (JS events)
10. Event delegation
11. WebSocket reconnection
12. Conditional rendering (JS-driven)

**Everything else → Use HTTP tests!**

---

## HTTP Test Pattern (Tier 1)

```go
//go:build http

package mypackage

import (
    "net/url"
    "testing"
    lvttest "github.com/livetemplate/lvt/testing"
)

func TestFeatureHTTP(t *testing.T) {
    test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
        AppPath: "../../cmd/myapp/main.go",
    })
    defer test.Cleanup()

    t.Run("Page Loads", func(t *testing.T) {
        resp := test.Get("/mypage")
        assert := lvttest.NewHTTPAssert(resp)

        assert.StatusOK(t)
        assert.Contains(t, "Expected Content")
        assert.NoTemplateErrors(t)
        assert.HasFormField(t, "fieldname")
    })

    t.Run("Form Submit", func(t *testing.T) {
        formData := url.Values{}
        formData.Set("name", "Test Value")

        resp := test.PostForm("/mypage", formData)
        assert := lvttest.NewHTTPAssert(resp)

        if resp.StatusCode >= 300 && resp.StatusCode < 400 {
            resp = test.FollowRedirect(resp)
            assert = lvttest.NewHTTPAssert(resp)
        }
        assert.StatusOK(t)
    })
}
```

### Available HTTP Assertions
- `StatusOK(t)`, `StatusCode(t, code)`, `StatusRedirect(t)`
- `Contains(t, text)`, `NotContains(t, text)`, `ContainsAll(t, ...texts)`
- `HasElement(t, selector)`, `ElementCount(t, selector, n)`
- `NoTemplateErrors(t)` - catches unflattened `{{.Field}}`
- `HasCSRFToken(t)`, `HasFormField(t, name)`
- `ContentTypeHTML(t)`, `ContentTypeJSON(t)`

---

## Browser Test Pattern (Tier 2)

Only when absolutely needed:

```go
//go:build browser

package e2e

import (
    "testing"
    "github.com/chromedp/chromedp"
)

func TestRendering_ModalLifecycle(t *testing.T) {
    // Use self-contained HTML to avoid path issues
    html := `<!DOCTYPE html>...embedded test page...`

    ctx, cancel := GetPooledChrome(t)
    defer cancel()

    // Test browser-specific behavior
    chromedp.Run(ctx,
        chromedp.Navigate(testURL),
        chromedp.Click("#open-modal"),
        chromedp.WaitVisible("#modal"),
        // ... assertions
    )
}
```

---

## Makefile Targets

| Target | What It Runs | When to Use |
|--------|--------------|-------------|
| `make test-fast` | Unit tests only | Quick feedback (~30s) |
| `make test-commit` | Unit + HTTP | **Before every commit** (~75s) |
| `make test-http` | HTTP tests only | Testing server logic |
| `make test-browser` | Browser tests | Weekly / pre-release (~45s) |
| `make test-all` | Everything | Full validation |

---

## Build Tag Reference

```go
//go:build http        // HTTP-only, no browser
//go:build browser     // Requires chromedp/browser
//go:build deployment  // Docker/Fly infrastructure tests
// (no tag)            // Always runs
```

---

## Remember

**Default to HTTP tests** - they're 10x faster
**Use `make test-commit`** before committing
Browser tests only for the 12 scenarios listed above
Add correct build tag to every new test file

Don't use browser tests for CRUD/form validation
Don't skip the build tag
Don't add browser tests for server-side logic
