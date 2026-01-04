---
name: lvt-debug-rendering
description: Debug LiveTemplate rendering issues - template not updating, partial renders, action dispatch errors, WebSocket problems
category: maintenance
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:debug-rendering

Systematic debugging guide for LiveTemplate rendering issues. Covers the full rendering pipeline from server-side tree generation to client-side DOM patching, helping diagnose and fix issues where templates don't update, render incorrectly, or actions fail to dispatch.

## Activation Rules

### Context Detection

This skill typically runs in **existing LiveTemplate projects** (.lvtrc exists).

**Context Established By:**
1. **Project context** - `.lvtrc` exists (most common scenario)
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Keyword matching** (case-insensitive): `lvt`, `livetemplate`, `lt`

### Trigger Patterns

**With Context:**
- "Template not updating"
- "UI won't refresh"
- "Actions not working"
- "Partial render"
- "WebSocket disconnected"
- "State not syncing"

**Without Context (needs keywords):**
- Must mention "lvt", "livetemplate", or "lt"

---

## Quick Symptom Lookup

| Symptom | Likely Phase | First Check | Key Files |
|---------|--------------|-------------|-----------|
| Template not updating after data change | Diff | Tree cache state | `tree_compare.go`, `template.go` |
| Partial/broken HTML | Build | Statics array alignment | `parse.go`, `types.go` |
| Action dispatch error | Dispatch | Method signature | `dispatch.go` |
| Range items not updating | Diff/Range | Key detection | `range_ops.go`, `tree-renderer.ts` |
| Empty response from server | Build | Tree generation errors | `template.go` |
| morphdom not applying changes | Client | Tree state merging | `tree-renderer.ts` |
| WebSocket message not received | Transport | Connection state | `websocket.ts`, `mount.go` |
| Statics missing on update | Registry | Client structure tracking | `signature/` package |

---

## Rendering Pipeline Overview

LiveTemplate uses a 5-phase rendering pipeline:

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  Parse  │ -> │  Build  │ -> │  Diff   │ -> │ Render  │ -> │  Send   │
└─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘
     │              │              │              │              │
 Template     TreeNode with    Changed       HTML string    WebSocket/
 -> AST       statics/dynamics  fields only  (full render)   HTTP
```

**Phase 1: Parse** - Template string -> AST -> Tree structure
**Phase 2: Build** - AST + Data -> TreeNode (statics + dynamics)
**Phase 3: Diff** - Compare oldTree vs newTree -> Minimal changes
**Phase 4: Render** - TreeNode -> HTML (for full renders)
**Phase 5: Send** - JSON over WebSocket or HTTP response

---

## Priority Issue Workflows

### Issue 1: Template Not Updating

**Symptoms:**
- Data changes but UI stays the same
- Actions execute (server logs show) but no visual update
- `window.__wsMessages` shows messages but DOM unchanged

**Investigation Checklist:**

```markdown
1. [ ] Verify action reached server
   - Check server logs for action dispatch
   - Look for DispatchError in logs

2. [ ] Verify state mutation
   - Add logging in controller action method
   - Confirm state changed before return
   - Check: are you returning the new state?

3. [ ] Check tree generation
   - Log buildTree() output
   - Verify dynamics contain changed values

4. [ ] Check diff generation
   - Log output of CompareTreesAndGetChangesWithPath()
   - Verify changed fields appear in result
   - Check if fingerprint changed

5. [ ] Check WebSocket delivery (Browser Console)
   window.__wsMessages  // See all messages
   window.__lastWSMessage  // Last message content

6. [ ] Check client-side merge
   liveTemplateClient.getTreeState()  // Current state

7. [ ] Check DOM update
   - Inspect element in browser DevTools
   - Check if morphdom applied changes
```

**Common Root Causes:**

| Cause | Detection | Fix |
|-------|-----------|-----|
| State mutation without returning new state | Logging shows state changed but not returned | Return new state from controller method |
| TreeNode fingerprint unchanged | Diff returns empty object | Ensure data actually differs |
| Client rangeState out of sync | Range items update wrong elements | Check range key consistency |
| morphdom key mismatch | Elements replaced instead of updated | Use consistent `data-key` or `data-lvt-key` |

**Example Fix - State Not Returned:**

```go
// WRONG - state mutated but not returned
func (c *Counter) Increment(state CounterState, ctx *livetemplate.Context) (CounterState, error) {
    state.Count++
    return state, nil  // Missing return was: return CounterState{}, nil
}

