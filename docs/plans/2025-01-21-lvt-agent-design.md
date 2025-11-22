# LVT Skill-Aware Agent Design

**Date:** 2025-01-21
**Status:** Design Complete
**Version:** 1.0.0

## Executive Summary

This document describes the design for a markdown-based AI agent (`lvt-assistant`) that intelligently orchestrates LiveTemplate's 19 Claude Code skills. The agent provides smart routing, workflow chaining, persistent progress tracking, and resumability, making LVT development accessible to beginners while accelerating power users.

**Key Value Propositions:**
- **Beginners:** Conversational guidance and automatic skill discovery
- **Power Users:** Fast multi-skill workflow automation
- **Everyone:** Error prevention and best practices enforcement

## Background

### Current State

The `add-claude-skills` branch contains:
- **19 production-ready skills** organized in 5 categories:
  - Core (14 skills): new-app, add-resource, add-view, add-migration, gen-auth, gen-schema, resource-inspect, manage-kits, validate-templates, run-and-test, customize, seed-data, deploy, manage-env
  - Workflows (3 skills): quickstart, production-ready, add-related-resources
  - Maintenance (3 skills): analyze, suggest, troubleshoot
  - Meta (1 skill): add-skill
- **100% test coverage** with chromedp e2e tests
- **Comprehensive documentation** (README, TESTING, SKILL_DEVELOPMENT guides)

### Problem Statement

While skills are powerful, users face challenges:
1. **Discovery:** Hard to know which of 19 skills to use for a given task
2. **Complexity:** Multi-step workflows require manual skill chaining
3. **Context Loss:** No persistence between Claude Code sessions
4. **Learning Curve:** New users must read extensive documentation

### Solution

Create an AI agent that:
1. Analyzes user intent and routes to appropriate skills
2. Automatically chains skills for complex workflows
3. Tracks progress persistently for resumability
4. Provides interactive guidance through natural language

## Architecture

### High-Level Design

```
User Request
    ↓
Claude Code
    ↓
lvt-assistant Agent
    ↓
┌─────────────────────────────────┐
│ Intent Classification           │
│ (analyze user request)          │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ Skill Selection & Routing       │
│ (map to 1+ skills)              │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ Execution Strategy              │
│ (direct, guided, or interactive)│
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ Progress Tracking               │
│ (TodoWrite + markdown file)     │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ Validation & Error Handling     │
│ (verify outputs, suggest fixes) │
└─────────────────────────────────┘
    ↓
Result + Next Steps
```

### Component Architecture

```
lvt Binary
│
├─ internal/claude/
│  ├─ embed.go          (embeds skills + agent via embed.FS)
│  ├─ install.go        (extraction logic)
│  └─ version.go        (version management)
│
├─ cmd/lvt/commands/
│  ├─ new.go            (calls InstallClaudeResources)
│  └─ install.go        (NEW: manual install/upgrade commands)
│
└─ Embedded Resources:
   ├─ .claude/agents/lvt-assistant.md
   ├─ .claude/skills/lvt/core/*.md (14 skills)
   ├─ .claude/skills/lvt/workflows/*.md (3 skills)
   ├─ .claude/skills/lvt/maintenance/*.md (3 skills)
   └─ .claude/skills/lvt/meta/*.md (1 skill)

User Project (after lvt new)
│
└─ .claude/
   ├─ agents/
   │  └─ lvt-assistant.md
   ├─ skills/lvt/
   │  ├─ core/*.md
   │  ├─ workflows/*.md
   │  ├─ maintenance/*.md
   │  └─ meta/*.md
   ├─ progress/          (gitignored, persistent state)
   │  └─ YYYY-MM-DD-HH-MM-SS-task-name.md
   └─ VERSION
```

## Agent Design

### Agent File Structure

**Location:** `.claude/agents/lvt-assistant.md`

**Format:**
```markdown
---
name: lvt-assistant
description: Intelligent orchestrator for LiveTemplate development
color: purple
version: 1.0.0
---

# LVT Assistant Agent

[Agent prompt content organized in sections...]
```

**Content Sections:**

1. **Agent Capabilities**
   - Overview of agent's role
   - Catalog of all 19 skills with brief descriptions
   - Routing and chaining capabilities

2. **Context Awareness**
   - How to read project state (lvt.yaml, schema.sql, migrations/)
   - Pre-execution validation checklist
   - Detection of existing resources

3. **Intent Classification**
   - Decision tree for mapping requests to skills
   - Pattern matching examples
   - Ambiguity resolution strategies

4. **Execution Protocol**
   - Direct execution (simple, single-command tasks)
   - Guided execution (multi-step workflows)
   - Interactive execution (ambiguous requests)
   - Progress tracking system

5. **Skill Invocation Templates**
   - For each of 19 skills:
     - Trigger patterns
     - Prerequisites
     - Execution template
     - Validation steps
     - Common errors and fixes

6. **Workflow Chains**
   - Pre-defined sequences:
     - Quickstart: new-app → quickstart → run-and-test
     - Full Auth: gen-auth → seed-data → run-and-test
     - Production: analyze → production-ready → deploy → manage-env
     - Resource Family: add-related-resources → (add-resource ×N) → customize

