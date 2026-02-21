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

// CompilationCheck runs sqlc generate (optional), go mod tidy (opt-in), and
// go build on an app directory. By default only go build runs; set
// RunGoModTidy or RunSqlc to true to enable those steps.
type CompilationCheck struct {
	// RunGoModTidy runs go mod tidy before building. Off by default because
	// it mutates go.mod/go.sum in the target directory.
	RunGoModTidy bool
	// RunSqlc runs sqlc generate before building (if sqlc.yaml exists and
	// queries.sql has content). Off by default because it requires a locally
	// installed sqlc binary.
	RunSqlc bool
}

func (c *CompilationCheck) Name() string { return "compilation" }

func (c *CompilationCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	env := envWithGOWORKOff()

	// sqlc generate (opt-in)
	if c.RunSqlc {
		c.runSqlc(ctx, appPath, env, result)
		if result.HasErrors() {
			return result // sqlc failure means build results would be unreliable
		}
	}

	// go mod tidy (opt-in — mutates go.mod/go.sum)
	if c.RunGoModTidy {
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
	sqlcCfg := filepath.Join(appPath, "database/sqlc.yaml")
	if _, err := os.Stat(sqlcCfg); err != nil {
		return // no sqlc config — nothing to do
	}
	if !hasQueries(filepath.Join(appPath, "database/queries.sql")) {
		return // no actual queries
	}

	// Require a locally installed sqlc binary; skip if not found.
	sqlcBin, err := exec.LookPath("sqlc")
	if err != nil {
		result.AddInfo("sqlc not found in PATH, skipping sqlc generate", "database/sqlc.yaml", 0)
		return
	}

	cmd := exec.CommandContext(ctx, sqlcBin, "generate", "-f", sqlcCfg)
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
	// If scanner hit an I/O error, conservatively return false.
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
