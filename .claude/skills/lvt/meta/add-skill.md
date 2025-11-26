---
name: lvt:add-skill
description: Use when creating new Claude Code skills for lvt CLI commands - covers TDD methodology for documentation, research process, skill structure, testing, and refinement
---

# lvt:add-skill

Create high-quality Claude Code skills for lvt CLI commands using TDD methodology.

## Overview

This skill documents the process for creating new lvt command skills. It follows **TDD for documentation**: RED (understand baseline behavior) → GREEN (write skill) → REFACTOR (close loopholes).

**When to use this skill:**
- Adding skills for new lvt commands
- Documenting lvt workflows
- Creating comprehensive reference documentation

**Result:** Production-ready skill that prevents common mistakes and accelerates user workflows.

## TDD Methodology for Skills

Apply RED-GREEN-REFACTOR to documentation:

| Phase | Goal | Activities |
|-------|------|------------|
| **RED** | Understand baseline (no skill) | Research code, identify mistakes users would make |
| **GREEN** | Create skill that prevents mistakes | Write comprehensive skill with examples |
| **REFACTOR** | Close loopholes | Test scenarios, add missing details, fix gaps |

## Phase 1: RESEARCH (RED)

### 1. Identify the Command/Feature

```bash
# Example: Creating skill for "lvt seed" command
# Find the implementation
find . -name "*seed*"
# Find: commands/seed.go, internal/seeder/
```

### 2. Read Core Implementation

**Read in this order:**
1. **Command file** - CLI structure, flags, validation
2. **Implementation files** - Business logic, core functionality
3. **Related files** - Dependencies, helpers, generators
4. **Tests** - Usage patterns, edge cases

**Example research for lvt seed:**
```
Read files:
1. commands/seed.go           ← Command structure
2. internal/seeder/seeder.go  ← Core seeding logic
3. internal/seeder/generator.go ← Data generation
```

### 3. Document Baseline Mistakes

**Ask yourself:** "Without a skill, what would users struggle with?"

**Common baseline mistakes:**
- Missing prerequisites (migrations not run)
- Wrong command syntax
- Incorrect flag usage
- Not understanding error messages
- Skipping important steps
- Misunderstanding command behavior

**Example for lvt seed:**
- ❌ Forgetting to run migrations before seeding
- ❌ Using wrong resource name (typo)
- ❌ Not knowing about --cleanup flag
- ❌ Confusion about test record markers

## Phase 2: WRITE SKILL (GREEN)

### Skill File Structure

**File naming:**
- Location: `~/.claude/skills/lvt/[category]/[command-name].md`
- Category: `core/`, `workflows/`, `advanced/`, `meta/`
- Name: Kebab-case matching command (e.g., `add-resource.md`, `seed-data.md`)

**Frontmatter:**
```markdown
---
name: lvt-command-name       # Must match: lvt-[command]
description: One-line description of when to use this skill - covers [topics]
---

# lvt:command-name

One-sentence summary.

## Overview
- What the command does
- Key features (3-5 bullets)
- When to use it

## Prerequisites
What must be true before using this command

## Basic Usage
Simplest examples that work

## Commands/Sections
Detailed command reference

## Common Issues
❌ Error scenarios with fixes

## Examples
Real-world usage patterns

## Quick Reference
Table: "I want to..." → Command

## Remember
✓ Do's
✗ Don'ts
```

### Writing Guidelines

**1. Start with simplest example:**
```markdown
## Basic Usage

```bash
# Generate 50 products
lvt seed products --count 50
```
```

**2. Add prerequisites explicitly:**
```markdown
**Before seeding:**
1. ✓ Resource generated (`lvt gen resource`)
2. ✓ Migrations applied (`lvt migration up`)
3. ✓ Database exists (`app.db`)
```

**3. Show command output:**
```markdown
**Progress tracking:**
```
Seeding products with 50 rows...
  Progress: 10/50
  Progress: 20/50
✅ Successfully seeded 50 rows
```
```

