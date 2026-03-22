package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/kits"
)

// AuthzConfig holds configuration for authorization generation.
type AuthzConfig struct {
	ModuleName string // Module path from go.mod
	TableName  string // Users table name (e.g., "users")
}

// GenerateAuthz generates the authorization system: role migration,
// role queries, and patches the schema for sqlc.
func GenerateAuthz(projectRoot string, cfg *AuthzConfig) error {
	if cfg.TableName == "" {
		cfg.TableName = "users"
	}

	authDir := filepath.Join(projectRoot, "app", "auth")
	if _, err := os.Stat(authDir); os.IsNotExist(err) {
		return fmt.Errorf("auth system not found at %s — run 'lvt gen auth' first", authDir)
	}

	projectConfig, err := config.LoadProjectConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	kitName := projectConfig.GetKit()
	kitLoader := kits.DefaultLoader()

	funcMap := template.FuncMap{
		"singular": singularize,
	}

	// Generate migration
	migrationsDir := filepath.Join(projectRoot, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now()
	var migrationPath string
	const maxRetries = 3600
	for i := 0; i < maxRetries; i++ {
		timestampStr := timestamp.Format("20060102150405")
		migrationPath = filepath.Join(migrationsDir, fmt.Sprintf("%s_add_user_roles.sql", timestampStr))
		matches, err := filepath.Glob(filepath.Join(migrationsDir, timestampStr+"_*.sql"))
		if err != nil {
			return fmt.Errorf("failed to check for existing migrations: %w", err)
		}
		if len(matches) == 0 {
			break
		}
		timestamp = timestamp.Add(1 * time.Second)
	}

	templateContent, err := kitLoader.LoadKitTemplate(kitName, "authz/migration.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load authz migration template: %w", err)
	}

	tmpl, err := template.New("migration").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse authz migration template: %w", err)
	}

	file, err := os.Create(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, cfg); err != nil {
		return fmt.Errorf("failed to execute migration template: %w", err)
	}

	// Append role queries to queries.sql
	queriesPath := filepath.Join(projectRoot, "database", "queries.sql")
	templateContent, err = kitLoader.LoadKitTemplate(kitName, "authz/queries.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load authz queries template: %w", err)
	}

	tmpl, err = template.New("queries").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse authz queries template: %w", err)
	}

	qFile, err := os.OpenFile(queriesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open queries.sql: %w", err)
	}

	stat, _ := qFile.Stat()
	if stat != nil && stat.Size() > 0 {
		qFile.WriteString("\n\n")
	}

	if err := tmpl.Execute(qFile, cfg); err != nil {
		qFile.Close()
		return fmt.Errorf("failed to execute authz queries template: %w", err)
	}
	qFile.Close()

	// Patch schema.sql to add role column to the users CREATE TABLE
	if err := patchSchemaWithRole(projectRoot, cfg.TableName); err != nil {
		return fmt.Errorf("failed to patch schema.sql: %w", err)
	}

	return nil
}

// patchSchemaWithRole adds the role column to the existing users
// CREATE TABLE in schema.sql so that sqlc generates the correct struct.
func patchSchemaWithRole(projectRoot, tableName string) error {
	schemaPath := filepath.Join(projectRoot, "database", "schema.sql")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	content := string(data)

	// Check if role column already exists
	if strings.Contains(content, "role TEXT") {
		return nil // Already patched
	}

	// Find the users table and add role before created_at
	// Pattern: "    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"
	// in the users table (identified by the table name)
	target := "    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n    updated_at"
	replacement := "    role TEXT NOT NULL DEFAULT 'user',\n    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n    updated_at"

	if !strings.Contains(content, target) {
		// Try alternative pattern without updated_at on next line
		target = "    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP\n)"
		replacement = "    role TEXT NOT NULL DEFAULT 'user',\n    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP\n)"
	}

	if !strings.Contains(content, target) {
		return fmt.Errorf("could not find insertion point for role column in schema.sql — add manually: role TEXT NOT NULL DEFAULT 'user'")
	}

	content = strings.Replace(content, target, replacement, 1)

	return os.WriteFile(schemaPath, []byte(content), 0644)
}
