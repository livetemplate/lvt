# LVT Evolution System: Technical Specification

**Version**: 1.0
**Status**: Proposed

## Overview

This document specifies the technical design for a self-improving generation system where LLMs can:
1. Detect when generation produces broken apps
2. Capture failure context for analysis
3. Propose fixes to templates, skills, and tests
4. Automatically apply high-confidence fixes
5. Track improvement over time

---

## 1. Core Components

### 1.1 Component Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         EVOLUTION SYSTEM                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │    Telemetry    │  │   Validation    │  │    Knowledge    │             │
│  │    Collector    │  │    Engine       │  │      Base       │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
│           │                    │                    │                       │
│           ▼                    ▼                    ▼                       │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │                      Event Store                             │           │
│  │              (SQLite: $XDG_CONFIG_HOME/lvt/evolution.db)            │           │
│  └─────────────────────────────────────────────────────────────┘           │
│           │                    │                    │                       │
│           ▼                    ▼                    ▼                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │      Fix        │  │      Fix        │  │   Skill         │             │
│  │    Proposer     │  │    Tester       │  │   Improver      │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
│           │                    │                    │                       │
│           ▼                    ▼                    ▼                       │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │                    Review Queue                              │           │
│  │         (Human or LLM review for uncertain fixes)            │           │
│  └─────────────────────────────────────────────────────────────┘           │
│                                │                                            │
│                                ▼                                            │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │                    Git Integration                           │           │
│  │      (Branch, Commit, PR creation for approved fixes)        │           │
│  └─────────────────────────────────────────────────────────────┘           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Package Structure

```
internal/
├── validation/               # Separate validation package (Issue 2.1)
│   ├── validation.go         # Main engine and types
│   ├── go_mod.go             # go.mod validation
│   ├── compilation.go        # Go build validation
│   ├── templates.go          # Template parse validation
│   ├── migrations.go         # Migration validation
│   └── runtime.go            # Runtime validation (Issue 2.4)
├── telemetry/                # Event capture (Issue 3.1)
│   ├── telemetry.go          # Main API
│   ├── events.go             # Event types
│   ├── store.go              # Storage interface
│   └── sqlite.go             # SQLite implementation
├── evolution/                # Self-improvement system
│   ├── evolution.go          # Main entry point
│   ├── proposer.go           # Fix proposal logic (Issue 3.4)
│   ├── tester.go             # Fix testing in isolation (Issue 3.5)
│   ├── reviewer.go           # Review queue management
│   ├── skill_metrics.go      # Skill effectiveness tracking
│   ├── skill_improver.go     # Skill improvement proposals
│   └── git_integration.go    # Git operations
├── evolution/knowledge/      # Pattern knowledge base (Issue 3.2)
│   ├── knowledge.go          # Knowledge base API
│   ├── parser.go             # Markdown parser for patterns.md
│   ├── matcher.go            # Error-to-pattern matching
│   └── types.go              # Pattern, Fix types
└── evolution/store/          # Event persistence
    ├── store.go              # Event store interface
    ├── sqlite.go             # SQLite implementation
    └── schema.sql            # Database schema

evolution/
└── patterns.md               # Source of truth - git tracked (see Section 4)
```

> **Note:** This package structure aligns with IMPLEMENTATION_PLAN.md milestones 2 and 3.

---

## 2. Telemetry System

### 2.1 Event Types

```go
// internal/evolution/types/events.go

package types

import "time"

// GenerationEvent captures everything about a generation attempt
type GenerationEvent struct {
    ID              string                 `json:"id"`
    Timestamp       time.Time              `json:"timestamp"`
    SessionID       string                 `json:"session_id"`

    // What was requested
    Command         string                 `json:"command"`       // e.g., "gen resource"
    Inputs          map[string]interface{} `json:"inputs"`        // Command arguments
    Skill           string                 `json:"skill,omitempty"`  // Skill that triggered this

    // Context
    Kit             string                 `json:"kit"`
    TemplateVersions map[string]string     `json:"template_versions"`
    LiveTemplateVer  string                `json:"livetemplate_version"`
    LvtVersion       string                `json:"lvt_version"`

    // Results
    Success         bool                   `json:"success"`
    Validation      *ValidationResult      `json:"validation,omitempty"`
    Errors          []GenerationError      `json:"errors,omitempty"`

    // Metadata
    Duration        time.Duration          `json:"duration"`
    FilesGenerated  []string               `json:"files_generated"`
    FilesModified   []string               `json:"files_modified"`
}

// GenerationError captures detailed error information
type GenerationError struct {
    Phase       string `json:"phase"`        // "generation", "compilation", "runtime", "template"
    File        string `json:"file"`
    Line        int    `json:"line,omitempty"`
    Column      int    `json:"column,omitempty"`
    Message     string `json:"message"`
    ErrorCode   string `json:"error_code,omitempty"`  // For pattern matching
    Context     string `json:"context"`               // Surrounding code
    StackTrace  string `json:"stack_trace,omitempty"`
}

// ValidationResult from post-generation validation
type ValidationResult struct {
    Compilation  ValidationCheck `json:"compilation"`
    Templates    ValidationCheck `json:"templates"`
    Migrations   ValidationCheck `json:"migrations"`
    GoMod        ValidationCheck `json:"go_mod"`
    Runtime      ValidationCheck `json:"runtime,omitempty"`
}

type ValidationCheck struct {
    Passed  bool     `json:"passed"`
    Errors  []string `json:"errors,omitempty"`
    Warnings []string `json:"warnings,omitempty"`
}
```

### 2.2 Event Capture Integration

```go
// internal/evolution/telemetry.go

package evolution

import (
    "context"
    "time"

    "github.com/livetemplate/lvt/internal/evolution/types"
    "github.com/livetemplate/lvt/internal/evolution/store"
)

type TelemetryCollector struct {
    store   store.EventStore
    enabled bool
}

func NewTelemetryCollector() *TelemetryCollector {
    return &TelemetryCollector{
        store:   store.NewSQLiteStore(),
        enabled: true,  // Can be disabled via config
    }
}

// StartGeneration begins tracking a generation event
func (t *TelemetryCollector) StartGeneration(cmd string, inputs map[string]interface{}) *GenerationCapture {
    if !t.enabled {
        return &GenerationCapture{noop: true}
    }

    return &GenerationCapture{
        collector: t,
        event: &types.GenerationEvent{
            ID:        generateUUID(),
            Timestamp: time.Now(),
            SessionID: getSessionID(),
            Command:   cmd,
            Inputs:    inputs,
            Kit:       getCurrentKit(),
            TemplateVersions: getTemplateVersions(),
            LiveTemplateVer:  getLiveTemplateVersion(),
            LvtVersion:       Version,
        },
    }
}

// GenerationCapture tracks a single generation attempt
type GenerationCapture struct {
    collector *TelemetryCollector
    event     *types.GenerationEvent
    noop      bool
}

func (c *GenerationCapture) SetSkill(skill string) {
    if c.noop { return }
    c.event.Skill = skill
}

func (c *GenerationCapture) RecordError(err types.GenerationError) {
    if c.noop { return }
    c.event.Errors = append(c.event.Errors, err)
}

func (c *GenerationCapture) RecordFileGenerated(path string) {
    if c.noop { return }
    c.event.FilesGenerated = append(c.event.FilesGenerated, path)
}

func (c *GenerationCapture) Finalize(success bool, validation *types.ValidationResult) {
    if c.noop { return }

    c.event.Success = success
    c.event.Validation = validation
    c.event.Duration = time.Since(c.event.Timestamp)

    // Store event
    c.collector.store.Save(c.event)

    // Trigger analysis if failure
    if !success {
        go c.collector.analyzeFailure(c.event)
    }
}
```

