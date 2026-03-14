# LiveTemplate Quick Reference for AI Assistants

## Setup

```bash
# Install
go install github.com/livetemplate/lvt@latest

# Set up AI agent integration
lvt install-agent
```

## Available CLI Commands

### Generation (5)
- `lvt new` - Create app (name, kit, css)
- `lvt gen resource` - Add CRUD (name, fields)
- `lvt gen view` - Add page (name)
- `lvt gen auth` - Add authentication
- `lvt gen schema` - Add DB schema (table, fields)

### Database (4)
- `lvt migration up` - Apply migrations
- `lvt migration down` - Rollback
- `lvt migration status` - Check status
- `lvt migration create` - Create migration (name)

### Development (7)
- `lvt seed` - Generate data (resource, count, cleanup)
- `lvt resource list` - List resources
- `lvt resource describe` - Show schema (resource)
- `lvt parse` - Validate template (template_file)
- `lvt env generate` - Generate .env.example
- `lvt kits list` - List kits
- `lvt kits info` - Kit details (name)

## Field Types

```
string  -> TEXT         (title, name, email)
text    -> TEXT         (content, description)
int     -> INTEGER      (quantity, count)
bool    -> BOOLEAN      (published, active)
float   -> REAL         (price, rating)
time    -> DATETIME     (created_at, due_date)
references:table -> FK  (user_id, post_id)
```

## Common Patterns

### New Project
```
lvt new -> lvt gen auth -> lvt gen resource -> lvt migration up -> lvt seed
```

### Add Feature
```
lvt resource list -> lvt gen resource -> lvt migration up -> lvt seed
```

### Database Check
```
lvt migration status -> lvt resource list -> lvt resource describe
```

### Pre-Deploy
```
lvt migration status -> lvt parse -> lvt env generate
```

## Example: Create Blog

```bash
lvt new blog --kit multi
cd blog
lvt gen auth
lvt gen resource posts title:string content:text user_id:references:users published:bool
lvt migration up
lvt seed posts --count 10
```

## Best Practices

1. **Check before migrate**: `lvt migration status` then `lvt migration up`
2. **Auth first**: `lvt gen auth` before resources with user_id
3. **Use cleanup**: `lvt seed --cleanup` for fresh data
4. **Validate templates**: `lvt parse` before deploy
5. **List before generate**: `lvt resource list` to avoid duplicates

## Error Handling

### Migration fails
```
lvt migration status -> Review -> lvt resource describe -> Fix -> lvt migration up
```

### Template invalid
```
lvt parse -> Fix -> Re-validate
```

## Documentation

- **Workflows**: `docs/WORKFLOWS.md`
- **Generic Guide**: `agents/generic/README.md`

## LLM-Specific Agents

- Claude Code: `.claude/`
- GitHub Copilot: `.github/copilot-instructions.md`
- Cursor: `.cursor/rules/lvt.md`
- Aider: `.aider/`
