# LVT Deterministic Generation: Implementation Plan

**Goal:** Transform lvt from a fragile code generator into a self-improving system that reliably produces working LiveTemplate apps.

**Target Audience:** LLM agents implementing these issues

---

## Overview

```
Milestone 1: Stop the Bleeding (Foundation)
     ↓
Milestone 2: Validation Layer (Feedback Infrastructure)
     ↓
Milestone 3: Telemetry & Evolution (Self-Improvement)
     ↓
Milestone 4: Components Integration (Monorepo - Template Consolidation)
     ↓
Milestone 5: Style System (CSS Swappability)
     ↓
Milestone 6: Components Evolution (Unified Feedback Loop - Single Repo)
```

---

## Milestone 1: Stop the Bleeding

**Goal:** Fix known issues and add basic safeguards to prevent shipping broken code.

**Duration:** 1 week

### Issue 1.1: Add Compilation Validation to E2E Tests

**Priority:** Critical
**Labels:** `testing`, `quick-win`

**Description:**
Currently, E2E tests verify that files are generated but never check if the generated code compiles. This allows syntax errors and type mismatches to ship.

**Tasks:**
- [ ] Add `go build ./...` step after each resource generation in E2E tests
- [ ] Add `go mod tidy` step before build (catches dependency issues)
- [ ] Ensure build errors fail the test with clear error messages
- [ ] Update all existing E2E tests to include compilation check

**Acceptance Criteria:**
- All E2E tests run `go build` on generated apps
- A test that generates invalid Go code fails with compilation error
- CI catches compilation issues before merge

**Files to Modify:**
- `e2e/resource_test.go`
- `e2e/app_test.go`
- `e2e/auth_test.go`
- Create `e2e/helpers/validation.go` for shared validation logic

**Technical Notes:**
```go
// Example validation helper
func ValidateGeneratedApp(t *testing.T, appPath string) {
    // Run go mod tidy
    cmd := exec.Command("go", "mod", "tidy")
    cmd.Dir = appPath
    if out, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("go mod tidy failed: %s\n%s", err, out)
    }

    // Run go build
    cmd = exec.Command("go", "build", "-o", "/dev/null", "./...")
    cmd.Dir = appPath
    cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
    if out, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("generated code does not compile: %s\n%s", err, out)
    }
}
```

---

### Issue 1.2: Remove SKIP_GO_MOD_TIDY from Agent Tests

**Priority:** High
**Labels:** `testing`, `quick-win`
**Depends On:** None

**Description:**
Agent tests set `SKIP_GO_MOD_TIDY=1` environment variable, which masks dependency resolution issues. This should be removed so tests catch real-world problems.

**Tasks:**
- [ ] Find all occurrences of `SKIP_GO_MOD_TIDY` in test code
- [ ] Remove the environment variable setting
- [ ] Fix any tests that fail due to dependency issues (these are real bugs)
- [ ] Document why this was removed (prevent re-introduction)

**Acceptance Criteria:**
- No `SKIP_GO_MOD_TIDY` in codebase
- All agent tests pass with `go mod tidy` running
- Any dependency issues discovered are fixed

**Files to Modify:**
- `agenttest/harness.go`
- Any other files setting this env var

---

### Issue 1.3: Fix EditingID Type Inconsistency

**Priority:** High
**Labels:** `bug`, `templates`, `quick-win`
**Depends On:** None

**Description:**
The `EditingID` field is compared as integer in some templates and string in others, causing type errors.

**Evidence from git history:**
- Commit `5700a82`: "single kit template bug" - EditingID comparison issue

**Tasks:**
- [ ] Audit all templates for EditingID usage
- [ ] Standardize on string type (empty string = no editing)
- [ ] Update all comparisons: `{{if ne .EditingID 0}}` → `{{if ne .EditingID ""}}`
- [ ] Update handler generation to use string type
- [ ] Add test case that catches this regression

**Acceptance Criteria:**
- All kits use consistent EditingID type (string)
- Generated handlers declare EditingID as string
- Test exists that fails if types mismatch

**Files to Modify:**
- `internal/kits/system/multi/templates/resource/template.tmpl.tmpl`
- `internal/kits/system/single/templates/resource/template.tmpl.tmpl`
- `internal/generator/templates/resource/handler.go.tmpl`

---

### Issue 1.4: Add Template Parse Validation to Generation

**Priority:** High
**Labels:** `validation`, `generator`
**Depends On:** None

**Description:**
After generating template files, validate they parse correctly before reporting success. Currently, template syntax errors are only discovered at runtime.

**Tasks:**
- [ ] Create template validation function in generator package
- [ ] Call validation after writing each `.tmpl` file
- [ ] Return clear error if template doesn't parse
- [ ] Include line number and context in error message

