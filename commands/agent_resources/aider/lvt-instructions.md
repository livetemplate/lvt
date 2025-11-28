# LiveTemplate Instructions for Aider

You are assisting with LiveTemplate development. Use `lvt` CLI commands to generate and manage code.

## Available Commands

LiveTemplate provides these CLI commands for code generation:

### Project Setup
```bash
lvt new <name>           # Create new app
  --kit multi|single|simple
  --css tailwind|bulma|pico|none
  --module <go-module-name>
```

### Code Generation
```bash
lvt gen resource <name> <field:type>...  # Add CRUD resource
lvt gen view <name>                      # Add view-only page
lvt gen auth                             # Add authentication
lvt gen schema <table> <field:type>...   # Add DB schema only
```

### Database
```bash
lvt migration up         # Apply all pending migrations
lvt migration down       # Rollback last migration
lvt migration status     # Check migration state
lvt migration create <name>  # Create empty migration
```

### Development
```bash
lvt seed <resource> --count N --cleanup  # Generate test data
lvt resource list                        # List all resources
lvt resource describe <name>             # Show resource schema
lvt parse <template-file>                # Validate template
lvt env generate                         # Create .env.example
```

### Kits
```bash
lvt kits list            # List available kits
lvt kits info <name>     # Show kit details
```

## Common Workflows

### 1. Create New App
```bash
# Start new blog project
lvt new myblog --kit multi --css tailwind

# Navigate to directory
cd myblog

# Add authentication
lvt gen auth

# Add resources
lvt gen resource posts title:string content:text user_id:references:users published:bool
lvt gen resource comments content:text post_id:references:posts user_id:references:users

# Apply database changes
lvt migration up

# Add test data
lvt seed posts --count 10
lvt seed comments --count 30

# Generate environment template
lvt env generate
```

### 2. Add Feature to Existing App
```bash
# Check what exists
lvt resource list

# Add new resource
lvt gen resource tasks title:string description:text completed:bool user_id:references:users due_date:time

# Apply changes
lvt migration up

# Add test data
lvt seed tasks --count 20
```

### 3. Database Management
```bash
# Check migration status
lvt migration status

# Apply pending migrations
lvt migration up

# Verify schema
lvt resource describe tasks
```

## Field Types

When generating resources, use these field types:

- `string` → TEXT (short text)
- `text` → TEXT (long text, multiline)
- `int` → INTEGER
- `bool` → BOOLEAN
- `float` → REAL
- `time` → DATETIME
- `references:table` → Foreign key to table

Examples:
```bash
lvt gen resource products \
  name:string \
  description:text \
  price:float \
  stock:int \
  available:bool \
  category_id:references:categories
```

## Best Practices

1. **Always run migrations after generation**
   ```bash
   lvt gen resource <name> <fields...>
   lvt migration up
   ```

2. **Check status before migrations**
   ```bash
   lvt migration status
   lvt migration up
   ```

3. **Use cleanup when re-seeding**
   ```bash
   lvt seed <resource> --count N --cleanup
   ```

4. **Generate auth first**
   ```bash
   lvt gen auth
   lvt gen resource items name:string user_id:references:users
   ```

5. **Validate templates before deploying**
   ```bash
   lvt parse internal/app/posts/posts.tmpl
   ```

## Project Structure

LiveTemplate generates this structure:

```
app/
├── cmd/app/main.go           # Entry point
├── internal/
│   ├── app/
│   │   ├── auth/            # Auth system (if generated)
│   │   ├── posts/           # Each resource:
│   │   │   ├── posts.go     #   Handler
│   │   │   ├── posts.tmpl   #   Template
│   │   │   └── posts_test.go #   Tests
│   │   └── ...
│   └── database/
│       ├── migrations/       # SQL migrations
│       ├── queries.sql       # SQL queries
│       └── schema.sql        # Complete schema
├── go.mod
└── .env.example             # Environment template
```

## Recommendations for Aider

When user asks to:

- **"Add a feature"** → Use `lvt gen resource`
- **"Add login"** → Use `lvt gen auth`
- **"Add a page"** → Use `lvt gen view`
- **"Update database"** → Use `lvt migration up`
- **"Add test data"** → Use `lvt seed`
- **"What's in the database?"** → Use `lvt resource list`

Always prefer using `lvt` commands over manual file creation.

## Error Handling

### Migration Issues
```bash
# Check current state
lvt migration status

# Check specific resource
lvt resource describe <name>

# Review migration files in internal/database/migrations/

# Fix and retry
lvt migration up
```

### Template Issues
```bash
# Validate template
lvt parse internal/app/<resource>/<resource>.tmpl

# Review error output
# Fix template
# Re-validate
```

## Example Session

```bash
# User: "Create a blog with authentication"

# 1. Create app
lvt new blog --kit multi

# 2. Enter directory
cd blog

# 3. Add authentication
lvt gen auth

# 4. Add posts
lvt gen resource posts \
  title:string \
  content:text \
  user_id:references:users \
  published:bool \
  created_at:time

# 5. Add comments
lvt gen resource comments \
  content:text \
  post_id:references:posts \
  user_id:references:users

# 6. Apply all changes
lvt migration up

# 7. Generate test data
lvt seed users --count 5
lvt seed posts --count 20
lvt seed comments --count 50

# 8. Generate environment template
lvt env generate

# 9. Ready to run
go run cmd/blog/main.go
```

## Advanced Usage

### Custom Migrations
```bash
# Create empty migration
lvt migration create add_indexes

# Edit the generated file in internal/database/migrations/

# Apply
lvt migration up
```

### Refresh Development Data
```bash
# Clean and reseed
lvt seed posts --count 50 --cleanup
lvt seed comments --count 150 --cleanup
```

### Pre-Deployment Checklist
```bash
# 1. Check migrations
lvt migration status

# 2. Validate templates
lvt parse internal/app/posts/posts.tmpl
lvt parse internal/app/comments/comments.tmpl

# 3. Generate env template
lvt env generate

# 4. Review .env.example
cat .env.example
```

## Documentation

- Full CLI Reference: `lvt --help`
- Tool Documentation: `docs/MCP_TOOLS.md`
- Workflow Patterns: `docs/WORKFLOWS.md`
- LiveTemplate: https://github.com/livetemplate/lvt

## Tips for AI-Assisted Development

1. **Suggest commands, don't execute directly**
   - Show the `lvt` command to run
   - Explain what it will do
   - Let user execute it

2. **Chain commands logically**
   - Generate → Migrate → Seed
   - List → Describe → Generate

3. **Always check before modifying**
   - Use `lvt resource list` before adding
   - Use `lvt migration status` before migrating

4. **Focus on workflows**
   - Guide user through complete workflows
   - Don't just give single commands
   - Explain the "why" behind each step

5. **Validate assumptions**
   - Check what exists with list/describe
   - Don't assume structure
   - Verify before suggesting changes
