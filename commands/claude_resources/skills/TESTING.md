# Testing lvt Skills

Comprehensive test scenarios for all 20 LiveTemplate skills + Phase 6 enhancements (environment management).

## How to Test Skills

1. **Copy a test prompt** from the scenarios below
2. **Paste into Claude Code** 
3. **Verify Claude:**
   - Mentions using the appropriate skill
   - Follows the skill's guidance accurately
   - Provides correct commands and parameters
   - Prevents common mistakes documented in the skill
   - Completes the workflow successfully

## Quick Smoke Test

**Complete production workflow (tests multiple skills):**
```
Create a new blog app, add posts with title and content, add authentication,
add an about page, generate test data, run the app, analyze it, suggest
improvements, and show me how to deploy to Fly.io
```

**Expected:** Uses 9+ skills in sequence, completes full workflow

---

## Core Command Skills (14 skills)

### 1. lvt:new-app

**Test 1 - Basic app creation:**
```
Create a new LiveTemplate app called "taskmanager"
```

**Expected:**
- Uses lvt:new-app skill
- Asks about kit preference (multi/single/simple)
- Runs `lvt new taskmanager --kit [choice]`
- Verifies structure created
- Shows next steps

**Test 2 - Kit selection guidance:**
```
Create a new app but I'm not sure which CSS framework to use
```

**Expected:**
- Explains multi (Tailwind, full layout) vs single (Tailwind SPA) vs simple (no CSS)
- Helps choose based on app type
- Creates with chosen kit

**Test 3 - Custom module:**
```
Create app "myshop" with module github.com/user/myshop
```

**Expected:**
- Uses --module flag
- Explains module naming
- Verifies go.mod correctness

---

### 2. lvt:add-resource

**Test 1 - Basic resource:**
```
Add a products resource with name and price
```

**Expected:**
- Uses lvt:add-resource skill
- Infers types (name:string, price:float)
- Runs `lvt gen resource products name price:float`
- Runs `lvt migration up`
- Generates sqlc models
- Shows created files

**Test 2 - Complex types:**
```
Add users with email, age (integer), active (boolean), joined (datetime)
```

**Expected:**
- Maps types correctly: email:string, age:int, active:bool, joined:time
- Shows type mapping table
- Handles timestamps

**Test 3 - Foreign keys:**
```
Add comments with post_id referencing posts, content, and author
```

**Expected:**
- Uses foreign key syntax: post_id:references:posts:CASCADE
- Explains CASCADE vs RESTRICT
- Creates proper relationship

**Test 4 - Troubleshooting:**
```
I got "undefined: queries" after adding a resource
```

**Expected:**
- References Common Issues
- Diagnoses: missing sqlc generate
- Provides fix: cd internal/database && sqlc generate

---

### 3. lvt:add-view

**Test 1 - Basic view:**
```
Add a dashboard view
```

**Expected:**
- Uses lvt:add-view skill
- Runs `lvt gen view dashboard`
- Explains no database needed
- Shows handler location

**Test 2 - Multiple views:**
```
Add analytics, reports, and settings views
```

**Expected:**
- Creates 3 separate views
- Verifies all created
- Suggests customization

**Test 3 - View vs resource:**
```
Should I use a view or resource for a counter?
```

**Expected:**
- Explains difference
- Recommends view (no persistence needed)
- Shows WebSocket state example

---

### 4. lvt:add-migration

**Test 1 - Add index:**
```
Add an index on products.price for faster sorting
```

**Expected:**
- Uses lvt:add-migration skill
- Creates migration file
- Shows SQL: CREATE INDEX idx_products_price ON products(price)
- Runs lvt migration up

**Test 2 - Unique constraint:**
```
Make user emails unique
```

**Expected:**
- Creates migration with UNIQUE constraint
- Shows Up/Down SQL
- Warns about existing duplicates

**Test 3 - Migration status:**
```
How do I check which migrations have run?
```

**Expected:**
- Runs `lvt migration status`
- Shows applied/pending migrations

---

### 5. lvt:gen-schema

**Test 1 - Backend table:**
```
Create a schema-only analytics table with event, timestamp, user_id
```

**Expected:**
- Uses lvt:gen-schema skill
- Runs `lvt gen schema analytics event timestamp:time user_id`
- No handler/template generated
- Explains use case (logs, cache, backend)

