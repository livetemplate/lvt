package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// envWithGOWORKOff returns a copy of the current environment with GOWORK
// reliably set to "off", filtering out any pre-existing GOWORK entries to
// avoid platform-dependent duplicate-key behavior.
func envWithGOWORKOff() []string {
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GOWORK=") {
			env = append(env, e)
		}
	}
	return append(env, "GOWORK=off")
}

// ValidateCompilation runs sqlc generate (if applicable), go mod tidy, and
// go build ./... in the given app directory. This ensures that generated code
// actually compiles and catches syntax errors, type mismatches, and missing
// dependencies before they reach production.
func ValidateCompilation(t *testing.T, appDir string) {
	t.Helper()

	env := envWithGOWORKOff()

	// Run sqlc generate if sqlc.yaml exists, since generated handlers
	// depend on sqlc-generated query code.
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err == nil {
		t.Log("Running sqlc generate...")
		sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
		sqlcCmd.Dir = appDir
		sqlcCmd.Env = env
		if output, err := sqlcCmd.CombinedOutput(); err != nil {
			t.Fatalf("sqlc generate failed in %s: %v\nOutput: %s", appDir, err, output)
		}
	}

	// Run go mod tidy to catch dependency problems introduced by generated code.
	t.Log("Running go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = appDir
	tidyCmd.Env = env
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed in %s: %v\nOutput: %s", appDir, err, output)
	}

	// Run go build ./... to validate that all generated code compiles.
	t.Log("Validating compilation with go build ./...")
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = appDir
	buildCmd.Env = env
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Compilation validation failed in %s: %v\nOutput: %s", appDir, err, output)
	}

	t.Log("Compilation validation passed")
}
