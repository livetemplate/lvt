---
name: lvt-ui-polish
description: Generate test apps, capture screenshots, analyze UI for issues, and recursively fix problems in kit templates
category: maintenance
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt", "ui", "polish", "screenshot", "test"]
---

# lvt:ui-polish

Automated UI testing and polishing skill for LiveTemplate kit templates. Generates test applications, captures screenshots of all UI states, analyzes them for functional/UX issues using vision capabilities, and recursively fixes problems in the embedded kit templates.

## ACTIVATION RULES

### Context Detection

This skill runs in **LiveTemplate project directories** (where lvt CLI is available).

**Context Established By:**
1. **Project context** - `.lvtrc` exists or working in lvt repo
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt" with UI/polish keywords

**Keyword matching** (case-insensitive): `lvt`, `livetemplate`, `lt` + `ui`, `polish`, `test`, `screenshot`, `fix`

### Trigger Patterns

**With Context:**
- "polish the UI"
- "test UI and fix issues"
- "analyze screenshots"
- "fix UI bugs in templates"

**Without Context (needs keywords):**
- "polish lvt UI for multi kit"
- "test livetemplate UI"
- "fix lvt template UI issues"

---

## User Prompts

**When to use:**
- "Polish the UI for multi kit"
- "Test and fix UI issues"
- "Run UI tests on all kits"
- "Analyze the generated UI"
- "Fix UI bugs in templates"
- "Screenshot test the single kit"

**Examples with arguments:**
- `/ui-polish --kit multi` - Test specific kit
- `/ui-polish --all` - Test all kits
- `/ui-polish --kit multi --analyze-only` - No fixes, just report issues
- `/ui-polish --kit multi --resource "products name price:float"` - Custom resource

---

## Prerequisites

Before running this skill, ensure:
1. **lvt CLI** - Must be installed and accessible
2. **Docker** - Required for Chrome container (screenshot capture)
3. **Go** - For building test apps
4. **ANTHROPIC_API_KEY** - Already configured (for vision analysis)

---

## CONFIRM BEFORE EXECUTING

**STOP and show the user:**

```
UI Polish Plan
==============
Kit(s) to test: {kit or "all"}
Resource: posts (title, content:text, published:bool)
UI States: list, add_modal, edit_modal, delete_confirm, validation

Steps:
1. Generate test app in /tmp/lvt-ui-polish-{timestamp}/
2. Capture screenshots (5-7 per kit)
3. Analyze for functional/UX issues
4. Fix issues in internal/kits/system/{kit}/ templates
5. Verify fixes, repeat if needed (max 3 attempts/issue)
6. Run e2e tests to check for regressions

Ready to proceed? (yes/no/change kit/change resource)
```

**Wait for user confirmation before executing.**

---

## Checklist

### Phase 1: Setup (Steps 1-4)

1. [ ] **Verify prerequisites**
   ```bash
   which lvt && echo "lvt: OK" || echo "lvt: MISSING"
   docker info > /dev/null 2>&1 && echo "Docker: OK" || echo "Docker: MISSING"
   go version && echo "Go: OK" || echo "Go: MISSING"
   ```

2. [ ] **Create temp directory**
   ```bash
   WORK_DIR="/tmp/lvt-ui-polish-$(date +%s)"
   mkdir -p "$WORK_DIR/screenshots"
   echo "Working in: $WORK_DIR"
   ```

3. [ ] **Determine kit to test** (from user input or default to "multi")
   - Options: `multi`, `single`, `simple`
   - Default: `multi`

4. [ ] **Confirm with user** (show plan, wait for approval)

### Phase 2: Generate Test App (Steps 5-7)

5. [ ] **Create test application**
   ```bash
   cd "$WORK_DIR"
   lvt new testapp --kit {kit}
   cd testapp
   ```

6. [ ] **Add test resource** (skip for simple kit)
   ```bash
   # For multi/single kits:
   lvt gen resource posts title content:text published:bool

   # Run migrations
   lvt migration up
   ```

7. [ ] **Build and verify**
   ```bash
   go build ./cmd/testapp
   echo "Build: OK"
   ```

### Phase 3: Capture Screenshots (Steps 8-9)

8. [ ] **Start test app server**
   ```bash
   PORT=9999 ./cmd/testapp/testapp &
   APP_PID=$!
   sleep 2  # Wait for server to start
   curl -s http://localhost:9999 > /dev/null && echo "Server: OK"
   ```

