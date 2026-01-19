# Deterministic LiveTemplate Generation: Analysis and Path Forward

**Date**: January 2026
**Status**: Deep Investigation Complete

## Executive Summary

After extensive investigation of the lvt codebase, git history, testing infrastructure, and upstream dependencies, I've identified the root causes of non-deterministic app generation and propose a comprehensive evolution mechanism for continuous LLM-driven improvement.

**Key Findings:**
1. **77% of commits are fixes** - systemic issues with template drift and state management
2. **Template duplication** - 3+ copies of the same template maintained separately
3. **Testing gaps** - generated code is never compiled or run in tests
4. **Upstream volatility** - LiveTemplate (v0.8.0) is in alpha, requiring frequent template updates
5. **Skill probabilism** - LLM behavior requires stronger constraints than currently provided

**Core Proposal**: Implement a **self-healing generation pipeline** where LLMs can:
1. Detect broken generation through automated validation
2. Propose fixes to templates, skills, and tests
3. Submit improvements through a structured review process
4. Measure and track quality over time

---

## Part 1: Current State Analysis

### 1.1 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     lvt CLI (main.go)                           │
├─────────────────────────────────────────────────────────────────┤
│  Commands Layer (commands/)                                     │
│  - new.go, gen.go, kits.go, mcp_server.go                      │
├─────────────────────────────────────────────────────────────────┤
│  Generation Layer (internal/generator/)                         │
│  - App scaffolding, Resource CRUD, Views, Auth                  │
├─────────────────────────────────────────────────────────────────┤
│  Kit System (internal/kits/)                                    │
│  - Template loading, CSS helpers, Component assembly            │
├─────────────────────────────────────────────────────────────────┤
│  Template Sources (Cascade Priority)                            │
│  1. Project (.lvt/kits/)                                        │
│  2. User (~/.config/lvt/kits/)                                  │
│  3. System (embedded in binary)                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 What lvt Currently Does Well

| Capability | Status | Notes |
|------------|--------|-------|
| App scaffolding | ✅ Solid | Creates working Go module structure |
| CRUD generation | ⚠️ Fragile | Works but template drift causes issues |
| Migration system | ✅ Solid | goose-based, reliable |
| Kit abstraction | ⚠️ Complex | CSS helpers work but templates diverge |
| MCP integration | ✅ Solid | 16 tools, well-structured |
| Skill system | ⚠️ Probabilistic | 23 skills, but LLM doesn't always follow |

### 1.3 Generation Flow Analysis

```
User Request: "Create a blog with posts and comments"
        │
        ▼
┌───────────────────┐
│  LLM Skill Layer  │  ← Skills guide LLM behavior
│  (23 skills)      │     Problem: Probabilistic interpretation
└───────────────────┘
        │
        ▼
┌───────────────────┐
│  MCP Tool Layer   │  ← Tools execute commands
│  (16 tools)       │     Problem: No validation feedback loop
└───────────────────┘
        │
        ▼
┌───────────────────┐
│  Generator Layer  │  ← Templates assembled
│  (3 kits × N      │     Problem: Templates out of sync
│   templates)      │
└───────────────────┘
        │
        ▼
┌───────────────────┐
│  Output Layer     │  ← Files written
│  (Go code, SQL,   │     Problem: No compilation check
│   templates)      │
└───────────────────┘
        │
        ▼
        ❌ No automated validation that generated app works
```

---

## Part 2: Identified Gaps

### 2.1 Template System Gaps

**Gap 1: Template Duplication Without Sync**
```
Same template exists in 3+ locations:
├── internal/generator/templates/components/form.tmpl
├── internal/kits/system/multi/components/form.tmpl
└── internal/kits/system/single/components/form.tmpl

No mechanism ensures these stay in sync.
```

**Gap 2: Type Inconsistencies**
```go
// single kit (template.tmpl.tmpl:119)
{{if ne .EditingID ""}}   // STRING comparison

// multi kit (template.tmpl.tmpl:119)
{{if ne .EditingID 0}}    // INTEGER comparison
```

**Gap 3: Missing Features Across Kits**
- Single kit: Missing delete button in edit modal
- Simple kit: No modal support at all
- Multi kit: Has features other kits lack

### 2.2 Testing Infrastructure Gaps

