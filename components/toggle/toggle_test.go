package toggle

import (
	"testing"
)

// =============================================================================
// Toggle Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates toggle with defaults", func(t *testing.T) {
		tg := New("dark-mode")
		if tg.ID() != "dark-mode" {
			t.Errorf("expected ID 'dark-mode', got %q", tg.ID())
		}
		if tg.Namespace() != "toggle" {
			t.Errorf("expected namespace 'toggle', got %q", tg.Namespace())
		}
		if tg.Checked {
			t.Error("expected unchecked by default")
		}
		if tg.Disabled {
			t.Error("expected not disabled by default")
		}
		if tg.LabelPosition != "right" {
			t.Errorf("expected label position 'right', got %q", tg.LabelPosition)
		}
		if tg.Size != SizeMd {
			t.Errorf("expected size md, got %v", tg.Size)
		}
		if tg.Value != "on" {
			t.Errorf("expected value 'on', got %q", tg.Value)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		tg := New("notifications",
			WithChecked(true),
			WithLabel("Email Notifications"),
			WithSize(SizeLg),
			WithDisabled(true),
		)
		if !tg.Checked {
			t.Error("expected checked")
		}
		if tg.Label != "Email Notifications" {
			t.Errorf("expected label 'Email Notifications', got %q", tg.Label)
		}
		if tg.Size != SizeLg {
			t.Errorf("expected size lg, got %v", tg.Size)
		}
		if !tg.Disabled {
			t.Error("expected disabled")
		}
	})
}

// =============================================================================
// Toggle Option Tests
// =============================================================================

func TestWithChecked(t *testing.T) {
	tg := New("test", WithChecked(true))
	if !tg.Checked {
		t.Error("expected checked")
	}
}

func TestWithDisabled(t *testing.T) {
	tg := New("test", WithDisabled(true))
	if !tg.Disabled {
		t.Error("expected disabled")
	}
}

func TestWithLabel(t *testing.T) {
	tg := New("test", WithLabel("My Label"))
	if tg.Label != "My Label" {
		t.Errorf("expected label 'My Label', got %q", tg.Label)
	}
}

func TestWithLabelPosition(t *testing.T) {
	tg := New("test", WithLabelPosition("left"))
	if tg.LabelPosition != "left" {
		t.Errorf("expected label position 'left', got %q", tg.LabelPosition)
	}
}

func TestWithSize(t *testing.T) {
	sizes := []Size{SizeSm, SizeMd, SizeLg}
	for _, size := range sizes {
		t.Run(string(size), func(t *testing.T) {
			tg := New("test", WithSize(size))
			if tg.Size != size {
				t.Errorf("expected size %v, got %v", size, tg.Size)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	tg := New("test", WithName("dark_mode"))
	if tg.Name != "dark_mode" {
		t.Errorf("expected name 'dark_mode', got %q", tg.Name)
	}
}

func TestWithValue(t *testing.T) {
	tg := New("test", WithValue("enabled"))
	if tg.Value != "enabled" {
		t.Errorf("expected value 'enabled', got %q", tg.Value)
	}
}

func TestWithRequired(t *testing.T) {
	tg := New("test", WithRequired(true))
	if !tg.Required {
		t.Error("expected required")
	}
}

func TestWithDescription(t *testing.T) {
	tg := New("test", WithDescription("Helper text"))
	if tg.Description != "Helper text" {
		t.Errorf("expected description 'Helper text', got %q", tg.Description)
	}
}

func TestWithStyled(t *testing.T) {
	tg := New("test", WithStyled(true))
	if !tg.IsStyled() {
		t.Error("expected styled true")
	}
}

// =============================================================================
// Toggle Method Tests
// =============================================================================

func TestToggle(t *testing.T) {
	tg := New("test")
	tg.Toggle()
	if !tg.Checked {
		t.Error("expected checked after toggle")
	}
	tg.Toggle()
	if tg.Checked {
		t.Error("expected unchecked after second toggle")
	}
}

func TestToggleWhenDisabled(t *testing.T) {
	tg := New("test", WithDisabled(true))
	tg.Toggle()
	if tg.Checked {
		t.Error("should not toggle when disabled")
	}
}

func TestCheck(t *testing.T) {
	tg := New("test")
	tg.Check()
	if !tg.Checked {
		t.Error("expected checked after Check")
	}
	tg.Check()
	if !tg.Checked {
		t.Error("expected still checked after second Check")
	}
}

func TestUncheck(t *testing.T) {
	tg := New("test", WithChecked(true))
	tg.Uncheck()
	if tg.Checked {
		t.Error("expected unchecked after Uncheck")
	}
}

func TestSetChecked(t *testing.T) {
	tg := New("test")
	tg.SetChecked(true)
	if !tg.Checked {
		t.Error("expected checked")
	}
	tg.SetChecked(false)
	if tg.Checked {
		t.Error("expected unchecked")
	}
}

func TestIsOn(t *testing.T) {
	tg := New("test", WithChecked(true))
	if !tg.IsOn() {
		t.Error("expected IsOn true")
	}
	tg.Checked = false
	if tg.IsOn() {
		t.Error("expected IsOn false")
	}
}

func TestIsOff(t *testing.T) {
	tg := New("test")
	if !tg.IsOff() {
		t.Error("expected IsOff true")
	}
	tg.Checked = true
	if tg.IsOff() {
		t.Error("expected IsOff false")
	}
}

func TestHasLabel(t *testing.T) {
	tg := New("test")
	if tg.HasLabel() {
		t.Error("expected HasLabel false")
	}
	tg.Label = "Test"
	if !tg.HasLabel() {
		t.Error("expected HasLabel true")
	}
}

func TestHasDescription(t *testing.T) {
	tg := New("test")
	if tg.HasDescription() {
		t.Error("expected HasDescription false")
	}
	tg.Description = "Help"
	if !tg.HasDescription() {
		t.Error("expected HasDescription true")
	}
}

func TestIsLabelLeft(t *testing.T) {
	tg := New("test", WithLabelPosition("left"))
	if !tg.IsLabelLeft() {
		t.Error("expected IsLabelLeft true")
	}
}

func TestIsLabelRight(t *testing.T) {
	tg := New("test", WithLabelPosition("right"))
	if !tg.IsLabelRight() {
		t.Error("expected IsLabelRight true")
	}
}

// =============================================================================
// Toggle CSS Class Tests
// =============================================================================

func TestSizeClasses(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "w-8 h-4"},
		{SizeMd, "w-11 h-6"},
		{SizeLg, "w-14 h-8"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			tg := New("test", WithSize(tt.size))
			if tg.SizeClasses() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tg.SizeClasses())
			}
		})
	}
}

