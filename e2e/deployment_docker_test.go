package e2e

import (
	"os"
	"testing"

	lvttesting "github.com/livetemplate/lvt/testing"
	"github.com/livetemplate/lvt/testing/providers"
)

// TestDockerInstalled checks if docker is installed and running (always runs)
func TestDockerInstalled(t *testing.T) {
	err := providers.CheckDockerInstalled()
	if err != nil {
		t.Skipf("docker not installed or not running: %v", err)
	}

	version, err := providers.GetDockerVersion()
	if err != nil {
		t.Fatalf("Failed to get docker version: %v", err)
	}

	t.Logf("docker version: %s", version)
}

// TestDockerClientCreation tests creating a Docker client
func TestDockerClientCreation(t *testing.T) {
	// This test doesn't require docker to be running, just tests client creation
	client := providers.NewDockerClient("test-container", 8080)

	if client == nil {
		t.Fatal("Failed to create Docker client")
	}

	if client.ContainerName != "test-container" {
		t.Errorf("Expected container name 'test-container', got '%s'", client.ContainerName)
	}

	if client.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", client.Port)
	}

	t.Log("Docker client created successfully")
}

// TestDockerDeployment tests actual deployment to Docker
// This test requires Docker to be installed and running
func TestDockerDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker deployment test in short mode")
	}

	// Check docker is installed and running
	if err := providers.CheckDockerInstalled(); err != nil {
		t.Skipf("docker not available: %v", err)
	}

	// Setup deployment with minimal test app
	opts := &lvttesting.DeploymentOptions{
		Provider: lvttesting.ProviderDocker,
		Kit:      "simple",
		// Let it generate a unique app name
	}

	dt := lvttesting.SetupDeployment(t, opts)
	t.Logf("Created test deployment: %s", dt.AppName)

	// Deploy to Docker
	t.Log("Starting deployment to Docker...")
	if err := dt.Deploy(); err != nil {
		t.Fatalf("Deployment failed: %v", err)
	}

	// Verify deployment
	t.Log("Verifying deployment...")
	if dt.AppURL == "" {
		t.Fatal("App URL not set after deployment")
	}
	t.Logf("App deployed at: %s", dt.AppURL)

	// Run smoke tests
	t.Log("Running smoke tests...")
	smokeOpts := lvttesting.DefaultSmokeTestOptions()
	smokeOpts.SkipBrowser = true // Skip browser tests for simplicity

	suite, err := lvttesting.RunSmokeTests(dt.AppURL, smokeOpts)
	if err != nil {
		t.Errorf("Smoke tests failed: %v", err)
	}

	if suite != nil {
		suite.PrintResults()
		if !suite.AllPassed() {
			t.Error("Some smoke tests failed")
		}
	}

	// Cleanup is handled automatically by t.Cleanup() registered in SetupDeployment
	t.Log("Test completed, cleanup will run automatically")
}

// TestDockerDeploymentWithResources tests deployment with generated resources
func TestDockerDeploymentWithResources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker deployment test in short mode")
	}

	// Only run if explicitly requested via environment variable
	if os.Getenv("RUN_DOCKER_DEPLOYMENT_TESTS") != "true" {
		t.Skip("Skipping Docker deployment test (set RUN_DOCKER_DEPLOYMENT_TESTS=true to enable)")
	}

	if err := providers.CheckDockerInstalled(); err != nil {
		t.Skipf("docker not available: %v", err)
	}

	opts := &lvttesting.DeploymentOptions{
		Provider:  lvttesting.ProviderDocker,
		Kit:       "multi",
		Resources: []string{"posts title:string content:text"},
	}

	dt := lvttesting.SetupDeployment(t, opts)
	t.Logf("Created test deployment with resources: %s", dt.AppName)

	if err := dt.Deploy(); err != nil {
		t.Fatalf("Deployment failed: %v", err)
	}

	t.Logf("App deployed successfully at: %s", dt.AppURL)

	// Run smoke tests
	smokeOpts := lvttesting.DefaultSmokeTestOptions()
	smokeOpts.SkipBrowser = true

	suite, err := lvttesting.RunSmokeTests(dt.AppURL, smokeOpts)
	if err != nil {
		t.Errorf("Smoke tests failed: %v", err)
	}

	if suite != nil && !suite.AllPassed() {
		suite.PrintResults()
		t.Error("Some smoke tests failed")
	}
}

// TestDockerDeploymentQuickSmoke is a fast smoke test for Docker deployment
// This test is skipped by default unless RUN_DOCKER_DEPLOYMENT_TESTS=true
func TestDockerDeploymentQuickSmoke(t *testing.T) {
	// Only run if explicitly requested
	if os.Getenv("RUN_DOCKER_DEPLOYMENT_TESTS") != "true" {
		t.Skip("Skipping Docker deployment test (set RUN_DOCKER_DEPLOYMENT_TESTS=true to enable)")
	}

	// Check docker availability
	if err := providers.CheckDockerInstalled(); err != nil {
		t.Skipf("docker not available: %v", err)
	}

	// Use simple kit for faster test
	opts := &lvttesting.DeploymentOptions{
		Provider: lvttesting.ProviderDocker,
		Kit:      "simple",
	}

	dt := lvttesting.SetupDeployment(t, opts)
	t.Logf("Created test deployment: %s", dt.AppName)

	// Deploy
	if err := dt.Deploy(); err != nil {
		t.Fatalf("Deployment failed: %v", err)
	}

	t.Logf("App deployed at: %s", dt.AppURL)

	// Just verify health endpoint
	if err := dt.VerifyHealth(); err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	t.Log("Quick smoke test passed")
}