**Acceptance Criteria:**
- Invalid template syntax causes generation to fail
- Error message includes file path, line number, and syntax issue
- User can immediately see what's wrong

**Files to Create/Modify:**
- Create `internal/generator/validate.go`
- Modify `internal/generator/resource.go` to call validation
- Modify `internal/generator/view.go` to call validation

**Technical Notes:**
```go
func ValidateTemplate(path string) error {
    content, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    _, err = template.New(filepath.Base(path)).Parse(string(content))
    if err != nil {
        return fmt.Errorf("template %s: %w", path, err)
    }
    return nil
}
```

---

### Issue 1.5: Ensure All Kits Have Feature Parity

**Priority:** Medium
**Labels:** `templates`, `kits`
**Depends On:** 1.3

**Description:**
Different kits have different features. The single kit is missing delete functionality that exists in multi kit. This causes confusion and inconsistent behavior.

**Tasks:**
- [ ] Audit all three kits (multi, single, simple) for feature differences
- [ ] Document which features each kit should have
- [ ] Add missing delete button to single kit edit modal
- [ ] Ensure cancel button works consistently across all kits
- [ ] Add tests that verify feature parity

**Acceptance Criteria:**
- Single kit has delete functionality in edit modal
- All kits have consistent button behavior
- Feature matrix documented in kit README

**Files to Modify:**
- `internal/kits/system/single/templates/resource/template.tmpl.tmpl`
- Create `internal/kits/README.md` with feature matrix

---

## Milestone 2: Validation Layer

**Goal:** Build infrastructure that validates generated code and provides feedback to the generation system.

**Duration:** 2 weeks

### Issue 2.1: Create Validation Engine Package

**Priority:** Critical
**Labels:** `infrastructure`, `validation`
**Depends On:** Milestone 1 complete

**Description:**
Create a dedicated package for validating generated applications. This will be used by both the generator and the evolution system.

**Tasks:**
- [ ] Create `internal/validation/` package
- [ ] Implement `ValidationEngine` struct with configurable checks
- [ ] Implement `GoModValidation` - checks go.mod is valid
- [ ] Implement `CompilationValidation` - runs go build
- [ ] Implement `TemplateValidation` - parses all .tmpl files
- [ ] Implement `MigrationValidation` - validates SQL syntax
- [ ] Create `ValidationResult` type with structured errors
- [ ] Add timeout support for long-running validations

**Acceptance Criteria:**
- `validation.Validate(appPath)` returns comprehensive result
- Each validation type can be run independently
- Results include specific file/line information for errors
- Validation respects configurable timeout

**Files to Create:**
```
internal/validation/
├── validation.go      # Main engine and types
├── go_mod.go          # go.mod validation
├── compilation.go     # go build validation
├── templates.go       # template parse validation
├── migrations.go      # SQL validation
└── validation_test.go # Tests
```

**Technical Notes:**
```go
type ValidationEngine struct {
    Timeout time.Duration
}

type ValidationResult struct {
    Passed     bool
    Checks     map[string]CheckResult
    Duration   time.Duration
}

type CheckResult struct {
    Name     string
    Passed   bool
    Errors   []ValidationError
    Warnings []string
}

type ValidationError struct {
    File    string
    Line    int
    Column  int
    Message string
    Context string // surrounding code
}
```

---

### Issue 2.2: Integrate Validation into Generator

**Priority:** Critical
**Labels:** `generator`, `validation`
**Depends On:** 2.1

**Description:**
Call the validation engine after generation and before reporting success. This ensures users get immediate feedback when something is wrong.

**Tasks:**
- [ ] Add validation call at end of `GenerateResource`
- [ ] Add validation call at end of `GenerateView`
- [ ] Add validation call at end of `GenerateAuth`
- [ ] Make validation configurable (can be disabled for speed)
- [ ] Include validation results in command output
- [ ] Return non-zero exit code if validation fails

**Acceptance Criteria:**
- `lvt gen resource` validates after generation
- Validation errors are clearly displayed
- User can skip validation with `--skip-validation` flag
- Exit code is non-zero when validation fails

**Files to Modify:**
- `internal/generator/resource.go`
- `internal/generator/view.go`
- `internal/generator/auth.go`
- `commands/gen.go` (add flag)

---

### Issue 2.3: Add Validation to MCP Tools

**Priority:** High
**Labels:** `mcp`, `validation`
**Depends On:** 2.1

**Description:**
MCP tools should return validation results so LLMs know immediately if generation succeeded or failed. This is critical for the evolution feedback loop.

