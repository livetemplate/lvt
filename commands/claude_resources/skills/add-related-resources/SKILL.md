---
name: lvt-add-related-resources
description: Intelligently suggest and add related resources based on domain and existing schema - uses context to recommend complementary resources
keywords: ["lvt", "livetemplate", "lt"]
category: workflows
version: 1.0.0
---

# lvt-add-related-resources

Analyzes existing app and intelligently suggests related resources to add. Detects domain patterns and recommends complementary resources with proper relationships.

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
✅ "suggest related resources"
✅ "what should I add next?"
✅ "add common blog resources"

**Without Context (needs keywords):**
✅ "suggest resources for my lvt app"
✅ "what should I add to my livetemplate blog?"
❌ "suggest related resources" (no context, no keywords)

---

## User Prompts

**When to use:**
- "What resources should I add next?"
- "Suggest related resources for my app"
- "What else does my blog need?"
- "Add related resources to my app"
- "What's missing from my schema?"

**Examples:**
- "I have posts, what else should I add?"
- "Suggest resources for my e-commerce app"
- "What resources go well with users?"
- "Add common blog resources"

## How It Works

### Step 1: Analyze Existing Schema

Use **lvt:resource-inspect** skill:
```bash
lvt resource list
```

**Detect domain from existing resources:**
- posts → Blog domain
- products → E-commerce domain
- tasks → Project management domain
- events → Event management domain
- users → User-centric app

### Step 2: Suggest Related Resources

Based on domain analysis, suggest complementary resources:

#### Blog Domain

**If has: posts**
Suggest:
- `comments` (post_id:references:posts:CASCADE, content, author, created_at)
- `categories` (name, slug, description)
- `post_categories` (post_id:references:posts, category_id:references:categories) [junction]
- `tags` (name, slug)
- `post_tags` (post_id:references:posts, tag_id:references:tags) [junction]

#### E-commerce Domain

**If has: products**
Suggest:
- `orders` (user_email, total:float, status, created_at)
- `order_items` (order_id:references:orders, product_id:references:products, quantity:int, price:float)
- `cart_items` (session_id, product_id:references:products, quantity:int)
- `customers` (email, name, address, phone)
- `reviews` (product_id:references:products, user_email, rating:int, content)

#### Project Management Domain

**If has: projects**
Suggest:
- `tasks` (project_id:references:projects, title, description, status, due_date)
- `team_members` (project_id:references:projects, user_email, role)
- `milestones` (project_id:references:projects, title, due_date, completed:bool)
- `comments` (task_id:references:tasks, user_email, content)

**If has: tasks**
Suggest:
- `projects` (name, description, status, owner_email)
- `labels` (name, color)
- `task_labels` (task_id:references:tasks, label_id:references:labels) [junction]
- `time_entries` (task_id:references:tasks, user_email, hours:float, date:time)

#### Event Management Domain

**If has: events**
Suggest:
- `registrations` (event_id:references:events, user_email, status, registered_at)
- `venues` (name, address, capacity:int)
- `speakers` (name, bio, email)
- `event_speakers` (event_id:references:events, speaker_id:references:speakers) [junction]

#### Social/Community Domain

**If has: users**
Suggest:
- `posts` (user_id:references:users, content, created_at)
- `comments` (post_id:references:posts, user_id:references:users, content)
- `likes` (post_id:references:posts, user_id:references:users)
- `follows` (follower_id:references:users, following_id:references:users)

### Step 3: Present Suggestions Using AskUserQuestion Tool

**IMPORTANT**: Use the `AskUserQuestion` tool to let user select which resources to add.

**Present suggestions as checkboxes** (multiSelect: true):

```markdown
## Example: Blog Domain

Question: "Which related resources would you like to add to your blog?"

Options:
1. comments
   Description: "Allow readers to discuss posts (post_id, content, author)"

2. categories
   Description: "Organize posts by topic (name, slug, description)"

3. tags
   Description: "Add flexible labeling with many-to-many (name, slug + junction table)"

4. None
   Description: "I'll add resources manually later"
```

**Use AskUserQuestion tool:**
```
AskUserQuestion({
  questions: [{
    question: "Which related resources would you like to add?",
    header: "Add Resources",
    multiSelect: true,
    options: [
      {label: "comments", description: "Allow readers to discuss posts (post_id, content, author)"},
      {label: "categories", description: "Organize posts by topic (name, slug, description)"},
      {label: "tags", description: "Flexible labeling with many-to-many relationship"}
    ]
  }]
})
```

**Benefits of AskUserQuestion:**
- ✅ User explicitly chooses which resources they want
- ✅ Can select multiple resources at once
- ✅ Can select "Other" to describe custom resource
- ✅ Clear, interactive UI instead of auto-adding everything

### Step 4: Add Selected Resources

For each selected resource, use **lvt:add-resource** skill:

```bash
# User selected: comments, categories

# Add comments
lvt gen resource comments post_id:references:posts:CASCADE content author created_at

# Add categories
lvt gen resource categories name slug description

# Apply migrations
lvt migration up
cd database && sqlc generate && cd ../..
go mod tidy
```

### Step 5: Suggest Relationships

After adding resources, suggest how to use them:

**For comments:**
```
Comments added! You can now:
1. Display comments under each post
2. Add comment form to post detail page
3. Moderate comments (add approved:bool field)

To customize:
- Edit app/comments/comments.go
- Modify app/comments/comments.tmpl
```

**For categories:**
```
Categories added! Next steps:
1. Create post_categories junction table for many-to-many
2. Add category filter to posts list
3. Show category on post cards

Junction table command:
lvt gen schema post_categories post_id:references:posts category_id:references:categories
```

## Domain Pattern Library

### Blog Patterns

**Basic blog:**
- posts (required)
- comments (common)
- categories (common)

**Advanced blog:**
- tags + post_tags (many-to-many)
- authors (if multi-author)
- media/attachments

**Full-featured blog:**
- drafts (status field in posts)
- scheduled posts (publish_at field)
- revisions/versions
- SEO metadata

### E-commerce Patterns

**Basic store:**
- products (required)
- orders (required)
- order_items (junction)

**Advanced store:**
- customers
- cart_items
- reviews + ratings
- product_images

**Full-featured store:**
- inventory tracking
- shipping addresses
- payment methods
- coupons/discounts
- wishlists

### SaaS Patterns

**Basic SaaS:**
- users (via auth)
- organizations/teams
- subscriptions

**Advanced SaaS:**
- usage_tracking
- invoices
- api_keys
- webhooks

### CMS Patterns

**Basic CMS:**
- pages
- media/assets
- menus

**Advanced CMS:**
- content_types (flexible)
- revisions
- workflows (draft→review→publish)

## Checklist

- [ ] Use lvt:resource-inspect to analyze existing schema
- [ ] Detect domain from existing resources
- [ ] Generate suggestions based on domain patterns
- [ ] Present suggestions with rationale
- [ ] Ask user which resources to add
- [ ] For each selected resource:
  - [ ] Use lvt:add-resource to generate
  - [ ] Apply migrations
  - [ ] Generate sqlc models
- [ ] Suggest next steps and relationships
- [ ] Offer to add junction tables if needed
- [ ] Recommend UI/UX improvements

## Smart Suggestions Algorithm

1. **Detect primary resources** (e.g., posts, products, events)
2. **Identify domain** from resource names and fields
3. **Match to pattern library** (blog, e-commerce, etc.)
4. **Suggest complementary resources** in priority order:
   - Essential (core functionality)
   - Common (frequently used)
   - Advanced (nice-to-have)
5. **Detect missing relationships** (foreign keys)
6. **Suggest junction tables** for many-to-many

## Common Relationships

### One-to-Many
```bash
# Parent → Child
posts → comments
products → reviews
projects → tasks
```

### Many-to-Many (via junction)
```bash
# posts ↔ tags
lvt gen schema post_tags post_id:references:posts tag_id:references:tags

# users ↔ roles
lvt gen schema user_roles user_id:references:users role_id:references:roles

# products ↔ categories
lvt gen schema product_categories product_id:references:products category_id:references:categories
```

### Self-Referencing
```bash
# Users following users
lvt gen schema follows follower_id:references:users following_id:references:users

# Categories with parent categories
lvt gen resource categories name parent_id:references:categories
```

## Example Conversations

### Example 1: Blog App
**User:** "I have posts, what should I add next?"

**Response:**
```
Analyzing your blog app... I see you have posts.

Essential additions:
1. comments - Let readers engage with posts
2. categories - Organize posts by topic

Would make your blog better:
3. tags - Flexible post labeling
4. authors - If you have multiple writers
5. media - Image/file uploads

Should I add comments and categories? (Most blogs need these)
```

### Example 2: E-commerce App
**User:** "Suggest resources for my shop"

**Response:**
```
I see you have products. For a complete shop, you need:

Core (required for orders):
1. orders - Track customer purchases
2. order_items - Line items for each order

Customer management:
3. customers - Store customer info
4. cart_items - Shopping cart functionality

Quality features:
5. reviews - Product reviews and ratings
6. inventory - Stock tracking

Shall I add the core resources (orders + order_items)?
```

## Success Criteria

Suggestions are good when:
1. ✅ Aligned with detected domain
2. ✅ Complement existing resources
3. ✅ Include proper relationships (foreign keys)
4. ✅ Prioritized by importance
5. ✅ Explain benefit of each suggestion
6. ✅ Provide next steps after adding

## Notes

- This is an intelligent workflow skill
- Uses pattern matching to detect domain
- Suggests industry-standard relationships
- Adapts to user's specific needs
- Can be combined with lvt:quickstart
- Helps users avoid missing critical resources
- Teaches good schema design practices