| Gap | Impact | Severity |
|-----|--------|----------|
| No `go build` on generated code | Syntax errors ship | CRITICAL |
| No runtime test (app startup) | Crash-on-start apps ship | CRITICAL |
| `go mod tidy` skipped in tests | Dependency issues ship | HIGH |
| Browser tests not run on PRs | UI bugs ship | HIGH |
| Single/simple kits under-tested | Kit-specific bugs ship | MEDIUM |
| No concurrency testing | Race conditions ship | MEDIUM |

**Evidence from git history:**
- Commit `5700a82`: "single kit template bug" - shipped despite tests
- Commit `db9794b`: "fix: eliminate data race" - race shipped
- Commit `3af9a87`: "exclude EditingItem from JSON" - JSON structure bug shipped

### 2.3 Skill System Gaps

**Gap 1: Insufficient Constraint Strength**

The skill system evolved through 4+ iterations trying to prevent LLMs from asking questions:
```
Commit 54ac8ff: "Forbid AskUserQuestion, add examples"
Commit 1e7c680: "Stronger enforcement"
Commit ef53635: "Override superpowers:brainstorming"
Commit bee78f7: "Use generic language"
```

Despite all these attempts, LLMs still deviate from expected behavior.

**Gap 2: No Skill Testing Integration**

Skills have manual testing procedures but:
- No automated verification that skills produce working apps
- No regression testing when skills change
- No measurement of skill effectiveness

**Gap 3: No Feedback Mechanism**

When an LLM-generated app fails:
- No way to capture the failure
- No way to trace back to which skill/template caused it
- No way to automatically improve

### 2.4 Upstream Dependency Gaps

**LiveTemplate (v0.8.0) is in Alpha**

Recent changelog shows frequent breaking changes:
```
v0.8.0: "preserve statics for conditional blocks"
v0.7.12: "recognize append/prepend patterns"
v0.7.11: "prevent statics resend"
v0.7.10: "handle range→else transitions"
v0.7.9: "invalidate registry when conditional empty"
v0.7.8: "detect tree node changes"
```

Each of these required template updates in lvt.

---

## Part 3: Root Causes of Non-Determinism

### 3.1 The Template Drift Problem

```
Root Cause: Multiple template sources with no single source of truth

Impact Chain:
Template A updated → Template B forgotten → Kit X works, Kit Y breaks
                                                    ↓
                                         LLM uses broken template
                                                    ↓
                                         Generated app fails
                                                    ↓
                                         User blames "LLM"
```

### 3.2 The Validation Gap Problem

```
Root Cause: Tests verify file creation, not correctness

Current Test Flow:
1. Run lvt command ✓
2. Check files exist ✓
3. (MISSING) Build generated code
4. (MISSING) Run generated code
5. (MISSING) Test generated functionality

Result: "Tests pass" but generated apps broken
```

### 3.3 The Skill Probabilism Problem

```
Root Cause: Skills are suggestions, not constraints

LLM reads skill:
├── "Always run migrations after resource generation"
├── LLM decides: "I'll do it later" (50% of the time)
└── Result: Non-deterministic workflow

Compare to MCP tools:
├── Tool has required parameters
├── Tool enforces validation
└── Result: Deterministic execution
```

### 3.4 The Upstream Coupling Problem

```
Root Cause: Tight coupling to unstable upstream

LiveTemplate changes → Templates must change
                              ↓
                    Manual update process
                              ↓
                    Some templates missed
                              ↓
                    Broken generation
```

---

## Part 4: Architectural Proposals

### 4.1 Proposal: Single Source Template System

**Current State:**
```
3 kits × N templates = 3N separate files to maintain
```

**Proposed State:**
```
internal/templates/
├── base/                    # Shared base templates
│   ├── form.base.tmpl       # Single source of truth
│   ├── table.base.tmpl
│   └── modal.base.tmpl
├── overrides/               # Kit-specific overrides (minimal)
│   ├── multi/
│   │   └── layout.tmpl      # Only layout differs
│   └── single/
│       └── layout.tmpl
└── generated/               # Build-time generated
    ├── multi/               # Generated from base + overrides
    └── single/
```