**Tasks:**
- [ ] Update `lvt_gen_resource` tool to include validation in response
- [ ] Update `lvt_gen_view` tool to include validation in response
- [ ] Update `lvt_gen_auth` tool to include validation in response
- [ ] Format validation errors clearly in tool response
- [ ] Include suggestions for common errors

**Acceptance Criteria:**
- MCP tool responses include validation status
- Validation errors are formatted for LLM consumption
- LLM can understand what failed and potentially fix it

**Files to Modify:**
- `commands/mcp_handlers.go` (or equivalent MCP handler file)

**Example Response Format:**
```json
{
  "success": false,
  "validation": {
    "compilation": {
      "passed": false,
      "errors": [
        {
          "file": "handlers/posts.go",
          "line": 45,
          "message": "undefined: models.Post",
          "suggestion": "Run migrations or check model generation"
        }
      ]
    },
    "templates": {
      "passed": true
    }
  }
}
```

---

### Issue 2.4: Create Runtime Validation (App Startup Test)

**Priority:** Medium
**Labels:** `validation`, `testing`
**Depends On:** 2.1

**Description:**
Beyond compilation, verify the generated app can actually start. This catches runtime errors like missing environment variables or database connection issues.

**Tasks:**
- [ ] Add `RuntimeValidation` to validation engine
- [ ] Start the app in a subprocess
- [ ] Wait for HTTP server to be ready (health check)
- [ ] Verify basic routes respond
- [ ] Gracefully shutdown after test
- [ ] Timeout if app doesn't start within threshold

**Acceptance Criteria:**
- Validation can optionally test app startup
- Startup failures are reported with error details
- Test doesn't leave orphan processes
- Works with SQLite (no external DB needed)

**Files to Modify:**
- Add `internal/validation/runtime.go`
- Update `internal/validation/validation.go`

---

## Milestone 3: Telemetry & Evolution

**Goal:** Capture generation events and build the foundation for self-improvement, including ability to propose fixes to upstream LiveTemplate ecosystem repos.

**Duration:** 2 weeks
**Issues:** 7 (3.1-3.7)

### Issue 3.1: Create Telemetry Package

**Priority:** Critical
**Labels:** `infrastructure`, `evolution`
**Depends On:** Milestone 2 complete

**Description:**
Create a package to capture and store generation events. This data drives the evolution system.

**Tasks:**
- [ ] Create `internal/telemetry/` package
- [ ] Define `GenerationEvent` struct with all relevant fields
- [ ] Create SQLite storage backend
- [ ] Implement event capture with minimal overhead
- [ ] Add privacy controls (can be disabled)
- [ ] Create retention policy (auto-delete old events)

**Acceptance Criteria:**
- Every generation attempt is captured (when enabled)
- Events include: command, inputs, kit, success/failure, errors, duration
- Data stored in `~/.config/lvt/telemetry.db`
- User can disable with `LVT_TELEMETRY=false`

**Files to Create:**
```
internal/telemetry/
├── telemetry.go       # Main API
├── events.go          # Event types
├── store.go           # Storage interface
├── sqlite.go          # SQLite implementation
├── schema.sql         # Database schema
└── telemetry_test.go
```

**Schema:**
```sql
CREATE TABLE generation_events (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    command TEXT NOT NULL,
    inputs TEXT NOT NULL,  -- JSON
    kit TEXT,
    lvt_version TEXT,
    success BOOLEAN NOT NULL,
    validation_result TEXT,  -- JSON
    errors TEXT,  -- JSON
    duration_ms INTEGER,
    files_generated TEXT  -- JSON
);
```

---

### Issue 3.2: Integrate Telemetry into Generator

**Priority:** High
**Labels:** `generator`, `telemetry`
**Depends On:** 3.1

**Description:**
Hook telemetry capture into the generation flow so every attempt is recorded.

**Tasks:**
- [ ] Start capture at beginning of generation
- [ ] Record all inputs and configuration
- [ ] Record validation results
- [ ] Record any errors with context
- [ ] Finalize capture at end (success or failure)
- [ ] Ensure capture doesn't slow down generation noticeably

**Acceptance Criteria:**
- All `lvt gen` commands record telemetry
- Failed generations include error details
- Capture adds <10ms overhead
- Works correctly even if generation crashes

**Files to Modify:**
- `internal/generator/resource.go`
- `internal/generator/view.go`
- `internal/generator/auth.go`

---

### Issue 3.3: Create Knowledge Base for Common Errors

**Priority:** High
**Labels:** `evolution`, `knowledge-base`
**Depends On:** 3.1

**Description:**
Create a knowledge base that maps error patterns to known fixes. Patterns are stored in a **markdown file** for git tracking, easy reading, and LLM editability. This is the foundation for automated fix proposals.

