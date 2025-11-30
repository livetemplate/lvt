# Generic LiveTemplate AI Assistant

This directory contains LLM-agnostic documentation for integrating LiveTemplate with any AI assistant. The documentation assumes the AI can call command-line tools or use the LiveTemplate MCP server.

## Overview

LiveTemplate provides a `lvt` MCP (Model Context Protocol) server that exposes 16 tools for building full-stack Go web applications. Any AI assistant that supports:

1. **MCP protocol** (recommended) - Use `lvt mcp-server` directly
2. **Command-line execution** - Use `lvt` CLI commands via shell
3. **Tool calling** - Adapt MCP tools to your LLM's tool format

## Quick Start

### Option 1: MCP Server (Recommended)

If your AI assistant supports MCP:

```bash
# Install LiveTemplate
go install github.com/livetemplate/lvt@latest

# Start MCP server
lvt mcp-server
```

Configure your AI to connect to the MCP server. See [MCP Integration](#mcp-integration) below.

### Option 2: CLI Commands

If your AI can execute shell commands:

```bash
# Install LiveTemplate
go install github.com/livetemplate/lvt@latest

# Use CLI commands directly
lvt new myapp --kit multi
lvt gen resource posts title:string content:text
lvt migration up
```

See [CLI Reference](#cli-reference) below.

## Available Tools

The LiveTemplate MCP server exposes 16 tools across 6 categories:

### Generation (5 tools)
- `lvt_new` - Create new LiveTemplate application
- `lvt_gen_resource` - Generate CRUD resource with database
- `lvt_gen_view` - Generate view-only page (no database)
- `lvt_gen_auth` - Generate authentication system
- `lvt_gen_schema` - Generate database schema only

### Database (4 tools)
- `lvt_migration_up` - Apply pending migrations
- `lvt_migration_down` - Rollback last migration
- `lvt_migration_status` - Check migration status
- `lvt_migration_create` - Create empty migration file

### Development (7 tools)
- `lvt_seed` - Generate test data
- `lvt_resource_list` - List all resources
- `lvt_resource_describe` - Show resource schema
- `lvt_validate_template` - Validate template syntax
- `lvt_env_generate` - Generate .env.example
- `lvt_kits_list` - List available kits
- `lvt_kits_info` - Get kit details

## MCP Integration

### Server Configuration

The LiveTemplate MCP server runs as a JSON-RPC service over stdio:

```json
{
  "mcpServers": {
    "lvt": {
      "command": "lvt",
      "args": ["mcp-server"]
    }
  }
}
```

### Tool Schema Format

All tools follow standard JSON Schema format. Example for `lvt_new`:

```json
{
  "name": "lvt_new",
  "description": "Create a new LiveTemplate application",
  "inputSchema": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "description": "Application name (required)"
      },
      "kit": {
        "type": "string",
        "enum": ["multi", "single", "simple"],
        "description": "Application kit template (default: multi)"
      },
      "css": {
        "type": "string",
        "enum": ["tailwind", "bulma", "pico", "none"],
        "description": "CSS framework (default: none)"
      }
    },
    "required": ["name"]
  }
}
```

See [docs/MCP_TOOLS.md](../../docs/MCP_TOOLS.md) for complete tool documentation.

## CLI Reference

If your AI executes shell commands instead of using MCP:

### Project Setup
```bash
lvt new <name>              # Create new app
  --kit multi|single|simple # Application template
  --css tailwind|bulma|pico # CSS framework
  --module <go-module-name> # Go module path
```

### Code Generation
```bash
lvt gen resource <name> <field:type>...  # Add CRUD resource
lvt gen view <name>                      # Add view-only page
lvt gen auth                             # Add authentication
lvt gen schema <table> <field:type>...   # Add DB schema only
```

### Database Management
```bash
lvt migration up           # Apply pending migrations
lvt migration down         # Rollback last migration
lvt migration status       # Check migration state
lvt migration create <name> # Create empty migration
```

### Development Tools
```bash
lvt seed <resource> --count N --cleanup  # Generate test data
lvt resource list                        # List all resources
lvt resource describe <name>             # Show resource schema
lvt parse <template-file>                # Validate template
lvt env generate                         # Create .env.example
```

### Kits
```bash
lvt kits list              # List available kits
lvt kits info <name>       # Show kit details
```

## Field Types

When generating resources, use these field types:

| Type | SQL Type | Use For | Example |
|------|----------|---------|---------|
| `string` | TEXT | Short text | title, name, email |
| `text` | TEXT | Long text | content, description |
| `int` | INTEGER | Numbers | quantity, count |
| `bool` | BOOLEAN | Flags | published, active |
| `float` | REAL | Decimals | price, rating |
| `time` | DATETIME | Timestamps | created_at, due_date |
| `references:table` | FK | Relations | user_id, post_id |

## Common Workflows

### 1. Create New Blog

```bash
# Using MCP tools
lvt_new(name: "myblog", kit: "multi")
lvt_gen_auth()
lvt_gen_resource(name: "posts", fields: {
  title: "string",
  content: "text",
  user_id: "references:users",
  published: "bool"
})
lvt_migration_up()
lvt_seed(resource: "posts", count: 10)
```

```bash
# Using CLI
lvt new myblog --kit multi
cd myblog
lvt gen auth
lvt gen resource posts title:string content:text user_id:references:users published:bool
lvt migration up
lvt seed posts --count 10
```

### 2. Add Feature to Existing App

```bash
# Using MCP tools
lvt_resource_list()  # Check existing resources
lvt_gen_resource(name: "comments", fields: {
  content: "text",
  post_id: "references:posts",
  user_id: "references:users"
})
lvt_migration_up()
lvt_seed(resource: "comments", count: 30)
```

```bash
# Using CLI
lvt resource list
lvt gen resource comments content:text post_id:references:posts user_id:references:users
lvt migration up
lvt seed comments --count 30
```

### 3. Database Management

```bash
# Using MCP tools
lvt_migration_status()  # Check pending migrations
lvt_migration_up()      # Apply changes
lvt_resource_describe(resource: "posts")  # Verify schema
```

```bash
# Using CLI
lvt migration status
lvt migration up
lvt resource describe posts
```

## Best Practices

### 1. Always Run Migrations After Generation

```bash
lvt gen resource posts title:string
lvt migration up  # Don't forget this!
```

### 2. Check Status Before Migrations

```bash
lvt migration status  # Review pending migrations
lvt migration up      # Then apply
```

### 3. Use Cleanup When Re-seeding

```bash
lvt seed posts --count 50 --cleanup  # Removes existing data first
```

### 4. Generate Auth First

```bash
lvt gen auth  # Do this first
lvt gen resource posts user_id:references:users  # Then reference users
```

### 5. Validate Templates Before Deployment

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

## Error Handling

### Migration Issues

```bash
# Check current state
lvt migration status

# Review specific resource
lvt resource describe <name>

# Fix and retry
lvt migration up
```

### Template Issues

```bash
# Validate template
lvt parse internal/app/<resource>/<resource>.tmpl

# Fix template based on error output
# Re-validate
lvt parse internal/app/<resource>/<resource>.tmpl
```

## Adapting to Your LLM

### Tool Format Conversion

If your LLM doesn't support MCP natively, convert the tools:

**From MCP Schema:**
```json
{
  "name": "lvt_new",
  "inputSchema": {
    "type": "object",
    "properties": {
      "name": {"type": "string"}
    }
  }
}
```

**To Your LLM's Format:**
```python
# Example: OpenAI function calling
{
  "name": "lvt_new",
  "parameters": {
    "type": "object",
    "properties": {
      "name": {"type": "string"}
    }
  }
}
```

### Workflow Adaptation

The workflows in this documentation are LLM-agnostic. Adapt them to your LLM's capabilities:

- **MCP-enabled LLMs**: Use tools directly
- **CLI-capable LLMs**: Execute shell commands
- **Hybrid LLMs**: Combine both approaches

## Documentation

Complete documentation is available in the `docs/` directory:

- **[MCP_TOOLS.md](../../docs/MCP_TOOLS.md)** - Complete reference for all 16 MCP tools
- **[WORKFLOWS.md](../../docs/WORKFLOWS.md)** - 17 detailed development workflows
- **[README.md](../../README.md)** - Project overview and installation

## LLM-Specific Implementations

Pre-built integrations for popular AI assistants:

- **[Claude Code](../.claude/)** - Native Claude Code agent with skills
- **[GitHub Copilot](../.github/copilot-instructions.md)** - GitHub Copilot instructions
- **[Cursor](../.cursor/rules/lvt.md)** - Cursor AI rules
- **[Aider](../.aider/)** - Aider configuration

## Support

- **Issues**: https://github.com/livetemplate/lvt/issues
- **Documentation**: https://github.com/livetemplate/lvt
- **MCP Protocol**: https://modelcontextprotocol.io

## Installation

Install LiveTemplate globally:

```bash
go install github.com/livetemplate/lvt@latest
```

Verify installation:

```bash
lvt --help
```

## Testing MCP Integration

Test that your MCP integration works:

```bash
# Start MCP server
lvt mcp-server

# In another terminal, test with MCP inspector
npx @modelcontextprotocol/inspector lvt mcp-server
```

## Example: Complete Blog Creation

Step-by-step example showing both MCP and CLI approaches:

### Using MCP Tools

```javascript
// 1. Create app
await callTool("lvt_new", {
  name: "myblog",
  kit: "multi",
  css: "tailwind"
});

// 2. Add authentication
await callTool("lvt_gen_auth", {});

// 3. Add posts
await callTool("lvt_gen_resource", {
  name: "posts",
  fields: {
    title: "string",
    content: "text",
    user_id: "references:users",
    published: "bool",
    created_at: "time"
  }
});

// 4. Add comments
await callTool("lvt_gen_resource", {
  name: "comments",
  fields: {
    content: "text",
    post_id: "references:posts",
    user_id: "references:users"
  }
});

// 5. Apply migrations
await callTool("lvt_migration_up", {});

// 6. Seed data
await callTool("lvt_seed", {resource: "posts", count: 10});
await callTool("lvt_seed", {resource: "comments", count: 30});
```

### Using CLI

```bash
# 1. Create app
lvt new myblog --kit multi --css tailwind
cd myblog

# 2. Add authentication
lvt gen auth

# 3. Add posts
lvt gen resource posts \
  title:string \
  content:text \
  user_id:references:users \
  published:bool \
  created_at:time

# 4. Add comments
lvt gen resource comments \
  content:text \
  post_id:references:posts \
  user_id:references:users

# 5. Apply migrations
lvt migration up

# 6. Seed data
lvt seed posts --count 10
lvt seed comments --count 30

# 7. Run app
go run cmd/myblog/main.go
```

## Tips for AI Assistants

When helping developers with LiveTemplate:

1. **Always check before generating** - Use `lvt resource list` to avoid duplicates
2. **Chain commands logically** - Generate → Migrate → Seed
3. **Use cleanup for development** - `lvt seed --cleanup` keeps data fresh
4. **Generate auth early** - Do this first, then add resources with user_id references
5. **Validate before deploying** - Use `lvt parse` on all templates

## Contributing

To add support for your LLM:

1. Create a directory in `agents/<your-llm>/`
2. Add configuration files in your LLM's format
3. Reference the MCP server or CLI commands
4. Submit a pull request

See existing implementations for examples.
