package helpers

import (
	"bufio"
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

// ValidationOptions controls optional steps in ValidateCompilation.
type ValidationOptions struct {
	SkipGoModTidy bool // Skip go mod tidy (useful when the caller already ran it)
}

// ValidateCompilation runs sqlc generate (if applicable), go mod tidy, and
// go build ./... in the given app directory. This ensures that generated code
// actually compiles and catches syntax errors, type mismatches, and missing
// dependencies before they reach production.
func ValidateCompilation(t *testing.T, appDir string, opts ...ValidationOptions) {
	t.Helper()

	var opt ValidationOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	env := envWithGOWORKOff()

	// Run sqlc generate if sqlc.yaml exists and queries.sql has actual queries.
	// Freshly created apps have an empty queries.sql (comments only) which
	// causes sqlc to fail with "no queries contained in paths".
	sqlcPath := filepath.Join(appDir, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err == nil {
		if hasQueries(filepath.Join(appDir, "database/queries.sql")) {
			t.Log("Running sqlc generate...")
			sqlcCmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
			sqlcCmd.Dir = appDir
			sqlcCmd.Env = env
			if output, err := sqlcCmd.CombinedOutput(); err != nil {
				t.Fatalf("sqlc generate failed in %s: %v\nOutput: %s", appDir, err, output)
			}
		} else {
			t.Log("Skipping sqlc generate (no queries in queries.sql)")
		}
	}

	// Run go mod tidy to catch dependency problems introduced by generated code.
	// Skip when the caller already ran it (e.g., createTestApp runs tidy after lvt new).
	if opt.SkipGoModTidy {
		t.Log("Skipping go mod tidy (already run by caller)")
	} else {
		t.Log("Running go mod tidy...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = appDir
		tidyCmd.Env = env
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			t.Fatalf("go mod tidy failed in %s: %v\nOutput: %s", appDir, err, output)
		}
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

// hasQueries checks whether a queries.sql file contains at least one actual
// SQL query (a non-empty, non-comment line).
func hasQueries(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "--") {
			return true
		}
	}
	return false
}
