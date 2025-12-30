package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateAuth_Handler(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Generate auth files
	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword:  true,
		EnableMagicLink: false,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	// Check auth.go exists in app/auth/
	authPath := filepath.Join(tmpDir, "app", "auth", "auth.go")
	if _, err := os.Stat(authPath); os.IsNotExist(err) {
		t.Errorf("auth.go not generated at %s", authPath)
	}

	// Read and verify content imports from pkg/
	content, err := os.ReadFile(authPath)
	if err != nil {
		t.Fatalf("failed to read auth.go: %v", err)
	}

	contentStr := string(content)

	// Verify imports from lvt/pkg/ packages (not shared/)
	if !strings.Contains(contentStr, "github.com/livetemplate/lvt/pkg/password") {
		t.Error("auth.go should import password from lvt/pkg/password")
	}
	if !strings.Contains(contentStr, "github.com/livetemplate/lvt/pkg/email") {
		t.Error("auth.go should import email from lvt/pkg/email")
	}

	// Verify does NOT import from shared/
	if strings.Contains(contentStr, "/shared/password") {
		t.Error("auth.go should NOT import from shared/password (use lvt/pkg/password instead)")
	}
	if strings.Contains(contentStr, "/shared/email") {
		t.Error("auth.go should NOT import from shared/email (use lvt/pkg/email instead)")
	}
}

func TestGenerateAuth_NoSharedDirectory(t *testing.T) {
	// Verify that shared/ directory is no longer generated
	tmpDir := t.TempDir()

	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword:      true,
		EnableEmailConfirm:  true,
		EnablePasswordReset: true,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	// Verify shared/ directory does NOT exist
	sharedPath := filepath.Join(tmpDir, "shared")
	if _, err := os.Stat(sharedPath); err == nil {
		t.Error("shared/ directory should NOT be generated (utilities are now in lvt/pkg/)")
	}
}

func TestGenerateAuth_Migration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create migrations directory
	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("failed to create migrations directory: %v", err)
	}

	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword: true,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	// Find migration file (should be timestamped)
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	var migrationFile string
	for _, f := range files {
		if strings.Contains(f.Name(), "create_auth_tables") {
			migrationFile = filepath.Join(migrationsDir, f.Name())
			break
		}
	}

	if migrationFile == "" {
		t.Fatal("auth migration file not found")
	}

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		t.Fatalf("failed to read migration file: %v", err)
	}

	contentStr := string(content)

	// Check for users table
	if !strings.Contains(contentStr, "CREATE TABLE IF NOT EXISTS users") {
		t.Error("migration missing users table")
	}

	// Check for users_tokens table
	if !strings.Contains(contentStr, "CREATE TABLE IF NOT EXISTS users_tokens") {
		t.Error("migration missing users_tokens table")
	}

	// Check for goose directives
	if !strings.Contains(contentStr, "-- +goose Up") {
		t.Error("migration missing goose Up directive")
	}

	if !strings.Contains(contentStr, "-- +goose Down") {
		t.Error("migration missing goose Down directive")
	}

	// Check for case-insensitive email (COLLATE NOCASE for SQLite)
	if !strings.Contains(contentStr, "COLLATE NOCASE") {
		t.Error("migration missing case-insensitive email")
	}
}

func TestGenerateAuth_Queries(t *testing.T) {
	tmpDir := t.TempDir()

	// Create database directory
	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword:      true,
		EnableEmailConfirm:  true,
		EnablePasswordReset: true,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	queriesPath := filepath.Join(dbDir, "queries.sql")
	if _, err := os.Stat(queriesPath); os.IsNotExist(err) {
		t.Errorf("queries.sql not generated/updated at %s", queriesPath)
	}

	content, err := os.ReadFile(queriesPath)
	if err != nil {
		t.Fatalf("failed to read queries.sql: %v", err)
	}

	contentStr := string(content)

	requiredQueries := []string{
		"-- name: CreateUser :one",
		"-- name: GetUserByEmail :one",
		"-- name: GetUserByID :one",
		"-- name: CreateUserToken :one",
		"-- name: GetUserToken :one",
		"-- name: DeleteUserToken :exec",
	}

	for _, query := range requiredQueries {
		if !strings.Contains(contentStr, query) {
			t.Errorf("queries.sql missing: %s", query)
		}
	}
}

func TestGenerateAuth_Queries_Append(t *testing.T) {
	tmpDir := t.TempDir()

	// Create database directory with existing queries.sql
	dbDir := filepath.Join(tmpDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	existingContent := "-- Existing queries\n-- name: GetSomething :one\nSELECT * FROM something WHERE id = ?;"
	queriesPath := filepath.Join(dbDir, "queries.sql")
	if err := os.WriteFile(queriesPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to write queries.sql: %v", err)
	}

	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword:   true,
		EnableSessionsUI: true,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	content, err := os.ReadFile(queriesPath)
	if err != nil {
		t.Fatalf("failed to read queries.sql: %v", err)
	}

	contentStr := string(content)

	// Verify existing content is preserved
	if !strings.Contains(contentStr, existingContent) {
		t.Error("queries.sql did not preserve existing content")
	}

	// Verify new content was appended
	if !strings.Contains(contentStr, "-- name: CreateUser :one") {
		t.Error("queries.sql missing new auth queries")
	}

	// Verify separator was added
	if !strings.Contains(contentStr, "\n\n-- Auth Queries") {
		t.Error("queries.sql missing separator between existing and new content")
	}
}

func TestGenerateAuth_UpdateDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal go.mod without dependencies
	goModContent := `module testapp

go 1.21
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Create go.sum to avoid issues
	if err := os.WriteFile(filepath.Join(tmpDir, "go.sum"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create go.sum: %v", err)
	}

	err := GenerateAuth(tmpDir, &AuthConfig{
		EnablePassword: true,
		EnableCSRF:     true,
	})
	if err != nil {
		t.Fatalf("GenerateAuth failed: %v", err)
	}

	// Read updated go.mod
	content, err := os.ReadFile(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}

	contentStr := string(content)

	// Check for bcrypt
	if !strings.Contains(contentStr, "golang.org/x/crypto") {
		t.Error("go.mod missing golang.org/x/crypto dependency")
	}

	// Check for gorilla/csrf
	if !strings.Contains(contentStr, "github.com/gorilla/csrf") {
		t.Error("go.mod missing github.com/gorilla/csrf dependency")
	}
}