**Tasks:**
- [ ] Create `internal/evolution/knowledge/` package
- [ ] Create `evolution/patterns.md` - the source of truth for all patterns
- [ ] Define `Pattern` struct (matcher + fixes)
- [ ] Implement markdown parser to load patterns from file
- [ ] Implement pattern matching against errors
- [ ] Seed with known patterns from git history:
  - EditingID type mismatch
  - Modal state persistence
  - Form sync issues
  - Session not cleared
  - Import path errors
- [ ] Add CLI command: `lvt evolution patterns` to list all patterns

**Acceptance Criteria:**
- Patterns stored in `evolution/patterns.md` (git-tracked)
- At least 10 patterns seeded from known issues
- Pattern matching is fast (<1ms per error)
- Patterns include confidence scores
- LLMs can propose pattern additions via PR to markdown file
- `lvt evolution patterns` lists all patterns with stats

**Files to Create:**
```
evolution/
└── patterns.md            # Source of truth - git tracked

internal/evolution/knowledge/
├── knowledge.go           # Knowledge base API
├── parser.go              # Markdown parser for patterns.md
├── patterns.go            # Pattern types
├── matcher.go             # Pattern matching logic
└── knowledge_test.go
```

**Patterns Markdown Format (`evolution/patterns.md`):**
```markdown
# Evolution Patterns

Known error patterns and their fixes. This file is the source of truth
for the evolution system's knowledge base.

## Pattern: editing-id-type

**Name:** EditingID Type Mismatch
**Confidence:** 0.95
**Added:** 2026-01-15
**Fix Count:** 0
**Success Rate:** -

### Description
EditingID compared as integer but is string type in handler.

### Error Pattern
- **Phase:** compilation
- **Message Regex:** `cannot convert .* to type int.*EditingID`
- **Context Regex:** `EditingID`

### Fix
- **File:** `*/template.tmpl.tmpl`
- **Find:** `{{if ne .EditingID 0}}`
- **Replace:** `{{if ne .EditingID ""}}`

---

## Pattern: modal-state-persistence

**Name:** Modal State Persists After Close
**Confidence:** 0.90
**Added:** 2026-01-15
**Fix Count:** 0
**Success Rate:** -

### Description
Modal editing state (IsAdding, EditingID) persists on page reload
because fields are not marked as transient.

### Error Pattern
- **Phase:** runtime
- **Message Regex:** `modal (open|visible) after (reload|refresh)`

### Fix
- **File:** `*/handler.go.tmpl`
- **Find:** `EditingID string`
- **Replace:** `EditingID string \`lvt:"transient"\``

---

## Pattern: form-sync-morphdom

**Name:** Form Values Revert After Update
**Confidence:** 0.88
**Added:** 2026-01-15
**Fix Count:** 0
**Success Rate:** -

### Description
Select dropdown values revert to previous state after morphdom
DOM patching because expected value not preserved.

### Error Pattern
- **Phase:** runtime
- **Message Regex:** `(select|dropdown) value (reverted|reset|changed)`

### Fix
- **File:** `*/components/form.tmpl`
- **Find:** `<select`
- **Replace:** `<select data-expected-value="{{.Value}}"`

---