// CORRECT
func (c *Counter) Increment(state CounterState, ctx *livetemplate.Context) (CounterState, error) {
    state.Count++
    return state, nil
}
```

---

### Issue 2: Partial/Broken Renders

**Symptoms:**
- Some parts of template render, others don't
- HTML looks malformed or truncated
- Nested elements missing
- Range items render incorrectly

**Investigation Checklist:**

```markdown
1. [ ] Check statics array alignment
   - len(statics) should equal len(dynamics) + 1
   - Statics interleave with dynamics: s[0] + d[0] + s[1] + d[1] + s[2]

2. [ ] Inspect TreeNode structure
   - JSON serialize tree for inspection
   - Check for missing "s" key in nested nodes
   - Verify numeric keys are sequential ("0", "1", "2"...)

3. [ ] Check template parsing
   lvt parse app/resource/resource.tmpl
   - Look for parse errors
   - Check template action syntax

4. [ ] Check data type matches
   - Verify template expects correct data structure
   - Check for nil pointers in data

5. [ ] Client-side reconstruction (Browser Console)
   // Get tree state
   liveTemplateClient.getTreeState()

   // Check if statics present
   // First render should have "s" arrays
   // Updates may omit statics (cached on client)
```

**Common Root Causes:**

| Cause | Detection | Fix |
|-------|-----------|-----|
| Nested TreeNode missing statics | First nested render shows as empty | Ensure registry marks structures as seen after first render |
| Range statics stripped prematurely | Range items missing HTML wrapper | Check client-side range caching |
| Dynamic key numbering gap | Tree has "0", "2" but no "1" | Fix tree building logic |
| Template composition not flattened | `{{template}}` references fail | Ensure templates are properly composed |

**Debug: Serialize TreeNode (Server)**

```go
// Add to handler for debugging
import "encoding/json"

func debugTree(tree *TreeNode) {
    data, _ := json.MarshalIndent(tree, "", "  ")
    log.Printf("Tree:\n%s", string(data))
}
```

---

### Issue 3: Action Dispatch Errors

**Symptoms:**
- "method not found" errors in logs
- Button clicks do nothing
- Form submissions fail silently
- DispatchError in server output

**Investigation Checklist:**

```markdown
1. [ ] Check method signature
   Controller+State pattern requires:
   func(state State, ctx *Context) (State, error)

2. [ ] Check action name mapping
   - Template: lvt-click="increment"
   - Go: func (c *Counter) Increment(...)
   - Note: snake_case in template -> PascalCase method

3. [ ] Check method visibility
   - Method must be exported (uppercase first letter)
   - Method must be on correct receiver type

4. [ ] Check DispatchError details
   - Action name attempted
   - Store type searched
   - Available methods on type
```

**Common Root Causes:**

| Cause | Detection | Fix |
|-------|-----------|-----|
| Method not exported | DispatchError shows lowercase method | Capitalize method name |
| Wrong signature | Method found but not matching | Use exact `(state State, ctx *Context) (State, error)` |
| Action name typo | DispatchError shows misspelled action | Fix template attribute |
| Wrong receiver type | Method on value receiver, store uses pointer | Match receiver to store usage |

**Example - DispatchError Output:**

```
Error: method not found: action 'incremnt' not found on type *Counter
Available methods: Increment, Decrement, Reset
```

**Fix:** Correct typo in template from `lvt-click="incremnt"` to `lvt-click="increment"`

**Example - Wrong Signature:**

```go
// WRONG - missing context parameter
func (c *Counter) Increment(state CounterState) (CounterState, error) {
    ...
}

