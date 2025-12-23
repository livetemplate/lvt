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

### Phase 2: Tier 2 Browser Tests (12 tests) [COMPLETED]
- [x] Create `/e2e/rendering_test.go`
- [x] TestRendering_DOM_ListOperations
- [x] TestRendering_DOM_TableRendering
- [x] TestRendering_Form_SubmitValidation
- [x] TestRendering_Modal_Lifecycle
- [x] TestRendering_Pagination_Navigation
- [x] TestRendering_InfiniteScroll
- [x] TestRendering_Focus_Preservation
- [x] TestRendering_Scroll_Directives
- [x] TestRendering_Lifecycle_Hooks
- [x] TestRendering_Event_Delegation
- [x] TestRendering_WebSocket_Reconnect
- [x] TestRendering_Conditional_Rendering

### Phase 3: Migrate Tests with Build Tags [COMPLETED]
- [x] Add `//go:build http` to HTTP-compatible tests:
  - app_creation_test.go
  - resource_generation_test.go
  - serve_test.go
  - migration_test.go
  - view_generation_test.go
  - agent_doc_validation_test.go
  - agent_skills_validation_test.go
  - css_frameworks_test.go
  - seeding_test.go
  - parsing_test.go
  - type_inference_test.go
  - textarea_fields_test.go
  - pagination_modes_test.go
  - resource_inspection_test.go
  - skill_debug_rendering_test.go
  - kit_runtime_test.go
  - kit_workflow_test.go
  - kit_management_test.go
  - editmode_test.go
- [x] Add `//go:build browser` to browser tests:
  - modal_test.go
  - pagemode_test.go
  - url_routing_test.go
  - tutorial_test.go
  - livetemplate_core_test.go
  - complete_workflow_test.go
  - delete_multi_post_test.go
  - shared_test.go
  - common_test.go
  - chrome_pool.go
  - helpers.go
  - test_main_test.go
- [x] Add `//go:build deployment` to deployment tests:
  - deployment_docker_test.go
  - deployment_fly_test.go
  - deployment_mock_test.go
- [x] Refactor helper files to work with build tags

### Phase 4: Template Updates [COMPLETED]
- [x] Convert `resource/e2e_test.go.tmpl` to HTTP (multi, single, generator)
- [x] Convert `auth/e2e_test.go.tmpl` to HTTP (multi)

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

### Phase 2 (Completed)
- `e2e/rendering_test.go` - 12 focused browser tests for rendering library validation
  - Uses self-contained HTML pattern (embedded test pages)
  - Validates DOM operations, forms, modals, focus, events, scrolling, lifecycle, pagination, infinite scroll, WebSocket

### Phase 4 (Completed)
- `internal/kits/system/multi/templates/resource/e2e_test.go.tmpl` - HTTP-based resource tests
- `internal/kits/system/single/templates/resource/e2e_test.go.tmpl` - HTTP-based resource tests
- `internal/generator/templates/resource/e2e_test.go.tmpl` - HTTP-based resource tests
- `internal/kits/system/multi/templates/auth/e2e_test.go.tmpl` - HTTP-based auth tests
  - All templates now use `//go:build http` tag
  - Use `SetupHTTP()`, `test.Get()`, `test.PostForm()`, `NewHTTPAssert()` pattern
  - Tests: Initial Load, Add Resource, Search, Form Validation, CSRF Protection

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
