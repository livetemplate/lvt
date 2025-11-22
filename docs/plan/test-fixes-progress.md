# E2E Test Fixes - Progress Tracker

**Goal:** Fix all failing e2e tests after Docker build refactor

**Started:** 2025-11-14

---

## Current Status

**Total Tests:** 57
- ✅ **Passing:** 52
- ❌ **Failing:** 9
- 🔵 **Skipped:** 5

---

## Failing Tests to Fix

### 1. TestKitRuntime_AllKits ❌
**Status:** Not Started
**File:** `e2e/kit_runtime_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 2. TestKitRuntime_TemplateRendering ❌
**Status:** Not Started
**File:** `e2e/kit_runtime_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 3. TestPageModeRendering ❌
**Status:** Not Started
**File:** `e2e/pagemode_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 4. TestServe_Defaults ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 5. TestServe_CustomPort ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 6. TestServe_ModeApp ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 7. TestServe_NoBrowser ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 8. TestServe_NoReload ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 9. TestServe_VerifyServerResponds ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

### 10. TestServe_ContextCancellation ❌
**Status:** Not Started
**File:** `e2e/serve_test.go`
**Error:** Unknown - need to investigate
**Plan:**
- [ ] Read test file and understand what it does
- [ ] Run test in isolation to see specific error
- [ ] Identify root cause
- [ ] Apply fix
- [ ] Verify test passes

---

## Test Groups

### Group 1: Kit Runtime Tests (2 tests)
- TestKitRuntime_AllKits
- TestKitRuntime_TemplateRendering

### Group 2: Page Mode Tests (1 test)
- TestPageModeRendering

### Group 3: Serve Tests (7 tests)
- TestServe_Defaults
- TestServe_CustomPort
- TestServe_ModeApp
- TestServe_NoBrowser
- TestServe_NoReload
- TestServe_VerifyServerResponds
- TestServe_ContextCancellation

---

## Strategy

1. **Analyze failures by group** - Run each group to understand common patterns
2. **Fix one test at a time** - Verify each fix before moving to next
3. **Run full suite after each group** - Ensure no regressions

---

## Notes

- All failures appear after Docker build refactor
- May be related to how tests build/run apps
- Check if tests need Docker build updates similar to what we did for tutorial/pagemode/url_routing tests

---

**Last Updated:** 2025-11-14
