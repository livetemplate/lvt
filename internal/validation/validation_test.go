package validation

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livetemplate/lvt/internal/validator"
)

// ---------------------------------------------------------------------------
// Engine tests
// ---------------------------------------------------------------------------

func TestEngine_SelectiveChecks(t *testing.T) {
	dir := t.TempDir()

	// Engine with only the GoModCheck — no templates or migrations needed.
	writeFile(t, dir, "go.mod", "module example.com/test\n\ngo 1.21\n")
	e := NewEngine(WithCheck(&GoModCheck{}))
	result := e.Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid result, got: %s", result.Format())
	}
}

func TestEngine_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancelled — no sleep needed

	// Multiple checks: cancellation should stop after first check boundary.
	e := NewEngine(WithCheck(&GoModCheck{}), WithCheck(&TemplateCheck{}))
	result := e.Run(ctx, t.TempDir())

	if result.Valid {
		t.Error("expected invalid result due to cancellation")
	}

	var found bool
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "cancelled") {
			found = true
		}
	}
	if !found {
		t.Error("expected a cancellation error issue")
	}

	// Should record exactly one cancellation error, not one per check.
	cancelCount := 0
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelError && strings.Contains(issue.Message, "cancelled") {
			cancelCount++
		}
	}
	if cancelCount != 1 {
		t.Errorf("expected 1 cancellation error, got %d", cancelCount)
	}
}

// stubCheck is a test check that sleeps, used to exercise WithTimeout.
type stubCheck struct {
	delay time.Duration
}

func (s *stubCheck) Name() string { return "stub" }
func (s *stubCheck) Run(ctx context.Context, _ string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	select {
	case <-time.After(s.delay):
		// completed normally
	case <-ctx.Done():
		result.AddError("check cancelled: "+ctx.Err().Error(), "", 0)
	}
	return result
}

func TestEngine_WithTimeout(t *testing.T) {
	// Engine timeout should cancel a slow check.
	e := NewEngine(
		WithTimeout(10*time.Millisecond),
		WithCheck(&stubCheck{delay: 5 * time.Second}),
	)
	result := e.Run(context.Background(), t.TempDir())

	if result.Valid {
		t.Error("expected invalid result due to timeout")
	}
	var found bool
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "cancelled") || strings.Contains(issue.Message, "deadline") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected timeout/cancel error, got: %+v", result.Issues)
	}
}

// ---------------------------------------------------------------------------
// GoModCheck tests
// ---------------------------------------------------------------------------

func TestGoModCheck_Valid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.21\n\nrequire (\n\tgithub.com/foo/bar v1.0.0\n)\n")

	result := (&GoModCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestGoModCheck_ValidNoRequires(t *testing.T) {
	dir := t.TempDir()
	// stdlib-only app — valid with no require directives.
	writeFile(t, dir, "go.mod", "module example.com/test\n\ngo 1.21\n")

	result := (&GoModCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid for stdlib-only app, got: %s", result.Format())
	}
}

func TestGoModCheck_Missing(t *testing.T) {
	dir := t.TempDir()

	result := (&GoModCheck{}).Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for missing go.mod")
	}
	assertHasError(t, result, "not found")
}

func TestGoModCheck_InvalidSyntax(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "this is not valid go.mod syntax {\n")

	result := (&GoModCheck{}).Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for malformed go.mod")
	}
	assertHasError(t, result, "parse error")
}

func TestGoModCheck_MissingModulePath(t *testing.T) {
	dir := t.TempDir()
	// Valid syntax but no module directive.
	writeFile(t, dir, "go.mod", "go 1.21\n")

	result := (&GoModCheck{}).Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for missing module path")
	}
	assertHasError(t, result, "no module path")
}

func TestGoModCheck_MissingGoVersion(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/test\n")

	result := (&GoModCheck{}).Run(context.Background(), dir)

	// Missing go version is a warning, not an error.
	if !result.Valid {
		t.Errorf("expected valid (warning only), got: %s", result.Format())
	}
	assertHasWarning(t, result, "no go version")
}

// ---------------------------------------------------------------------------
// TemplateCheck tests
// ---------------------------------------------------------------------------