// CORRECT
func (c *Counter) Increment(state CounterState, ctx *livetemplate.Context) (CounterState, error) {
    ...
}
```

---

## Server-Side Debugging Guide

### Phase 1: Parse (`internal/parse/`)

**What happens:** Template string -> AST -> Tree structure

**Key functions:**
- `Parse()` - Main entry point
- `evaluatePipe()` - Expression evaluation

**Debug points:**
```go
// Check parse result
tree, err := parse.Parse(templateString)
if err != nil {
    log.Printf("Parse error: %v", err)
}
```

**Common issues:**
- Unhandled node type in AST
- Template not flattened before parsing ({{define}}/{{template}})
- Function map not available

---

### Phase 2: Build (`internal/build/`)

**What happens:** AST + Data -> TreeNode with statics/dynamics

**Key types:**
```go
type TreeNode struct {
    Statics  []string               // Static HTML fragments
    Dynamics map[string]interface{} // Dynamic values at positions
    Range    *RangeData             // Range metadata if present
    Metadata *TreeMetadata          // Additional metadata
}
```

**Debug points:**
```go
// Log tree structure
tree := buildTree(ast, data)
treeJSON, _ := json.MarshalIndent(tree, "", "  ")
log.Printf("Built tree:\n%s", string(treeJSON))
```

**Common issues:**
- Statics/dynamics count mismatch
- Nested TreeNode not properly constructed
- Range items not keyed correctly

---

### Phase 3: Diff (`internal/diff/`)

**What happens:** oldTree vs newTree -> Minimal changes

**Key functions:**
- `CompareTreesAndGetChangesWithPath()` - Main entry point
- `FindRangeConstructMatches()` - Range item matching

**Debug points:**
```go
// Log diff result
changes := diff.CompareTreesAndGetChangesWithPath(oldTree, newTree, path)
log.Printf("Old dynamics: %+v", oldTree.Dynamics)
log.Printf("New dynamics: %+v", newTree.Dynamics)
log.Printf("Changes: %+v", changes.Dynamics)
```

**Range operations format:**
```json
{
  "d": [
    ["u", "item-3", {"0": "updated value"}],   // Update
    ["i", "item-2", "0", {"s": [...], "0": ...}], // Insert after
    ["r", "item-5"],                            // Remove
    ["o", ["id-1", "id-2", "id-3"]]            // Reorder
  ]
}
```

**Common issues:**
- Range operations not generated (empty diff)
- Wrong items marked as changed
- Statics incorrectly stripped from updates

---

### Phase 4: Render (`internal/render/`, `internal/context/`)

**What happens:** Execute Go template, produce HTML

**Key functions:**
- `renderHTML()` - Template execution
- `MinifyHTML()` - HTML minification

**Template context available:**
```go
// In templates, available via .lvt namespace
.lvt.Error("fieldName")     // Validation error for field
.lvt.HasError("fieldName")  // Check if field has error
.lvt.AllErrors()            // All validation errors
.lvt.Uploads("fieldName")   // Access uploads
.lvt.DevMode                // Development mode flag
```

**Common issues:**
- Template execution error (field not found)
- Minification error (fallback to original)
- Wrapper div injection conflicts

---

### Phase 5: Send (`internal/send/`, `internal/session/`)

**What happens:** TreeNode -> JSON -> WebSocket/HTTP

**Key files:**
- `mount.go` - HTTP/WebSocket handler
- `session/registry.go` - Connection management

**WebSocket architecture:**
```
Send() → Queue to sendChan → writePump goroutine → WebSocket.WriteMessage()
```

**Metrics to check:**
```go
// Prometheus metrics available
wsBufferFull       // Buffer overflow events
wsSlowClientCloses // Slow client disconnects
wsWriteErrors      // WebSocket write failures
```

**Common issues:**
- WebSocket buffer full (`ErrClientTooSlow`)
- Connection closed (`ErrConnectionClosed`)
- Write errors logged but not propagated

---

## Client-Side Debugging Guide

### WebSocket Layer (`transport/websocket.ts`)

**Debug variables (Browser Console):**
```javascript
// All received WebSocket messages
window.__wsMessages

// Last message received
window.__lastWSMessage

// Connection state: 0=CONNECTING, 1=OPEN, 2=CLOSING, 3=CLOSED
liveTemplateClient.ws?.readyState

// Transport used for last message
window.__lvtSendPath  // "websocket" | "http" | "http-fallback"
```

**Connection issues:**
```javascript
// Check if WebSocket connected
if (liveTemplateClient.ws?.readyState !== 1) {
    console.log("WebSocket not connected");
}

