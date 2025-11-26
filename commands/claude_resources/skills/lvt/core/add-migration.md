---
name: lvt:add-migration
description: Use when adding database migrations to LiveTemplate apps - guides both auto-generated migrations (from lvt gen resource) and custom migrations (indexes, constraints, data transformations)
---

# lvt:add-migration

Add and manage database migrations in LiveTemplate applications using goose + sqlc.

## Overview

LiveTemplate uses **goose** for migration version control and **sqlc** for type-safe code generation. Migrations come from two sources:

1. **Auto-generated**: `lvt gen resource` creates migrations automatically
2. **Custom migrations**: Manual SQL for indexes, constraints, data fixes, etc.

## Migration Workflow

### Auto-Generated Migrations (Resources)

When you generate a resource, lvt creates:
- `internal/database/schema.sql` - Updated with new table
- `internal/database/queries.sql` - CRUD operations
- Migration file - Timestamped in `migrations/`

```bash
# 1. Generate resource (creates migration automatically)
lvt gen resource products name price:float

# 2. Apply migration + generate Go code
lvt migration up

# 3. Verify
lvt migration status
```

**Note:** `lvt migration up` automatically runs `sqlc generate` - no need to run it manually.

### Custom Migrations

For indexes, constraints, or data transformations:

```bash
# 1. Create migration file
lvt migration create add_products_price_index

# 2. Edit generated file
# internal/database/migrations/YYYYMMDDHHMMSS_add_products_price_index.sql

# 3. Add SQL in goose format:
-- +goose Up
CREATE INDEX idx_products_price ON products(price);

-- +goose Down
DROP INDEX idx_products_price;

# 4. Apply migration
lvt migration up
```

## Commands Reference

| Command | Purpose | Auto-runs sqlc? |
|---------|---------|-----------------|
| `lvt migration create <name>` | Create empty migration file | No |
| `lvt migration up` | Apply pending migrations | Yes ✓ |
| `lvt migration down` | Rollback last migration | Yes ✓ |
| `lvt migration status` | Show migration state | No |

## Common Scenarios

### Adding an Index

```bash
lvt migration create add_user_email_index

# Edit migration file:
-- +goose Up
CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP INDEX idx_users_email;

lvt migration up
```

### Adding a Constraint

```bash
lvt migration create add_price_check

# Edit migration file:
-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD CONSTRAINT price_positive CHECK (price > 0);
-- +goose StatementEnd

-- +goose Down
ALTER TABLE products DROP CONSTRAINT price_positive;

lvt migration up
```

### Data Transformation

```bash
lvt migration create normalize_user_emails

# Edit migration file:
-- +goose Up
UPDATE users SET email = LOWER(TRIM(email));

-- +goose Down
-- Data transformations usually can't be reversed
-- Document this in comment

lvt migration up
```

## Migration File Format

Goose uses special comments to mark sections:

```sql
-- +goose Up
-- SQL for applying migration
CREATE TABLE users (id INTEGER PRIMARY KEY);

-- +goose Down
-- SQL for rolling back
DROP TABLE users;
```

**For multi-statement migrations:**

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (id INTEGER PRIMARY KEY);
CREATE INDEX idx_users_id ON users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_id;
DROP TABLE users;
-- +goose StatementEnd
```

## When to Use Custom Migrations

**Use custom migrations for:**
- Indexes for query performance
- Foreign key constraints
- Unique constraints
- Check constraints
- Data fixes or transformations
- Complex schema changes

**DON'T use for:**
- Adding new tables → Use `lvt gen resource`
- Adding columns to existing tables → Edit schema.sql + create migration manually

## Common Mistakes

### ❌ Editing schema.sql Without Migration

```bash
# WRONG - loses version control
vim internal/database/schema.sql  # add index manually
lvt migration up  # nothing happens
```

**Why wrong:** Schema changes need migrations for version control and team coordination.

### ❌ Running sqlc Manually

```bash
# WRONG - unnecessary step
lvt migration up
cd internal/database && sqlc generate  # already done!
```

**Why wrong:** `lvt migration up` auto-runs sqlc generate.

### ❌ Wrong Goose Format

```sql
-- WRONG - missing goose comments
CREATE INDEX idx_products_price ON products(price);
```

**Why wrong:** Goose won't recognize this as a migration.

### ❌ Forgetting Down Migration

```sql
-- +goose Up
CREATE INDEX idx_products_price ON products(price);

-- +goose Down
-- (empty)
```

**Why wrong:** Can't rollback if needed. Always write reversible migrations when possible.

## Rollback Strategy

Check status before rolling back:

```bash
# 1. See what's applied
lvt migration status

# 2. Rollback last migration
lvt migration down

# 3. Verify
lvt migration status
```

**Note:** Rolling back drops the last migration. If you need to roll back multiple migrations, run `lvt migration down` multiple times.

## Troubleshooting

### Migration Fails

```bash
# Error: migration failed
lvt migration up
# Error: near "INDEX": syntax error

# Fix:
# 1. Check SQL syntax in migration file
# 2. Fix the file
# 3. Try again - goose will resume
lvt migration up
```

### Out of Sync State

```bash
# Database has migrations that code doesn't
lvt migration status

# Options:
# 1. Rollback to code state: lvt migration down
# 2. Pull latest migrations from team
```

## Integration with Development Workflow

**Typical flow:**

```bash
# 1. Generate resource (auto-creates migration)
lvt gen resource products name price:float

# 2. Add custom migration for index
lvt migration create add_products_price_index
# (edit file with CREATE INDEX...)

# 3. Apply all pending migrations
lvt migration up

# 4. Code is generated, start developing
vim internal/app/products/products.go
```

## File Locations

```
project/
├── internal/
│   └── database/
│       ├── migrations/           ← Migration files here
│       │   ├── 20240101120000_create_products.sql
│       │   └── 20240101120100_add_products_price_index.sql
│       ├── schema.sql            ← Current schema (updated by gen resource)
│       ├── queries.sql           ← SQL queries (updated by gen resource)
│       └── models/              ← Generated Go code (auto-generated)
└── app.db                        ← SQLite database
```

## Quick Reference

**I need to...** | **Command**
---|---
Add a new table | `lvt gen resource <name> <fields>`
Add an index | `lvt migration create add_<table>_<col>_index`
Add a constraint | `lvt migration create add_<name>_constraint`
Apply migrations | `lvt migration up`
Rollback last | `lvt migration down`
Check status | `lvt migration status`
Fix data | `lvt migration create fix_<description>`

## Remember

✓ `lvt gen resource` creates migrations automatically
✓ `lvt migration up` auto-runs sqlc generate
✓ Always write Down migrations when possible
✓ Use goose format (`-- +goose Up/Down`)
✓ Custom migrations go in `internal/database/migrations/`

✗ Don't edit schema.sql directly without creating migration
✗ Don't manually run sqlc after `migration up`
✗ Don't skip Down migration section
