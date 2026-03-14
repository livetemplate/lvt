# LiveTemplate AI Agent Setup Guide

This guide shows how to set up LiveTemplate with different AI assistants using agent installation.

## Quick Start

```bash
# Install LiveTemplate
go install github.com/livetemplate/lvt@latest

# Install agent for your LLM (project-specific)
lvt install-agent                    # Interactive menu (recommended)
lvt install-agent --llm <type>       # Direct installation
```

## Available Agents

List all available agent types:

```bash
lvt install-agent --list
```

Output:
```
Available AI agents for LiveTemplate:

  claude    - Claude Code agent with skills and project management
  aider     - Aider CLI configuration for LiveTemplate
  copilot   - GitHub Copilot instructions for LiveTemplate development
  cursor    - Cursor AI rules for LiveTemplate development
  generic   - LLM-agnostic documentation and guides
```

---

## Claude Code Setup

**Recommended Approach:** Agent Installation (autonomous skills + specialist agent)

### Setup

Install Claude Code agent and skills in your project:

```bash
# In your LiveTemplate project
lvt install-agent              # Interactive menu - select Claude Code
lvt install-agent --llm claude # Direct installation
```

**What gets installed:**
- `.claude/skills/` - 22 autonomous skills (each in its own directory)
- `.claude/agents/lvt-assistant.md` - LiveTemplate specialist agent
- `.claude/settings.json` - Configuration and permissions

**Understanding what you get:**

**Skills (22 total):** Autonomous capabilities that Claude uses automatically
- NOT slash commands - they activate based on your requests
- Claude reads their descriptions and decides when to use them
- Examples: `lvt-new-app`, `lvt-add-resource`, `lvt-gen-auth`, `lvt-brainstorm`
- Verify they're loaded by asking: "What skills are available?"
- Skills use **keyword-gating** to prevent false positives (require "lvt", "livetemplate", or "lt" unless working in a LiveTemplate project)

**Agent:** Specialist available for complex workflows
- The `lvt-assistant` agent in `.claude/agents/` provides LiveTemplate expertise
- Claude can use it automatically for complex multi-step tasks
- You can view it in the `/agents` menu but don't need to explicitly invoke it
- Can use all 21 skills plus general capabilities

**Start Claude Code:**
```bash
claude
```

**IMPORTANT:** Restart Claude Code completely after installation for skills to load.

---

## GitHub Copilot Setup

**Recommended Approach:** Agent Installation

### Setup

1. **Install the Copilot agent:**
   ```bash
   lvt install-agent --llm copilot
   ```

   This creates `.github/copilot-instructions.md` with LiveTemplate instructions.

2. **Open your project in VS Code/IDE** with Copilot enabled

3. **Copilot automatically reads the instructions** and understands LiveTemplate

### Usage

- Open any file and start coding
- Use `@workspace` to ask questions about LiveTemplate
- Copilot will suggest using `lvt` commands
- In chat, ask: "How do I add a posts resource?"

---

## Cursor Setup

**Recommended Approach:** Agent Installation

### Setup

1. **Install the Cursor agent:**
   ```bash
   lvt install-agent --llm cursor
   ```

   This creates `.cursor/rules/lvt.md` with LiveTemplate rules.

2. **Open project in Cursor**

3. **Rules apply automatically** to `*.go` files

### Usage

- Use **Composer mode** for best results
- Use **Agent mode** for autonomous workflows
- Ask: "Add a blog with authentication"
- Cursor follows LiveTemplate patterns automatically

### Features

- File-specific rules (applies to `**/*.go`)
- Common workflow patterns
- Error handling guidance
- Best practices enforcement

---

## Aider Setup

**Recommended Approach:** Agent Installation

### Setup

1. **Install the Aider agent:**
   ```bash
   lvt install-agent --llm aider
   ```

   This creates:
   - `.aider/.aider.conf.yml` - Configuration
   - `.aider/lvt-instructions.md` - Instructions

2. **Start Aider:**
   ```bash
   aider
   ```

   Configuration loads automatically.

### Usage

```bash
# Aider automatically knows about lvt commands
aider> Add a posts resource with title and content

# Aider will use: lvt gen resource posts title:string content:text
# Then: lvt migration up
```

### Configuration

The `.aider.conf.yml` includes:
- Auto-commits enabled
- Custom commit messages
- LiveTemplate instructions

---

## Generic LLM Setup

For LLMs not listed above (ChatGPT, Claude API, local models, etc.):

### Setup

1. **Install generic documentation:**
   ```bash
   lvt install-agent --llm generic
   ```

   This creates `lvt-agent/` with:
   - `README.md` - Complete integration guide
   - `QUICK_REFERENCE.md` - Quick command reference

2. **Point your LLM to execute `lvt` CLI commands via shell.**

### Documentation

The generic agent includes:
- CLI command reference
- Field type documentation
- Common workflow patterns
- Example sessions

---

## Upgrading Agents

Upgrade an existing agent installation:

```bash
# Upgrade specific agent
lvt install-agent --llm <type> --upgrade

# Upgrade Claude agent
lvt install-agent --upgrade
```

**What gets preserved:**
- Custom settings (`settings.local.json`)
- Local configuration files

**What gets updated:**
- Agent files and skills
- Documentation
- Workflows

---

## Choosing Your Setup

### Use Agent Installation If:
- You want guided workflows and best practices
- You prefer project-specific configuration
- You want your LLM to understand LiveTemplate patterns
- You're using Claude Code, Cursor, Copilot, or Aider

---

## Troubleshooting

### Agent Not Working

```bash
# Reinstall agent
lvt install-agent --llm <type> --force

# Check installation
ls -la .claude/  # or .cursor/, .aider/, etc.
```

### Wrong Agent Type

```bash
# List available types
lvt install-agent --list

# Install correct type
lvt install-agent --llm <correct-type> --force
```

---

## Next Steps

After setup:

1. **Read the documentation:**
   - `docs/WORKFLOWS.md` - Common development workflows
   - `.claude/`, `.cursor/`, etc. - LLM-specific guides

2. **Try a workflow:**
   ```bash
   # Create new app
   lvt new myblog --kit multi

   # Install agent
   lvt install-agent --llm <your-llm>

   # Start your AI assistant
   # Ask it: "Add a posts resource with authentication"
   ```

3. **Explore the tools:**
   ```bash
   lvt --help
   lvt gen --help
   lvt migration --help
   ```

---

## Support

- **Documentation:** `docs/` directory
- **Issues:** https://github.com/livetemplate/lvt/issues
