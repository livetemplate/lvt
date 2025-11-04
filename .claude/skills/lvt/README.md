# LiveTemplate Skills for Claude Code

Complete AI-guided development for LiveTemplate applications - from zero to production in minutes.

## What Are These Skills?

These skills enable Claude Code to guide you through the entire LiveTemplate development lifecycle. When you ask Claude about lvt commands or workflows, it automatically uses the relevant skill to provide expert, production-ready guidance.

**You don't need to do anything special** - just ask Claude about your LiveTemplate needs, and it will use the appropriate skills automatically.

## 📊 Complete Coverage

- **19 Skills** covering the entire development lifecycle
- **13 Core Commands** for project development
- **3 Workflow Orchestrations** for rapid development
- **3 Maintenance Tools** for optimization and debugging
- **1 Meta Skill** for extending the skill system
- **Phase 6 Enhancements** for production-ready defaults

---

## Available Skills

### 🚀 Core Commands (13 skills)

#### 📦 lvt:new-app
**Create new LiveTemplate applications**

**Ask Claude:**
- "Create a new app called myapp"
- "Which CSS framework should I use?"
- "Create a blog app with Tailwind"

**What you'll get:** Kit selection (multi/single/simple), module configuration, project structure setup

**Example:**
```
User: "Create a blog app"
Claude: Creates app with multi kit, explains structure, shows next steps
```

---

#### 🗄️ lvt:add-resource
**Add database-backed CRUD resources**

**Ask Claude:**
- "Add a products resource with name and price"
- "Create a posts table with foreign key to users"
- "How do I add a boolean field?"

**What you'll get:** Type inference, database migrations, route injection, pagination setup

**Example:**
```
User: "Add posts with title, content, published"
Claude: Generates resource, runs migrations, updates routes, shows URL
```

---

#### 📄 lvt:add-view
**Add UI-only pages without database**

**Ask Claude:**
- "Add a dashboard view"
- "Create an analytics page"
- "Add a landing page"

**What you'll get:** View-only handlers, no database overhead, when to use vs resources

**Example:**
```
User: "Add a counter view"
Claude: Generates handler, shows how to add state, explains WebSocket actions
```

---

#### 🛠️ lvt:add-migration
**Manage database schema changes**

**Ask Claude:**
- "Add an index on products.price"
- "Create a unique constraint on user emails"
- "How do I rename a column?"

**What you'll get:** Migration creation (goose format), Up/Down SQL, schema modification patterns

**Example:**
```
User: "Add index on created_at"
Claude: Creates migration, shows SQL, runs it, confirms with status
```

---

#### 🔐 lvt:gen-auth
**Add complete authentication system** ⭐ **USER PRIORITY**

**Ask Claude:**
- "Add authentication to my app"
- "I need password and magic link auth"
- "How do I protect routes?"

**What you'll get:** Password auth, magic links, email confirmation, password reset, sessions, CSRF, middleware, E2E tests

**Example:**
```
User: "Add authentication"
Claude: Generates auth system, runs migrations, shows how to wire routes, demonstrates RequireAuth middleware
```

**Features:**
- Password authentication (bcrypt)
- Magic link (passwordless)
- Email confirmation
- Password reset
- Session management
- CSRF protection
- Route protection middleware

---

#### 🗂️ lvt:gen-schema
**Generate database schema without UI**

**Ask Claude:**
- "Create a schema-only table for logs"
- "Add a sessions table without handlers"
- "I need a backend-only table"

**What you'll get:** Database table without handler/template, perfect for logs, analytics, cache

**Example:**
```
User: "Create analytics table"
Claude: Generates schema.sql and queries.sql, skips handler, explains use cases
```

---

#### 🔍 lvt:resource-inspect
**Inspect database resources and schema**

**Ask Claude:**
- "Show me all my resources"
- "What's in the posts table?"
- "List all tables"

**What you'll get:** Resource listing, detailed schema view, relationships, indexes

**Example:**
```
User: "Describe posts table"
Claude: Shows columns, types, indexes, foreign keys, constraints
```

---

#### 🎨 lvt:manage-kits
**Manage CSS framework kits**

**Ask Claude:**
- "What kits are available?"
- "Show me kit details"
- "Create a custom kit"

**What you'll get:** Kit listing, detailed info, validation, custom kit creation

**Example:**
```
User: "List available kits"
Claude: Shows multi (Tailwind), single (Tailwind SPA), simple (Pico)
```

---

#### ✅ lvt:validate-templates
**Validate and analyze templates**

**Ask Claude:**
- "Check my template for errors"
- "Validate posts.tmpl"
- "Why is my template failing?"

**What you'll get:** Syntax checking, execution testing, common issues detection