**Implementation:**
```go
// Template composition at build time
type TemplateBuilder struct {
    Base      *template.Template
    Overrides map[string]*template.Template
}

func (b *TemplateBuilder) Build(kit string) *template.Template {
    result := b.Base.Clone()
    if override, ok := b.Overrides[kit]; ok {
        result = result.Override(override)
    }
    return result
}
```

**Benefits:**
- Single source of truth for shared functionality
- Automatic sync - change base, all kits update
- Clear override points for kit-specific behavior
- Eliminates template drift

### 4.2 Proposal: Compilation-Verified Generation

**Add Post-Generation Validation:**

```go
// internal/generator/validator.go

type GenerationValidator struct {
    TempDir string
}

func (v *GenerationValidator) Validate(appPath string) *ValidationResult {
    result := &ValidationResult{}

    // Step 1: Go module validation
    if err := v.validateGoMod(appPath); err != nil {
        result.AddError("go.mod", err)
    }

    // Step 2: Compilation check
    if err := v.compileApp(appPath); err != nil {
        result.AddError("compilation", err)
    }

    // Step 3: Template parse check
    if err := v.parseTemplates(appPath); err != nil {
        result.AddError("templates", err)
    }

    // Step 4: Migration validity
    if err := v.validateMigrations(appPath); err != nil {
        result.AddError("migrations", err)
    }

    return result
}

func (v *GenerationValidator) compileApp(appPath string) error {
    cmd := exec.Command("go", "build", "-o", "/dev/null", "./...")
    cmd.Dir = appPath
    cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
    return cmd.Run()
}
```

**Integration Points:**

1. **Post-generation hook:**
```go
func GenerateResource(...) error {
    // ... generation code ...

    // NEW: Validate before returning success
    validator := NewGenerationValidator()
    if result := validator.Validate(appPath); !result.IsValid() {
        return fmt.Errorf("generation validation failed: %v", result.Errors())
    }

    return nil
}
```

2. **MCP tool with validation:**
```go
func handleLvtGenResource(...) (*mcp.CallToolResult, error) {
    // ... generation ...

    // Run validation
    validation := validator.Validate(appPath)

    return &mcp.CallToolResult{
        Content: []interface{}{
            mcp.TextContent{
                Text: fmt.Sprintf(
                    "Generated resource: %s\n"+
                    "Compilation: %s\n"+
                    "Templates: %s\n"+
                    "Migrations: %s",
                    resourceName,
                    validation.CompilationStatus(),
                    validation.TemplateStatus(),
                    validation.MigrationStatus(),
                ),
            },
        },
    }, nil
}
```

### 4.3 Proposal: Skill Enforcement Through Tools

**Current Problem:** Skills are natural language suggestions that LLMs may ignore.

**Solution:** Encode critical workflows as MCP tool chains with enforced ordering.

```go
// MCP tool that enforces workflow
func registerWorkflowTool(server *mcp.Server) {
    tool := &mcp.Tool{
        Name: "lvt_workflow_add_resource",
        Description: `Complete resource workflow with validation.
This tool handles: generation, migration, validation, and optional seeding.
Use this instead of calling individual tools.`,
        InputSchema: WorkflowInput{
            ResourceName: "required",
            Fields:       "required",
            RunMigration: "default:true",
            SeedCount:    "default:0",
            Validate:     "default:true",
        },
    }

    handler := func(ctx context.Context, input WorkflowInput) (*Result, error) {
        // Step 1: Generate resource
        if err := generateResource(input); err != nil {
            return nil, fmt.Errorf("generation failed: %w", err)
        }

        // Step 2: Run migration (enforced, not optional)
        if input.RunMigration {
            if err := runMigration(input.AppPath); err != nil {
                return nil, fmt.Errorf("migration failed: %w", err)
            }
        }

        // Step 3: Validate (enforced)
        if input.Validate {
            validation := validator.Validate(input.AppPath)
            if !validation.IsValid() {
                return nil, fmt.Errorf("validation failed: %v", validation.Errors())
            }
        }

        // Step 4: Seed if requested
        if input.SeedCount > 0 {
            if err := seedData(input.ResourceName, input.SeedCount); err != nil {
                return nil, fmt.Errorf("seeding failed: %w", err)
            }
        }

        return &Result{
            Success: true,
            Message: "Resource created, migrated, validated, and seeded",
        }, nil
    }

    mcp.AddTool(server, tool, handler)
}
```

