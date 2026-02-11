//go:build http

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/e2e/helpers"
)

// TestTypeInference tests field type inference
func TestTypeInference(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "testapp")

	// Create app
	if err := runLvtCommand(t, tmpDir, "new", "testapp"); err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Generate resource with inferred types (no :type specified)
	// Note: avoid "created_at" as a user field since the schema template
	// auto-generates it; use "published_at" to test _at suffix inference.
	if err := runLvtCommand(t, appDir, "gen", "resource", "users", "name", "email", "age", "price", "published", "published_at"); err != nil {
		t.Fatalf("Failed to generate resource with type inference: %v", err)
	}

	// Verify schema has correct inferred types
	schemaFile := filepath.Join(appDir, "database", "schema.sql")
	content, err := os.ReadFile(schemaFile)
	if err != nil {
		t.Fatalf("Failed to read schema: %v", err)
	}

	contentStr := string(content)

	// Check inferred types
	checks := map[string]string{
		"name":         "TEXT",     // string
		"email":        "TEXT",     // string
		"age":          "INTEGER",  // int
		"price":        "REAL",     // float
		"published":    "INTEGER",  // bool
		"published_at": "DATETIME", // time
	}

	for field, expectedType := range checks {
		if !strings.Contains(contentStr, field) || !strings.Contains(contentStr, expectedType) {
			t.Errorf("❌ Field '%s' not inferred as %s", field, expectedType)
		}
	}

	// Validate generated code compiles
	helpers.ValidateCompilation(t, appDir)

	t.Log("✅ Type inference working correctly")
}