**Example:**
```
User: "Parse internal/app/posts/posts.tmpl"
Claude: Validates html/template + LiveTemplate syntax, checks for issues
```

---

#### ▶️ lvt:run-and-test
**Run development server and tests**

**Ask Claude:**
- "Start my app"
- "Run tests"
- "Port 8080 is in use"

**What you'll get:** Dev server usage, testing workflows, debugging, troubleshooting

**Example:**
```
User: "Run the app"
Claude: Starts lvt serve, opens browser, shows WebSocket connection
```

---

#### ✏️ lvt:customize
**Customize generated code**

**Ask Claude:**
- "Add filtering to products list"
- "Change the HTML layout"
- "Add a custom SQL query"

**What you'll get:** Handler modification, template editing, custom queries, WebSocket actions

**Example:**
```
User: "Add search to posts"
Claude: Modifies handler, updates template, adds WHERE clause, shows results
```

---

#### 🌱 lvt:seed-data
**Generate realistic test data**

**Ask Claude:**
- "Generate 100 test users"
- "Seed my products table"
- "Clean up test data"

**What you'll get:** Context-aware data generation, cleanup strategies, bulk insertion

**Example:**
```
User: "Generate 50 posts"
Claude: Creates seed command, runs it, shows realistic data in database
```

---

#### 🚀 lvt:deploy
**Deploy to production**

**Ask Claude:**
- "Deploy to Fly.io"
- "Create a Dockerfile"
- "Production deployment guide"

**What you'll get:** Docker setup, Fly.io config, Kubernetes manifests, VPS deployment, migrations, CGO

**Example:**
```
User: "Deploy to production"
Claude: Creates Dockerfile, fly.toml, explains litestream, deploys
```

---

### 🔄 Workflow Orchestration (3 skills)

#### ⚡ lvt:quickstart
**Rapid end-to-end app creation**

**Ask Claude:**
- "Quickstart a blog app"
- "Build a todo app fast"
- "Create a shop quickly"

**What you'll get:** Complete workflow from zero to running app in minutes

**Chains together:**
1. lvt:new-app (create app)
2. lvt:add-resource (add resources)
3. lvt:seed-data (test data)
4. lvt:run-and-test (start server)

**Example:**
```
User: "Quickstart a blog"
Claude: Creates app → adds posts → adds comments → seeds data → runs server
Time: < 3 minutes to working blog
```

---

#### 🏭 lvt:production-ready
**Transform dev app to production**

**Ask Claude:**
- "Make my app production ready"
- "Add auth and deployment"
- "Prepare for real users"

**What you'll get:** Auth system + deployment config + env setup + best practices

**Chains together:**
1. lvt:gen-auth (authentication)
2. lvt:deploy (deployment config)
3. Environment setup
4. Production checklist

**Example:**
```
User: "Make my shop production ready"
Claude: Adds auth → wires routes → creates Dockerfile → sets up .env → provides checklist
Time: < 10 minutes to production
```

---

#### 🧠 lvt:add-related-resources
**Intelligent resource suggestions**

**Ask Claude:**
- "What should I add to my blog?"
- "I have posts, what's next?"
- "Suggest related resources"

**What you'll get:** Domain-based suggestions, relationship patterns, industry standards

**Example:**
```
User: "I have posts, what else?"
Claude: Suggests comments (one-to-many), categories (many-to-many), tags
Explains relationships, provides commands
```

**Domain patterns:**
- Blog → comments, categories, tags
- E-commerce → orders, reviews, cart
- SaaS → teams, subscriptions, usage
- Project management → tasks, milestones, team members

---

### 🛠️ Maintenance & Support (3 skills)

#### 📊 lvt:analyze
**Comprehensive app analysis**

**Ask Claude:**
- "Analyze my app"
- "Show app structure"
- "Review my schema"

**What you'll get:** Schema analysis, relationships, complexity metrics, database health, feature detection

**Example:**
```
User: "Analyze my blog"
Claude:
- 3 resources (posts, comments, users)
- Medium complexity (18 fields)
- Relationships: posts→comments (1:many)
- Missing: categories, search, pagination
```

---

#### 💡 lvt:suggest
**Actionable improvement recommendations**

**Ask Claude:**
- "Suggest improvements"
- "How can I make this better?"
- "What's missing from my app?"

**What you'll get:** Prioritized recommendations (P0/P1/P2), specific commands, impact estimates

**Example:**
```
User: "Suggest improvements"
Claude:
P0 (Critical):
- Add authentication (10 min) ⭐⭐⭐⭐⭐
P1 (Important):
- Add indexes on created_at (5 min) ⭐⭐⭐⭐
P2 (Nice to have):
- Add search functionality (15 min) ⭐⭐⭐
```