**Benefits:**
- Workflow is deterministic - tools enforce ordering
- Validation is mandatory, not optional
- Errors are caught before user sees "success"
- LLM cannot skip steps

### 4.4 Proposal: Template Version Tracking

**Add Version Metadata to Templates:**

```yaml
# template.tmpl.tmpl header
# lvt-template-version: 2.3.0
# livetemplate-compat: >=0.8.0
# depends-on: [form.tmpl:1.2.0, table.tmpl:1.0.0]
# last-sync: 2026-01-15
```

**Compatibility Checker:**

```go
type TemplateVersion struct {
    Version      semver.Version
    Compat       semver.Constraint
    Dependencies map[string]semver.Version
    LastSync     time.Time
}

func CheckTemplateCompatibility(kit string, livetemplateVersion string) []Warning {
    var warnings []Warning

    templates := loadKitTemplates(kit)
    for _, t := range templates {
        version := parseTemplateVersion(t)

        // Check LiveTemplate compatibility
        if !version.Compat.Check(livetemplateVersion) {
            warnings = append(warnings, Warning{
                Template: t.Name,
                Message:  fmt.Sprintf("requires livetemplate %s, got %s",
                    version.Compat, livetemplateVersion),
            })
        }

        // Check dependency versions
        for dep, requiredVer := range version.Dependencies {
            actualVer := getTemplateVersion(kit, dep)
            if actualVer.LessThan(requiredVer) {
                warnings = append(warnings, Warning{
                    Template: t.Name,
                    Message:  fmt.Sprintf("requires %s@%s, got %s",
                        dep, requiredVer, actualVer),
                })
            }
        }
    }

    return warnings
}
```

---

## Part 5: Evolution Mechanism Design

### 5.1 The Self-Healing Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                    SELF-HEALING GENERATION PIPELINE                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐         │
│  │ Generate │───▶│ Validate│───▶│ Report  │───▶│ Analyze │         │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘         │
│       │              │              │              │                │
│       │              │              │              ▼                │
│       │              │              │       ┌─────────┐            │
│       │              │              │       │ Propose │            │
│       │              │              │       │  Fixes  │            │
│       │              │              │       └─────────┘            │
│       │              │              │              │                │
│       │              │              │              ▼                │
│       │              │              │       ┌─────────┐            │
│       │              │              └──────▶│  Apply  │            │
│       │              │                      │  Fixes  │            │
│       │              │                      └─────────┘            │
│       │              │                           │                  │
│       │              │                           ▼                  │
│       │              │                      ┌─────────┐            │
│       └──────────────┴─────────────────────▶│ Verify  │            │
│                                             └─────────┘            │
│                                                  │                  │
│                                                  ▼                  │
│                                             ┌─────────┐            │
│                                             │ Commit  │            │
│                                             └─────────┘            │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 Failure Capture System

**New Component: Generation Telemetry**

```go
// internal/telemetry/generation.go

type GenerationEvent struct {
    ID            string
    Timestamp     time.Time
    Command       string
    Inputs        map[string]interface{}
    Kit           string
    TemplateVers  map[string]string
    Success       bool
    ValidationRes *ValidationResult
    Errors        []Error
    Duration      time.Duration
}

type Error struct {
    Phase    string  // "generation", "compilation", "runtime", "template"
    File     string
    Line     int
    Message  string
    Context  string  // surrounding code/template
}

func CaptureGeneration(cmd string, inputs map[string]interface{}) *GenerationCapture {
    return &GenerationCapture{
        Event: &GenerationEvent{
            ID:        uuid.New().String(),
            Timestamp: time.Now(),
            Command:   cmd,
            Inputs:    inputs,
        },
    }
}

func (c *GenerationCapture) RecordError(phase, file string, line int, msg, context string) {
    c.Event.Errors = append(c.Event.Errors, Error{
        Phase:   phase,
        File:    file,
        Line:    line,
        Message: msg,
        Context: context,
    })
}

func (c *GenerationCapture) Finalize(success bool, validation *ValidationResult) {
    c.Event.Success = success
    c.Event.ValidationRes = validation
    c.Event.Duration = time.Since(c.Event.Timestamp)

    // Store for analysis
    telemetryStore.Save(c.Event)
}
```

