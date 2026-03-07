package telemetry

import (
	"testing"
)

func TestAttributeErrors_FilePathMatching(t *testing.T) {
	errors := []GenerationError{
		{Phase: "generation", File: "components/modal/modal.go", Message: "undefined: modal.Size"},
		{Phase: "template", File: "app/posts/posts.tmpl", Message: "unrelated error"},
		{Phase: "generation", File: "components/toast/toast.go", Message: "missing field"},
	}
	componentsUsed := []string{"modal", "toast"}

	result := AttributeErrors(errors, componentsUsed)

	if len(result) != 2 {
		t.Fatalf("expected 2 component errors, got %d", len(result))
	}

	if result[0].Component != "modal" {
		t.Errorf("expected first component 'modal', got %q", result[0].Component)
	}
	if result[0].File != "components/modal/modal.go" {
		t.Errorf("expected file 'components/modal/modal.go', got %q", result[0].File)
	}

	if result[1].Component != "toast" {
		t.Errorf("expected second component 'toast', got %q", result[1].Component)
	}
}

func TestAttributeErrors_MessageMatching(t *testing.T) {
	errors := []GenerationError{
		{Phase: "runtime", Message: "modal.New: invalid size parameter"},
		{Phase: "runtime", Message: "toast: container not found"},
	}
	componentsUsed := []string{"modal", "toast", "dropdown"}

	result := AttributeErrors(errors, componentsUsed)

	if len(result) != 2 {
		t.Fatalf("expected 2 component errors, got %d", len(result))
	}
	if result[0].Component != "modal" {
		t.Errorf("expected 'modal', got %q", result[0].Component)
	}
	if result[1].Component != "toast" {
		t.Errorf("expected 'toast', got %q", result[1].Component)
	}
}

func TestAttributeErrors_NoMatch(t *testing.T) {
	errors := []GenerationError{
		{Phase: "generation", File: "app/posts/handler.go", Message: "syntax error"},
	}
	componentsUsed := []string{"modal", "toast"}

	result := AttributeErrors(errors, componentsUsed)

	if len(result) != 0 {
		t.Fatalf("expected 0 component errors, got %d", len(result))
	}
}

func TestAttributeErrors_EmptyInputs(t *testing.T) {
	// No errors
	result := AttributeErrors(nil, []string{"modal"})
	if result != nil {
		t.Error("expected nil for nil errors")
	}

	// No components
	result = AttributeErrors([]GenerationError{{Message: "test"}}, nil)
	if result != nil {
		t.Error("expected nil for nil components")
	}

	// Both empty
	result = AttributeErrors(nil, nil)
	if result != nil {
		t.Error("expected nil for both nil")
	}
}

func TestComponentsFromUsage(t *testing.T) {
	type ComponentUsage struct {
		UseModal    bool
		UseToast    bool
		UseDropdown bool
	}

	tests := []struct {
		name     string
		usage    any
		expected []string
	}{
		{
			name:     "all true",
			usage:    ComponentUsage{UseModal: true, UseToast: true, UseDropdown: true},
			expected: []string{"modal", "toast", "dropdown"},
		},
		{
			name:     "partial",
			usage:    ComponentUsage{UseModal: true, UseToast: false, UseDropdown: true},
			expected: []string{"modal", "dropdown"},
		},
		{
			name:     "none",
			usage:    ComponentUsage{},
			expected: nil,
		},
		{
			name:     "nil",
			usage:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComponentsFromUsage(tt.usage)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d components, got %d: %v", len(tt.expected), len(result), result)
			}
			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("component[%d]: expected %q, got %q", i, exp, result[i])
				}
			}
		})
	}
}

func TestComponentsFromUsage_Pointer(t *testing.T) {
	type Usage struct {
		UseModal bool
		UseToast bool
	}
	u := &Usage{UseModal: true}
	result := ComponentsFromUsage(u)
	if len(result) != 1 || result[0] != "modal" {
		t.Errorf("expected [modal], got %v", result)
	}
}

func TestAttributeErrors_CaseInsensitive(t *testing.T) {
	errors := []GenerationError{
		{Phase: "generation", File: "Components/Modal/Modal.go", Message: "error"},
	}
	result := AttributeErrors(errors, []string{"modal"})

	if len(result) != 1 {
		t.Fatalf("expected 1 component error (case-insensitive), got %d", len(result))
	}
	if result[0].Component != "modal" {
		t.Errorf("expected 'modal', got %q", result[0].Component)
	}
}

func TestAttributeErrors_PrefixNoFalsePositive(t *testing.T) {
	// Verify that "toast." does NOT match "toaster." — the dot acts as a
	// natural word boundary so component-name prefixes don't false-positive.
	errors := []GenerationError{
		{Phase: "runtime", Message: "toaster.New: invalid config"},
	}
	result := AttributeErrors(errors, []string{"toast"})

	if len(result) != 0 {
		t.Fatalf("expected 0 (no false positive), got %d", len(result))
	}
}

func TestAttributeErrors_SuffixNoFalsePositive(t *testing.T) {
	// "premodal.init: failed" must NOT match component "modal" — the
	// left-boundary guard requires a non-word character before the match.
	errors := []GenerationError{
		{Phase: "runtime", Message: "premodal.init: failed"},
	}
	result := AttributeErrors(errors, []string{"modal"})

	if len(result) != 0 {
		t.Fatalf("expected 0 (left-boundary guard should prevent match), got %d", len(result))
	}
}

func TestAttributeErrors_BoundaryMatch(t *testing.T) {
	// "modal.init: failed" should still match at a word boundary
	errors := []GenerationError{
		{Phase: "runtime", Message: "modal.init: failed"},
	}
	result := AttributeErrors(errors, []string{"modal"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestComponentsFromUsage_IgnoresNonUseFields(t *testing.T) {
	// Struct with bool fields that don't start with "Use" must be skipped.
	type MixedUsage struct {
		UseModal bool
		Enabled  bool
		Active   bool
		UseToast bool
	}
	result := ComponentsFromUsage(MixedUsage{
		UseModal: true,
		Enabled:  true,
		Active:   true,
		UseToast: true,
	})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(result), result)
	}
	if result[0] != "modal" || result[1] != "toast" {
		t.Errorf("expected [modal, toast], got %v", result)
	}
}

func TestAttributeErrors_PreservesPhaseAndMessage(t *testing.T) {
	errors := []GenerationError{
		{Phase: "compilation", File: "components/dropdown/dropdown.go", Message: "undefined: Searchable", Context: "line 42"},
	}
	result := AttributeErrors(errors, []string{"dropdown"})

	if len(result) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result))
	}
	if result[0].Phase != "compilation" {
		t.Errorf("expected phase 'compilation', got %q", result[0].Phase)
	}
	if result[0].Message != "undefined: Searchable" {
		t.Errorf("expected original message, got %q", result[0].Message)
	}
}
