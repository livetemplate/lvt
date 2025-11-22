---
name: lvt:resource-inspect
description: Inspect database resources and schema - list all tables, view table structure, analyze columns and constraints
category: core
version: 1.0.0
---

# lvt:resource-inspect

Inspect existing database resources (tables) in your LiveTemplate application. View schema structure, column details, constraints, and relationships without modifying anything.

## User Prompts

**When to use:**
- "What resources exist in my app?"
- "Show me the users table structure"
- "List all my database tables"
- "What fields does the posts table have?"
- "Inspect the products resource"

**Examples:**
- "List all resources"
- "Describe the users table"
- "Show me what's in my schema"
- "What columns are in the orders table?"

## Quick Reference

```bash
# List all resources
lvt resource list
lvt resource ls  # alias

# Describe specific resource
lvt resource describe users
lvt resource desc users   # alias
lvt resource show users   # alias
```

## What It Shows

**List command (`lvt resource list`):**
- All table names
- Field count per table
- Summary of resources

**Describe command (`lvt resource describe <name>`):**
- Resource/table name
- All columns with:
  - Column name
  - Data type (SQL)
  - Constraints (NOT NULL, PRIMARY KEY, etc.)
  - Default values
  - Foreign key relationships
- Indexes
- Full table structure

## Checklist

- [ ] Verify in lvt project (`.lvtrc` exists)
- [ ] Determine if user wants list or describe
- [ ] For list: Run `lvt resource list`
- [ ] For describe: Extract table name from request
- [ ] Run: `lvt resource describe <table_name>`
- [ ] Explain output to user

## Common Issues

**Issue: "failed to parse schema"**
- schema.sql might be corrupted
- Fix: Check `internal/database/schema.sql`

**Issue: "No resources found in schema"**
- No tables created yet
- Fix: Generate resources first with `lvt gen resource`

**Issue: "resource 'X' not found"**
- Table name doesn't exist or misspelled
- Fix: Run `lvt resource list` first to see available tables

## Example Output

**List command:**
```
Available resources:
  users                (7 fields)
  posts                (6 fields)
  comments             (5 fields)
  sessions             (4 fields)

Use 'lvt resource describe <name>' to see details
```

**Describe command:**
```
Resource: users
Table: users

Fields:
  id               INTEGER      PRIMARY KEY
  email            TEXT         NOT NULL UNIQUE
  password_hash    TEXT
  email_verified   INTEGER      DEFAULT 0
  created_at       DATETIME     DEFAULT CURRENT_TIMESTAMP
  updated_at       DATETIME     DEFAULT CURRENT_TIMESTAMP

Indexes:
  idx_users_email  (email)

Foreign Keys:
  None
```

## Use Cases

1. **Schema exploration:** Understand existing database structure
2. **Documentation:** Generate schema documentation
3. **Planning:** Before adding new resources, see what exists
4. **Debugging:** Verify table structure matches expectations
5. **Migration planning:** See current state before changes

## Notes

- Read-only operation (never modifies database)
- Parses internal/database/schema.sql
- Works with all table types (resources, views, auth)
- Helpful before customizations
- No database connection needed (reads schema file)
