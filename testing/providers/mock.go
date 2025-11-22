package providers

import (
	"fmt"
	"time"
)

// MockFlyClient simulates Fly.io API calls for fast testing without actual deployments
type MockFlyClient struct {
	// Simulated state
	Apps    map[string]*MockAppStatus
	Volumes map[string]*MockVolume

	// Configuration
	SimulateDelay    bool          // If true, adds realistic delays
	SimulateFailures bool          // If true, randomly fails operations
	FailureRate      float64       // Probability of failure (0.0 - 1.0)
	DeployDuration   time.Duration // How long to simulate deployment
}

// MockAppStatus represents the status of a mock app
type MockAppStatus struct {
	Name      string
	Status    string // "running", "stopped", "deploying"
	URL       string
	Region    string
	CreatedAt time.Time
	Deployed  bool
}

// MockVolume represents a mock volume
type MockVolume struct {
	ID     string
	Name   string
	Region string
	SizeGB int
	AppID  string
}

// NewMockFlyClient creates a new mock Fly.io client
func NewMockFlyClient() *MockFlyClient {
	return &MockFlyClient{
		Apps:           make(map[string]*MockAppStatus),
		Volumes:        make(map[string]*MockVolume),
		SimulateDelay:  true,
		DeployDuration: 2 * time.Second, // Fast mock deployment
	}
}

// Launch creates a new mock app
func (m *MockFlyClient) Launch(appName, region string) error {
	if m.simulateFailure() {
		return fmt.Errorf("mock failure: unable to launch app")
	}

	if m.SimulateDelay {
		time.Sleep(500 * time.Millisecond)
	}

	if _, exists := m.Apps[appName]; exists {
		return fmt.Errorf("app %s already exists", appName)
	}

	m.Apps[appName] = &MockAppStatus{
		Name:      appName,
		Status:    "stopped",
		URL:       fmt.Sprintf("https://%s.fly.dev", appName),
		Region:    region,
		CreatedAt: time.Now(),
		Deployed:  false,
	}

	return nil
}

// Deploy simulates deploying an app
func (m *MockFlyClient) Deploy(appName, appDir string) error {
	if m.simulateFailure() {
		return fmt.Errorf("mock failure: deployment failed")
	}

	app, exists := m.Apps[appName]
	if !exists {
		return fmt.Errorf("app %s not found", appName)
	}

	// Simulate deployment time
	if m.SimulateDelay {
		time.Sleep(m.DeployDuration)
	}

	app.Status = "running"
	app.Deployed = true

	return nil
}

// Status returns the status of a mock app
func (m *MockFlyClient) Status(appName string) (*MockAppStatus, error) {
	if m.simulateFailure() {
		return nil, fmt.Errorf("mock failure: unable to get status")
	}

	if m.SimulateDelay {
		time.Sleep(200 * time.Millisecond)
	}

	app, exists := m.Apps[appName]
	if !exists {
		return nil, fmt.Errorf("app %s not found", appName)
	}

	return app, nil
}

// CreateVolume creates a mock volume
func (m *MockFlyClient) CreateVolume(appName, region string, sizeGB int) (string, error) {
	if m.simulateFailure() {
		return "", fmt.Errorf("mock failure: unable to create volume")
	}

	if m.SimulateDelay {
		time.Sleep(1 * time.Second)
	}

	volumeID := fmt.Sprintf("vol_%s_%d", appName, time.Now().Unix())
	m.Volumes[volumeID] = &MockVolume{
		ID:     volumeID,
		Name:   fmt.Sprintf("%s-data", appName),
		Region: region,
		SizeGB: sizeGB,
		AppID:  appName,
	}

	return volumeID, nil
}

// Destroy removes a mock app and its volumes
func (m *MockFlyClient) Destroy(appName string) error {
	if m.simulateFailure() {
		return fmt.Errorf("mock failure: unable to destroy app")
	}

	if m.SimulateDelay {
		time.Sleep(500 * time.Millisecond)
	}

	if _, exists := m.Apps[appName]; !exists {
		return fmt.Errorf("app %s not found", appName)
	}

	// Remove app
	delete(m.Apps, appName)

	// Remove associated volumes
	for id, vol := range m.Volumes {
		if vol.AppID == appName {
			delete(m.Volumes, id)
		}
	}

	return nil
}

// GetAppURL returns the mock URL for an app
func (m *MockFlyClient) GetAppURL(appName string) (string, error) {
	app, exists := m.Apps[appName]
	if !exists {
		return "", fmt.Errorf("app %s not found", appName)
	}

	return app.URL, nil
}

// WaitForAppReady simulates waiting for app to be ready
func (m *MockFlyClient) WaitForAppReady(appName string, timeout time.Duration) error {
	if m.simulateFailure() {
		return fmt.Errorf("mock failure: app never became ready")
	}

	app, exists := m.Apps[appName]
	if !exists {
		return fmt.Errorf("app %s not found", appName)
	}

	// Simulate deployment completing
	if m.SimulateDelay {
		time.Sleep(100 * time.Millisecond)
	}

	app.Status = "running"
	return nil
}

// simulateFailure returns true if a failure should be simulated
func (m *MockFlyClient) simulateFailure() bool {
	if !m.SimulateFailures {
		return false
	}
	// Simple deterministic failure for testing
	// In real usage, you might use rand.Float64() < m.FailureRate
	return false
}

// Reset clears all mock state
func (m *MockFlyClient) Reset() {
	m.Apps = make(map[string]*MockAppStatus)
	m.Volumes = make(map[string]*MockVolume)
}

// ListApps returns all mock apps
func (m *MockFlyClient) ListApps() []*MockAppStatus {
	apps := make([]*MockAppStatus, 0, len(m.Apps))
	for _, app := range m.Apps {
		apps = append(apps, app)
	}
	return apps
}

// ListVolumes returns all mock volumes
func (m *MockFlyClient) ListVolumes() []*MockVolume {
	volumes := make([]*MockVolume, 0, len(m.Volumes))
	for _, vol := range m.Volumes {
		volumes = append(volumes, vol)
	}
	return volumes
}
