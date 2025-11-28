---
name: lvt-add-resource
description: Add a CRUD resource to an existing LiveTemplate app with database schema, queries, handler, and template
keywords: ["lvt", "livetemplate", "lt"]
category: core
version: 1.0.0
---

# lvt-add-resource

Adds a full CRUD (Create, Read, Update, Delete) resource to an existing LiveTemplate application. This skill intelligently infers field types, generates database migrations, SQL queries, Go handlers, and HTML templates.

## 🎯 ACTIVATION RULES

### Context Detection

This skill typically runs in **existing LiveTemplate projects** (.lvtrc exists).

**✅ Context Established By:**
1. **Project context** - `.lvtrc` exists (most common scenario)
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Keyword matching** (case-insensitive): `lvt`, `livetemplate`, `lt`

### Trigger Patterns

**With Context:**
✅ "add a posts resource"
✅ "create a products resource with name and price"
✅ "generate CRUD for tasks"

**Without Context (needs keywords):**
✅ "add a posts resource to my lvt app"
✅ "use livetemplate to create a products resource"
❌ "add a posts resource" (no context, no keywords)

---

## User Prompts

This skill should activate when the user requests to add a resource to their app:

**Explicit prompts:**
- "Add a [resource] resource to my app"
- "Generate CRUD for [resource] with [fields]"
- "Create a [resource] with [field1], [field2], [field3]"
- "I need a [resource] resource with [field descriptions]"

**Implicit prompts:**
- "Let's add [resource] to the app"
- "I want to track [resources]"
- "Can you add [resource] functionality?"
- "Create [resource] management"

**Examples:**
- "Add a posts resource with title, content, and published"
- "Generate CRUD for products with name, price, quantity"
- "Create a users resource with name email password"
- "I need tasks with title description due_date completed"

## Context Awareness

Before executing this skill, verify:

1. **In Project Directory:**
   - Check for `.lvtrc` file (confirms it's an lvt project)
   - Check for `go.mod` (confirms it's a Go project)
   - Check for `internal/database/` directory

2. **Dependencies Available:**
   - `lvt` binary is installed and accessible
   - Project was created with `lvt new`

3. **Not Already Exists:**
   - Check if resource directory already exists
   - Warn user if resource name conflicts

## Checklist

- [ ] **Step 1:** Verify we're in an lvt project directory
  - Check for `.lvtrc` file
  - Check for `go.mod` file
  - Check for `internal/database/` directory
  - If missing, inform user they need to create an app first (use lvt:new-app skill)

- [ ] **Step 2:** Validate prerequisites
  - Verify `lvt` command is available
  - Check current directory is project root

- [ ] **Step 3:** Extract resource details from user request
  - Resource name (singular form preferred, e.g., "post", "user", "task")
  - Field list with types (explicit or to be inferred)

- [ ] **Step 4:** Parse and organize fields
  - If user provides field:type format, use as-is
  - If user provides just field names, rely on lvt's type inference
  - Validate field names are valid Go identifiers

- [ ] **Step 5:** Check for naming conflicts
  - Check if `internal/app/[resource]/` already exists
  - If exists, ask user if they want to:
    - Overwrite existing resource (warning: data loss)
    - Choose a different name
    - Cancel operation

- [ ] **Step 6:** Determine resource options
  - Pagination mode (default: infinite, options: infinite, load-more, prev-next, numbers)
  - Page size (default: 20)
  - Edit mode (default: modal, options: modal, page)
  - Ask user only if they have specific preferences, otherwise use defaults

- [ ] **Step 7:** Build and run the `lvt gen resource` command
  - Format: `lvt gen resource <name> <field1:type1> <field2:type2> ...`
  - Include optional flags if user specified preferences:
    - `--pagination <mode>` if not default
    - `--page-size <num>` if not default
    - `--edit-mode <mode>` if not default

- [ ] **Step 8:** Verify resource generation succeeded
  - Check for success message from lvt
  - Verify files created:
    - `internal/app/<resource>/<resource>.go` (handler)
    - `internal/app/<resource>/<resource>.tmpl` (template)
    - `internal/app/<resource>/<resource>_test.go` (tests)
  - Verify files updated:
    - `internal/database/schema.sql` (schema updated)
    - `internal/database/queries.sql` (queries added)
    - `internal/database/migrations/<timestamp>_create_<table>.sql` (migration created)
    - `cmd/<app>/main.go` or `main.go` (route injected)

- [ ] **Step 9:** Run database migration
  - Execute: `lvt migration up`
  - Verify migration succeeded
  - Handle errors if migration fails

- [ ] **Step 10:** Generate sqlc models
  - Navigate to `internal/database/` directory
  - Run: `go run github.com/sqlc-dev/sqlc/cmd/sqlc generate`
  - Verify models generated successfully
  - Return to project root

- [ ] **Step 11:** Run `go mod tidy`
  - Ensure all dependencies are up to date
  - Verify no errors

- [ ] **Step 12:** Verify app builds successfully
  - For multi/single kits: `go build ./cmd/<app>`
  - For simple kit: `go build`
  - If build fails, diagnose and fix issues

- [ ] **Step 13:** Provide user with success summary
  - List files created
  - List files updated
  - Show the generated route
  - Provide next steps (run app, test functionality)

## Type Inference Reference

If the user doesn't specify types, lvt will infer them:

**String types:**
- name, email, title, username, password, token
- url, slug, path, address, city, state, country, phone, status

**Text types (textarea):**
- description, content, body

**Integer types:**
- age, count, quantity, views, likes, shares, year, rating
- Any field ending with: `_count`, `_number`, `_id`, or ending with `id`

**Float types:**
- price, amount, total, latitude, longitude
- Any field ending with: `_price`, `_amount`, `_total`, or ending with `price`

**Boolean types:**
- enabled, active, visible, published, deleted, featured
- Any field starting with: `is_`, `has_`, `can_`, `should_`

**Time types:**
- created_at, updated_at, deleted_at, published_at, expires_at
- Any field ending with: `_at`, `_date`, `_time`, or ending with `date`

## Foreign Key References

To create foreign key relationships:
- Format: `field_name:references:table_name`
- With custom ON DELETE: `field_name:references:table_name:CASCADE`
- Options: CASCADE, SET NULL, RESTRICT, NO ACTION

**Example:** `user_id:references:users:CASCADE`

## Example Commands

```bash
# Simple resource with type inference
lvt gen resource posts title content published

# Explicit types
lvt gen resource products name:string price:float quantity:int enabled:bool

# With foreign key
lvt gen resource comments post_id:references:posts:CASCADE content author

# With pagination options
lvt gen resource articles title content --pagination numbers --page-size 10

# With edit mode
lvt gen resource tasks title description due_date --edit-mode page
```

## Error Handling

**If resource directory already exists:**
1. Warn user about potential data loss
2. Ask for confirmation before proceeding
3. Suggest using a different name

**If not in an lvt project:**
1. Check for `.lvtrc` file
2. If missing, inform user they need to create an app first
3. Suggest using lvt:new-app skill

**If field names are invalid:**
1. Check for Go identifier validity (alphanumeric + underscore, no spaces)
2. Suggest corrections if needed
3. Avoid SQL reserved keywords

**If migration fails:**
1. Show migration error output
2. Check for common issues (syntax, constraints, duplicates)
3. Suggest manual review of migration file

**If build fails after generation:**
1. Run `go mod tidy` again
2. Check for import errors
3. Verify sqlc generated models exist
4. Check route injection didn't break main.go syntax

## Success Response

After successful resource generation, provide:

```
✅ Resource '[name]' generated successfully!

📁 Files created:
  - internal/app/[resource]/[resource].go
  - internal/app/[resource]/[resource].tmpl
  - internal/app/[resource]/[resource]_test.go
  - internal/database/migrations/[timestamp]_create_[table].sql

📝 Files updated:
  - internal/database/schema.sql
  - internal/database/queries.sql
  - cmd/[app]/main.go (route: /[resource])

✅ Migration applied successfully
✅ Models generated successfully
✅ App builds successfully

🚀 Next steps:
  1. Start your app: lvt serve (or go run main.go)
  2. Visit: http://localhost:8080/[resource]
  3. Test CRUD operations in your browser
```

## Common User Scenarios

**Scenario 1: Blog posts**
- User: "Add a posts resource with title content published"
- Fields inferred: title→string, content→text, published→bool
- Auto-generates: posts table, CRUD handlers, list/detail templates

**Scenario 2: E-commerce products**
- User: "Create products with name, price, quantity, image URL"
- Fields inferred: name→string, price→float, quantity→int, image_url→string
- Pagination: default infinite scroll
- Edit mode: default modal

**Scenario 3: Task management with relationships**
- User: "Add tasks with title, description, user_id references users, due_date, completed"
- Recognizes foreign key: user_id→references:users
- Fields inferred: title→string, description→text, due_date→time, completed→bool
- Generates: CASCADE constraint on user_id

**Scenario 4: Comments with explicit types**
- User: "Generate CRUD for comments with post_id:int, author:string, content:text, created_at:time"
- Uses explicit types as specified
- No type inference needed

## Validation Criteria

Resource generation is successful if:
1. ✅ All files created without errors
2. ✅ Routes injected correctly in main.go
3. ✅ Migration applied successfully
4. ✅ sqlc models generated
5. ✅ `go build` succeeds
6. ✅ No compilation errors
7. ✅ App can be started and accessed

## Notes

- Resource names are automatically pluralized for table names (post → posts, user → users)
- Field names are automatically converted to snake_case for SQL, PascalCase for Go
- Auto-generated fields: id (primary key), created_at, updated_at
- Type inference is smart but explicit types are always respected
- Foreign keys automatically create indexes for performance