### 5.3 Automated Fix Proposal System

**New Component: Fix Proposer**

```go
// internal/evolution/proposer.go

type FixProposer struct {
    KnowledgeBase *KnowledgeBase  // Pattern → Fix mapping
    LLMClient     *LLMClient      // For novel issues
}

type Fix struct {
    ID          string
    TargetFile  string
    ChangeType  string  // "edit", "add", "delete"
    OldContent  string
    NewContent  string
    Confidence  float64
    Rationale   string
}

func (p *FixProposer) ProposeFixesFor(event *GenerationEvent) []Fix {
    var fixes []Fix

    for _, err := range event.Errors {
        // First, check knowledge base for known patterns
        if knownFixes := p.KnowledgeBase.LookupFixes(err); len(knownFixes) > 0 {
            fixes = append(fixes, knownFixes...)
            continue
        }

        // For unknown patterns, use LLM to propose fix
        if fix := p.proposeWithLLM(err, event); fix != nil {
            fix.Confidence = 0.7  // Lower confidence for LLM-generated
            fixes = append(fixes, *fix)
        }
    }

    return fixes
}

func (p *FixProposer) proposeWithLLM(err Error, event *GenerationEvent) *Fix {
    prompt := fmt.Sprintf(`
Analyze this generation error and propose a fix:

Error Phase: %s
Error File: %s
Error Line: %d
Error Message: %s
Context:
%s

Kit: %s
Template Versions: %v

Propose a specific code change to fix this error.
Output format:
FILE: <path>
OLD:
<old content>
NEW:
<new content>
RATIONALE: <explanation>
`, err.Phase, err.File, err.Line, err.Message, err.Context,
        event.Kit, event.TemplateVers)

    response := p.LLMClient.Generate(prompt)
    return parseFix(response)
}
```

### 5.4 Knowledge Base for Common Fixes

```go
// internal/evolution/knowledge.go

type KnowledgeBase struct {
    Patterns []Pattern
}

type Pattern struct {
    Name        string
    Matcher     func(Error) bool
    Fixes       []FixTemplate
    Confidence  float64
}

var DefaultPatterns = []Pattern{
    {
        Name: "EditingID type mismatch",
        Matcher: func(e Error) bool {
            return strings.Contains(e.Message, "cannot convert") &&
                   strings.Contains(e.Context, "EditingID")
        },
        Fixes: []FixTemplate{
            {
                File:    "template.tmpl.tmpl",
                Pattern: `{{if ne .EditingID 0}}`,
                Replace: `{{if ne .EditingID ""}}`,
            },
        },
        Confidence: 0.95,
    },
    {
        Name: "Missing form sync script",
        Matcher: func(e Error) bool {
            return strings.Contains(e.Message, "select value reverted") ||
                   strings.Contains(e.Message, "morphdom")
        },
        Fixes: []FixTemplate{
            {
                File:    "layout.tmpl",
                Pattern: `</body>`,
                Replace: `<script>/* form sync script */</script></body>`,
            },
        },
        Confidence: 0.9,
    },
    // ... more patterns
}
```

### 5.5 Automated Testing of Proposed Fixes

```go
// internal/evolution/tester.go

type FixTester struct {
    TempDir   string
    Validator *GenerationValidator
}

func (t *FixTester) TestFix(fix Fix, originalEvent *GenerationEvent) *TestResult {
    // Create isolated test environment
    testDir := t.createIsolatedEnv()
    defer os.RemoveAll(testDir)

    // Apply fix to templates
    if err := t.applyFix(testDir, fix); err != nil {
        return &TestResult{Success: false, Error: err}
    }

    // Re-run the same generation that failed
    if err := t.regenerate(testDir, originalEvent); err != nil {
        return &TestResult{Success: false, Error: err}
    }

    // Validate
    validation := t.Validator.Validate(testDir)

    // Run additional tests
    testResults := t.runTests(testDir)

    return &TestResult{
        Success:    validation.IsValid() && testResults.AllPassed(),
        Validation: validation,
        Tests:      testResults,
    }
}
```

### 5.6 Fix Review and Application Process

