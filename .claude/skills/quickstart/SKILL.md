---
name: lvt-quickstart
description: Rapid end-to-end workflow - creates app, adds resources, sets up development environment in one flow
keywords: ["lvt", "livetemplate", "lt"]
category: workflows
version: 1.1.0
---

# lvt-quickstart

Get from zero to working app in minutes. This workflow chains multiple skills to create a complete working application with resources and development environment ready.

## ⚠️ MANDATORY: PLAN FIRST, THEN EXECUTE

**NEVER start executing commands without user approval.**

```
User request
     ↓
┌─────────────────────────────┐
│ 1. INFER from request       │
│    - App name               │
│    - Domain (blog/shop/etc) │
│    - Resources needed       │
│    - Auth (if mentioned)    │
└─────────────────────────────┘
     ↓
┌─────────────────────────────┐
│ 2. PRESENT complete plan    │
│    - All settings in table  │
│    - Commands to execute    │
│    - "Ready?" prompt        │
└─────────────────────────────┘
     ↓
┌─────────────────────────────┐
│ 3. WAIT for user response   │  ← MANDATORY STOP
│    - "yes" → Execute        │
│    - "change X" → Update    │
│    - "no" → Cancel          │
└─────────────────────────────┘
     ↓
Execute ONLY after approval
```

**Example - What you MUST do:**

User: "create a [domain] app with [features] using livetemplate"

You respond with a plan like this:
```
📋 **Plan for your [domain] app**

| Setting | Value |
|---------|-------|
| App name | [inferred from domain] |
| Primary resource | [domain-appropriate resource] |
| Fields | [domain-appropriate fields] |
| Authentication | [if requested: Password/Magic Link/None] |
| Kit | multi (Tailwind CSS) |
| Test data | 50 records |

**Commands I'll run:**

lvt new [name] --kit multi
cd [name]
lvt gen resource [resource] [fields...]
[lvt gen auth - if auth requested]
lvt migration up
go mod tidy
lvt seed [resource] --count 50

**Ready to create?**
- **yes** - proceed with this plan
- **change X** - modify a setting (e.g., "change name to myblog", "add comments resource")
- **advanced** - explore more options (kit, pagination, edit mode, CSS framework)
- **no** - cancel
```

Then WAIT. Do not execute until user approves.

**If user says "advanced"**, show additional options:
```
⚙️ **Advanced Options**

| Option | Current | Alternatives |
|--------|---------|--------------|
| Kit | multi | single, simple |
| CSS Framework | tailwind | pico, bulma, bootstrap |
| Pagination | infinite scroll | page numbers |
| Edit Mode | modal | inline, page |
| Database | sqlite | postgres (requires setup) |

What would you like to change?
```

**Concrete examples by domain:**
- "shop app" → products resource with name, price, quantity
- "blog app" → posts resource with title, content, published
- "todo app" → tasks resource with title, description, completed, due_date
- "crm app" → contacts resource with name, email, company, phone

---

## 🎯 ACTIVATION RULES

### Context Detection

This skill activates when **LiveTemplate context is established**:

**✅ Context Established By:**

1. **Project context** - `.lvtrc` file exists in current directory
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Priority**: Project context > Agent context > Keyword context

### Keyword Matching

**Accepted keywords** (case-insensitive, whole words):
- `lvt`
- `livetemplate`
- `lt`

**Accepted patterns:**
- "create * {with|using|via} {lvt|livetemplate|lt}"
- "{lvt|lt} {quickstart|new|create|build} *"
- "use {livetemplate|lvt} to *"
- "quickstart * {with|using} {livetemplate|lvt}"

### Trigger Patterns

**With Context (any of: .lvtrc OR agent OR keywords):**
✅ "quickstart a blog"
✅ "create a quick shop"
✅ "build me a working todo app"

**Without Context (must include keywords):**
✅ "quickstart a blog with lvt"
✅ "use livetemplate to create a shop"
✅ "lt quickstart for todo app"
❌ "create a blog app" (no context, no keywords)