7. **Resumption Logic**
   - Detect incomplete workflows
   - Parse progress files
   - Restore state
   - Offer resumption to user

8. **Error Handling**
   - Common error patterns
   - Recovery procedures
   - When to invoke troubleshoot skill

**Estimated Size:** 800-1200 lines

### Intent Classification

The agent classifies user requests into categories and routes to appropriate skills:

| Category | Keywords/Patterns | Target Skills | Example |
|----------|------------------|---------------|---------|
| **Project Init** | "new app", "start project", "create application" | new-app, quickstart | "create a blog app" |
| **Resource Management** | "add [entity]", "create model", "CRUD for [thing]" | add-resource, add-view | "add products with name and price" |
| **Authentication** | "auth", "login", "signup", "magic link" | gen-auth | "add magic link authentication" |
| **Database** | "migration", "schema", "alter table" | add-migration, gen-schema | "add created_at to users" |
| **Development** | "run", "test", "start server" | run-and-test | "run the app" |
| **Production** | "deploy", "production", "docker", "fly.io" | production-ready, deploy, manage-env | "deploy to production" |
| **Analysis** | "analyze", "check", "problem", "slow" | analyze, troubleshoot, suggest | "why is my app slow?" |
| **Customization** | "customize", "modify template", "change style" | customize, manage-kits | "change to Tailwind CSS" |
| **Inspection** | "show resources", "what models", "list tables" | resource-inspect | "what resources do I have?" |

**Ambiguity Resolution:**

When intent is unclear, agent asks targeted questions using AskUserQuestion tool:
- "add users" → "Do you want authentication too?" (yes → gen-auth, no → add-resource)
- "deploy" → "Where?" (Docker/Fly.io/K8s/VPS options)
- "add blog" → "With comments and tags?" (suggests related resources)

### Execution Strategies

The agent uses three execution modes based on task complexity:

#### 1. Direct Execution

**When:** Single-skill, clear parameters, no validation needed

**How:** Agent directly runs lvt CLI commands via Bash tool

**Examples:**
```
User: "add a products resource with name and price"
Agent: lvt gen products name price:float

User: "run the app"
Agent: PORT=8080 timeout 30s go run main.go
```

#### 2. Guided Execution

**When:** Multi-step workflows, validation between steps

**How:** Agent invokes skill, follows checklist, validates outputs

**Examples:**
```
User: "set up authentication with magic links"
Agent:
  1. Read lvt:gen-auth skill
  2. Run lvt gen auth --magic-links
  3. Validate email templates created
  4. Check SMTP config in lvt.yaml
  5. Run tests
  6. Report success + next steps

User: "make this production-ready"
Agent:
  1. Run lvt:analyze
  2. Review issues
  3. Run lvt:production-ready
  4. Apply fixes
  5. Run lvt:deploy
  6. Validate deployment
```

#### 3. Interactive Execution

**When:** Ambiguous request, missing information

**How:** Agent asks clarifying questions, then proceeds with direct or guided

**Examples:**
```
User: "improve my app"
Agent: (AskUserQuestion) "What aspect?"
  - Performance
  - Security
  - UX
  - Code Quality
User selects → Agent routes to appropriate skill

User: "add products"
Agent: "What fields should products have?"
User responds → Agent runs lvt gen products [fields]
```

**Decision Matrix:**

| Task Complexity | Command Count | Validation Needed | Mode |
|----------------|---------------|-------------------|------|
| Simple | 1 | No | Direct |
| Moderate | 1-3 | Yes | Guided |
| Complex | 4+ | Yes | Guided + Interactive |
| Unclear | Any | N/A | Interactive → Direct/Guided |

## Progress Tracking System

### Dual-Track Architecture

The agent maintains progress in two places for different purposes:

**1. TodoWrite (Real-time UI)**
- Updates Claude Code interface during execution
- Shows current step status (pending/in_progress/completed)
- Provides immediate feedback to users
- Ephemeral (lost between sessions)

**2. Markdown Files (Persistent State)**
- Written to `.claude/progress/YYYY-MM-DD-HH-MM-SS-task-name.md`
- Survives session restarts
- Auditable history
- Enables workflow resumption

### Progress File Format

```markdown
---
workflow: production-deployment
started: 2025-01-21T14:30:00Z
status: in_progress
user_request: "deploy my app to production"
agent_version: 1.0.0
---

# Production Deployment Workflow

## Progress
- [x] Step 1: Analyze app (lvt:analyze) - 14:30:15 ✓
- [x] Step 2: Run production-ready checks (lvt:production-ready) - 14:32:40 ✓
- [ ] Step 3: Deploy to Fly.io (lvt:deploy) - IN PROGRESS
- [ ] Step 4: Verify deployment
- [ ] Step 5: Update environment variables (lvt:manage-env)

## Execution Log

### Step 1: Analyze app (14:30:15 - 14:32:10)
**Command:** `lvt analyze`
**Output:**
```
[analyzer output]
```
**Validation:** ✓ No critical issues found

### Step 3: Deploy to Fly.io (14:35:00 - ONGOING)
**Command:** `lvt deploy --platform fly.io`
**Status:** Building Docker image...

## State
current_step: 3
platform: fly.io
project_root: /path/to/project
```

