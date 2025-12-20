# Agent Testing Framework

This directory contains automated tests that validate the LiveTemplate agent usage documentation.

## Purpose

The agent test framework ensures that all examples and workflows documented in the agent usage guide actually work as described. This prevents documentation drift and gives users confidence that the guide is accurate.

## Test Structure

### Test Harness (`internal/agenttest/harness.go`)

The test harness provides:
- Isolated test environments using temp directories
- Command execution tracking
- Skill usage tracking
- File assertion helpers
- Database table verification
- Simulated agent workflows

### Documentation Validation Tests (`e2e/agent_doc_validation_test.go`)

These tests validate every major example from the usage guide:

1. **Quick Start** - Basic blog with posts
2. **Full Stack** - Task manager with auth, projects, tasks
3. **Step-by-Step** - Incremental blog building
4. **Common Patterns** - Auth, CRUD resources, relationships
5. **All Kits** - multi, single, simple kit creation
6. **Incremental Features** - Building blog piece by piece
7. **Recipe Example** - Complete example session

## Running the Tests

```bash
# Run all agent documentation validation tests
go test -v -run TestAgentDocValidation ./e2e

# Run a specific test
go test -v -run TestAgentDocValidation_QuickStart ./e2e

# Run with race detection
go test -race -v -run TestAgentDocValidation ./e2e
```

##  What Gets Tested

Each test verifies:

- ✅ Correct lvt commands execute successfully
- ✅ Expected files are created in correct locations
- ✅ Database tables are created properly
- ✅ Generated code contains expected content
- ✅ Skills are tracked correctly
- ✅ Command history is accurate

## Test Environment

Tests run in complete isolation:
- Each test gets a fresh temp directory
- Apps are created with `SKIP_GO_MOD_TIDY=1` for speed
- No shared state between tests
- Automatic cleanup

## Adding New Tests

To add a new test:

1. Add example to usage documentation
2. Create test in `agent_doc_validation_test.go`:
   ```go
   func TestAgentDocValidation_NewFeature(t *testing.T) {
       env := agenttest.Setup(t, &agenttest.SetupOptions{
           AppName: "testapp",
           Kit:     "multi",
       })

       // Simulate the workflow
       err := env.RunLvtCommand("gen", "resource", "items", "name:string")
       require.NoError(t, err)

       // Verify outcomes
       env.AssertFileExists("app/items/items.go")
   }
   ```
3. Run the test to verify it works
4. Commit both documentation and test

## Test Philosophy

**Tests follow documentation, not the other way around.**

If a test fails, it means either:
1. The documentation is inaccurate (update docs)
2. The implementation changed (update tests + docs)
3. There's a bug (fix it)

Never update tests to pass without also updating documentation.

## Current Limitations

- Tests validate file creation and basic structure
- Do not start browsers or servers (kept simple intentionally)
- Do not test actual user interactions
- Focus on command-line interface and file generation

For full E2E browser testing, see the regular e2e test suite.

## MCP Server Testing

### Overview

The MCP (Model Context Protocol) server provides tool-based access to lvt functionality for AI agents like Claude Desktop. These tests ensure all MCP tools work correctly.

### MCP Server Tests (`commands/mcp_server_test.go`)

Tests the MCP server lifecycle, tool registration, and helper flags:

#### Core Server Tests (5 tests)
1. **Server Initialization** - Verifies MCP server can be created
2. **Tool Registration** - Tests all 11 registration functions work without panicking
3. **All Tools Registration** - Verifies all 16 tools can be registered together
4. **Context Cancellation** - Tests server context handling
5. **Input/Output Structures** - Validates all tool input and output structs are defined

#### Flag Tests (6 tests)
6. **Help Flag** - Tests `--help` and `-h` flags display setup instructions
7. **Setup Flag** - Tests `--setup` interactive wizard (skipped in automated tests)
8. **List Tools Flag** - Tests `--list-tools` shows all 16 tools
9. **Version Flag** - Tests `--version` shows MCP protocol version
10. **Invalid Flag** - Tests handling of invalid flags
11. **Multiple Flags** - Tests that first flag takes precedence

