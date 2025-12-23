# Testing Migration Guide

Guide for updating branches to the 4-tier testing strategy.

## Quick Checklist

- [ ] Add build tags to e2e test files
- [ ] Convert auth template to HTTP tests
- [ ] Update `testing/chrome.go` with resource limits
- [ ] Update `e2e/chrome_pool.go` pool size
- [ ] Add signal handler to `e2e/shared_test.go`
- [ ] Update `testing/testing.go` ChromeLocal cleanup

---

## 1. Add Build Tags to E2E Tests

### HTTP Tests (Tier 1)
Tests that don't need JavaScript execution:

```go
//go:build http

package e2e
```

**Which tests get `//go:build http`:**
- `app_creation_test.go`
- `resource_generation_test.go`
- `serve_test.go`
- `migration_test.go`
- `view_generation_test.go`
- `agent_doc_validation_test.go`
- `agent_skills_validation_test.go`
- `css_frameworks_test.go`
- `seeding_test.go`
- `parsing_test.go`
- `type_inference_test.go`
- `textarea_fields_test.go`
- `pagination_modes_test.go`
- `resource_inspection_test.go`
- `skill_debug_rendering_test.go`
- `kit_runtime_test.go`
- `kit_workflow_test.go`
- `kit_management_test.go`
- `editmode_test.go`

### Browser Tests (Tier 2)
Tests requiring JavaScript/DOM:

```go
//go:build browser

package e2e
```

**Which tests get `//go:build browser`:**
- `modal_test.go`
- `pagemode_test.go`
- `url_routing_test.go`
- `tutorial_test.go`
- `livetemplate_core_test.go`
- `complete_workflow_test.go`
- `delete_multi_post_test.go`
- `shared_test.go`
- `common_test.go`
- `chrome_pool.go`
- `helpers.go`
- `test_main_test.go`
- `rendering_test.go`

### Deployment Tests (Tier 3)

```go
//go:build deployment

package e2e
```

**Which tests get `//go:build deployment`:**
- `deployment_docker_test.go`
- `deployment_fly_test.go`
- `deployment_mock_test.go`

---

## 2. Convert Auth Template to HTTP Tests

**File:** `internal/kits/system/multi/templates/auth/e2e_test.go.tmpl`

### Before (Browser-based)
```go
package auth

import (
    "github.com/chromedp/chromedp"
    // ...
)

func TestAuthE2E(t *testing.T) {
    // Uses chromedp, StartDockerChrome, etc.
}
```

### After (HTTP-based)
```go
//go:build http

package auth

import (
    "net/url"
    "testing"
    lvttest "github.com/livetemplate/lvt/testing"
)

func TestAuthHTTP(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping HTTP test in short mode")
    }

    test := lvttest.SetupHTTP(t, &lvttest.HTTPSetupOptions{
        AppPath: "../../cmd/[[.ModuleName]]/main.go",
    })
    defer test.Cleanup()

    t.Run("Auth Page Loads", func(t *testing.T) {
        resp := test.Get("/auth")
        assert := lvttest.NewHTTPAssert(resp)

        assert.StatusOK(t)
        assert.ContentTypeHTML(t)
        assert.Contains(t, "Email")
        assert.NoTemplateErrors(t)
        assert.HasFormField(t, "Email")
    })

    // ... more tests
}
```

**Key changes:**
- Add `//go:build http` tag
- Use `lvttest.SetupHTTP()` instead of `chromedp`
- Use `test.Get()`, `test.PostForm()` for requests
- Use `lvttest.NewHTTPAssert()` for assertions
- Use `[[` delimiters (not `{{`) for template variables

---

## 3. Update testing/chrome.go

### Add Resource Limits

```go
const (
    dockerImage           = "chromedp/headless-shell:latest"
    chromeContainerPrefix = "chrome-e2e-test-"
    staleContainerGrace   = 1 * time.Minute  // Was 10 minutes
)
```

### Add Docker Flags

