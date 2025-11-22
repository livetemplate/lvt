package e2e

import (
	"os"
	"testing"

	lvttesting "github.com/livetemplate/lvt/testing"
	"github.com/livetemplate/lvt/testing/providers"
)

// TestFlyctlInstalled checks if flyctl is installed (always runs)
func TestFlyctlInstalled(t *testing.T) {
	err := providers.CheckFlyctlInstalled()
	if err != nil {
		t.Skipf("flyctl not installed or not in PATH: %v", err)
	}

	version, err := providers.GetFlyctlVersion()
	if err != nil {
		t.Fatalf("Failed to get flyctl version: %v", err)
	}

	t.Logf("flyctl version: %s", version)
}

// TestFlyClientCreation tests creating a Fly.io client
func TestFlyClientCreation(t *testing.T) {
	// This test doesn't require credentials, just tests client creation
	client := providers.NewFlyClient("test-token", "test-org")

	if client == nil {
		t.Fatal("Failed to create Fly.io client")
	}

	if client.APIToken != "test-token" {
		t.Errorf("Expected API token 'test-token', got '%s'", client.APIToken)
	}

	if client.OrgSlug != "test-org" {
		t.Errorf("Expected org slug 'test-org', got '%s'", client.OrgSlug)
	}

	t.Log("Fly.io client created successfully")
}

// TestRealFlyDeployment tests actual deployment to Fly.io
// This test requires FLY_API_TOKEN to be set and will be skipped otherwise
func TestRealFlyDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real deployment test in short mode")
	}

	// Require Fly.io credentials
	lvttesting.RequireFlyCredentials(t)

	// Check flyctl is installed
	if err := providers.CheckFlyctlInstalled(); err != nil {
		t.Skipf("flyctl not installed: %v", err)
	}

	// Setup deployment with minimal test app
	opts := &lvttesting.DeploymentOptions{
		Provider: lvttesting.ProviderFly,
		Region:   "sjc", // San Jose
		Kit:      "simple",
		// Let it generate a unique app name
	}

	dt := lvttesting.SetupDeployment(t, opts)
	t.Logf("Created test deployment: %s", dt.AppName)

	// Deploy to Fly.io
	t.Log("Starting deployment to Fly.io...")
	if err := dt.Deploy(); err != nil {
		t.Fatalf("Deployment failed: %v", err)
	}

	// Verify health
	t.Log("Verifying deployment health...")
	if dt.AppURL == "" {
		t.Fatal("App URL not set after deployment")
	}
	t.Logf("App deployed at: %s", dt.AppURL)

	// Run smoke tests
	t.Log("Running smoke tests...")
	smokeOpts := lvttesting.DefaultSmokeTestOptions()
	smokeOpts.SkipBrowser = true // Skip browser tests in CI

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

// TestFlyDeploymentWithResources tests deployment with generated resources
func TestFlyDeploymentWithResources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real deployment test in short mode")
	}

	// Only run if explicitly requested via environment variable
	if os.Getenv("RUN_FLY_DEPLOYMENT_TESTS") != "true" {
		t.Skip("Skipping Fly.io deployment test (set RUN_FLY_DEPLOYMENT_TESTS=true to enable)")
	}

	lvttesting.RequireFlyCredentials(t)

	if err := providers.CheckFlyctlInstalled(); err != nil {
		t.Skipf("flyctl not installed: %v", err)
	}

	opts := &lvttesting.DeploymentOptions{
		Provider:  lvttesting.ProviderFly,
		Region:    "sjc",
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
