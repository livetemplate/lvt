---
name: lvt-suggest
description: Suggest improvements and next steps based on app analysis - recommends features, optimizations, security enhancements
category: maintenance
version: 1.0.0
keywords: ["lvt", "livetemplate", "lt"]
---

# lvt:suggest

Analyzes app structure and suggests actionable improvements. Provides prioritized recommendations for features, performance, security, and user experience.

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
- "What should I add to my app?"
- "Suggest improvements for my application"
- "What's missing from my app?"
- "How can I make this better?"
- "What features should I add next?"

**Examples:**
- "Suggest improvements for my blog"
- "What should I add to make this production-ready?"
- "How can I improve performance?"
- "What security improvements do I need?"

## How It Works

### Step 1: Analyze Current State

Use **lvt:resource-inspect** to understand the app:
```bash
lvt resource list
lvt resource describe <table>
```

**Gather context:**
- What resources exist?
- What relationships are present?
- What features are implemented?
- What's the domain/purpose?

### Step 2: Categorize Suggestions

Generate recommendations in these categories:

#### 1. Missing Core Features (P0 - Critical)

**Authentication:**
```
⚠️ No authentication system detected
→ Add user authentication to protect resources
→ Command: lvt gen auth
→ Impact: Security, user management, sessions
```

**Essential Relationships:**
```
⚠️ posts table exists but no comments
→ Add comments to enable reader engagement
→ Command: lvt gen resource comments post_id:references:posts:CASCADE content author
→ Impact: User engagement, community building
```

#### 2. Performance Optimizations (P1 - Important)

**Missing Indexes:**
```
⚠️ No index on posts(created_at)
→ Add index for date-based queries
→ Command: lvt migration create add_posts_created_at_index
→ SQL: CREATE INDEX idx_posts_created_at ON posts(created_at);
→ Impact: Faster list/pagination queries
```

**Query Optimization:**
```
⚠️ N+1 query potential in posts → comments
→ Consider eager loading or denormalization
→ Impact: Reduced database roundtrips
```

#### 3. Security Enhancements (P1 - Important)

**CSRF Protection:**
```
⚠️ Forms lack CSRF protection
→ Add CSRF tokens to all forms
→ Already built-in with auth system
→ Impact: Prevent cross-site request forgery
```

**Input Validation:**
```
⚠️ No field length limits
→ Add validation to prevent abuse
→ Add constraints in migrations
→ Impact: Data integrity, DoS prevention
```

#### 4. UX Improvements (P2 - Nice to have)

**Pagination:**
```
💡 Large lists without pagination
→ Add infinite scroll or page-based navigation
→ Already supported by LiveTemplate
→ Impact: Better performance, usability
```

**Search:**
```
💡 No search functionality
→ Add search to help users find content
→ SQL: WHERE title LIKE ? OR content LIKE ?
→ Impact: Improved user experience
```

**Sorting/Filtering:**
```
💡 No sort or filter controls
→ Add UI controls for sorting/filtering
→ Impact: Better content discovery
```

#### 5. Data Management (P2 - Nice to have)

**Soft Deletes:**
```
💡 Hard deletes with CASCADE
→ Consider soft deletes for audit trail
→ Add deleted_at:time field
→ Impact: Data recovery, audit compliance
```

**Timestamps:**
```
💡 Missing updated_at fields
→ Add updated_at for change tracking
→ Impact: Better data lifecycle management
```

### Step 3: Prioritize Recommendations

**Priority matrix:**

**P0 (Critical - Do now):**
- Security vulnerabilities
- Missing authentication (if handling user data)
- Broken relationships
- Data integrity issues

**P1 (Important - Do soon):**
- Performance bottlenecks
- Missing core features for domain
- UX pain points
- Production readiness gaps

**P2 (Nice to have - Do later):**
- Advanced features
- Polish and refinement
- Optional optimizations

### Step 4: Generate Action Plan

**Format suggestions as actionable steps:**