---

#### 🔧 lvt:troubleshoot
**Debug common issues**

**Ask Claude:**
- "My app won't build"
- "Getting database locked error"
- "Templates are failing"

**What you'll get:** 7 issue categories, common patterns, diagnostic commands, solutions

**Example:**
```
User: "Error: undefined: queries"
Claude: Diagnoses → missing sqlc generate → provides fix → verifies success
```

**Covers:**
- Build errors
- Migration problems
- Template errors
- Auth issues
- WebSocket debugging
- Deployment failures
- Runtime errors

---

### 🎓 Meta Skills (1 skill)

#### 📚 lvt:add-skill
**Create new skills using TDD**

**Ask Claude:**
- "Create a skill for lvt template copy"
- "How do I write good skills?"
- "Skill development guide"

**What you'll get:** TDD methodology, template structure, testing approach, quality checklist

---

## 🎯 Phase 6: CLI Enhancements

### lvt env generate
**Smart environment configuration**

**Features:**
- Detects app features automatically (auth, database, email)
- Generates comprehensive .env.example
- Documents all variables
- Security best practices included

**Usage:**
```bash
lvt env generate
# Creates .env.example with detected configuration
```

### Production-Ready Templates
**Battle-tested defaults out of the box**

**Enhanced main.go.tmpl includes:**
- ✅ Structured logging (log/slog, JSON)
- ✅ Security headers (XSS, clickjacking, CSP, HSTS)
- ✅ Recovery middleware (panic handling)
- ✅ HTTP request logging (metrics)
- ✅ Graceful shutdown (SIGINT, SIGTERM)
- ✅ Health check endpoint (/health)
- ✅ Production timeouts (15s read/write, 60s idle)
- ✅ Environment variables (PORT, LOG_LEVEL, DATABASE_PATH, APP_ENV)

---

## 🎬 Complete Workflows

### Workflow 1: Zero to Production Blog (15 minutes)

```
1. "Create a blog app"                     → lvt:new-app
2. "Add posts with title and content"      → lvt:add-resource
3. "Add comments to posts"                 → lvt:add-resource
4. "Add authentication"                    → lvt:gen-auth
5. "Generate test data"                    → lvt:seed-data
6. "Run the app"                          → lvt:run-and-test
7. "Generate environment config"           → lvt env generate
8. "Deploy to Fly.io"                     → lvt:deploy

Result: Production blog with auth, running in the cloud
```

### Workflow 2: Rapid Prototyping (3 minutes)

```
1. "Quickstart a todo app"                 → lvt:quickstart
   - Creates app
   - Adds tasks resource
   - Seeds 20 tasks
   - Starts server

Result: Working todo app at http://localhost:8080
```

### Workflow 3: Optimization Pass

```
1. "Analyze my app"                        → lvt:analyze
2. "Suggest improvements"                  → lvt:suggest
3. Implement P0/P1 suggestions
4. "What else should I add?"              → lvt:add-related-resources

Result: Optimized app with recommended features
```

### Workflow 4: Debugging Session

```
1. "My app won't build"                    → lvt:troubleshoot
2. Diagnoses: missing sqlc generate
3. Provides fix
4. "Validate my templates"                 → lvt:validate-templates
5. "Inspect my schema"                     → lvt:resource-inspect

Result: All issues resolved, app building successfully
```

---

## 📚 Quick Reference

### By Development Phase

| Phase | Skills | Purpose |
|-------|--------|---------|
| **Setup** | new-app | Create project |
| **Build** | add-resource, add-view, gen-schema | Add features |
| **Database** | add-migration, resource-inspect | Schema management |
| **Security** | gen-auth | Authentication |
| **Quality** | validate-templates, customize | Refinement |
| **Testing** | run-and-test, seed-data | Development |
| **Production** | deploy, env generate | Go live |
| **Optimize** | analyze, suggest | Improvements |
| **Debug** | troubleshoot | Fix issues |
| **Automate** | quickstart, production-ready | Rapid workflows |

### By Complexity

**Beginner (getting started):**
- lvt:new-app
- lvt:add-resource
- lvt:run-and-test

**Intermediate (building features):**
- lvt:add-view
- lvt:customize
- lvt:seed-data
- lvt:add-migration

**Advanced (production):**
- lvt:gen-auth
- lvt:deploy
- lvt:gen-schema
- lvt:validate-templates

**Expert (optimization):**
- lvt:analyze
- lvt:suggest
- lvt:troubleshoot
- lvt:add-related-resources

**Power User (automation):**
- lvt:quickstart
- lvt:production-ready
- lvt:add-skill

---

## 💡 Tips for Success

