package testing

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateTestAppName generates a unique app name for deployment testing
//
// Format: lvt-test-{random}-{timestamp}
// Example: lvt-test-a3f4b2-1699123456
//
// The name is guaranteed to:
// - Be unique (random + timestamp)
// - Be valid for Fly.io (lowercase, alphanumeric, hyphens only)
// - Be traceable (includes "test" and timestamp)
// - Be short enough for Fly.io limits (<32 chars)
func GenerateTestAppName(prefix string) string {
	if prefix == "" {
		prefix = "lvt-test"
	}

	// Generate 6 random hex characters (3 bytes)
	randomBytes := make([]byte, 3)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp-based random if crypto/rand fails
		randomBytes = []byte{
			byte(time.Now().UnixNano() & 0xFF),
			byte((time.Now().UnixNano() >> 8) & 0xFF),
			byte((time.Now().UnixNano() >> 16) & 0xFF),
		}
	}
	randomStr := hex.EncodeToString(randomBytes)

	// Use Unix timestamp (10 digits)
	timestamp := time.Now().Unix()

	// Format: prefix-random-timestamp
	name := fmt.Sprintf("%s-%s-%d", prefix, randomStr, timestamp)

	// Ensure it's valid for Fly.io
	name = strings.ToLower(name)
	name = sanitizeAppName(name)

	return name
}

// sanitizeAppName ensures the name is valid for Fly.io
// - Lowercase only
// - Alphanumeric and hyphens only
// - No leading/trailing hyphens
// - Max 32 characters
func sanitizeAppName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	name = reg.ReplaceAllString(name, "-")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Remove leading/trailing hyphens
	name = strings.Trim(name, "-")

	// Truncate to 32 characters (Fly.io limit)
	if len(name) > 32 {
		name = name[:32]
		name = strings.TrimRight(name, "-")
	}

	return name
}

// ValidateAppName checks if an app name is valid for Fly.io
func ValidateAppName(name string) error {
	if name == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if len(name) > 32 {
		return fmt.Errorf("app name too long (max 32 characters): %s", name)
	}

	// Must be lowercase alphanumeric and hyphens only
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(name) {
		return fmt.Errorf("app name must be lowercase alphanumeric with hyphens only: %s", name)
	}

	// Cannot start or end with hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("app name cannot start or end with hyphen: %s", name)
	}

	return nil
}

// IsTestAppName checks if a name appears to be a test app
func IsTestAppName(name string) bool {
	return strings.HasPrefix(name, "lvt-test-") ||
		strings.HasPrefix(name, "test-") ||
		strings.Contains(name, "-test-")
}