### 2.3 Integration into Generator

```go
// internal/generator/resource.go (modified)

func GenerateResource(basePath, moduleName, resourceName string,
                      fields []Field, kit string, ...) error {

    // START: Telemetry capture
    capture := telemetry.StartGeneration("gen resource", map[string]interface{}{
        "resource_name": resourceName,
        "fields":        fields,
        "kit":           kit,
    })
    defer func() {
        validation := validator.Validate(basePath)
        capture.Finalize(validation.AllPassed(), validation)
    }()

    // ... existing generation code ...

    // On error, capture context
    if err != nil {
        capture.RecordError(types.GenerationError{
            Phase:   "generation",
            Message: err.Error(),
            Context: getErrorContext(err),
        })
        return err
    }

    // Track generated files
    capture.RecordFileGenerated(handlerPath)
    capture.RecordFileGenerated(templatePath)

    return nil
}
```

---

## 3. Validation Engine

### 3.1 Validation Pipeline

```go
// internal/evolution/validator.go

package evolution

import (
    "os/exec"
    "path/filepath"

    "github.com/livetemplate/lvt/internal/evolution/types"
)

type ValidationEngine struct {
    timeout time.Duration
}

func NewValidationEngine() *ValidationEngine {
    return &ValidationEngine{
        timeout: 60 * time.Second,
    }
}

func (v *ValidationEngine) Validate(appPath string) *types.ValidationResult {
    result := &types.ValidationResult{}

    // Run validations in parallel where possible
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        result.GoMod = v.validateGoMod(appPath)
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        result.Templates = v.validateTemplates(appPath)
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        result.Migrations = v.validateMigrations(appPath)
    }()

    wg.Wait()

    // Compilation depends on go.mod being valid
    if result.GoMod.Passed {
        result.Compilation = v.validateCompilation(appPath)
    } else {
        result.Compilation = types.ValidationCheck{
            Passed: false,
            Errors: []string{"Skipped: go.mod validation failed"},
        }
    }

    return result
}

func (v *ValidationEngine) validateGoMod(appPath string) types.ValidationCheck {
    // Check go.mod exists
    goModPath := filepath.Join(appPath, "go.mod")
    if _, err := os.Stat(goModPath); err != nil {
        return types.ValidationCheck{
            Passed: false,
            Errors: []string{"go.mod not found"},
        }
    }

    // Run go mod tidy (catches dependency issues)
    cmd := exec.Command("go", "mod", "tidy")
    cmd.Dir = appPath

    output, err := cmd.CombinedOutput()
    if err != nil {
        return types.ValidationCheck{
            Passed: false,
            Errors: []string{string(output)},
        }
    }

    return types.ValidationCheck{Passed: true}
}

func (v *ValidationEngine) validateCompilation(appPath string) types.ValidationCheck {
    // Use os.DevNull for cross-platform compatibility
    cmd := exec.Command("go", "build", "-o", os.DevNull, "./...")
    cmd.Dir = appPath
    cmd.Env = append(os.Environ(), "CGO_ENABLED=1")

    output, err := cmd.CombinedOutput()
    if err != nil {
        // Parse compilation errors
        errors := parseCompilationErrors(string(output))
        return types.ValidationCheck{
            Passed: false,
            Errors: errors,
        }
    }

    return types.ValidationCheck{Passed: true}
}

func (v *ValidationEngine) validateTemplates(appPath string) types.ValidationCheck {
    // Find all .tmpl files
    var templates []string
    filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".tmpl") {
            templates = append(templates, path)
        }
        return nil
    })

    var errors []string
    for _, tmpl := range templates {
        if err := parseTemplate(tmpl); err != nil {
            errors = append(errors, fmt.Sprintf("%s: %v", tmpl, err))
        }
    }

    return types.ValidationCheck{
        Passed: len(errors) == 0,
        Errors: errors,
    }
}

func (v *ValidationEngine) validateMigrations(appPath string) types.ValidationCheck {
    migrationsPath := filepath.Join(appPath, "database", "migrations")

    // Check migrations directory exists
    if _, err := os.Stat(migrationsPath); err != nil {
        return types.ValidationCheck{Passed: true}  // No migrations = valid
    }

    // Validate SQL syntax in each migration
    var errors []string
    filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".sql") {
            if err := validateSQL(path); err != nil {
                errors = append(errors, fmt.Sprintf("%s: %v", path, err))
            }
        }
        return nil
    })

    return types.ValidationCheck{
        Passed: len(errors) == 0,
        Errors: errors,
    }
}
```

---

## 4. Knowledge Base

**Key Design Decision:** Patterns are stored in a **markdown file** (`evolution/patterns.md`) rather than in code or a database. This enables:
- Git tracking of pattern history
- Easy reading by humans and LLMs
- LLM-proposed pattern additions via PR
- No recompilation needed to add patterns

### 4.1 Patterns Markdown File

```markdown
<!-- evolution/patterns.md -->
# Evolution Patterns

Known error patterns and their fixes. This file is the source of truth
for the evolution system's knowledge base.

Format: Each pattern is an H2 section with structured subsections.

---

## Pattern: editing-id-type

**Name:** EditingID Type Mismatch
**Confidence:** 0.95
**Added:** 2026-01-15
**Fix Count:** 12
**Success Rate:** 0.92

### Description
EditingID compared as integer but is string type in handler.
This causes compilation errors when templates expect string comparison.

### Error Pattern
- **Phase:** compilation
- **Message Regex:** `cannot convert .* to type int.*EditingID`
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
**Added:** 2026-01-15
**Fix Count:** 8
**Success Rate:** 0.88

### Description
Modal editing state (IsAdding, EditingID) persists on page reload
because fields are not marked as transient.

### Error Pattern
- **Phase:** runtime
- **Message Regex:** `modal (open|visible) after (reload|refresh)`
- **Context Regex:** `(IsAdding|EditingID)`

### Fix
- **File:** `*/handler.go.tmpl`
- **Find:** `EditingID string`
- **Replace:** `EditingID string \`lvt:"transient"\``
- **Is Regex:** false

---

## Pattern: form-sync-morphdom

**Name:** Form Values Revert After Update
**Confidence:** 0.88
**Added:** 2026-01-15
**Fix Count:** 5
**Success Rate:** 0.80

### Description
Select dropdown values revert to previous state after morphdom
DOM patching because expected value not preserved.

### Error Pattern
- **Phase:** runtime
- **Message Regex:** `(select|dropdown) value (reverted|reset|changed)`