func TestKnobSizeClasses(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "w-3 h-3"},
		{SizeMd, "w-5 h-5"},
		{SizeLg, "w-6 h-6"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			tg := New("test", WithSize(tt.size))
			if tg.KnobSizeClasses() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tg.KnobSizeClasses())
			}
		})
	}
}

func TestTrackColorClass(t *testing.T) {
	t.Run("unchecked enabled", func(t *testing.T) {
		tg := New("test")
		if tg.TrackColorClass() != "bg-gray-200" {
			t.Errorf("expected 'bg-gray-200', got %q", tg.TrackColorClass())
		}
	})

	t.Run("checked enabled", func(t *testing.T) {
		tg := New("test", WithChecked(true))
		if tg.TrackColorClass() != "bg-blue-600" {
			t.Errorf("expected 'bg-blue-600', got %q", tg.TrackColorClass())
		}
	})

	t.Run("checked disabled", func(t *testing.T) {
		tg := New("test", WithChecked(true), WithDisabled(true))
		if tg.TrackColorClass() != "bg-blue-300" {
			t.Errorf("expected 'bg-blue-300', got %q", tg.TrackColorClass())
		}
	})
}

// =============================================================================
// Checkbox Constructor Tests
// =============================================================================

func TestNewCheckbox(t *testing.T) {
	t.Run("creates checkbox with defaults", func(t *testing.T) {
		c := NewCheckbox("terms")
		if c.ID() != "terms" {
			t.Errorf("expected ID 'terms', got %q", c.ID())
		}
		if c.Checked {
			t.Error("expected unchecked by default")
		}
		if c.Indeterminate {
			t.Error("expected not indeterminate by default")
		}
		if c.Value != "on" {
			t.Errorf("expected value 'on', got %q", c.Value)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		c := NewCheckbox("agree",
			WithCheckboxChecked(true),
			WithCheckboxLabel("I agree"),
			WithCheckboxRequired(true),
		)
		if !c.Checked {
			t.Error("expected checked")
		}
		if c.Label != "I agree" {
			t.Errorf("expected label 'I agree', got %q", c.Label)
		}
		if !c.Required {
			t.Error("expected required")
		}
	})
}

// =============================================================================
// Checkbox Option Tests
// =============================================================================

func TestWithCheckboxIndeterminate(t *testing.T) {
	c := NewCheckbox("test", WithCheckboxIndeterminate(true))
	if !c.Indeterminate {
		t.Error("expected indeterminate")
	}
}

func TestWithCheckboxDisabled(t *testing.T) {
	c := NewCheckbox("test", WithCheckboxDisabled(true))
	if !c.Disabled {
		t.Error("expected disabled")
	}
}

func TestWithCheckboxStyled(t *testing.T) {
	c := NewCheckbox("test", WithCheckboxStyled(true))
	if !c.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// Checkbox Method Tests
// =============================================================================

func TestCheckboxToggle(t *testing.T) {
	c := NewCheckbox("test")
	c.Toggle()
	if !c.Checked {
		t.Error("expected checked after toggle")
	}
}

func TestCheckboxSetIndeterminate(t *testing.T) {
	c := NewCheckbox("test", WithCheckboxChecked(true))
	c.SetIndeterminate(true)
	if !c.Indeterminate {
		t.Error("expected indeterminate")
	}
	if c.Checked {
		t.Error("expected unchecked when indeterminate")
	}
}

func TestCheckboxHasLabel(t *testing.T) {
	c := NewCheckbox("test")
	if c.HasLabel() {
		t.Error("expected HasLabel false")
	}
	c.Label = "Test"
	if !c.HasLabel() {
		t.Error("expected HasLabel true")
	}
}

func TestCheckboxStateClass(t *testing.T) {
	t.Run("unchecked enabled", func(t *testing.T) {
		c := NewCheckbox("test")
		expected := "bg-white border-gray-300"
		if c.CheckboxStateClass() != expected {
			t.Errorf("expected %q, got %q", expected, c.CheckboxStateClass())
		}
	})

	t.Run("checked enabled", func(t *testing.T) {
		c := NewCheckbox("test", WithCheckboxChecked(true))
		expected := "bg-blue-600 border-blue-600"
		if c.CheckboxStateClass() != expected {
			t.Errorf("expected %q, got %q", expected, c.CheckboxStateClass())
		}
	})
}

// =============================================================================
// Template Tests
// =============================================================================

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Templates() returned nil")
	}
}