func TestTemplateCheck_Valid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "views/index.tmpl", `<h1>{{.Title}}</h1>`)

	result := (&TemplateCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
}

func TestTemplateCheck_Invalid(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "views/bad.tmpl", "<h1>{{.Title}</h1>") // missing closing }}

	result := (&TemplateCheck{}).Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for broken template")
	}
	if result.ErrorCount() == 0 {
		t.Error("expected at least one error")
	}

	// Verify the error has file info.
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelError {
			if issue.File == "" {
				t.Error("expected file path in error issue")
			}
			if issue.Line > 0 && issue.Hint == "" {
				t.Error("expected hint with source context when line is known")
			}
		}
	}
}

func TestTemplateCheck_NoTemplates(t *testing.T) {
	dir := t.TempDir()

	result := (&TemplateCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Error("expected valid when no templates exist")
	}
	// Should have an info issue.
	var found bool
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelInfo && strings.Contains(issue.Message, "no .tmpl") {
			found = true
		}
	}
	if !found {
		t.Error("expected info about no templates found")
	}
}

func TestTemplateCheck_MismatchedDelimiters(t *testing.T) {
	dir := t.TempDir()
	// Template that parses fine but has extra }} in plain text content
	// (e.g. from inline JS or CSS). Parse succeeds, delimiter count differs.
	writeFile(t, dir, "views/warn.tmpl", `<h1>{{.Title}}</h1><script>if(x){y={}}}</script>`)

	result := (&TemplateCheck{}).Run(context.Background(), dir)

	// Parse succeeds, so the structural delimiter check runs and warns.
	assertHasWarning(t, result, "mismatched delimiters")
}

func TestTemplateCheck_SkipsHiddenAndVendorDirs(t *testing.T) {
	dir := t.TempDir()
	// Valid template in the app.
	writeFile(t, dir, "views/index.tmpl", `<h1>{{.Title}}</h1>`)
	// Invalid template in .git — should be skipped.
	writeFile(t, dir, ".git/hooks/bad.tmpl", `<h1>{{.Title}</h1>`)
	// Invalid template in vendor — should be skipped.
	writeFile(t, dir, "vendor/lib/bad.tmpl", `<h1>{{.Title}</h1>`)

	result := (&TemplateCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid (hidden/vendor dirs skipped), got: %s", result.Format())
	}
}

// ---------------------------------------------------------------------------
// MigrationCheck tests
// ---------------------------------------------------------------------------

func TestMigrationCheck_ValidSQL(t *testing.T) {
	dir := t.TempDir()
	migration := `-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS users;
`
	writeFile(t, dir, "database/migrations/001_create_users.sql", migration)

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
}

func TestMigrationCheck_InvalidSQL(t *testing.T) {
	dir := t.TempDir()
	migration := `-- +goose Up
CREATE TABL users (
    id INTEGER PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS users;
`
	writeFile(t, dir, "database/migrations/001_bad.sql", migration)

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	// Invalid SQL is a warning (SQLite dialect may differ from target DB).
	assertHasWarning(t, result, "SQL error")
}

func TestMigrationCheck_NoMigrations(t *testing.T) {
	dir := t.TempDir()

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Error("expected valid when no migrations directory")
	}
	var found bool
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected info about no migrations")
	}
}

func TestMigrationCheck_MultiMigrationOrdering(t *testing.T) {
	dir := t.TempDir()
	m1 := `-- +goose Up
CREATE TABLE authors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authors;
`
	m2 := `-- +goose Up
CREATE TABLE books (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER NOT NULL REFERENCES authors(id)
);

-- +goose Down
DROP TABLE IF EXISTS books;
`
	writeFile(t, dir, "database/migrations/001_authors.sql", m1)
	writeFile(t, dir, "database/migrations/002_books.sql", m2)

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid for ordered migrations, got: %s", result.Format())
	}
}

func TestMigrationCheck_MissingGooseDirective(t *testing.T) {
	dir := t.TempDir()
	// No goose directive — just raw SQL.
	migration := `CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
`
	writeFile(t, dir, "database/migrations/001_no_goose.sql", migration)

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	assertHasWarning(t, result, "missing -- +goose Up")
}

