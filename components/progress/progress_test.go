package progress

import (
	"testing"
)

// =============================================================================
// Progress Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates progress with defaults", func(t *testing.T) {
		p := New("download")
		if p.ID() != "download" {
			t.Errorf("expected ID 'download', got %q", p.ID())
		}
		if p.Namespace() != "progress" {
			t.Errorf("expected namespace 'progress', got %q", p.Namespace())
		}
		if p.Value != 0 {
			t.Errorf("expected value 0, got %f", p.Value)
		}
		if p.Max != 100 {
			t.Errorf("expected max 100, got %f", p.Max)
		}
		if p.Size != SizeMd {
			t.Errorf("expected size md, got %v", p.Size)
		}
		if p.Color != ColorPrimary {
			t.Errorf("expected color primary, got %v", p.Color)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		p := New("upload",
			WithValue(75),
			WithMax(200),
			WithSize(SizeLg),
			WithColor(ColorSuccess),
			WithShowLabel(true),
			WithStriped(true),
		)
		if p.Value != 75 {
			t.Errorf("expected value 75, got %f", p.Value)
		}
		if p.Max != 200 {
			t.Errorf("expected max 200, got %f", p.Max)
		}
		if p.Size != SizeLg {
			t.Errorf("expected size lg, got %v", p.Size)
		}
		if p.Color != ColorSuccess {
			t.Errorf("expected color success, got %v", p.Color)
		}
		if !p.ShowLabel {
			t.Error("expected ShowLabel true")
		}
		if !p.Striped {
			t.Error("expected Striped true")
		}
	})
}

// =============================================================================
// Progress Option Tests
// =============================================================================

func TestWithValue(t *testing.T) {
	p := New("test", WithValue(50))
	if p.Value != 50 {
		t.Errorf("expected value 50, got %f", p.Value)
	}
}

func TestWithMax(t *testing.T) {
	p := New("test", WithMax(500))
	if p.Max != 500 {
		t.Errorf("expected max 500, got %f", p.Max)
	}
}

func TestWithSize(t *testing.T) {
	sizes := []Size{SizeXs, SizeSm, SizeMd, SizeLg}

	for _, size := range sizes {
		t.Run(string(size), func(t *testing.T) {
			p := New("test", WithSize(size))
			if p.Size != size {
				t.Errorf("expected size %v, got %v", size, p.Size)
			}
		})
	}
}

func TestWithColor(t *testing.T) {
	colors := []Color{ColorPrimary, ColorSuccess, ColorWarning, ColorDanger, ColorInfo}

	for _, color := range colors {
		t.Run(string(color), func(t *testing.T) {
			p := New("test", WithColor(color))
			if p.Color != color {
				t.Errorf("expected color %v, got %v", color, p.Color)
			}
		})
	}
}

func TestWithShowLabel(t *testing.T) {
	p := New("test", WithShowLabel(true))
	if !p.ShowLabel {
		t.Error("expected ShowLabel true")
	}
}

func TestWithLabel(t *testing.T) {
	p := New("test", WithLabel("Uploading..."))
	if p.Label != "Uploading..." {
		t.Errorf("expected label 'Uploading...', got %q", p.Label)
	}
}

func TestWithStriped(t *testing.T) {
	p := New("test", WithStriped(true))
	if !p.Striped {
		t.Error("expected Striped true")
	}
}

func TestWithAnimated(t *testing.T) {
	p := New("test", WithAnimated(true))
	if !p.Animated {
		t.Error("expected Animated true")
	}
}

func TestWithIndeterminate(t *testing.T) {
	p := New("test", WithIndeterminate(true))
	if !p.Indeterminate {
		t.Error("expected Indeterminate true")
	}
}

func TestWithStyled(t *testing.T) {
	p := New("test", WithStyled(true))
	if !p.IsStyled() {
		t.Error("expected styled true")
	}
}

// =============================================================================
// Progress Method Tests
// =============================================================================

func TestPercentage(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		max      float64
		expected float64
	}{
		{"zero", 0, 100, 0},
		{"half", 50, 100, 50},
		{"full", 100, 100, 100},
		{"over max", 150, 100, 100},
		{"negative", -10, 100, 0},
		{"zero max", 50, 0, 0},
		{"custom max", 25, 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New("test", WithValue(tt.value), WithMax(tt.max))
			if p.Percentage() != tt.expected {
				t.Errorf("expected percentage %f, got %f", tt.expected, p.Percentage())
			}
		})
	}
}

