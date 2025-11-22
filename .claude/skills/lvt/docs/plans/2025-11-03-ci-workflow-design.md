# CI/CD Workflow Design for LVT

**Date:** 2025-11-03
**Status:** Approved
**Purpose:** Add GitHub Actions workflow to verify tests pass on all pull requests

## Requirements

### Triggers
- Run on pull requests: `opened`, `synchronize`, `reopened`
- Do not run on direct main branch pushes (saves CI minutes)

### Checks Required
- Code formatting verification (`go fmt`)
- Linting with golangci-lint (matching pre-commit hook config)
- Full test suite including chromedp E2E tests
- Must block PR merging if any check fails

### Constraints
- Use latest stable Go version
- Run chromedp tests with Docker headless Chrome
- Match pre-commit hook behavior (fmt, lint, test with 120s timeout)

## Design

### Workflow Structure

**Single comprehensive job** running steps sequentially:
1. Checkout code
2. Setup Go environment
3. Start Chrome container for E2E tests
4. Install dependencies
5. Verify code formatting
6. Run linting
7. Execute all tests

**Rationale:** Simple, maintainable, easy to debug. Sequential execution ensures clear failure points.

### Implementation Details

#### Job Configuration
```yaml
name: test
runs-on: ubuntu-latest
permissions:
  contents: read
  pull-requests: write
```

#### Step 1: Checkout
- Action: `actions/checkout@v4`
- Fetch full history for git operations

#### Step 2: Setup Go
- Action: `actions/setup-go@v5`
- Version: `'stable'` (auto-updates to latest stable)
- Cache: `true` (speeds up module downloads)

#### Step 3: Chrome Container
- Start: `docker run -d -p 9222:9222 --name chrome chromedp/headless-shell:latest`
- Environment: `CHROMEDP_REMOTE_URL=ws://localhost:9222`
- Purpose: Provides browser for chromedp E2E tests

#### Step 4: Install Dependencies
- Command: `go mod download`
- Ensures all modules available before checks

#### Step 5: Format Check
- Command: `go fmt ./...`
- Verify: Check if files modified after formatting
- Fail if unformatted code detected

#### Step 6: Linting
- Action: `golangci/golangci-lint-action@v6`
- Linters: `errcheck`, `unused`, `staticcheck`, `gosimple`, `ineffassign`
- Matches pre-commit hook configuration

#### Step 7: Run Tests
- Command: `go test -v ./... -timeout=120s`
- Includes unit, integration, and E2E tests
- Timeout matches pre-commit hook

### Error Handling

- Each step fails fast with non-zero exit code
- Clear error messages for each failure type:
  - Formatting: Shows which files need formatting
  - Linting: Displays specific issues
  - Tests: Shows test failures with verbose output
- Chrome container failures caught before tests run

### Branch Protection

**Post-deployment setup:**
1. Navigate to Settings → Branches → Branch protection rules
2. Add rule for `main` branch
3. Enable "Require status checks to pass before merging"
4. Select "test" workflow as required check
5. Enable "Require branches to be up to date before merging"

This prevents merging PRs until CI passes.

### Maintenance

- Go version: `'stable'` auto-updates to latest
- golangci-lint: Action manages linter version updates
- Chrome: Uses `latest` tag (can pin specific version if stability issues arise)

## Timeline

- **Implementation:** ~15 minutes
  - Create `.github/workflows/test.yml`
  - Commit and push to new branch
  - Open PR to test workflow
  - Verify all checks pass

- **Branch Protection:** ~5 minutes
  - Configure after first successful run
  - Verify PR blocking behavior

## Success Criteria

- ✅ Workflow runs on every PR
- ✅ All three checks execute (fmt, lint, test)
- ✅ chromedp E2E tests pass with Docker Chrome
- ✅ PRs cannot merge with failing tests
- ✅ Workflow completes in under 8 minutes

## Alternatives Considered

### Parallel Jobs
- **Pros:** Faster feedback (~3-5 min)
- **Cons:** More complex, uses more CI resources
- **Decision:** Rejected - simplicity preferred for single-maintainer project

### Matrix Strategy (Multiple Go Versions)
- **Pros:** Tests compatibility across versions
- **Cons:** 2x CI time and resources, overkill for CLI tool
- **Decision:** Rejected - CLI tracks core library version, not library itself