<!-- Add new patterns above this line -->
```

**Benefits of Markdown Storage:**
1. **Git tracked** - full history of pattern changes
2. **Human readable** - easy to review and understand
3. **LLM editable** - evolution system can propose pattern additions
4. **PR reviewable** - pattern changes go through normal review
5. **No recompile** - add patterns without rebuilding lvt

---

### Issue 3.4: Create Fix Proposer

**Priority:** High
**Labels:** `evolution`
**Depends On:** 3.3

**Description:**
Create a component that analyzes generation failures and proposes fixes using the knowledge base.

**Tasks:**
- [ ] Create `internal/evolution/proposer.go`
- [ ] Implement `ProposeFixesFor(event GenerationEvent) []Fix`
- [ ] Query knowledge base for matching patterns
- [ ] Rank fixes by confidence
- [ ] Include rationale in each proposed fix
- [ ] Support multiple fixes for single error

**Acceptance Criteria:**
- Known error patterns get fix proposals
- Fixes include file, change, and rationale
- Multiple fixes can be proposed for one failure
- Confidence scores help prioritize

**Files to Create:**
- `internal/evolution/proposer.go`
- `internal/evolution/types.go` (Fix, Proposal types)

---

### Issue 3.5: Create Fix Tester

**Priority:** Medium
**Labels:** `evolution`, `testing`
**Depends On:** 3.4, 2.1

**Description:**
Before applying a fix, test it in isolation to ensure it actually works.

**Tasks:**
- [ ] Create `internal/evolution/tester.go`
- [ ] Implement isolated test environment (temp directory)
- [ ] Apply proposed fix to templates
- [ ] Re-run the failed generation
- [ ] Validate the result
- [ ] Clean up test environment
- [ ] Report test results

**Acceptance Criteria:**
- Fixes are tested before application
- Test uses same inputs as original failure
- Test environment is isolated (no side effects)
- Test results include validation details

**Files to Create:**
- `internal/evolution/tester.go`
- `internal/evolution/tester_test.go`

---

### Issue 3.6: Create Evolution CLI Commands

**Priority:** Medium
**Labels:** `cli`, `evolution`
**Depends On:** 3.1, 3.3, 3.4

**Description:**
Add CLI commands to interact with the evolution system.

**Tasks:**
- [ ] Add `lvt evolution status` - show system status
- [ ] Add `lvt evolution metrics` - show success rates
- [ ] Add `lvt evolution failures` - list recent failures
- [ ] Add `lvt evolution propose <event-id>` - propose fixes for failure
- [ ] Add `lvt evolution apply <fix-id>` - apply a fix
- [ ] Add `lvt evolution learn <pattern-file>` - add new pattern

**Acceptance Criteria:**
- All commands work and have help text
- Status shows meaningful metrics
- Failures are listed with enough context
- Apply actually modifies templates

**Files to Create/Modify:**
- Create `commands/evolution.go`
- Modify `main.go` to register commands

---

### Issue 3.7: Add Upstream Library Evolution Support

**Priority:** Medium
**Labels:** `evolution`, `upstream`
**Depends On:** 3.2, 3.4, 3.5

**Description:**
Extend the evolution system to propose fixes not just to lvt templates, but also to upstream LiveTemplate ecosystem repos. Some bugs originate in the core library or client code, and the evolution system should be able to track and propose fixes there.

**Evidence from git history:**
- Multiple session-related fixes trace to livetemplate/livetemplate
- morphdom sync issues trace to livetemplate/client
- These patterns are documented in `evolution/patterns.md` under "Upstream Patterns"

**Tasks:**
- [ ] Add `UpstreamRepo` field parsing to pattern markdown parser
- [ ] Create `internal/evolution/upstream.go` with upstream fix logic
- [ ] Implement upstream repo cloning/updating for fix testing
- [ ] Implement PR creation to upstream repos (requires GitHub token)
- [ ] Add `lvt evolution upstream-status` command to track pending upstream PRs
- [ ] Add tracking for when upstream fixes are merged and released
- [ ] Auto-update go.mod when upstream releases include our fixes

**Acceptance Criteria:**
- Patterns with `UpstreamRepo` field are correctly parsed
- Upstream fixes can be tested against local clone of upstream repo
- PRs can be created to upstream repos with proper evidence
- CLI shows status of pending upstream PRs
- go.mod is updated when upstream fixes are released

**Files to Create/Modify:**
- Create `internal/evolution/upstream.go`
- Modify `internal/evolution/knowledge/parser.go` to handle UpstreamRepo
- Modify `commands/evolution.go` to add upstream commands

**Technical Notes:**
```go
// Pattern with upstream repo field
pattern := &Pattern{
    ID:           "morphdom-select-sync",
    UpstreamRepo: "github.com/livetemplate/client",
    Fix: Fix{
        File: "src/morphdom-config.js",
        // ...
    },
}

// Upstream fix workflow
proposer := evolution.NewUpstreamProposer(git, gh, knowledge)
fix, _ := proposer.ProposeUpstreamFix(pattern, event)
pr, _ := proposer.CreateUpstreamPR(fix)
// Track PR status and handle merge
```

---

## Milestone 4: Components Integration

**Goal:** Integrate the components library to eliminate template drift.

**Duration:** 2 weeks

### Issue 4.1: Move Components into lvt Monorepo

**Priority:** High
**Labels:** `components`, `architecture`, `monorepo`
**Depends On:** Milestone 1 complete

**Description:**
Move the components library into lvt as a nested Go module. This enables atomic changes, faster iteration, and a single feedback loop for the evolution system while maintaining independent importability.

**Tasks:**
- [ ] Use `git subtree add` to import existing components repo into `components/`
- [ ] Create `components/go.mod` with module path `github.com/livetemplate/lvt/components`
- [ ] Ensure components/go.mod only depends on livetemplate, NOT on lvt
- [ ] Update lvt's go.mod to reference local components (for development)
- [ ] Add CI workflow to verify components independence (no lvt imports)
- [ ] Create helper functions for component imports in generated code
- [ ] Update handler templates to support component state fields
- [ ] Document component usage patterns

**Acceptance Criteria:**
- Components live in `components/` directory with own go.mod
- External apps can import `github.com/livetemplate/lvt/components/modal`
- Components build and test independently (`cd components && go test ./...`)
- CI fails if components import anything from lvt
- Generated apps can import and use components

**Files to Create/Modify:**
- Create `components/go.mod`
- Create `components/modal/`, `components/toast/`, etc. (from subtree)
- Create `.github/workflows/components-independence.yml`
- Modify `go.mod` (add replace directive for local dev)
- Modify `internal/generator/templates/resource/handler.go.tmpl`

**Technical Notes:**
```bash
# Import existing components
git subtree add --prefix=components \
    git@github.com:livetemplate/components.git main --squash