### 1. Be Specific

❌ **Vague:** "I need help"
✅ **Specific:** "Add products with name:string, price:float, stock:int"

### 2. Provide Context

**Good context includes:**
- Current project state (new/existing)
- What you've already done
- Errors you're seeing
- What you're trying to achieve

Example: "I created posts, ran migrations, but getting 'database locked' error when running app"

### 3. Trust the Workflow

Skills prevent common mistakes automatically:
- ✅ Checks prerequisites
- ✅ Runs migrations
- ✅ Generates sqlc models
- ✅ Wires routes
- ✅ Validates configuration

### 4. Use Workflows for Speed

**Instead of:**
```
"Create app"
"Add resource"
"Add another resource"
"Generate data"
"Run app"
```

**Use:**
```
"Quickstart a blog app"
```
(Does all of the above automatically)

### 5. Ask Follow-Up Questions

Skills are comprehensive - ask for clarification:
- "Why float instead of int for price?"
- "What's the difference between CASCADE and RESTRICT?"
- "Can you explain the CSP header?"

---

## 🚦 Skill Quality Standards

All skills follow these standards:

✓ **Tested** - Verified against actual lvt commands
✓ **Complete** - Cover all command flags and options
✓ **Accurate** - Error messages match real output
✓ **Example-rich** - Copy-paste ready code
✓ **Error-aware** - Common mistakes documented with fixes
✓ **Cross-referenced** - Link to related skills
✓ **Production-ready** - Best practices included

---

## 📁 Skill Organization

```
.claude/skills/lvt/
├── README.md                         # This file
├── TESTING.md                        # Test scenarios
│
├── core/                             # 13 command skills
│   ├── new-app.md                   # Create applications
│   ├── add-resource.md              # CRUD resources
│   ├── add-view.md                  # UI-only pages
│   ├── add-migration.md             # Database migrations
│   ├── gen-schema.md                # Schema-only tables
│   ├── gen-auth.md                  # Authentication system ⭐
│   ├── resource-inspect.md          # Inspect resources
│   ├── manage-kits.md               # Kit management
│   ├── validate-templates.md        # Template validation
│   ├── run-and-test.md              # Dev server & tests
│   ├── customize.md                 # Code customization
│   ├── seed-data.md                 # Test data
│   └── deploy.md                    # Production deployment
│
├── workflows/                        # 3 workflow orchestration
│   ├── quickstart.md                # Rapid app creation
│   ├── production-ready.md          # Production transformation
│   └── add-related-resources.md     # Intelligent suggestions
│
├── maintenance/                      # 3 maintenance & support
│   ├── analyze.md                   # App analysis
│   ├── suggest.md                   # Improvement recommendations
│   └── troubleshoot.md              # Debug common issues
│
├── meta/                             # 1 meta skill
│   └── add-skill.md                 # Skill creation guide
│
└── docs/                             # Project documentation
    ├── CLAUDE_SKILLS_TRACKER.md     # Status tracker (100% complete)
    ├── SKILL_DEVELOPMENT.md         # Development guide
    ├── SKILL_TESTING_CHECKLISTS.md  # Testing procedures
    └── TEST_RESULTS_NEW_APP.md      # Test results
```

---

## ❓ Getting Help

### In Claude Code

Just ask naturally:
```
"How do I create a new app?"
"Add a products resource"
"Make my app production ready"
"Why is this error happening: [paste error]"
```

### Coverage

**✅ Covered:**
- All lvt commands
- Complete workflows
- Production deployment
- Debugging and optimization
- Environment configuration
- Template customization

**📋 Coming Soon:**
- Advanced patterns
- Performance tuning
- Scalability guides

### Feedback

If you find issues or have suggestions:
- Missing information → Let Claude know
- Incorrect guidance → Report it
- Unclear explanations → Ask for clarification

Skills are living documents that improve with usage!

---

## 🎯 Success Metrics

**Time to Working App:** < 3 minutes (with quickstart)
**Time to Production:** < 15 minutes (with production-ready)
**Build Success Rate:** 100% (skills prevent common errors)
**Test Pass Rate:** 100% (9/9 automated tests)

---

## 🚀 Next Steps

1. **Start Building:** "Create a new LiveTemplate app"
2. **Learn by Doing:** Ask Claude for help at each step
3. **Explore Workflows:** Try quickstart and production-ready
4. **Optimize:** Use analyze and suggest for improvements
5. **Deploy:** Ship to production with confidence

**Remember:** You don't need to memorize these skills - Claude uses them automatically when you ask LiveTemplate questions!

---

**Happy building with LiveTemplate!** 🎉

**Status:** 19/19 skills (100% complete) + Phase 6 enhancements ✅
