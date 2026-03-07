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

func TestBase_StyleData(t *testing.T) {
	b := NewBase("myid", "dropdown")

	// Initially nil
	if b.StyleData() != nil {
		t.Error("expected nil StyleData initially")
	}

	// Set and get
	data := "test-style-data"
	b.SetStyleData(data)
	if b.StyleData() != data {
		t.Errorf("StyleData() = %v, want %v", b.StyleData(), data)
	}

	// SetStyled clears cached style data
	b.SetStyled(false)
	if b.StyleData() != nil {
		t.Error("expected nil StyleData after SetStyled(false)")
	}
}

func TestBase_StyleDataWithStruct(t *testing.T) {
	type testStyles struct {
		Panel string
		Body  string
	}

	b := NewBase("myid", "modal")
	s := testStyles{Panel: "bg-white", Body: "p-4"}
	b.SetStyleData(s)

	got, ok := b.StyleData().(testStyles)
	if !ok {
		t.Fatal("StyleData type assertion failed")
	}
	if got.Panel != "bg-white" {
		t.Errorf("Panel = %q, want %q", got.Panel, "bg-white")
	}
	if got.Body != "p-4" {
		t.Errorf("Body = %q, want %q", got.Body, "p-4")
	}
}
