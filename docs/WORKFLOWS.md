# LiveTemplate Common Workflows

This guide provides step-by-step workflows for common development tasks using LiveTemplate MCP tools. Each workflow is designed to be executable by AI assistants through the MCP protocol.

## Table of Contents

- [Quick Start Workflows](#quick-start-workflows)
- [Authentication Workflows](#authentication-workflows)
- [Resource Management Workflows](#resource-management-workflows)
- [Database Workflows](#database-workflows)
- [Development Workflows](#development-workflows)
- [Deployment Workflows](#deployment-workflows)

---

## Quick Start Workflows

### 1. Create a Simple Blog

**Goal:** Build a blog with posts and comments in under 2 minutes.

**Steps:**

1. **Create new app:**
   ```
   Tool: lvt_new
   Input: {name: "myblog", kit: "multi"}
   ```

2. **Add authentication:**
   ```
   Tool: lvt_gen_auth
   Input: {}
   ```

3. **Add posts resource:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "posts",
     fields: {
       title: "string",
       content: "text",
       user_id: "references:users",
       published: "bool"
     }
   }
   ```

4. **Add comments resource:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "comments",
     fields: {
       content: "text",
       post_id: "references:posts",
       user_id: "references:users"
     }
   }
   ```

5. **Apply database changes:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

6. **Generate test data:**
   ```
   Tool: lvt_seed
   Input: {resource: "posts", count: 10}

   Tool: lvt_seed
   Input: {resource: "comments", count: 30}
   ```

**Result:** Fully functional blog with auth, posts, and comments.

---

### 2. Build a Task Manager

**Goal:** Create a task management app with projects and tasks.

**Steps:**

1. **Create app:**
   ```
   Tool: lvt_new
   Input: {name: "taskmanager", kit: "single", css: "tailwind"}
   ```

2. **Add authentication:**
   ```
   Tool: lvt_gen_auth
   Input: {}
   ```

3. **Add projects:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "projects",
     fields: {
       name: "string",
       description: "text",
       user_id: "references:users",
       deadline: "time"
     }
   }
   ```

4. **Add tasks:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "tasks",
     fields: {
       title: "string",
       description: "text",
       project_id: "references:projects",
       completed: "bool",
       priority: "int",
       due_date: "time"
     }
   }
   ```

5. **Run migrations:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

6. **Seed test data:**
   ```
   Tool: lvt_seed
   Input: {resource: "projects", count: 5}

   Tool: lvt_seed
   Input: {resource: "tasks", count: 25}
   ```

**Result:** Task manager with projects, tasks, and user authentication.

---

### 3. E-Commerce Store Setup

**Goal:** Set up a product catalog with categories and orders.

**Steps:**

1. **Create store app:**
   ```
   Tool: lvt_new
   Input: {name: "mystore", kit: "multi"}
   ```

2. **Add user authentication:**
   ```
   Tool: lvt_gen_auth
   Input: {struct_name: "Customer"}
   ```

3. **Add categories:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "categories",
     fields: {
       name: "string",
       description: "text",
       slug: "string"
     }
   }
   ```

4. **Add products:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "products",
     fields: {
       name: "string",
       description: "text",
       price: "float",
       stock: "int",
       category_id: "references:categories",
       image_url: "string"
     }
   }
   ```

5. **Add orders:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "orders",
     fields: {
       customer_id: "references:customers",
       total: "float",
       status: "string",
       created_at: "time"
     }
   }
   ```

6. **Add order items:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "order_items",
     fields: {
       order_id: "references:orders",
       product_id: "references:products",
       quantity: "int",
       price: "float"
     }
   }
   ```

7. **Apply all migrations:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

8. **Seed catalog:**
   ```
   Tool: lvt_seed
   Input: {resource: "categories", count: 5}

   Tool: lvt_seed
   Input: {resource: "products", count: 50}
   ```

**Result:** Complete e-commerce foundation.

---

## Authentication Workflows

### 4. Add Authentication to Existing App

**Scenario:** You have an app without auth, need to add it.

**Steps:**

1. **Check existing resources:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

2. **Generate auth system:**
   ```
   Tool: lvt_gen_auth
   Input: {}
   ```

3. **Apply auth migration:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

4. **Create test users:**
   ```
   Tool: lvt_seed
   Input: {resource: "users", count: 10}
   ```

5. **Verify auth tables:**
   ```
   Tool: lvt_resource_describe
   Input: {resource: "users"}
   ```

---

### 5. Custom Authentication Setup

**Scenario:** Need custom user table and fields.

**Steps:**

1. **Generate custom auth:**
   ```
   Tool: lvt_gen_auth
   Input: {
     struct_name: "Account",
     table_name: "admin_users",
     no_magic_link: true,
     no_email_confirm: true
   }
   ```

2. **Add custom profile fields via migration:**
   ```
   Tool: lvt_migration_create
   Input: {name: "add_admin_profile_fields"}
   ```

   *Then manually edit migration file to add:*
   - role column
   - department column
   - permissions column

3. **Run migrations:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

---

## Resource Management Workflows

### 6. Add Resource with Relationships

**Scenario:** Adding a resource that references multiple other resources.

**Steps:**

1. **List existing resources:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

2. **Generate related resource:**
   ```
   Tool: lvt_gen_resource
   Input: {
     name: "articles",
     fields: {
       title: "string",
       content: "text",
       author_id: "references:users",
       category_id: "references:categories",
       published_at: "time",
       status: "string"
     }
   }
   ```

3. **Verify relationships:**
   ```
   Tool: lvt_resource_describe
   Input: {resource: "articles"}
   ```

4. **Apply changes:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

---

### 7. Add View-Only Pages

**Scenario:** Need static pages (about, contact, dashboard).

**Steps:**

1. **Add about page:**
   ```
   Tool: lvt_gen_view
   Input: {name: "about"}
   ```

2. **Add contact page:**
   ```
   Tool: lvt_gen_view
   Input: {name: "contact"}
   ```

3. **Add dashboard:**
   ```
   Tool: lvt_gen_view
   Input: {name: "dashboard"}
   ```

4. **Verify files created:**
   *Check app/about/, app/contact/, app/dashboard/*

---

## Database Workflows

### 8. Safe Migration Workflow

**Best practice:** Always check before applying migrations.

**Steps:**

1. **Check current status:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

2. **Review pending migrations:**
   *Examine migration files listed as pending*

3. **Apply migrations:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

4. **Verify success:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

5. **Inspect updated schemas:**
   ```
   Tool: lvt_resource_list
   Input: {}

   Tool: lvt_resource_describe
   Input: {resource: "<newly added resource>"}
   ```

---

### 9. Rollback Migration

**Scenario:** Need to undo a migration.

**Warning:** May cause data loss. Use with caution.

**Steps:**

1. **Check what will be rolled back:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

2. **Backup database** (manual step)

3. **Rollback:**
   ```
   Tool: lvt_migration_down
   Input: {}
   ```

4. **Verify rollback:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

---

### 10. Add Custom Migration

**Scenario:** Need to add indexes, modify existing tables, or run custom SQL.

**Steps:**

1. **Create migration file:**
   ```
   Tool: lvt_migration_create
   Input: {name: "add_performance_indexes"}
   ```

2. **Edit migration file** (manual step)
   *Add custom SQL for indexes, constraints, etc.*

3. **Verify syntax:**
   *Review SQL before applying*

4. **Apply migration:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

---

## Development Workflows

### 11. Set Up Development Environment

**Scenario:** New developer joining the project.

**Steps:**

1. **List available resources:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

2. **Inspect each resource:**
   ```
   Tool: lvt_resource_describe
   Input: {resource: "users"}

   Tool: lvt_resource_describe
   Input: {resource: "posts"}

   ... (for each resource)
   ```

3. **Run all migrations:**
   ```
   Tool: lvt_migration_up
   Input: {}
   ```

4. **Seed development data:**
   ```
   Tool: lvt_seed
   Input: {resource: "users", count: 10}

   Tool: lvt_seed
   Input: {resource: "posts", count: 50, cleanup: true}

   Tool: lvt_seed
   Input: {resource: "comments", count: 100, cleanup: true}
   ```

5. **Generate environment template:**
   ```
   Tool: lvt_env_generate
   Input: {}
   ```

6. **Create .env from .env.example** (manual step)

---

### 12. Validate Templates Before Deployment

**Scenario:** Ensure all templates are valid before deploying.

**Steps:**

1. **List resources to find templates:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

2. **Validate each template:**
   ```
   Tool: lvt_validate_template
   Input: {template_file: "app/posts/posts.tmpl"}

   Tool: lvt_validate_template
   Input: {template_file: "app/comments/comments.tmpl"}

   ... (for each resource)
   ```

3. **Fix any validation errors** (manual step)

4. **Re-validate after fixes**

---

### 13. Refresh Test Data

**Scenario:** Need fresh test data for development.

**Steps:**

1. **Clean and reseed all resources:**
   ```
   Tool: lvt_seed
   Input: {resource: "users", count: 20, cleanup: true}

   Tool: lvt_seed
   Input: {resource: "posts", count: 100, cleanup: true}

   Tool: lvt_seed
   Input: {resource: "comments", count: 300, cleanup: true}
   ```

2. **Verify data:**
   *Check application UI or database*

---

## Deployment Workflows

### 14. Pre-Deployment Checklist

**Goal:** Ensure app is ready for production.

**Steps:**

1. **Verify all migrations applied:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

2. **List all resources:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

3. **Validate all templates:**
   ```
   Tool: lvt_validate_template
   Input: {template_file: "app/posts/posts.tmpl"}

   ... (for each template)
   ```

4. **Generate production env template:**
   ```
   Tool: lvt_env_generate
   Input: {}
   ```

5. **Review generated .env.example**
   *Ensure all required variables are present*

6. **Manual checks:**
   - Review security settings
   - Check CSRF protection enabled
   - Verify email configuration
   - Test authentication flows

---

### 15. Post-Deployment Verification

**Scenario:** App deployed, need to verify it's working.

**Steps:**

1. **Check migration status on production:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```
   *Run in production environment*

2. **Verify all resources exist:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

3. **DO NOT seed production data** (use real data only)

4. **Manual verification:**
   - Test critical user flows
   - Verify database connections
   - Check logs for errors
   - Test authentication

---

## Troubleshooting Workflows

### 16. Diagnose Migration Issues

**Problem:** Migrations failing or database out of sync.

**Steps:**

1. **Check migration status:**
   ```
   Tool: lvt_migration_status
   Input: {}
   ```

2. **List resources:**
   ```
   Tool: lvt_resource_list
   Input: {}
   ```

3. **Describe problematic resource:**
   ```
   Tool: lvt_resource_describe
   Input: {resource: "<failing resource>"}
   ```

4. **Common fixes:**
   - Rollback last migration if safe
   - Check migration file syntax
   - Verify database permissions
   - Check for conflicting migrations

---

### 17. Fix Template Errors

**Problem:** Templates not rendering correctly.

**Steps:**

1. **Validate template:**
   ```
   Tool: lvt_validate_template
   Input: {template_file: "<path to template>"}
   ```

2. **Review validation output**

3. **Common issues:**
   - Unclosed tags
   - Missing template definitions
   - Syntax errors in Go template expressions
   - Incorrect component references

4. **Fix and re-validate:**
   ```
   Tool: lvt_validate_template
   Input: {template_file: "<path to template>"}
   ```

---

## Best Practices

### General Guidelines

1. **Always check status before migrations:**
   ```
   lvt_migration_status → Review → lvt_migration_up
   ```

2. **Use cleanup when reseeding:**
   ```
   lvt_seed with cleanup: true
   ```

3. **Validate templates before deployment:**
   ```
   lvt_validate_template → Fix → Deploy
   ```

4. **Generate auth first, then resources:**
   ```
   lvt_gen_auth → lvt_gen_resource with user_id references
   ```

5. **List before describe:**
   ```
   lvt_resource_list → Pick resource → lvt_resource_describe
   ```

### Workflow Patterns

**Pattern: Check → Generate → Apply → Verify**

```
1. Check current state (list, status)
2. Generate new code (gen commands)
3. Apply changes (migration up)
4. Verify success (describe, validate)
```

**Pattern: Clean Slate Development**

```
1. Seed with cleanup: true
2. Test thoroughly
3. Seed with cleanup: true again
4. Repeat
```

**Pattern: Safe Production Deployment**

```
1. Test in dev environment
2. Validate all templates
3. Check migration status
4. Generate env template
5. Deploy code
6. Run migrations on production
7. Verify with resource commands
8. Never seed production
```

---

## Quick Reference

Common command sequences:

```bash
# New project setup
lvt_new → lvt_gen_auth → lvt_gen_resource → lvt_migration_up → lvt_seed

# Add feature
lvt_gen_resource → lvt_migration_up → lvt_seed

# Pre-deployment
lvt_migration_status → lvt_validate_template → lvt_env_generate

# Development refresh
lvt_seed (cleanup: true) → Test → Repeat

# Troubleshooting
lvt_resource_list → lvt_resource_describe → lvt_migration_status
```

---

## Support

For detailed tool documentation, see [MCP_TOOLS.md](./MCP_TOOLS.md).

For issues: https://github.com/livetemplate/lvt/issues
