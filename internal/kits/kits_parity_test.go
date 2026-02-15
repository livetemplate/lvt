package kits

import (
	"strings"
	"testing"
)

// TestKitFeatureParity verifies that multi and single kit monolithic templates
// contain the same critical UI features. This catches regressions where a kit
// template drifts from the expected feature set.
func TestKitFeatureParity(t *testing.T) {
	templatePath := "system/%s/templates/resource/template.tmpl.tmpl"

	kits := []string{"multi", "single"}

	// Features that both CRUD kits must have in their monolithic template.
	requiredFeatures := []struct {
		name    string
		pattern string
	}{
		{"delete button in edit modal", `lvt-click="delete"`},
		{"cancel edit button", `lvt-click="cancel_edit"`},
		{"update form submission", `lvt-submit="update"`},
		{"add form submission", `lvt-submit="add"`},
		{"add modal open button", `lvt-modal-open="add-modal"`},
		{"edit button in table", `lvt-click="edit"`},
	}

	for _, kit := range kits {
		t.Run(kit, func(t *testing.T) {
			path := strings.Replace(templatePath, "%s", kit, 1)
			data, err := systemKits.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read template for kit %q: %v", kit, err)
			}
			content := string(data)

			for _, feat := range requiredFeatures {
				if !strings.Contains(content, feat.pattern) {
					t.Errorf("Kit %q missing feature %q (pattern: %s)", kit, feat.name, feat.pattern)
				}
			}
		})
	}
}
