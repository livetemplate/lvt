# Skipped E2E Tests - Investigation and Remediation Plan

**Last Updated**: 2025-11-19
**Status**: 6 tests currently skipped (2 need fixes, 4 intentionally disabled)

## Summary

This document tracks all skipped tests in the e2e suite, their root causes, and remediation plans.

### Quick Stats
- **Total skipped**: 6 tests
- **Need fixing**: 2 tests (TestCompleteWorkflow_BlogApp + 1 TutorialE2E subtest)
- **Intentionally disabled**: 4 tests (deployment tests - by design)
- **Recently fixed**: 2 tests (Modal_Delete_with_Confirmation, Validation_Errors)

---

## Tests That Need Fixing

### 1. TestCompleteWorkflow_BlogApp

**File**: `e2e/complete_workflow_test.go:20`

**Status**: ❌ Skipped - Client library loading failure

**Current Skip Message**:
```go
t.Skip("Temporarily skipped: Client library loading issue in Docker - needs unpkg CDN fix or local embed")
```

**Root Cause**:
- Test builds complete Docker image of blog app
- App attempts to load client library from unpkg CDN: `https://unpkg.com/@livetemplate/client@latest/dist/livetemplate-client.browser.js`
- Client library fails to load (CDN unavailable or wrong version)
- All UI tests fail with: `client: missing`, `data-lvt-loading: true`

**Failed Subtests When Enabled**:
- ❌ WebSocket_Connection (20s timeout)
- ✅ Posts_Page_Loads (passes without client)
- ❌ Create_and_Edit_Post (5s timeout)
- ❌ Delete_Post (5s timeout)
- ❌ Validation_Errors (5s timeout)
- ✅ Infinite_Scroll (passes - static check)
- ✅ Server_Logs_Check (passes)
- ✅ Console_Logs_Check (passes)

**Remediation Options**:

**Option A: Embed Client Library in Docker Image** (Recommended)
- Modify Dockerfile to copy embedded client library
- Update app to serve from local file instead of unpkg
- Ensures test independence from external CDN
- Aligns with approach used in other tests

**Option B: Fix unpkg CDN Reference**
- Verify correct unpkg URL and version
- Add retry logic for CDN failures
- Less reliable (external dependency)

**Option C: Use Native Build Instead of Docker**
- Convert test to use `buildAndRunNative()` like other tests
- Faster execution (~50s vs Docker build)
- Loses Docker deployment validation

**Estimated Effort**: Medium (4-8 hours)
- Requires Dockerfile modifications
- Client library embedding setup
- Test verification

**Priority**: Medium - This is a comprehensive integration test but other tests cover most functionality

---

### 2. TestTutorialE2E/Modal_Delete_with_Confirmation

**File**: `e2e/tutorial_test.go:323`

**Status**: ✅ **FIXED** (2025-11-18)

**Changes Made**:
1. ✅ Removed `t.Skip()` statement
2. ✅ Added `ensureTutorialPostExists()` call to make test independent
3. ✅ Added delete button with `lvt-confirm` attribute to edit modal template (internal/kits/system/multi/templates/resource/template.tmpl.tmpl:156)
4. ✅ Fixed modal selector to wait for correct edit modal (`form[lvt-submit="update"]`)

**Root Cause**:
The test had TWO issues:
1. **Test dependency**: Test depended on data from previous "Add Post" test
2. **Wrong modal selector**: Test was waiting for `input[name="title"]` which exists in BOTH add and edit modals. When "Add Post" test ran first, its add modal was still in the DOM, causing the selector to find the wrong modal.

**Solution**:
- Made test independent by calling `ensureTutorialPostExists()` to create its own data
- Changed wait condition from `input[name="title"]` to `form[lvt-submit="update"]` to specifically wait for the edit modal
- Added delete button with `lvt-confirm` to edit modal template

