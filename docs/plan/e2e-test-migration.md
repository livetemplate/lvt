# E2E Test Migration Plan: 4-Tier Testing Strategy

## Problem

E2E chromedp browser tests freeze the laptop (CPU + memory + time exhaustion). Current suite: 98 browser-based tests across 33 files, running 8 Chrome containers in parallel.

## Solution

Migrate to a 4-tier testing strategy that maintains rendering library confidence while minimizing browser usage.

---

## Tier Overview

| Tier | What | How | Speed | When to Run |
|------|------|-----|-------|-------------|
| **0** | Client JS logic | Jest + jsdom (exists) | ~5s | Always |
| **1** | Server HTML output | Go httptest (new) | ~10s | Always |
| **2** | Rendering paths | Browser (12 tests) | ~45s | Weekly / Pre-release |
| **3** | Application logic | HTTP tests (migrate) | ~60s | Always |

**Result:** `make test-commit` runs in ~75 seconds without browser (vs 10+ minutes currently)

---

## Progress Tracker

### Phase 1: HTTP Testing Framework [COMPLETED]
- [x] Create `/testing/http.go` - HTTPTest struct and setup
- [x] Create `/testing/http_assertions.go` - HTTP assertions (25+ methods)
- [x] Add form submission helpers (POST, multipart)
- [x] Add HTML response parsing
- [x] Add template expression validation
- [x] Add database state verification
- [x] Write unit tests for HTTP framework (18 tests passing)

### Phase 5: CI/Makefile [COMPLETED]
- [x] Update Makefile with new targets (test-http, test-browser)
- [x] Add build tags to HTTP test files
- [x] Create example HTTP tests in e2e/

### Phase 2: Tier 2 Browser Tests (12 tests) [PENDING]
- [ ] Create `/e2e/rendering_test.go`
- [ ] TestDOM_ListOperations
- [ ] TestDOM_TableRendering
- [ ] TestForm_SubmitValidation
- [ ] TestModal_Lifecycle
- [ ] TestPagination_Navigation
- [ ] TestInfiniteScroll
- [ ] TestFocus_Preservation
- [ ] TestScroll_Directives
- [ ] TestLifecycle_Hooks
- [ ] TestEvent_Delegation
- [ ] TestWebSocket_Reconnect
- [ ] TestConditional_Rendering

### Phase 3: Migrate 98 Tests [PENDING]
- [ ] Migrate app_creation_test.go (4 tests)
- [ ] Migrate resource_generation_test.go (12 tests)
- [ ] Migrate serve_test.go (9 tests)
- [ ] Migrate kit_management_test.go (9 tests)
- [ ] Migrate remaining tests

### Phase 4: Template Updates [PENDING]
- [ ] Convert `resource/e2e_test.go.tmpl` to HTTP
- [ ] Convert `auth/e2e_test.go.tmpl` to HTTP

---

## New Makefile Targets

```makefile
# Fast feedback - unit tests only (~30 seconds)
make test-fast

# Before commit - unit + HTTP tests, NO browser (~75 seconds)
make test-commit

# HTTP tests only
make test-http

# Browser rendering tests only - Tier 2 (~45 seconds)
make test-browser

# Full validation - all tiers including browser
make test-all
```

---

## Files Created/Modified

### Phase 1 (Completed)
- `testing/http.go` - Core HTTP testing framework
- `testing/http_assertions.go` - 25+ assertion methods
- `testing/http_test.go` - Unit tests

### Phase 5 (Completed)
- `Makefile` - Updated with tiered test targets
- `e2e/http_example_test.go` - Example HTTP tests with `//go:build http` tag

---

## Expected Outcome After Full Migration

| Metric | Before | After |
|--------|--------|-------|
| `make test-commit` time | 10+ minutes | ~75 seconds |
| Browser tests count | 98 | 12 (Tier 2) + ~10 (retained) |
| HTTP tests count | 0 | ~75 |
| Memory usage | 8 Chrome containers | 1 Chrome (when needed) |
| CPU usage | Freezes laptop | Normal |
| Coverage | 100% | 100% (same) |
