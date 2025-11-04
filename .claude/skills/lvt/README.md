# LiveTemplate Skills for Claude Code

Comprehensive skills for building LiveTemplate applications with the lvt CLI.

## What Are These Skills?

These skills help Claude Code provide expert guidance on LiveTemplate development. When you ask Claude about lvt commands or workflows, it automatically uses the relevant skill to give you accurate, comprehensive answers.

**You don't need to do anything special** - just ask Claude about your LiveTemplate needs, and it will use the appropriate skill automatically.

## Available Skills

### Core Commands

#### 📦 lvt:new-app
**When to use:** Creating a new LiveTemplate application

**Ask Claude:**
- "Create a new app called myapp"
- "I want to start a new LiveTemplate project"
- "Which CSS framework should I use for my app?"

**What you'll get:** Guidance on `lvt new`, kit selection, project structure verification

---

#### 🗄️ lvt:add-resource
**When to use:** Adding database-backed CRUD resources

**Ask Claude:**
- "Add a products resource with name and price"
- "I need a users table with email and password"
- "How do I specify field types like integers and booleans?"

**What you'll get:** Field type mapping, database generation, migration running, route registration

---

#### 📄 lvt:add-view
**When to use:** Adding UI pages without database (dashboards, analytics, static pages)

**Ask Claude:**
- "Add a dashboard view to my app"
- "Create an analytics page"
- "What's the difference between views and resources?"

**What you'll get:** View generation, understanding when to use views vs resources, no-database patterns

---

#### 🔧 lvt:add-migration
**When to use:** Managing database schema changes (indexes, constraints, data fixes)

**Ask Claude:**
- "Add an index on the products price field"
- "Create a unique constraint on user emails"
- "I need to add a check constraint"

**What you'll get:** Migration creation, goose format, Up/Down SQL, migration workflow

---

#### ▶️ lvt:run-and-test
**When to use:** Running your app locally and testing it

**Ask Claude:**
- "How do I run my app?"
- "My tests are failing"
- "Port 8080 is already in use"

**What you'll get:** Development server usage, testing guidance, troubleshooting common issues

---

#### ✏️ lvt:customize
**When to use:** Modifying generated code (handlers, templates, queries)

**Ask Claude:**
- "Add filtering to my products list"
- "Change the HTML structure of my page"
- "Add a real-time delete button"

**What you'll get:** Handler customization, template editing, custom queries, WebSocket actions

---

#### 🌱 lvt:seed-data
**When to use:** Generating realistic test data for development

**Ask Claude:**
- "I need test data for my products"
- "Generate 100 users with realistic data"
- "How do I clean up test data?"

**What you'll get:** Seeding with --count, cleanup strategies, context-aware generation, test record management

---

#### 🚀 lvt:deploy
**When to use:** Deploying your app to production

**Ask Claude:**
- "Deploy my app to production"
- "Create a Dockerfile for my app"
- "How do I deploy to Fly.io?"

**What you'll get:** Docker setup, Fly.io deployment, Kubernetes config, VPS deployment, production best practices

---

### Meta Skills

#### 🎓 lvt:add-skill
**When to use:** Creating new skills for lvt commands

**Ask Claude:**
- "I want to create a skill for lvt template copy"
- "How do I write high-quality skills?"
- "What's the TDD process for documentation?"

**What you'll get:** Skill creation methodology, template structure, testing approach, quality checklist

---

## How Skills Work

### Automatic Activation

Skills are automatically activated when you ask relevant questions:

```
You: "Create a new app called blog"
Claude: "I'm using the lvt:new-app skill to help you create a new LiveTemplate application..."
```

### Multi-Skill Workflows

Claude uses multiple skills for complex tasks:

```
You: "Build a task manager app and deploy it"
Claude uses:
1. lvt:new-app → Create application
2. lvt:add-resource → Add tasks resource
3. lvt:seed-data → Generate test data
4. lvt:run-and-test → Verify locally
5. lvt:deploy → Deploy to production
```

### Preventing Common Mistakes

Skills prevent typical errors:

❌ **Without skill:** Forgets to run migrations, uses wrong field types, skips prerequisites
✅ **With skill:** Checks prerequisites, uses correct types, runs migrations automatically

---

## Quick Start Examples

### Example 1: Build a Product Catalog

```
You: "I want to build a product catalog app with products and categories"

Claude will:
1. Create app with lvt:new-app
2. Add resources with lvt:add-resource
3. Set up relationships
4. Generate test data with lvt:seed-data
5. Show how to run with lvt:run-and-test
```

### Example 2: Add Dashboard Analytics

```
You: "Add an analytics dashboard to my existing app"

Claude will:
1. Use lvt:add-view (not resource - no database)
2. Show how to customize template
3. Add real-time updates with WebSocket
```

### Example 3: Deploy to Production

```
You: "Deploy my app to the cloud"

Claude will:
1. Use lvt:deploy to recommend Fly.io
2. Create Dockerfile if needed
3. Set up volumes for SQLite
4. Handle migrations in production
5. Configure backups
```

---

## Skill Categories

### By Development Phase

**Starting Out:**
- lvt:new-app → Create project

