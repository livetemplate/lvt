# Generic LiveTemplate AI Assistant

This directory contains LLM-agnostic documentation for integrating LiveTemplate with any AI assistant. The documentation assumes the AI can call command-line tools.

## Overview

LiveTemplate provides the `lvt` CLI for building full-stack Go web applications. Any AI assistant that supports command-line execution can use `lvt` commands directly.

## Quick Start

### Installation

```bash
# Install LiveTemplate
go install github.com/livetemplate/lvt@latest

# Set up AI agent integration for your project
lvt install-agent
```

### Using CLI Commands

```bash
# Use CLI commands directly
lvt new myapp --kit multi
lvt gen resource posts title:string content:text
lvt migration up
```

See [CLI Reference](#cli-reference) below.

## Available Commands

The `lvt` CLI provides commands across several categories:

### Generation (5 commands)
- `lvt new` - Create new LiveTemplate application
- `lvt gen resource` - Generate CRUD resource with database
- `lvt gen view` - Generate view-only page (no database)
- `lvt gen auth` - Generate authentication system
- `lvt gen schema` - Generate database schema only

### Database (4 commands)
- `lvt migration up` - Apply pending migrations
- `lvt migration down` - Rollback last migration
- `lvt migration status` - Check migration status
- `lvt migration create` - Create empty migration file

### Development (7 commands)
- `lvt seed` - Generate test data
- `lvt resource list` - List all resources
- `lvt resource describe` - Show resource schema
- `lvt parse` - Validate template syntax
- `lvt env generate` - Generate .env.example
- `lvt kits list` - List available kits
- `lvt kits info` - Get kit details

## CLI Reference

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
lvt new myblog --kit multi
cd myblog
lvt gen auth
lvt gen resource posts title:string content:text user_id:references:users published:bool
lvt migration up
lvt seed posts --count 10
```

### 2. Add Feature to Existing App

```bash
lvt resource list
lvt gen resource comments content:text post_id:references:posts user_id:references:users
lvt migration up
lvt seed comments --count 30
```

### 3. Database Management

```bash
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
lvt parse app/posts/posts.tmpl
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
lvt parse app/<resource>/<resource>.tmpl

# Fix template based on error output
# Re-validate
lvt parse app/<resource>/<resource>.tmpl
```

## Documentation

Complete documentation is available in the `docs/` directory:

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

## Installation

Install LiveTemplate globally:

```bash
go install github.com/livetemplate/lvt@latest
```

Verify installation:

```bash
lvt --help
```

Set up AI agent integration:

```bash
lvt install-agent
```

## Example: Complete Blog Creation

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
2. **Chain commands logically** - Generate -> Migrate -> Seed
3. **Use cleanup for development** - `lvt seed --cleanup` keeps data fresh
4. **Generate auth early** - Do this first, then add resources with user_id references
5. **Validate before deploying** - Use `lvt parse` on all templates

## Contributing

To add support for your LLM:

1. Create a directory in `agents/<your-llm>/`
2. Add configuration files in your LLM's format
3. Reference the `lvt` CLI commands
4. Submit a pull request

See existing implementations for examples.
