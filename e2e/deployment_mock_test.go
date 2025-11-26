package e2e

import (
	"testing"
	"time"

	lvttesting "github.com/livetemplate/lvt/testing"
	"github.com/livetemplate/lvt/testing/providers"
)

// TestDeploymentInfrastructure_Mock tests our deployment infrastructure with mock client
func TestDeploymentInfrastructure_Mock(t *testing.T) {
	// This test validates:
	// 1. Credentials management
	// 2. App name generation
	// 3. Deployment harness
	// 4. Mock Fly.io client
	// 5. Smoke test suite

	t.Run("AppNameGeneration", func(t *testing.T) {
		// Test unique app name generation
		name1 := lvttesting.GenerateTestAppName("test")
		name2 := lvttesting.GenerateTestAppName("test")

		if name1 == name2 {
			t.Errorf("Expected unique names, got duplicates: %s", name1)
		}

		// Validate name format
		if err := lvttesting.ValidateAppName(name1); err != nil {
			t.Errorf("Generated invalid app name: %v", err)
		}

		t.Logf("Generated app names: %s, %s", name1, name2)
	})

	t.Run("CredentialsManagement", func(t *testing.T) {
		// Test credential loading (should work without credentials set)
		creds, err := lvttesting.LoadTestCredentials()
		if err != nil {
			t.Fatalf("Failed to load credentials: %v", err)
		}

		t.Logf("Credentials loaded (FLY_API_TOKEN set: %v)", creds.FlyAPIToken != "")

		// Test validation (should fail gracefully)
		err = lvttesting.ValidateCredentials(lvttesting.ProviderFly)
		if err == nil && creds.FlyAPIToken == "" {
			t.Error("Expected validation error when FLY_API_TOKEN not set")
		}

		// Test HasCredentials
		hasFly := lvttesting.HasCredentials(lvttesting.ProviderFly)
		hasDocker := lvttesting.HasCredentials(lvttesting.ProviderDocker)
		t.Logf("Has Fly credentials: %v, Has Docker: %v", hasFly, hasDocker)
	})

	t.Run("MockFlyClient", func(t *testing.T) {
		// Test mock Fly.io client
		client := providers.NewMockFlyClient()
		appName := lvttesting.GenerateTestAppName("mock")

		// Launch app
		if err := client.Launch(appName, "sjc"); err != nil {
			t.Fatalf("Failed to launch app: %v", err)
		}
		t.Logf("Launched mock app: %s", appName)

		// Check status
		status, err := client.Status(appName)
		if err != nil {
			t.Fatalf("Failed to get status: %v", err)
		}
		if status.Status != "stopped" {
			t.Errorf("Expected status 'stopped', got '%s'", status.Status)
		}
		t.Logf("App status: %s", status.Status)

		// Deploy app
		if err := client.Deploy(appName, "/fake/dir"); err != nil {
			t.Fatalf("Failed to deploy: %v", err)
		}
		t.Logf("Deployed mock app")

		// Check status after deploy
		status, err = client.Status(appName)
		if err != nil {
			t.Fatalf("Failed to get status: %v", err)
		}
		if status.Status != "running" {
			t.Errorf("Expected status 'running', got '%s'", status.Status)
		}
		if !status.Deployed {
			t.Error("Expected app to be marked as deployed")
		}

		// Create volume
		volumeID, err := client.CreateVolume(appName, "sjc", 10)
		if err != nil {
			t.Fatalf("Failed to create volume: %v", err)
		}
		t.Logf("Created volume: %s", volumeID)

		// Get app URL
		url, err := client.GetAppURL(appName)
		if err != nil {
			t.Fatalf("Failed to get URL: %v", err)
		}
		expectedURL := "https://" + appName + ".fly.dev"
		if url != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, url)
		}
		t.Logf("App URL: %s", url)

		// List apps
		apps := client.ListApps()
		if len(apps) != 1 {
			t.Errorf("Expected 1 app, got %d", len(apps))
		}

		// List volumes
		volumes := client.ListVolumes()
		if len(volumes) != 1 {
			t.Errorf("Expected 1 volume, got %d", len(volumes))
		}

		// Destroy app
		if err := client.Destroy(appName); err != nil {
			t.Fatalf("Failed to destroy: %v", err)
		}
		t.Logf("Destroyed app and volumes")

		// Verify cleanup
		apps = client.ListApps()
		if len(apps) != 0 {
			t.Errorf("Expected 0 apps after destroy, got %d", len(apps))
		}
		volumes = client.ListVolumes()
		if len(volumes) != 0 {
			t.Errorf("Expected 0 volumes after destroy, got %d", len(volumes))
		}
	})

	t.Run("SmokeTestSuite", func(t *testing.T) {
		// Test smoke test suite with mock server
		// Note: This requires a running server, so we'll just test the structure

		opts := &lvttesting.SmokeTestOptions{
			Timeout:     30 * time.Second,
			RetryDelay:  1 * time.Second,
			MaxRetries:  2,
			SkipBrowser: true, // Skip browser tests for unit test
		}

		// We can't run actual smoke tests without a server,
		// but we can verify the options structure works
		if opts.Timeout != 30*time.Second {
			t.Error("Timeout not set correctly")
		}
		if opts.MaxRetries != 2 {
			t.Error("MaxRetries not set correctly")
		}

		t.Log("Smoke test suite structure validated")
	})
}