// Force reconnect (if stuck)
liveTemplateClient.ws?.close();
// Auto-reconnect should trigger
```

---

### TreeRenderer (`state/tree-renderer.ts`)

**Key state:**
- `treeState` - Current tree structure
- `rangeState` - Range item tracking
- `rangeIdKeys` - ID-to-key mapping for ranges

**Debug (Browser Console):**
```javascript
// Get current tree state
liveTemplateClient.getTreeState()

// After an update, check what was applied
// Look for changes in dynamics
```

**Key methods:**
- `applyUpdate()` - Entry point for server updates
- `deepMergeTreeNodes()` - Merge nested updates
- `applyDifferentialOpsToRange()` - Handle range operations
- `reconstructFromTree()` - TreeNode -> HTML string

**Common issues:**
- `rangeState` out of sync with server
- Deep merge not handling nested updates
- HTML reconstruction misalignment

---

### DOM Update (`livetemplate-client.ts`)

**Uses morphdom for efficient DOM reconciliation.**

**Key callbacks:**
- `onBeforeElUpdated` - Focus preservation, skip equal nodes
- `onNodeAdded` - `lvt-mounted` lifecycle hook
- `onBeforeNodeDiscarded` - `lvt-destroyed` lifecycle hook

**Debug variables:**
```javascript
// Check if send was called
window.__lvtSendCalled

// Action name of last message
window.__lvtMessageAction

// Form submission tracking
window.__lvtSubmitListenerTriggered
window.__lvtActionFound
window.__lvtInWrapper
```

**Common issues:**
- Key mismatch causing full replace instead of update
- Focus lost during update (focus restoration failing)
- Lifecycle hooks not firing

---

### Event Delegation (`dom/event-delegation.ts`)

**Supported attributes:**
- `lvt-click`, `lvt-submit`, `lvt-change`, `lvt-input`
- `lvt-keydown`, `lvt-keyup`, `lvt-focus`, `lvt-blur`
- `lvt-window-*` for window-level events
- `lvt-click-away` for outside clicks

**Debug:**
```javascript
// Add before event to trace handling
document.addEventListener('click', (e) => {
    console.log('Click target:', e.target);
    console.log('Has lvt-click:', e.target.hasAttribute?.('lvt-click'));
}, true);
```

---

## Debug Commands & Techniques

### Server-Side Logging

```go
// Add to handler for comprehensive debugging
import (
    "encoding/json"
    "log"
)

// Before tree building
log.Printf("State before: %+v", state)

// After tree building
func debugTree(name string, tree *TreeNode) {
    data, _ := json.MarshalIndent(tree, "", "  ")
    log.Printf("%s tree:\n%s", name, string(data))
}

// Before/after diff
debugTree("Old", oldTree)
debugTree("New", newTree)
log.Printf("Changes: %+v", diff)
```

### Client-Side Debugging (Browser Console)

```javascript
// Comprehensive message logging
window.addEventListener('message', (e) => {
    console.log('Window message:', e.data);
});

// Trace all WebSocket messages
const origOnMessage = liveTemplateClient.ws.onmessage;
liveTemplateClient.ws.onmessage = (e) => {
    console.log('WS message:', JSON.parse(e.data));
    origOnMessage.call(liveTemplateClient.ws, e);
};

// Manual tree state inspection
JSON.stringify(liveTemplateClient.getTreeState(), null, 2);

// Force a re-render (testing)
liveTemplateClient.applyUpdate({"0": "test"});
```

### E2E Test Debugging (chromedp)

```go
// Access browser console logs
chromedp.ActionFunc(func(ctx context.Context) error {
    chromedp.ListenTarget(ctx, func(ev interface{}) {
        if msg, ok := ev.(*runtime.EventConsoleAPICalled); ok {
            for _, arg := range msg.Args {
                log.Printf("Console [%s]: %s", msg.Type, arg.Value)
            }
        }
    })
    return nil
})

// Get WebSocket state
var wsState int
chromedp.Evaluate(`liveTemplateClient.ws?.readyState || -1`, &wsState)

// Get tree state
var treeState map[string]interface{}
chromedp.Evaluate(`liveTemplateClient.getTreeState()`, &treeState)

