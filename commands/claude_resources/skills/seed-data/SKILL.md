---
name: lvt-seed-data
description: Use when generating test data for LiveTemplate apps - covers seeding resources with realistic fake data, cleanup, and understanding data generation patterns
---

# lvt:seed-data

Generate realistic test data for development and testing.

## 🎯 ACTIVATION RULES

### Context Detection

This skill typically runs in **existing LiveTemplate projects** (.lvtrc exists).

**✅ Context Established By:**
1. **Project context** - `.lvtrc` exists (most common scenario)
2. **Agent context** - User is working with `lvt-assistant` agent
3. **Keyword context** - User mentions "lvt", "livetemplate", or "lt"

**Keyword matching** (case-insensitive): `lvt`, `livetemplate`, `lt`

### Trigger Patterns

**With Context:**
✅ Generic prompts related to this skill's purpose

**Without Context (needs keywords):**
✅ Must mention "lvt", "livetemplate", or "lt"
❌ Generic requests without keywords

---

## Overview

LiveTemplate includes a seeder that generates realistic test data based on your schema. It uses field names to generate contextually appropriate values (e.g., "email" → fake email addresses).

**Key features:**
- Context-aware generation (field names determine data type)
- Bulk insert with progress tracking
- Cleanup of test data
- Marks test records for easy identification

## Basic Usage

```bash
# Generate 50 products
lvt seed products --count 50

# Clean up test data
lvt seed products --cleanup

# Clean up then seed fresh data
lvt seed products --cleanup --count 30
```

## Prerequisites

**Before seeding:**
1. ✓ Resource generated (`lvt gen resource`)
2. ✓ Migrations applied (`lvt migration up`)
3. ✓ Database exists (`app.db`)

```bash
# Complete setup before seeding
lvt gen resource products name price:float
lvt migration up
lvt seed products --count 50
```

## Commands

### Seed Data

```bash
# Basic seeding
lvt seed <resource-name> --count N

# Examples
lvt seed users --count 100
lvt seed products --count 50
lvt seed tasks --count 200
```

**Progress tracking:**
```
Seeding products with 50 rows...
  Progress: 10/50
  Progress: 20/50
  Progress: 30/50
  Progress: 40/50
  Progress: 50/50
✅ Successfully seeded 50 rows

Total test records in products: 50
```

### Clean Up Test Data

```bash
# Remove all test records
lvt seed <resource-name> --cleanup

# Example
lvt seed products --cleanup
```

**Output:**
```
Cleaning up test data for products...
✅ Deleted 50 test records from products
```

### Clean + Reseed

```bash
# Clean old data then seed fresh
lvt seed <resource-name> --cleanup --count N

# Example - replace old test data with 100 new records
lvt seed users --cleanup --count 100
```

## Context-Aware Generation

The seeder generates realistic data based on field names:

| Field Name | Generated Data |
|------------|----------------|
| `email` | john.doe@example.com |
| `name` | John Doe |
| `first_name` | John |
| `last_name` | Doe |
| `phone` | (555) 123-4567 |
| `address` | 123 Main St |
| `city` | Springfield |
| `state` | California |
| `country` | United States |
| `title` | Senior Developer |
| `company` | Tech Corp |
| `content` | Full paragraphs |
| `description` | 2-3 sentences |
| `price` | 10.00 - 10000.00 |
| `quantity` | 1 - 1000 |
| `age` | 18 - 99 |
| `rating` | 1.0 - 5.0 |
| `status` | active/inactive/pending |
| `url` | https://example.com |
| `image` | https://example.com/image.jpg |
| `color` | Red, Blue, etc. |
| `date` | 2024-01-15 |
| `enabled` / `active` / `is_*` | true/false (random) |

**Auto-managed fields:**
- `id` - Auto-generated (marked with test prefix)
- `created_at` - Current timestamp
- `updated_at` - Current timestamp

**Fallback for unknown fields:**
- String fields → Random text (5-10 words)
- Numeric fields → Random number (1-1000)
- Bool fields → Random true/false

## Examples

### User Data

```bash
# Generate 100 users
lvt seed users --count 100
```

**Generated fields:**
- `name` → "Alice Johnson"
- `email` → "alice.johnson@example.com"
- `phone` → "(555) 234-5678"
- `age` → 34
- `created_at` → 2024-11-04 18:30:00

### Product Data

```bash
# Generate 50 products
lvt seed products --count 50
```

**Generated fields:**
- `name` → "Innovative Widget"
- `description` → "High quality product with excellent features..."
- `price` → 149.99
- `quantity` → 87
- `category` → "Electronics"
- `status` → "active"

### Task Data

```bash
# Generate 200 tasks
lvt seed tasks --count 200
```

**Generated fields:**
- `title` → "Complete Project Documentation"
- `description` → "Finalize all documentation for the Q4 release..."
- `status` → "pending"
- `priority` → "high"

## Test Record Identification

All test records are marked with special IDs starting with `test-seed-`:

```sql
-- In database
id: test-seed-00001
id: test-seed-00002
id: test-seed-00003
```

This allows:
- Easy cleanup (`--cleanup` finds these records)
- Visual identification in database browser
- Prevents mixing with real user data

