# LiveTemplate AI Agent Setup Guide

This guide shows how to set up LiveTemplate with different AI assistants using either the MCP server (recommended) or agent installation.

## Quick Start

```bash
# Install LiveTemplate
go install github.com/livetemplate/lvt@latest

# Option 1: Install agent for your LLM (project-specific)
lvt install-agent                    # Interactive menu (recommended)
lvt install-agent --llm <type>       # Direct installation

# Option 2: Start MCP server (global, works with all LLMs)
lvt mcp-server                       # Interactive setup (recommended)
lvt mcp-server --help                # Show full documentation
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

### Method 1: Agent Installation (Recommended)

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

### Method 2: MCP Server Only

If you prefer using only the MCP server:

1. **Run the interactive setup:**
   ```bash
   lvt mcp-server           # Shows setup wizard
   ```

   Or see help for manual setup:
   ```bash
   lvt mcp-server --help    # Full documentation
   lvt mcp-server --setup   # Interactive wizard
   ```

2. **The wizard will guide you through:**
   - Selecting your AI client (Claude Desktop, VS Code, etc.)
   - Finding your config file location
   - Adding the MCP server configuration
   - Restarting your client

3. **Manual configuration (if preferred):**

   Add to your Claude Code MCP configuration:
   ```json
   {
     "mcpServers": {
       "lvt": {
         "command": "lvt",
         "args": ["mcp-server"]
       }
     }
   }
   ```

4. **Restart your AI client**

### Method 3: Both (Best Experience)

Combine agent installation + MCP server for the complete experience:

```bash
# Install agent (interactive menu)
lvt install-agent

# MCP server auto-configured in agent
claude
```

The agent uses MCP tools under the hood while providing workflow guidance.

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

### MCP Server (Optional)

If your Copilot environment supports MCP:

```bash
# Start MCP server
lvt mcp-server

# Configure in your IDE's MCP settings
```

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

### MCP Server (Optional)

For MCP support:

```yaml
# .aider/.aider.conf.yml
mcp_servers:
  - name: lvt
    command: lvt
    args: [mcp-server]
```

---

## Generic LLM Setup

**Recommended Approach:** MCP Server + Documentation

For LLMs not listed above (ChatGPT, Claude API, local models, etc.):

### Setup

1. **Install generic documentation:**
   ```bash
   lvt install-agent --llm generic
   ```

   This creates `lvt-agent/` with:
   - `README.md` - Complete integration guide
   - `QUICK_REFERENCE.md` - Quick command reference

2. **Start MCP server:**
   ```bash
   lvt mcp-server
   ```

3. **Configure your LLM:**

   **For MCP-enabled LLMs:**
   ```json
   {
     "mcpServers": {
       "lvt": {
         "command": "lvt",
         "args": ["mcp-server"]
       }
     }
   }
   ```

   **For CLI-capable LLMs:**

   Point your LLM to execute `lvt` commands via shell.

   **For Tool-calling LLMs:**

   Convert MCP tool schemas to your LLM's format. See `lvt-agent/README.md` for examples.

### Documentation

The generic agent includes:
- Complete MCP tool reference
- CLI command reference
- Field type documentation
- Common workflow patterns
- Example sessions

---

## MCP Server Reference

The LiveTemplate MCP server provides 16 tools globally:

### Starting the Server

```bash
# Start the server (for AI client use)
lvt mcp-server

# Show setup instructions
lvt mcp-server --help

# Interactive setup wizard
lvt mcp-server --setup

# List all available tools
lvt mcp-server --list-tools

# Show version information
lvt mcp-server --version
```

**Note:** The MCP server runs as a JSON-RPC service over stdio. Don't run it directly in your terminal - configure it in your AI client instead.

### Getting Help

If you're setting up the MCP server for the first time:

```bash
# Show detailed setup instructions for your platform
lvt mcp-server --setup
```

This will guide you through:
1. Locating your config file (platform-specific)
2. Adding the MCP server configuration
3. Restarting your AI client
4. Verifying the installation

### Available Tools

**Generation (5 tools):**
- `lvt_new` - Create new app
- `lvt_gen_resource` - Add CRUD resource
- `lvt_gen_view` - Add view-only page
- `lvt_gen_auth` - Add authentication
- `lvt_gen_schema` - Add database schema

**Database (4 tools):**
- `lvt_migration_up` - Apply migrations
- `lvt_migration_down` - Rollback migration
- `lvt_migration_status` - Check status
- `lvt_migration_create` - Create migration

**Development (7 tools):**
- `lvt_seed` - Generate test data
- `lvt_resource_list` - List resources
- `lvt_resource_describe` - Describe schema
- `lvt_validate_template` - Validate templates
- `lvt_env_generate` - Generate .env
- `lvt_kits_list` - List kits
- `lvt_kits_info` - Get kit info

### Testing the Server

```bash
# Start server
lvt mcp-server

# In another terminal, test with MCP inspector
npx @modelcontextprotocol/inspector lvt mcp-server
```

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

### Use MCP Server If:
- You want global access across all projects
- Your LLM supports MCP protocol
- You prefer minimal project setup
- You're using a custom/different LLM

### Use Both If:
- You want the best of both worlds
- You're using Claude Code (agent + MCP)
- You want workflows + global tools

---

## Troubleshooting

### Agent Not Working

```bash
# Reinstall agent
lvt install-agent --llm <type> --force

# Check installation
ls -la .claude/  # or .cursor/, .aider/, etc.
```

### MCP Server Not Connecting

```bash
# Get setup instructions
lvt mcp-server --setup

# Verify tools are available
lvt mcp-server --list-tools

# Check version
lvt mcp-server --version

# Test with inspector
npx @modelcontextprotocol/inspector lvt mcp-server
```

**Common Issues:**

1. **Running in terminal directly** - The server will show a warning. Configure it in your AI client instead.
2. **Wrong config path** - Use `lvt mcp-server --setup` to see the correct path for your platform.
3. **Invalid JSON** - Check your config file syntax with a JSON validator.
4. **lvt not in PATH** - Ensure `lvt` is installed and accessible from your terminal.

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
   - `docs/MCP_TOOLS.md` - Complete tool reference
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
- **MCP Protocol:** https://modelcontextprotocol.io
