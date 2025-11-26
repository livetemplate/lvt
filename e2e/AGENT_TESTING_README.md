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
       env.AssertFileExists("internal/app/items/items.go")
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

Tests the MCP server lifecycle and tool registration:

1. **Server Initialization** - Verifies MCP server can be created
2. **Tool Registration** - Tests all 11 registration functions work without panicking
3. **All Tools Registration** - Verifies all 16 tools can be registered together
4. **Context Cancellation** - Tests server context handling
5. **Input/Output Structures** - Validates all tool input and output structs are defined

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
    env.AssertFileExists("internal/app/items/items.go")
}
```

## Future Enhancements

Potential additions:
- Server startup validation
- Browser-based verification
- WebSocket connection testing
- Performance benchmarks
- Screenshot comparisons
- MCP protocol compliance tests
- Concurrent tool invocation tests

Currently kept simple to maximize maintainability and speed.