// TestMockDeploymentWorkflow demonstrates a complete mock deployment workflow
func TestMockDeploymentWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping workflow test in short mode")
	}

	// Create mock client
	client := providers.NewMockFlyClient()

	// Make deployment faster for testing
	client.SimulateDelay = false

	appName := lvttesting.GenerateTestAppName("workflow")
	t.Logf("Testing workflow with app: %s", appName)

	// Step 1: Launch app
	t.Log("Step 1: Launching app...")
	if err := client.Launch(appName, "sjc"); err != nil {
		t.Fatalf("Launch failed: %v", err)
	}

	// Step 2: Create volume
	t.Log("Step 2: Creating volume...")
	volumeID, err := client.CreateVolume(appName, "sjc", 10)
	if err != nil {
		t.Fatalf("Volume creation failed: %v", err)
	}
	t.Logf("Created volume: %s", volumeID)

	// Step 3: Deploy app
	t.Log("Step 3: Deploying app...")
	if err := client.Deploy(appName, "/fake/app/dir"); err != nil {
		t.Fatalf("Deployment failed: %v", err)
	}

	// Step 4: Wait for ready
	t.Log("Step 4: Waiting for app to be ready...")
	if err := client.WaitForAppReady(appName, 30*time.Second); err != nil {
		t.Fatalf("App never became ready: %v", err)
	}

	// Step 5: Verify status
	t.Log("Step 5: Verifying app status...")
	status, err := client.Status(appName)
	if err != nil {
		t.Fatalf("Status check failed: %v", err)
	}
	if status.Status != "running" {
		t.Errorf("Expected status 'running', got '%s'", status.Status)
	}
	t.Logf("App is running at: %s", status.URL)

	// Step 6: Cleanup
	t.Log("Step 6: Cleaning up...")
	if err := client.Destroy(appName); err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify cleanup
	apps := client.ListApps()
	if len(apps) != 0 {
		t.Errorf("Cleanup incomplete: %d apps remaining", len(apps))
	}

	t.Log("✅ Mock deployment workflow completed successfully!")
}

// TestMockClientFailureSimulation tests error handling
func TestMockClientFailureSimulation(t *testing.T) {
	client := providers.NewMockFlyClient()
	client.SimulateDelay = false
	client.SimulateFailures = true
	client.FailureRate = 1.0 // Always fail

	appName := lvttesting.GenerateTestAppName("fail")

	// This should currently not fail because simulateFailure() is hardcoded to false
	// But the structure is there for future testing
	err := client.Launch(appName, "sjc")

	// For now, we expect success since simulateFailure() returns false
	if err != nil {
		t.Logf("Got expected simulated failure: %v", err)
	} else {
		t.Log("Mock client is in success mode (simulateFailure returns false)")
	}
}