9. [ ] **Capture UI state screenshots**

   Use chromedp or browser e2e helpers to capture:

   | State | Action | Screenshot |
   |-------|--------|------------|
   | list_empty | Navigate to /posts | list_empty.png |
   | add_modal | Click "Add" button | add_modal.png |
   | form_filled | Fill form fields | form_filled.png |
   | list_with_data | After adding item | list_with_data.png |
   | edit_modal | Click "Edit" on item | edit_modal.png |
   | delete_confirm | Click "Delete" | delete_confirm.png |
   | validation_errors | Submit empty form | validation.png |

   Save to: `$WORK_DIR/screenshots/`

### Phase 4: Analyze Screenshots (Step 10)

10. [ ] **Analyze each screenshot for issues**

    For each screenshot file, use the Read tool to view it, then analyze:

    **Analysis Prompt:**
    ```
    Analyze this screenshot of a LiveTemplate-generated {kit} kit CRUD application.

    Context:
    - Kit: {kit} (CSS: {tailwind|none})
    - Resource: posts
    - UI State: {state_name}
    - Fields: title (string), content (text), published (bool)

    Look for:

    1. FUNCTIONAL ISSUES (Critical)
       - Template errors visible ({{.variable}} showing raw)
       - Broken layouts or overlapping elements
       - Buttons/forms that appear non-functional
       - Missing content that should be present

    2. UX ISSUES (High/Medium)
       - Poor alignment or spacing
       - Inconsistent styling
       - Confusing element placement
       - Missing visual feedback

    3. ACCESSIBILITY ISSUES (Medium)
       - Low contrast text
       - Missing focus indicators
       - Poor visual hierarchy

    Output as JSON:
    [
      {
        "type": "functional|ux|accessibility",
        "severity": "critical|high|medium|low",
        "description": "Clear description",
        "location": "Element location",
        "suggested_fix": "Fix suggestion"
      }
    ]

    Return [] if no issues found.
    ```

### Phase 5: Fix Loop (Steps 11-13)

11. [ ] **For each issue (priority: critical > high > medium > low):**

    a. [ ] **Spawn fix subagent** via Task tool:
       ```
       Task(
         subagent_type="general-purpose",
         description="Fix UI issue in {kit} kit",
         prompt="""
         You are fixing a UI issue in LiveTemplate kit templates.

         ## Issue
         - Type: {type}
         - Severity: {severity}
         - Description: {description}
         - Location: {location}
         - Suggested Fix: {suggested_fix}

         ## Screenshot
         Read the screenshot at: {screenshot_path}

         ## Template Files to Check
         1. internal/kits/system/{kit}/templates/resource/template.tmpl.tmpl
         2. internal/kits/system/{kit}/templates/resource/template_components.tmpl.tmpl
         3. internal/kits/system/{kit}/components/form.tmpl
         4. internal/kits/system/{kit}/components/table.tmpl

         ## Instructions
         1. Read the screenshot to understand the visual issue
         2. Read the relevant template files
         3. Identify the root cause
         4. Make the MINIMAL fix required
         5. Ensure fix is GENERIC (works for all resources)
         6. Do NOT add resource-specific code

         ## Constraints
         - Only modify files in internal/kits/system/{kit}/
         - Preserve existing functionality
         - Use the kit's CSS framework
         """
       )
       ```

    b. [ ] **Rebuild test app**
       ```bash
       cd "$WORK_DIR/testapp"
       # Re-generate with updated templates
       rm -rf app/posts database/migrations
       lvt gen resource posts title content:text published:bool
       lvt migration up
       go build ./cmd/testapp
       ```

    c. [ ] **Restart server and re-capture screenshot**

    d. [ ] **Verify fix** - Re-analyze the new screenshot:
       ```
       I previously identified this issue:
       {original_issue_json}

       Here is the screenshot AFTER a fix was applied.

       Is the original issue now FIXED? (yes/no)
       Any NEW issues introduced?
       ```

    e. [ ] **If not fixed** - Retry (max 3 attempts per issue)

    f. [ ] **If fixed** - Continue to next issue

12. [ ] **Track modified files**
    ```bash
    git -C /path/to/lvt/repo status --short internal/kits/system/
    ```

13. [ ] **Run regression tests**
    ```bash
    cd /path/to/lvt/repo
    go test -tags browser ./e2e -run Test{Kit} -v
    ```

### Phase 6: Cleanup & Report (Steps 14-15)

14. [ ] **Cleanup**
    ```bash
    # Stop test app server
    kill $APP_PID 2>/dev/null || true

    # Optionally remove temp directory
    # rm -rf "$WORK_DIR"
    ```

