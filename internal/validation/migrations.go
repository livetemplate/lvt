package validation

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livetemplate/lvt/internal/validator"
	_ "modernc.org/sqlite"
)

// MigrationCheck validates SQL migration files against an in-memory SQLite DB.
// Because apps may target PostgreSQL, SQL execution errors are reported as
// warnings (they may be false positives for non-SQLite dialects). Structural
// issues (missing files, missing goose directives) are always reported.
type MigrationCheck struct{}

func (c *MigrationCheck) Name() string { return "migrations (sqlite)" }

func (c *MigrationCheck) Run(ctx context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	migrationsDir := filepath.Join(appPath, "database", "migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			result.AddInfo("no database/migrations directory found", "", 0)
			return result
		}
		result.AddError("failed to read migrations directory: "+err.Error(), "database/migrations", 0)
		return result
	}

	// Collect .sql files sorted by name (goose ordering).
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		result.AddInfo("no .sql migration files found", "database/migrations", 0)
		return result
	}

	// Open an in-memory SQLite database.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		result.AddError("failed to open in-memory SQLite: "+err.Error(), "", 0)
		return result
	}
	defer db.Close()

	for _, name := range files {
		if ctx.Err() != nil {
			result.AddError("validation cancelled: "+ctx.Err().Error(), "", 0)
			break
		}
		c.validateMigration(ctx, db, migrationsDir, name, result)
	}

	return result
}

func (c *MigrationCheck) validateMigration(ctx context.Context, db *sql.DB, dir, name string, result *validator.ValidationResult) {
	relPath := filepath.Join("database", "migrations", name)
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		result.AddError("failed to read migration: "+err.Error(), relPath, 0)
		return
	}

	content := string(data)
	stmts, hasGooseUp := parseUpStatements(content)

	if !hasGooseUp {
		result.AddWarning("missing -- +goose Up directive", relPath, 0)
	}

	if len(stmts) == 0 {
		return
	}

	for _, stmt := range stmts {
		s := strings.TrimSpace(stmt.sql)
		if s == "" {
			continue
		}
		_, execErr := db.ExecContext(ctx, s)
		if execErr != nil {
			// Warn rather than error: the target DB may be PostgreSQL,
			// so SQLite-specific failures could be false positives.
			result.AddWarning(
				fmt.Sprintf("SQL error (sqlite): %s", execErr.Error()),
				relPath, stmt.line,
			)
		}
	}
}

// statement tracks a SQL statement and the line it starts on.
type statement struct {
	sql  string
	line int
}

// parseUpStatements extracts SQL statements from the -- +goose Up section,
// handling StatementBegin/End blocks and semicolon-delimited statements.
func parseUpStatements(content string) ([]statement, bool) {
	lines := strings.Split(content, "\n")

	var (
		inUp             bool
		hasGooseUp       bool
		inStatementBlock bool
		stmts            []statement
		current          strings.Builder
		currentLine      int
	)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Goose directives.
		if strings.HasPrefix(trimmed, "-- +goose Up") {
			inUp = true
			hasGooseUp = true
			continue
		}
		if strings.HasPrefix(trimmed, "-- +goose Down") {
			// Stop at the Down section.
			if inUp {
				// Flush any accumulated statement.
				if current.Len() > 0 {
					stmts = append(stmts, statement{sql: current.String(), line: currentLine})
				}
			}
			break
		}
		if strings.HasPrefix(trimmed, "-- +goose StatementBegin") {
			inStatementBlock = true
			continue
		}
		if strings.HasPrefix(trimmed, "-- +goose StatementEnd") {
			if current.Len() > 0 {
				stmts = append(stmts, statement{sql: current.String(), line: currentLine})
				current.Reset()
			}
			inStatementBlock = false
			continue
		}

		if !inUp {
			continue
		}

		// Skip comments.
		if strings.HasPrefix(trimmed, "--") {
			continue
		}

		if inStatementBlock {
			if current.Len() == 0 && trimmed != "" {
				currentLine = i + 1
			}
			current.WriteString(line)
			current.WriteString("\n")
			continue
		}

		// Semicolon-delimited mode.
		if trimmed == "" {
			continue
		}
		if current.Len() == 0 {
			currentLine = i + 1
		}
		current.WriteString(line)
		current.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			stmts = append(stmts, statement{sql: current.String(), line: currentLine})
			current.Reset()
		}
	}

	// Flush remaining (no Down section encountered).
	if inUp && current.Len() > 0 {
		stmts = append(stmts, statement{sql: current.String(), line: currentLine})
	}

	return stmts, hasGooseUp
}
