package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// FlyClient wraps flyctl CLI commands for Fly.io deployments
type FlyClient struct {
	APIToken string
	OrgSlug  string
}

// NewFlyClient creates a new Fly.io client
func NewFlyClient(apiToken, orgSlug string) *FlyClient {
	return &FlyClient{
		APIToken: apiToken,
		OrgSlug:  orgSlug,
	}
}

// FlyAppStatus represents the status of a Fly.io app
type FlyAppStatus struct {
	Name       string
	Status     string
	Hostname   string
	Version    int
	Deployed   bool
	Allocation string
}

// Launch creates a new Fly.io app
func (f *FlyClient) Launch(appName, region string) error {
	args := []string{
		"apps", "create", appName,
		"--org", f.OrgSlug,
	}

	if region != "" {
		// Note: region is set during deployment, not during app creation
		// We'll store it for later use
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flyctl apps create failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Deploy deploys an app to Fly.io
func (f *FlyClient) Deploy(appName, appDir string, region string) error {
	args := []string{
		"deploy",
		"--app", appName,
		"--config", appDir + "/fly.toml",
	}

	if region != "" {
		args = append(args, "--region", region)
	}

	cmd := f.buildCommand(args...)
	cmd.Dir = appDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flyctl deploy failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Status returns the status of a Fly.io app
func (f *FlyClient) Status(appName string) (*FlyAppStatus, error) {
	args := []string{
		"status",
		"--app", appName,
		"--json",
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("flyctl status failed: %w\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse status JSON: %w", err)
	}

	status := &FlyAppStatus{
		Name:     appName,
		Deployed: true,
	}

	// Extract relevant fields from JSON
	if hostname, ok := result["Hostname"].(string); ok {
		status.Hostname = hostname
	}
	if statusStr, ok := result["Status"].(string); ok {
		status.Status = statusStr
	}

	return status, nil
}

// CreateVolume creates a volume for a Fly.io app
func (f *FlyClient) CreateVolume(appName, region string, sizeGB int) (string, error) {
	volumeName := fmt.Sprintf("%s_data", appName)

	args := []string{
		"volumes", "create", volumeName,
		"--app", appName,
		"--region", region,
		"--size", fmt.Sprintf("%d", sizeGB),
		"--yes", // Auto-confirm
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("flyctl volumes create failed: %w\nOutput: %s", err, string(output))
	}

	// Extract volume ID from output
	// The output typically contains the volume ID
	return volumeName, nil
}

// GetAppURL returns the URL for a Fly.io app
func (f *FlyClient) GetAppURL(appName string) (string, error) {
	status, err := f.Status(appName)
	if err != nil {
		return "", err
	}

	if status.Hostname != "" {
		return "https://" + status.Hostname, nil
	}

	// Default URL format
	return fmt.Sprintf("https://%s.fly.dev", appName), nil
}

// WaitForAppReady waits for a Fly.io app to be ready
func (f *FlyClient) WaitForAppReady(appName string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		status, err := f.Status(appName)
		if err == nil && status.Status == "running" {
			return nil
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("app %s did not become ready within %v", appName, timeout)
}

// Destroy removes a Fly.io app and all its resources
func (f *FlyClient) Destroy(appName string) error {
	// First, list and delete all volumes
	if err := f.destroyVolumes(appName); err != nil {
		// Log error but continue with app deletion
		fmt.Printf("Warning: failed to destroy volumes: %v\n", err)
	}

	// Delete the app
	args := []string{
		"apps", "destroy", appName,
		"--yes", // Auto-confirm
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flyctl apps destroy failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// destroyVolumes removes all volumes associated with an app
func (f *FlyClient) destroyVolumes(appName string) error {
	// List volumes
	args := []string{
		"volumes", "list",
		"--app", appName,
		"--json",
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flyctl volumes list failed: %w\nOutput: %s", err, string(output))
	}

	// Parse volume list
	var volumes []map[string]interface{}
	if err := json.Unmarshal(output, &volumes); err != nil {
		// If parsing fails, it might be because there are no volumes
		return nil
	}

	// Delete each volume
	for _, vol := range volumes {
		volumeID, ok := vol["id"].(string)
		if !ok {
			continue
		}

		deleteArgs := []string{
			"volumes", "delete", volumeID,
			"--app", appName,
			"--yes",
		}

		deleteCmd := f.buildCommand(deleteArgs...)
		if _, err := deleteCmd.CombinedOutput(); err != nil {
			fmt.Printf("Warning: failed to delete volume %s: %v\n", volumeID, err)
		}
	}

	return nil
}

// ListApps returns all apps in the organization
func (f *FlyClient) ListApps() ([]string, error) {
	args := []string{
		"apps", "list",
		"--json",
	}

	cmd := f.buildCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("flyctl apps list failed: %w\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var apps []map[string]interface{}
	if err := json.Unmarshal(output, &apps); err != nil {
		return nil, fmt.Errorf("failed to parse apps JSON: %w", err)
	}

	var appNames []string
	for _, app := range apps {
		if name, ok := app["Name"].(string); ok {
			appNames = append(appNames, name)
		}
	}

	return appNames, nil
}

// buildCommand creates an exec.Cmd with authentication set up
func (f *FlyClient) buildCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("flyctl", args...)

	// Set FLY_API_TOKEN environment variable
	if f.APIToken != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("FLY_API_TOKEN=%s", f.APIToken))
	}

	return cmd
}

// CheckFlyctlInstalled verifies that flyctl is installed and accessible
func CheckFlyctlInstalled() error {
	cmd := exec.Command("flyctl", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flyctl not found or not executable: %w\nOutput: %s", err, string(output))
	}

	// Verify output contains version info
	if !strings.Contains(string(output), "flyctl") {
		return fmt.Errorf("unexpected flyctl version output: %s", string(output))
	}

	return nil
}

// GetFlyctlVersion returns the installed flyctl version
func GetFlyctlVersion() (string, error) {
	cmd := exec.Command("flyctl", "version")
	// Use Output() to properly close pipes and avoid I/O wait
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get flyctl version: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