## Development Workflow

```bash
# 1. Create and setup resource
lvt gen resource products name price:float description
lvt migration up

# 2. Seed initial test data
lvt seed products --count 50

# 3. Develop and test with realistic data
lvt serve
# Browse to http://localhost:8080/products

# 4. Need fresh data?
lvt seed products --cleanup --count 50

# 5. Need more data?
lvt seed products --count 100
# Now 150 total (50 + 100)
```

## Common Patterns

### Testing Pagination

```bash
# Generate enough data to test pagination
lvt seed products --count 100

# Default page size is 20, so 100 records = 5 pages
```

### Testing Search

```bash
# Generate diverse data for search testing
lvt seed products --count 50
lvt seed users --count 30
lvt seed tasks --count 100
```

### Clean State

```bash
# Start fresh for new feature
lvt seed products --cleanup
lvt seed users --cleanup
lvt seed tasks --cleanup

# Reseed with specific amounts
lvt seed products --count 20
lvt seed users --count 10
```

### Demo Data

```bash
# Prepare demo with realistic amounts
lvt seed users --cleanup --count 15
lvt seed products --cleanup --count 30
lvt seed orders --cleanup --count 50
```

## Choosing Record Counts

**Development (fast feedback):**
- Small: 10-20 records (test basic CRUD)
- Medium: 50-100 records (test pagination)
- Large: 500-1000 records (test performance)

**Testing:**
- Unit tests: 5-10 records (specific scenarios)
- E2E tests: 20-50 records (realistic flows)
- Load tests: 1000+ records (stress testing)

**Demo:**
- Small dataset: 10-30 records (easy to navigate)
- Diverse data: Multiple resources with 15-20 each
- Representative: Enough to show features without overwhelming

**Rule of thumb:**
- Start small (20-50 records)
- Increase as needed for specific tests
- Use `--cleanup` to reset between scenarios

## Common Issues

### ❌ Resource Not Found

```bash
# Error: resource 'product' not found in schema

# Cause: Resource doesn't exist in schema.yaml
# The seeder reads your app's schema.yaml to find resource definitions

# Fix 1: Check spelling (must match exactly)
lvt seed products --count 50  # not 'product' (case-sensitive)

# Fix 2: Generate resource first
lvt gen resource products name price:float
lvt migration up
lvt seed products --count 50

# Fix 3: Verify schema.yaml contains the resource
cat schema.yaml | grep products
```

### ❌ Database Not Found

```bash
# Error: database not found

# Cause: Haven't run migrations yet

# Fix: Run migrations first
lvt migration up
lvt seed products --count 50
```

### ❌ No Flags Specified

```bash
# Error: either --count or --cleanup must be specified

# Cause: Forgot to specify what to do

# Fix: Add --count or --cleanup
lvt seed products --count 50  # seed data
lvt seed products --cleanup   # cleanup data
```

### ❌ Foreign Key Constraints

```bash
# Error: FOREIGN KEY constraint failed

# Cause: Seeding child table without parent data

# Fix: Seed parent table first
lvt seed posts --count 20      # parent
lvt seed comments --count 100  # child (references posts)
```

## Field Name Best Practices

**For better generated data, use descriptive field names:**

```bash
# ✅ Good - context-aware generation
lvt gen resource users email phone_number full_name

# ❌ Less good - generic generation
lvt gen resource users field1 field2 field3
```

**Examples:**
- Use `email` not `user_email_address`
- Use `phone` not `contact_number`
- Use `price` not `item_cost`
- Use `description` not `desc`

The seeder recognizes common patterns in field names.

## Cleanup Strategy

**During development:**
```bash
# Clean and reseed often
lvt seed products --cleanup --count 50
```

**Before demo:**
```bash
# Clean all, then seed specific amounts
lvt seed users --cleanup --count 10
lvt seed products --cleanup --count 20
```

**Testing edge cases:**
```bash
# Small dataset
lvt seed products --cleanup --count 5

# Large dataset
lvt seed products --cleanup --count 500
```

## Integration with Testing

```bash
# Before running tests
lvt seed products --cleanup --count 30
go test ./internal/app/products

# After tests
lvt seed products --cleanup
```

## Quick Reference

**I want to...** | **Command**
---|---
Generate 50 products | `lvt seed products --count 50`
Generate 100 users | `lvt seed users --count 100`
Remove test data | `lvt seed products --cleanup`
Fresh start | `lvt seed products --cleanup --count 50`
Add more data | `lvt seed products --count 25` (adds to existing)
Check total | Output shows: "Total test records in products: 75"

## Remember

✓ Seed parent resources before children (foreign keys)
✓ Use descriptive field names for better data
✓ Test records marked with "test-seed-" prefix
✓ `--cleanup` only removes test records, not real data
✓ Seeding adds to existing data unless you use `--cleanup`

✗ Don't forget to run migrations before seeding
✗ Don't seed without checking foreign key dependencies
✗ Use resource name from schema.yaml, not database table name
  - Resource name: `products` (as defined in schema.yaml)
  - Table name: `products` (database table, should match but use resource name for seeding)