// Get all WebSocket messages
var wsMessages []interface{}
chromedp.Evaluate(`window.__wsMessages || []`, &wsMessages)
```

---

## Common Root Causes & Fix Patterns

### Tree Building Issues

| Cause | Symptom | Fix |
|-------|---------|-----|
| Statics/dynamics mismatch | Partial HTML, missing content | Check tree building logic alignment |
| Nested range not embedded | Range items inline instead of nested | Check `childTree.HasRange()` branch |
| Key generator reuse | Wrong keys on reconnect | Reset keyGen per session |

### Diff Issues

| Cause | Symptom | Fix |
|-------|---------|-----|
| Same fingerprint | No update sent | Verify data actually changed |
| Registry mistrack | Statics wrongly stripped | Check `MarkSeen()` calls |
| Range match failure | Full replace instead of diff | Check key extraction logic |

### Client Issues

| Cause | Symptom | Fix |
|-------|---------|-----|
| rangeState out of sync | Range operations fail | Initialize from first tree |
| morphdom key conflict | Wrong elements updated | Use consistent `data-key` |
| Tree merge incomplete | Old values persist | Check `deepMergeTreeNodes()` |

---

## Related Skills

- **lvt:troubleshoot** - General debugging for build, migration, template errors
- **lvt:run-and-test** - Running development server and tests
- **lvt:validate-templates** - Template syntax validation
- **lvt:customize** - Understanding generated handler/template structure

### Escalation Path

1. Start with this skill for rendering-specific issues
2. If template syntax issue -> `lvt:validate-templates`
3. If build/compilation error -> `lvt:troubleshoot`
4. If deployment issue -> `lvt:deploy`
5. If test failure -> `lvt:run-and-test`

---

## Source Code References

These are the key files and functions this skill references. If any change significantly, update the skill.

### Server-Side (livetemplate/)

| File | Key Functions/Types |
|------|---------------------|
| `template.go` | `buildTree()`, `ExecuteUpdates()`, `Clone()` |
| `dispatch.go` | `DispatchError`, `DispatchWithState()` |
| `internal/parse/parse.go` | `Parse()`, `evaluatePipe()` |
| `internal/build/types.go` | `TreeNode`, `RangeData`, `TreeMetadata` |
| `internal/diff/tree_compare.go` | `CompareTreesAndGetChangesWithPath()` |
| `internal/diff/range_ops.go` | Range differential operations |
| `internal/observe/metrics.go` | `actionsProcessed`, `treesBuilt`, `wsWriteErrors` |

### Client-Side (client/)

| File | Key Functions/Variables |
|------|-------------------------|
| `livetemplate-client.ts` | `updateDOM()`, `handleWebSocketPayload()`, `window.__wsMessages`, `window.__lvtSendPath` |
| `state/tree-renderer.ts` | `applyUpdate()`, `reconstructFromTree()`, `deepMergeTreeNodes()` |
| `transport/websocket.ts` | `WebSocketTransport`, `WebSocketManager`, `getReadyState()` |
| `dom/event-delegation.ts` | `lvt-*` attribute handlers |

---

## Maintenance

**Last Validated:** 2024-12-20 against commit 696b1fd

**Update Checklist:**

When updating this skill, verify:
- [ ] Pipeline phases still match (Parse -> Build -> Diff -> Render -> Send)
- [ ] Debug variables in client still exist (`__wsMessages`, `__lvtSendPath`, etc.)
- [ ] DispatchError structure unchanged
- [ ] TreeNode type unchanged
- [ ] Run validation test: `go test ./e2e -run TestDebugRenderingSkillReferences`

**Monitored Files:**

Server: `template.go`, `dispatch.go`, `internal/parse/parse.go`, `internal/build/types.go`, `internal/diff/tree_compare.go`

Client: `livetemplate-client.ts`, `state/tree-renderer.ts`, `transport/websocket.ts`

---

## Remember

**Do:**
- Check server logs first for DispatchError
- Use `window.__wsMessages` to verify WebSocket delivery
- Inspect TreeNode structure when debugging partial renders
- Verify method signatures match Controller+State pattern
- Check `rangeState` for range update issues

**Don't:**
- Assume client received updates without checking `__wsMessages`
- Ignore fingerprint unchanged errors
- Skip checking method signature when dispatch fails
- Forget to check both server AND client for rendering issues
