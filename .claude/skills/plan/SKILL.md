---
name: lvt-plan
description: "Present a complete plan for new LiveTemplate apps. NO QUESTIONS - just fill in the template below and show it."
keywords: ["lvt", "livetemplate", "lt", "app", "application", "create", "build", "make", "new", "plan"]
requires_keywords: true
category: workflows
version: 4.0.0
---

# OUTPUT THIS EXACTLY (fill in blanks from table below):

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

# FILL IN BLANKS FROM THIS TABLE:

| If user mentions | DOMAIN | NAME | RESOURCE | FIELDS | AUTH |
|------------------|--------|------|----------|--------|------|
| blog | blog | blog | posts | title:string content:text published:bool | Password |
| shop/store/ecommerce | e-commerce | shop | products | name:string description:text price:float | Password |
| todo | todo | todo | tasks | title:string completed:bool due_date:time | Password |
| crm/contacts | CRM | crm | contacts | name:string email:string company:string | Password |
| forum | forum | forum | topics | title:string content:text | Password |
| (other) | app | app | items | name:string description:text | None |

# CHAINED COMMAND TEMPLATES

**WITH auth** (use when AUTH = Password):
```
lvt new ___NAME___ --kit multi && \
  cd ___NAME___ && \
  lvt gen resource ___RESOURCE___ ___FIELDS___ && \
  lvt gen auth && \
  lvt migration up && \
  go mod tidy && \
  lvt seed ___RESOURCE___ --count 50
```

**WITHOUT auth** (use when AUTH = None):
```
lvt new ___NAME___ --kit multi && \
  cd ___NAME___ && \
  lvt gen resource ___RESOURCE___ ___FIELDS___ && \
  lvt migration up && \
  go mod tidy && \
  lvt seed ___RESOURCE___ --count 50
```

# RULES

- If user mentions "auth" or "authentication" → AUTH = Password, use WITH auth template
- If user gives a specific name like "myblog" → use that for NAME
- If no auth mentioned and domain doesn't require it → AUTH = None, use WITHOUT auth template

# IMPORTANT

- DO NOT ask questions
- DO NOT say "let me gather details"
- DO NOT ask for app name, auth type, or anything else
- JUST output the filled template above
- User can modify after seeing the plan