### MCP Tool Tests (`commands/mcp_tools_test.go`)

Comprehensive tests for all 16 MCP tools:

#### Core Generation Tools (5 tools)
- `lvt_new` - Create new apps with different configurations
- `lvt_gen_resource` - Generate CRUD resources with fields
- `lvt_gen_view` - Generate view-only handlers
- `lvt_gen_auth` - Generate authentication systems
- `lvt_gen_schema` - Generate database schema only

#### Migration Tools (4 tools)
- `lvt_migration_up` - Run pending migrations
- `lvt_migration_down` - Rollback last migration
- `lvt_migration_status` - Show migration status
- `lvt_migration_create` - Create new migration files

#### Resource Inspection Tools (2 tools)
- `lvt_resource_list` - List all available resources
- `lvt_resource_describe` - Show detailed schema for a resource

#### Data Management Tools (1 tool)
- `lvt_seed` - Generate test data for resources

#### Template Tools (1 tool)
- `lvt_validate_template` - Validate and analyze template files

#### Environment Tools (1 tool)
- `lvt_env_generate` - Generate .env.example with detected config

#### Kits Tools (2 tools)
- `lvt_kits_list` - List available CSS framework kits
- `lvt_kits_info` - Show detailed kit information

### Running MCP Tests

```bash
# Run all MCP server tests
go test -v -run TestMCPServer ./commands

# Run all MCP tool tests
go test -v -run TestMCPTool ./commands

# Run specific tool test
go test -v -run TestMCPTool_LvtNew ./commands

# Run with short mode (skips long-running tests)
SKIP_GO_MOD_TIDY=1 go test -short -run TestMCPTool ./commands
```

### What MCP Tests Validate

Each tool test verifies:
- ✅ Valid inputs produce successful results
- ✅ Files and directories are created correctly
- ✅ Commands execute without errors
- ✅ Output structures are properly defined
- ✅ Integration with underlying commands works

### MCP Test Environment

MCP tests use isolated environments:
- Each test gets a fresh temp directory
- Apps are created with `SKIP_GO_MOD_TIDY=1` for speed
- Tests clean up automatically
- No shared state between tests

## Agent Test Harness Extensions

### Supported Commands

The agent test harness (`internal/agenttest/harness.go`) now supports all major lvt commands:

- `gen` - Generate resources, views, auth, schemas
- `migration` - Database migration operations
- `seed` - Generate test data
- `resource` / `res` - Inspect resources and schemas
- `parse` - Validate template files
- `env` - Environment variable management
- `kits` / `kit` - CSS framework kit management

### Example Usage

```go
func TestMyWorkflow(t *testing.T) {
    env := agenttest.Setup(t, &agenttest.SetupOptions{
        AppName: "testapp",
        Kit:     "multi",
    })

    // Generate a resource
    err := env.RunLvtCommand("gen", "resource", "items", "name:string")
    require.NoError(t, err)

    // Seed data
    err = env.RunLvtCommand("seed", "items", "--count", "10")
    require.NoError(t, err)

    // List resources
    err = env.RunLvtCommand("resource", "list")
    require.NoError(t, err)

    // Verify files exist
    env.AssertFileExists("app/items/items.go")
}
```

## Multi-LLM Agent Installation Testing

### Overview

The multi-LLM agent installation tests ensure that LiveTemplate can be integrated with different AI assistants through various installation methods.

### Agent Installation Tests (`commands/install_agent_multi_llm_test.go`)

Comprehensive tests for multi-LLM agent support:

#### Agent Type Tests (5 agents)
- `TestInstallAgent_CopilotAgent` - GitHub Copilot instructions installation
- `TestInstallAgent_CursorAgent` - Cursor AI rules installation
- `TestInstallAgent_AiderAgent` - Aider CLI configuration installation
- `TestInstallAgent_GenericAgent` - Generic LLM documentation installation
- `TestInstallAgent_DefaultIsClaudeCode` - Verify Claude is the default

