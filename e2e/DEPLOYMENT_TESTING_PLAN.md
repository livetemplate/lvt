# Deployment Testing Implementation Plan & Tracker

**Goal**: Build production-ready end-to-end testing from `lvt new` through Fly.io deployment with full verification.

**Status**: 🚧 In Progress
**Started**: 2025-11-05
**Target**: Complete deployment testing infrastructure

---

## Phase 1: Core Infrastructure (Foundation)

### ✅ Task 1.1: Credentials Management
**File**: `testing/credentials.go`
**Status**: ⬜ Not Started
**Time**: 30 min

**Implementation**:
- [ ] Create TestCredentials struct
- [ ] LoadTestCredentials() from environment
- [ ] RequireCredentials(t, provider) helper
- [ ] ValidateCredentials(provider) checker
- [ ] Support FLY_API_TOKEN, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, S3_BUCKET

**Acceptance**:
- Tests skip gracefully when credentials missing
- Validates credentials are present for provider
- Loads from environment variables

---

### ✅ Task 1.2: App Name Generator
**File**: `testing/naming.go`
**Status**: ⬜ Not Started
**Time**: 15 min

**Implementation**:
- [ ] GenerateTestAppName(prefix) function
- [ ] Format: `lvt-test-{random}-{timestamp}`
- [ ] Validate Fly.io naming rules
- [ ] Generate cryptographically random suffix

**Acceptance**:
- Unique names for parallel tests
- Valid for Fly.io (lowercase, alphanumeric, hyphens)
- No collisions in 10,000 iterations

---

### ✅ Task 1.3: Deployment Test Harness
**File**: `testing/deployment.go`
**Status**: ⬜ Not Started
**Time**: 1 hour

**Implementation**:
- [ ] DeploymentTest struct (Provider, AppName, AppDir, AppURL, CleanupFuncs)
- [ ] SetupDeployment(t, opts) - Creates test app with unique name
- [ ] Deploy() method - Executes deployment
- [ ] VerifyHealth() method - Checks app is responding
- [ ] VerifyWebSocket() method - Tests WebSocket connection
- [ ] Cleanup() method - Destroys all resources
- [ ] Defer-based cleanup tracking
- [ ] Error aggregation for partial cleanups

**Acceptance**:
- Can deploy and cleanup successfully
- Cleanup runs even on test failure
- All resources destroyed (no cost leaks)
- Health check validates app is running

---

## Phase 2: Provider Implementations

### ✅ Task 2.1: Mock Fly.io Client
**File**: `testing/providers/mock.go`
**Status**: ⬜ Not Started
**Time**: 45 min