# components/go.mod
module github.com/livetemplate/lvt/components

go 1.25

require github.com/livetemplate/livetemplate v0.8.0
# NO lvt dependency allowed here
```

```yaml
# .github/workflows/components-independence.yml
- name: Verify no lvt imports
  run: |
    cd components
    if grep -r '"github.com/livetemplate/lvt"' --include="*.go" .; then
      echo "ERROR: components must not import lvt"
      exit 1
    fi
    if grep -r '"github.com/livetemplate/lvt/internal' --include="*.go" .; then
      echo "ERROR: components must not import lvt internals"
      exit 1
    fi
- name: Build standalone
  run: cd components && go build ./...
- name: Test standalone
  run: cd components && go test ./...
```

---

### Issue 4.2: Integrate Modal Component

**Priority:** High
**Labels:** `components`, `modal`
**Depends On:** 4.1

**Description:**
Replace kit modal templates with the modal component from the components library. This eliminates the most common source of bugs (40% of UI issues).

**Tasks:**
- [ ] Update handler template to use `modal.ConfirmState` for delete confirmation
- [ ] Update handler template to use `modal.State` for edit modal (if applicable)
- [ ] Update template generation to use `{{template "lvt:modal:confirm:v1" .DeleteConfirm}}`
- [ ] Register component templates in generated main.go
- [ ] Update all three kits to use component
- [ ] Remove old modal templates from kits

**Acceptance Criteria:**
- Generated apps use modal component
- Delete confirmation works correctly
- Edit modal works correctly (if using modal)
- Old modal templates removed from kits

**Files to Modify:**
- `internal/generator/templates/resource/handler.go.tmpl`
- `internal/generator/templates/app/main.go.tmpl`
- `internal/kits/system/multi/templates/resource/template.tmpl.tmpl`
- `internal/kits/system/single/templates/resource/template.tmpl.tmpl`

---

### Issue 4.3: Integrate Toast Component

**Priority:** Medium
**Labels:** `components`, `toast`
**Depends On:** 4.1

**Description:**
Add toast notifications to generated apps using the toast component. This provides standardized user feedback.

**Tasks:**
- [ ] Add `toast.ContainerState` to generated handler state
- [ ] Add toast container to layout template
- [ ] Add success toasts after create/update/delete actions
- [ ] Add error toasts on failures
- [ ] Register toast templates in main.go

**Acceptance Criteria:**
- Generated apps show toast on CRUD operations
- Toasts auto-dismiss after timeout
- Toast position is configurable

**Files to Modify:**
- `internal/generator/templates/resource/handler.go.tmpl`
- `internal/kits/system/*/templates/components/layout.tmpl`

---

### Issue 4.4: Integrate Dropdown Component

**Priority:** Medium
**Labels:** `components`, `dropdown`
**Depends On:** 4.1

**Description:**
Replace select elements with dropdown component to fix morphdom sync issues (30% of form bugs).

**Tasks:**
- [ ] Detect select fields in resource generation
- [ ] Use dropdown component for select fields
- [ ] Support both single and multi-select
- [ ] Ensure value sync works with LiveTemplate updates

**Acceptance Criteria:**
- Select fields use dropdown component
- Values persist correctly across updates
- Searchable dropdown works for long lists

**Files to Modify:**
- `internal/generator/templates/resource/handler.go.tmpl`
- `internal/generator/form.go` (form field generation)
- `internal/kits/system/*/templates/components/form.tmpl`

---

### Issue 4.5: Add Component Usage Detection

**Priority:** Low
**Labels:** `components`, `generator`
**Depends On:** 4.2, 4.3, 4.4

**Description:**
Automatically detect which components a generated app uses and only import/register what's needed.

**Tasks:**
- [ ] Analyze generated templates for component usage
- [ ] Generate minimal import list
- [ ] Generate minimal template registration
- [ ] Avoid importing unused components

**Acceptance Criteria:**
- Only used components are imported
- Generated go.mod only includes necessary dependencies
- No dead code in generated apps

---

## Milestone 5: Style System

**Goal:** Implement CSS swappability through style adapters.

**Duration:** 2 weeks

### Issue 5.1: Create Style Adapter Interface

**Priority:** High
**Labels:** `styles`, `architecture`
**Depends On:** Milestone 4 started

**Description:**
Define the style adapter interface that allows swapping CSS frameworks.

**Tasks:**
- [ ] Create `components/styles/` package in lvt monorepo
- [ ] Define `StyleAdapter` interface
- [ ] Define style structs: `ButtonStyles`, `ModalStyles`, `FormStyles`, etc.
- [ ] Create adapter registration mechanism
- [ ] Document how to create custom adapters

**Acceptance Criteria:**
- Interface is well-documented
- All common UI elements have style definitions
- Registration mechanism works

**Files to Create (in lvt monorepo):**
```
components/styles/
├── adapter.go         # Interface definition
├── types.go           # Style struct types
├── registry.go        # Adapter registration
└── styles_test.go
```

---

### Issue 5.2: Implement Tailwind Adapter

**Priority:** High
**Labels:** `styles`, `tailwind`
**Depends On:** 5.1

**Description:**
Create the default Tailwind adapter with all required styles.

**Tasks:**
- [ ] Implement all methods of StyleAdapter interface
- [ ] Use Tailwind best practices for each component
- [ ] Include responsive variants
- [ ] Include dark mode support
- [ ] Test with all components

**Acceptance Criteria:**
- All components render correctly with Tailwind
- Responsive design works
- Dark mode works (if enabled)

**Files to Create:**
```
components/styles/tailwind/
├── adapter.go
└── adapter_test.go
```

---

### Issue 5.3: Implement Unstyled Adapter

**Priority:** High
**Labels:** `styles`, `unstyled`
**Depends On:** 5.1

**Description:**
Create an unstyled adapter that outputs semantic class names for custom CSS.

**Tasks:**
- [ ] Implement all methods with semantic class names
- [ ] Follow BEM or similar naming convention
- [ ] Document all class names
- [ ] Create CSS scaffold generator

**Acceptance Criteria:**
- All components render with semantic classes
- Class names are documented
- `lvt styles scaffold` generates CSS template

**Files to Create:**
```
components/styles/unstyled/
├── adapter.go
├── adapter_test.go
└── scaffold.css.tmpl
```

---

### Issue 5.4: Integrate Style Adapters into Templates

**Priority:** High
**Labels:** `styles`, `templates`
**Depends On:** 5.2, 5.3

**Description:**
Update all component and kit templates to use style adapter references instead of hardcoded classes.

**Tasks:**
- [ ] Update component templates to use `{{.Styles.X.Y}}`
- [ ] Update kit templates to use same pattern
- [ ] Ensure styles are passed through template context
- [ ] Test with both Tailwind and Unstyled adapters

**Acceptance Criteria:**
- No hardcoded CSS classes in templates
- Same template works with any adapter
- All components render correctly with both adapters

---

### Issue 5.5: Add Style Selection to Generation

**Priority:** Medium
**Labels:** `styles`, `cli`
**Depends On:** 5.4

**Description:**
Allow users to select style adapter at generation time and change it later.

**Tasks:**
- [ ] Add `--styles` flag to `lvt new`
- [ ] Store style choice in `.lvtrc`
- [ ] Add `lvt styles set <adapter>` command
- [ ] Add `lvt styles eject` to dump adapter for customization
- [ ] Add `lvt styles scaffold` to generate CSS template

**Acceptance Criteria:**
- Can choose style at generation: `lvt new myapp --styles=unstyled`
- Can change style later: `lvt styles set tailwind`
- Can eject for customization
- Can generate CSS scaffold

---

## Milestone 6: Components Evolution

**Goal:** Extend evolution system to improve components alongside lvt kits (all in one repo).

**Duration:** 2 weeks

**Note:** With the monorepo approach (components inside lvt), this milestone is simpler than originally planned. No cross-repo coordination needed - all fixes are in single PRs.

### Issue 6.1: Add Component Attribution to Telemetry

**Priority:** High
**Labels:** `telemetry`, `components`, `evolution`
**Depends On:** 3.1, 4.2

**Description:**
Extend telemetry to track which components are involved in failures.

**Tasks:**
- [ ] Add `ComponentsUsed` field to GenerationEvent
- [ ] Add `ComponentErrors` field for component-specific issues
- [ ] Implement component attribution logic
- [ ] Update storage schema

**Acceptance Criteria:**
- Telemetry records which components were used
- Errors are attributed to specific components when possible
- Can query failures by component

---

### Issue 6.2: Create Component Health Dashboard

**Priority:** Medium
**Labels:** `evolution`, `components`, `cli`
**Depends On:** 6.1

**Description:**
Add CLI command to show component health metrics.

**Tasks:**
- [ ] Create `lvt evolution components` command
- [ ] Show usage count per component
- [ ] Show success rate per component
- [ ] Show common errors per component
- [ ] Highlight components needing attention

**Acceptance Criteria:**
- Can see which components are problematic
- Data drives prioritization of fixes
- Clear visualization of component health

**Example Output:**
```
Component Health (last 30 days)
───────────────────────────────────────────────
Component   │ Usage │ Success │ Common Errors
────────────┼───────┼─────────┼───────────────
modal       │  847  │  94.2%  │ State sync
toast       │  623  │  98.1%  │ -
dropdown    │  412  │  87.3%  │ Selection
toggle      │  534  │  99.2%  │ -
───────────────────────────────────────────────