### Synchronization Strategy

Agent keeps both in sync:

1. **Start workflow:**
   - Create TodoWrite todos
   - Create progress markdown file
   - Write metadata to both

2. **During execution:**
   - Update TodoWrite status (in_progress → completed)
   - Append to progress file execution log
   - Include timestamps, commands, outputs

3. **On completion:**
   - Mark all TodoWrite todos completed
   - Update progress file status to "completed"
   - Add completion timestamp

4. **On error:**
   - Keep TodoWrite todo as in_progress
   - Log error in progress file
   - Create new todo for recovery action

### Resumption Logic

**On agent startup:**

```
1. Check for incomplete progress files
   ls .claude/progress/*.md | grep 'status: in_progress'

2. If found, parse most recent file:
   - Extract workflow name
   - Find last completed step
   - Load state metadata

3. Ask user via AskUserQuestion:
   "Resume [workflow name] from step [N]?"
   - Yes → Restore TodoWrite state → Continue from checkpoint
   - No → Archive old file → Start fresh

4. If resuming:
   - Read progress file state section
   - Recreate TodoWrite with correct statuses
   - Continue execution from next pending step
```

**Archive strategy:**
- On new workflow start, move old incomplete files to `.claude/progress/archive/`
- Keep last 10 progress files, delete older

### Progress Tracking Modes

| Mode | TodoWrite | Markdown Progress | Use Case |
|------|-----------|-------------------|----------|
| **Direct** | Yes | No | Single command (e.g., "add products") |
| **Guided** | Yes | Yes | Multi-step 2-5 steps (e.g., "set up auth") |
| **Complex** | Yes | Yes (detailed) | Long workflows 6+ steps (e.g., "production deploy") |
| **Resumable** | Yes | Yes (checkpointed) | Multi-session tasks (e.g., "migrate 100 tables") |

## Installation & Distribution

### Embedding Strategy

**Go Embed Implementation:**

```go
// internal/claude/embed.go
package claude

import "embed"

//go:embed skills/**/*.md agents/*.md
var ClaudeResources embed.FS

const Version = "1.0.0"
```

**Directory Structure in Binary:**

```
skills/
├─ lvt/
│  ├─ core/
│  │  ├─ new-app.md
│  │  ├─ add-resource.md
│  │  └─ ... (12 more)
│  ├─ workflows/
│  │  ├─ quickstart.md
│  │  ├─ production-ready.md
│  │  └─ add-related-resources.md
│  ├─ maintenance/
│  │  ├─ analyze.md
│  │  ├─ suggest.md
│  │  └─ troubleshoot.md
│  └─ meta/
│     └─ add-skill.md
agents/
└─ lvt-assistant.md
```

### Installation Functions

**1. Extract Resources**

```go
// internal/claude/install.go
package claude

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
)

func InstallClaudeResources(projectRoot string) error {
    // Create directory structure
    dirs := []string{
        ".claude/skills/lvt/core",
        ".claude/skills/lvt/workflows",
        ".claude/skills/lvt/maintenance",
        ".claude/skills/lvt/meta",
        ".claude/agents",
        ".claude/progress",
    }

    for _, dir := range dirs {
        path := filepath.Join(projectRoot, dir)
        if err := os.MkdirAll(path, 0755); err != nil {
            return fmt.Errorf("create dir %s: %w", dir, err)
        }
    }

    // Extract embedded files
    if err := extractEmbeddedFS(projectRoot); err != nil {
        return fmt.Errorf("extract embedded files: %w", err)
    }

    // Write version metadata
    if err := writeVersionFile(projectRoot); err != nil {
        return fmt.Errorf("write version file: %w", err)
    }

    // Update .gitignore
    if err := updateGitignore(projectRoot); err != nil {
        return fmt.Errorf("update gitignore: %w", err)
    }

    return nil
}

func extractEmbeddedFS(projectRoot string) error {
    return fs.WalkDir(ClaudeResources, ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }

        // Read embedded file
        content, err := ClaudeResources.ReadFile(path)
        if err != nil {
            return err
        }

        // Write to project
        destPath := filepath.Join(projectRoot, ".claude", path)
        return os.WriteFile(destPath, content, 0644)
    })
}
```

**2. Version Management**

```go
// internal/claude/version.go
package claude

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type VersionInfo struct {
    Version      string    `yaml:"claude-resources-version"`
    SkillsHash   string    `yaml:"skills-hash"`
    AgentHash    string    `yaml:"agent-hash"`
    InstalledAt  time.Time `yaml:"installed-at"`
}

func writeVersionFile(projectRoot string) error {
    info := VersionInfo{
        Version:     Version,
        SkillsHash:  computeSkillsHash(),
        AgentHash:   computeAgentHash(),
        InstalledAt: time.Now(),
    }

    path := filepath.Join(projectRoot, ".claude", "VERSION")
    return writeYAML(path, info)
}

func CheckVersion(projectRoot string) (needsUpgrade bool, err error) {
    path := filepath.Join(projectRoot, ".claude", "VERSION")
    var installed VersionInfo
    if err := readYAML(path, &installed); err != nil {
        return false, err
    }

    return installed.Version != Version, nil
}
```

