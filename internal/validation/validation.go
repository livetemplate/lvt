// Package validation provides a unified engine for validating generated
// LiveTemplate applications. It aggregates multiple checks (go.mod structure,
// template syntax, SQL migrations, compilation) into a single pass and
// reports issues using the shared validator.ValidationResult type.
package validation

import (
	"context"
	"time"

	"github.com/livetemplate/lvt/internal/validator"
)

// Check is a single validation check that can be run against an app directory.
type Check interface {
	// Name returns a human-readable name for the check (e.g. "go.mod").
	Name() string
	// Run executes the check and returns issues found.
	Run(ctx context.Context, appPath string) *validator.ValidationResult
}

// Engine orchestrates multiple validation checks.
type Engine struct {
	checks  []Check
	timeout time.Duration
}

// Option configures an Engine.
type Option func(*Engine)

// WithCheck appends a check to the engine.
func WithCheck(c Check) Option {
	return func(e *Engine) {
		e.checks = append(e.checks, c)
	}
}

// WithTimeout sets the maximum duration for all checks combined.
func WithTimeout(d time.Duration) Option {
	return func(e *Engine) {
		e.timeout = d
	}
}

// NewEngine creates an engine with the given options.
func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		timeout: 5 * time.Minute,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// DefaultEngine returns an engine pre-loaded with all built-in checks.
func DefaultEngine() *Engine {
	return NewEngine(
		WithCheck(&GoModCheck{}),
		WithCheck(&TemplateCheck{}),
		WithCheck(&MigrationCheck{}),
		WithCheck(&CompilationCheck{}),
	)
}

// Run executes every registered check sequentially and merges their results.
// Checks run independently even if earlier ones report errors, so the caller
// gets a complete picture of all issues. It respects context cancellation
// between checks so a timeout or cancellation stops early.
func (e *Engine) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	if e.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	result := validator.NewValidationResult()
	for _, c := range e.checks {
		if err := ctx.Err(); err != nil {
			result.AddError("validation cancelled: "+err.Error(), "", 0)
			break
		}
		result.Merge(c.Run(ctx, appPath))
	}
	return result
}

// PostGenEngine returns an engine suited for post-generation validation.
// It checks structural aspects (go.mod, templates, migrations) but skips
// compilation because the app may not compile until sqlc generate is run.
func PostGenEngine() *Engine {
	return NewEngine(
		WithCheck(&GoModCheck{}),
		WithCheck(&TemplateCheck{}),
		WithCheck(&MigrationCheck{}),
	)
}

// ValidatePostGen runs structural checks (no compilation) after code generation.
// Use this right after lvt gen commands where the app may not yet compile.
func ValidatePostGen(ctx context.Context, appPath string) *validator.ValidationResult {
	return PostGenEngine().Run(ctx, appPath)
}

// FullEngine returns an engine with all default checks plus RuntimeCheck.
// RuntimeCheck is expensive (builds binary, starts subprocess, probes HTTP)
// so it is not included in DefaultEngine.
func FullEngine() *Engine {
	e := DefaultEngine()
	WithCheck(&RuntimeCheck{})(e)
	return e
}

// Validate runs all default checks including compilation (go build ./...).
// This may be slow if dependencies are not cached. Use NewEngine with
// selective checks for latency-sensitive call sites (e.g. file watchers).
func Validate(ctx context.Context, appPath string) *validator.ValidationResult {
	return DefaultEngine().Run(ctx, appPath)
}

// ValidateFull runs all checks including the runtime startup check.
// This is the most thorough validation and is the most expensive.
func ValidateFull(ctx context.Context, appPath string) *validator.ValidationResult {
	return FullEngine().Run(ctx, appPath)
}