```markdown
## Recommendations for [App Name]

### Critical (Do Now)

1. **Add Authentication System**
   - Why: Protect user data and enable personalization
   - How: `lvt gen auth`
   - Time: 10 minutes
   - Impact: ⭐⭐⭐⭐⭐

2. **Add Missing Foreign Keys**
   - Why: Data integrity and relationships
   - How: `lvt gen resource comments post_id:references:posts content`
   - Time: 5 minutes
   - Impact: ⭐⭐⭐⭐

### Important (Do Soon)

3. **Add Indexes for Performance**
   - Why: Faster queries on large datasets
   - How: Migration to add indexes
   - Time: 5 minutes
   - Impact: ⭐⭐⭐⭐

4. **Add Search Functionality**
   - Why: Better content discovery
   - How: Add search form + WHERE clause
   - Time: 15 minutes
   - Impact: ⭐⭐⭐

### Nice to Have (Do Later)

5. **Add Categories/Tags**
   - Why: Better content organization
   - How: `lvt gen resource categories name slug`
   - Time: 20 minutes
   - Impact: ⭐⭐⭐

6. **Add Soft Deletes**
   - Why: Data recovery and audit trail
   - How: Add deleted_at field + filter queries
   - Time: 30 minutes
   - Impact: ⭐⭐
```

## Domain-Specific Suggestions

### Blog Domain

**If has: posts**
Suggest:
- Comments (user engagement)
- Categories/tags (organization)
- Authors (multi-user)
- Search (content discovery)
- RSS feed (syndication)
- Related posts (engagement)
- View counter (analytics)

### E-commerce Domain

**If has: products**
Suggest:
- Reviews/ratings (social proof)
- Inventory tracking (stock management)
- Categories/filters (navigation)
- Search (product discovery)
- Cart persistence (UX)
- Order history (customer service)
- Product images (visual appeal)

### SaaS Domain

**If has: users (via auth)**
Suggest:
- Organizations/teams (multi-tenant)
- Role-based access control (permissions)
- Usage tracking (billing)
- API keys (integrations)
- Audit logs (compliance)
- Activity feed (transparency)
- Export functionality (data portability)

### Project Management Domain

**If has: projects, tasks**
Suggest:
- Labels/tags (organization)
- Time tracking (productivity)
- Comments/activity (collaboration)
- File attachments (documentation)
- Notifications (updates)
- Dashboard/analytics (insights)
- Kanban board view (workflow)

## Suggestion Patterns

### Pattern 1: Missing CRUD Relationships

**Detection:** Parent resource without children
**Example:** posts without comments
**Suggestion:**
```
Add comments to enable reader engagement
→ lvt gen resource comments post_id:references:posts:CASCADE content author created_at
→ Impact: Community building, feedback loop
```

### Pattern 2: No Search/Filter

**Detection:** Large datasets without search
**Example:** 100+ products without search
**Suggestion:**
```
Add search functionality to help users find products
→ Add search input to list template
→ Add WHERE clause: WHERE name LIKE ? OR description LIKE ?
→ Impact: Better user experience, faster discovery
```

### Pattern 3: Missing Authentication

**Detection:** User-specific resources without auth
**Example:** tasks, posts with "user_email" but no auth system
**Suggestion:**
```
Add authentication to protect user data
→ lvt gen auth
→ Use auth.RequireAuth middleware on routes
→ Impact: Security, user management
```

### Pattern 4: Performance Issues

**Detection:** Tables > 1000 rows without indexes on common queries
**Example:** posts sorted by created_at without index
**Suggestion:**
```
Add index on created_at for faster sorting
→ lvt migration create add_posts_created_at_index
→ CREATE INDEX idx_posts_created_at ON posts(created_at);
→ Impact: 10-100x faster queries
```

### Pattern 5: Missing Production Features

**Detection:** App ready to deploy but missing prod features
**Example:** No health checks, no monitoring
**Suggestion:**
```
Add production readiness features
→ Health check endpoint
→ Structured logging
→ Error monitoring (Sentry)
→ Database backups (Litestream)
→ Impact: Reliability, observability
```

## Checklist

- [ ] Use lvt:resource-inspect to analyze schema
- [ ] Detect domain from existing resources
- [ ] Identify missing core features (auth, relationships)
- [ ] Check for performance issues (indexes, N+1)
- [ ] Review security posture (CSRF, validation)
- [ ] Assess UX completeness (search, pagination)
- [ ] Categorize by priority (P0/P1/P2)
- [ ] Generate actionable recommendations
- [ ] Provide specific commands/code
- [ ] Estimate time and impact
- [ ] Present in order of priority

