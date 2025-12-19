---
applyTo: "**/*.go"
---

# LiveTemplate Development Rules for Cursor

When working with LiveTemplate projects, use the `lvt` MCP server tools to generate and manage code. This guide provides rules and patterns optimized for Cursor AI.

## MCP Tools Overview

The LiveTemplate MCP server exposes 16 tools across 6 categories:

### Generation (5)
- `lvt_new` - Create new LiveTemplate app
- `lvt_gen_resource` - Generate CRUD resource
- `lvt_gen_view` - Generate view-only page
- `lvt_gen_auth` - Generate auth system
- `lvt_gen_schema` - Generate database schema

### Database (4)
- `lvt_migration_up` - Apply migrations
- `lvt_migration_down` - Rollback migration
- `lvt_migration_status` - Check status
- `lvt_migration_create` - Create migration

### Inspection (2)
- `lvt_resource_list` - List resources
- `lvt_resource_describe` - Describe schema

### Data (1)
- `lvt_seed` - Generate test data

### Validation (1)
- `lvt_validate_template` - Validate templates

### Config (3)
- `lvt_env_generate` - Generate .env.example
- `lvt_kits_list` - List kits
- `lvt_kits_info` - Kit details

## Development Patterns

### Pattern 1: New Feature Flow
```
User asks: "Add posts feature"
→ Use lvt_gen_resource with fields
→ Use lvt_migration_up
→ Use lvt_seed for test data
→ Show user next steps
```

### Pattern 2: Database First
```
User asks: "What's in the database?"
→ Use lvt_resource_list
→ Use lvt_resource_describe for details
→ Suggest additions based on current schema
```

### Pattern 3: Auth Setup
```
User asks: "Add login"
→ Check lvt_resource_list for existing auth
→ If none, use lvt_gen_auth
→ Use lvt_migration_up
→ Explain auth flow to user
```

### Pattern 4: Safe Migrations
```
Before any migration:
→ Use lvt_migration_status first
→ Show user pending migrations
→ Get confirmation
→ Use lvt_migration_up
→ Verify with lvt_resource_describe
```

## File-Specific Rules

### When editing *.go files in app/
**DO:**
- Suggest using lvt_gen_resource for new resources
- Recommend lvt_gen_view for non-CRUD pages
- Check existing resources with lvt_resource_list first

**DON'T:**
- Manually write boilerplate handlers
- Create database code by hand
- Skip migrations after schema changes

### When editing migration files (*.sql)
**DO:**
- Validate syntax before applying
- Use lvt_migration_status to check state
- Suggest custom migrations with lvt_migration_create

**DON'T:**
- Apply migrations without checking status
- Edit applied migration files
- Skip sqlc regeneration (happens auto in migration_up)

### When editing template files (*.tmpl)
**DO:**
- Use lvt_validate_template before suggesting changes
- Follow LiveTemplate component patterns
- Suggest regenerating if structure is broken

**DON'T:**
- Break template syntax
- Remove LiveTemplate directives
- Ignore validation errors

## Common Workflows

### Workflow: Create Blog
```json
[
  {"tool": "lvt_new", "input": {"name": "blog", "kit": "multi"}},
  {"tool": "lvt_gen_auth", "input": {}},
  {"tool": "lvt_gen_resource", "input": {
    "name": "posts",
    "fields": {"title": "string", "content": "text", "user_id": "references:users"}
  }},
  {"tool": "lvt_gen_resource", "input": {
    "name": "comments",
    "fields": {"content": "text", "post_id": "references:posts", "user_id": "references:users"}
  }},
  {"tool": "lvt_migration_up", "input": {}},
  {"tool": "lvt_seed", "input": {"resource": "posts", "count": 10}},
  {"tool": "lvt_seed", "input": {"resource": "comments", "count": 30}}
]
```

### Workflow: Add Feature to Existing App
```json
[
  {"tool": "lvt_resource_list", "input": {}},
  {"tool": "lvt_gen_resource", "input": {"name": "tags", "fields": {...}}},
  {"tool": "lvt_migration_up", "input": {}},
  {"tool": "lvt_seed", "input": {"resource": "tags", "count": 20}}
]
```

### Workflow: Development Data Refresh
```json
[
  {"tool": "lvt_seed", "input": {"resource": "posts", "cleanup": true, "count": 50}},
  {"tool": "lvt_seed", "input": {"resource": "comments", "cleanup": true, "count": 150}}
]
```

## Field Type Reference

When suggesting fields for lvt_gen_resource:

| Go Type | SQL Type | Use For | Example Field |
|---------|----------|---------|---------------|
| string | TEXT | Short text | title, name, email |
| text | TEXT | Long text | content, description, bio |
| int | INTEGER | Numbers | quantity, priority, count |
| bool | BOOLEAN | Flags | published, completed, active |
| float | REAL | Decimals | price, rating, score |
| time | DATETIME | Timestamps | created_at, due_date |
| references:table | FK | Relations | user_id, post_id |

## Error Handling

### Migration Errors
```
1. Check: lvt_migration_status
2. Review error message
3. Describe: lvt_resource_describe <resource>
4. Fix migration file or schema conflict
5. Retry: lvt_migration_up
```

### Template Errors
```
1. Validate: lvt_validate_template
2. Review syntax errors
3. Fix template
4. Re-validate
```

### Generation Errors
```
1. Check: lvt_resource_list (name conflicts?)
2. Verify field syntax
3. Ensure migrations applied
4. Retry generation
```

## Best Practices

1. **Always check before generating**
   - Use lvt_resource_list to avoid duplicates
   - Check lvt_migration_status before migrations

2. **Chain operations logically**
   - Generate → Migrate → Seed
   - List → Describe → Generate

3. **Use cleanup for development**
   - lvt_seed with cleanup: true
   - Keeps test data fresh

4. **Generate auth early**
   - Run lvt_gen_auth first
   - Then add resources with user_id references

5. **Validate before deploying**
   - Use lvt_validate_template on all templates
   - Check lvt_migration_status
   - Generate .env.example with lvt_env_generate

## Quick Reference Commands

```bash
# New project
lvt_new → lvt_gen_auth → lvt_gen_resource → lvt_migration_up

# Add feature
lvt_resource_list → lvt_gen_resource → lvt_migration_up → lvt_seed

# Database check
lvt_migration_status → lvt_resource_list → lvt_resource_describe

# Dev refresh
lvt_seed (cleanup: true)

# Pre-deploy
lvt_migration_status → lvt_validate_template → lvt_env_generate
```

## Project Structure Awareness

LiveTemplate projects have this structure:

```
project/
├── cmd/app/main.go           # Entry point
├── internal/
│   ├── app/
│   │   ├── resource/         # Each resource has:
│   │   │   ├── resource.go   # Handler (suggest lvt_gen_resource)
│   │   │   ├── resource.tmpl # Template (validate with lvt_validate_template)
│   │   │   └── resource_test.go # Tests
│   │   └── auth/            # Auth system (lvt_gen_auth)
│   └── database/
│       ├── migrations/       # SQL files (lvt_migration_*)
│       ├── queries.sql       # Generated queries
│       └── schema.sql        # Complete schema
├── go.mod
└── .env.example             # Generated by lvt_env_generate
```

When user asks to modify any of these, suggest appropriate lvt tool instead of manual editing.

## Integration with Cursor Features

### Composer Mode
When in Composer, chain multiple MCP tools together:
```
"Add blog with auth" →
  lvt_gen_auth +
  lvt_gen_resource(posts) +
  lvt_gen_resource(comments) +
  lvt_migration_up +
  lvt_seed (both resources)
```

### Agent Mode
In Agent mode, proactively:
- Check resource list before suggesting features
- Validate migrations before applying
- Seed data after adding resources
- Suggest next logical steps

### Chat Mode
In Chat, provide:
- Tool explanations with examples
- Workflow guidance
- Error troubleshooting
- Best practice recommendations

## Advanced Patterns

### Multi-Resource Generation
```
For "create e-commerce store":
1. lvt_new (store, multi)
2. lvt_gen_auth
3. lvt_gen_resource (categories)
4. lvt_gen_resource (products with category_id)
5. lvt_gen_resource (orders with user_id)
6. lvt_gen_resource (order_items with order_id, product_id)
7. lvt_migration_up
8. Seed in order: categories → products → (users from auth) → orders → order_items
```

### Custom Migrations
```
For "add indexes":
1. lvt_migration_create (name: "add_indexes")
2. Edit generated SQL file manually
3. lvt_migration_up
4. lvt_resource_describe to verify
```

### Template Customization
```
For "customize post template":
1. lvt_validate_template (ensure valid first)
2. Suggest safe modifications
3. Re-validate after changes
4. Test in browser
```

## Documentation Links

- **Full Tool Reference**: See `docs/MCP_TOOLS.md` for complete tool documentation
- **Workflow Patterns**: See `docs/WORKFLOWS.md` for 17 detailed workflows
- **LiveTemplate Docs**: https://github.com/livetemplate/lvt

## MCP Server Setup

Ensure `lvt` MCP server is running for these tools to work:

```bash
# Install
go install github.com/livetemplate/lvt@latest

# Start server
lvt mcp-server
```

Configure in Cursor settings or use Auto-MCP if available.