15. [ ] **Report results**
    ```
    UI Polish Report
    ================
    Kit tested: {kit}
    Screenshots captured: {count}

    Issues Found: {total}
    - Critical: {critical_count}
    - High: {high_count}
    - Medium: {medium_count}
    - Low: {low_count}

    Issues Fixed: {fixed_count}
    Issues Remaining: {remaining_count}

    Modified Files:
    - internal/kits/system/{kit}/templates/resource/template.tmpl.tmpl
    - internal/kits/system/{kit}/components/form.tmpl

    Next Steps:
    1. Review changes: git diff internal/kits/system/
    2. Run full tests: go test -tags browser ./e2e
    3. Commit if satisfied: git add -A && git commit -m "fix: UI improvements for {kit} kit"
    ```

---

## Kit-Specific Notes

### multi kit
- CSS: Tailwind
- Layout: Multi-page with navigation
- Components: table, form, modal, pagination
- Focus: CRUD operations, modals, table layout

### single kit
- CSS: Tailwind
- Layout: Single-page application
- Components: Same as multi but SPA-style
- Focus: In-page updates, no full page reloads

### simple kit
- CSS: None (plain HTML)
- Layout: Single page with counter/clock
- Components: Minimal
- Focus: Basic rendering, WebSocket updates

---

## Error Handling

### Screenshot capture fails
```
Error: chromedp timeout or connection failed

Solutions:
1. Ensure Docker is running: docker info
2. Start Chrome container: docker run -d --name chrome zenika/alpine-chrome
3. Check port availability: lsof -i :9222
4. Retry with longer timeout
```

### Template syntax error after fix
```
Error: template parse error

Solutions:
1. Run: lvt parse internal/kits/system/{kit}/templates/resource/template.tmpl.tmpl
2. Check for unclosed {{ }} tags
3. Verify Go template syntax
4. Rollback: git checkout internal/kits/system/{kit}/
```

### Build fails after template change
```
Error: go build failed

Solutions:
1. Check error message for specific issue
2. Verify template doesn't break Go code generation
3. Rollback and try smaller fix
4. Run: go vet ./...
```

### E2E tests fail (regression)
```
Error: e2e test failed after fix

Solutions:
1. Check which test failed
2. Compare before/after screenshots
3. Rollback problematic fix
4. Fix must pass all existing tests
```

---

## Success Criteria

UI polish is complete when:
1. All critical/high issues are fixed
2. E2e tests pass (no regressions)
3. Modified templates are valid (lvt parse passes)
4. Changes committed to feature branch
5. User approves the visual improvements

---

## Example Session

```
User: polish the UI for multi kit

Claude: Let me run the UI polish workflow for the multi kit.

[Shows confirmation plan]

User: yes

Claude: Starting UI polish...

1. Creating test app in /tmp/lvt-ui-polish-1704567890/
   - lvt new testapp --kit multi ✓
   - lvt gen resource posts ✓
   - go build ✓

2. Capturing screenshots...
   - list_empty.png ✓
   - add_modal.png ✓
   - edit_modal.png ✓
   - delete_confirm.png ✓
   - validation.png ✓

3. Analyzing screenshots...
   Found 3 issues:
   - [HIGH] Modal close button misaligned
   - [MEDIUM] Table header lacks contrast
   - [LOW] Form labels inconsistent spacing

4. Fixing issues...
   - Issue 1: Spawning subagent to fix modal...
     [Subagent edits template_components.tmpl.tmpl]
     Verified: Fixed ✓

   - Issue 2: Spawning subagent to fix table...
     [Subagent edits table.tmpl]
     Verified: Fixed ✓

   - Issue 3: Spawning subagent to fix form...
     [Subagent edits form.tmpl]
     Verified: Fixed ✓

5. Running regression tests...
   go test -tags browser ./e2e -run TestMulti ✓

UI Polish Complete!
==================
Issues Found: 3
Issues Fixed: 3

Modified Files:
- internal/kits/system/multi/templates/resource/template_components.tmpl.tmpl
- internal/kits/system/multi/components/table.tmpl
- internal/kits/system/multi/components/form.tmpl

Review: git diff internal/kits/system/multi/
Commit: git add -A && git commit -m "fix: UI polish for multi kit"
```

---

## Notes

- Always create test app in /tmp to avoid polluting project
- Screenshots are essential for vision-based analysis
- Fixes must be GENERIC (work for any resource, not just test resource)
- Run e2e tests after each fix to catch regressions early
- Keep fixes minimal - don't over-engineer solutions
- Document what was changed in commit messages
- Multiple runs may be needed for complex UI issues
