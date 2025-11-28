---
name: lvt-brainstorm
description: Interactive planning and app creation - helps design and create LiveTemplate apps through progressive questions. Use for "create/build/make a [lvt/livetemplate] [domain] app" OR "plan/design/brainstorm a [lvt/livetemplate] app"
keywords: ["lvt", "livetemplate", "lt", "create", "build", "make", "plan", "design", "brainstorm"]
requires_keywords: true
category: workflows
version: 1.0.0
---

# lvt-brainstorm

Interactive planning skill that guides users through progressive questions to design their LiveTemplate application before generating any code.

## 🎯 ACTIVATION RULES

### Special Rule: Always Requires Keywords

**Unlike other skills**, brainstorming ALWAYS requires keywords, even if:
- `.lvtrc` exists (already in project)
- `lvt-assistant` agent is active
- Previous keywords established context

**Why**: "Help me plan an app" is too generic - user might want to plan any type of app (Next.js, Rails, etc.). Requiring keywords ensures they want LiveTemplate planning.

### Keyword Matching

**Accepted keywords** (case-insensitive, whole words):
- `lvt`
- `livetemplate`
- `lt`

**Accepted patterns:**

**Planning patterns:**
✅ "help me plan a **livetemplate** app"
✅ "walk me through **lvt** project design"
✅ "**lt** brainstorm for blog"
✅ "design a shop **with livetemplate**"
✅ "use **lvt** to help me plan"

**Creation patterns (WITHOUT specific resources):**
✅ "create a **lvt** blog app" (domain mentioned, no resources)
✅ "build a **livetemplate** shop" (domain mentioned, no resources)
✅ "make an **lvt** CRM" (domain mentioned, no resources)
✅ "start a **livetemplate** todo app" (domain mentioned, no resources)

❌ "help me plan an app" (too generic, even with context)
❌ "let's brainstorm" (no LiveTemplate keywords)
❌ "create an app" (no keywords, too generic)

### Examples

**Will Activate:**
- "Help me plan a LiveTemplate blog"
- "Walk me through creating an lvt e-commerce app"
- "I want to brainstorm a todo app using livetemplate"
- "Can you help me design a CRM with lvt?"
- **"Create a lvt blog app"** ✅ NEW: Domain without resources
- **"Build a livetemplate shop"** ✅ NEW: Domain without resources
- **"Make an lvt todo application"** ✅ NEW: Domain without resources

**Won't Activate:**
- "Help me plan a blog" (no keywords)
- "Let's brainstorm an app" (no keywords)
- "I need help designing my application" (no keywords)
- "Create blog app with posts(title,content) and multi kit" (detailed requirements, skip to new-app)

---

## Purpose

Guide users through progressive questions to understand requirements, then execute appropriate LiveTemplate commands.

**Not a CLI command** - purely conversational in Claude Code.

**Progressive disclosure**: Start with 3-5 core questions, offer "more options" for detailed configuration.

---

## Progressive Question Flow

### Phase 1: Core Questions (Always Ask)

Ask these 4-6 questions to establish basics:

#### Question 1: App Name
"What would you like to name your app?

This will be:
- The directory name (e.g., `myblog`)
- The Go module name (e.g., `github.com/yourname/myblog`)
- Used in your code and configuration

**Requirements:**
- Lowercase letters, numbers, hyphens
- No spaces or special characters
- Example: `myblog`, `todo-app`, `my-shop`

Your answer:"

#### Question 2: App Domain
"What type of application are you building?

Common types:
- **Blog** (posts, comments, categories)
- **E-commerce** (products, orders, cart)
- **SaaS** (users, organizations, subscriptions)
- **Todo/Tasks** (tasks, projects, labels)
- **CRM** (contacts, deals, activities)
- **Forum** (topics, posts, replies)
- **Other** (describe it)

Your answer:"

#### Question 3: Primary Resource
"What's the main thing you're tracking in your app?

For a blog → `posts`
For a shop → `products`
For todos → `tasks`
For CRM → `contacts`

Your answer:"

#### Question 3: Authentication
"Will users need to log in to your app?

Options:
- **Yes, with password** (email + password login)
- **Yes, passwordless** (magic link via email)
- **Yes, both** (password + magic link options)
- **No** (public app, no user accounts)

Your answer:"

#### Question 4: Related Resources
"Besides {primary_resource}, are you tracking anything else?

