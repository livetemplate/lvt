# Evolution Patterns

Known error patterns and their fixes. This file is the source of truth for the evolution system's knowledge base.

**How to use this file:**
- Each pattern is an H2 section (`## Pattern: id`)
- The evolution system parses this file to match errors to fixes
- Stats (Fix Count, Success Rate) are updated automatically when fixes are applied
- Add new patterns via PR - they become active immediately after merge

**How to add a pattern:**
1. Copy the template at the bottom of this file
2. Fill in the pattern details
3. Submit a PR for review

---

## Pattern: editing-id-type

**Name:** EditingID Type Mismatch
**Confidence:** 0.95
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

EditingID is compared as integer in some templates but is actually a string type in the handler. This causes compilation errors like "cannot convert X to type int".

Evidence: Commit `5700a82` fixed this in single kit template.

### Error Pattern

- **Phase:** compilation
- **Message Regex:** `cannot convert .* to type int`
- **Context Regex:** `EditingID`

### Fix

- **File:** `*/template.tmpl.tmpl`
- **Find:** `{{if ne .EditingID 0}}`
- **Replace:** `{{if ne .EditingID ""}}`
- **Is Regex:** false

---

## Pattern: modal-state-persistence

**Name:** Modal State Persists After Close
**Confidence:** 0.90
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Modal editing state (IsAdding, IsEditing, EditingID) persists on page reload because fields are not marked as transient. User sees modal open after refresh.

Evidence: Commit `8c434c4` added `lvt:"transient"` tags to fix this.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `modal (open|visible|showing) after (reload|refresh|navigation)`
- **Context Regex:** `(IsAdding|IsEditing|EditingID)`

### Fix

- **File:** `*/handler.go.tmpl`
- **Find:** `IsAdding bool`
- **Replace:** `IsAdding bool \`lvt:"transient"\``
- **Is Regex:** false

### Fix 2

- **File:** `*/handler.go.tmpl`
- **Find:** `EditingID string`
- **Replace:** `EditingID string \`lvt:"transient"\``
- **Is Regex:** false

---

## Pattern: form-select-sync

**Name:** Select Value Reverts After Update
**Confidence:** 0.88
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Select dropdown values revert to previous state after morphdom DOM patching. The expected value is not preserved during the update cycle.

Evidence: Commit `bd536bf` added `data-expected-value` attribute and sync script.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `(select|dropdown) value (reverted|reset|changed|wrong)`

### Fix

- **File:** `*/components/form.tmpl`
- **Find:** `<select name="{{.Name}}">`
- **Replace:** `<select name="{{.Name}}" data-expected-value="{{.Value}}">`
- **Is Regex:** false

---

## Pattern: session-not-cleared

**Name:** Session State Not Cleared on Auth Change
**Confidence:** 0.88
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

LiveTemplate session persists stale state after authentication changes. User sees cached IsLoggedIn=false after successful login, or remains "logged in" after logout.

Evidence: Commits `0589544`, `a55b2b5`, `fa25417` all fixed session clearing issues.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `(IsLoggedIn|logged.?in|session).*(stale|persisted|cached|wrong|incorrect)`

### Fix

- **File:** `*/auth/login.go.tmpl`
- **Find:** `return nil`
- **Replace:** `ctx.ClearSession()\n\treturn nil`
- **Is Regex:** false

---

## Pattern: modal-event-propagation

**Name:** Modal Event Propagation Breaks Buttons
**Confidence:** 0.85
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Adding `event.stopPropagation()` on modal inner div breaks event delegation, causing Cancel and other buttons inside the modal to stop working.

Evidence: Commit `c568ec1` removed stopPropagation to fix this.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `(cancel|close|button).*(not working|broken|unresponsive|no effect)`
- **Context Regex:** `modal`

### Fix

- **File:** `*/components/modal.tmpl`
- **Find:** `onclick="event.stopPropagation()"`
- **Replace:** ``
- **Is Regex:** false

---

## Pattern: update-clears-editing-state