### Fix
- **File:** `*/components/form.tmpl`
- **Find:** `<select name="{{.Name}}"`
- **Replace:** `<select name="{{.Name}}" data-expected-value="{{.Value}}"`
- **Is Regex:** false

---

## Pattern: session-not-cleared

**Name:** Session State Not Cleared on Login
**Confidence:** 0.88
**Added:** 2026-01-15
**Fix Count:** 4
**Success Rate:** 0.75

### Description
LiveTemplate session persists stale state after auth changes,
showing cached IsLoggedIn=false after successful login.

### Error Pattern
- **Phase:** runtime
- **Message Regex:** `(IsLoggedIn|session).*(stale|persisted|cached|wrong)`

### Fix
- **File:** `*/auth/login.go.tmpl`
- **Find:** `return nil`
- **Replace:** `ctx.ClearSession()\n\treturn nil`
- **Is Regex:** false

---

<!-- New patterns are added above this line -->
<!-- The evolution system can propose new patterns via PR -->
```

### 4.2 Pattern Types

```go
// internal/evolution/knowledge/types.go

package knowledge

import "regexp"

// Pattern represents a known error pattern and its fix
type Pattern struct {
    ID          string
    Name        string
    Description string

    // Matching
    ErrorPhase  string           // "compilation", "runtime", "template"
    MessageRe   *regexp.Regexp   // Regex to match error message
    ContextRe   *regexp.Regexp   // Regex to match context

    // Fix
    Fixes       []FixTemplate
    Confidence  float64

    // Metadata (updated by evolution system)
    Added       string   // Date added
    FixCount    int      // Times this fix was applied
    SuccessRate float64  // Success rate of applications
}

// FixTemplate describes how to apply a fix
type FixTemplate struct {
    File        string  // Glob pattern (e.g., "*/template.tmpl.tmpl")
    FindPattern string  // What to find
    Replace     string  // What to replace with
    IsRegex     bool    // Is FindPattern a regex?
}
```

### 4.3 Markdown Parser

```go
// internal/evolution/knowledge/parser.go

package knowledge

import (
    "bufio"
    "os"
    "regexp"
    "strings"
)

// ParsePatternsFile parses evolution/patterns.md
func ParsePatternsFile(path string) ([]*Pattern, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var patterns []*Pattern
    var current *Pattern
    var section string

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        // New pattern starts with "## Pattern: "
        if strings.HasPrefix(line, "## Pattern: ") {
            if current != nil {
                patterns = append(patterns, current)
            }
            current = &Pattern{
                ID: strings.TrimPrefix(line, "## Pattern: "),
            }
            section = ""
            continue
        }

        if current == nil {
            continue
        }

        // Parse metadata fields
        if strings.HasPrefix(line, "**Name:**") {
            current.Name = strings.TrimSpace(strings.TrimPrefix(line, "**Name:**"))
        } else if strings.HasPrefix(line, "**Confidence:**") {
            fmt.Sscanf(line, "**Confidence:** %f", &current.Confidence)
        } else if strings.HasPrefix(line, "**Added:**") {
            current.Added = strings.TrimSpace(strings.TrimPrefix(line, "**Added:**"))
        } else if strings.HasPrefix(line, "**Fix Count:**") {
            fmt.Sscanf(line, "**Fix Count:** %d", &current.FixCount)
        } else if strings.HasPrefix(line, "**Success Rate:**") {
            fmt.Sscanf(line, "**Success Rate:** %f", &current.SuccessRate)
        }

        // Track which section we're in
        if strings.HasPrefix(line, "### ") {
            section = strings.TrimPrefix(line, "### ")
            continue
        }

        // Parse section content
        switch section {
        case "Description":
            if line != "" && !strings.HasPrefix(line, "#") {
                current.Description += line + " "
            }
        case "Error Pattern":
            parseErrorPattern(current, line)
        case "Fix":
            parseFix(current, line)
        }
    }

    // Don't forget the last pattern
    if current != nil {
        patterns = append(patterns, current)
    }

    return patterns, scanner.Err()
}

func parseErrorPattern(p *Pattern, line string) error {
    if strings.HasPrefix(line, "- **Phase:**") {
        p.ErrorPhase = extractValue(line)
    } else if strings.HasPrefix(line, "- **Message Regex:**") {
        val := extractValue(line)
        re, err := regexp.Compile(strings.Trim(val, "`"))
        if err != nil {
            return fmt.Errorf("invalid message regex in pattern %s: %w", p.ID, err)
        }
        p.MessageRe = re
    } else if strings.HasPrefix(line, "- **Context Regex:**") {
        val := extractValue(line)
        re, err := regexp.Compile(strings.Trim(val, "`"))
        if err != nil {
            return fmt.Errorf("invalid context regex in pattern %s: %w", p.ID, err)
        }
        p.ContextRe = re
    }
    return nil
}

func parseFix(p *Pattern, line string) {
    // Initialize fix if needed
    if len(p.Fixes) == 0 {
        p.Fixes = []FixTemplate{{}}
    }
    fix := &p.Fixes[0]

    if strings.HasPrefix(line, "- **File:**") {
        fix.File = strings.Trim(extractValue(line), "`")
    } else if strings.HasPrefix(line, "- **Find:**") {
        fix.FindPattern = strings.Trim(extractValue(line), "`")
    } else if strings.HasPrefix(line, "- **Replace:**") {
        fix.Replace = strings.Trim(extractValue(line), "`")
    } else if strings.HasPrefix(line, "- **Is Regex:**") {
        fix.IsRegex = extractValue(line) == "true"
    }
}

func extractValue(line string) string {
    parts := strings.SplitN(line, ":**", 2)
    if len(parts) == 2 {
        return strings.TrimSpace(parts[1])
    }
    return ""
}
```

### 4.4 Knowledge Base API

```go
// internal/evolution/knowledge/knowledge.go

package knowledge

import (
    "os"
    "path/filepath"
    "sync"

    "github.com/livetemplate/lvt/internal/evolution/types"
)

type KnowledgeBase struct {
    patterns    map[string]*Pattern
    patternsFile string
    mu          sync.RWMutex
}

func NewKnowledgeBase() *KnowledgeBase {
    kb := &KnowledgeBase{
        patterns:     make(map[string]*Pattern),
        patternsFile: findPatternsFile(),
    }
    kb.Load()
    return kb
}

func findPatternsFile() string {
    // Look for evolution/patterns.md in repo root
    // or $XDG_CONFIG_HOME/lvt/patterns.md for user patterns
    candidates := []string{
        "evolution/patterns.md",
        filepath.Join(os.Getenv("HOME"), ".config/lvt/patterns.md"),
    }
    for _, c := range candidates {
        if _, err := os.Stat(c); err == nil {
            return c
        }
    }
    return "evolution/patterns.md"
}

func (kb *KnowledgeBase) Load() error {
    kb.mu.Lock()
    defer kb.mu.Unlock()

    patterns, err := ParsePatternsFile(kb.patternsFile)
    if err != nil {
        return err
    }

    kb.patterns = make(map[string]*Pattern)
    for _, p := range patterns {
        kb.patterns[p.ID] = p
    }
    return nil
}