#### List and Validation Tests
- `TestInstallAgent_ListAgents` - Test --list flag functionality
- `TestInstallAgent_InvalidLLMType` - Error handling for invalid LLM types

#### Upgrade Tests
- `TestInstallAgent_UpgradeCursorAgent` - Upgrade preserves custom files
- `TestInstallAgent_UpgradeAiderAgent` - Upgrade preserves .aider.local.yml

#### Force Installation Tests
- `TestInstallAgent_ExistingCopilotAgent` - --force flag overwrites existing installations

### What Multi-LLM Tests Validate

Each agent installation test verifies:
- ✅ Correct directory structure is created
- ✅ Required files are installed
- ✅ File content is not empty
- ✅ Files contain expected keywords/markers
- ✅ Special file handling (e.g., .aider.conf.yml renaming)
- ✅ Upgrade preserves custom user files
- ✅ Error handling for edge cases

### Running Multi-LLM Tests

```bash
# Run all multi-LLM agent tests
go test -v ./commands -run TestInstallAgent_

# Run specific agent test
go test -v ./commands -run TestInstallAgent_CopilotAgent

# Run only upgrade tests
go test -v ./commands -run TestInstallAgent_Upgrade

# Run with short mode
go test -short ./commands -run TestInstallAgent_
```

### Agent Types and Directories

| Agent | Directory | Key Files |
|-------|-----------|-----------|
| Claude Code | `.claude/` | settings.json, skills/, agents/ |
| GitHub Copilot | `.github/` | copilot-instructions.md |
| Cursor | `.cursor/rules/` | lvt.md |
| Aider | `.aider/` | .aider.conf.yml, lvt-instructions.md |
| Generic | `lvt-agent/` | README.md, QUICK_REFERENCE.md |

### Custom File Preservation

The upgrade mechanism preserves these custom files:
- `settings.local.json` (Claude Code)
- `.aider.local.yml` (Aider)
- `local.md` (Cursor)

### Special Handling

**Aider Configuration:**
The Aider agent has special file renaming logic because Go's embed.FS doesn't include files starting with `.`:
- Embedded as: `aider.conf.yml`
- Installed as: `.aider.conf.yml`

This is tested in `TestInstallAgent_AiderAgent` which verifies the rename happens correctly.

## Documentation Testing

### Agent Usage Guide (`docs/AGENT_USAGE_GUIDE.md`)

The comprehensive AI agent usage guide covers all integration approaches:

#### Multi-LLM Support (5 LLM types)
- **Claude Desktop / Claude Code** - 22 skills with workflow orchestration (includes brainstorming)
- **GitHub Copilot** - VS Code integration with MCP tools
- **Cursor AI** - Composer/Agent mode with file-specific rules
- **Aider CLI** - Terminal-based pair programming
- **Generic** - Custom LLMs, ChatGPT, Claude API, local models

#### Integration Approaches (2 methods)
- **MCP Server** - Global tool access (16 tools)
- **Agent Installation** - Project-specific guidance

#### Documentation Coverage
- ✅ Complete MCP tool reference with JSON schemas
- ✅ All 22 Claude skills documented by category (Core: 14, Workflow: 4, Maintenance: 3, Meta: 1)
- ✅ Keyword-gating explanation for skill activation
- ✅ Brainstorming skill for interactive planning
- ✅ Field type reference (string, int, bool, float, time, text, references)
- ✅ Common workflows (Quick Start, Full Stack, Incremental, Production)
- ✅ Best practices for each LLM type
- ✅ Troubleshooting guide
- ✅ Upgrade procedures

#### What Gets Tested
The documentation examples are validated through:
- Multi-LLM agent installation tests (10 tests in `commands/install_agent_multi_llm_test.go`)
- MCP server flag tests (6 tests in `commands/mcp_server_test.go`)
- MCP tool tests (16 tools in `commands/mcp_tools_test.go`)
- Agent workflow tests (planned - see Future Enhancements)

### Interactive MCP Setup

The `lvt mcp-server --setup` command provides an interactive wizard:

#### Features Tested
- ✅ Multi-LLM selection menu (5 options)
- ✅ Platform-specific config paths (macOS, Linux, Windows)
- ✅ LLM-specific setup instructions
- ✅ Agent vs MCP recommendations
- ✅ Graceful fallback for invalid choices

#### Setup Options
1. **Claude Desktop / Claude Code** - MCP server config with platform detection
2. **VS Code with Copilot** - Agent installation recommendation
3. **Cursor AI** - Agent installation with Composer guidance
4. **Aider CLI** - Both agent and optional MCP config
5. **Other MCP clients** - Generic MCP setup instructions

#### Test Coverage
Interactive setup is tested in `TestMCPServer_SetupFlag` (skipped in automated tests due to user input requirement). Manual testing covers:
- All 5 LLM choice paths
- Platform detection accuracy
- Config path correctness
- Instruction completeness

## Current Gaps

### Agent and Skills Validation

**Problem:** The AGENT_USAGE_GUIDE.md documents usage of the lvt-assistant agent and 21 skills, but we don't test:

1. ✗ lvt-assistant agent file exists
2. ✗ All 21 documented skills exist in `.claude/skills/`
3. ✗ Skill invocation syntax (`/lvt-assistant`, `/lvt-quickstart`) is correct
4. ✗ Examples in guide match actual skill capabilities
5. ✗ Skills can be invoked without errors
6. ✗ Agent provides correct guidance

**Current Testing:**
- ✅ We track that skills are "used" in workflows (`env.TrackSkill("lvt-quickstart")`)
- ✅ We verify CLI commands execute successfully
- ✅ **NEW:** We verify skills and agent actually exist (see `agent_skills_validation_test.go`)
- ✅ **NEW:** We test skill/agent invocation syntax

**Implemented Tests (in `e2e/agent_skills_validation_test.go`):**

```go
// Test that Claude agent installation includes all documented components
func TestClaudeAgent_Installation(t *testing.T) {
    // ✅ Verifies .claude/agents/lvt-assistant.md exists
    // ✅ Verifies .claude/skills/lvt/ directory exists
    // ✅ Verifies settings.json exists
}

// Test that all 22 documented skills exist
func TestClaudeAgent_AllSkillsExist(t *testing.T) {
    // ✅ Verifies all 22 documented skills are present
    // ✅ Parses each skill file and verifies frontmatter
    // ✅ Skills organized in 4 categories: core/, workflows/, maintenance/, meta/
}

// Test that agent frontmatter is correct
func TestClaudeAgent_AgentMetadata(t *testing.T) {
    // ✅ Verifies frontmatter has name: "lvt-assistant"
    // ✅ Verifies description exists
    // ✅ Verifies agent mentions LiveTemplate
}

// Test that skill names match directory structure
func TestClaudeAgent_SkillInvocationSyntax(t *testing.T) {
    // ✅ Verifies skill frontmatter names follow lvt-skill-name pattern
    // ✅ Verifies names match directory structure (e.g., new-app/ → lvt-new-app)
}

// Test skill count matches documentation
func TestClaudeAgent_SkillCount(t *testing.T) {
    // ✅ Verifies we have at least 22 skills as documented (Core: 14, Workflow: 4, Maintenance: 3, Meta: 1)
    // ✅ Counts across all 4 category directories
}
```

**Why This Matters:**
- ✅ **SOLVED:** Documentation drift prevented by validation tests
- ✅ **SOLVED:** Breaking changes (renaming skills) caught immediately
- ✅ **SOLVED:** User confusion prevented - examples are tested

## Future Enhancements

Potential additions:
- **Agent/Skills validation tests** (see Current Gaps above)
- Server startup validation
- Browser-based verification
- WebSocket connection testing
- Performance benchmarks
- Screenshot comparisons
- MCP protocol compliance tests
- Concurrent tool invocation tests
- Agent workflow integration tests
- Cross-LLM compatibility tests
- Interactive setup automated testing (with input mocking)
- AGENT_USAGE_GUIDE.md example validation tests
- Skill capability vs documentation matching tests

Currently kept simple to maximize maintainability and speed.
