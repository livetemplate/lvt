---
name: lvt:quickstart
description: Rapid end-to-end workflow - creates app, adds resources, sets up development environment in one flow
category: workflows
version: 1.0.0
---

# lvt:quickstart

Get from zero to working app in minutes. This workflow chains multiple skills to create a complete working application with resources and development environment ready.

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
1. **lvt:new-app** - Create application
2. **lvt:add-resource** - Add initial resource(s)
3. **lvt:run-and-test** - Start dev server
4. **lvt:seed-data** (optional) - Add test data

### Step 1: Understand Requirements

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
cd internal/database && sqlc generate && cd ../..
go mod tidy
```

### Step 4: Add Related Resources (if applicable)

For domains with multiple resources, suggest adding related ones:

**Blog:**
```bash
# Add comments with foreign key to posts
lvt gen resource comments post_id:references:posts:CASCADE content author
lvt migration up
cd internal/database && sqlc generate && cd ../..
```

**Shop:**
```bash
# Add orders
lvt gen resource orders user_email:string total:float status:string
lvt migration up
cd internal/database && sqlc generate && cd ../..
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
cd internal/database && sqlc generate && cd ../..
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
cd internal/database && sqlc generate && cd ../..
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
cd internal/database && sqlc generate && cd ../..
go mod tidy
lvt seed products --count 50
lvt seed orders --count 100
lvt serve
```

## Checklist

- [ ] Extract app name and domain from user request
- [ ] Detect domain type and suggest initial resources
- [ ] Use lvt:new-app to create application
- [ ] Verify app created successfully
- [ ] Use lvt:add-resource for primary resource
- [ ] Run migrations and generate models
- [ ] Suggest related resources based on domain
- [ ] Add related resources if user agrees
- [ ] Offer to seed test data
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