// LookupFixes finds applicable fixes for an error
func (kb *KnowledgeBase) LookupFixes(err types.GenerationError) []types.Fix {
    kb.mu.RLock()
    defer kb.mu.RUnlock()

    var fixes []types.Fix
    for _, pattern := range kb.patterns {
        if pattern.Matches(err) {
            for _, tmpl := range pattern.Fixes {
                fix := types.Fix{
                    ID:          generateFixID(),
                    PatternID:   pattern.ID,
                    TargetFile:  tmpl.File,
                    FindPattern: tmpl.FindPattern,
                    Replace:     tmpl.Replace,
                    IsRegex:     tmpl.IsRegex,
                    Confidence:  pattern.Confidence,
                    Rationale:   pattern.Description,
                }
                fixes = append(fixes, fix)
            }
        }
    }
    return fixes
}

func (p *Pattern) Matches(err types.GenerationError) bool {
    if p.ErrorPhase != "" && p.ErrorPhase != err.Phase {
        return false
    }
    if p.MessageRe != nil && !p.MessageRe.MatchString(err.Message) {
        return false
    }
    if p.ContextRe != nil && !p.ContextRe.MatchString(err.Context) {
        return false
    }
    return true
}

// ListPatterns returns all patterns for display
func (kb *KnowledgeBase) ListPatterns() []*Pattern {
    kb.mu.RLock()
    defer kb.mu.RUnlock()

    patterns := make([]*Pattern, 0, len(kb.patterns))
    for _, p := range kb.patterns {
        patterns = append(patterns, p)
    }
    return patterns
}
```

### 4.5 Updating Pattern Stats

When a fix is applied, update the pattern's stats in the markdown file:

```go
// internal/evolution/knowledge/updater.go

func (kb *KnowledgeBase) RecordFixApplication(patternID string, success bool) error {
    kb.mu.Lock()
    defer kb.mu.Unlock()

    pattern, ok := kb.patterns[patternID]
    if !ok {
        return fmt.Errorf("pattern not found: %s", patternID)
    }

    // Update stats
    oldCount := pattern.FixCount
    oldRate := pattern.SuccessRate

    pattern.FixCount++
    if oldCount == 0 {
        if success {
            pattern.SuccessRate = 1.0
        } else {
            pattern.SuccessRate = 0.0
        }
    } else {
        // Weighted average
        totalSuccess := oldRate * float64(oldCount)
        if success {
            totalSuccess++
        }
        pattern.SuccessRate = totalSuccess / float64(pattern.FixCount)
    }

    // Update the markdown file
    return kb.updatePatternsFile()
}

func (kb *KnowledgeBase) updatePatternsFile() error {
    // Re-generate the patterns.md file with updated stats
    // This keeps the file as source of truth
    return WritePatternsFile(kb.patternsFile, kb.ListPatterns())
}
```

### 4.6 Benefits of Markdown Storage

1. **Git Tracked** - Full history of pattern changes, who added what
2. **Human Readable** - Easy to review and understand patterns
3. **LLM Editable** - Evolution system can propose new patterns via PR
4. **PR Reviewable** - Pattern additions go through normal code review
5. **No Recompile** - Add patterns without rebuilding lvt binary
6. **Portable** - Users can copy/share pattern files
7. **Debuggable** - Easy to see why a fix was proposed

---

## 5. Fix Proposer

### 5.1 Proposal Logic

```go
// internal/evolution/proposer.go

package evolution

import (
    "context"

    "github.com/livetemplate/lvt/internal/evolution/types"
)

type FixProposer struct {
    knowledge *KnowledgeBase
    llm       *LLMClient
}

func NewFixProposer(kb *KnowledgeBase, llm *LLMClient) *FixProposer {
    return &FixProposer{
        knowledge: kb,
        llm:       llm,
    }
}

// ProposeFixesFor analyzes an event and proposes fixes
func (p *FixProposer) ProposeFixesFor(event *types.GenerationEvent) []types.Fix {
    var allFixes []types.Fix

    for _, err := range event.Errors {
        // 1. Check knowledge base first
        knownFixes := p.knowledge.LookupFixes(err)
        if len(knownFixes) > 0 {
            allFixes = append(allFixes, knownFixes...)
            continue
        }

        // 2. For unknown errors, use LLM
        if p.llm != nil {
            llmFixes := p.proposeLLMFix(err, event)
            allFixes = append(allFixes, llmFixes...)
        }
    }

    // Deduplicate and prioritize
    return deduplicateFixes(allFixes)
}

func (p *FixProposer) proposeLLMFix(err types.GenerationError, event *types.GenerationEvent) []types.Fix {
    prompt := p.buildFixPrompt(err, event)

    response, llmErr := p.llm.Generate(context.Background(), prompt)
    if llmErr != nil {
        return nil
    }

    fixes := parseFixResponse(response)

    // Mark LLM-generated fixes with lower confidence
    for i := range fixes {
        fixes[i].Confidence = 0.6  // Lower than knowledge base
        fixes[i].Source = "llm"
    }

    return fixes
}

func (p *FixProposer) buildFixPrompt(err types.GenerationError, event *types.GenerationEvent) string {
    return fmt.Sprintf(`You are an expert at fixing LiveTemplate code generation issues.

## Error Information
- **Phase**: %s
- **File**: %s
- **Line**: %d
- **Message**: %s

## Context
%s

## Generation Context
- **Command**: %s
- **Kit**: %s
- **LiveTemplate Version**: %s

## Instructions
Analyze this error and propose a specific code fix. Consider:
1. Is this a template syntax issue?
2. Is this a type mismatch?
3. Is this a missing import?
4. Is this a state management issue?

## Output Format
Respond with a JSON object:
{
  "file": "relative/path/to/file.go",
  "find": "exact text to find",
  "replace": "text to replace with",
  "rationale": "explanation of why this fix works",
  "is_regex": false
}

If multiple fixes are needed, return an array of objects.
`, err.Phase, err.File, err.Line, err.Message, err.Context,
        event.Command, event.Kit, event.LiveTemplateVer)
}
```

---

## 6. Fix Tester

### 6.1 Isolated Testing Environment

```go
// internal/evolution/tester.go

package evolution

import (
    "os"
    "os/exec"
    "path/filepath"

    "github.com/livetemplate/lvt/internal/evolution/types"
)

type FixTester struct {
    baseDir   string
    validator *ValidationEngine
}

func NewFixTester() *FixTester {
    return &FixTester{
        baseDir:   filepath.Join(os.TempDir(), "lvt-evolution"),
        validator: NewValidationEngine(),
    }
}