**Test Results**:
- ✅ Passes: `go test -run "TestTutorialE2E/Modal_Delete_with_Confirmation$"`  (isolation)
- ✅ Passes: `go test -run "^TestTutorialE2E$"` (full suite)

**Commits**:
- `9ed743c`: Add delete button to edit modal template
- `5f668f5`: Fix modal selector to wait for correct edit modal

**Priority**: ✅ COMPLETE

---

### 3. TestTutorialE2E/Validation_Errors

**File**: `e2e/tutorial_test.go:653`

**Status**: ✅ **FIXED** (2025-11-19)

**WebSocket Response Analysis** (v0.3.1):
```json
{
  "tree": {},
  "meta": {
    "success": false,
    "errors": {
      "Content": "Content is required",
      "Title": "Title is required"
    },
    "action": "add"
  }
}
```

**Client State After Response**:
```javascript
window.liveTemplateClient.errors = {}  // ❌ EMPTY!
```

**Root Cause**: The **JavaScript client library** is NOT handling `meta.errors` from the WebSocket response!

**What's Working**:
- ✅ Server captures validation errors via MultiError
- ✅ Server sends errors in `meta.errors` in WebSocket response

**What's Broken**:
- ❌ Client library doesn't extract `meta.errors` from response
- ❌ Client library doesn't store errors in `window.liveTemplateClient.errors`
- ❌ Client library doesn't inject error HTML (`<small>` tags) into the form

**Expected Flow**:
1. Server sends `meta.errors` in WebSocket response ✅
2. Client extracts errors from response and stores them ❌
3. Client dynamically injects `<small>` error tags into form ❌
4. User sees validation errors ❌

**Root Cause** (Found 2025-11-19):
The issue was a **field name casing mismatch** in the livetemplate library:

Templates check for errors using lowercase field names (matching HTML form inputs):
```html
{{if .lvt.HasError "title"}}
<small style="color: #c00;">{{.lvt.Error "title"}}</small>
{{end}}
```

But the server was sending validation errors with PascalCase field names (matching Go struct fields):
```json
{
  "errors": {
    "Title": "Title is required",    // ❌ PascalCase
    "Content": "Content is required"  // ❌ PascalCase
  }
}
```

This caused `.lvt.HasError "title"` to return false, so error conditionals never rendered.

**Fix Applied** (livetemplate v0.3.2):
Modified `ValidationToMultiError()` in `action.go` to convert struct field names to lowercase:

```go
for _, e := range validationErrs {
    structFieldName := e.Field()         // "Title"
    formFieldName := strings.ToLower(structFieldName)  // "title"

    fieldErrors = append(fieldErrors, FieldError{
        Field:   formFieldName,  // ✅ Use lowercase to match HTML input names
        Message: message,
    })
}
```

**Changes Made**:
1. ✅ Fixed `ValidationToMultiError` in livetemplate library (action.go:295-319)
2. ✅ Released livetemplate v0.3.2 with the fix
3. ✅ Updated lvt to use livetemplate v0.3.2
4. ✅ Verified all livetemplate library tests pass
5. ✅ Removed `t.Skip()` from test (was never needed after fix)

**Files Modified**:
- **Livetemplate Library** (`github.com/livetemplate/livetemplate`):
  - `action.go`: Fixed field name casing in ValidationToMultiError
- **LVT** (`github.com/livetemplate/lvt`):
  - `go.mod`: Updated to livetemplate v0.3.2

**Test Results**:
- ✅ Passes: `go test -run "TestTutorialE2E/Validation_Errors$"` (isolation)
- ✅ Passes: Full E2E suite `go test ./e2e/`
- ✅ Error messages correctly displayed in UI:
  ```
  Title is required
  Content is required
  ```

**WebSocket Response** (After Fix):
```json
{
  "tree": {
    "9": {"0": "Title is required", "s": [...]},
    "11": {"0": "Content is required", "s": [...]}
  },
  "meta": {
    "success": false,
    "errors": {
      "title": "Title is required",    // ✅ lowercase!
      "content": "Content is required"  // ✅ lowercase!
    }
  }
}
```