```go
cmd := exec.Command("docker", "run", "-d",
    "--rm",              // Auto-remove on stop
    "--memory", "512m",  // Memory limit
    "--cpus", "0.5",     // CPU limit
    "-p", portMapping,
    "--name", containerName,
    "--add-host", "host.docker.internal:host-gateway",
    dockerImage,
)
```

---

## 4. Update e2e/chrome_pool.go

### Reduce Pool Size

```go
// Pool size of 4 to stay under 4GB memory limit (4 × 512MB = 2GB max)
chromePool = NewChromePool(t, 4)  // Was 8
```

---

## 5. Add Signal Handler to e2e/shared_test.go

### Add Imports

```go
import (
    "os/signal"
    "syscall"
    // ... existing imports
)
```

### Add Handler in TestMain

```go
func TestMain(m *testing.M) {
    // Setup signal handling for cleanup on interrupt (Ctrl+C)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigCh
        log.Println("🛑 Interrupted - cleaning up Chrome containers...")
        chromePoolMu.Lock()
        if chromePool != nil {
            chromePool.Cleanup()
        }
        chromePoolMu.Unlock()
        cleanupChromeContainers()
        log.Println("✅ Cleanup complete")
        os.Exit(1)
    }()

    // ... rest of TestMain
}
```

---

## 6. Update testing/testing.go ChromeLocal Cleanup

```go
func (e *E2ETest) Cleanup() {
    // ... existing cleanup ...

    // Stop Chrome
    if e.ChromeMode == ChromeDocker {
        StopDockerChrome(e.T, e.ChromePort)
    }

    // For ChromeLocal, give time for process cleanup
    if e.ChromeMode == ChromeLocal {
        time.Sleep(100 * time.Millisecond)
    }

    // ... rest of cleanup
}
```

---

## HTTP Test Patterns

### Basic Page Load Test
```go
t.Run("Page Loads", func(t *testing.T) {
    resp := test.Get("/mypage")
    assert := lvttest.NewHTTPAssert(resp)

    assert.StatusOK(t)
    assert.Contains(t, "Expected Content")
    assert.NoTemplateErrors(t)
})
```

### Form Submission Test
```go
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
```

### CSRF Token Test
```go
t.Run("CSRF Protection", func(t *testing.T) {
    resp := test.Get("/mypage")
    assert := lvttest.NewHTTPAssert(resp)

    assert.StatusOK(t)
    assert.HasCSRFToken(t)
})
```

---

## Available HTTP Assertions

| Method | Description |
|--------|-------------|
| `StatusOK(t)` | Assert 200 status |
| `StatusCode(t, code)` | Assert specific status |
| `StatusRedirect(t)` | Assert 3xx redirect |
| `Contains(t, text)` | Body contains text |
| `NotContains(t, text)` | Body doesn't contain text |
| `ContainsAll(t, ...texts)` | Body contains all texts |
| `HasElement(t, selector)` | Has CSS selector match |
| `ElementCount(t, selector, n)` | Exact element count |
| `NoTemplateErrors(t)` | No `{{.Field}}` in output |
| `HasCSRFToken(t)` | Has CSRF token |
| `HasFormField(t, name)` | Has form field |
| `ContentTypeHTML(t)` | Content-Type is HTML |
| `ContentTypeJSON(t)` | Content-Type is JSON |

---

## Makefile Targets

```makefile
test-fast      # Unit tests only (~30s)
test-commit    # Unit + HTTP (~75s) - run before commits
test-http      # HTTP tests only
test-browser   # Browser tests only (~45s)
test-all       # Everything
```

---

## Decision Guide: HTTP vs Browser

| Testing... | Use |
|------------|-----|
| Form submission | HTTP |
| Database CRUD | HTTP |
| Template rendering | HTTP |
| CSRF protection | HTTP |
| Page loads/redirects | HTTP |
| Modal open/close | Browser |
| Focus preservation | Browser |
| WebSocket reconnection | Browser |
| JavaScript interactions | Browser |
| CSS animations | Browser |

**Default to HTTP tests** - they're 10x faster.