**3. .gitignore Updates**

```go
func updateGitignore(projectRoot string) error {
    gitignorePath := filepath.Join(projectRoot, ".gitignore")

    // Read existing
    content, err := os.ReadFile(gitignorePath)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    // Check if already has .claude/progress/
    if strings.Contains(string(content), ".claude/progress/") {
        return nil
    }

    // Append
    f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.WriteString("\n# Claude Code agent progress (local state)\n.claude/progress/\n")
    return err
}
```

### CLI Integration

**1. Automatic Installation (lvt new)**

```go
// cmd/lvt/commands/new.go
func runNew(cmd *cobra.Command, args []string) error {
    projectName := args[0]

    // ... existing project creation logic ...

    // Install Claude Code resources
    if err := claude.InstallClaudeResources(projectPath); err != nil {
        return fmt.Errorf("install claude resources: %w", err)
    }

    fmt.Printf("✓ Installed Claude Code skills and agent\n")

    return nil
}
```

**2. Manual Installation Commands**

```go
// cmd/lvt/commands/install.go (NEW FILE)
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "path/to/internal/claude"
)

var installCmd = &cobra.Command{
    Use:   "install [target]",
    Short: "Install optional components",
    Long:  "Install Claude Code skills, agent, or other optional components",
}

var installClaudeCmd = &cobra.Command{
    Use:   "claude-code",
    Short: "Install Claude Code skills and agent",
    Long:  "Install or update Claude Code skills and lvt-assistant agent",
    RunE:  runInstallClaude,
}

func runInstallClaude(cmd *cobra.Command, args []string) error {
    projectRoot, err := os.Getwd()
    if err != nil {
        return err
    }

    // Check if already installed
    versionPath := filepath.Join(projectRoot, ".claude", "VERSION")
    if _, err := os.Stat(versionPath); err == nil {
        // Already installed, ask about upgrade
        needsUpgrade, err := claude.CheckVersion(projectRoot)
        if err != nil {
            return err
        }
        if !needsUpgrade {
            fmt.Println("Claude Code resources are already up to date")
            return nil
        }
        fmt.Println("Upgrading Claude Code resources...")
    } else {
        fmt.Println("Installing Claude Code resources...")
    }

    if err := claude.InstallClaudeResources(projectRoot); err != nil {
        return fmt.Errorf("install failed: %w", err)
    }

    fmt.Println("✓ Claude Code skills and agent installed successfully")
    fmt.Println("\nNext steps:")
    fmt.Println("  1. Open this project in Claude Code")
    fmt.Println("  2. Try: 'help me build a blog'")
    fmt.Println("  3. The lvt-assistant agent will guide you")

    return nil
}

func init() {
    installCmd.AddCommand(installClaudeCmd)
    rootCmd.AddCommand(installCmd)
}
```

**3. Upgrade Command**

```go
var upgradeCmd = &cobra.Command{
    Use:   "upgrade [target]",
    Short: "Upgrade components to latest version",
}

var upgradeClaudeCmd = &cobra.Command{
    Use:   "claude-code",
    Short: "Upgrade Claude Code skills and agent",
    RunE:  runUpgradeClaude,
}

func runUpgradeClaude(cmd *cobra.Command, args []string) error {
    projectRoot, err := os.Getwd()
    if err != nil {
        return err
    }

    // Check current version
    versionPath := filepath.Join(projectRoot, ".claude", "VERSION")
    if _, err := os.Stat(versionPath); os.IsNotExist(err) {
        return fmt.Errorf("Claude Code resources not installed. Run: lvt install claude-code")
    }

    needsUpgrade, err := claude.CheckVersion(projectRoot)
    if err != nil {
        return err
    }

    if !needsUpgrade {
        fmt.Println("✓ Already running latest version")
        return nil
    }

    fmt.Println("Upgrading Claude Code resources...")

    // Backup current progress files
    if err := claude.BackupProgressFiles(projectRoot); err != nil {
        return err
    }

    // Reinstall (overwrites skills/agent, preserves progress)
    if err := claude.InstallClaudeResources(projectRoot); err != nil {
        return fmt.Errorf("upgrade failed: %w", err)
    }

    fmt.Println("✓ Upgraded successfully")
    return nil
}
```

### Installation Workflow

```
User runs: lvt new myapp
  ↓
lvt creates project structure
  ↓
lvt calls: claude.InstallClaudeResources("myapp")
  ↓
┌─────────────────────────────────────────┐
│ 1. Create .claude/ directory structure  │
│ 2. Extract embedded skills (19 files)   │
│ 3. Extract agent (lvt-assistant.md)     │
│ 4. Write VERSION metadata               │
│ 5. Update .gitignore                    │
└─────────────────────────────────────────┘
  ↓
User opens Claude Code in myapp/
  ↓
Claude Code auto-discovers:
  - .claude/agents/lvt-assistant.md
  - .claude/skills/lvt/**/*.md
  ↓
User types: "add authentication"
  ↓
Agent activates and orchestrates skills
```