⚠️  dropdown needs attention (87.3% success rate)
```

---

### Issue 6.3: Implement Component-Aware Fix Proposals

**Priority:** Medium
**Labels:** `evolution`, `components`
**Depends On:** 3.4, 6.1

**Description:**
Extend fix proposer to detect whether errors are in components or kit templates and generate appropriate fixes. With the monorepo, all fixes are in the same PR.

**Tasks:**
- [ ] Detect when error is in `components/` vs `internal/kits/`
- [ ] Generate fix with correct file path within monorepo
- [ ] For component bugs: propose fix in `components/modal/` etc.
- [ ] For kit bugs using components: propose fix in `internal/kits/`
- [ ] Single PR can include both component fix and kit update (atomic)

**Acceptance Criteria:**
- Fix proposer correctly identifies component vs kit errors
- Fixes target correct paths within lvt monorepo
- Single PR can fix component and update kit usage together

---

### Issue 6.4: Verify Component Independence in CI

**Priority:** Medium
**Labels:** `ci`, `components`
**Depends On:** 4.1

**Description:**
Ensure CI verifies that `components/` remains independently importable - it must not import anything from the parent lvt module. This is created in Issue 4.1 but should be enhanced here.

**Tasks:**
- [ ] Enhance the independence check from Issue 4.1
- [ ] Add test that imports components from a separate test module
- [ ] Verify components work without any lvt code
- [ ] Add badge showing components independence status
- [ ] Document the independence guarantee for external users

**Acceptance Criteria:**
- CI prevents lvt imports in components
- Separate test module can import and use components
- Clear documentation for external users
- Badge in README showing independence status

**Files to Modify:**
```yaml
# .github/workflows/ci.yml (add job)
components-independence:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Check for lvt imports
      run: |
        cd components
        if grep -r '"github.com/livetemplate/lvt"' --include="*.go" .; then
          echo "ERROR: components must not import lvt"
          exit 1
        fi
        if grep -r '"github.com/livetemplate/lvt/internal' --include="*.go" .; then
          echo "ERROR: components must not import lvt internals"
          exit 1
        fi
    - name: Build standalone
      run: cd components && go build ./...
    - name: Test standalone
      run: cd components && go test ./...
    - name: Test external import
      run: |
        mkdir -p /tmp/test-import
        cd /tmp/test-import
        go mod init test-import
        echo 'package main
        import "github.com/livetemplate/lvt/components/modal"
        func main() { _ = modal.New("test") }' > main.go
        go mod tidy
        go build .