func TestPercentageStr(t *testing.T) {
	p := New("test", WithValue(75))
	expected := "75%"
	if p.PercentageStr() != expected {
		t.Errorf("expected %q, got %q", expected, p.PercentageStr())
	}
}

func TestDisplayLabel(t *testing.T) {
	t.Run("returns custom label when set", func(t *testing.T) {
		p := New("test", WithValue(50), WithLabel("Custom"))
		if p.DisplayLabel() != "Custom" {
			t.Errorf("expected 'Custom', got %q", p.DisplayLabel())
		}
	})

	t.Run("returns percentage when no label", func(t *testing.T) {
		p := New("test", WithValue(50))
		if p.DisplayLabel() != "50%" {
			t.Errorf("expected '50%%', got %q", p.DisplayLabel())
		}
	})
}

func TestSetValue(t *testing.T) {
	p := New("test")
	p.SetValue(75)
	if p.Value != 75 {
		t.Errorf("expected value 75, got %f", p.Value)
	}
}

func TestIncrement(t *testing.T) {
	p := New("test", WithValue(50))
	p.Increment(10)
	if p.Value != 60 {
		t.Errorf("expected value 60, got %f", p.Value)
	}

	// Should not exceed max
	p.Increment(100)
	if p.Value != 100 {
		t.Errorf("expected value 100, got %f", p.Value)
	}
}

func TestDecrement(t *testing.T) {
	p := New("test", WithValue(50))
	p.Decrement(10)
	if p.Value != 40 {
		t.Errorf("expected value 40, got %f", p.Value)
	}

	// Should not go below 0
	p.Decrement(100)
	if p.Value != 0 {
		t.Errorf("expected value 0, got %f", p.Value)
	}
}

func TestReset(t *testing.T) {
	p := New("test", WithValue(75))
	p.Reset()
	if p.Value != 0 {
		t.Errorf("expected value 0, got %f", p.Value)
	}
}

func TestComplete(t *testing.T) {
	p := New("test", WithValue(50))
	p.Complete()
	if p.Value != 100 {
		t.Errorf("expected value 100, got %f", p.Value)
	}
}

func TestIsComplete(t *testing.T) {
	p := New("test", WithValue(50))
	if p.IsComplete() {
		t.Error("expected IsComplete false at 50%")
	}

	p.Value = 100
	if !p.IsComplete() {
		t.Error("expected IsComplete true at 100%")
	}

	p.Value = 150
	if !p.IsComplete() {
		t.Error("expected IsComplete true when over max")
	}
}

func TestSizeClass(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeXs, "h-1"},
		{SizeSm, "h-2"},
		{SizeMd, "h-4"},
		{SizeLg, "h-6"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			p := New("test", WithSize(tt.size))
			if p.SizeClass() != tt.expected {
				t.Errorf("expected SizeClass %q, got %q", tt.expected, p.SizeClass())
			}
		})
	}
}

func TestColorClass(t *testing.T) {
	tests := []struct {
		color    Color
		expected string
	}{
		{ColorPrimary, "bg-blue-500"},
		{ColorSuccess, "bg-green-500"},
		{ColorWarning, "bg-yellow-500"},
		{ColorDanger, "bg-red-500"},
		{ColorInfo, "bg-cyan-500"},
	}

	for _, tt := range tests {
		t.Run(string(tt.color), func(t *testing.T) {
			p := New("test", WithColor(tt.color))
			if p.ColorClass() != tt.expected {
				t.Errorf("expected ColorClass %q, got %q", tt.expected, p.ColorClass())
			}
		})
	}
}

// =============================================================================
// Circular Progress Tests
// =============================================================================