func TestMigrationCheck_StatementBeginEnd(t *testing.T) {
	dir := t.TempDir()
	migration := `-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;
`
	writeFile(t, dir, "database/migrations/001_events.sql", migration)

	result := (&MigrationCheck{}).Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
}

// ---------------------------------------------------------------------------
// CompilationCheck tests
// ---------------------------------------------------------------------------

func TestCompilationCheck_ValidProject(t *testing.T) {
	if testing.Short() {
		t.Skip("compilation check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

	c := &CompilationCheck{}
	result := c.Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
}

func TestCompilationCheck_InvalidCode(t *testing.T) {
	if testing.Short() {
		t.Skip("compilation check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {\n\tundefined()\n}\n")

	c := &CompilationCheck{}
	result := c.Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for code that doesn't compile")
	}

	// Should have a structured error with file info.
	var foundFileError bool
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelError && issue.File != "" && issue.Line > 0 {
			foundFileError = true
		}
	}
	if !foundFileError {
		t.Error("expected at least one error with file:line info")
	}
}

func TestCompilationCheck_OptInTidy(t *testing.T) {
	if testing.Short() {
		t.Skip("compilation check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

	c := &CompilationCheck{RunGoModTidy: true}
	result := c.Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid with tidy opt-in, got: %s", result.Format())
	}
}

// ---------------------------------------------------------------------------
// Integration: Validate() convenience function
// ---------------------------------------------------------------------------

func TestValidate_FullApp(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test runs external commands")
	}

	dir := t.TempDir()

	// Self-contained go.mod (no external dependencies).
	writeFile(t, dir, "go.mod", "module example.com/fullapp\n\ngo 1.21\n")

	// main.go
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {}\n")

	// Template
	writeFile(t, dir, "views/index.tmpl", `<h1>{{.Title}}</h1>`)

	// Migration
	writeFile(t, dir, "database/migrations/001_init.sql", `-- +goose Up
CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS items;
`)

	// Run with all default checks (including compilation).
	result := Validate(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid full app, got: %s", result.Format())
	}
}

// ---------------------------------------------------------------------------
// parseUpStatements unit tests
// ---------------------------------------------------------------------------

func TestParseUpStatements_Basic(t *testing.T) {
	content := `-- +goose Up
CREATE TABLE foo (id INTEGER PRIMARY KEY);

-- +goose Down
DROP TABLE foo;
`
	stmts, hasGooseUp := parseUpStatements(content)
	if !hasGooseUp {
		t.Error("expected hasGooseUp to be true")
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	if !strings.Contains(stmts[0].sql, "CREATE TABLE foo") {
		t.Errorf("unexpected statement: %s", stmts[0].sql)
	}
}

func TestParseUpStatements_StatementBlock(t *testing.T) {
	content := `-- +goose Up
-- +goose StatementBegin
CREATE TABLE bar (
    id INTEGER PRIMARY KEY,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE bar;
`
	stmts, hasGooseUp := parseUpStatements(content)
	if !hasGooseUp {
		t.Error("expected hasGooseUp to be true")
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestParseUpStatements_NoGooseDirective(t *testing.T) {
	content := `CREATE TABLE raw (id INTEGER PRIMARY KEY);
`
	stmts, hasGooseUp := parseUpStatements(content)
	if hasGooseUp {
		t.Error("expected hasGooseUp to be false")
	}
	if len(stmts) != 0 {
		t.Errorf("expected no statements without goose directive, got %d", len(stmts))
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func writeFile(t *testing.T, dir, relPath, content string) {
	t.Helper()
	path := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertHasError(t *testing.T, result *validator.ValidationResult, substr string) {
	t.Helper()
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelError && strings.Contains(issue.Message, substr) {
			return
		}
	}
	t.Errorf("expected error containing %q, got issues: %+v", substr, result.Issues)
}

func assertHasWarning(t *testing.T, result *validator.ValidationResult, substr string) {
	t.Helper()
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelWarning && strings.Contains(issue.Message, substr) {
			return
		}
	}
	t.Errorf("expected warning containing %q, got issues: %+v", substr, result.Issues)
}
