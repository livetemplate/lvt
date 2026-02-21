package validation

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/livetemplate/lvt/internal/validator"
)

// compilerErrorPattern matches Go compiler errors like "file.go:10:5: message".
var compilerErrorPattern = regexp.MustCompile(`^(.+?\.go):(\d+):\d+:\s+(.+)$`)

// CompilationCheck runs sqlc generate, go mod tidy, and go build on an app.
type CompilationCheck struct {
	SkipGoModTidy bool
	SkipSqlc      bool
}

func (c *CompilationCheck) Name() string { return "compilation" }

func (c *CompilationCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	env := envWithGOWORKOff()

	// sqlc generate (if applicable)
	if !c.SkipSqlc {
		c.runSqlc(ctx, appPath, env, result)
	}

	// go mod tidy
	if !c.SkipGoModTidy {
		tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
		tidyCmd.Dir = appPath
		tidyCmd.Env = env
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			result.AddError("go mod tidy failed: "+strings.TrimSpace(string(output)), "go.mod", 0)
			return result
		}
	}

	// go build ./...
	buildCmd := exec.CommandContext(ctx, "go", "build", "./...")
	buildCmd.Dir = appPath
	buildCmd.Env = env
	if output, err := buildCmd.CombinedOutput(); err != nil {
		parseCompilerErrors(string(output), result)
		if result.ErrorCount() == 0 {
			// Couldn't parse structured errors — add the raw output.
			result.AddError("compilation failed: "+strings.TrimSpace(string(output)), "", 0)
		}
	}

	return result
}

func (c *CompilationCheck) runSqlc(ctx context.Context, appPath string, env []string, result *validator.ValidationResult) {
	sqlcPath := filepath.Join(appPath, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcPath); err != nil {
		return // no sqlc config — nothing to do
	}
	if !hasQueries(filepath.Join(appPath, "database/queries.sql")) {
		return // no actual queries
	}

	cmd := exec.CommandContext(ctx, "go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest", "generate", "-f", sqlcPath)
	cmd.Dir = appPath
	cmd.Env = env
	if output, err := cmd.CombinedOutput(); err != nil {
		result.AddError("sqlc generate failed: "+strings.TrimSpace(string(output)), "database/sqlc.yaml", 0)
	}
}

// envWithGOWORKOff returns the current env with GOWORK=off.
func envWithGOWORKOff() []string {
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GOWORK=") {
			env = append(env, e)
		}
	}
	return append(env, "GOWORK=off")
}

// hasQueries returns true if the file contains at least one non-comment line.
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

// parseCompilerErrors extracts file:line errors from go build output.
func parseCompilerErrors(output string, result *validator.ValidationResult) {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if m := compilerErrorPattern.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[2])
			result.AddError(m[3], m[1], lineNum)
		}
	}
}