**Test 2 - Sessions table:**
```
Add a sessions table for auth without UI
```

**Expected:**
- Creates schema + queries
- Skips handler
- Mentions use with gen-auth

---

### 6. lvt:gen-auth ⭐

**Test 1 - Full auth:**
```
Add authentication to my app
```

**Expected:**
- Uses lvt:gen-auth skill
- Runs `lvt gen auth`
- Lists generated files
- Shows route wiring example
- Explains email configuration
- Runs migrations

**Test 2 - Custom auth:**
```
Add auth with password only, no magic link
```

**Expected:**
- Uses --no-magic-link flag
- Generates password-only auth
- Shows disabled features

**Test 3 - Protect routes:**
```
How do I protect my /dashboard route?
```

**Expected:**
- Shows RequireAuth middleware
- Example: auth.RequireAuth(queries, dashboardHandler)
- Explains session checking

**Test 4 - Email setup:**
```
How do I configure email for magic links?
```

**Expected:**
- Shows ConsoleEmailSender for dev
- Shows SMTPSender for production
- Explains env vars (SMTP_HOST, SMTP_USER, etc.)

---

### 7. lvt:resource-inspect

**Test 1 - List resources:**
```
Show me all my resources
```

**Expected:**
- Uses lvt:resource-inspect skill
- Runs `lvt resource list`
- Shows all tables

**Test 2 - Describe table:**
```
What's in the posts table?
```

**Expected:**
- Runs `lvt resource describe posts`
- Shows columns, types, indexes, foreign keys

---

### 8. lvt:manage-kits

**Test 1 - List kits:**
```
What kits are available?
```

**Expected:**
- Uses lvt:manage-kits skill
- Runs `lvt kits list`
- Shows multi, single, simple

**Test 2 - Kit details:**
```
Show me details about the multi kit
```

**Expected:**
- Runs `lvt kits info multi`
- Shows components, CSS framework

**Test 3 - Custom kit:**
```
Can I create a custom kit?
```

**Expected:**
- Explains lvt kits create
- Shows directory structure
- Mentions validation

---

### 9. lvt:validate-templates

**Test 1 - Validate template:**
```
Check my posts template for errors
```

**Expected:**
- Uses lvt:validate-templates skill
- Runs `lvt parse internal/app/posts/posts.tmpl`
- Shows validation results

**Test 2 - Template error:**
```
I'm getting "template: unexpected EOF"
```

**Expected:**
- Diagnoses: unclosed {{ }}
- Suggests running lvt parse
- Shows how to find the error

---

### 10. lvt:run-and-test

**Test 1 - Start server:**
```
Run my app
```

