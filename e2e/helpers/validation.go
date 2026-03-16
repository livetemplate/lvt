package helpers

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// sqlcPackage is the pinned sqlc version for reproducible CI builds.
const sqlcPackage = "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0"

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
			sqlcCmd := exec.Command("go", "run", sqlcPackage, "generate", "-f", sqlcPath)
			sqlcCmd.Dir = appDir
			sqlcCmd.Env = env
			if output, err := sqlcCmd.CombinedOutput(); err != nil {
				t.Fatalf("sqlc generate failed in %s: %v\nOutput: %s", appDir, err, output)
			}
		} else {
			t.Log("Skipping sqlc generate (no queries in queries.sql)")
		}
	}

	// Inject components module if not already present (components sub-module isn't published yet)
	InjectComponentsForTest(t, appDir)

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

// InjectComponentsForTest copies the components module into the test app directory
// and adds require/replace directives so go mod tidy can resolve component imports.
// This is needed because the components sub-module isn't published to the Go module proxy yet.
func InjectComponentsForTest(t *testing.T, appDir string) {
	t.Helper()

	// Find project root (e2e/helpers/ is two levels deep from root)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Log("⚠️  Could not determine project root for components injection")
		return
	}
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	componentsSrc := filepath.Join(projectRoot, "components")

	// Check if components source exists
	if _, err := os.Stat(componentsSrc); os.IsNotExist(err) {
		t.Logf("⏭️  Components directory not found at %s, skipping injection", componentsSrc)
		return
	}

	// Skip if already injected
	componentsDst := filepath.Join(appDir, "components")
	if _, err := os.Stat(componentsDst); err == nil {
		return // Already injected
	}

	// Copy components directory into app
	cpCmd := exec.Command("cp", "-r", componentsSrc, componentsDst)
	if output, err := cpCmd.CombinedOutput(); err != nil {
		t.Logf("⚠️  Failed to copy components: %v\nOutput: %s", err, output)
		return
	}

	// Add require and replace directives to app's go.mod
	goModPath := filepath.Join(appDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Logf("⚠️  Failed to read go.mod for components injection: %v", err)
		return
	}

	goModStr := string(goModContent)
	modified := false

	if !strings.Contains(goModStr, "github.com/livetemplate/lvt/components") {
		goModStr += "\nrequire github.com/livetemplate/lvt/components v0.0.0\n"
		goModStr += "\nreplace github.com/livetemplate/lvt/components => ./components\n"
		modified = true
	} else if !strings.Contains(goModStr, "replace github.com/livetemplate/lvt/components") {
		goModStr += "\nreplace github.com/livetemplate/lvt/components => ./components\n"
		modified = true
	}

	if modified {
		if err := os.WriteFile(goModPath, []byte(goModStr), 0644); err != nil {
			t.Logf("⚠️  Failed to update go.mod for components injection: %v", err)
			return
		}
	}
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