## Testing Strategy

### Test Categories

#### 1. Agent Logic Tests (Unit)

**File:** `internal/claude/agent_test.go`

**Tests:**
- Intent classification accuracy
- Skill routing correctness
- Progress file parsing/writing
- Version detection and upgrade logic
- Resumption state restoration

**Example:**
```go
func TestIntentClassification(t *testing.T) {
    tests := []struct {
        input    string
        expected string // skill name
    }{
        {"add products resource", "lvt:add-resource"},
        {"set up authentication", "lvt:gen-auth"},
        {"deploy to production", "lvt:deploy"},
        {"my app is slow", "lvt:analyze"},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            skill := classifyIntent(tt.input)
            assert.Equal(t, tt.expected, skill)
        })
    }
}
```

#### 2. Installation Tests (Integration)

**File:** `internal/claude/install_test.go`

**Tests:**
- Embedding correctness (all files present)
- Directory structure creation
- Version file writing
- .gitignore updates
- Upgrade preserves progress files

**Example:**
```go
func TestInstallClaudeResources(t *testing.T) {
    tmpDir := t.TempDir()

    err := InstallClaudeResources(tmpDir)
    require.NoError(t, err)

    // Verify structure
    assert.DirExists(t, filepath.Join(tmpDir, ".claude/skills/lvt/core"))
    assert.DirExists(t, filepath.Join(tmpDir, ".claude/agents"))

    // Verify files
    assert.FileExists(t, filepath.Join(tmpDir, ".claude/agents/lvt-assistant.md"))
    assert.FileExists(t, filepath.Join(tmpDir, ".claude/skills/lvt/core/new-app.md"))

    // Verify version
    assert.FileExists(t, filepath.Join(tmpDir, ".claude/VERSION"))

    // Count skills
    skills, _ := filepath.Glob(filepath.Join(tmpDir, ".claude/skills/lvt/**/*.md"))
    assert.Len(t, skills, 19)
}
```

#### 3. E2E Tests (chromedp)

**File:** `e2e/agent_integration_test.go`

**Tests:**
- Simple request routing (direct execution)
- Complex workflow chaining (guided execution)
- Progress tracking (TodoWrite + markdown sync)
- Resumption after interruption
- Error recovery

**Test Scenarios:**

1. **Simple Resource Addition**
```go
func TestAgentSimpleResource(t *testing.T) {
    // 1. Create test project with agent
    // 2. Start Claude Code session
    // 3. Send: "add products with name and price"
    // 4. Verify:
    //    - Agent routes to lvt:add-resource
    //    - Command executes: lvt gen products name price:float
    //    - Files created: internal/app/products/*.go
    //    - No progress file (simple task)
}
```

2. **Authentication Workflow**
```go
func TestAgentAuthWorkflow(t *testing.T) {
    // 1. Create test project with agent
    // 2. Send: "add magic link authentication"
    // 3. Verify:
    //    - Agent routes to lvt:gen-auth
    //    - Progress file created
    //    - TodoWrite shows steps
    //    - Auth files generated
    //    - Email templates created
    //    - Tests run
    //    - Progress file marked complete
}
```

3. **Workflow Resumption**
```go
func TestAgentResumption(t *testing.T) {
    // 1. Start production deployment workflow
    // 2. Interrupt after step 2 (simulate crash)
    // 3. Verify progress file shows in_progress
    // 4. Restart agent
    // 5. Verify:
    //    - Agent detects incomplete workflow
    //    - Asks: "Resume production deployment?"
    //    - On yes, continues from step 3
    //    - TodoWrite restored correctly
}
```

4. **Error Recovery**
```go
func TestAgentErrorRecovery(t *testing.T) {
    // 1. Send request that will fail validation
    // 2. Verify:
    //    - Agent detects error
    //    - Logs error in progress file
    //    - Suggests fix via lvt:troubleshoot
    //    - Offers to retry after fix
}
```

5. **Workflow Chain**
```go
func TestAgentWorkflowChain(t *testing.T) {
    // 1. Send: "create a production-ready blog with auth"
    // 2. Verify agent chains:
    //    - lvt:new-app blog
    //    - lvt:add-resource posts title:string content:text
    //    - lvt:add-resource comments post_id:int text:string
    //    - lvt:gen-auth
    //    - lvt:seed-data
    //    - lvt:production-ready
    // 3. Verify progress file tracks all steps
    // 4. Verify TodoWrite updates in real-time
}
```

#### 4. Skill Integration Tests

**Goal:** Verify agent works with all 19 existing skills

**Approach:** Reuse test scenarios from `CLAUDE_SKILLS_TRACKER.md`:

1. **Session 1 Tests** → Run through agent
   - Test 1: new-app + add-resource + run-and-test
   - Test 2: add-migration + validate-templates
   - Test 3: gen-auth (password + magic link)

