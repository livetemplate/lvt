---
applyTo: "**/*.go"
---

# LiveTemplate Development Rules for Cursor

When working with LiveTemplate projects, use `lvt` CLI commands to generate and manage code. This guide provides rules and patterns optimized for Cursor AI.

## CLI Commands Overview

The `lvt` CLI provides commands across several categories:

### Generation (5)
- `lvt new` - Create new LiveTemplate app
- `lvt gen resource` - Generate CRUD resource
- `lvt gen view` - Generate view-only page
- `lvt gen auth` - Generate auth system
- `lvt gen schema` - Generate database schema

### Database (4)
- `lvt migration up` - Apply migrations
- `lvt migration down` - Rollback migration
- `lvt migration status` - Check status
- `lvt migration create` - Create migration

### Inspection (2)
- `lvt resource list` - List resources
- `lvt resource describe` - Describe schema

### Data (1)
- `lvt seed` - Generate test data

### Validation (1)
- `lvt parse` - Validate templates

### Config (3)
- `lvt env generate` - Generate .env.example
- `lvt kits list` - List kits
- `lvt kits info` - Kit details

## Development Patterns

### Pattern 1: New Feature Flow
```
User asks: "Add posts feature"
-> Use lvt gen resource with fields
-> Use lvt migration up
-> Use lvt seed for test data
-> Show user next steps
```

### Pattern 2: Database First
```
User asks: "What's in the database?"
-> Use lvt resource list
-> Use lvt resource describe for details
-> Suggest additions based on current schema
```

### Pattern 3: Auth Setup
```
User asks: "Add login"
-> Check lvt resource list for existing auth
-> If none, use lvt gen auth
-> Use lvt migration up
-> Explain auth flow to user
```

### Pattern 4: Safe Migrations
```
Before any migration:
-> Use lvt migration status first
-> Show user pending migrations
-> Get confirmation
-> Use lvt migration up
-> Verify with lvt resource describe
```

## File-Specific Rules

### When editing *.go files in app/
**DO:**
- Suggest using lvt gen resource for new resources
- Recommend lvt gen view for non-CRUD pages
- Check existing resources with lvt resource list first

**DON'T:**
- Manually write boilerplate handlers
- Create database code by hand
- Skip migrations after schema changes

### When editing migration files (*.sql)
**DO:**
- Validate syntax before applying
- Use lvt migration status to check state
- Suggest custom migrations with lvt migration create

**DON'T:**
- Apply migrations without checking status
- Edit applied migration files
- Skip sqlc regeneration (happens auto in migration up)

### When editing template files (*.tmpl)
**DO:**
- Use lvt parse before suggesting changes
- Follow LiveTemplate component patterns
- Suggest regenerating if structure is broken

**DON'T:**
- Break template syntax
- Remove LiveTemplate directives
- Ignore validation errors

## Common Workflows

### Workflow: Create Blog
```bash
lvt new blog --kit multi
cd blog
lvt gen auth
lvt gen resource posts title:string content:text user_id:references:users
lvt gen resource comments content:text post_id:references:posts user_id:references:users
lvt migration up
lvt seed posts --count 10
lvt seed comments --count 30
```

### Workflow: Add Feature to Existing App
```bash
lvt resource list
lvt gen resource tags name:string color:string
lvt migration up
lvt seed tags --count 20
```

### Workflow: Development Data Refresh
```bash
lvt seed posts --cleanup --count 50
lvt seed comments --cleanup --count 150
```

## Field Type Reference

When suggesting fields for lvt gen resource:

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
1. Check: lvt migration status
2. Review error message
3. Describe: lvt resource describe <resource>
4. Fix migration file or schema conflict
5. Retry: lvt migration up
```

### Template Errors
```
1. Validate: lvt parse <template>
2. Review syntax errors
3. Fix template
4. Re-validate
```

### Generation Errors
```
1. Check: lvt resource list (name conflicts?)
2. Verify field syntax
3. Ensure migrations applied
4. Retry generation
```

## Best Practices

1. **Always check before generating**
   - Use lvt resource list to avoid duplicates
   - Check lvt migration status before migrations

2. **Chain operations logically**
   - Generate -> Migrate -> Seed
   - List -> Describe -> Generate

3. **Use cleanup for development**
   - lvt seed with --cleanup flag
   - Keeps test data fresh

4. **Generate auth early**
   - Run lvt gen auth first
   - Then add resources with user_id references

5. **Validate before deploying**
   - Use lvt parse on all templates
   - Check lvt migration status
   - Generate .env.example with lvt env generate

## Quick Reference Commands

```bash
# New project
lvt new -> lvt gen auth -> lvt gen resource -> lvt migration up

# Add feature
lvt resource list -> lvt gen resource -> lvt migration up -> lvt seed

# Database check
lvt migration status -> lvt resource list -> lvt resource describe

# Dev refresh
lvt seed --cleanup

# Pre-deploy
lvt migration status -> lvt parse -> lvt env generate
```

## Project Structure Awareness

LiveTemplate projects have this structure:

```
project/
├── cmd/app/main.go           # Entry point
├── internal/
│   ├── app/
│   │   ├── resource/         # Each resource has:
│   │   │   ├── resource.go   # Handler (suggest lvt gen resource)
│   │   │   ├── resource.tmpl # Template (validate with lvt parse)
│   │   │   └── resource_test.go # Tests
│   │   └── auth/            # Auth system (lvt gen auth)
│   └── database/
│       ├── migrations/       # SQL files (lvt migration *)
│       ├── queries.sql       # Generated queries
│       └── schema.sql        # Complete schema
├── go.mod
└── .env.example             # Generated by lvt env generate
```

When user asks to modify any of these, suggest appropriate lvt command instead of manual editing.

## Integration with Cursor Features

### Composer Mode
When in Composer, chain multiple CLI commands together:
```
"Add blog with auth" ->
  lvt gen auth +
  lvt gen resource posts ... +
  lvt gen resource comments ... +
  lvt migration up +
  lvt seed (both resources)
```

### Agent Mode
In Agent mode, proactively:
- Check resource list before suggesting features
- Validate migrations before applying
- Seed data after adding resources
- Suggest next logical steps

### Chat Mode
In Chat, provide:
- Command explanations with examples
- Workflow guidance
- Error troubleshooting
- Best practice recommendations

## Advanced Patterns

### Multi-Resource Generation
```
For "create e-commerce store":
1. lvt new store --kit multi
2. lvt gen auth
3. lvt gen resource categories name:string
4. lvt gen resource products name:string price:float category_id:references:categories
5. lvt gen resource orders status:string user_id:references:users
6. lvt gen resource order_items quantity:int order_id:references:orders product_id:references:products
7. lvt migration up
8. Seed in order: categories -> products -> orders -> order_items
```

### Custom Migrations
```
For "add indexes":
1. lvt migration create add_indexes
2. Edit generated SQL file manually
3. lvt migration up
4. lvt resource describe to verify
```

### Template Customization
```
For "customize post template":
1. lvt parse (ensure valid first)
2. Suggest safe modifications
3. Re-validate after changes
4. Test in browser
```

## Documentation Links

- **Workflow Patterns**: See `docs/WORKFLOWS.md` for 17 detailed workflows
- **LiveTemplate Docs**: https://github.com/livetemplate/lvt

## Setup

Ensure `lvt` is installed:

```bash
# Install
go install github.com/livetemplate/lvt@latest

# Set up agent integration
lvt install-agent
```