**4. Document errors with context:**
```markdown
### ❌ Resource Not Found

```bash
# Error: resource 'product' not found in schema

# Cause: Typo or not generated yet

# Fix: Check spelling
lvt seed products --count 50  # not 'product'
```
```

**5. Add "why wrong" explanations:**
```markdown
**Why wrong:** SQLite requires CGO. Without it, database operations fail.
```

### Skill Template

Save this template for new skills:

```markdown
---
name: lvt-COMMAND
description: Use when [doing X] - covers [key topics]
---

# lvt:COMMAND

[One sentence description]

## Overview

[2-3 sentences explaining what this command does]

**Key features:**
- [Feature 1]
- [Feature 2]
- [Feature 3]

## Prerequisites

**Before using this command:**
1. ✓ [Prerequisite 1]
2. ✓ [Prerequisite 2]

## Basic Usage

```bash
# Simplest example
lvt COMMAND [args]
```

## Commands Reference

| Command | Purpose |
|---------|---------|
| `lvt COMMAND` | [Description] |

## Common Scenarios

### [Scenario 1]

```bash
# [Description]
lvt COMMAND [example]
```

## Common Issues

### ❌ [Error Name]

```bash
# Error: [error message]

# Cause: [Why this happens]

# Fix:
lvt COMMAND [correct usage]
```

**Why wrong:** [Explanation]

## Examples

### [Use Case 1]

```bash
# [Description]
lvt COMMAND [example]
```

**Generated/Output:**
- [What happens]

## Quick Reference

**I want to...** | **Command**
---|---
[Task 1] | `lvt COMMAND [args]`
[Task 2] | `lvt COMMAND [args]`

## Remember

✓ [Do this]
✓ [Do that]

✗ Don't [avoid this]
✗ Don't [avoid that]
```

## Phase 3: TEST SKILL (GREEN)

### Testing Approaches

**Option A: Run actual commands (preferred)**
```bash
# Test the documented workflow
cd /tmp
mkdir test_skill && cd test_skill
lvt new testapp
cd testapp

# Follow skill instructions exactly
lvt gen resource products name price:float
lvt migration up
lvt seed products --count 50

# Verify output matches skill documentation

# IMPORTANT: Clean up test files when done
cd /tmp
rm -rf test_skill
```

**Option B: Code analysis (when bash unavailable)**
```bash
# Read implementation to verify skill accuracy
cat commands/seed.go
cat internal/seeder/seeder.go

# Check skill against code:
# - Are flags documented correctly?
# - Are error messages accurate?
# - Are examples valid?
```

### Test Checklist

- [ ] Basic example works as documented
- [ ] Prerequisites are accurate
- [ ] Error scenarios produce documented errors
- [ ] Command flags are correct
- [ ] Examples are realistic and valid
- [ ] Output samples match actual output
- [ ] Quick reference commands work
- [ ] Test files cleaned up after testing

## Phase 4: REFACTOR

### Finding Loopholes

**Common loopholes:**
1. **Missing context** - Assumes knowledge user doesn't have
2. **Incomplete error coverage** - Missing common errors
3. **Vague prerequisites** - Not specific enough
4. **Missing "why" explanations** - Commands without context
5. **Unclear defaults** - What happens when flags omitted
6. **Edge cases** - Unusual but valid scenarios

### Loophole Checklist

**Ask these questions:**
- [ ] Are all flags/options documented?
- [ ] Does it explain WHY, not just WHAT?
- [ ] Are error messages copy-pasted from actual output?
- [ ] Does it cover edge cases?
- [ ] Are examples copy-paste ready?
- [ ] Does quick reference match detailed sections?
- [ ] Are prerequisites testable?
- [ ] Does it link to related skills?

### Example Refactoring

**Before (loophole):**
```markdown
Run migrations before seeding.
```

