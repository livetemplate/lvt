---
name: lvt:gen-schema
description: Generate database schema without UI - creates migration, schema SQL, and sqlc queries for custom tables
category: core
version: 1.0.0
---

# lvt:gen-schema

Generates database schema (migration + SQL + queries) without creating handlers or templates. Perfect for:
- Data-only tables (logs, analytics, sessions)
- Backend tables without UI
- Custom database structures
- Tables used by multiple resources

## User Prompts

**When to use:**
- "Create a database table without UI"
- "Add a logs table to the database"
- "Generate schema for analytics data"
- "I need a table but no CRUD interface"

**Examples:**
- "Create an audit_logs table with user_id, action, timestamp"
- "Add a sessions table to store user sessions"
- "Generate schema for notifications table"

## Quick Reference

```bash
# Basic usage (with type inference)
lvt gen schema table_name field1 field2 field3

# Explicit types
lvt gen schema products name:string price:float quantity:int

# With foreign keys
lvt gen schema orders user_id:references:users:CASCADE total:float

# Complex example
lvt gen schema audit_logs user_id:references:users action:string details:text created_at:time
```

## What It Generates

**Files created:**
- `internal/database/migrations/<timestamp>_create_<table>.sql` - Migration file
- Updates `internal/database/schema.sql` - Schema definition
- Updates `internal/database/queries.sql` - CRUD queries for sqlc

**Does NOT create:**
- Handler files (*.go)
- Template files (*.tmpl)
- Test files (*_test.go)
- Routes (no auto-injection)

## Checklist

- [ ] Verify in lvt project (`.lvtrc` exists)
- [ ] Parse table name and fields from user request
- [ ] Apply type inference if needed
- [ ] Run: `lvt gen schema <table> <fields...>`
- [ ] Verify migration created
- [ ] Run: `lvt migration up`
- [ ] Run: `cd internal/database && sqlc generate && cd ../..`
- [ ] Run: `go mod tidy`
- [ ] Verify build succeeds

## Common Issues

**Issue: "table name required"**
- Missing table name argument
- Fix: `lvt gen schema table_name field1 field2`

**Issue: "at least one field required"**
- No fields specified
- Fix: Provide at least one field

**Issue: Build fails after generation**
- Forgot to run sqlc generate
- Fix: `cd internal/database && sqlc generate`

## Success Response

```
✅ Schema generated successfully!

Files created/updated:
  - internal/database/migrations/<timestamp>_create_<table>.sql
  - internal/database/schema.sql (updated)
  - internal/database/queries.sql (updated)

Next steps:
  1. Run migration: lvt migration up
  2. Generate models: cd internal/database && sqlc generate
  3. Use generated types in your handlers
```

## Use Cases

1. **Audit Logs:** Track user actions without UI
2. **Sessions:** Custom session storage
3. **Analytics:** Data collection tables
4. **Cache:** Temporary data storage
5. **Queue:** Job/task tables
6. **Settings:** App configuration storage

## Notes

- Uses same type inference as `lvt gen resource`
- Auto-generates id, created_at, updated_at
- Generates standard CRUD queries for sqlc
- Perfect for backend-only tables
- Can be used with custom handlers later
