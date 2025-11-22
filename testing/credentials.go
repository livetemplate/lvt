package testing

import (
	"fmt"
	"os"
	stdtesting "testing"
)

// Provider represents a deployment provider
type Provider string

const (
	ProviderFly          Provider = "fly"
	ProviderDocker       Provider = "docker"
	ProviderKubernetes   Provider = "kubernetes"
	ProviderDigitalOcean Provider = "digitalocean"
)

// TestCredentials holds all credentials needed for deployment testing
type TestCredentials struct {
	// Fly.io credentials
	FlyAPIToken string

	// AWS credentials (for litestream S3 backups)
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	S3Bucket           string
	S3Region           string

	// DigitalOcean credentials
	DOAPIToken string

	// Kubernetes credentials (optional, uses local kubeconfig by default)
	KubeConfig string
}

// LoadTestCredentials loads credentials from environment variables
func LoadTestCredentials() (*TestCredentials, error) {
	creds := &TestCredentials{
		FlyAPIToken:        os.Getenv("FLY_API_TOKEN"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		S3Region:           os.Getenv("S3_REGION"),
		DOAPIToken:         os.Getenv("DO_API_TOKEN"),
		KubeConfig:         os.Getenv("KUBECONFIG"),
	}

	// Set defaults
	if creds.S3Region == "" {
		creds.S3Region = "us-east-1"
	}

	return creds, nil
}

// ValidateCredentials checks if required credentials are available for a provider
func ValidateCredentials(provider Provider) error {
	creds, err := LoadTestCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	switch provider {
	case ProviderFly:
		if creds.FlyAPIToken == "" {
			return fmt.Errorf("FLY_API_TOKEN environment variable not set")
		}

	case ProviderDocker:
		// Docker doesn't require credentials (uses local daemon)
		return nil

	case ProviderKubernetes:
		// Kubernetes uses kubeconfig by default, optional override
		// No strict validation needed here
		return nil

	case ProviderDigitalOcean:
		if creds.DOAPIToken == "" {
			return fmt.Errorf("DO_API_TOKEN environment variable not set")
		}

	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}

	return nil
}

// RequireCredentials skips the test if credentials for the provider are not available
func RequireCredentials(t *stdtesting.T, provider Provider) {
	t.Helper()

	if err := ValidateCredentials(provider); err != nil {
		t.Skipf("Skipping test: %v", err)
	}
}

// RequireFlyCredentials is a convenience function for Fly.io tests
func RequireFlyCredentials(t *stdtesting.T) {
	t.Helper()
	RequireCredentials(t, ProviderFly)
}

// RequireAWSCredentials checks if AWS credentials are available (for litestream)
func RequireAWSCredentials(t *stdtesting.T) {
	t.Helper()

	creds, err := LoadTestCredentials()
	if err != nil {
		t.Skipf("Failed to load credentials: %v", err)
	}

	if creds.AWSAccessKeyID == "" || creds.AWSSecretAccessKey == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY required")
	}

	if creds.S3Bucket == "" {
		t.Skip("Skipping test: S3_BUCKET required for litestream backups")
	}
}

// HasCredentials returns true if credentials are available for the provider (non-testing version)
func HasCredentials(provider Provider) bool {
	return ValidateCredentials(provider) == nil
}

// GetFlyAPIToken returns the Fly.io API token from environment
func GetFlyAPIToken() (string, error) {
	token := os.Getenv("FLY_API_TOKEN")
	if token == "" {
		return "", fmt.Errorf("FLY_API_TOKEN environment variable not set")
	}
	return token, nil
}

// GetAWSCredentials returns AWS credentials for litestream
func GetAWSCredentials() (accessKeyID, secretAccessKey, bucket, region string, err error) {
	creds, err := LoadTestCredentials()
	if err != nil {
		return "", "", "", "", err
	}

	if creds.AWSAccessKeyID == "" || creds.AWSSecretAccessKey == "" {
		return "", "", "", "", fmt.Errorf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY required")
	}

	if creds.S3Bucket == "" {
		return "", "", "", "", fmt.Errorf("S3_BUCKET required")
	}

	return creds.AWSAccessKeyID, creds.AWSSecretAccessKey, creds.S3Bucket, creds.S3Region, nil
}