### Examples by Scenario

**Scenario 1: New conversation, no .lvtrc, no agent**
- User: "Create a quick blog app"
  → ❌ No context, no keywords → Don't activate

- User: "Quickstart a blog with livetemplate"
  → ✅ Keywords found → Activate skill
  → ✅ Context now established for conversation

**Scenario 2: Existing project (.lvtrc exists)**
- User: "Quickstart a blog"
  → ✅ Project context → Activate skill

**Scenario 3: Using lvt-assistant agent**
- User (in agent): "Build a quick shop"
  → ✅ Agent context → Activate skill

**Scenario 4: Context persistence**
- User: "Use lvt to build a blog"
  → ✅ Keywords → Activate skill
  → ✅ Context established

- User: "Add authentication"
  → ✅ Context persists → Other skills activate

---

## 💡 Quickstart vs Brainstorming

**Quickstart (this skill):**
- User: "quickstart a blog with lvt"
- → Shows plan with sensible defaults
- → Waits for approval, then executes
- → Best when user knows what they want

**Brainstorming (guided discovery):**
- User: "help me plan a livetemplate blog"
- → Asks questions to understand requirements
- → Explores options together
- → Then shows plan and executes

**Both skills present a plan and wait for approval before executing.**

---

## User Prompts

**When to use:**
- "Create a quick [type] app"
- "I want to start a [name] project fast"
- "Quickstart a [domain] application"
- "Build me a working [type] app"
- "I need a [name] app up and running"

**Examples:**
- "Quickstart a blog app"
- "Create a quick todo application"
- "I want to start a shop project fast"
- "Build me a working task manager"

## Workflow Steps

This skill chains together:
1. **Present Plan** - Show complete plan and wait for approval
2. **lvt:new-app** - Create application
3. **lvt:add-resource** - Add initial resource(s)
4. **lvt:run-and-test** - Start dev server
5. **lvt:seed-data** (optional) - Add test data

### Step 1: Present Plan and Wait for Approval

**This step is MANDATORY. Do NOT skip it.**

1. Extract from user request: app name, domain, resources, auth needs
2. Apply domain defaults (see Domain Detection below)
3. Present complete plan in table format
4. End with "Ready to create?" prompt
5. **WAIT for user response before proceeding**

### Step 2: Understand Requirements (after approval)

Extract from user request:
- App name
- Domain/type (blog, todo, shop, tasks, etc.)
- Initial resources needed

**Domain Detection:**
- "blog" → posts, comments
- "todo/tasks" → tasks
- "shop/store" → products, orders
- "project management" → projects, tasks
- "social" → users, posts, likes
- "forum" → topics, replies

### Step 2: Create Application

Use **lvt:new-app** skill:
```bash
lvt new <app-name>
cd <app-name>
```

Choose kit based on requirements:
- Complex apps → multi kit (Tailwind)
- SPAs → single kit (Tailwind)
- Simple/prototypes → simple kit (Pico)

### Step 3: Add Initial Resource

Use **lvt:add-resource** skill:

**For blog:**
```bash
lvt gen resource posts title content published
```

**For todo app:**
```bash
lvt gen resource tasks title description due_date completed
```

**For shop:**
```bash
lvt gen resource products name price quantity image_url
```

Apply migrations:
```bash
lvt migration up
cd database && sqlc generate && cd ../..
go mod tidy
```

### Step 4: Add Related Resources (if applicable)

For domains with multiple resources, suggest adding related ones:

**Blog:**
```bash
# Add comments with foreign key to posts
lvt gen resource comments post_id:references:posts:CASCADE content author
lvt migration up
cd database && sqlc generate && cd ../..
```

**Shop:**
```bash
# Add orders
lvt gen resource orders user_email:string total:float status:string
lvt migration up
cd database && sqlc generate && cd ../..
```

### Step 5: Seed Test Data (Optional)

Use **lvt:seed-data** skill:
```bash
lvt seed <resource> --count 10
```