**Implementation**:
- [ ] MockFlyClient struct
- [ ] Launch(appName) - Returns mock success
- [ ] Deploy(appDir) - Returns mock success
- [ ] Status(appName) - Returns mock status
- [ ] CreateVolume() - Returns mock volume ID
- [ ] Destroy(appName) - Returns mock success
- [ ] GetAppURL() - Returns mock URL (http://mock-app.fly.dev)

**Acceptance**:
- All methods return predictable mock data
- No actual API calls made
- Fast execution (<100ms)
- Can run offline

---

### ✅ Task 2.2: Real Fly.io Helpers
**File**: `testing/providers/fly.go`
**Status**: ⬜ Not Started
**Time**: 1.5 hours

**Implementation**:
- [ ] DeployToFly(appName, appDir) - Wrapper around `flyctl deploy`
- [ ] CreateFlyVolume(appName, region, sizeGB) - Creates persistent volume
- [ ] DestroyFlyApp(appName) - Destroys app and volumes
- [ ] GetFlyAppURL(appName) - Retrieves app URL from flyctl
- [ ] GetFlyAppStatus(appName) - Checks app status
- [ ] WaitForFlyAppReady(appName, timeout) - Waits for deployment
- [ ] Retry logic for network failures
- [ ] Proper error handling and messages

**Acceptance**:
- Actually deploys to Fly.io
- Creates and destroys volumes
- Waits for deployment to complete
- Returns correct app URL
- Cleans up resources properly

---

### ✅ Task 2.3: Docker Helpers
**File**: `testing/providers/docker.go`
**Status**: ⬜ Not Started
**Time**: 45 min

**Implementation**:
- [ ] DockerComposeUp(composeFile) - Starts containers
- [ ] DockerComposeDown(composeFile) - Stops and removes
- [ ] DockerHealthCheck(containerName) - Verifies health
- [ ] DockerGetContainerURL(containerName) - Returns container URL
- [ ] WaitForDockerHealthy(containerName, timeout)

**Acceptance**:
- Can start/stop docker-compose apps
- Health checks work correctly
- Cleanup removes all containers
- Works on macOS, Linux, Windows

---

## Phase 3: Verification & Testing

### ✅ Task 3.1: Smoke Test Suite
**File**: `testing/smoke.go`
**Status**: ⬜ Not Started
**Time**: 1 hour

**Implementation**:
- [ ] RunSmokeTests(appURL) function
- [ ] Test 1: HTTP 200 on root path
- [ ] Test 2: /health endpoint responds {"status":"ok"}
- [ ] Test 3: WebSocket connection establishes
- [ ] Test 4: Database CRUD (create test record)
- [ ] Test 5: Static assets load (livetemplate-client.js)
- [ ] Test 6: Templates render without errors
- [ ] Test 7: Server logs show no errors
- [ ] Configurable timeout and retry logic
- [ ] Detailed error messages on failure

**Acceptance**:
- All 7 smoke tests pass against deployed app
- Tests retry transient failures
- Clear error messages when tests fail
- Works with both local and remote deployments

---

### ✅ Task 3.2: E2E Deployment Workflow Test (Mock)
**File**: `e2e/deployment_workflow_mock_test.go`
**Status**: ⬜ Not Started
**Time**: 45 min

**Implementation**:
- [ ] TestCompleteWorkflow_Mock() function
- [ ] Create test app (`lvt new blogapp`)
- [ ] Add resources (`lvt gen resource posts title content`)
- [ ] Add auth (`lvt gen auth`)
- [ ] Generate Fly stack (`lvt gen stack fly`)
- [ ] Deploy using mock client
- [ ] Verify deployment succeeded
- [ ] Run smoke tests against mock URL
- [ ] Cleanup test app

**Acceptance**:
- Complete workflow executes successfully
- Uses mock deployment (fast)
- All steps validated
- Runs in CI without credentials

---

### ✅ Task 3.3: E2E Deployment Workflow Test (Real Fly.io)
**File**: `e2e/deployment_workflow_flyio_test.go`
**Status**: ⬜ Not Started
**Time**: 1 hour

**Implementation**:
- [ ] TestCompleteWorkflow_FlyIO_Real() with build tag
- [ ] RequireCredentials check
- [ ] Create test app with unique name
- [ ] Add posts resource
- [ ] Add authentication
- [ ] Generate Fly.io stack with litestream
- [ ] Actually deploy to Fly.io
- [ ] Wait for deployment to complete
- [ ] Run smoke tests against real URL
- [ ] Browser test against deployed app
- [ ] Cleanup (destroy app and volumes)

**Acceptance**:
- Actually deploys to Fly.io
- Smoke tests pass on live deployment
- Browser can interact with deployed app
- Cleanup destroys all resources
- Test skips without FLY_API_TOKEN

---

### ✅ Task 3.4: Docker Deployment Test
**File**: `e2e/deployment_docker_test.go`
**Status**: ⬜ Not Started
**Time**: 45 min

**Implementation**:
- [ ] TestDeployment_Docker() function
- [ ] Create test app
- [ ] Generate Docker stack
- [ ] Docker compose up
- [ ] Wait for containers healthy
- [ ] Run smoke tests against localhost
- [ ] Browser test
- [ ] Docker compose down (cleanup)

**Acceptance**:
- Deploys locally with Docker
- All containers start successfully
- Smoke tests pass
- Cleanup removes all containers

---

## Phase 4: CI/CD Integration

### ✅ Task 4.1: GitHub Actions Workflow
**File**: `.github/workflows/deployment.yml`
**Status**: ⬜ Not Started
**Time**: 45 min

**Implementation**:
- [ ] Nightly deployment test workflow
- [ ] Run mock tests on every PR
- [ ] Run real Fly.io tests nightly (scheduled)
- [ ] Use GitHub secrets for FLY_API_TOKEN
- [ ] Cleanup on failure (destroy test apps)
- [ ] Notify on failure (GitHub issues or Slack)
- [ ] Upload test logs as artifacts

**Acceptance**:
- Mock tests run on every PR
- Real tests run nightly
- Failures create notifications
- Test apps cleaned up properly

---

### ✅ Task 4.2: Test Documentation
**File**: `e2e/DEPLOYMENT_TESTING.md`
**Status**: ⬜ Not Started
**Time**: 30 min

**Implementation**:
- [ ] Overview of deployment testing
- [ ] How to run tests locally
- [ ] Setting up credentials
- [ ] Mock vs real test modes
- [ ] Troubleshooting guide
- [ ] Cost management tips
- [ ] Adding new deployment providers
- [ ] Example test output

**Acceptance**:
- Clear documentation for users
- Setup instructions work
- Troubleshooting covers common issues

---

## Progress Tracker

| Phase | Tasks | Completed | Status |
|-------|-------|-----------|--------|
| 1. Core Infrastructure | 3 | 0/3 | ⬜ Not Started |
| 2. Provider Implementations | 3 | 0/3 | ⬜ Not Started |
| 3. Verification & Testing | 4 | 0/4 | ⬜ Not Started |
| 4. CI/CD Integration | 2 | 0/2 | ⬜ Not Started |
| **TOTAL** | **12** | **0/12** | **0%** |

---

## Testing Commands

```bash
# Run all tests (uses mocks)
go test ./...

# Run e2e tests only (mocks)
go test ./e2e

# Run with real Fly.io deployments (requires FLY_API_TOKEN)
FLY_API_TOKEN=xxx go test -tags=deployment_integration ./e2e

# Run specific deployment test
go test -tags=deployment_integration -run TestCompleteWorkflow_FlyIO ./e2e

# Run Docker deployment tests
go test -run TestDeployment_Docker ./e2e

# Verbose output
go test -v ./e2e
```

---

## Success Criteria

- [x] Plan documented and tracked
- [ ] All 12 tasks completed
- [ ] Mock deployment tests pass
- [ ] Real Fly.io deployment tests pass
- [ ] Docker deployment tests pass
- [ ] Smoke tests verify deployed apps work
- [ ] Cleanup properly destroys resources
- [ ] Tests skip gracefully without credentials
- [ ] CI pipeline integrated
- [ ] Documentation complete

---

## Notes

- Mock tests should run fast (<30s total)
- Real Fly.io tests will be slower (2-5 min per deploy)
- Docker tests are medium speed (30-60s)
- All tests must cleanup resources (no cost leaks)
- Use `t.Cleanup()` for automatic cleanup on failure

---

## Timeline

**Estimated Total Time**: 7-8 hours
**Started**: 2025-11-05
**Target Completion**: TBD

Last Updated: 2025-11-05
