---
name: lvt-quickstart
description: Rapid end-to-end workflow - creates app, adds resources, sets up development environment in one flow
keywords: ["lvt", "livetemplate", "lt"]
category: workflows
version: 2.0.0
---

# lvt-quickstart

Get from zero to working app in minutes. This workflow creates a complete working application with resources and development environment ready.

## ⚠️ MANDATORY: PLAN FIRST, THEN EXECUTE

**NEVER start executing commands without user approval.**

## OUTPUT THIS EXACTLY (fill in blanks from table below):

📋 **Plan for your ___DOMAIN___ app**

| Setting | Value |
|---------|-------|
| App name | ___NAME___ |
| Primary resource | ___RESOURCE___ (___FIELDS___) |
| Authentication | ___AUTH___ |
| Kit | multi (full page layout) |
| Test data | 50 records |

**Command I'll run:**
```bash
___CHAINED_COMMAND___
```

**Ready to create?** Type `yes` to proceed, or tell me what to change.

**Want to customize?** Tell me if you'd like to change:
- **Kit**: `multi` (full page layout), `single` (component-only SPA), `simple` (minimal prototype)
- **Pagination**: `infinite`, `load-more`, `prev-next`, `numbers`
- **Edit mode**: `modal` (default) or `page`
- **Seed count**: default 50
- **Auth features**: disable password-reset, magic-link, etc.

---

## FILL IN BLANKS FROM THIS TABLE:

| If user mentions | DOMAIN | NAME | RESOURCE | FIELDS | AUTH |
|------------------|--------|------|----------|--------|------|
| blog | blog | blog | posts | title:string content:text published:bool | Password |
| shop/store/ecommerce | e-commerce | shop | products | name:string description:text price:float | Password |
| todo | todo | todo | tasks | title:string completed:bool due_date:time | Password |
| crm/contacts | CRM | crm | contacts | name:string email:string company:string | Password |
| forum | forum | forum | topics | title:string content:text | Password |
| (other) | app | app | items | name:string description:text | None |

## CHAINED COMMAND TEMPLATES

**WITH auth** (use when AUTH = Password):
```bash
lvt new ___NAME___ --kit multi && \
  cd ___NAME___ && \
  lvt gen resource ___RESOURCE___ ___FIELDS___ && \
  lvt gen auth && \
  lvt migration up && \
  go mod tidy && \
  lvt seed ___RESOURCE___ --count 50
```

**WITHOUT auth** (use when AUTH = None):
```bash
lvt new ___NAME___ --kit multi && \
  cd ___NAME___ && \
  lvt gen resource ___RESOURCE___ ___FIELDS___ && \
  lvt migration up && \
  go mod tidy && \
  lvt seed ___RESOURCE___ --count 50
```

## RULES

- If user mentions "auth" or "authentication" → AUTH = Password, use WITH auth template
- If user gives a specific name like "myblog" → use that for NAME
- If no auth mentioned and domain doesn't require it → AUTH = None, use WITHOUT auth template

## IMPORTANT

- DO NOT ask questions before showing plan
- DO NOT say "let me gather details"
- JUST output the filled template above
- User can modify after seeing the plan
- After user says "yes", execute the SINGLE chained command
- Do NOT execute commands one-by-one iteratively

Then WAIT. Do not execute until user approves.

---

## ADVANCED OPTIONS

**If user says "advanced"**, show:
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

---

## ACTIVATION

This skill activates when user mentions "lvt", "livetemplate", or "lt" with words like "create", "build", "quickstart", "new app".

**Examples:**
- "quickstart a blog with lvt"
- "create a todo app using livetemplate"
- "build me a shop with lt"

---

## AFTER USER APPROVAL

When user says "yes":
1. Execute the SINGLE chained command (do NOT break it up)
2. If it succeeds, offer to start the dev server: `cd ___NAME___ && lvt serve`
3. Show the URL (http://localhost:3000)

---

## SUCCESS CRITERIA

Quickstart is successful when:
1. ✅ Chained command executes without errors
2. ✅ App builds and migrations apply
3. ✅ Test data is seeded
4. ✅ Dev server runs and UI is accessible