2. **Session 2 Tests** → Run through agent
   - Test 4: quickstart workflow
   - Test 5: add-related-resources + add-view

3. **Session 3 Tests** → Run through agent
   - Test 6: production-ready + deploy
   - Test 7: analyze + troubleshoot + suggest

4. **Session 4 Tests** → Run through agent
   - Test 8: customize + manage-kits
   - Test 9: seed-data + resource-inspect

**Validation:** All 9 scenarios must pass through agent with same outcomes as direct skill usage.

### Test Infrastructure

**Reuse Existing:**
- `e2e/test_helpers.go` - chromedp utilities
- `e2e/browser_console.go` - console log capture
- `e2e/websocket_monitor.go` - WS message tracking
- `scripts/test_cli_smoke.sh` - CLI validation

**New Additions:**
- `e2e/agent_helpers.go` - Agent-specific test utilities
  - `createTestProjectWithAgent()`
  - `sendAgentRequest()`
  - `verifyProgressFile()`
  - `verifyTodoWriteState()`
  - `simulateInterruption()`

### Success Criteria

Agent tests must achieve:
- ✅ 95%+ intent classification accuracy
- ✅ All 19 skills accessible through agent
- ✅ All 9 existing skill test scenarios pass
- ✅ Workflow resumption works after interruption
- ✅ Progress files correctly track state
- ✅ TodoWrite and markdown stay synchronized
- ✅ Error recovery suggests appropriate fixes
- ✅ No skill hallucination (only real skills invoked)

## User Experience

### Beginner Flow

```
User: "I want to build a task management app"

Agent: "Great! I'll help you create a task management app. Let me start by creating
        a new LiveTemplate project and adding the necessary resources.

        I'm using the lvt:quickstart skill to set this up efficiently."

[Agent creates TodoWrite todos]:
1. Create new project
2. Add tasks resource
3. Add users resource
4. Set up authentication
5. Generate seed data
6. Run and test

Agent: "Creating project 'taskmanager'..."
[Executes: lvt new taskmanager --template standard]

Agent: "Adding tasks resource with title, description, status, due_date..."
[Executes: lvt gen tasks title description status due_date:date]

Agent: "Adding users resource..."
[Executes: lvt gen users name email]

Agent: "Setting up authentication with password and magic links..."
[Invokes lvt:gen-auth skill, follows checklist]

Agent: "Generating realistic seed data..."
[Executes: lvt seed tasks users]

Agent: "Starting development server..."
[Executes: PORT=8080 go run main.go]

Agent: "✓ Your task management app is running at http://localhost:8080!

        Try logging in with test credentials:
        - Email: alice@example.com
        - Password: password123

        Next steps you might want to try:
        - 'add comments to tasks'
        - 'customize the UI'
        - 'deploy to production'"
```

### Power User Flow

```
User: "quickstart blog with auth then deploy to fly.io"

Agent: "I'll create a production-ready blog with authentication and deploy it to Fly.io.
        This will take several steps."

[Creates progress file: .claude/progress/2025-01-21-14-30-00-blog-deployment.md]
[Creates TodoWrite todos for full workflow]

Agent: [chains multiple skills without user intervention]
- lvt:new-app → lvt:quickstart
- lvt:add-resource posts, comments, tags
- lvt:gen-auth
- lvt:seed-data
- lvt:production-ready (validates, applies fixes)
- lvt:deploy --platform fly.io

[Each step logged in progress file with outputs]

Agent: "✓ Blog deployed successfully to https://my-blog.fly.dev

        Deployment summary:
        - 3 resources: posts, comments, tags
        - Authentication: password + magic links
        - Database: PostgreSQL on Fly
        - CDN: Cloudflare (for assets)
        - Monitoring: Fly metrics enabled

        Credentials saved to: .env.production

        Next: 'configure custom domain' or 'add search functionality'"
```

### Resumption Flow

```
[User starts deployment, then closes laptop]

User: [reopens Claude Code next day]

Agent: "I noticed you have an incomplete workflow from yesterday:
        'Production Deployment' (started 2025-01-21 14:30)

        Progress so far:
        ✓ Analyzed app
        ✓ Applied production-ready fixes
        ✗ Deploy to Fly.io (interrupted during Docker build)

        Would you like to resume from the deployment step?"

User: "yes"

Agent: "Resuming deployment. I'll pick up where we left off..."

[Reads progress file, restores state, continues from step 3]

Agent: [completes remaining steps]
- Rebuild Docker image
- Push to Fly.io registry
- Deploy containers
- Run migrations
- Verify deployment

Agent: "✓ Deployment complete! Your app is live at https://myapp.fly.dev"
```

### Error Handling Flow

```
User: "deploy to production"

Agent: "I'll deploy your app to production. First, let me analyze it..."

[Runs lvt:analyze]

Agent: "⚠ Found 3 issues that need to be fixed before deployment:

        1. Missing environment variables: DATABASE_URL, SMTP_HOST
        2. No database migrations in production mode
        3. Static assets not optimized

        I recommend:
        1. Use 'lvt:manage-env' to set up environment variables
        2. Run 'lvt:production-ready' to fix these issues automatically

        Would you like me to run production-ready checks and fixes?"

User: "yes"

Agent: "Running production-ready checks..."

[Invokes lvt:production-ready, applies fixes, validates]

Agent: "✓ All issues fixed! Ready to deploy.

        Proceeding with deployment..."

[Continues with deployment workflow]
```