**Commits**:
- livetemplate: `a6343bae` - fix: convert validation error field names to lowercase
- lvt: Pending - Update to livetemplate v0.3.2

**Priority**: ✅ COMPLETE - Server-side validation rendering now working correctly

---

## Intentionally Disabled Tests (By Design)

These tests are correctly disabled and require explicit opt-in via environment variables.

### 4. TestDockerDeploymentWithResources

**File**: `e2e/deployment_docker_test.go:110`

**Condition**: Requires `RUN_DOCKER_DEPLOYMENT_TESTS=true`

**Purpose**: Long-running Docker deployment integration test

**Status**: ✅ Correctly disabled - CI/CD opt-in by design

---

### 5. TestDockerDeploymentQuickSmoke

**File**: `e2e/deployment_docker_test.go:188`

**Condition**: Requires `RUN_DOCKER_DEPLOYMENT_TESTS=true`

**Purpose**: Quick smoke test for Docker deployments

**Status**: ✅ Correctly disabled - CI/CD opt-in by design

---

### 6. TestRealFlyDeployment

**File**: `e2e/deployment_fly_test.go:54`

**Condition**: Requires `FLY_API_TOKEN` environment variable

**Purpose**: Real Fly.io deployment test (requires cloud credentials)

**Status**: ✅ Correctly disabled - requires cloud credentials

---

### 7. TestFlyDeploymentWithResources

**File**: `e2e/deployment_fly_test.go:114`

**Condition**: Requires `RUN_FLY_DEPLOYMENT_TESTS=true`

**Purpose**: Fly.io deployment with resources

**Status**: ✅ Correctly disabled - CI/CD opt-in by design

---

## Remediation Priority

### Completed ✅
1. ✅ **TestTutorialE2E/Modal_Delete_with_Confirmation** - Fixed (2025-11-18)
2. ✅ **TestTutorialE2E/Validation_Errors** - Fixed (2025-11-19)

### Medium Priority (Moderate Effort)
3. **TestCompleteWorkflow_BlogApp** - Embed client library in Docker (4-8 hours)

### No Action Needed
4-7. Deployment tests - Intentionally disabled by design ✅

---

## Test Suite Performance

**Before Optimizations**: 130 seconds
**After Optimizations**: 63.6 seconds
**Improvement**: 51% faster

**Recent Fixes** (2025-11-18):
- ✅ Fixed TestModalFunctionality flakiness (working directory issue)
- ✅ Optimized server startup detection (exponential backoff)
- ✅ Fixed I/O hanging in tests

---

## How to Enable Deployment Tests

### Docker Deployment Tests
```bash
export RUN_DOCKER_DEPLOYMENT_TESTS=true
go test ./e2e -v -run "Docker"
```

### Fly.io Deployment Tests
```bash
export FLY_API_TOKEN="your-token-here"
export RUN_FLY_DEPLOYMENT_TESTS=true
go test ./e2e -v -run "Fly"
```

---

## Notes for Future Maintainers

1. **Test Independence**: Always ensure tests can run independently. Avoid depending on data from previous tests.

2. **Client Library Loading**:
   - Use embedded client library (`e2etest.GetClientLibraryJS()`) for reliability
   - Avoid CDN dependencies in tests
   - See `e2e/modal_test.go` for reference implementation

3. **Skip Messages**: Keep skip messages accurate and descriptive. Include:
   - What is broken
   - Root cause if known
   - What needs to be fixed

4. **Docker Tests**: When testing Docker deployments:
   - Embed dependencies in image
   - Don't rely on external CDNs
   - Test can be slow (60s+ for build)

5. **Performance**: Current suite runs in ~64 seconds. Target is <60 seconds for all tests.
