package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// DockerClient wraps docker CLI commands for local Docker deployments
type DockerClient struct {
	ImageTag      string
	ContainerName string
	Port          int
}

// NewDockerClient creates a new Docker client
func NewDockerClient(containerName string, port int) *DockerClient {
	return &DockerClient{
		ImageTag:      fmt.Sprintf("lvt-test-%s:latest", containerName),
		ContainerName: containerName,
		Port:          port,
	}
}

// DockerContainerStatus represents the status of a Docker container
type DockerContainerStatus struct {
	ID      string
	Name    string
	Status  string
	Running bool
	Port    int
}

// Build builds a Docker image from the app directory
func (d *DockerClient) Build(appDir string) error {
	args := []string{
		"build",
		"-t", d.ImageTag,
		".",
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = appDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker build failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Run starts a Docker container from the built image
func (d *DockerClient) Run() error {
	args := []string{
		"run",
		"-d",                                   // Detached mode
		"-p", fmt.Sprintf("%d:8080", d.Port),  // Port mapping
		"--name", d.ContainerName,             // Container name
		"-e", "PORT=8080",                     // Environment variable
		"-e", "APP_ENV=test",                  // Test environment
		d.ImageTag,                            // Image to run
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Status returns the status of the Docker container
func (d *DockerClient) Status() (*DockerContainerStatus, error) {
	args := []string{
		"inspect",
		d.ContainerName,
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker inspect failed: %w\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var containers []map[string]interface{}
	if err := json.Unmarshal(output, &containers); err != nil {
		return nil, fmt.Errorf("failed to parse inspect JSON: %w", err)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("container not found: %s", d.ContainerName)
	}

	container := containers[0]
	status := &DockerContainerStatus{
		Name: d.ContainerName,
		Port: d.Port,
	}

	// Extract ID
	if id, ok := container["Id"].(string); ok {
		status.ID = id
	}

	// Extract state
	if stateMap, ok := container["State"].(map[string]interface{}); ok {
		if statusStr, ok := stateMap["Status"].(string); ok {
			status.Status = statusStr
		}
		if running, ok := stateMap["Running"].(bool); ok {
			status.Running = running
		}
	}

	return status, nil
}

// GetContainerURL returns the URL to access the container
func (d *DockerClient) GetContainerURL() string {
	return fmt.Sprintf("http://localhost:%d", d.Port)
}

// WaitForReady waits for the container to be ready and responding to HTTP requests
func (d *DockerClient) WaitForReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	url := d.GetContainerURL() + "/health"

	for time.Now().Before(deadline) {
		// Check container status
		status, err := d.Status()
		if err == nil && status.Running {
			// Try HTTP health check
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

	// Get logs for debugging
	logs, _ := d.Logs(50)
	return fmt.Errorf("container %s did not become ready within %v\nLogs:\n%s", d.ContainerName, timeout, logs)
}

// Logs retrieves the last N lines of logs from the container
func (d *DockerClient) Logs(lines int) (string, error) {
	args := []string{
		"logs",
		"--tail", fmt.Sprintf("%d", lines),
		d.ContainerName,
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker logs failed: %w", err)
	}

	return string(output), nil
}

// Stop stops the Docker container
func (d *DockerClient) Stop() error {
	args := []string{
		"stop",
		d.ContainerName,
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Container might already be stopped, check if it exists
		if strings.Contains(string(output), "No such container") {
			return nil
		}
		return fmt.Errorf("docker stop failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Remove removes the Docker container
func (d *DockerClient) Remove() error {
	args := []string{
		"rm",
		"-f", // Force removal even if running
		d.ContainerName,
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Container might not exist
		if strings.Contains(string(output), "No such container") {
			return nil
		}
		return fmt.Errorf("docker rm failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// RemoveImage removes the Docker image
func (d *DockerClient) RemoveImage() error {
	args := []string{
		"rmi",
		"-f", // Force removal
		d.ImageTag,
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Image might not exist
		if strings.Contains(string(output), "No such image") {
			return nil
		}
		return fmt.Errorf("docker rmi failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Destroy removes both the container and the image
func (d *DockerClient) Destroy() error {
	// Stop and remove container
	if err := d.Stop(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	if err := d.Remove(); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	// Remove image
	if err := d.RemoveImage(); err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}

	return nil
}

// CheckDockerInstalled verifies that docker is installed and accessible
func CheckDockerInstalled() error {
	cmd := exec.Command("docker", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker not found or not executable: %w\nOutput: %s", err, string(output))
	}

	// Verify docker daemon is running
	cmd = exec.Command("docker", "ps")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon not running: %w", err)
	}

	return nil
}

// GetDockerVersion returns the installed docker version
func GetDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get docker version: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