```go
// internal/evolution/reviewer.go

type FixReviewer struct {
    AutoApplyThreshold float64  // e.g., 0.95
    ReviewQueue        *ReviewQueue
}

func (r *FixReviewer) ProcessFix(fix Fix, testResult *TestResult) {
    if !testResult.Success {
        // Don't apply failed fixes
        log.Printf("Fix %s failed testing, discarding", fix.ID)
        return
    }

    if fix.Confidence >= r.AutoApplyThreshold && testResult.AllPassed() {
        // High confidence + all tests pass = auto-apply
        r.autoApply(fix)
    } else {
        // Queue for human/LLM review
        r.ReviewQueue.Add(ReviewItem{
            Fix:        fix,
            TestResult: testResult,
            Reason:     "Confidence below threshold or some tests failed",
        })
    }
}

func (r *FixReviewer) autoApply(fix Fix) {
    // Create branch
    branch := fmt.Sprintf("auto-fix/%s", fix.ID)
    git.CreateBranch(branch)

    // Apply fix
    applyFix(fix)

    // Commit
    git.Commit(fmt.Sprintf("fix: %s\n\nAuto-applied fix with confidence %.2f\n\nRationale: %s",
        fix.ID, fix.Confidence, fix.Rationale))

    // Create PR
    createPR(branch, fix)
}
```

### 5.7 Skill Evolution Mechanism

**Skill Quality Metrics:**