### Step 6: Start Development

Use **lvt:run-and-test** skill:
```bash
lvt serve
# Opens browser automatically at http://localhost:3000
```

## Quick Reference

### Blog App (2 resources)
```bash
lvt new myblog
cd myblog
lvt gen resource posts title content published
lvt gen resource comments post_id:references:posts:CASCADE content author
lvt migration up
cd database && sqlc generate && cd ../..
go mod tidy
lvt seed posts --count 10
lvt seed comments --count 30
lvt serve
```

### Todo App (1 resource)
```bash
lvt new mytodos
cd mytodos
lvt gen resource tasks title description due_date completed
lvt migration up
cd database && sqlc generate && cd ../..
go mod tidy
lvt seed tasks --count 20
lvt serve
```

### Shop App (2 resources)
```bash
lvt new myshop
cd myshop
lvt gen resource products name price:float quantity:int image_url
lvt gen resource orders user_email total:float status
lvt migration up
cd database && sqlc generate && cd ../..
go mod tidy
lvt seed products --count 50
lvt seed orders --count 100
lvt serve
```

## Checklist

**Plan Phase (MANDATORY):**
- [ ] Extract app name and domain from user request
- [ ] Detect domain type and determine resources
- [ ] Apply domain defaults (fields, auth, kit, seed count)
- [ ] Present complete plan in table format
- [ ] Show exact commands that will be run
- [ ] End with "Ready to create?" prompt
- [ ] **WAIT for user approval before proceeding**

**Execution Phase (only after approval):**
- [ ] Use lvt:new-app to create application
- [ ] Verify app created successfully
- [ ] Use lvt:add-resource for primary resource
- [ ] Run migrations and generate models
- [ ] Add auth if included in plan
- [ ] Seed test data
- [ ] Use lvt:run-and-test to start dev server
- [ ] Verify app runs and is accessible
- [ ] Show user the URL and next steps

## Domain-Specific Guidance

### Blog Domain
**Primary resource:** posts (title, content, published)
**Related resources:** comments, categories, tags, authors
**Common views:** about, contact
**Auth needed:** Yes (for admin/author features)

### Todo/Tasks Domain
**Primary resource:** tasks (title, description, due_date, completed)
**Related resources:** projects, labels, users
**Common views:** dashboard (stats)
**Auth needed:** Yes (for user tasks)

### E-commerce Domain
**Primary resource:** products (name, price, quantity, image_url)
**Related resources:** orders, cart_items, customers
**Common views:** about, shipping, returns
**Auth needed:** Yes (for checkout)

### Project Management Domain
**Primary resource:** projects (name, description, status)
**Related resources:** tasks, team_members, milestones
**Common views:** dashboard, analytics
**Auth needed:** Yes (for teams)

## Success Criteria

Quickstart is successful when:
1. ✅ App created and builds without errors
2. ✅ Primary resource generated and working
3. ✅ Migrations applied successfully
4. ✅ Dev server running
5. ✅ User can see working CRUD interface
6. ✅ Test data populated (if requested)

## Time Estimates

- **Simple (1 resource):** 2-3 minutes
- **Medium (2 resources):** 4-5 minutes
- **Complex (3+ resources):** 6-8 minutes

## Common Patterns

### Pattern 1: Parent-Child Resources
```bash
# Parent
lvt gen resource posts title content

# Child with FK
lvt gen resource comments post_id:references:posts:CASCADE content
```

### Pattern 2: Many-to-Many
```bash
# Create junction table with gen schema
lvt gen schema post_tags post_id:references:posts tag_id:references:tags
```

### Pattern 3: User-Owned Resources
```bash
# Resource with user ownership
lvt gen resource tasks user_id:references:users title description
```

## Notes

- This is a meta-skill that chains other skills
- Always verify each step before proceeding to next
- Offer choices when multiple options exist
- Keep user informed of progress
- Show final URL and next steps
- Perfect for demos and prototypes
- Can be extended with auth (see lvt:production-ready)