**Name:** Update Action Should Clear EditingID
**Confidence:** 0.85
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

After successfully updating a record, the EditingID remains set, keeping the edit modal open. The update action should clear EditingID on success.

Evidence: Commit `b5d3b91` added clearing EditingID after update.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `(edit modal|editing).*(stays open|remains open|won't close|doesn't close)`

### Fix

- **File:** `*/handler.go.tmpl`
- **Find:** `// Update successful`
- **Replace:** `// Update successful\n\ts.EditingID = ""`
- **Is Regex:** false

---

## Pattern: hardcoded-import-path

**Name:** Hardcoded Import Path Instead of Module Name
**Confidence:** 0.92
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Import statements use hardcoded paths like "/database/models" instead of using the module name variable, causing "cannot find package" errors.

Evidence: Commit `a745964` fixed hardcoded import paths.

### Error Pattern

- **Phase:** compilation
- **Message Regex:** `cannot find package.*/database/models`

### Fix

- **File:** `*/handler.go.tmpl`
- **Find:** `"/database/models"`
- **Replace:** `"{{.ModuleName}}/database/models"`
- **Is Regex:** false

---

## Pattern: auth-receiver-type

**Name:** Auth Middleware Wrong Receiver Type
**Confidence:** 0.90
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Auth middleware generated with `*Handler` receiver type instead of `*Controller`, causing compilation errors.

Evidence: Commit `a745964` fixed receiver type in auth middleware.

### Error Pattern

- **Phase:** compilation
- **Message Regex:** `(undefined|invalid).*\*Handler`
- **Context Regex:** `auth|middleware`

### Fix

- **File:** `*/auth/middleware.go.tmpl`
- **Find:** `func (h *Handler)`
- **Replace:** `func (c *Controller)`
- **Is Regex:** false

---

## Pattern: textarea-not-rendered

**Name:** Text Fields Rendered as Input Instead of Textarea
**Confidence:** 0.85
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

Fields marked as "text" type render as single-line `<input>` instead of multi-line `<textarea>`, losing the ability to enter long content.

Evidence: Commit `84132a2` fixed textarea rendering in standalone kit.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `(text|textarea|multiline).*(single.?line|input|not textarea|wrong element)`

### Fix

- **File:** `*/components/form.tmpl`
- **Find:** `{{if eq .Type "text"}}<input`
- **Replace:** `{{if eq .Type "text"}}<textarea`
- **Is Regex:** false

---

## Pattern: json-editing-item

**Name:** EditingItem Should Not Be in JSON Response
**Confidence:** 0.80
**Added:** 2026-01-19
**Fix Count:** 0
**Success Rate:** -

### Description

EditingItem struct sent in JSON response causes client-side tree renderer error "Object value reached string conversion". EditingItem should be excluded from JSON.

Evidence: Commit `3af9a87` excluded EditingItem from JSON serialization.

### Error Pattern

- **Phase:** runtime
- **Message Regex:** `Object value reached string conversion`
- **Context Regex:** `EditingItem`

### Fix

- **File:** `*/handler.go.tmpl`
- **Find:** `EditingItem *{{.ResourceName}}`
- **Replace:** `EditingItem *{{.ResourceName}} \`json:"-"\``
- **Is Regex:** false

---

<!-- TEMPLATE FOR NEW PATTERNS

Copy everything below to add a new pattern:

## Pattern: your-pattern-id

**Name:** Human Readable Name
**Confidence:** 0.85
**Added:** YYYY-MM-DD
**Fix Count:** 0
**Success Rate:** -

### Description

Describe the problem and when it occurs.
Include evidence (commit hash, issue number) if available.

### Error Pattern

- **Phase:** compilation | runtime | template | generation
- **Message Regex:** `regex to match error message`
- **Context Regex:** `optional regex to match surrounding code`

### Fix

- **File:** `glob/pattern/to/file.ext`
- **Find:** `exact text or regex to find`
- **Replace:** `replacement text`
- **Is Regex:** false | true

-->
