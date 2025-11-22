.PHONY: test-fast test-commit test-all test-e2e test-unit test-clean

# Fast feedback - unit tests only (~30 seconds)
test-fast:
	@echo "Running unit tests..."
	GOWORK=off go test -short ./internal/... ./commands/...

# Before commit - unit + quick e2e (~3-4 minutes)
test-commit:
	@echo "Running commit validation tests..."
	GOWORK=off go test -p 4 -timeout=20m ./...

# Full validation - all tests including deployment (~5-6 minutes)
test-all:
	@echo "Running full test suite..."
	GOWORK=off go test -p 4 -timeout=10m -tags=deployment ./...

# E2E only
test-e2e:
	@echo "Running e2e tests..."
	GOWORK=off go test -v -timeout=10m ./e2e/...

# Unit tests only
test-unit:
	@echo "Running unit tests..."
	GOWORK=off go test ./internal/... ./commands/...

# Clean up Docker resources
test-clean:
	@echo "Cleaning up test Docker resources..."
	docker rm -f $$(docker ps -aq --filter "name=lvt-test-") 2>/dev/null || true
	docker rm -f $$(docker ps -aq --filter "name=chrome-pool-") 2>/dev/null || true
	docker rmi $$(docker images -q "lvt-test-*") 2>/dev/null || true
	docker rmi lvt-base:latest 2>/dev/null || true
