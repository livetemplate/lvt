# Skipped E2E Tests - Investigation and Remediation Plan

**Last Updated**: 2025-11-18
**Status**: 7 tests currently skipped (3 need fixes, 4 intentionally disabled)

## Summary

This document tracks all skipped tests in the e2e suite, their root causes, and remediation plans.

### Quick Stats
- **Total skipped**: 7 tests
- **Need fixing**: 3 tests (TestCompleteWorkflow_BlogApp + 2 TutorialE2E subtests)
- **Intentionally disabled**: 4 tests (deployment tests - by design)

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

**Status**: 🔄 In Progress - Partially Fixed

**Changes Made** (2025-11-18):
1. ✅ Removed `t.Skip()` statement
2. ✅ Added `ensureTutorialPostExists()` call to make test independent
3. ✅ Added delete button with `lvt-confirm` attribute to edit modal template (internal/kits/system/multi/templates/resource/template.tmpl.tmpl:156)
4. ✅ Verified embedded template contains the delete button

**Current Issue**:
Test passes when run in isolation but fails when run as part of full `TestTutorialE2E` suite. The delete button is not being found in the modal during full suite execution.

**Root Cause Analysis**:
- Original issue: Test depended on data from previous test ("My First Blog Post")
- Original issue: Delete button was missing from edit modal template
- **Fixed**: Added post creation at start of test via `ensureTutorialPostExists()`
- **Fixed**: Added delete button to edit modal template
- **Remaining**: Test cache or build cache issue preventing template changes from being consistently applied in full suite runs

**Test Results**:
- ✅ Passes: `go test -run "TestTutorialE2E/Modal_Delete_with_Confirmation$"`  (isolation)
- ❌ Fails: `go test -run "^TestTutorialE2E$"` (full suite)
- ✅ Embedded template verified to contain `lvt-confirm` attribute

**Next Steps**:
1. Investigate why full suite run doesn't pick up template changes
2. May need to add explicit cache busting or build steps
3. Consider adding debug output to capture actual modal HTML during test

**Priority**: Medium - Core functionality is implemented, need to resolve test execution issue

---

### 3. TestTutorialE2E/Validation_Errors

**File**: `e2e/tutorial_test.go:643`

**Status**: ❌ Skipped - Product bug

**Current Skip Message**:
```go
t.Skip("Skipping until conditional rendering bug is fixed")
```

**Root Cause**:
- Conditional rendering bug in product code
- Validation errors not displaying correctly
- Likely related to template conditional rendering

**Context**:
This appears to be a known product bug rather than a test issue. The test is correctly written but the feature it tests is broken.

**Remediation**:

**Step 1: Identify the Bug**
- Run test with skip removed
- Capture exact failure mode
- Check if validation errors render in browser
- Review template conditional logic

**Step 2: Fix Product Bug**
- Fix conditional rendering in templates
- Ensure validation errors display properly
- May require template engine changes

**Step 3: Re-enable Test**
- Remove skip
- Verify test passes
- Add to CI

**Estimated Effort**: Unknown (depends on bug complexity)
- Could be simple template fix (2-4 hours)
- Could require engine changes (8-16 hours)
- Needs investigation first

**Priority**: Medium - Validation is important UX feature

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

### High Priority (Simple Fixes)
1. ✅ **TestTutorialE2E/Modal_Delete_with_Confirmation** - Make test independent (1-2 hours)

### Medium Priority (Moderate Effort)
2. **TestCompleteWorkflow_BlogApp** - Embed client library in Docker (4-8 hours)
3. **TestTutorialE2E/Validation_Errors** - Investigate and fix product bug (2-16 hours)

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