For a blog → maybe `comments`, `categories`, `tags`?
For a shop → maybe `orders`, `reviews`, `cart_items`?
For todos → maybe `projects`, `labels`, `users`?

List them separated by commas, or say 'none':
Your answer:"

#### Question 5: Test Data
"Want me to generate sample data for testing?

Options:
- **10 records** (quick testing)
- **50 records** (realistic amount)
- **100 records** (stress testing)
- **None** (I'll add my own data)

Your answer:"

---

### Phase 2: Summary & Offer More Options

After collecting core answers, show summary and offer to dive deeper:

"📋 **Summary So Far**

- **App name**: {app_name}
- **App type**: {domain}
- **Primary resource**: {primary_resource}
- **Authentication**: {yes/no + method}
- **Related data**: {related_resources or "none"}
- **Test data**: {count} records

I have enough information to create a working app!

**Want to configure more details?**
- Field types (I'll infer smart defaults, but you can customize)
- Pagination style (infinite scroll, load more, page numbers)
- Edit mode (modal popup vs dedicated page)
- CSS framework (Tailwind, Bulma, Pico, none)

Configure more options? (yes/no):"

**If NO** → Skip to Phase 4 (Preview)
**If YES** → Continue to Phase 3

---

### Phase 3: Detailed Configuration (Optional)

Only ask these if user wants detailed configuration:

#### Question 6: Field Types
"What fields should **{primary_resource}** have?

**Option A: Let me infer** (recommended for speed):
Just list field names: `title content published`

**Option B: Specify types explicitly**:
`title:string content:text published:bool published_at:time`

**Available types:**
- `string` - Short text (VARCHAR)
- `text` - Long text, textarea (TEXT)
- `int` - Integer numbers
- `float` - Decimal numbers
- `bool` - True/false checkbox
- `time` - Timestamp (created_at, updated_at)
- `references:{table}` - Foreign key to another table

Your answer:"

#### Question 7: Pagination Style
"How should the **{primary_resource}** list display items?

Options:
1. **Infinite scroll** - Load more as you scroll down (like Twitter)
2. **Load more button** - Manual click to load more (like Instagram)
3. **Prev/Next buttons** - Simple forward/back navigation
4. **Page numbers** - Numbered pages 1, 2, 3... (like Google)

Your answer (1-4 or name):"

#### Question 8: Edit Mode
"Where should users edit **{primary_resource}**?

Options:
1. **Modal popup** - Overlay dialog, stays on same page (default)
2. **Dedicated page** - Separate URL like `/posts/123/edit`

Your answer (1-2 or name):"

#### Question 9: CSS Framework
"Which CSS framework would you like?

Options:
1. **Tailwind CSS** - Utility-first, highly customizable (default)
2. **Bulma** - Component-based, clean classes
3. **Pico CSS** - Classless, minimal, semantic HTML
4. **None** - I'll bring my own styles

Your answer (1-4 or name):"

#### Question 10: Related Resources Details

For each related resource mentioned, ask:

"For **{related_resource}**:

**Fields**: List fields (or 'auto' to let me infer)
Your answer:

**Relationship to {primary_resource}**:
- **belongs_to** - Each {related_resource} has ONE {primary_resource}
- **has_many** - Each {primary_resource} has MANY {related_resource}

Your answer:"

---

### Phase 4: Complete Preview & Confirmation

Show the complete plan before executing anything:

"📋 **Complete Plan**

## App Configuration

- **Name**: {app_name}
- **Domain**: {domain}
- **CSS Framework**: {css_framework}
- **Authentication**: {auth_method or "None"}

## Resources

### 1. {primary_resource}
- **Fields**: {fields}
- **Pagination**: {pagination_style}
- **Edit Mode**: {edit_mode}
- **Test Records**: {count}

### 2. {related_resource} (if any)
- **Fields**: {fields}
- **Relationship**: {relationship} {primary_resource}
- **Test Records**: {count}

---

## Commands I Will Execute

**IMPORTANT**: Create the app in the current working directory (CWD), NOT in /tmp or a git worktree.
New applications should be created directly in the user's current location.

\`\`\`bash
# 1. Create app in current directory
lvt new {app_name} --kit {kit}
cd {app_name}

# 2. Generate primary resource
lvt gen resource {primary_resource} {fields} \\
  --pagination {pagination_style} \\
  --edit-mode {edit_mode}

# 3. Generate related resource(s)
lvt gen resource {related_resource} {fields_with_foreign_key}

# 4. Add authentication (if requested)
lvt gen auth {options}

# 5. Apply database migrations
lvt migration up
cd internal/database && sqlc generate && cd ../..
go mod tidy

# 6. Generate test data
lvt seed {primary_resource} --count {count}
lvt seed {related_resource} --count {count}

# 7. Start development server
PORT=8080 go run cmd/{app_name}/main.go
\`\`\`

---

## Files That Will Be Created

\`\`\`
{app_name}/
├── cmd/{app_name}/
│   └── main.go                    ← Entry point
├── internal/
│   ├── app/
│   │   ├── home/                  ← Homepage
│   │   ├── {primary_resource}/
│   │   │   ├── handler.go         ← HTTP handlers
│   │   │   ├── routes.go          ← URL routes
│   │   │   └── {primary_resource}.tmpl  ← HTML template
│   │   ├── {related_resource}/
│   │   │   └── ... (same structure)
│   │   └── auth/ (if enabled)
│   │       ├── handler.go
│   │       ├── middleware.go
│   │       ├── login.tmpl
│   │       └── register.tmpl
│   └── database/
│       ├── schema.sql              ← Migrations
│       ├── queries/{resource}.sql  ← SQL queries
│       └── models.go               ← Generated models
├── web/
│   └── assets/
│       ├── css/                    ← Stylesheets
│       └── js/                     ← JavaScript
├── go.mod
├── go.sum
└── .lvtrc                          ← Project config
\`\`\`

---

**Ready to proceed?**

**1. Yes, create it!** → I'll execute all commands
**2. Let me customize** → Modify steps before executing
**3. Cancel** → Don't create anything

Your choice:"

**Based on answer:**
- **1 (Yes)** → Execute Phase 5
- **2 (Customize)** → Ask what to change, update plan, show preview again
- **3 (Cancel)** → Confirm cancellation, stop

---

### Phase 5: Execution

Execute commands sequentially, showing progress:

"🚀 **Creating your {app_name} app...**

⏳ Step 1/7: Creating app structure...
✅ Created app with `lvt new {app_name}`

⏳ Step 2/7: Generating {primary_resource} resource...
✅ Generated {primary_resource} with {field_count} fields

⏳ Step 3/7: Generating {related_resource} resource...
✅ Generated {related_resource} with foreign key to {primary_resource}

⏳ Step 4/7: Adding authentication system...
✅ Generated auth with {auth_method}

⏳ Step 5/7: Applying database migrations...
✅ Migrations applied, models generated with sqlc

⏳ Step 6/7: Generating test data...
✅ Seeded {count} {primary_resource} records
✅ Seeded {count} {related_resource} records

⏳ Step 7/7: Starting development server...
✅ Server running at http://localhost:8080

---

🎉 **Success! Your {app_name} is ready!**

## What Just Happened

I created a complete LiveTemplate application with:
- ✅ {resource_count} database-backed resources
- ✅ {auth_status} authentication system
- ✅ {count} test records for each resource
- ✅ Full CRUD operations (Create, Read, Update, Delete)
- ✅ Real-time updates via WebSockets
- ✅ {css_framework} styling

## Try It Out

1. **Visit**: http://localhost:8080/{primary_resource}
2. **Create** new {primary_resource}
3. **Edit** existing {primary_resource}
4. **Delete** {primary_resource}
5. **Check database**: `cat internal/database/db.sqlite | head`

## Next Steps

Want to add more features? Just ask:
- \"Add search to {primary_resource}\"
- \"Add file upload for {primary_resource} images\"
- \"Add a dashboard view\"
- \"Add email notifications\"
- \"Deploy to production\"

---

**Note**: LiveTemplate context is now established. You can use generic prompts like \"add search\" without saying \"lvt\" or \"livetemplate\" each time."

---

## Error Handling

If any command fails during execution:

### 1. Show Exact Error
Display the command output verbatim in a code block.

### 2. Explain in Plain Language
"❌ **Failed at Step {N}: {step_name}**

The `{command}` command failed with this error:

\`\`\`
{error_output}
\`\`\`

**What this means**: {plain_language_explanation}"

### 3. Suggest Fix
Based on common error patterns:

**Error: `lvt: command not found`**
- **Cause**: LiveTemplate CLI not installed
- **Fix**: `go install github.com/livetemplate/lvt@latest`
- **Then**: Add `$GOPATH/bin` to PATH, or run `export PATH=$PATH:$(go env GOPATH)/bin`

**Error: `directory already exists`**
- **Cause**: An app with name `{app_name}` already exists
- **Fix Option 1**: Choose different name: \"Let's use {app_name}2 instead\"
- **Fix Option 2**: Delete existing: `rm -rf {app_name}` (⚠️ loses data!)

**Error: `go mod tidy failed`**
- **Cause**: Network issues or Go proxy problems
- **Fix**: Check internet connection, retry with `GOPROXY=direct go mod tidy`

**Error: `go build` failed**
- **Cause**: Code generation or dependency issues
- **Fix**: Run `cd internal/database && sqlc generate` manually
- **Then**: Try `go mod tidy` again

**Error: `port already in use`**
- **Cause**: Another process using port 8080
- **Fix Option 1**: Kill process: `lsof -ti:8080 | xargs kill`
- **Fix Option 2**: Use different port: `PORT=8081 go run cmd/{app_name}/main.go`

### 4. Offer to Retry
"Would you like me to:
1. **Retry** after you fix the issue
2. **Skip this step** and continue with next steps
3. **Cancel** the whole operation

Your choice:"

---

## Context Persistence

After brainstorming completes successfully:

✅ **LiveTemplate context established**
- You're now in a LiveTemplate project (`.lvtrc` exists)
- Generic prompts will use LiveTemplate skills automatically
- No need to say "lvt" or "livetemplate" in every message

**Examples of what you can now say:**
- "Add a dashboard view"
- "Add search functionality"
- "Generate more test data"
- "Add file uploads"
- "Deploy this to production"

All of these will use the appropriate `lvt-*` skills automatically.

---

## Domain-Specific Intelligence

### Blog Domain
**Typical structure:**
- Primary: `posts` (title, content, published_at, author_id)
- Related: `comments` (post_id, content, author), `categories`, `tags`
- Auth: Yes (authors need accounts)
- Features: Rich text editing, categories, tags, comments

### E-commerce Domain
**Typical structure:**
- Primary: `products` (name, price, quantity, image_url)
- Related: `orders` (user_email, total, status), `cart_items`, `reviews`
- Auth: Optional (guest checkout vs user accounts)
- Features: Shopping cart, checkout, payment integration, inventory

### Todo/Task Management
**Typical structure:**
- Primary: `tasks` (title, description, due_date, completed, user_id)
- Related: `projects` (name, description), `labels`, `assignments`
- Auth: Yes (each user has their own tasks)
- Features: Filtering, sorting, due dates, priorities

### SaaS/Multi-tenant
**Typical structure:**
- Primary: `organizations` or `workspaces`
- Related: `users` (org_id), `projects`, `subscriptions`
- Auth: Yes (complex: users belong to orgs)
- Features: Roles, permissions, billing, teams

### CRM
**Typical structure:**
- Primary: `contacts` (name, email, company, phone)
- Related: `deals` (contact_id, value, stage), `activities`, `notes`
- Auth: Yes (sales team accounts)
- Features: Pipeline stages, activity tracking, reporting

**Use domain detection** to provide smart defaults when asking questions.

---

## Testing Considerations

**Manual Verification** (skill invocation cannot be automated):

Testing checklist for developers:

- [ ] Skill activates for "help me plan a livetemplate app"
- [ ] Skill DOES NOT activate for "help me plan an app" (no keywords)
- [ ] Phase 1: Core questions are asked (5 questions)
- [ ] Phase 2: Summary shown, "more options" offered
- [ ] Phase 3: Detailed questions asked if user says yes
- [ ] Phase 4: Complete preview shown with all commands and files
- [ ] User can cancel at preview phase → nothing created
- [ ] Phase 5: Commands execute successfully
- [ ] Error handling shows helpful messages
- [ ] Context persists after completion (generic prompts work)

**Structure Validation** (can be automated):

See `e2e/agent_skills_validation_test.go`:
- Skill exists in `skills/brainstorm/SKILL.md`
- Skill has valid frontmatter (name, description, keywords)
- Skill name follows format: `lvt-brainstorm`

---

## Version History

- **v1.0.0** (2025-11-28): Initial implementation
  - Progressive disclosure (3-5 core → offer more → 8-12 detailed)
  - Always requires keywords (prevents false positives)
  - Domain-specific intelligence for common app types
  - Complete preview before execution
  - Error handling with plain-language explanations
