#!/bin/bash
# Run deployment tests inside a Docker container
# This avoids macOS permission issues and provides consistent environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building test container...${NC}"
docker build -f Dockerfile.test -t lvt-test-runner .

# Default test pattern
TEST_PATTERN="${1:-TestDeploymentInfrastructure_Mock}"
TEST_FLAGS="${2:--v -timeout=5m}"

echo -e "${GREEN}Running tests: ${TEST_PATTERN}${NC}"
echo -e "${YELLOW}Mounting Docker socket for deployment tests...${NC}"

# Run tests with Docker socket mounted
# This allows the container to build Docker images for deployment tests
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e RUN_DOCKER_DEPLOYMENT_TESTS="${RUN_DOCKER_DEPLOYMENT_TESTS:-false}" \
  -e RUN_FLY_DEPLOYMENT_TESTS="${RUN_FLY_DEPLOYMENT_TESTS:-false}" \
  -e FLY_API_TOKEN="${FLY_API_TOKEN:-}" \
  lvt-test-runner \
  go test ./e2e -run "${TEST_PATTERN}" ${TEST_FLAGS}

echo -e "${GREEN}✅ Tests completed${NC}"
