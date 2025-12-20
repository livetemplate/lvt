---
name: lvt-plan
description: "Use when creating new LiveTemplate/lvt apps - this is THE definitive skill for lvt app creation. Supersedes generic brainstorming/planning skills. Triggers: 'create/build/make a lvt/livetemplate [type] app', 'plan a livetemplate app'."
keywords: ["lvt", "livetemplate", "lt", "app", "application", "create", "build", "make", "new", "plan"]
requires_keywords: true
category: workflows
version: 1.3.0
---

# lvt-plan

Plan-first skill for creating LiveTemplate applications. Presents a complete plan with sensible defaults upfront, lets user modify before execution.

## SKILL PRIORITY

**This is THE definitive skill for LiveTemplate app creation. Use it directly - do NOT use generic brainstorming or planning skills first.**

When a user says "create a lvt blog app" or similar:
- Use `lvt-plan` directly - it handles the planning process
- Do NOT use generic creative/brainstorming skills first

**Why**: This skill is domain-specific for LiveTemplate and provides a tailored experience with lvt-specific defaults (resources, auth, kits, pagination, etc.). Generic brainstorming skills lack this context.

---

## ACTIVATION RULES

### Keyword Matching

**Accepted keywords** (case-insensitive):
- `lvt`
- `livetemplate`
- `lt`

**Will Activate:**
- "Create a **lvt** blog app"
- "Build a **livetemplate** shop with auth"
- "Make an **lt** todo application"
- "Help me plan a **livetemplate** app"

**Won't Activate:**
- "Help me plan a blog" (no keywords)
- "Create an app" (no keywords)
- "Create blog app with posts(title,content)" (detailed spec → use new-app skill directly)

---

## THE FLOW: Plan-First

```
User request
     ↓
┌─────────────────────────────┐
│ 1. INFER from request       │
│    - App name               │
│    - Domain (blog/shop/etc) │
│    - Auth (if mentioned)    │
└─────────────────────────────┘
     ↓
┌─────────────────────────────┐
│ 2. APPLY domain defaults    │
│    - Primary resource       │
│    - Fields                 │
│    - Typical settings       │
└─────────────────────────────┘
     ↓
┌─────────────────────────────┐
│ 3. PRESENT complete plan    │
│    - All settings in table  │
│    - Commands to execute    │
│    - "Ready?" prompt        │
└─────────────────────────────┘
     ↓
User: "yes" → Execute
User: "change X" → Update plan, show again
User: "no" → Cancel
```

---

## Step 1: Infer from User Request

Extract as much as possible from the initial prompt:

| Pattern | Extract |
|---------|---------|
| "create a lvt **blog** app" | name=blog, domain=Blog |
| "build **myblog** with lvt" | name=myblog, domain=Blog |
| "lvt **shop** with auth" | name=shop, domain=E-commerce, auth=yes |
| "make a **todo-app** using lt" | name=todo-app, domain=Todo |
| "create a lvt app" | name=?, domain=? (need to ask) |

**If domain is unclear**: Ask ONE question: "What type of app are you building? (blog, shop, todo, crm, or describe it)"

---

## Step 2: Apply Domain Defaults

Use domain-specific intelligence to fill in all settings:

### Blog Domain
```
name: blog (or extracted)
resource: posts
fields: title:string, content:text, published:bool
auth: optional (default: no)
kit: multi
seed: 50
```

### E-commerce Domain
```
name: shop (or extracted)
resource: products
fields: name:string, description:text, price:float, quantity:int
auth: optional (default: no)
kit: multi
seed: 50
```

### Todo Domain
```
name: todo (or extracted)
resource: tasks
fields: title:string, description:text, completed:bool, due_date:time
auth: yes (each user has own tasks)
kit: multi
seed: 50
```

### CRM Domain
```
name: crm (or extracted)
resource: contacts
fields: name:string, email:string, company:string, phone:string
auth: yes (sales team accounts)
kit: multi
seed: 50
```

### Forum Domain
```
name: forum (or extracted)
resource: topics
fields: title:string, content:text, pinned:bool
auth: yes
kit: multi
seed: 50
```

### Generic/Unknown Domain
```
name: app (or extracted)
resource: items
fields: name:string, description:text
auth: no
kit: multi
seed: 50
```

---

## Step 3: Present Complete Plan

Show the FULL plan immediately. Do not ask questions one-by-one.

**Template:**