## Implementation Plan

### Phase 1: Agent Creation (Week 1)

**Tasks:**
1. Write `lvt-assistant.md` agent file
   - Front matter and metadata
   - Agent capabilities section
   - Intent classification logic
   - All 19 skill invocation templates
   - Workflow chain definitions
   - Progress tracking protocol
   - Resumption logic
   - Error handling procedures

2. Create example progress file templates

3. Write agent documentation
   - User guide: How to interact with agent
   - Developer guide: How to extend agent with new skills

**Deliverables:**
- `.claude/agents/lvt-assistant.md` (complete)
- `docs/LVT_AGENT_USER_GUIDE.md`
- `docs/LVT_AGENT_DEVELOPER_GUIDE.md`

**Success Criteria:**
- Agent markdown is valid (parseable front matter)
- All 19 skills referenced correctly
- Workflow chains are logically sound
- Progress tracking protocol is clear

### Phase 2: Go Infrastructure (Week 1-2)

**Tasks:**
1. Create `internal/claude/` package
   - `embed.go` - Embed agent + skills
   - `install.go` - Extraction logic
   - `version.go` - Version management

2. Implement installation functions
   - `InstallClaudeResources()`
   - `CheckVersion()`
   - `BackupProgressFiles()`
   - `updateGitignore()`

3. Create `cmd/lvt/commands/install.go`
   - `lvt install claude-code` command
   - `lvt upgrade claude-code` command

4. Integrate with `cmd/lvt/commands/new.go`
   - Add InstallClaudeResources() call
   - Add success message

5. Write unit tests
   - `internal/claude/install_test.go`
   - Test embedding, extraction, version management

**Deliverables:**
- `internal/claude/*.go` (complete)
- `cmd/lvt/commands/install.go` (new)
- Updated `cmd/lvt/commands/new.go`
- Unit tests with 100% coverage

**Success Criteria:**
- `lvt new myapp` installs agent + skills
- `lvt install claude-code` works for existing projects
- `lvt upgrade claude-code` preserves progress files
- All embedded files extract correctly
- .gitignore updated automatically

### Phase 3: E2E Testing (Week 2)

**Tasks:**
1. Create `e2e/agent_integration_test.go`
   - Simple request routing test
   - Complex workflow chain test
   - Progress tracking test
   - Resumption test
   - Error recovery test

2. Create `e2e/agent_helpers.go`
   - Test utilities for agent
   - Progress file validation
   - TodoWrite state verification

3. Run all 9 existing skill scenarios through agent
   - Verify same outcomes
   - Document any differences

4. Create golden files for agent outputs
   - Intent classification examples
   - Progress file templates
   - Workflow execution logs

**Deliverables:**
- `e2e/agent_integration_test.go` (complete)
- `e2e/agent_helpers.go` (new)
- Updated test documentation

**Success Criteria:**
- All 5 agent E2E tests pass
- All 9 skill scenarios pass through agent
- 95%+ intent classification accuracy
- Progress tracking works correctly
- Resumption works after interruption

### Phase 4: Documentation & Polish (Week 2)

**Tasks:**
1. Update README with agent information
2. Create tutorial videos/GIFs showing agent in action
3. Update existing skill documentation to mention agent
4. Add troubleshooting guide for common agent issues
5. Create release notes

**Deliverables:**
- Updated `README.md`
- `docs/LVT_AGENT_QUICKSTART.md`
- Tutorial content
- Release notes for agent feature

**Success Criteria:**
- Documentation is clear and comprehensive
- Examples work as shown
- New users can get started in < 5 minutes

## Security Considerations

### Embedded Resource Security

- **Integrity:** Skills and agent are embedded at compile time, tamper-proof
- **Version Tracking:** VERSION file allows detection of manual modifications
- **No Network:** Installation is purely local, no external dependencies

### Progress File Security

- **Location:** `.claude/progress/` is gitignored, not committed
- **Content:** May contain sensitive data (DB URLs, API keys in logs)
- **Permissions:** Created with 0644 (user read/write, group/other read)
- **Cleanup:** Old progress files auto-archived to prevent accumulation

### Agent Execution Security

- **No Arbitrary Code:** Agent only invokes predefined lvt CLI commands
- **Input Validation:** User requests validated before execution
- **Privilege Escalation:** Agent runs with same permissions as user
- **Audit Trail:** All commands logged in progress files

### Best Practices

1. Never commit `.claude/progress/` files
2. Review progress files before sharing (may contain secrets)
3. Use `lvt:manage-env` for environment variables (not hardcoded)
4. Regularly run `lvt upgrade claude-code` for security updates

## Performance Considerations

### Binary Size Impact

**Current lvt binary:** ~8-10 MB (estimate)