**After (closed):**
```markdown
**Before seeding:**
1. ✓ Resource generated (`lvt gen resource`)
2. ✓ Migrations applied (`lvt migration up`)
3. ✓ Database exists (`app.db`)

```bash
# Complete setup before seeding
lvt gen resource products name price:float
lvt migration up
lvt seed products --count 50
```
```

## Real Example: Creating lvt:seed-data

**Phase 1: Research**
```bash
# Read implementation
Read commands/seed.go
Read internal/seeder/seeder.go
Read internal/seeder/generator.go

# Identified:
- Context-aware generation based on field names
- --count and --cleanup flags
- Test record markers (test-seed- prefix)
- Progress tracking
```

**Phase 2: Baseline mistakes**
```
Without skill, users would:
❌ Forget to run migrations first
❌ Not know about --cleanup flag
❌ Confusion about resource name vs table name
❌ Not understand context-aware generation
❌ Miss foreign key dependencies
```

**Phase 3: Write skill**
```markdown
Created sections:
1. Basic Usage - Simple examples
2. Prerequisites - Explicit requirements
3. Commands - --count, --cleanup, both
4. Context-Aware Generation - Field name table
5. Common Issues - With fixes
6. Quick Reference - Task → Command
```

**Phase 4: Test**
```bash
# Would test (bash was unavailable):
lvt new testapp
cd testapp
lvt gen resource products name price:float
lvt migration up
lvt seed products --count 50
```

**Phase 5: Refactor (loopholes found)**
```
1. ❌ Boolean field generation not documented
   ✅ Added enabled/active/is_* pattern

2. ❌ Fallback behavior for unknown fields missing
   ✅ Added fallback explanation

3. ❌ Schema.yaml dependency not explained
   ✅ Clarified seeder reads from schema.yaml

4. ❌ No guidance on choosing record counts
   ✅ Added "Choosing Record Counts" section
```

## Cross-Referencing Skills

**Link to related skills:**
```markdown
## Related Skills

- **lvt:new-app** - Create app before adding resources
- **lvt:add-migration** - Custom migrations after resource generation
- **lvt:run-and-test** - Testing your resources

## See Also

For production deployment, see **lvt:deploy** skill.
```

**When to cross-reference:**
- Prerequisites: "Complete lvt:new-app first"
- Next steps: "See lvt:deploy for production"
- Alternatives: "For views without database, see lvt:add-view"
- Advanced topics: "For custom queries, see lvt:customize"

## Common Patterns

### Pattern 1: Prerequisites Section

Always show complete setup:
```markdown
**Before [command]:**
1. ✓ [Prerequisite with verification command]
2. ✓ [Next prerequisite]

```bash
# Complete setup before [command]
[step 1]
[step 2]
[the command]
```
```

### Pattern 2: Error Documentation

Follow this format:
```markdown
### ❌ [Error Name]

```bash
# Error: [exact error message]

# Cause: [Why this happens]

# Fix:
[correct command]
```

**Why wrong:** [Technical explanation]
```

### Pattern 3: Quick Reference

Always include task-oriented table:
```markdown
**I want to...** | **Command**
---|---
[User goal] | `[exact command]`
```

### Pattern 4: Remember Section

Structure as do's and don'ts:
```markdown
## Remember

✓ [Positive action]
✓ [Positive action]

✗ Don't [avoid this]
✗ Don't [avoid that]
```

## Skill Organization

**Directory structure:**
```
~/.claude/skills/lvt/
├── core/              ← Core commands
│   ├── new-app.md
│   ├── add-resource.md
│   ├── add-view.md
│   ├── add-migration.md
│   ├── run-and-test.md
│   ├── customize.md
│   ├── seed-data.md
│   └── deploy.md
├── workflows/         ← Multi-command workflows
│   └── full-crud.md
├── advanced/          ← Advanced topics
│   └── custom-templates.md
└── meta/              ← Skills about skills
    └── add-skill.md   ← This file
```

## Example Session Workflow

**Complete skill creation session:**

