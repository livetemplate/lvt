# LiveTemplate Quick Reference for AI Assistants

## MCP Server Setup

```bash
# Install
go install github.com/livetemplate/lvt@latest

# Start MCP server
lvt mcp-server
```

## 16 Available Tools

### Generation (5)
- `lvt_new` - Create app (name, kit, css)
- `lvt_gen_resource` - Add CRUD (name, fields)
- `lvt_gen_view` - Add page (name)
- `lvt_gen_auth` - Add authentication
- `lvt_gen_schema` - Add DB schema (table, fields)

### Database (4)
- `lvt_migration_up` - Apply migrations
- `lvt_migration_down` - Rollback
- `lvt_migration_status` - Check status
- `lvt_migration_create` - Create migration (name)

### Development (7)
- `lvt_seed` - Generate data (resource, count, cleanup)
- `lvt_resource_list` - List resources
- `lvt_resource_describe` - Show schema (resource)
- `lvt_validate_template` - Validate (template_file)
- `lvt_env_generate` - Generate .env.example
- `lvt_kits_list` - List kits
- `lvt_kits_info` - Kit details (name)

## Field Types

```
string  → TEXT         (title, name, email)
text    → TEXT         (content, description)
int     → INTEGER      (quantity, count)
bool    → BOOLEAN      (published, active)
float   → REAL         (price, rating)
time    → DATETIME     (created_at, due_date)
references:table → FK  (user_id, post_id)
```

## Common Patterns

### New Project
```
lvt_new → lvt_gen_auth → lvt_gen_resource → lvt_migration_up → lvt_seed
```

### Add Feature
```
lvt_resource_list → lvt_gen_resource → lvt_migration_up → lvt_seed
```

### Database Check
```
lvt_migration_status → lvt_resource_list → lvt_resource_describe
```

### Pre-Deploy
```
lvt_migration_status → lvt_validate_template → lvt_env_generate
```

## Example: Create Blog

```json
// MCP calls
lvt_new({name: "blog", kit: "multi"})
lvt_gen_auth({})
lvt_gen_resource({
  name: "posts",
  fields: {
    title: "string",
    content: "text",
    user_id: "references:users",
    published: "bool"
  }
})
lvt_migration_up({})
lvt_seed({resource: "posts", count: 10})
```

```bash
# CLI equivalent
lvt new blog --kit multi
lvt gen auth
lvt gen resource posts title:string content:text user_id:references:users published:bool
lvt migration up
lvt seed posts --count 10
```

## Best Practices

1. **Check before migrate**: `lvt_migration_status` → `lvt_migration_up`
2. **Auth first**: `lvt_gen_auth` before resources with user_id
3. **Use cleanup**: `lvt_seed` with `cleanup: true` for fresh data
4. **Validate templates**: `lvt_validate_template` before deploy
5. **List before generate**: `lvt_resource_list` to avoid duplicates

## Error Handling

### Migration fails
```
lvt_migration_status → Review → lvt_resource_describe → Fix → lvt_migration_up
```

### Template invalid
```
lvt_validate_template → Fix → Re-validate
```

## Documentation

- **Complete Tool Docs**: `docs/MCP_TOOLS.md`
- **Workflows**: `docs/WORKFLOWS.md`
- **Generic Guide**: `agents/generic/README.md`

## LLM-Specific Agents

- Claude Code: `.claude/`
- GitHub Copilot: `.github/copilot-instructions.md`
- Cursor: `.cursor/rules/lvt.md`
- Aider: `.aider/`
