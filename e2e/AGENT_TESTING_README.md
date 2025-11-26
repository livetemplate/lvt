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

## Future Enhancements

Potential additions:
- Server startup validation
- Browser-based verification
- WebSocket connection testing
- Performance benchmarks
- Screenshot comparisons

Currently kept simple to maximize maintainability and speed.
