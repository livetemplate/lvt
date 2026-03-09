# LiveTemplate AI Agent Usage Guide

This guide explains how to use LiveTemplate with AI assistants to create and develop Go web applications.

## Overview

LiveTemplate supports **5 different AI assistants** via agent installation:

**Supported AI Assistants:**
- Claude Code
- GitHub Copilot (VS Code)
- Cursor AI
- Aider CLI
- Any LLM (Generic)

---

## Quick Start

### Installation

```bash
# Install project-specific agent
lvt install-agent                    # Interactive menu (recommended)
lvt install-agent --llm <type>       # Direct installation
```

---

## LLM-Specific Guides

### Claude Desktop / Claude Code

**Best For:** Full-featured development with autonomous skills that activate based on your requests

**Setup:**
```bash
# Install Claude agent and skills in your project
lvt install-agent              # Interactive menu - select Claude Code
lvt install-agent --llm claude # Direct installation
```

**What You Get:**
- `.claude/skills/` - 22 autonomous skills (each in its own directory)
- `.claude/agents/lvt-assistant.md` - LiveTemplate specialist agent
- `.claude/settings.json` - Configuration and permissions

**Understanding Skills vs Agents:**

**Skills** are autonomous capabilities that Claude uses automatically:
- **NOT invoked with slash commands** - they activate based on your request
- Claude reads skill descriptions and decides when to use them
- Work seamlessly in the background
- Listed when you ask: "What skills are available?"

**Agents** are specialists you can explicitly invoke:
- Use via `/agents` picker or by mentioning them
- Example: Select "lvt-assistant" from the agents menu
- Have their own context window for complex tasks

**Skills Available (22 total):**

**Core Skills (14):**
- `lvt-new-app` - Create new application
- `lvt-add-resource` - Add CRUD resource with database
- `lvt-add-view` - Add standalone view (no database)
- `lvt-add-migration` - Create database migration
- `lvt-gen-schema` - Generate database schema
- `lvt-gen-auth` - Generate authentication system
- `lvt-resource-inspect` - List and inspect resources
- `lvt-manage-kits` - Manage CSS framework kits
- `lvt-manage-env` - Manage environment variables
- `lvt-validate-templates` - Validate template syntax
- `lvt-run-and-test` - Start server and run tests
- `lvt-customize` - Customize generated code
- `lvt-seed-data` - Generate test data
- `lvt-deploy` - Deploy to production

**Workflow Skills (4):**
- `lvt-quickstart` - Fast-track common patterns
- `lvt-production-ready` - Production deployment checklist
- `lvt-add-related-resources` - Add multiple related resources
- `lvt-brainstorm` - Interactive planning through progressive questions

**Maintenance Skills (3):**
- `lvt-analyze` - Analyze project structure
- `lvt-suggest` - Suggest improvements
- `lvt-troubleshoot` - Debug issues

**Meta Skill (1):**
- `lvt-add-skill` - Create custom skills

**How to Use:**

**Option 1: Natural Language (Skills activate automatically)**
```
You: Create a blog app with posts and comments
Claude: [Automatically uses lvt-quickstart and lvt-add-resource skills]

You: Add user authentication
Claude: [Automatically uses lvt-gen-auth skill]

You: Generate 50 test posts
Claude: [Automatically uses lvt-seed-data skill]
```

**Option 2: Use the lvt-assistant Agent**

The `lvt-assistant` agent is available in your `.claude/agents/` directory. Like skills, Claude will use it automatically when appropriate for complex LiveTemplate workflows. You don't need to explicitly invoke it - it's available as a resource for handling LiveTemplate-specific tasks.

**Important:** Skills are NOT slash commands. Don't try to type `/lvt-new-app`. Just describe what you want to do in natural language, and Claude will automatically use the appropriate skills.

**Understanding Skill Activation:**

Skills use **keyword-gating** to prevent false positives. A skill activates when **LiveTemplate context is established** through:

1. **Project context** - `.lvtrc` file exists in your project (most common)
2. **Agent context** - You're working with the `lvt-assistant` agent
3. **Keyword context** - You mention "lvt", "livetemplate", or "lt" in your request

**Examples:**

✅ **With Context** (generic prompts work once context is established):
```
You: Create a blog with posts and comments
Claude: [Uses lvt-quickstart skill automatically]

You: Add authentication
Claude: [Uses lvt-gen-auth skill automatically]
```

✅ **Without Context** (keywords required):
```
You: Create a LiveTemplate blog app with posts
Claude: [Uses lvt-new-app skill]

You: Help me build an lvt task manager
Claude: [Uses lvt-quickstart skill]
```

❌ **Without Context or Keywords** (skills won't activate):
```
You: Create a blog app
Claude: [Won't use LiveTemplate skills - could be any framework]
```

**Why Keyword-Gating?**

This prevents skills from activating on generic requests when you might mean a different framework (Next.js, Rails, Django, etc.). Once you mention LiveTemplate or start working in a LiveTemplate project, skills activate seamlessly.

**Using the Brainstorming Skill:**

The `lvt-brainstorm` skill helps you plan LiveTemplate apps through progressive questions:

```
You: I want to build a blog with lvt
Claude: [Uses lvt-brainstorm skill]
        I'll help you plan your LiveTemplate blog. Let me ask a few questions...

        1. What's your primary use case?
        2. Do you need user authentication?
        3. What CSS framework do you prefer?
        ...
```

The brainstorming skill:
- Always requires keywords ("lvt", "livetemplate", or "lt")
- Uses progressive disclosure (starts with 3-5 questions)
- Offers to ask more detailed questions if needed
- Helps you make informed decisions about your app structure

**Verifying Installation:**

After installation and restarting Claude Code:

**Step 1: Verify skills are loaded**
```
You: What skills are available?
Claude: [Lists all 22 lvt skills plus any others]
```

**Step 2: Test skill activation (manual verification)**

Try making a natural language request that should trigger a skill. If Claude uses the Skill tool and references a specific lvt skill, it's working:

```
You: Create a blog app with posts and comments
Expected: Claude should use the Skill tool and invoke lvt-quickstart or lvt-new-app

You: Add a users table with email and password
Expected: Claude should use the Skill tool and invoke lvt-add-resource or lvt-gen-auth
```

**Note:** There's currently no automated way to test skill invocation since it depends on Claude's autonomous decision-making. Skill usage is validated through:
1. Checking that skill files exist and are properly formatted (automated tests)
2. Manual verification by making requests and observing which skills Claude uses

---

### GitHub Copilot (VS Code)

**Best For:** Code-focused development with inline suggestions

**Setup:**
```bash
# Install Copilot agent in your project
lvt install-agent --llm copilot
```

**What You Get:**
- `.github/copilot-instructions.md` - Complete tool reference
- Automatic tool discovery in VS Code
- 16 MCP tools available via Copilot Chat

**Usage:**
```
You: @workspace How do I add a posts resource?
Copilot: Use lvt gen resource posts title:string content:text

You: Add authentication
Copilot: [Suggests using lvt gen auth with options]
```

**MCP Integration:**
Copilot can use MCP tools directly if your environment supports it. See MCP Server section below.

---

### Cursor AI

**Best For:** Full-stack Go development with Composer mode

**Setup:**
```bash
# Install Cursor agent in your project
lvt install-agent --llm cursor
```

**What You Get:**
- `.cursor/rules/lvt.md` - Cursor-specific rules
- File-specific guidance (applies to `**/*.go`)
- Automatic pattern recognition

**Usage:**

**Composer Mode (Recommended):**
```
You: Add a blog with authentication
Cursor: [Chains multiple tools]
  1. lvt new blog
  2. lvt gen resource posts title content
  3. lvt gen auth
  4. lvt migration up
```

**Agent Mode:**
```
You: Create a task management system
Cursor: [Autonomous workflow execution]
```

**Features:**
- Multi-step workflows in Composer
- Context-aware suggestions
- Production best practices enforcement

---

### Aider CLI

**Best For:** Terminal-based pair programming

**Setup:**
```bash
# Install Aider agent in your project
lvt install-agent --llm aider
```

**What You Get:**
- `.aider/.aider.conf.yml` - Aider configuration
- `.aider/lvt-instructions.md` - Command reference
- Auto-commit enabled with conventional commits

**Usage:**
```bash
$ aider

You: Add a posts resource with title and content
Aider: I'll generate the posts resource.
       Running: lvt gen resource posts title:string content:text
       [Shows output]
       Running: lvt migration up
       [Commits changes]
```

**MCP Support (Optional):**
If your Aider version supports MCP, uncomment the server config in `.aider/.aider.conf.yml`:
```yaml
mcp_servers:
  - name: lvt
    command: lvt
    args: [mcp-server]
```

---

### Generic / Other LLMs

**Best For:** Custom LLMs, ChatGPT, Claude API, local models

**Setup:**
```bash
# Install generic documentation
lvt install-agent --llm generic
```

**What You Get:**
- `lvt-agent/README.md` - Complete integration guide (520+ lines)
- `lvt-agent/QUICK_REFERENCE.md` - Quick command reference
- Tool format conversion examples

**Integration Approaches:**

**1. MCP Server (If Supported):**
```bash
lvt mcp-server
```
Configure in your LLM client's MCP settings.

**2. CLI Execution:**
Your LLM executes `lvt` commands via shell.

**3. Tool Calling:**
Convert MCP tool schemas to your LLM's format.

**Example: OpenAI Function Calling**
```json
{
  "name": "lvt_gen_resource",
  "description": "Generate a CRUD resource with database integration",
  "parameters": {
    "type": "object",
    "properties": {
      "name": {"type": "string", "description": "Resource name (singular)"},
      "fields": {
        "type": "object",
        "description": "Field definitions",
        "additionalProperties": {"type": "string"}
      }
    },
    "required": ["name"]
  }
}
```

See `lvt-agent/README.md` for complete conversion examples.

---

## CLI Command Reference

### Generation Commands

```bash
lvt new <name> [--kit multi|single|simple] [--module <path>]
lvt gen resource <name> <field:type>...
lvt gen view <name>
lvt gen auth [StructName] [table_name]
lvt gen schema <table> <field:type>...
```

### Database Commands

```bash
lvt migration up
lvt migration down
lvt migration status
lvt migration create <name>
```

### Resource & Data Commands

```bash
lvt resource list
lvt resource describe <name>
lvt seed <resource> --count <N> [--cleanup]
lvt parse <template-file>
lvt env generate
lvt kits list
lvt kits info <name>
```

### Field Types Reference

```
string     → Go: string,     SQL: TEXT
int        → Go: int64,      SQL: INTEGER
bool       → Go: bool,       SQL: BOOLEAN
float      → Go: float64,    SQL: REAL
time       → Go: time.Time,  SQL: DATETIME
text       → Go: string,     SQL: TEXT (for longer content)
references → Go: int64,      SQL: INTEGER (foreign key)
```

**Foreign Key Example:** `author_id:references:users` creates a foreign key to the `users` table.

---

## Common Workflows

### 1. Quick Start: Blog in 5 Minutes

```bash
lvt new blog --kit multi
cd blog
lvt gen resource posts title:string content:text
lvt migration up
lvt seed posts --count 10
go run cmd/blog/main.go
```

**Result:** Working blog with 10 sample posts at http://localhost:8080

### 2. Full Stack: Task Manager with Auth

```bash
# 1. Create app
lvt new tasks --kit multi

# 2. Add authentication
lvt gen auth

# 3. Add resources
lvt gen resource projects name:string description:text user_id:references:users
lvt gen resource tasks title:string done:bool project_id:references:projects

# 4. Apply migrations
lvt migration up

# 5. Seed data
lvt seed users --count 5
lvt seed projects --count 10
lvt seed tasks --count 50

# 6. Run
go run cmd/tasks/main.go
```

**Result:** Multi-user task manager with projects and tasks

### 3. Incremental Development

**Step 1:** Create basic app
```
You: Create a simple blog
AI: [Creates blog app with lvt new]
```

**Step 2:** Add core resource
```
You: Add posts with title and content
AI: [Generates posts resource]
```

**Step 3:** Add relationships
```
You: Add categories and link posts to categories
AI: [Creates categories resource with foreign key]
```

**Step 4:** Add authentication
```
You: Add user accounts so posts have authors
AI: [Generates auth system, links posts to users]
```

### 4. Production Deployment

```bash
# 1. Generate deployment files
lvt gen stack docker --db postgres

# 2. Build production
go build -o myapp cmd/myapp/main.go

# 3. Deploy
docker-compose up -d
```

---

## Best Practices

### When to Use Which LLM

**Claude Code:**
- Complex projects requiring orchestration
- Need workflow guidance and best practices
- Want natural language interaction

**GitHub Copilot:**
- Working primarily in VS Code
- Want inline code suggestions
- Code-focused development

**Cursor:**
- Full-stack Go development
- Need Composer mode for multi-step workflows
- Want autonomous agent capabilities

**Aider:**
- Terminal-based workflow
- Pair programming style
- Want auto-commits with proper messages

**Generic:**
- Using ChatGPT, Claude API, or local models
- Need custom integration
- Want maximum flexibility

### Common Patterns

**1. CRUD Resource:**
```
Add products with name, price, and stock quantity
```

**2. Authentication:**
```
I need user authentication with email and password
```

**3. Relationships:**
```
Add orders that belong to users and have many products
```

**4. Custom Views:**
```
Add a dashboard view showing statistics
```

**5. Real-time:**
```
Add live updates when new messages arrive
```

---

## Troubleshooting

### Agent Not Working

```bash
# Reinstall agent
lvt install-agent --llm <type> --force

# Check installation
ls -la .claude/  # or .cursor/, .aider/, etc.
```

### Skills Not Showing (Claude)

Ensure you're in a project with `.claude/` directory:
```bash
ls .claude/skills/lvt/
```

If missing, reinstall:
```bash
lvt install-agent --llm claude --force
```

### Commands Failing

Check you're in the correct directory:
```bash
# Should be in project root with go.mod
ls go.mod internal/ cmd/
```

---

## Upgrading Agents

```bash
# Upgrade specific agent type
lvt install-agent --llm <type> --upgrade

# Upgrade Claude agent (default)
lvt install-agent --upgrade
```

**What Gets Preserved:**
- Custom settings (`settings.local.json`)
- Local configuration files
- User-created content

**What Gets Updated:**
- Agent files and skills
- Documentation
- Workflow definitions

---

## Reference Documentation

For more detailed information, see:

- **Setup Guide:** `docs/AGENT_SETUP.md` - Complete setup instructions for all LLMs
- **Workflows:** `docs/WORKFLOWS.md` - Advanced development patterns
- **User Guide:** `docs/guides/user-guide.md` - General usage and concepts

---

## Getting Help

### From Your AI Assistant

```
You: How do I add a new feature?
You: What's the best way to structure my app?
You: I'm getting an error, can you help?
```

### From Documentation

```bash
# General help
lvt --help

# Command-specific help
lvt gen --help
lvt migration --help

# List available agents
lvt install-agent --list
```

### From Community

- **Issues:** https://github.com/livetemplate/lvt/issues
- **Discussions:** https://github.com/livetemplate/lvt/discussions
- **Documentation:** https://livetemplate.io/docs

---

## What's Next?

After setting up your AI assistant:

1. **Try the Quick Start** - Create a blog in 5 minutes
2. **Explore Workflows** - See `docs/WORKFLOWS.md` for patterns
3. **Build Something** - Use your AI assistant to create a real project

The AI assistant will guide you through the entire development process, from project creation to production deployment.