```

---

## Summary

| Milestone | Issues | Duration | Key Outcome |
|-----------|--------|----------|-------------|
| 1. Stop the Bleeding | 5 | 1 week | No more shipping broken code |
| 2. Validation Layer | 4 | 2 weeks | Immediate feedback on failures |
| 3. Telemetry & Evolution | 6 | 2 weeks | Self-improvement foundation |
| 4. Components Integration | 5 | 2 weeks | Template drift eliminated |
| 5. Style System | 5 | 2 weeks | CSS swappability |
| 6. Components Evolution | 4 | 2 weeks | Unified feedback loop (monorepo) |

**Total: 29 issues across 6 milestones**

## Dependencies Graph

```
Milestone 1 (Foundation)
    │
    ├──▶ Milestone 2 (Validation)
    │        │
    │        └──▶ Milestone 3 (Evolution)
    │                 │
    │                 └──▶ Milestone 6 (Components Evolution)
    │
    └──▶ Milestone 4 (Components Monorepo)
             │
             └──▶ Milestone 5 (Styles)
                      │
                      └──▶ Milestone 6 (Components Evolution)
```

Milestones 1-3 and 1-4-5 can proceed in parallel after Milestone 1 is complete.

**Monorepo Benefit:** With components inside lvt, Milestone 6 is much simpler - no cross-repo coordination needed. All fixes are atomic within a single repository.