**Expected:**
- Uses lvt:run-and-test skill
- Runs `lvt serve`
- Shows URL (http://localhost:8080)
- Mentions browser auto-open

**Test 2 - Custom port:**
```
Start on port 3000
```

**Expected:**
- Uses --port flag or PORT env var
- Runs server on 3000

**Test 3 - Port conflict:**
```
Getting "port already in use"
```

**Expected:**
- Diagnoses port conflict
- Shows how to find process: lsof -i :8080
- Suggests kill or different port

---

### 11. lvt:customize

**Test 1 - Add filtering:**
```
Add search to my posts list
```

**Expected:**
- Uses lvt:customize skill
- Modifies handler (WHERE clause)
- Updates template (search form)
- Shows SQL and HTML

**Test 2 - Custom query:**
```
Add a query to get recent posts
```

**Expected:**
- Edits queries.sql
- Runs sqlc generate
- Shows handler usage

**Test 3 - WebSocket action:**
```
Add a delete button with WebSocket
```

**Expected:**
- Shows lf-action syntax
- Handler case statement
- Partial template

---

### 12. lvt:seed-data

**Test 1 - Generate data:**
```
Generate 100 test users
```

**Expected:**
- Uses lvt:seed-data skill
- Runs `lvt seed users --count 100`
- Shows realistic data created

**Test 2 - Cleanup:**
```
Remove all test data
```

**Expected:**
- Runs `lvt seed users --cleanup`
- Confirms deletion

**Test 3 - Seed and cleanup:**
```
Clean up old data and generate 50 new posts
```

**Expected:**
- Runs `lvt seed posts --count 50 --cleanup`
- Cleans then seeds

---

### 13. lvt:deploy

**Test 1 - Fly.io deployment:**
```
Deploy my app to Fly.io
```

**Expected:**
- Uses lvt:deploy skill
- Runs `lvt gen stack fly`
- Shows fly.toml
- Explains litestream
- Runs fly deploy

**Test 2 - Docker:**
```
Create a Dockerfile for my app
```

**Expected:**
- Runs `lvt gen stack docker`
- Shows Dockerfile and docker-compose.yml
- Explains CGO_ENABLED

**Test 3 - Migrations in production:**
```
How do I handle migrations when deploying?
```

**Expected:**
- Shows migration workflow
- Explains running before deploy
- Mentions backup strategy

---

## Workflow Orchestration Skills (3 skills)

### 14. lvt:manage-env ⭐ NEW

**Test 1 - Initial setup:**
```
Set up environment variables for my app
```

**Expected:**
- Runs `lvt env generate` to create .env.example
- Copies to .env
- Runs `lvt env validate` to show missing vars
- Guides user to set APP_ENV, DATABASE_PATH, SESSION_SECRET
- Shows how to generate secrets: `openssl rand -hex 32`
- Validates again to confirm success

**Test 2 - List variables:**
```
What environment variables are configured?
```

**Expected:**
- Uses lvt:manage-env skill
- Runs `lvt env list`
- Shows masked sensitive values
- Marks required variables with [REQUIRED]
- Suggests --show-values flag

**Test 3 - Validate configuration:**
```
Am I ready to deploy?
```

**Expected:**
- Uses lvt:manage-env skill
- Runs `lvt env validate --strict`
- Checks all required vars are set
- Validates APP_ENV is valid
- Checks secrets are strong (32+ chars)
- Detects placeholder values
- Reports ready or lists issues

**Test 4 - Configure SMTP:**
```
Set up SMTP for production
```

**Expected:**
- Asks for SMTP provider (Gmail, SendGrid, etc.)
- Guides setting SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS
- Sets EMAIL_PROVIDER=smtp
- Shows provider-specific instructions (e.g., Gmail app password)
- Validates configuration with strict mode

**Test 5 - Missing required variables:**
```
Validate environment
```

**Expected (when .env incomplete):**
- Lists missing required variables
- Explains why each is required (with feature context)
- Provides exact set commands
- Shows how to fix placeholder values
- Re-validates after fixes

**Test 6 - Update existing variable:**
```
Change SESSION_SECRET
```

**Expected:**
- Runs `lvt env set SESSION_SECRET $(openssl rand -hex 32)`
- Shows old value (masked) and new value (masked)
- Confirms update
- Maintains other variables

**Test 7 - Remove variable:**
```
Remove SMTP_HOST
```

**Expected:**
- Runs `lvt env unset SMTP_HOST`
- Confirms removal
- Shows updated list without that variable

**Test 8 - Security check:**
```
Check if my secrets are secure
```

**Expected:**
- Runs `lvt env validate --strict`
- Validates secret length (32+ chars)
- Detects "change-me" or other placeholders
- Confirms .env is in .gitignore
- Checks file permissions (0600)

---

### 15. lvt:quickstart

**Test 1 - Blog quickstart:**
```
Quickstart a blog app
```

**Expected:**
- Uses lvt:quickstart skill
- Creates app
- Adds posts resource
- Adds comments resource
- Seeds data
- Starts server
- Time: < 3 minutes

**Test 2 - Todo quickstart:**
```
Build a todo app quickly
```

**Expected:**
- Detects "todo" domain
- Creates app
- Adds tasks resource
- Seeds tasks
- Runs server

**Test 3 - Custom quickstart:**
```
Quickstart an e-commerce shop
```

**Expected:**
- Detects "shop" domain
- Suggests products + orders
- Creates both resources
- Sets up relationships

---

### 16. lvt:production-ready

**Test 1 - Full production:**
```
Make my blog production ready
```

**Expected:**
- Uses lvt:production-ready skill
- Adds authentication
- Generates deployment config
- Creates .env.example
- Provides production checklist
- Time: < 10 minutes

**Test 2 - Deployment target:**
```
Prepare for Docker deployment
```

**Expected:**
- Adds auth
- Creates Dockerfile
- Explains environment setup
- Security checklist

---

### 17. lvt:add-related-resources

**Test 1 - Blog suggestions:**
```
I have posts, what should I add next?
```

**Expected:**
- Uses lvt:add-related-resources skill
- Detects blog domain
- Suggests: comments, categories, tags
- Explains relationships
- Provides commands

**Test 2 - E-commerce suggestions:**
```
Suggest resources for my shop
```

**Expected:**
- Detects e-commerce from "shop"
- Suggests: orders, order_items, reviews, cart
- Prioritizes (essential vs nice-to-have)
- Shows junction tables

**Test 3 - Domain detection:**
```
What goes well with a tasks resource?
```

**Expected:**
- Detects project management
- Suggests: projects, labels, time_entries
- Shows many-to-many patterns

---

## Maintenance & Support Skills (3 skills)

### 18. lvt:analyze

**Test 1 - App analysis:**
```
Analyze my blog app
```

**Expected:**
- Uses lvt:analyze skill
- Runs lvt resource list
- Shows schema summary
- Calculates complexity
- Detects relationships
- Lists features (auth, CRUD, etc.)
- Provides recommendations

**Test 2 - Schema review:**
```
Review my database schema
```

**Expected:**
- Analyzes tables
- Checks indexes
- Reviews foreign keys
- Suggests improvements

---

### 19. lvt:suggest

**Test 1 - Improvement suggestions:**
```
Suggest improvements for my app
```

**Expected:**
- Uses lvt:suggest skill
- Analyzes current state
- Provides P0/P1/P2 recommendations
- Shows specific commands
- Estimates time and impact

**Test 2 - Missing features:**
```
What's missing from my blog?
```

**Expected:**
- Detects domain (blog)
- Suggests: search, categories, comments
- Prioritizes by importance
- Shows implementation

---

### 20. lvt:troubleshoot

**Test 1 - Build error:**
```
My app won't build, getting "undefined: queries"
```

**Expected:**
- Uses lvt:troubleshoot skill
- Diagnoses: missing sqlc generate
- Provides fix
- Verifies success

**Test 2 - Database error:**
```
Getting "database is locked"
```

**Expected:**
- Explains SQLite limitation
- Shows WAL mode solution
- Provides db.SetMaxOpenConns(1)

**Test 3 - Template error:**
```
Template error: "can't evaluate field Title"
```

**Expected:**
- Checks field names
- Compares to models.go
- Shows case sensitivity

**Test 4 - Auth error:**
```
Users can't log in, password always fails
```

**Expected:**
- Checks bcrypt implementation
- Verifies hash storage
- Shows correct comparison

**Test 5 - Deployment error:**
```
Port 8080 already in use
```

**Expected:**
- Shows lsof command
- Kill process option
- Change port option

---

## Meta Skills (1 skill)

### 21. lvt:add-skill

**Test:**
```
How do I create a skill for lvt template copy?
```

**Expected:**
- Uses lvt:add-skill skill
- Explains TDD methodology
- Shows template structure
- Provides quality checklist

---

## Phase 6: CLI Enhancements

### lvt env generate

**Test 1 - Generate env file:**
```
Generate environment configuration for my app
```

**Expected:**
- Detects features (auth, database, email)
- Runs `lvt env generate`
- Creates .env.example
- Shows all sections
- Explains security practices

**Test 2 - Auth app env:**
```
Create .env.example for my app with authentication
```

**Expected:**
- Includes SESSION_SECRET
- Includes CSRF_SECRET
- Includes SMTP config
- Shows secret generation command

---

### Production-Ready Templates

**Test:**
```
Create a new app and verify production features
```

**Expected:**
- Generates app with enhanced main.go
- Shows structured logging
- Shows security headers
- Shows graceful shutdown
- Shows health endpoint
- Verifies /health returns 200

---

## Integration Tests (Multi-Skill Workflows)

### Workflow 1: Complete Blog (9 skills)

**Prompt:**
```
Walk me through creating a complete production blog app from scratch
```

**Expected skills used in order:**
1. lvt:new-app
2. lvt:add-resource (posts)
3. lvt:add-resource (comments)
4. lvt:gen-auth
5. lvt:seed-data
6. lvt:run-and-test
7. lvt env generate
8. lvt:deploy
9. lvt:analyze

**Verification:**
- App created and running
- Auth working
- Test data present
- .env.example generated
- Deployment config created
- Analysis shows complete app

---

### Workflow 2: Rapid Prototype (Quickstart)

**Prompt:**
```
Quickstart a todo app
```

**Expected skills:**
- lvt:quickstart (chains new-app, add-resource, seed-data, run-and-test)

**Verification:**
- App running in < 3 minutes
- Tasks resource exists
- Test data loaded
- Server accessible

---

### Workflow 3: Optimization Pass

**Prompt:**
```
Analyze my app, suggest improvements, and add recommended features
```

**Expected skills:**
1. lvt:analyze
2. lvt:suggest
3. lvt:add-related-resources
4. Implements suggestions

**Verification:**
- Analysis complete
- Prioritized suggestions
- Domain-appropriate recommendations
- Implementation guidance

---

### Workflow 4: Debug and Fix

**Prompt:**
```
My app has errors and won't build
```

**Expected skills:**
1. lvt:troubleshoot (diagnose)
2. lvt:validate-templates (if template issue)
3. lvt:resource-inspect (if schema issue)
4. Provides fixes

**Verification:**
- Issue diagnosed
- Fix provided
- Build successful

---

## Edge Cases & Error Scenarios

### Test: Wrong directory
```
I'm getting "go.mod not found"
```

**Expected:**
- Skills check prerequisites
- Explain need to be in app root
- Show how to verify location

---

### Test: Missing dependencies
```
Getting "sqlc: command not found"
```

**Expected:**
- Explains sqlc installation
- Shows go install command
- Verifies PATH

---

### Test: Database migration conflict
```
Migration failed: "table already exists"
```

**Expected:**
- Check lvt migration status
- Explains migration state
- Shows rollback if needed

---

### Test: Template syntax error
```
Template has {{ without closing
```

**Expected:**
- Runs lvt parse
- Shows exact line with error
- Suggests fix

---

## Negative Tests (Skills Should NOT Be Used)

### Test: General Go question
```
How do I create a struct in Go?
```

**Expected:**
- Does NOT use lvt skills
- Answers directly (general Go knowledge)

---

### Test: CSS question
```
How do I center a div?
```

**Expected:**
- Does NOT use lvt skills
- Answers with CSS knowledge

---

### Test: Different framework
```
How do I deploy a React app?
```

**Expected:**
- Does NOT use lvt skills
- General deployment advice

---

## Performance Benchmarks

**Metric:** Time to working app
**Target:** < 3 minutes with quickstart
**Test:** "Quickstart a blog app"

**Metric:** Time to production
**Target:** < 15 minutes with production-ready
**Test:** "Make my app production ready"

**Metric:** Build success rate
**Target:** 100% (skills prevent errors)
**Test:** Follow any core skill workflow

---

## Quality Checklist

After testing each skill, verify:

- ✅ Skill is mentioned by name
- ✅ Commands match actual lvt CLI
- ✅ File paths are accurate
- ✅ Error messages match real output
- ✅ Examples are copy-pasteable
- ✅ Common issues are prevented
- ✅ Cross-references work
- ✅ Prerequisites are checked
- ✅ Best practices included
- ✅ Workflows complete successfully

---

## Reporting Issues

If a skill test fails:

1. **Document the failure:**
   - Prompt used
   - Expected behavior
   - Actual behavior
   - Error messages

2. **Check skill file:**
   - Is information outdated?
   - Is example incorrect?
   - Is guidance incomplete?

3. **Report:**
   - Note in tracker
   - Update skill if needed
   - Re-test after fix

---

## Test Coverage Summary

**Core Commands:** 13/13 skills tested
**Workflow Orchestration:** 3/3 skills tested
**Maintenance & Support:** 3/3 skills tested
**Meta Skills:** 1/1 skills tested
**Phase 6 Enhancements:** 2/2 tested
**Integration Workflows:** 4 multi-skill scenarios
**Edge Cases:** 4 scenarios
**Negative Tests:** 3 scenarios

**Total Test Scenarios:** 60+

---

**Testing Status:** Comprehensive coverage for all 19 skills + Phase 6 ✅
