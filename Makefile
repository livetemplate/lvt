.PHONY: test-fast test-commit test-all test-e2e test-unit test-http test-browser test-clean install upgrade-skills

# =============================================================================
# TIERED TESTING STRATEGY
# =============================================================================
# Tier 0: Client JS tests (Jest + jsdom) - run separately via npm
# Tier 1: HTTP tests (no browser) - fast server-side validation
# Tier 2: Browser rendering tests - validates rendering library
# Tier 3: Application tests - HTTP for most, browser for UI-critical
# =============================================================================

# Fast feedback - unit tests only (~30 seconds)
test-fast:
	@echo "Running unit tests (Tier 0)..."
	GOWORK=off go test -short ./internal/... ./commands/... ./testing/...

# Before commit - unit + HTTP tests, NO browser (~75 seconds)
# This is the recommended pre-push target - no laptop freezing!
test-commit:
	@echo "Running commit validation tests (no browser)..."
	GOWORK=off go test -p 4 -timeout=10m ./internal/... ./commands/... ./testing/...
	@echo "Running HTTP-based e2e tests..."
	GOWORK=off go test -p 4 -timeout=10m -tags=http ./e2e/...

# HTTP tests only - Tier 1 + Tier 3 HTTP tests
test-http:
	@echo "Running HTTP-based tests (no browser)..."
	GOWORK=off go test -v -timeout=10m -tags=http ./e2e/...

# Browser rendering tests only - Tier 2 (~45 seconds)
# Run weekly or before releases to validate rendering library
test-browser:
	@echo "Running browser rendering tests (Tier 2)..."
	GOWORK=off go test -v -timeout=5m -tags=browser ./e2e/...

# Full validation - all tiers including browser (~3-5 minutes)
test-all:
	@echo "Running full test suite (all tiers)..."
	GOWORK=off go test -p 4 -timeout=15m -tags="http,browser,deployment" ./...

# Legacy e2e target - runs all e2e tests (browser-based)
test-e2e:
	@echo "Running all e2e tests (legacy, includes browser)..."
	GOWORK=off go test -v -timeout=10m ./e2e/...

# Unit tests only
test-unit:
	@echo "Running unit tests..."
	GOWORK=off go test ./internal/... ./commands/... ./testing/...

# Clean up Docker resources
test-clean:
	@echo "Cleaning up test Docker resources..."
	docker rm -f $$(docker ps -aq --filter "name=lvt-test-") 2>/dev/null || true
	docker rm -f $$(docker ps -aq --filter "name=chrome-pool-") 2>/dev/null || true
	docker rmi $$(docker images -q "lvt-test-*") 2>/dev/null || true
	docker rmi lvt-base:latest 2>/dev/null || true

# Install lvt CLI to GOPATH/bin
install:
	@echo "Installing lvt..."
	GOWORK=off go install .

# Install lvt and upgrade Claude Code skills (for local testing)
upgrade-skills: install
	@echo "Upgrading Claude Code skills..."
	lvt install-agent --llm claude --upgrade
