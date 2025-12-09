# LiveTemplate MCP Tools Reference

This document provides a complete reference for all Model Context Protocol (MCP) tools provided by the `lvt` CLI. These tools enable AI assistants to help you build LiveTemplate applications.

## Table of Contents

- [Getting Started](#getting-started)
- [Core Generation Tools](#core-generation-tools)
- [Database Migration Tools](#database-migration-tools)
- [Resource Inspection Tools](#resource-inspection-tools)
- [Data Management Tools](#data-management-tools)
- [Template Tools](#template-tools)
- [Environment Tools](#environment-tools)
- [Kits Management Tools](#kits-management-tools)
- [Tool Catalog](#tool-catalog)

## Getting Started

### Prerequisites

1. Install `lvt` CLI
2. Start the MCP server:
   ```bash
   lvt mcp-server
   ```

### For LLM Configuration

**Claude Desktop:**
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

**OpenAI Agents SDK:**
```python
from composio_openai import ComposioToolSet

toolset = ComposioToolSet()
toolset.add_mcp_server("lvt", "lvt mcp-server")
```

**LangChain:**
```python
from langchain_mcp import MCPClient

mcp = MCPClient("lvt", command=["lvt", "mcp-server"])
tools = await mcp.list_tools()
```

---

## Core Generation Tools

### 1. lvt_new

Create a new LiveTemplate application.

**Purpose:** Bootstrap a new full-stack Go web application with LiveTemplate.

**Input Schema:**
```json
{
  "name": "string (required)",
  "kit": "string (optional: multi|single|simple, default: multi)",
  "css": "string (optional: tailwind|none)",
  "module": "string (optional, default: app name)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "app_dir": "string (absolute path)"
}
```

**Example:**
```json
// Input
{
  "name": "myblog",
  "kit": "multi",
  "css": "tailwind"
}

// Output
{
  "success": true,
  "message": "App created successfully",
  "app_dir": "/path/to/myblog"
}
```

**Common Use Cases:**
- Starting a new project
- Creating prototypes
- Generating boilerplate for specific kits

---

### 2. lvt_gen_resource

Generate a full CRUD resource with database integration.

**Purpose:** Create a complete resource with handlers, templates, database schema, and migrations.

**Input Schema:**
```json
{
  "name": "string (required)",
  "fields": "object<string, string> (required)"
}
```

**Field Type Options:**
- `string` → Go: string, SQL: TEXT
- `int` → Go: int64, SQL: INTEGER
- `bool` → Go: bool, SQL: BOOLEAN
- `float` → Go: float64, SQL: REAL
- `time` → Go: time.Time, SQL: DATETIME
- `text` / `textarea` → Go: string, SQL: TEXT (multiline)
- `references:table` → Foreign key

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "files_created": "array<string>"
}
```

**Example:**
```json
// Input
{
  "name": "posts",
  "fields": {
    "title": "string",
    "content": "text",
    "published": "bool",
    "author_id": "references:users"
  }
}

// Output
{
  "success": true,
  "message": "Resource generated successfully",
  "files_created": [
    "internal/app/posts/posts.go",
    "internal/app/posts/posts.tmpl",
    "internal/database/migrations/20240101120000_create_posts.sql"
  ]
}
```

**Generated Files:**
- `internal/app/{resource}/{resource}.go` - CRUD handler
- `internal/app/{resource}/{resource}.tmpl` - UI template
- `internal/app/{resource}/{resource}_test.go` - E2E tests
- `internal/database/migrations/*.sql` - Migration file
- `internal/database/queries.sql` - SQL queries (appended)

---

### 3. lvt_gen_view

Generate a view-only handler (no database).

**Purpose:** Create pages without CRUD operations (dashboard, about, landing, etc.).

**Input Schema:**
```json
{
  "name": "string (required)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string"
}
```

**Example:**
```json
// Input
{
  "name": "dashboard"
}

// Output
{
  "success": true,
  "message": "View generated successfully"
}
```

**Use Cases:**
- Static pages (about, contact)
- Dashboards
- Landing pages
- Custom views without database

---

### 4. lvt_gen_auth

Generate a complete authentication system.

**Purpose:** Add user authentication with login, signup, password reset, sessions, etc.

**Input Schema:**
```json
{
  "struct_name": "string (optional, default: User)",
  "table_name": "string (optional, default: users)",
  "no_password": "boolean (optional)",
  "no_magic_link": "boolean (optional)",
  "no_email_confirm": "boolean (optional)",
  "no_password_reset": "boolean (optional)",
  "no_sessions_ui": "boolean (optional)",
  "no_csrf": "boolean (optional)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "features_enabled": "array<string>"
}
```

**Example:**
```json
// Input
{
  "struct_name": "Account",
  "table_name": "admin_users"
}

// Output
{
  "success": true,
  "message": "Auth system generated",
  "features_enabled": [
    "password_auth",
    "magic_link",
    "email_confirmation",
    "password_reset",
    "sessions_ui",
    "csrf_protection"
  ]
}
```

**Generated Components:**
- User struct and handlers
- Login/Signup pages
- Password hashing
- Session management
- Email verification
- Password reset flow
- CSRF protection

---

### 5. lvt_gen_schema

Generate database schema only (no handlers or UI).

**Purpose:** Create database tables without generating application code.

**Input Schema:**
```json
{
  "table": "string (required)",
  "fields": "object<string, string> (required)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string"
}
```

**Example:**
```json
// Input
{
  "table": "analytics",
  "fields": {
    "event_name": "string",
    "user_id": "int",
    "timestamp": "time"
  }
}

// Output
{
  "success": true,
  "message": "Schema generated for analytics"
}
```

**Use Cases:**
- External data sources
- Analytics tables
- Logging tables
- Data warehousing

---

## Database Migration Tools

### 6. lvt_migration_up

Run all pending database migrations.

**Purpose:** Apply schema changes to the database.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "migrations_applied": "array<string>"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": "Migrations completed",
  "migrations_applied": [
    "20240101120000_create_users.sql",
    "20240102140000_create_posts.sql"
  ]
}
```

**Side Effects:**
- Updates database schema
- Generates Go code with sqlc
- Creates type-safe query interfaces

---

### 7. lvt_migration_down

Rollback the last migration.

**Purpose:** Undo the most recent schema change.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "migration_rolled_back": "string"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": "Migration rolled back",
  "migration_rolled_back": "20240102140000_create_posts.sql"
}
```

**Warning:** Use with caution - may result in data loss.

---

### 8. lvt_migration_status

Show current migration status.

**Purpose:** Check which migrations are applied and which are pending.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "applied": "array<string>",
  "pending": "array<string>"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": "Migration status retrieved",
  "applied": [
    "20240101120000_create_users.sql"
  ],
  "pending": [
    "20240102140000_create_posts.sql"
  ]
}
```

---

### 9. lvt_migration_create

Create a new empty migration file.

**Purpose:** Generate a timestamped migration file for custom SQL.

**Input Schema:**
```json
{
  "name": "string (required)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "file_path": "string"
}
```

**Example:**
```json
// Input
{
  "name": "add_indexes"
}

// Output
{
  "success": true,
  "message": "Migration created",
  "file_path": "internal/database/migrations/20240103160000_add_indexes.sql"
}
```

---

## Resource Inspection Tools

### 10. lvt_resource_list

List all available resources in the project.

**Purpose:** Discover what resources exist in the application.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "resources": "array<string>",
  "details": "string"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": "Resources listed",
  "resources": ["users", "posts", "comments"],
  "details": "Found 3 resources:\n- users\n- posts\n- comments"
}
```

---

### 11. lvt_resource_describe

Show detailed schema for a specific resource.

**Purpose:** Inspect the database schema, fields, and types for a resource.

**Input Schema:**
```json
{
  "resource": "string (required)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "details": "string"
}
```

**Example:**
```json
// Input
{
  "resource": "posts"
}

// Output
{
  "success": true,
  "message": "Resource posts described",
  "details": "Table: posts\nFields:\n- id (TEXT PRIMARY KEY)\n- title (TEXT)\n- content (TEXT)\n- created_at (DATETIME)"
}
```

---

## Data Management Tools

### 12. lvt_seed

Generate test data for a resource.

**Purpose:** Populate database with fake data for testing and development.

**Input Schema:**
```json
{
  "resource": "string (required)",
  "count": "int (optional, default: 10)",
  "cleanup": "boolean (optional)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string"
}
```

**Example:**
```json
// Input
{
  "resource": "posts",
  "count": 50,
  "cleanup": true
}

// Output
{
  "success": true,
  "message": "Seeded 50 posts"
}
```

**Options:**
- `count` - Number of records to generate
- `cleanup` - Remove existing test data first

---

## Template Tools

### 13. lvt_validate_template

Validate and analyze a template file.

**Purpose:** Check template syntax and structure before use.

**Input Schema:**
```json
{
  "template_file": "string (required, path to .tmpl file)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "details": "string"
}
```

**Example:**
```json
// Input
{
  "template_file": "internal/app/posts/posts.tmpl"
}

// Output
{
  "success": true,
  "message": "Template is valid",
  "details": "Parsed successfully\nComponents: 3\nActions: 5"
}
```

---

## Environment Tools

### 14. lvt_env_generate

Generate .env.example with detected configuration.

**Purpose:** Create environment variable template based on code analysis.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": ".env.example generated"
}
```

**Generated Variables:**
- Database connection strings
- Email service configuration
- API keys placeholders
- Port settings
- Environment (dev/prod)

---

## Kits Management Tools

### 15. lvt_kits_list

List all available CSS framework kits.

**Purpose:** Discover available kits for new projects.

**Input Schema:**
```json
{}  // No input required
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "details": "string"
}
```

**Example:**
```json
// Output
{
  "success": true,
  "message": "Kits listed",
  "details": "Available kits:\n- multi (Tailwind CSS)\n- single (Tailwind CSS)\n- simple (Pico CSS)"
}
```

---

### 16. lvt_kits_info

Show detailed information about a specific kit.

**Purpose:** Get kit details before creating a new app.

**Input Schema:**
```json
{
  "name": "string (required)"
}
```

**Output Schema:**
```json
{
  "success": "boolean",
  "message": "string",
  "details": "string"
}
```

**Example:**
```json
// Input
{
  "name": "multi"
}

// Output
{
  "success": true,
  "message": "Kit info retrieved",
  "details": "Kit: multi\nCSS: Tailwind CSS\nMode: Multi-page\nComponents: Full HTML layout"
}
```

---

## Tool Catalog

Quick reference table of all 16 MCP tools:

| Tool | Category | Purpose | Requires Project |
|------|----------|---------|-----------------|
| lvt_new | Generation | Create new app | No |
| lvt_gen_resource | Generation | Add CRUD resource | Yes |
| lvt_gen_view | Generation | Add view-only page | Yes |
| lvt_gen_auth | Generation | Add authentication | Yes |
| lvt_gen_schema | Generation | Add database schema | Yes |
| lvt_migration_up | Migration | Run migrations | Yes |
| lvt_migration_down | Migration | Rollback migration | Yes |
| lvt_migration_status | Migration | Check migration status | Yes |
| lvt_migration_create | Migration | Create migration file | Yes |
| lvt_seed | Data | Generate test data | Yes |
| lvt_resource_list | Inspection | List resources | Yes |
| lvt_resource_describe | Inspection | Describe resource | Yes |
| lvt_validate_template | Validation | Validate template | Yes |
| lvt_env_generate | Config | Generate .env.example | Yes |
| lvt_kits_list | Kits | List available kits | No |
| lvt_kits_info | Kits | Get kit details | No |

---

## Error Handling

All tools return structured errors:

```json
{
  "success": false,
  "message": "Error description",
  "details": "Additional error context (optional)"
}
```

Common error scenarios:
- Missing required fields → Field validation error
- File not found → Path error with suggestions
- Command execution failure → Command output included
- Invalid input → Specific validation message

---

## Best Practices

1. **Always run migrations after generation:**
   ```
   lvt_gen_resource → lvt_migration_up
   ```

2. **Check status before migrations:**
   ```
   lvt_migration_status → Review → lvt_migration_up
   ```

3. **Use cleanup when re-seeding:**
   ```
   lvt_seed with cleanup: true
   ```

4. **Validate templates before deployment:**
   ```
   lvt_validate_template → Fix issues → Deploy
   ```

5. **List resources before describing:**
   ```
   lvt_resource_list → Pick resource → lvt_resource_describe
   ```

---

## Complete Workflow Example

Building a blog application:

```json
// 1. Create app
{"tool": "lvt_new", "input": {"name": "myblog", "kit": "multi"}}

// 2. Add users (auth)
{"tool": "lvt_gen_auth", "input": {}}

// 3. Add posts resource
{"tool": "lvt_gen_resource", "input": {
  "name": "posts",
  "fields": {
    "title": "string",
    "content": "text",
    "user_id": "references:users"
  }
}}

// 4. Add comments resource
{"tool": "lvt_gen_resource", "input": {
  "name": "comments",
  "fields": {
    "content": "text",
    "post_id": "references:posts",
    "user_id": "references:users"
  }
}}

// 5. Check migration status
{"tool": "lvt_migration_status", "input": {}}

// 6. Apply migrations
{"tool": "lvt_migration_up", "input": {}}

// 7. Seed test data
{"tool": "lvt_seed", "input": {"resource": "posts", "count": 20}}
{"tool": "lvt_seed", "input": {"resource": "comments", "count": 50}}

// 8. Verify resources
{"tool": "lvt_resource_list", "input": {}}

// 9. Generate environment template
{"tool": "lvt_env_generate", "input": {}}
```

---

## Support

For issues or questions:
- Documentation: https://github.com/livetemplate/lvt
- GitHub Issues: https://github.com/livetemplate/lvt/issues
- MCP Server: `lvt mcp-server --help`