> **Note:** See [EVOLUTION_SYSTEM_SPEC.md](./EVOLUTION_SYSTEM_SPEC.md#8-skill-evolution-pipeline) for the full `SkillMetrics` type definition and implementation details.

The evolution system tracks metrics for each skill to identify which need improvement:

```go
// internal/evolution/skill_metrics.go
// Simplified example - see EVOLUTION_SYSTEM_SPEC.md for full implementation

metrics := evolution.ComputeSkillMetrics("add-resource", events)
// metrics.SuccessRate       - Apps that compiled successfully
// metrics.CompilationSuccessRate - Compilation success rate
// metrics.AverageErrors     - Average errors per generation
// metrics.CommonErrors      - Most frequent error patterns
```

**Skill Improvement Proposer:**

> **Note:** See [EVOLUTION_SYSTEM_SPEC.md](./EVOLUTION_SYSTEM_SPEC.md#82-skill-improvement-proposals) for the full `SkillImprover` implementation.

The skill improver uses LLM analysis to propose changes when metrics indicate poor performance:

```go
// internal/evolution/skill_improver.go
// Simplified example - see EVOLUTION_SYSTEM_SPEC.md for full implementation

improver := evolution.NewSkillImprover(llmClient)
changes := improver.ProposeImprovements(metrics, failures)
// Returns []SkillChange with:
// - SkillName, Section, Old/New content, Rationale
// - Estimated impact on success rate
```

---

## Part 6: Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)

**1.1 Add Compilation Validation**
- [ ] Create `internal/generator/validator.go`
- [ ] Add `go build` check after generation
- [ ] Add template parse check
- [ ] Integrate into MCP tools

**1.2 Add Generation Telemetry**
- [ ] Create `internal/telemetry/` package
- [ ] Capture all generation events
- [ ] Store locally (SQLite initially)
- [ ] Add error context capture

### Phase 2: Template Consolidation (Weeks 3-4)

**2.1 Single Source Templates**
- [ ] Create `internal/templates/base/` directory
- [ ] Consolidate form.tmpl to single source
- [ ] Consolidate table.tmpl to single source
- [ ] Create kit-specific overrides directory
- [ ] Add template composition build step

**2.2 Template Versioning**
- [ ] Add version headers to all templates
- [ ] Create compatibility checker
- [ ] Add version mismatch warnings

### Phase 3: Evolution System (Weeks 5-6)

**3.1 Fix Proposal System**
- [ ] Create `internal/evolution/` package
- [ ] Implement pattern-based fix proposal
- [ ] Implement LLM-based fix proposal
- [ ] Create knowledge base with known patterns

**3.2 Automated Fix Testing**
- [ ] Create isolated test environment system
- [ ] Implement fix application and testing
- [ ] Add regression test suite for fixes

### Phase 4: Skill Evolution (Weeks 7-8)

**4.1 Skill Metrics**
- [ ] Track skill usage and success rates
- [ ] Track LLM deviations from skills
- [ ] Create skill effectiveness dashboard

**4.2 Skill Improvement**
- [ ] Implement skill improvement proposer
- [ ] Create skill change review process
- [ ] Add A/B testing for skill versions

### Phase 5: Continuous Improvement (Ongoing)

**5.1 Feedback Loop**
- [ ] Daily telemetry analysis
- [ ] Weekly fix proposal review
- [ ] Monthly skill effectiveness review

**5.2 Quality Gates**
- [ ] Block releases with <95% compilation rate
- [ ] Require all kits to pass same test suite
- [ ] Automated regression detection

---

## Part 7: Immediate Action Items

### 7.1 Quick Wins (This Week)

1. **Add `go build` to E2E tests:**
```go
// In e2e tests, after generation:
func TestResourceGeneration(t *testing.T) {
    // ... generation ...

    // NEW: Verify compilation
    cmd := exec.Command("go", "build", "-o", "/dev/null", "./...")
    cmd.Dir = appPath
    if err := cmd.Run(); err != nil {
        t.Fatalf("Generated code does not compile: %v", err)
    }
}
```

2. **Fix EditingID type inconsistency:**
```
File: internal/kits/system/multi/templates/resource/template.tmpl.tmpl
Change: {{if ne .EditingID 0}} → {{if ne .EditingID ""}}
```

3. **Add delete button to single kit:**
```
File: internal/kits/system/single/templates/resource/template.tmpl.tmpl
Add delete button matching multi kit template
```

### 7.2 Medium-Term (This Month)

1. **Remove agent test `go mod tidy` skip:**
```go
// Remove these lines from agenttest/harness.go:
// os.Setenv("SKIP_GO_MOD_TIDY", "1")
```

2. **Enable browser tests on PRs:**
```yaml
# .github/workflows/test.yml
- name: Run browser tests
  run: go test ./e2e -tags=browser -v
```

3. **Create template sync check:**
```bash
# Script to detect template drift
diff internal/kits/system/multi/components/form.tmpl \
     internal/kits/system/single/components/form.tmpl
```

### 7.3 Long-Term (This Quarter)

1. Implement single-source template system
2. Deploy generation telemetry
3. Build fix proposal pipeline
4. Launch skill metrics dashboard

---

## Appendix A: Git History Evidence

### Common Bug Patterns (from 77% fix commits)

| Pattern | Occurrences | Root Cause |
|---------|-------------|------------|
| Modal state issues | 4+ | Server vs client state confusion |
| Form sync issues | 2+ | morphdom value preservation |
| Template type errors | 3+ | String vs int comparisons |
| Auth session bugs | 4+ | Session not cleared properly |
| Kit-specific bugs | 3+ | Templates out of sync |

### Skill Iteration Evidence

```
54ac8ff: "Forbid AskUserQuestion"
1e7c680: "Stronger enforcement"
ef53635: "Override superpowers:brainstorming"
bee78f7: "Use generic language"
```

Each iteration attempted to make skills more deterministic, but LLM behavior remains probabilistic.

---

## Appendix B: LiveTemplate Upstream Changelog

Critical recent changes requiring template updates:

```
v0.8.0 - preserve statics for conditional blocks
v0.7.12 - recognize append/prepend patterns
v0.7.11 - prevent statics resend
v0.7.10 - handle range→else transitions
v0.7.9 - invalidate registry when conditional empty
v0.7.8 - detect tree node changes
v0.7.5 - handle non-TreeNode transitions
v0.7.4 - ensure Range.Statics for empty→items
v0.7.3 - support heterogeneous range items
```

Each of these could break existing templates. Need automated compatibility testing.

---

## Conclusion

Making lvt generation deterministic requires:

1. **Single source of truth for templates** - eliminate drift
2. **Mandatory compilation validation** - catch errors before "success"
3. **Workflow tools over skills** - enforce ordering through tools
4. **Telemetry and feedback** - capture failures, learn from them
5. **Automated evolution** - propose, test, and apply fixes

The proposed evolution mechanism creates a self-healing system where:
- Every failure is captured
- Patterns are learned
- Fixes are proposed and tested
- Improvements are applied automatically when safe
- Human review handles uncertain cases

This transforms lvt from a static tool into an evolving system that gets better with every use.