```bash
# 1. Research phase
Read implementation files
Document baseline mistakes
Create TodoWrite todos for tracking

# 2. Write phase
Create skill file
Fill in sections following template
Add examples and error scenarios

# 3. Test phase
Run commands to verify accuracy
Check examples work
Verify error messages match

# 4. Refactor phase
Find loopholes
Add missing details
Close gaps
Update Remember section

# 5. Complete
Mark todos complete
Verify skill is comprehensive
```

**TodoWrite pattern:**
```
1. Research [command] and implementation
2. Understand how [feature] works
3. Document baseline mistakes (no skill)
4. GREEN phase: Write skill
5. Test skill with real scenarios
6. REFACTOR: Close loopholes
```

## Quality Checklist

**Before marking skill complete:**

- [ ] File named correctly (`~/.claude/skills/lvt/[category]/[name].md`)
- [ ] Frontmatter has name and description
- [ ] Overview explains what/when/why
- [ ] Prerequisites are explicit and testable
- [ ] Basic usage has simplest working example
- [ ] All command flags documented
- [ ] Common errors have fixes with "why wrong"
- [ ] Examples are realistic and copy-paste ready
- [ ] Quick reference table is complete
- [ ] Remember section has do's and don'ts
- [ ] Cross-references to related skills added
- [ ] No assumptions about user knowledge
- [ ] Tested with actual commands (or code analysis)
- [ ] Test files cleaned up
- [ ] Loopholes identified and closed
- [ ] TodoWrite todos marked complete

## Common Mistakes When Creating Skills

### ❌ Assuming User Knowledge

```markdown
# WRONG
Run the seeder.

# CORRECT
**Before seeding:**
1. ✓ Resource generated (`lvt gen resource`)
2. ✓ Migrations applied (`lvt migration up`)

```bash
lvt seed products --count 50
```
```

### ❌ Vague Error Messages

```markdown
# WRONG
You might get an error about the database.

# CORRECT
### ❌ Database Not Found

```bash
# Error: database not found

# Cause: Haven't run migrations yet

# Fix: Run migrations first
lvt migration up
lvt seed products --count 50
```
```

### ❌ Missing "Why" Explanations

```markdown
# WRONG
Use CGO_ENABLED=1 when building.

# CORRECT
```bash
# Build with CGO enabled
CGO_ENABLED=1 go build ./cmd/myapp
```

**Why:** SQLite requires CGO. Without it, database operations fail with "undefined: sqlite3.Open" errors.
```

### ❌ Incomplete Quick Reference

```markdown
# WRONG
**I want to...** | **Command**
---|---
Seed data | `lvt seed`

# CORRECT
**I want to...** | **Command**
---|---
Generate 50 products | `lvt seed products --count 50`
Remove test data | `lvt seed products --cleanup`
Fresh start | `lvt seed products --cleanup --count 50`
```

## Quick Reference

**I want to...** | **How**
---|---
Create new lvt skill | Follow RED-GREEN-REFACTOR process
Research command | Read commands/, internal/, tests
Document baseline | List mistakes without skill
Write skill | Use skill template above
Test skill | Run actual commands or analyze code
Find loopholes | Review quality checklist
Organize skills | Use core/workflows/advanced/meta structure

## Remember

✓ Use TDD methodology (RED-GREEN-REFACTOR)
✓ Start with research (read implementation)
✓ Document baseline mistakes first
✓ Use skill template for consistency
✓ Follow naming convention: lvt-[command].md
✓ Test with actual commands when possible
✓ Clean up test files after testing
✓ Close loopholes before marking complete
✓ Add "why wrong" explanations for errors
✓ Make examples copy-paste ready
✓ Cross-reference related skills
✓ Use TodoWrite to track progress

✗ Don't skip research phase
✗ Don't assume user knowledge
✗ Don't write vague error messages
✗ Don't forget prerequisite verification
✗ Don't skip testing phase
✗ Don't leave test files behind
✗ Don't leave loopholes open
✗ Don't forget cross-references to related skills
✗ Don't mark complete without quality checklist