// TestFix applies a fix in isolation and validates the result
func (t *FixTester) TestFix(fix types.Fix, originalEvent *types.GenerationEvent) *types.TestResult {
    // Create isolated environment
    testID := generateTestID()
    testDir := filepath.Join(t.baseDir, testID)

    if err := os.MkdirAll(testDir, 0755); err != nil {
        return &types.TestResult{
            Success: false,
            Error:   fmt.Sprintf("failed to create test dir: %v", err),
        }
    }
    defer os.RemoveAll(testDir)

    // Copy lvt binary and templates to test environment
    if err := t.setupTestEnv(testDir, fix); err != nil {
        return &types.TestResult{
            Success: false,
            Error:   fmt.Sprintf("failed to setup test env: %v", err),
        }
    }

    // Apply the fix to templates
    if err := t.applyFix(testDir, fix); err != nil {
        return &types.TestResult{
            Success: false,
            Error:   fmt.Sprintf("failed to apply fix: %v", err),
        }
    }

    // Re-run the original generation
    genResult := t.regenerate(testDir, originalEvent)
    if !genResult.Success {
        return &types.TestResult{
            Success: false,
            Error:   fmt.Sprintf("regeneration failed: %s", genResult.Error),
        }
    }

    // Validate the result
    appPath := filepath.Join(testDir, "testapp")
    validation := t.validator.Validate(appPath)

    // Run additional regression tests
    regressionResults := t.runRegressionTests(appPath)

    return &types.TestResult{
        Success:    validation.AllPassed() && regressionResults.AllPassed(),
        Validation: validation,
        Regression: regressionResults,
    }
}

func (t *FixTester) applyFix(testDir string, fix types.Fix) error {
    // Find files matching the fix pattern
    var targetFiles []string
    filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
        if matchesPattern(path, fix.TargetFile) {
            targetFiles = append(targetFiles, path)
        }
        return nil
    })

    for _, file := range targetFiles {
        content, err := os.ReadFile(file)
        if err != nil {
            return err
        }

        var newContent string
        if fix.IsRegex {
            re, err := regexp.Compile(fix.FindPattern)
            if err != nil {
                return fmt.Errorf("invalid fix regex pattern: %w", err)
            }
            newContent = re.ReplaceAllString(string(content), fix.Replace)
        } else {
            newContent = strings.Replace(string(content), fix.FindPattern, fix.Replace, -1)
        }

        if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
            return err
        }
    }

    return nil
}

func (t *FixTester) regenerate(testDir string, event *types.GenerationEvent) *GenerationResult {
    // Build command from event
    args := buildCommandArgs(event)

    cmd := exec.Command(filepath.Join(testDir, "lvt"), args...)
    cmd.Dir = testDir

    output, err := cmd.CombinedOutput()
    if err != nil {
        return &GenerationResult{
            Success: false,
            Error:   string(output),
        }
    }

    return &GenerationResult{Success: true}
}

func (t *FixTester) runRegressionTests(appPath string) *types.RegressionResults {
    // Run a suite of sanity checks
    results := &types.RegressionResults{}

    // Test 1: App builds
    buildResult := t.testBuild(appPath)
    results.Add("build", buildResult)

    // Test 2: Server starts
    startResult := t.testServerStart(appPath)
    results.Add("server_start", startResult)

    // Test 3: Basic HTTP request succeeds
    if startResult.Passed {
        httpResult := t.testBasicHTTP(appPath)
        results.Add("basic_http", httpResult)
    }

    return results
}
```

---

## 7. Review and Application

### 7.1 Review Queue

```go
// internal/evolution/reviewer.go

package evolution

import (
    "github.com/livetemplate/lvt/internal/evolution/types"
)

type FixReviewer struct {
    autoApplyThreshold float64
    queue              *ReviewQueue
    git                *GitIntegration
}

func NewFixReviewer(threshold float64) *FixReviewer {
    return &FixReviewer{
        autoApplyThreshold: threshold,  // Default: 0.95 (conservative for v1)
        queue:              NewReviewQueue(),
        git:                NewGitIntegration(),
    }
}

// ProcessTestedFix decides whether to auto-apply or queue for review
func (r *FixReviewer) ProcessTestedFix(fix types.Fix, testResult *types.TestResult) error {
    if !testResult.Success {
        // Fix failed testing - discard
        return r.recordFailedFix(fix, testResult)
    }

    if fix.Confidence >= r.autoApplyThreshold {
        // High confidence + tests pass = auto-apply
        return r.autoApply(fix, testResult)
    }

    // Queue for review
    return r.queueForReview(fix, testResult)
}

func (r *FixReviewer) autoApply(fix types.Fix, testResult *types.TestResult) error {
    // Create feature branch
    branchName := fmt.Sprintf("auto-fix/%s", fix.ID)
    if err := r.git.CreateBranch(branchName); err != nil {
        return err
    }

    // Apply fix to actual templates
    if err := r.applyFixToRepo(fix); err != nil {
        r.git.DeleteBranch(branchName)
        return err
    }

    // Commit
    commitMsg := fmt.Sprintf(`fix(evolution): %s

Auto-applied fix with confidence %.2f

Pattern: %s
Rationale: %s

Test results:
- Compilation: %v
- Templates: %v
- Regression: %v
`, fix.ID, fix.Confidence, fix.PatternID, fix.Rationale,
        testResult.Validation.Compilation.Passed,
        testResult.Validation.Templates.Passed,
        testResult.Regression.AllPassed())

    if err := r.git.Commit(commitMsg); err != nil {
        return err
    }

    // Create PR
    pr, err := r.git.CreatePR(branchName, "main", fix.ID, commitMsg)
    if err != nil {
        return err
    }

    // Log the auto-application
    r.recordAutoApply(fix, testResult, pr)

    return nil
}

func (r *FixReviewer) queueForReview(fix types.Fix, testResult *types.TestResult) error {
    item := &ReviewItem{
        Fix:        fix,
        TestResult: testResult,
        Status:     "pending",
        CreatedAt:  time.Now(),
        Reason:     fmt.Sprintf("Confidence %.2f below threshold %.2f",
                               fix.Confidence, r.autoApplyThreshold),
    }

    return r.queue.Add(item)
}
```

### 7.2 LLM Review Option

```go
// internal/evolution/llm_reviewer.go

// For fixes that need review, an LLM can act as reviewer
func (r *FixReviewer) RequestLLMReview(item *ReviewItem) (*ReviewDecision, error) {
    prompt := fmt.Sprintf(`You are reviewing a proposed fix for the LiveTemplate code generator.

## Fix Details
- **ID**: %s
- **Target File**: %s
- **Pattern**: %s
- **Confidence**: %.2f
- **Rationale**: %s

## Change
Find:
%s

Replace with:
%s

## Test Results
- Compilation: %v
- Templates: %v
- Regression: %v

## Question
Should this fix be applied? Consider:
1. Does the fix address the root cause?
2. Could it introduce regressions?
3. Is the confidence appropriate?

Respond with:
{
  "decision": "approve" | "reject" | "modify",
  "confidence_adjustment": <float between -0.3 and 0.3>,
  "rationale": "explanation",
  "suggested_modification": "if decision is modify, provide the modification"
}
`, item.Fix.ID, item.Fix.TargetFile, item.Fix.PatternID,
        item.Fix.Confidence, item.Fix.Rationale,
        item.Fix.FindPattern, item.Fix.Replace,
        item.TestResult.Validation.Compilation.Passed,
        item.TestResult.Validation.Templates.Passed,
        item.TestResult.Regression.AllPassed())

    response, err := r.llm.Generate(context.Background(), prompt)
    if err != nil {
        return nil, err
    }

    return parseReviewDecision(response)
}
```

### 7.3 Upstream Library Evolution

The evolution system can propose fixes not just to lvt templates, but also to upstream repos in the LiveTemplate ecosystem:

**Target Repositories:**
- `github.com/livetemplate/livetemplate` - Core Go library (session, rendering, WebSocket)
- `github.com/livetemplate/client` - Client-side JavaScript (morphdom config, reconnection)

> **Note on Components:** The components library lives inside lvt as a nested module at `github.com/livetemplate/lvt/components` (monorepo approach). Component fixes are handled directly in lvt PRs, not as upstream fixes. See [COMPONENTS_INTEGRATION_STRATEGY.md](./COMPONENTS_INTEGRATION_STRATEGY.md) for details.

```go
// internal/evolution/upstream.go