func TestNewCircular(t *testing.T) {
	t.Run("creates circular with defaults", func(t *testing.T) {
		c := NewCircular("loading")
		if c.ID() != "loading" {
			t.Errorf("expected ID 'loading', got %q", c.ID())
		}
		if c.Value != 0 {
			t.Errorf("expected value 0, got %f", c.Value)
		}
		if c.Max != 100 {
			t.Errorf("expected max 100, got %f", c.Max)
		}
		if c.Size != 48 {
			t.Errorf("expected size 48, got %d", c.Size)
		}
		if c.StrokeWidth != 4 {
			t.Errorf("expected stroke width 4, got %d", c.StrokeWidth)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		c := NewCircular("progress",
			WithCircularValue(75),
			WithCircularSize(80),
			WithCircularStrokeWidth(6),
			WithCircularColor(ColorSuccess),
			WithCircularShowLabel(true),
		)
		if c.Value != 75 {
			t.Errorf("expected value 75, got %f", c.Value)
		}
		if c.Size != 80 {
			t.Errorf("expected size 80, got %d", c.Size)
		}
		if c.StrokeWidth != 6 {
			t.Errorf("expected stroke width 6, got %d", c.StrokeWidth)
		}
		if c.Color != ColorSuccess {
			t.Errorf("expected color success, got %v", c.Color)
		}
		if !c.ShowLabel {
			t.Error("expected ShowLabel true")
		}
	})
}

func TestCircularPercentage(t *testing.T) {
	c := NewCircular("test", WithCircularValue(50))
	if c.Percentage() != 50 {
		t.Errorf("expected percentage 50, got %f", c.Percentage())
	}
}

func TestCircularRadius(t *testing.T) {
	c := NewCircular("test", WithCircularSize(48), WithCircularStrokeWidth(4))
	expected := (48 - 4) / 2 // 22
	if c.Radius() != expected {
		t.Errorf("expected radius %d, got %d", expected, c.Radius())
	}
}

func TestCircularCenter(t *testing.T) {
	c := NewCircular("test", WithCircularSize(48))
	if c.Center() != 24 {
		t.Errorf("expected center 24, got %d", c.Center())
	}
}

func TestCircularCircumference(t *testing.T) {
	c := NewCircular("test", WithCircularSize(48), WithCircularStrokeWidth(4))
	// 2 * pi * 22 ≈ 138.23
	circ := c.Circumference()
	if circ < 138 || circ > 139 {
		t.Errorf("expected circumference ~138.23, got %f", circ)
	}
}

func TestCircularColorClass(t *testing.T) {
	c := NewCircular("test", WithCircularColor(ColorDanger))
	if c.ColorClass() != "text-red-500" {
		t.Errorf("expected 'text-red-500', got %q", c.ColorClass())
	}
}

func TestWithCircularStyled(t *testing.T) {
	c := NewCircular("test", WithCircularStyled(true))
	if !c.IsStyled() {
		t.Error("expected circular styled true")
	}
}

func TestWithCircularIndeterminate(t *testing.T) {
	c := NewCircular("test", WithCircularIndeterminate(true))
	if !c.Indeterminate {
		t.Error("expected Indeterminate true")
	}
}

func TestWithCircularLabel(t *testing.T) {
	c := NewCircular("test", WithCircularLabel("Custom"))
	if c.Label != "Custom" {
		t.Errorf("expected label 'Custom', got %q", c.Label)
	}
}

// =============================================================================
// Spinner Tests
// =============================================================================

func TestNewSpinner(t *testing.T) {
	t.Run("creates spinner with defaults", func(t *testing.T) {
		s := NewSpinner("loading")
		if s.ID() != "loading" {
			t.Errorf("expected ID 'loading', got %q", s.ID())
		}
		if s.Size != "md" {
			t.Errorf("expected size 'md', got %q", s.Size)
		}
		if s.Color != ColorPrimary {
			t.Errorf("expected color primary, got %v", s.Color)
		}
		if s.Label != "Loading..." {
			t.Errorf("expected label 'Loading...', got %q", s.Label)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		s := NewSpinner("wait",
			WithSpinnerSize("lg"),
			WithSpinnerColor(ColorInfo),
			WithSpinnerLabel("Please wait..."),
		)
		if s.Size != "lg" {
			t.Errorf("expected size 'lg', got %q", s.Size)
		}
		if s.Color != ColorInfo {
			t.Errorf("expected color info, got %v", s.Color)
		}
		if s.Label != "Please wait..." {
			t.Errorf("expected label 'Please wait...', got %q", s.Label)
		}
	})
}

func TestSpinnerSizeClass(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"sm", "w-4 h-4"},
		{"md", "w-6 h-6"},
		{"lg", "w-8 h-8"},
		{"xl", "w-12 h-12"},
	}

	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			s := NewSpinner("test", WithSpinnerSize(tt.size))
			if s.SizeClass() != tt.expected {
				t.Errorf("expected SizeClass %q, got %q", tt.expected, s.SizeClass())
			}
		})
	}
}

func TestSpinnerColorClass(t *testing.T) {
	s := NewSpinner("test", WithSpinnerColor(ColorSuccess))
	if s.ColorClass() != "text-green-500" {
		t.Errorf("expected 'text-green-500', got %q", s.ColorClass())
	}
}

func TestWithSpinnerStyled(t *testing.T) {
	s := NewSpinner("test", WithSpinnerStyled(true))
	if !s.IsStyled() {
		t.Error("expected spinner styled true")
	}
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
