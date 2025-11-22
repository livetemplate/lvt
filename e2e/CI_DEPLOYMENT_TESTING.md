# CI/CD Deployment Testing

This document describes the automated testing infrastructure for the lvt project, including continuous integration and deployment testing workflows.

## Overview

The project has two GitHub Actions workflows for testing:

1. **`test.yml`** - Main CI workflow (runs on every PR)
2. **`deployment-tests.yml`** - Deployment testing workflow (manual/scheduled)

## Main CI Workflow (`test.yml`)

**Triggers**: Pull requests (opened, synchronize, reopened)

**Tests Run**:
- ✅ Code formatting check
- ✅ Unit tests (internal packages)
- ✅ Commands package tests
- ✅ E2E tests (short mode - skips deployment tests)
- ✅ Mock deployment infrastructure tests

**Duration**: ~3-5 minutes

**Requirements**:
- Docker (pre-installed on GitHub runners)
- Go 1.25+

**What gets skipped**:
- Real Fly.io deployments (requires `FLY_API_TOKEN`)
- Real Docker deployments (requires `RUN_DOCKER_DEPLOYMENT_TESTS=true`)
- Long-running integration tests

## Deployment Testing Workflow (`deployment-tests.yml`)

**Triggers**:
- Manual dispatch (Actions tab in GitHub)
- Weekly schedule (Sundays at 2 AM UTC)

**Test Modes**:

### 1. Mock Deployment Tests
- **Always runs**: Yes
- **Duration**: ~2 minutes
- **Credentials**: None required
- **What it tests**: Deployment infrastructure with mocked Fly.io API

### 2. Docker Deployment Tests
- **Enabled by**: `run_docker_tests` input (default: true) or scheduled run
- **Duration**: ~10 minutes
- **Credentials**: None required
- **What it tests**:
  - Docker image building
  - Container lifecycle management
  - Health checks and readiness
  - Smoke tests on deployed containers

### 3. Fly.io Deployment Tests
- **Enabled by**: `run_fly_tests` input AND `FLY_API_TOKEN` secret
- **Duration**: ~15 minutes
- **Credentials**: Required (`FLY_API_TOKEN`)
- **What it tests**:
  - Real deployment to Fly.io
  - App creation and volume management
  - Health checks and smoke tests
  - Automatic cleanup

## Running Deployment Tests Manually

### Via GitHub Actions UI

1. Go to **Actions** tab in GitHub
2. Select **Deployment Tests** workflow
3. Click **Run workflow**
4. Configure options:
   - ☑️ Run Fly.io deployment tests (requires credentials)
   - ☑️ Run Docker deployment tests (default: enabled)
5. Click **Run workflow**

### Locally

```bash
# Run all tests (unit + e2e in short mode)
go test ./... -short

# Run only e2e tests (short mode)
go test ./e2e -short -v

# Run Docker deployment tests
RUN_DOCKER_DEPLOYMENT_TESTS=true go test ./e2e -run TestDockerDeployment -v

# Run Fly.io deployment tests (requires credentials)
export FLY_API_TOKEN="your_token_here"
RUN_FLY_DEPLOYMENT_TESTS=true go test ./e2e -run TestRealFlyDeployment -v

# Run mock deployment tests (fast, no credentials)
go test ./e2e -run TestDeploymentInfrastructure_Mock -v
```

## Setting Up Credentials

### For CI/CD

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Add repository secret:
   - **Name**: `FLY_API_TOKEN`
   - **Value**: Your Fly.io API token

### For Local Development

```bash
# Fly.io
export FLY_API_TOKEN="your_fly_token"

# Or use .env file (not committed)
echo "FLY_API_TOKEN=your_fly_token" >> .env
source .env
```

## Test Coverage

### Unit Tests Coverage
- ✅ Internal packages (config, generator, parser, serve, stack, ui, validator)
- ✅ Commands package
- ✅ Deployment infrastructure (credentials, naming, smoke tests)

### E2E Tests Coverage
- ✅ App creation and generation
- ✅ Resource generation
- ✅ CSS framework integration
- ✅ Edit mode functionality
- ✅ Kit management
- ✅ Migration workflows
- ✅ Modal interactions
- ✅ Page mode rendering
- ✅ URL routing
- ✅ View generation

### Deployment Tests Coverage
- ✅ Mock Fly.io deployments (fast, no credentials)
- ✅ Real Fly.io deployments (with credentials)
- ✅ Docker deployments (local and CI)
- ✅ Smoke tests (HTTP, health, WebSocket, templates)
- ✅ Resource-based deployments
- ✅ Cleanup and lifecycle management

## Troubleshooting

### E2E Tests Failing in CI

**Issue**: E2E tests pass locally but fail in CI

**Common causes**:
1. Docker not available
2. Port conflicts
3. Timing issues (increase timeouts)
4. Missing dependencies

**Solution**:
```bash
# Check CI logs for specific errors
# Ensure Docker is running in CI (it should be pre-installed)
# Increase test timeouts if needed
```

### Deployment Tests Timing Out

**Issue**: Deployment tests exceed timeout

**Solutions**:
- Increase timeout in workflow: `-timeout=20m`
- Check if deployment region is responsive
- Verify credentials are correct
- Check Fly.io status page for outages

### Docker Tests Failing

**Issue**: Docker build or run failures

**Common causes**:
1. Missing Dockerfile (auto-generated in tests)
2. go.sum missing (runs `go mod tidy` automatically)
3. Port conflicts

**Solution**:
```bash
# Clean up Docker resources
docker system prune -a

# Verify Docker is running
docker ps

# Run with verbose output
go test ./e2e -run TestDockerDeployment -v
```

## Performance Optimization

### Fast CI Runs
- Use `-short` flag to skip slow tests
- Run only changed packages
- Cache Go modules and build artifacts
- Use matrix builds for parallel execution

### Efficient Deployment Testing
- Use mock tests for quick feedback
- Schedule real deployments during off-hours
- Reuse test apps when possible
- Implement proper cleanup to avoid costs

## Best Practices

1. **Always run tests locally before pushing**
   ```bash
   go test ./... -short
   ```

2. **Use descriptive test names**
   - Good: `TestDockerDeploymentWithResources`
   - Bad: `TestDeployment`

3. **Clean up resources**
   - Use `t.Cleanup()` for automatic cleanup
   - Implement defer-based cleanup
   - Handle errors in cleanup functions

4. **Make tests idempotent**
   - Tests should pass regardless of order
   - Don't rely on external state
   - Generate unique names for test resources

5. **Document test requirements**
   - List required credentials
   - Note expected duration
   - Explain skip conditions

## Future Enhancements

- [ ] Matrix testing across Go versions
- [ ] Parallel test execution
- [ ] Test result caching
- [ ] Performance benchmarking
- [ ] Multi-region deployment testing
- [ ] Load testing for deployed apps
- [ ] Security scanning in CI

## Related Documentation

- [Deployment Testing Plan](./DEPLOYMENT_TESTING_PLAN.md)
- [Deployment Testing Progress](./DEPLOYMENT_TESTING_PLAN_UPDATE.md)
- [Main Documentation](./DEPLOYMENT_TESTING.md)
