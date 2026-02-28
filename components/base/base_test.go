package base

import (
	"testing"
)

func TestNewBase(t *testing.T) {
	b := NewBase("myid", "dropdown")

	if b.ID() != "myid" {
		t.Errorf("expected ID 'myid', got '%s'", b.ID())
	}

	if b.Namespace() != "dropdown" {
		t.Errorf("expected Namespace 'dropdown', got '%s'", b.Namespace())
	}

	// Default should be styled
	if !b.IsStyled() {
		t.Error("expected default Styled to be true")
	}
}

func TestBase_ActionName(t *testing.T) {
	b := NewBase("user-dropdown", "dropdown")

	tests := []struct {
		action   string
		expected string
	}{
		{"toggle", "toggle_user-dropdown"},
		{"select", "select_user-dropdown"},
		{"close", "close_user-dropdown"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			result := b.ActionName(tt.action)
			if result != tt.expected {
				t.Errorf("ActionName(%q) = %q, want %q", tt.action, result, tt.expected)
			}
		})
	}
}

func TestBase_SetStyled(t *testing.T) {
	b := NewBase("myid", "dropdown")

	// Default is styled
	if !b.IsStyled() {
		t.Error("expected default to be styled")
	}

	// Set to unstyled
	b.SetStyled(false)
	if b.IsStyled() {
		t.Error("expected unstyled after SetStyled(false)")
	}

	// Set back to styled
	b.SetStyled(true)
	if !b.IsStyled() {
		t.Error("expected styled after SetStyled(true)")
	}
}