## Output Format

```markdown
# Suggestions for [App Name]

## Overview
Your app is a [domain] with [X] resources. Here are prioritized recommendations:

## 🔴 Critical (Do Now)

### 1. [Suggestion Title]
**Why:** [Business/technical reason]
**How:**
```bash
[Exact commands to run]
```
**Time:** [Estimate]
**Impact:** [What changes]

## 🟡 Important (Do Soon)

### 2. [Suggestion Title]
...

## 🟢 Nice to Have (Do Later)

### 3. [Suggestion Title]
...

## Next Steps

1. Start with critical items (🔴)
2. Test each change before moving to next
3. Re-run suggestions after implementing top items
```

## Example Suggestions

### Example 1: Simple Blog

**Analysis:**
- Has: posts (title, content, created_at)
- Missing: comments, categories, auth

**Suggestions:**
```
## 🔴 Critical

1. Add User Authentication
   - Why: Protect post creation/editing
   - How: lvt gen auth
   - Time: 10 min
   - Impact: Security, author attribution

## 🟡 Important

2. Add Comments
   - Why: Reader engagement
   - How: lvt gen resource comments post_id:references:posts:CASCADE content author
   - Time: 5 min
   - Impact: Community building

3. Add Categories
   - Why: Content organization
   - How: lvt gen resource categories name slug
   - Time: 5 min
   - Impact: Better navigation

## 🟢 Nice to Have

4. Add Search
   - Why: Content discovery
   - How: Add WHERE title LIKE ? OR content LIKE ?
   - Time: 15 min
   - Impact: Better UX
```

### Example 2: E-commerce App

**Analysis:**
- Has: products (name, price, quantity)
- Missing: orders, cart, reviews, categories

**Suggestions:**
```
## 🔴 Critical

1. Add Order Management
   - Why: Can't process sales without orders
   - How:
     lvt gen resource orders user_email total:float status
     lvt gen schema order_items order_id:references:orders product_id:references:products quantity:int price:float
   - Time: 15 min
   - Impact: Enable sales

2. Add Authentication
   - Why: Track customer orders
   - How: lvt gen auth
   - Time: 10 min
   - Impact: Customer accounts

## 🟡 Important

3. Add Product Reviews
   - Why: Social proof increases sales
   - How: lvt gen resource reviews product_id:references:products rating:int content
   - Time: 10 min
   - Impact: Trust, conversions

4. Add Shopping Cart
   - Why: Multi-item purchases
   - How: lvt gen schema cart_items session_id product_id:references:products quantity:int
   - Time: 20 min
   - Impact: Better UX, higher AOV
```

## Success Criteria

Good suggestions should:
1. ✅ Be specific and actionable
2. ✅ Include exact commands
3. ✅ Explain business value
4. ✅ Estimate time and impact
5. ✅ Be prioritized by urgency
6. ✅ Be domain-appropriate
7. ✅ Be achievable quickly

## Common Suggestion Categories

**Features:**
- Authentication/authorization
- Search and filtering
- Pagination/infinite scroll
- Comments/reviews
- Categories/tags
- File uploads
- Email notifications

**Performance:**
- Database indexes
- Query optimization
- Caching strategies
- Asset optimization
- Connection pooling

**Security:**
- CSRF protection
- Input validation
- Rate limiting
- SQL injection prevention
- XSS prevention

**Data Management:**
- Soft deletes
- Audit trails
- Data exports
- Backups
- Migrations

**Production:**
- Health checks
- Monitoring
- Error tracking
- Logging
- Deployment config

## Notes

- Suggestions should be contextual (based on actual analysis)
- Always explain "why" not just "what"
- Provide specific commands, not generic advice
- Estimate time realistically (err on high side)
- Focus on high-impact, low-effort wins first
- Re-run suggestions after implementing to find new opportunities
- Combine with lvt:add-related-resources for resource suggestions
- Combine with lvt:production-ready for deployment suggestions