```
📋 **Plan for your {domain} app**

| Setting | Value |
|---------|-------|
| App name | {name} |
| Primary resource | {resource} |
| Fields | {fields} |
| Authentication | {auth_description} |
| Kit | {kit} ({css_framework}) |
| Test data | {seed_count} records |

**Commands I'll run:**

lvt new {name} --kit {kit}
cd {name}
lvt gen resource {resource} {fields}
{auth_command if auth}
lvt migration up
go mod tidy
lvt seed {resource} --count {seed_count}

**Ready to create?**
- **yes** - proceed with this plan
- **change X** - modify a setting
- **advanced** - explore more options
- **no** - cancel
```

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

**Example - Blog with auth:**

```
📋 **Plan for your blog app**

| Setting | Value |
|---------|-------|
| App name | blog |
| Primary resource | posts |
| Fields | title:string, content:text, published:bool |
| Authentication | Password (email + password) |
| Kit | multi (Tailwind CSS) |
| Test data | 50 records |

**Commands I'll run:**

lvt new blog --kit multi
cd blog
lvt gen resource posts title:string content:text published:bool
lvt gen auth
lvt migration up
go mod tidy
lvt seed posts --count 50

**Ready to create?**
- **yes** - proceed with this plan
- **change X** - modify a setting
- **advanced** - explore more options
- **no** - cancel
```

---

## Handling User Responses

### "yes" / "y" / "go" / "create" / "do it"
Execute the plan immediately. Proceed to Execution phase.

### "no" / "cancel" / "stop"
Acknowledge and stop: "No problem. Let me know when you're ready to create an app."

### "advanced" / "options" / "customize"
Show the advanced options table:

| Option | Description | Values |
|--------|-------------|--------|
| Kit | Project structure | multi (recommended), single, simple |
| CSS Framework | Styling | tailwind (default), pico, bulma, bootstrap |
| Pagination | List navigation | infinite (default), page |
| Edit Mode | How items are edited | modal (default), inline, page |
| Database | Data storage | sqlite (default), postgres |

After user selects options, update the plan and show it again.

### Modification Requests
Update the plan and show it again. Examples:

| User says | Action |
|-----------|--------|
| "no auth" | Remove auth, show updated plan |
| "add auth" | Add `lvt gen auth`, show updated plan |
| "add comments" | Add comments resource with post_id reference |
| "call it myblog" | Change app name to myblog |
| "use pico" / "simple kit" | Change kit to simple (Pico CSS) |
| "100 records" | Change seed count to 100 |
| "add categories" | Add categories resource |

After updating, show the plan again with: "Updated plan: ... Ready to create?"

---

## Execution Phase

Execute commands sequentially, showing progress:

```
🚀 **Creating your {name} app...**

⏳ Creating app structure...
✅ Created with `lvt new {name} --kit {kit}`

⏳ Generating {resource} resource...
✅ Generated {resource} with {field_count} fields

⏳ Adding authentication... (if applicable)
✅ Auth system added

⏳ Running migrations...
✅ Database ready

⏳ Seeding test data...
✅ Added {seed_count} {resource} records

🎉 **Done! Your {name} app is ready.**

Start the server:
  cd {name}
  go run cmd/{name}/main.go

Then visit: http://localhost:8080/{resource}
```

---

## Error Handling

If any command fails:

1. **Show exact error** in code block
2. **Explain in plain language**
3. **Suggest fix**
4. **Offer to retry**

Common errors:

| Error | Cause | Fix |
|-------|-------|-----|
| `lvt: command not found` | CLI not installed | `go install github.com/livetemplate/lvt@latest` |
| `directory already exists` | Name taken | Choose different name or delete existing |
| `go mod tidy failed` | Network issue | Retry with `GOPROXY=direct go mod tidy` |
| `port already in use` | Port 8080 busy | Use `PORT=8081 go run ...` |

---

## Context Persistence

After successful creation:
- `.lvtrc` exists in the new app directory
- User can use generic prompts: "add search", "add comments", "deploy"
- No need to say "lvt" or "livetemplate" anymore

---

## Version History

- **v1.3.0** (2025-12-18): Plan-first approach
  - Present complete plan with defaults upfront
  - Single confirmation point instead of 5+ questions
  - User modifies plan with natural language

- **v1.2.0** (2025-12-18): Added skill priority language
  - Supersedes generic brainstorming skills

- **v1.0.0** (2025-11-28): Initial implementation
  - Progressive question flow (deprecated)