**Building Features:**
- lvt:add-resource → Database-backed features
- lvt:add-view → UI-only pages
- lvt:add-migration → Schema changes
- lvt:customize → Modify generated code

**Development:**
- lvt:run-and-test → Local testing
- lvt:seed-data → Test data

**Production:**
- lvt:deploy → Go live

**Meta:**
- lvt:add-skill → Create more skills

---

## Common Workflows

### Full CRUD Application

1. **Create app:** "Create new app called myapp"
2. **Add resource:** "Add users resource with name and email"
3. **Test data:** "Generate 50 test users"
4. **Run locally:** "Start the development server"
5. **Customize:** "Add email validation to users"
6. **Deploy:** "Deploy to Fly.io"

### Dashboard Application

1. **Create app:** "Create analytics app"
2. **Add views:** "Add dashboard, reports, and settings views"
3. **Customize:** "Add real-time metrics to dashboard"
4. **Deploy:** "Containerize with Docker"

### Adding Features to Existing App

1. **Add resource:** "Add comments to my blog app"
2. **Migration:** "Add index on comment timestamps"
3. **Test data:** "Generate 200 test comments"
4. **Customize:** "Add comment voting functionality"

---

## Tips for Working with Skills

### 1. Be Specific

❌ **Vague:** "I need help with my app"
✅ **Specific:** "Add a products resource with name, price, and description"

### 2. Mention Your Context

Good to include:
- Where you are in the project (new, existing, deploying)
- What you've already done
- Errors you're seeing

Example: "I created a products resource, ran migrations, but now getting 'database locked' error"

### 3. Ask Follow-Up Questions

Skills are comprehensive - ask for clarification:
- "Why did you use float instead of int for price?"
- "What's the difference between --cleanup and just deleting the database?"
- "Can you explain the CGO_ENABLED requirement?"

### 4. Let Claude Choose Skills

You don't need to name skills - Claude picks the right ones:

```
You: "I want to add test data"
Claude: (automatically uses lvt:seed-data)

You: "Deploy my app"
Claude: (automatically uses lvt:deploy)
```

---

## Troubleshooting with Skills

### Common Issues

**"I'm getting an error..."**
→ Skills have "Common Issues" sections with fixes

**"How do I...?"**
→ Skills have "Quick Reference" tables

**"What's the difference between...?"**
→ Skills explain concepts and trade-offs

**"This isn't working..."**
→ Skills check prerequisites and common mistakes

---

## What's NOT Covered

These skills are specifically for **lvt CLI and LiveTemplate**. For other topics:

- **General Go questions** → Ask Claude directly (no skill needed)
- **CSS/JavaScript help** → Ask Claude directly
- **Database design theory** → Ask Claude directly
- **Other frameworks** → Not covered by lvt skills

---

## Skill Quality Guarantees

All skills follow these standards:

✓ **Tested** - Verified against actual lvt commands
✓ **Complete** - Cover all command flags and options
✓ **Accurate** - Error messages match real output
✓ **Example-rich** - Copy-paste ready code
✓ **Error-aware** - Common mistakes documented with fixes
✓ **Cross-referenced** - Link to related skills

---

## Getting Help

### In Claude Code

Just ask naturally:
```
"How do I create a new app?"
"Add a products resource"
"Deploy to production"
"I'm getting this error: [paste error]"
```

### Skill Coverage

**Covered:** All lvt commands, workflows, deployment
**Not covered yet:** Advanced topics (coming soon)

### Feedback

If you find:
- Missing information
- Incorrect guidance
- Unclear explanations

Let Claude know! Skills can be updated based on feedback.

---

## Quick Reference

**I want to...** | **Ask Claude** | **Skill Used**
---|---|---
Start a new project | "Create a new app" | lvt:new-app
Add database tables | "Add products resource" | lvt:add-resource
Add UI pages | "Add dashboard view" | lvt:add-view
Modify database | "Add an index" | lvt:add-migration
Run my app | "Start the server" | lvt:run-and-test
Change generated code | "Add filtering" | lvt:customize
Generate test data | "Create 100 users" | lvt:seed-data
Go to production | "Deploy my app" | lvt:deploy
Create new skills | "Make a skill for..." | lvt:add-skill

---

## Next Steps

1. **Start building:** "Create a new LiveTemplate app called [name]"
2. **Learn by doing:** Ask Claude for help at each step
3. **Explore skills:** Try different commands and see skills in action
4. **Trust the process:** Skills prevent common mistakes automatically

**Remember:** You don't need to know these skills exist - Claude uses them automatically when you ask LiveTemplate questions!

---

## Skill Files

Location: `~/.claude/skills/lvt/`

```
lvt/
├── core/              ← Command skills
│   ├── new-app.md
│   ├── add-resource.md
│   ├── add-view.md
│   ├── add-migration.md
│   ├── run-and-test.md
│   ├── customize.md
│   ├── seed-data.md
│   └── deploy.md
├── meta/              ← Skills about skills
│   └── add-skill.md
├── README.md          ← This file
└── TESTING.md         ← Testing guide
```

---

**Happy building with LiveTemplate!** 🚀
