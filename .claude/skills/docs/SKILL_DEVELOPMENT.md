# LVT Claude Code Skills - Development Guide

This guide explains how to develop, test, and refine Claude Code skills for the lvt CLI tool.

---

## Table of Contents

1. [Overview](#overview)
2. [Skill Development Cycle](#skill-development-cycle)
3. [Writing Skills](#writing-skills)
4. [Testing Strategy](#testing-strategy)
5. [Gap Discovery Process](#gap-discovery-process)
6. [Fix & Retest Cycle](#fix--retest-cycle)
7. [Best Practices](#best-practices)

---

## Overview

### Goals

- **Speed:** Rapid iteration from idea → skill → test → fix
- **Quality:** High confidence through automated + manual validation
- **Discovery:** Find gaps in skills, CLI, and templates early
- **Documentation:** Track all issues and resolutions

### Core Philosophy

**Test early, test often, fix immediately**

Every skill goes through multiple test iterations with both automated validation and human verification before being considered complete.

---

## Skill Development Cycle

### The 15-30 Minute Loop

```
1. Write Skill Draft (5-10 min)
   ↓
2. Test with Real App (5 min)
   ↓
3. Automated Validation (1 min)
   ↓
4. Manual Testing (5 min)
   ↓
5. Document Issues (2 min)
   ↓
6. Fix & Repeat (5 min)
```

### Iteration Until Success

**Success Criteria:**
- 10+ consecutive automated test passes
- 5+ manual test sessions with 4+ star rating
- Zero P0/P1 issues remaining
- Clear documentation

---

## Writing Skills

### Skill File Structure

Skills are markdown files with front matter and checklists.

**Location:** `skills/lvt/[category]/[skill-name].md`

**Template:**

```markdown
---
name: lvt:skill-name
description: Brief description of what this skill does
---

# Skill Name

## User Prompts

This skill activates when the user says:
- "Exact prompt 1"
- "Exact prompt 2"
- "Pattern with {variable}"

## Context Awareness

Before executing, this skill:
- Checks if already in lvt project
- Reads existing configuration
- Validates prerequisites

## Checklist

- [ ] Step 1: Check prerequisites
- [ ] Step 2: Gather inputs
- [ ] Step 3: Run lvt command
- [ ] Step 4: Validate success
- [ ] Step 5: Report results
- [ ] Step 6: Suggest next steps

## Commands Executed

```bash
# Example commands this skill runs
lvt new {app_name} --kit {kit}
go mod tidy
```

## Validation

This skill verifies:
- File exists: `.lvtrc`
- Directory exists: `{app}/.`
- Command succeeds: `go build`

## Error Handling

If errors occur:
- Parse error message
- Explain to user
- Suggest fix
- Offer to retry

## Next Steps

After successful execution, suggest:
- Next logical action
- Testing command
- Related skills
```

### Skill Categories

**Core Skills** (`skills/lvt/core/`)
- Atomic operations
- Single lvt command
- Building blocks

**Workflow Skills** (`skills/lvt/workflows/`)
- Multi-step processes
- Chain multiple core skills
- Complete scenarios

**Maintenance Skills** (`skills/lvt/maintenance/`)
- Analysis and inspection
- Troubleshooting
- Suggestions

---

## Testing Strategy

### Phase 1: Automated Validation

**Goal:** Catch obvious failures fast

**Process:**
1. Create test session
2. Run skill
3. Execute validation script
4. Record pass/fail

**Script:** `/tmp/lvt-skill-tests/validate-generated-app.sh`

**Checks:**
- ✅ Build success
- ✅ Tests pass
- ✅ Migrations valid
- ✅ Code quality
- ✅ Configuration present

**Time:** 30-60 seconds per test

---

### Phase 2: Manual Validation

**Goal:** Verify user experience and catch UX issues

**Process:**
1. Start dev server
2. Open browser
3. Follow manual checklist
4. Test all CRUD operations
5. Check console for errors
6. Rate experience (1-5 stars)

**Time:** 3-5 minutes per test

**See:** `SKILL_TESTING_CHECKLISTS.md` for detailed checklists

---

### Phase 3: Issue Collection

**When issues found:**

1. **Stop** - Don't continue testing
2. **Document** - Record in `gaps.md`
3. **Categorize** - Skill/CLI/Template/Docs
4. **Prioritize** - P0/P1/P2/P3
5. **Fix** - Apply resolution
6. **Retest** - Verify fix works

**Use AskUserQuestion for structured feedback:**

```markdown
What specifically didn't work?
- [ ] Build failed
- [ ] Tests failed
- [ ] UI broken
- [ ] Feature not working
- [ ] Other

What did you observe?
[Free text]

What did you expect?
[Free text]
```

---

## Gap Discovery Process

### Issue Categories

**1. Skill Logic Issues**
- Prompt doesn't match
- Missing validation
- Incorrect command
- Unclear guidance

**Fix Location:** Skill markdown file

**2. LVT CLI Issues**
- Command fails
- Poor error messages
- Missing features
- Unexpected behavior

**Fix Location:** lvt CLI codebase

**3. Template Issues**
- Generated code has bugs
- Missing best practices
- Poor defaults
- Test failures

**Fix Location:** Template files in `internal/kits/`

**4. Documentation Issues**
- Unclear instructions
- Missing examples
- Confusing output

**Fix Location:** READMEs, docs, skill descriptions

---

### Priority Levels

**P0 (Blocker) - Fix Immediately**
- Breaks core workflow
- App won't compile
- Tests can't run
- Critical security issue

**Example:** "Generated app has undefined variable"

**P1 (Critical) - Fix This Iteration**
- Significantly degrades UX
- Confusing error messages
- Missing essential validation
- Poor performance

**Example:** "No feedback when command fails"

**P2 (Important) - Fix Next Iteration**
- Workaround exists
- Nice-to-have features
- Documentation gaps
- Minor UX issues

**Example:** "Would be nice to show progress"

**P3 (Enhancement) - Backlog**
- Quality of life improvements
- Edge case handling
- Polish
- Advanced features

**Example:** "Add color to output"

---

### Gap Tracking Format

```markdown
## GAP-001: Issue Title

**Discovery Method:** Manual Testing
**Skill:** lvt:skill-name
**Scenario:** Specific test scenario
**Date:** 2025-11-03
**Priority:** P0

**User Feedback:**
> "What the user reported"

**Automated Checks:**
- ✅ Build: Success
- ❌ Tests: Failed
- ✅ Vet: Clean

**Root Cause:**
Detailed analysis of what's wrong and why

**Fixes Required:**
1. Fix template X
2. Update skill Y
3. Add validation Z

**Status:** 🔴 Open | 🟡 In Progress | 🟢 Fixed
**Resolution:** What was done to fix it
**Retest:** ✅ Confirmed working
```

---

## Fix & Retest Cycle

### When Issue Found

1. **Log the Issue**
   - Add to `CLAUDE_SKILLS_TRACKER.md` gaps section
   - Create GAP-XXX entry
   - Include all context

2. **Analyze Root Cause**
   - Read relevant code
   - Understand why it failed
   - Identify fix location

3. **Propose Fix**
   - What needs to change?
   - Which files affected?
   - Any side effects?

4. **Apply Fix**
   - Edit skill file, or
   - Edit CLI code, or
   - Edit template, or
   - Update documentation

5. **Retest in Clean Environment**
   ```bash
   cd /tmp/lvt-skill-tests
   ./new-test-session.sh lvt:skill-name "retest-GAP-001"
   cd current
   # Run same test scenario
   ../validate-generated-app.sh testapp
   ```

6. **Verify Fix**
   - Automated validation passes?
   - Manual testing confirms fix?
   - No new issues introduced?

7. **Update Tracker**
   - Mark GAP as fixed
   - Record resolution
   - Update metrics

8. **Move to Next Issue**
   - If more P0/P1 issues: Fix next
   - If all critical fixed: Continue development
   - If all issues fixed: Skill complete!

---

## Best Practices

### Do's

✅ **Test with Real Prompts**
- Use conversational language
- Test ambiguous inputs
- Try variations

✅ **Isolate Each Test**
- Fresh directory per test
- Clean Go module cache
- No shared state

✅ **Test Error Paths**
- Invalid inputs
- Missing prerequisites
- Edge cases

✅ **Document Everything**
- What you tested
- What broke
- What you fixed
- Why it broke

✅ **Fix P0/P1 Immediately**
- Don't accumulate critical bugs
- Stop and fix blockers
- Verify fix before continuing

✅ **Use Validation Scripts**
- Automate what you can
- Consistent checking
- Fast feedback

✅ **Ask for Human Feedback**
- Manual testing catches UX issues
- Screenshots help debugging
- User ratings indicate quality

---

### Don'ts

❌ **Don't Skip Testing**
- Every skill needs 10+ automated tests
- Every skill needs 5+ manual tests
- No exceptions

❌ **Don't Test in Main Project**
- Use `/tmp/lvt-skill-tests/`
- Keep project directory clean
- Prevent accidental commits

❌ **Don't Batch Fixes**
- Fix and verify immediately
- One issue at a time
- Confirm before moving on

❌ **Don't Ignore P0/P1**
- Critical issues block progress
- Fix before continuing
- Don't work around them

❌ **Don't Assume It Works**
- Always validate
- Always manual test
- Always get feedback

❌ **Don't Skip Documentation**
- Update tracker after each task
- Document decisions
- Track all gaps

---

## Testing Workflow Example

### Complete Cycle for `lvt:new-app`

**1. Create Skill File**
```bash
cd /Users/adnaan/code/livetemplate/lvt/.worktrees/claude-code-skills
mkdir -p skills/lvt/core
vim skills/lvt/core/new-app.md
# Write skill with checklist
```

**2. Start Test Session**
```bash
cd /tmp/lvt-skill-tests
./new-test-session.sh lvt:new-app "test-1-basic-multi-kit"
cd current
```

**3. Run Skill (via Claude Code)**
```
User: "Create a new lvt app called testblog"
Skill: lvt:new-app executes
```

**4. Automated Validation**
```bash
../validate-generated-app.sh testblog
# Result: ✅ All checks passed
```

**5. Manual Testing**
```bash
cd testblog
PORT=8080 go run cmd/testblog/main.go &
open http://localhost:8080

# Follow checklist:
# ✅ Homepage loads
# ✅ No console errors
# ✅ WebSocket connects
# ✅ Layout renders correctly

pkill -f "testblog"
```

**6. Document Results**
```bash
cd /tmp/lvt-skill-tests/current
vim SESSION.md
# Update with test results
# Rating: 5/5 stars
```

**7. Repeat 9 More Times**
```bash
# Different scenarios:
# - Simple kit
# - Single kit
# - Custom module
# - Invalid name (error handling)
# etc.
```

**8. Update Tracker**
```bash
cd /Users/adnaan/code/livetemplate/lvt/.worktrees/claude-code-skills
vim docs/CLAUDE_SKILLS_TRACKER.md
# Mark skill complete
# Update metrics
# Document any gaps found
```

**9. Skill Complete!**
- 10/10 automated tests passed
- 5/5 manual tests passed
- 0 P0/P1 issues
- Average rating: 4.8/5.0

---

## Troubleshooting

### Common Issues

**Issue: Tests fail in worktree**
```bash
# Solution: Use GOWORK=off
GOWORK=off go test ./...
```

**Issue: Port already in use**
```bash
# Solution: Kill existing process
lsof -ti:8080 | xargs kill -9
```

**Issue: Database locked**
```bash
# Solution: Remove database file
rm *.db
```

**Issue: Stale Go module cache**
```bash
# Solution: Clear cache
go clean -modcache
go mod download
```

---

## Metrics to Track

### Per Skill
- Automated test pass rate (target: 100%)
- Manual test success rate (target: 80%+)
- Average user rating (target: 4.0+/5.0)
- Number of gaps discovered
- Time to fix cycle (target: <1 hour)

### Overall Project
- Total skills complete (target: 13)
- Total gaps discovered
- P0/P1 resolution rate (target: 100%)
- Time to working app (target: <2 min)

---

## Next Steps

After reading this guide:

1. **Review** `SKILL_TESTING_CHECKLISTS.md` for detailed test procedures
2. **Set up** testing infrastructure (already done!)
3. **Create** your first skill
4. **Test** following this guide
5. **Document** results in tracker
6. **Iterate** until success criteria met

---

**Remember:** The goal is not to write perfect skills on the first try. The goal is to iterate rapidly, discover gaps early, and fix them systematically. Every test makes the skills better!