**Agent + Skills:**
- 19 skill files × ~5 KB average = ~95 KB
- 1 agent file × ~50 KB = ~50 KB
- **Total addition:** ~150 KB (1.5% increase)

**Verdict:** Negligible impact on binary size

### Installation Performance

**Extraction time:**
- 20 files × ~5 KB = ~100 KB to write
- **Estimated time:** < 100ms on modern SSD

**Verdict:** Imperceptible to users

### Runtime Performance

- Agent is invoked by Claude Code (not lvt CLI)
- No performance impact on lvt commands themselves
- Progress file I/O is minimal (append-only logs)

**Verdict:** No measurable runtime overhead

### Optimization Opportunities

1. **Lazy Progress Tracking:** Only create markdown files for complex workflows
2. **Progress File Compression:** Archive old files with gzip
3. **Agent Caching:** Cache intent classification results (if applicable)

## Future Enhancements

### Phase 2 Features (Post-Launch)

1. **Skill Analytics**
   - Track which skills are used most
   - Identify workflow patterns
   - Optimize agent routing based on usage

2. **Custom Workflows**
   - Allow users to define custom workflow chains
   - Save to `.claude/workflows/my-workflow.md`
   - Agent discovers and suggests custom workflows

3. **Multi-Project Agent**
   - Agent that works across multiple LVT projects
   - Shares knowledge between projects
   - Suggests patterns from other projects

4. **Agent Learning**
   - Remember user preferences
   - Adapt routing based on past interactions
   - Personalized skill recommendations

5. **Collaborative Features**
   - Share progress files with team
   - Team workflow templates
   - Collaborative debugging via agent

### Integration Opportunities

1. **CI/CD Integration**
   - Agent-driven deployment pipelines
   - Automated testing workflows
   - Production monitoring suggestions

2. **IDE Plugins**
   - VS Code extension with agent integration
   - JetBrains plugin
   - Inline skill suggestions

3. **Web Dashboard**
   - Visualize workflow history
   - Progress tracking UI
   - Team analytics

## Appendix

### Skill Catalog Reference

**Core Skills (14):**
1. `lvt:new-app` - Create new LiveTemplate application
2. `lvt:add-resource` - Add database-backed CRUD resource
3. `lvt:add-view` - Add UI-only page without database
4. `lvt:add-migration` - Create and run database migrations
5. `lvt:gen-auth` - Generate complete authentication system
6. `lvt:gen-schema` - Generate database schema without UI
7. `lvt:resource-inspect` - Inspect database resources and schema
8. `lvt:manage-kits` - Manage CSS framework kits
9. `lvt:validate-templates` - Validate and analyze templates
10. `lvt:run-and-test` - Run development server and tests
11. `lvt:customize` - Customize generated code
12. `lvt:seed-data` - Generate realistic test data
13. `lvt:deploy` - Deploy to production (Docker, Fly.io, K8s, VPS)
14. `lvt:manage-env` - Manage environment variables

**Workflow Skills (3):**
1. `lvt:quickstart` - Rapid end-to-end app creation
2. `lvt:production-ready` - Transform dev app to production
3. `lvt:add-related-resources` - Intelligent resource suggestions

**Maintenance Skills (3):**
1. `lvt:analyze` - Comprehensive app analysis
2. `lvt:suggest` - Actionable improvement recommendations
3. `lvt:troubleshoot` - Debug common issues

**Meta Skills (1):**
1. `lvt:add-skill` - Create new skills using TDD

### Example Intent Patterns

| User Input | Intent | Skill | Execution Mode |
|-----------|--------|-------|----------------|
| "create a blog app" | project-init | new-app, quickstart | Guided |
| "add products with name and price" | resource-add | add-resource | Direct |
| "set up authentication" | auth-setup | gen-auth | Guided |
| "add login page" | auth-setup | gen-auth | Guided |
| "deploy to fly.io" | deployment | production-ready, deploy | Complex |
| "my app is slow" | troubleshooting | analyze, troubleshoot | Guided |
| "run the app" | development | run-and-test | Direct |
| "add migration for users table" | database | add-migration | Direct |
| "change to tailwind" | customization | customize, manage-kits | Direct |
| "what resources do I have?" | inspection | resource-inspect | Direct |
| "make production ready" | production-prep | production-ready | Guided |
| "seed test data" | development | seed-data | Direct |

### Progress File Examples

See Section "Progress Tracking System" → "Progress File Format" for detailed examples.

### Glossary

- **Agent:** AI-powered assistant that orchestrates skills
- **Skill:** Markdown document describing how to perform a specific task
- **Workflow:** Sequence of multiple skills chained together
- **Progress File:** Markdown file tracking workflow execution state
- **TodoWrite:** Claude Code tool for real-time progress UI
- **Intent Classification:** Process of mapping user request to skill(s)
- **Direct Execution:** Agent runs CLI commands directly
- **Guided Execution:** Agent follows skill checklist with validation
- **Resumption:** Restarting interrupted workflow from checkpoint

---

**Document Version:** 1.0.0
**Last Updated:** 2025-01-21
**Next Review:** After implementation completion
