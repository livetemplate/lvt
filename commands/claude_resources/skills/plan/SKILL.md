---
name: lvt-plan
description: "Use when creating new LiveTemplate/lvt apps - this is THE definitive skill for lvt app creation. Supersedes generic brainstorming/planning skills. Triggers: 'create/build/make a lvt/livetemplate [type] app', 'plan a livetemplate app'."
keywords: ["lvt", "livetemplate", "lt", "app", "application", "create", "build", "make", "new", "plan"]
requires_keywords: true
category: workflows
version: 2.0.0
---

# lvt-plan

## YOUR FIRST RESPONSE

When user asks to create an app, respond with this IMMEDIATELY (no questions):

```
📋 **Plan for your {domain} app**

| Setting | Value |
|---------|-------|
| App name | {name} |
| Primary resource | {resource} ({fields}) |
| Authentication | {auth_type} |
| Kit | multi (Tailwind CSS) |
| Test data | 50 records |

**Commands I'll run:**
lvt new {name} --kit multi
lvt gen resource {resource} {fields}
{lvt gen auth  # if auth requested}
lvt migration up && go mod tidy
lvt seed {resource} --count 50

**Ready to create?** (yes / change X / advanced)
```

Fill in `{placeholders}` using the Domain Defaults below.

---

## Domain Defaults

| Domain | name | resource | fields | auth |
|--------|------|----------|--------|------|
| Blog | blog | posts | title:string content:text published:bool | if mentioned |
| Shop/E-commerce | shop | products | name:string description:text price:float quantity:int | if mentioned |
| Todo | todo | tasks | title:string description:text completed:bool due_date:time | yes |
| CRM | crm | contacts | name:string email:string company:string phone:string | yes |
| Forum | forum | topics | title:string content:text pinned:bool | yes |
| Unknown | app | items | name:string description:text | no |

**Auth types:** Password (email + password), Magic Link (email only), None

---

## After User Responds

| User says | Action |
|-----------|--------|
| "yes" | Execute the commands |
| "change X" | Update plan, show again |
| "no auth" | Remove auth, show updated plan |
| "add comments" | Add resource with foreign key |
| "advanced" | Show kit/pagination/database options |
| "no" | Cancel |

---

## Execution

Show progress as you run each command:

```
🚀 Creating your {name} app...

⏳ Creating app structure...
✅ Created with lvt new {name} --kit multi

⏳ Generating {resource} resource...
✅ Generated {resource}

⏳ Running migrations...
✅ Database ready

🎉 Done! Run: cd {name} && go run cmd/{name}/main.go
```

---

## Version History

- **v2.0.0** (2025-12-21): Complete rewrite - minimal structure
  - Root cause: 450-line skill with "Step 1, 2, 3" implied sequential process
  - Now just: template + defaults table + response handling
  - No numbered steps, no phases, no wizard-like structure
  - Claude follows structure - so structure is now "show plan immediately"

- **v1.7.0** (2025-12-21): Add CRITICAL-STOP-AND-READ section (didn't fix it)
- **v1.6.0** (2025-12-20): Override brainstorming skills (didn't fix it)
- **v1.5.0** (2025-12-20): FORBIDDEN list (didn't fix it)
- **v1.4.0** (2025-12-20): Reinforce no-questions (didn't fix it)
- **v1.3.0** (2025-12-18): Plan-first approach (partial fix)