package evolution

type UpstreamFix struct {
    Fix                            // Embed standard fix
    UpstreamRepo    string         // Target repo
    UpstreamBranch  string         // Branch to base PR on (usually "main")
    RequiresRelease bool           // Whether fix needs upstream release to take effect
    MinVersion      string         // If RequiresRelease, what version constraint to add
}

// UpstreamProposer creates PRs to upstream repos
type UpstreamProposer struct {
    git       *GitClient
    gh        *GitHubClient
    knowledge *Knowledge
}

func (p *UpstreamProposer) ProposeUpstreamFix(pattern *Pattern, event types.GenerationEvent) (*UpstreamFix, error) {
    if pattern.UpstreamRepo == "" {
        return nil, nil // Not an upstream pattern
    }

    fix := &UpstreamFix{
        Fix: Fix{
            PatternID:   pattern.ID,
            TargetFile:  pattern.Fix.File,
            FindPattern: pattern.Fix.Find,
            Replace:     pattern.Fix.Replace,
            IsRegex:     pattern.Fix.IsRegex,
            Confidence:  pattern.Confidence,
            Rationale:   fmt.Sprintf("Fix pattern %s: %s", pattern.ID, pattern.Name),
        },
        UpstreamRepo:    pattern.UpstreamRepo,
        UpstreamBranch:  "main",
        RequiresRelease: true,
    }

    return fix, nil
}

func (p *UpstreamProposer) CreateUpstreamPR(fix *UpstreamFix) (*PullRequest, error) {
    // Clone or update upstream repo
    repoPath, err := p.git.EnsureRepo(fix.UpstreamRepo)
    if err != nil {
        return nil, err
    }

    // Create branch
    branchName := fmt.Sprintf("evolution/%s", fix.PatternID)
    if err := p.git.CreateBranch(repoPath, branchName); err != nil {
        return nil, err
    }

    // Apply fix
    if err := applyFix(repoPath, &fix.Fix); err != nil {
        return nil, err
    }

    // Run upstream tests
    if err := p.runUpstreamTests(repoPath); err != nil {
        return nil, fmt.Errorf("upstream tests failed: %w", err)
    }

    // Create PR
    body := fmt.Sprintf(`## Auto-generated by lvt evolution system

This fix was identified by the lvt code generator's evolution system.

**Pattern**: %s
**Confidence**: %.2f

### Rationale
%s

### Evidence
- Error message: %s
- Context: %s

---
*This PR was created automatically. Please review carefully.*
`, fix.PatternID, fix.Confidence, fix.Rationale,
   fix.OriginalError, fix.ErrorContext)

    return p.gh.CreatePR(fix.UpstreamRepo, branchName, fix.UpstreamBranch,
        fmt.Sprintf("[evolution] %s", fix.PatternID), body)
}
```

**Workflow for Upstream Fixes:**

1. **Detection**: Pattern with `UpstreamRepo` field matches an error
2. **Validation**: Fix is tested against a local clone of the upstream repo
3. **PR Creation**: Automated PR is created with evidence and rationale
4. **Tracking**: lvt tracks the PR status and updates pattern when merged
5. **Version Bump**: After upstream release, lvt updates its go.mod dependency

```go
// After upstream fix is merged and released
func (p *UpstreamProposer) HandleUpstreamMerge(fix *UpstreamFix, newVersion string) error {
    // Update go.mod to require the new version
    if fix.RequiresRelease {
        return p.updateGoMod(fix.UpstreamRepo, newVersion)
    }
    return nil
}
```

---

## 8. Skill Evolution

### 8.1 Skill Metrics

```go
// internal/evolution/skill_metrics.go

package evolution

type SkillMetrics struct {
    SkillName        string

    // Usage
    TotalUses        int
    Last30DaysUses   int

    // Success rates
    GenerationSuccess float64  // Apps that generated without error
    CompilationSuccess float64 // Apps that compiled
    RuntimeSuccess     float64 // Apps that ran

    // Deviation tracking
    DeviationRate     float64  // How often LLM deviated from skill
    DeviationTypes    map[string]int  // What kinds of deviations

    // Time metrics
    AverageSteps      float64  // Steps to complete
    AverageDuration   time.Duration

    // Error patterns
    CommonErrors      []ErrorPattern
}

func ComputeSkillMetrics(skillName string, events []types.GenerationEvent) *SkillMetrics {
    relevant := filterEventsBySkill(events, skillName)

    metrics := &SkillMetrics{
        SkillName:      skillName,
        TotalUses:      len(relevant),
        Last30DaysUses: countLast30Days(relevant),
    }

    // Compute success rates
    var genSuccess, compSuccess, runSuccess int
    for _, e := range relevant {
        if e.Success {
            genSuccess++
        }
        if e.Validation != nil && e.Validation.Compilation.Passed {
            compSuccess++
        }
        if e.Validation != nil && e.Validation.Runtime.Passed {
            runSuccess++
        }
    }

    if len(relevant) > 0 {
        metrics.GenerationSuccess = float64(genSuccess) / float64(len(relevant))
        metrics.CompilationSuccess = float64(compSuccess) / float64(len(relevant))
        metrics.RuntimeSuccess = float64(runSuccess) / float64(len(relevant))
    }

    // Compute error patterns
    metrics.CommonErrors = computeErrorPatterns(relevant)

    return metrics
}
```

### 8.2 Skill Improver

```go
// internal/evolution/skill_improver.go

package evolution

type SkillImprover struct {
    llm           *LLMClient
    metricsStore  *MetricsStore
}

func (i *SkillImprover) ProposeImprovements(skillName string) ([]SkillChange, error) {
    // Get current metrics
    metrics := i.metricsStore.GetMetrics(skillName)

    // Get recent failures
    failures := i.metricsStore.GetFailures(skillName, 50)

    // Check if improvement is needed
    if metrics.CompilationSuccess > 0.95 && metrics.DeviationRate < 0.05 {
        return nil, nil  // Skill is performing well
    }

    // Get current skill content
    skillContent, err := loadSkillContent(skillName)
    if err != nil {
        return nil, err
    }

    // Build improvement prompt
    prompt := i.buildImprovementPrompt(metrics, failures, skillContent)

    // Generate improvements
    response, err := i.llm.Generate(context.Background(), prompt)
    if err != nil {
        return nil, err
    }

    return parseSkillChanges(response)
}

