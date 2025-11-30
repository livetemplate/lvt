---
name: lvt-new-app
description: Create a new LiveTemplate application with proper setup and validation
keywords: ["lvt", "livetemplate", "lt"]
category: core
version: 1.0.0
---

# lvt-new-app - Create New Application

This skill creates a new LiveTemplate application using the lvt CLI, validates the setup, and guides the user through next steps.

## ⚠️ IMPORTANT: Check for Brainstorming First

**Before using this skill**, check if the user needs brainstorming:

### When to Use `lvt-plan` Instead

If the user's request is:
- **Vague about requirements** (e.g., "create a blog app" without specifying resources)
- **A new project from scratch** (not adding to existing app)
- **Mentions a domain** (blog, shop, CRM, todo) **without specific resources**

Then **STOP** and use the `lvt-plan` skill instead:
```
Skill("lvt-plan")
```

The lvt-plan skill will:
1. Ask progressive questions to understand requirements
2. Gather resource definitions, auth needs, pagination style, etc.
3. **Then call this skill** (`lvt-new-app`) with complete requirements

### When to Use This Skill Directly

Use `lvt-new-app` directly only when:
- ✅ User has **already completed brainstorming** (you just finished brainstorm workflow)
- ✅ User provides **specific, detailed requirements** (app name, resources with fields, kit preference)
- ✅ User explicitly says **"just create the app"** or **"skip the questions"**

**Example:**
- ❌ "create a blog app" → Use `lvt-plan` (vague, needs planning)
- ✅ "create blog app with posts(title, content), comments(text), use multi kit" → Use `lvt-new-app` (detailed)
- ✅ "now create it" (after planning completed) → Use `lvt-new-app`

## 🎯 ACTIVATION RULES

### Context Detection

This skill activates when **LiveTemplate context is established**:

**✅ Context Established By:**

