---
name: lvt-plan
description: "Present a complete plan for new LiveTemplate apps. NO QUESTIONS - just fill in the template below and show it."
keywords: ["lvt", "livetemplate", "lt", "app", "application", "create", "build", "make", "new", "plan"]
requires_keywords: true
category: workflows
version: 3.0.0
---

# OUTPUT THIS EXACTLY (fill in blanks from table below):

📋 **Plan for your ___DOMAIN___ app**

| Setting | Value |
|---------|-------|
| App name | ___NAME___ |
| Primary resource | ___RESOURCE___ (___FIELDS___) |
| Authentication | ___AUTH___ |
| Kit | multi (Tailwind CSS) |
| Test data | 50 records |

**Commands I'll run:**
```
lvt new ___NAME___ --kit multi
lvt gen resource ___RESOURCE___ ___FIELDS___
___AUTH_COMMAND___
lvt migration up && go mod tidy
lvt seed ___RESOURCE___ --count 50
```

**Ready to create?** Type `yes` to proceed, or tell me what to change.

---

# FILL IN BLANKS FROM THIS TABLE:

| If user mentions | DOMAIN | NAME | RESOURCE | FIELDS | AUTH | AUTH_COMMAND |
|------------------|--------|------|----------|--------|------|--------------|
| blog | blog | blog | posts | title:string content:text published:bool | Password | lvt gen auth |
| shop/store/ecommerce | e-commerce | shop | products | name:string description:text price:float | Password | lvt gen auth |
| todo | todo | todo | tasks | title:string completed:bool due_date:time | Password | lvt gen auth |
| crm/contacts | CRM | crm | contacts | name:string email:string company:string | Password | lvt gen auth |
| forum | forum | forum | topics | title:string content:text | Password | lvt gen auth |
| (other) | app | app | items | name:string description:text | None | (remove line) |

- If user mentions "auth" or "authentication" → AUTH = Password, AUTH_COMMAND = lvt gen auth
- If user gives a specific name like "myblog" → use that for NAME
- If no auth mentioned and domain doesn't require it → AUTH = None, remove AUTH_COMMAND line

# IMPORTANT

- DO NOT ask questions
- DO NOT say "let me gather details"
- DO NOT ask for app name, auth type, or anything else
- JUST output the filled template above
- User can modify after seeing the plan