func (i *SkillImprover) buildImprovementPrompt(
    metrics *SkillMetrics,
    failures []types.GenerationEvent,
    skillContent string,
) string {
    return fmt.Sprintf(`You are improving a Claude Code skill for the lvt CLI.

## Current Skill Performance
- **Name**: %s
- **Usage (30 days)**: %d
- **Generation Success**: %.1f%%
- **Compilation Success**: %.1f%%
- **Deviation Rate**: %.1f%%

## Common Errors
%s

## Recent Failures
%s

## Current Skill Content
%s

## Instructions
Propose specific improvements to this skill to:
1. Increase compilation success rate
2. Reduce LLM deviation from instructions
3. Better handle error cases
4. Add missing validation steps

Focus on:
- Adding stronger constraints (MUST, NEVER, CRITICAL)
- Adding explicit examples of what NOT to do
- Adding validation checkpoints
- Clarifying ambiguous instructions

## Output Format
Return an array of changes:
[
  {
    "section": "section name or line number",
    "type": "add" | "replace" | "delete",
    "old": "old content (for replace)",
    "new": "new content",
    "rationale": "why this improves the skill"
  }
]
`, metrics.SkillName, metrics.Last30DaysUses,
        metrics.GenerationSuccess*100, metrics.CompilationSuccess*100,
        metrics.DeviationRate*100,
        formatErrorPatterns(metrics.CommonErrors),
        formatFailures(failures),
        skillContent)
}
```

---

## 9. Database Schema

```sql
-- internal/evolution/store/schema.sql

-- Generation events
CREATE TABLE IF NOT EXISTS generation_events (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    session_id TEXT,
    command TEXT NOT NULL,
    inputs TEXT NOT NULL,  -- JSON
    skill TEXT,
    kit TEXT,
    template_versions TEXT,  -- JSON
    livetemplate_version TEXT,
    lvt_version TEXT,
    success BOOLEAN NOT NULL,
    validation TEXT,  -- JSON
    errors TEXT,  -- JSON
    duration_ms INTEGER,
    files_generated TEXT,  -- JSON
    files_modified TEXT   -- JSON
);

CREATE INDEX idx_events_timestamp ON generation_events(timestamp);
CREATE INDEX idx_events_success ON generation_events(success);
CREATE INDEX idx_events_skill ON generation_events(skill);

-- Proposed fixes
CREATE TABLE IF NOT EXISTS proposed_fixes (
    id TEXT PRIMARY KEY,
    event_id TEXT REFERENCES generation_events(id),
    pattern_id TEXT,
    target_file TEXT NOT NULL,
    find_pattern TEXT NOT NULL,
    replace_text TEXT NOT NULL,
    is_regex BOOLEAN DEFAULT FALSE,
    confidence REAL NOT NULL,
    source TEXT,  -- 'knowledge_base' or 'llm'
    rationale TEXT,
    status TEXT DEFAULT 'proposed',  -- proposed, testing, tested, applied, rejected
    test_result TEXT,  -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    applied_at DATETIME
);

CREATE INDEX idx_fixes_status ON proposed_fixes(status);
CREATE INDEX idx_fixes_pattern ON proposed_fixes(pattern_id);

