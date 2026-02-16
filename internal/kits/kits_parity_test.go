package kits

import (
	"fmt"
	"strings"
	"testing"
)

// TestKitFeatureParity is a smoke test that ensures the multi and single kit
// monolithic templates all include a shared set of critical UI features. This
// catches regressions where a kit template loses a required feature.
func TestKitFeatureParity(t *testing.T) {
	kits := []string{"multi", "single"}

	// Features that both CRUD kits must have in their monolithic template.
	requiredFeatures := []struct {
		name    string
		pattern string
	}{
		{"delete button in edit modal", `lvt-data-id="{{.EditingID}}"`},
		{"cancel edit button", `lvt-click="cancel_edit"`},
		{"update form submission", `lvt-submit="update"`},
		{"add form submission", `lvt-submit="add"`},
		{"add modal open button", `lvt-modal-open="add-modal"`},
		{"edit button in table", `lvt-click="edit"`},
	}

	for _, kit := range kits {
		t.Run(kit, func(t *testing.T) {
			path := fmt.Sprintf("system/%s/templates/resource/template.tmpl.tmpl", kit)
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
