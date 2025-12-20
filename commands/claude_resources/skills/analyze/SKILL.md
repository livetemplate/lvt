---
name: lvt-analyze
description: Analyze LiveTemplate app structure - examine schema, resources, relationships, complexity, and provide insights
category: maintenance
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:analyze

Comprehensive analysis of LiveTemplate application structure. Examines database schema, resources, relationships, code organization, and provides actionable insights.

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

## User Prompts

**When to use:**
- "Analyze my app"
- "Show me my app structure"
- "What's in my application?"
- "Review my schema"
- "Tell me about my app's architecture"

**Examples:**
- "Analyze my blog app"
- "Show me the structure of my application"
- "Review my database schema"
- "What resources do I have?"

## Analysis Components

### 1. Schema Analysis

Use **lvt:resource-inspect** to examine:
- All tables and their purposes
- Field counts and complexity
- Data types used
- Relationships (foreign keys)
- Indexes

**Output format:**
```
=== Schema Analysis ===

Resources: 5 tables found

Core Resources:
- users (7 fields) - User accounts and profiles
- posts (6 fields) - Blog posts with content
- comments (5 fields) - User comments on posts

Supporting Tables:
- sessions (4 fields) - User session management
- categories (3 fields) - Post categorization

Complexity: Medium (5 tables, 25 total fields)
```

### 2. Relationship Analysis

Detect and document:
- One-to-many relationships
- Many-to-many relationships
- Self-referencing relationships
- Missing relationships (potential gaps)

**Output format:**
```
=== Relationships ===

One-to-Many:
- posts → comments (via post_id)
- users → posts (via user_id)
- users → sessions (via user_id)

Many-to-Many:
- (none detected)

Potential Missing Relationships:
- posts could relate to categories
- comments could relate to users

Foreign Key Analysis:
✅ All foreign keys use CASCADE delete
✅ Relationships properly indexed
```

### 3. Resource Complexity

Analyze each resource:
- Field count (simple: <5, medium: 5-10, complex: >10)
- Field types (diversity)
- Business logic hints
- UI complexity

**Output format:**
```
=== Resource Complexity ===

Simple Resources:
- categories (3 fields) - Basic lookup table
- sessions (4 fields) - Session tracking

Medium Resources:
- comments (5 fields) - Standard CRUD
- posts (6 fields) - Content management

Complex Resources:
- users (7 fields) - Auth + profile data

Average Complexity: 5 fields per resource
```

### 4. Code Organization

Check file structure:
```
=== Code Organization ===

App Structure:
✅ app/ - Handler organization
✅ database/ - Database layer
✅ shared/ - Shared utilities

Resources with Full Stack:
- posts (handler, template, tests)
- comments (handler, template, tests)
- users (auth system)

Views:
- home (landing page)
- about (static page)

Missing:
⚠️  No E2E tests for posts
⚠️  No custom middleware
```

### 5. Database Health

Analyze database structure:
- Migration count and history
- Index coverage
- Potential performance issues
- Schema consistency

**Output format:**
```
=== Database Health ===

Migrations: 5 applied
Latest: 20251104_create_comments.sql

Index Coverage:
✅ Primary keys on all tables
✅ Foreign keys indexed
⚠️  High-volume table 'posts' missing index on created_at

Performance Recommendations:
- Add index on posts(created_at) for date sorting
- Add index on comments(post_id, created_at) for pagination

Schema Consistency:
✅ All timestamps use DATETIME
✅ Consistent naming (snake_case)
✅ ID fields use INTEGER PRIMARY KEY
```

### 6. Feature Detection

Identify features in use:
- Authentication (password, magic link, email confirm)
- Authorization (middleware, protected routes)
- CRUD operations
- Search/filter
- Pagination
- Sorting

**Output format:**
```
=== Features Detected ===

Authentication:
✅ Password authentication
✅ Magic link authentication
✅ Email confirmation
✅ Session management
✅ CSRF protection

CRUD Operations:
✅ posts - Full CRUD
✅ comments - Full CRUD
✅ categories - Full CRUD

Advanced Features:
✅ Pagination (infinite scroll)
⚠️  No search functionality
⚠️  No sorting controls
⚠️  No filters
```

## Checklist

- [ ] Run lvt resource list to get all tables
- [ ] For each resource, run lvt resource describe
- [ ] Analyze relationships and foreign keys
- [ ] Calculate complexity metrics
- [ ] Check code organization
- [ ] Review database migrations
- [ ] Detect features in use
- [ ] Identify missing indexes
- [ ] Suggest improvements
- [ ] Generate comprehensive report

## Analysis Report Template

```markdown
# Application Analysis Report

## Overview
- **App Name:** [name]
- **Resources:** [count] tables
- **Complexity:** [simple/medium/complex]
- **Features:** [auth, CRUD, pagination, etc.]

## Schema Summary
[List all resources with field counts]

## Relationships
[Document all foreign keys and relationships]

## Complexity Analysis
[Breakdown by resource]

## Database Health
- **Migrations:** [count]
- **Index Coverage:** [percentage]
- **Performance:** [issues/recommendations]

## Feature Coverage
[Which features are implemented]

## Recommendations
[Ordered list of suggestions]

## Next Steps
[Actionable items for improvement]
```

## Example Analysis

### Small Blog App
```
=== App Analysis ===

Overview:
- 3 core resources (posts, comments, users)
- Medium complexity (18 total fields)
- Basic CRUD + Auth

Strengths:
✅ Clean schema with proper relationships
✅ Full authentication system
✅ Consistent naming conventions

Opportunities:
⚠️  Add categories for better organization
⚠️  Add indexes for better performance
⚠️  Add search functionality

Recommendation: Add categories and post_categories junction table
```

### E-commerce App
```
=== App Analysis ===

Overview:
- 6 core resources (products, orders, customers, etc.)
- Complex (45 total fields)
- Full e-commerce features

Strengths:
✅ Complete order workflow
✅ Customer management
✅ Product catalog

Opportunities:
⚠️  Missing inventory tracking
⚠️  No product reviews
⚠️  Cart abandonment not tracked

Recommendation: Add reviews resource and inventory fields
```

## Metrics to Calculate

### Schema Metrics
- Total tables
- Total fields
- Average fields per table
- Relationship count
- Index count

### Code Metrics
- Handlers with tests
- Template coverage
- Middleware count
- Shared utilities

### Feature Metrics
- Auth completeness (0-100%)
- CRUD coverage (0-100%)
- Advanced features (search, sort, filter)

## Common Insights

### Well-Structured App
- Consistent naming
- Proper relationships
- Good index coverage
- Comprehensive tests
- Clear separation of concerns

### Needs Improvement
- Missing relationships
- No indexes on high-volume queries
- Inconsistent naming
- Missing tests
- Monolithic handlers

## Success Criteria

Analysis is complete when:
1. ✅ All resources documented
2. ✅ Relationships mapped
3. ✅ Complexity calculated
4. ✅ Code organization reviewed
5. ✅ Database health checked
6. ✅ Features detected
7. ✅ Recommendations provided

## Notes

- Purely informational (no modifications)
- Helps users understand their app
- Identifies improvement opportunities
- Great for onboarding new developers
- Useful before refactoring
- Can inform testing priorities
- Combines data from multiple inspection tools