1. **Project context** - `.lvtrc` file exists in current directory (usually won't apply for "new" app)
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Priority**: Project context > Agent context > Keyword context

### Keyword Matching

**Accepted keywords** (case-insensitive, whole words):
- `lvt`
- `livetemplate`
- `lt`

**Accepted patterns:**
- "create * {app|application|project} {with|using|via} {lvt|livetemplate|lt}"
- "{lvt|lt} new *"
- "new {lvt|livetemplate} {app|application|project}"
- "start * {with|using} {livetemplate|lvt}"

### Trigger Patterns

**With Context (any of: .lvtrc OR agent OR keywords):**
✅ "create a new app called blog"
✅ "start a new project for shop"
✅ "make an app named todos"

**Without Context (must include keywords):**
✅ "create a new lvt app called blog"
✅ "start a livetemplate project for shop"
✅ "make a new app with lt named todos"
❌ "create a new app" (no context, no keywords)

### Examples by Scenario

**Scenario 1: New conversation, no .lvtrc, no agent**
- User: "Create a new app called blog"
  → ❌ No context, no keywords → Don't activate

- User: "Create a new livetemplate app called blog"
  → ✅ Keywords found → Activate skill
  → ✅ Context now established

**Scenario 2: Using lvt-assistant agent**
- User (in agent): "Create an app called shop"
  → ✅ Agent context → Activate skill

**Scenario 3: Context persistence**
- User: "Use lvt to make a blog"
  → ✅ Keywords → Activate, establish context

- User: "Actually, make it called myblog"
  → ✅ Context persists → Can adjust

---

## User Prompts

This skill activates when the user says:

**Explicit prompts:**
- "Create a new lvt app called {name}"
- "Generate a new LiveTemplate application named {name}"
- "Start a new lvt project for {name}"
- "I want to create a new app called {name}"
- "Make a new lvt app: {name}"

**Implicit prompts:**
- "Let's build a {type} app with LiveTemplate"
- "I need to start a new project using lvt for {description}"
- "Create a {name} application"
- "Start a new project"

**Kit-specific prompts:**
- "Create a {name} app with Tailwind" → multi kit
- "Create a SPA called {name}" → single kit
- "Create a simple {name} app" → simple kit

---

## Context Awareness

Before executing, check the following:

### 1. Check Current Directory
- [ ] Are we already in an lvt project?
- [ ] Does `.lvtrc` exist in current directory?
- [ ] **If yes:** Warn user and ask if they want to create in a different location

### 2. Validate Prerequisites
- [ ] Is `lvt` command available? (Check with `which lvt`)
- [ ] Is Go installed? (Check with `go version`)
- [ ] **If missing:** Explain what's needed and how to install

### 3. Validate App Name
- [ ] Is the app name provided?
- [ ] **If not:** Ask user for app name
- [ ] Is it a valid Go module name? (lowercase, alphanumeric, hyphens ok)
- [ ] **If invalid:** Suggest valid alternative

### 4. Determine Kit
- [ ] Did user specify kit in prompt?
- [ ] **If not:** Ask user or use default (multi)
- [ ] Valid options: multi (Tailwind), single (Tailwind SPA), simple (Pico CSS)

---

## Checklist

- [ ] **Step 1:** Validate prerequisites (lvt installed, Go available)
- [ ] **Step 2:** Check if already in lvt project (warn if yes)
- [ ] **Step 3:** Gather app name from user input
- [ ] **Step 4:** Validate app name format
- [ ] **Step 5:** Determine kit selection (ask if not specified)
- [ ] **Step 6:** Run `lvt new {app_name} --kit {kit}`
- [ ] **Step 7:** Monitor command output for errors
- [ ] **Step 8:** Verify directory created
- [ ] **Step 9:** Verify `.lvtrc` exists
- [ ] **Step 10:** Run `go mod tidy` in new directory
- [ ] **Step 11:** Verify app builds successfully
- [ ] **Step 12:** Report success with structure summary
- [ ] **Step 13:** Suggest next steps

---

## Commands Executed

### Check Prerequisites
```bash
# Check if lvt is installed
which lvt

# Check Go version
go version

# Check if already in lvt project
[ -f .lvtrc ] && echo "Already in lvt project" || echo "Not in lvt project"
```

### Create Application
```bash
# Basic app creation (multi kit is default)
lvt new {app_name}

# Or with explicit kit
lvt new {app_name} --kit multi
lvt new {app_name} --kit single
lvt new {app_name} --kit simple

# With custom module name
lvt new {app_name} --module github.com/user/{app_name}
```

### Validate Creation
```bash
# Enter directory
cd {app_name}

# Check .lvtrc exists
cat .lvtrc

# Tidy dependencies
go mod tidy

# Test build
go build ./cmd/{app_name}/main.go
```

---

## Validation

### Success Criteria

This skill verifies successful creation by checking:

1. **Directory Created**
   - Directory `{app_name}/` exists
   - Can `cd` into directory

2. **Configuration File**
   - File `.lvtrc` exists
   - Contains valid kit specification
   - Contains module name

3. **Go Module**
   - File `go.mod` exists
   - Has correct module name
   - Dependencies listed

4. **Project Structure**
   - `cmd/{app_name}/main.go` exists
   - `internal/app/home/` exists
   - `internal/database/` exists
   - `web/assets/` exists

5. **Build Success**
   - `go build` completes without errors
   - No missing imports
   - No compilation errors

6. **Dependencies Resolved**
   - `go mod tidy` succeeds
   - All dependencies downloadable

---

## Error Handling

### Common Errors and Fixes

**Error: `lvt: command not found`**
- **Cause:** lvt CLI not installed
- **Fix:** Install lvt: `go install github.com/livetemplate/lvt/cmd/lvt@latest`
- **Action:** Explain installation, offer to continue after user installs

**Error: `directory already exists`**
- **Cause:** App name conflicts with existing directory
- **Fix:** Choose different name or delete existing directory
- **Action:** Ask user which approach to take

**Error: `invalid module name`**
- **Cause:** App name contains invalid characters
- **Fix:** Suggest valid alternative (lowercase, no special chars except hyphen)
- **Action:** Auto-correct and ask user to confirm

**Error: `go: command not found`**
- **Cause:** Go not installed
- **Fix:** Install Go from https://go.dev/dl/
- **Action:** Explain requirement, cannot proceed without Go

**Error: `go mod tidy failed`**
- **Cause:** Network issues, dependency problems
- **Fix:** Check internet connection, retry
- **Action:** Show error, suggest retry, offer to debug

**Error: Build fails with import errors**
- **Cause:** Dependencies not downloaded properly
- **Fix:** Run `go mod download` then retry
- **Action:** Automatically run fix, revalidate

---

## Success Response Template

When skill completes successfully:

```markdown
✅ Successfully created lvt app: {app_name}

## App Structure
- **Kit:** {kit} (Tailwind CSS / Pico CSS)
- **Module:** {module_name}
- **Location:** {full_path}

## Generated Files
- ✅ cmd/{app_name}/main.go - Application entry point
- ✅ internal/app/home/ - Homepage handler
- ✅ internal/database/ - Database setup
- ✅ web/assets/ - Static assets
- ✅ .lvtrc - Configuration
- ✅ go.mod - Go module

## Validation
- ✅ Build successful
- ✅ Dependencies resolved
- ✅ Ready to run

## Next Steps

You can now:

1. **Start dev server:** `cd {app_name} && lvt serve`
2. **Add your first resource:** Use the `lvt:add-resource` skill
3. **Add authentication:** Use the `lvt:add-auth` skill
4. **Run tests:** `go test ./...`

Would you like me to help with any of these?
```

---

## Next Step Suggestions

Based on context, suggest:

### Always Suggest
1. "Start the development server to see your app running"
2. "Add your first resource (like posts, products, users)"

### Conditional Suggestions
- **If kit is multi or single:** "This uses Tailwind CSS - check the templates in `internal/app/home/`"
- **If kit is simple:** "This uses Pico CSS for minimal styling - perfect for prototyping"
- **If custom module:** "Your module name is `{module}` - remember this for imports"

### Workflow Suggestions
- "Want to jump straight to a working CRUD app? Try `lvt:quickstart`"
- "Want to make this production-ready? Try `lvt:production-ready`"

---

## Testing Notes

### Automated Test Scenarios

1. **Test 1: Basic app with multi kit (default)**
   - Prompt: "Create a new lvt app called testblog"
   - Expected: App created with multi kit, Tailwind CSS
   - Validation: All checks pass, builds successfully

2. **Test 2: App with single kit**
   - Prompt: "Create a SPA called myapp using lvt"
   - Expected: App created with single kit (SPA mode)
   - Validation: Component-based structure

3. **Test 3: App with simple kit**
   - Prompt: "Create a simple app called quicktest"
   - Expected: App created with simple kit, Pico CSS
   - Validation: Minimal structure, counter example

4. **Test 4: App with custom module**
   - Prompt: "Create an app called shop with module github.com/myuser/shop"
   - Expected: App created with custom module name
   - Validation: go.mod has correct module

5. **Test 5: Error - invalid name**
   - Prompt: "Create an app called My-App!"
   - Expected: Validation error, suggest valid alternative
   - Validation: Graceful error handling

### Manual Testing Checklist

See `docs/SKILL_TESTING_CHECKLISTS.md` - "lvt:new-app" section

Key things to verify manually:
- [ ] Homepage loads in browser
- [ ] No console errors
- [ ] CSS framework applied correctly
- [ ] WebSocket connects
- [ ] Layout renders properly

---

## Implementation Notes

### Kit Mapping

When user mentions:
- "Tailwind" or "multi-page" → `--kit multi`
- "SPA" or "single-page" → `--kit single`
- "simple" or "minimal" or "Pico" → `--kit simple`
- No mention → Ask or use default (multi)

### App Name Validation

Valid app names:
- ✅ `blog`, `myapp`, `task-manager`, `api`
- ❌ `My-App`, `app!`, `123app`, `APP`

Validation regex: `^[a-z][a-z0-9-]*$`

### Module Name Inference

If not specified:
- Default pattern: `{app_name}/...` (relative module)
- Better: `github.com/{user}/{app_name}` (if user info available)
- Ask user if ambiguous

---

## Example Interactions

### Example 1: Basic Creation

**User:** "Create a new lvt app called blog"

**Skill:**
1. Validates prerequisites ✅
2. Checks not in lvt project ✅
3. Runs: `lvt new blog`
4. Runs: `cd blog && go mod tidy`
5. Validates build ✅
6. Reports success
7. Suggests: "Start dev server" or "Add posts resource"

---

### Example 2: With Kit Selection

**User:** "Create a SPA called taskapp"

**Skill:**
1. Detects "SPA" → kit = single
2. Runs: `lvt new taskapp --kit single`
3. Validates ✅
4. Explains: "Created single-page app with Tailwind CSS"
5. Suggests: "This is component-based - check out the structure"

---

### Example 3: Error Handling

**User:** "Create an app called My-Blog!"

**Skill:**
1. Validates name → INVALID
2. Responds: "App name 'My-Blog!' contains invalid characters. Go module names must be lowercase alphanumeric with hyphens."
3. Suggests: "How about 'my-blog' instead?"
4. Waits for user confirmation
5. Proceeds with corrected name

---

## Success Metrics

Track for this skill:
- **Completion rate:** % of attempts that succeed
- **Error rate:** % that hit validation errors
- **Average duration:** Time from prompt to success
- **User satisfaction:** Rating after completion
- **Next action:** What users do after (add-resource, dev server, etc.)

---

## Gap Tracking

Issues discovered during testing will be documented in `docs/CLAUDE_SKILLS_TRACKER.md` under the "Discovered Gaps" section with format:

```markdown
### GAP-XXX: Issue title
- **Skill:** lvt:new-app
- **Scenario:** Test N
- **Priority:** P0/P1/P2/P3
- **Status:** Open/Fixed
```

---

## Version History

- **v1.0** (2025-11-03): Initial implementation
- Focus on happy path + basic error handling
- Supports all 3 kit types
- Comprehensive validation