-- Knowledge base patterns
CREATE TABLE IF NOT EXISTS patterns (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    error_phase TEXT,
    message_regex TEXT,
    context_regex TEXT,
    fixes TEXT NOT NULL,  -- JSON array of FixTemplate
    confidence REAL NOT NULL,
    fix_count INTEGER DEFAULT 0,
    success_rate REAL DEFAULT 0.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Skill metrics
CREATE TABLE IF NOT EXISTS skill_metrics (
    skill_name TEXT PRIMARY KEY,
    total_uses INTEGER DEFAULT 0,
    generation_success_rate REAL DEFAULT 0.0,
    compilation_success_rate REAL DEFAULT 0.0,
    runtime_success_rate REAL DEFAULT 0.0,
    deviation_rate REAL DEFAULT 0.0,
    common_errors TEXT,  -- JSON
    last_computed DATETIME
);

-- Review queue
CREATE TABLE IF NOT EXISTS review_queue (
    id TEXT PRIMARY KEY,
    fix_id TEXT REFERENCES proposed_fixes(id),
    status TEXT DEFAULT 'pending',  -- pending, approved, rejected, modified
    reason TEXT,
    reviewer TEXT,  -- 'human' or 'llm'
    review_result TEXT,  -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    reviewed_at DATETIME
);
```

---

## 10. CLI Commands

### 10.1 Evolution Commands

```go
// Add to main.go

case "evolution":
    if len(os.Args) < 3 {
        showEvolutionHelp()
        return
    }
    switch os.Args[2] {
    case "status":
        // Show evolution system status
        commands.EvolutionStatus()
    case "metrics":
        // Show skill and template metrics
        commands.EvolutionMetrics(os.Args[3:])
    case "review":
        // Review pending fixes
        commands.EvolutionReview(os.Args[3:])
    case "apply":
        // Apply a specific fix
        commands.EvolutionApply(os.Args[3:])
    case "learn":
        // Add a pattern to knowledge base
        commands.EvolutionLearn(os.Args[3:])
    case "analyze":
        // Analyze recent failures
        commands.EvolutionAnalyze(os.Args[3:])
    default:
        showEvolutionHelp()
    }
```

### 10.2 Command Implementations

```go
// commands/evolution.go

func EvolutionStatus() {
    store := evolution.NewEventStore()

    // Get stats
    total := store.CountEvents()
    successes := store.CountSuccessful()
    failures := store.CountFailed()

    // Get pending fixes
    pendingFixes := store.CountPendingFixes()

    // Get recent improvements
    recentFixes := store.GetRecentAppliedFixes(7)

    fmt.Printf("Evolution System Status\n")
    fmt.Printf("=======================\n\n")
    fmt.Printf("Generation Events (last 30 days):\n")
    fmt.Printf("  Total:     %d\n", total)
    fmt.Printf("  Successes: %d (%.1f%%)\n", successes, float64(successes)/float64(total)*100)
    fmt.Printf("  Failures:  %d (%.1f%%)\n", failures, float64(failures)/float64(total)*100)
    fmt.Printf("\n")
    fmt.Printf("Pending Fixes: %d\n", pendingFixes)
    fmt.Printf("Applied This Week: %d\n", len(recentFixes))
}

func EvolutionMetrics(args []string) {
    // Show skill or template metrics
    if len(args) > 0 {
        skillName := args[0]
        metrics := evolution.GetSkillMetrics(skillName)
        displaySkillMetrics(metrics)
    } else {
        // Show all skills
        allMetrics := evolution.GetAllSkillMetrics()
        displayMetricsTable(allMetrics)
    }
}

func EvolutionReview(args []string) {
    queue := evolution.NewReviewQueue()
    items := queue.GetPending()

    if len(items) == 0 {
        fmt.Println("No fixes pending review")
        return
    }

    for _, item := range items {
        displayReviewItem(item)

        fmt.Print("Action [a]pprove, [r]eject, [s]kip, [l]lm review: ")
        var action string
        fmt.Scanln(&action)

        switch action {
        case "a":
            queue.Approve(item.ID)
        case "r":
            queue.Reject(item.ID)
        case "l":
            decision := queue.RequestLLMReview(item)
            displayLLMDecision(decision)
        }
    }
}
```

---

## 11. Configuration

### 11.1 Evolution Config

```yaml
# $XDG_CONFIG_HOME/lvt/evolution.yaml

enabled: true

# Telemetry settings
telemetry:
  capture_all: true
  retention_days: 90

# Validation settings
validation:
  run_compilation: true
  run_templates: true
  run_migrations: true
  timeout_seconds: 60

# Fix proposal settings
fix_proposal:
  use_llm: true
  llm_model: "claude-3-5-sonnet"
  max_llm_calls_per_day: 100

# Auto-apply settings
# IMPORTANT: For v1, recommend enabled: false until system is proven
auto_apply:
  enabled: false  # Start with human review only
  confidence_threshold: 0.95  # Conservative threshold when enabled
  require_all_tests_pass: true
  create_pr: true  # Always create PR, never direct commit
  auto_merge: false  # Require human approval to merge

# Skill improvement settings
skill_improvement:
  enabled: true
  check_interval_hours: 24
  improvement_threshold: 0.85  # Below this success rate, propose improvements
```

---

## 12. Security Considerations

### 12.1 Threat Model

The evolution system introduces several security considerations:

| Threat | Risk Level | Mitigation |
|--------|------------|------------|
| Malicious patterns in knowledge base | Medium | Pattern PRs require human review; patterns.md is git-tracked |
| LLM-generated fixes introduce vulnerabilities | High | All fixes must pass validation; auto-apply disabled by default |
| Telemetry data exposure | Medium | Local SQLite storage; no external transmission without consent |
| Code injection via regex patterns | Medium | Validate regex patterns at parse time; sandbox fix testing |
| Privilege escalation in fix application | Low | Fixes only modify template files; no system commands |

### 12.2 Security Checklist

Before enabling auto-apply or LLM-generated fixes:

- [ ] **Pattern Review**: All patterns in `evolution/patterns.md` reviewed by security team
- [ ] **Sandbox Testing**: Fix tester runs in isolated environment (temp directory, no network)
- [ ] **Input Validation**: All user inputs and pattern regexes validated and bounded
- [ ] **Output Sanitization**: Generated code scanned for common vulnerabilities (SQLi, XSS)
- [ ] **Audit Logging**: All fix applications logged with before/after diffs
- [ ] **Rollback Capability**: Every auto-applied fix can be reverted via git
- [ ] **Rate Limiting**: LLM API calls rate-limited to prevent abuse/cost overrun
- [ ] **Secrets Handling**: Telemetry excludes sensitive data (env vars, credentials)

### 12.3 Recommended Security Posture by Phase

| Phase | Auto-Apply | LLM Fixes | Human Review |
|-------|------------|-----------|--------------|
| v1.0 (Initial) | Disabled | Disabled | All fixes |
| v1.1 (Validated) | High-confidence only (≥0.98) | Disabled | Medium/low confidence |
| v2.0 (Mature) | ≥0.95 confidence | Enabled with review | Low confidence only |

---

## 13. Migration Strategy

### 13.1 For Existing lvt Users

Users with existing lvt-generated apps can adopt the evolution system gradually:

**Phase 1: Opt-in Telemetry**
```bash
# Enable telemetry to start collecting data
export LVT_TELEMETRY=true
lvt gen resource ...
```

**Phase 2: Validation Only**
```bash
# Enable validation without evolution
lvt gen resource --validate
```

**Phase 3: Evolution Monitoring**
```bash
# View evolution metrics without applying fixes
lvt evolution status
lvt evolution failures --last 30d
```

**Phase 4: Reviewed Fixes**
```bash
# Review and apply fixes manually
lvt evolution propose <event-id>
lvt evolution apply <fix-id> --dry-run
lvt evolution apply <fix-id>
```

### 13.2 Breaking Changes

The evolution system introduces no breaking changes to the lvt CLI. All new features are:
- **Opt-in**: Telemetry disabled by default
- **Additive**: New commands don't affect existing workflows
- **Backward Compatible**: Generated apps unchanged unless fixes applied

### 13.3 Rollback Procedure

If a fix causes issues:

```bash
# View recent auto-applied fixes
lvt evolution history --auto-applied

# Revert a specific fix
git revert <commit-hash>

# Or restore from backup
lvt evolution restore --before <fix-id>
```

---

## 14. Success Metrics

### 14.1 Baseline Metrics (Current State)

Before implementing the evolution system, establish baselines:

| Metric | Current (Estimated) | Target |
|--------|---------------------|--------|
| Generation success rate | ~60% | 95% |
| Compilation success rate | ~75% | 99% |
| Template parse success | ~85% | 100% |
| Time to fix known issues | Days-weeks | Hours |
| Pattern coverage | 0 patterns | 50+ patterns |

### 14.2 Key Performance Indicators (KPIs)

Track these metrics weekly:

```sql
-- Generation success rate (last 7 days)
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN success THEN 1 ELSE 0 END) as successful,
    ROUND(100.0 * SUM(CASE WHEN success THEN 1 ELSE 0 END) / COUNT(*), 2) as success_rate
FROM generation_events
WHERE timestamp > datetime('now', '-7 days');

-- Fix effectiveness
SELECT
    pattern_id,
    COUNT(*) as fix_count,
    SUM(CASE WHEN status = 'applied' THEN 1 ELSE 0 END) as applied,
    ROUND(100.0 * SUM(CASE WHEN status = 'applied' THEN 1 ELSE 0 END) / COUNT(*), 2) as apply_rate
FROM proposed_fixes
GROUP BY pattern_id;
```

### 14.3 Success Criteria by Milestone

| Milestone | Success Criteria |
|-----------|------------------|
| 1. Stop the Bleeding | All E2E tests include compilation check; 0 SKIP_GO_MOD_TIDY |
| 2. Validation Layer | 100% of generations validated; <100ms validation overhead |
| 3. Telemetry & Evolution | 10+ patterns with >80% fix rate; telemetry captures all failures |
| 4. Components Integration | 0 template drift; modal/toast/dropdown working |
| 5. Style System | Tailwind and unstyled adapters complete; style switching works |
| 6. Components Evolution | Component health dashboard shows >95% success rate |

### 14.4 Long-term Goals

- **6 months**: 90% generation success rate, 50 patterns, auto-apply enabled
- **12 months**: 95% success rate, 100 patterns, LLM fixes enabled with review
- **18 months**: 98% success rate, pattern suggestions from community

---

## Summary

This evolution system creates a self-improving feedback loop:

1. **Capture** - Every generation attempt is recorded with full context
2. **Validate** - Generated code is validated (compilation, templates, migrations)
3. **Analyze** - Failures are analyzed against knowledge base patterns
4. **Propose** - Fixes are proposed (from patterns or LLM)
5. **Test** - Fixes are tested in isolation
6. **Review** - High-confidence fixes auto-applied; others queued
7. **Learn** - Successful fixes become new patterns

This transforms lvt from a static tool into an evolving system that improves with every use, getting closer to deterministic, one-shot generation of working apps.
