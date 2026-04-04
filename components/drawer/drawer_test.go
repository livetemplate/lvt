package drawer

import (
	"testing"

	// Register the tailwind adapter so CSS class methods return expected values.
	_ "github.com/livetemplate/lvt/components/styles/tailwind"
)

// =============================================================================
// Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates drawer with defaults", func(t *testing.T) {
		d := New("sidebar")
		if d.ID() != "sidebar" {
			t.Errorf("expected ID 'sidebar', got %q", d.ID())
		}
		if d.Namespace() != "drawer" {
			t.Errorf("expected namespace 'drawer', got %q", d.Namespace())
		}
		if d.Position != PositionLeft {
			t.Errorf("expected position left, got %v", d.Position)
		}
		if d.Size != SizeMd {
			t.Errorf("expected size md, got %v", d.Size)
		}
		if !d.ShowClose {
			t.Error("expected ShowClose true by default")
		}
		if !d.ShowOverlay {
			t.Error("expected ShowOverlay true by default")
		}
		if !d.CloseOnOverlay {
			t.Error("expected CloseOnOverlay true by default")
		}
		if !d.CloseOnEscape {
			t.Error("expected CloseOnEscape true by default")
		}
		if d.Persistent {
			t.Error("expected Persistent false by default")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		d := New("panel",
			WithPosition(PositionRight),
			WithSize(SizeLg),
			WithTitle("Settings"),
			WithShowClose(false),
		)
		if d.Position != PositionRight {
			t.Errorf("expected position right, got %v", d.Position)
		}
		if d.Size != SizeLg {
			t.Errorf("expected size lg, got %v", d.Size)
		}
		if d.Title != "Settings" {
			t.Errorf("expected title 'Settings', got %q", d.Title)
		}
		if d.ShowClose {
			t.Error("expected ShowClose false")
		}
	})
}

// =============================================================================
// Option Tests
// =============================================================================

func TestWithPosition(t *testing.T) {
	tests := []struct {
		name     string
		position Position
	}{
		{"left", PositionLeft},
		{"right", PositionRight},
		{"top", PositionTop},
		{"bottom", PositionBottom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New("test", WithPosition(tt.position))
			if d.Position != tt.position {
				t.Errorf("expected position %v, got %v", tt.position, d.Position)
			}
		})
	}
}

func TestWithSize(t *testing.T) {
	tests := []struct {
		name string
		size Size
	}{
		{"sm", SizeSm},
		{"md", SizeMd},
		{"lg", SizeLg},
		{"xl", SizeXl},
		{"full", SizeFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New("test", WithSize(tt.size))
			if d.Size != tt.size {
				t.Errorf("expected size %v, got %v", tt.size, d.Size)
			}
		})
	}
}

func TestWithTitle(t *testing.T) {
	d := New("test", WithTitle("Navigation"))
	if d.Title != "Navigation" {
		t.Errorf("expected title 'Navigation', got %q", d.Title)
	}
}

func TestWithShowClose(t *testing.T) {
	d := New("test", WithShowClose(false))
	if d.ShowClose {
		t.Error("expected ShowClose false")
	}
}

func TestWithShowOverlay(t *testing.T) {
	d := New("test", WithShowOverlay(false))
	if d.ShowOverlay {
		t.Error("expected ShowOverlay false")
	}
}

func TestWithCloseOnOverlay(t *testing.T) {
	d := New("test", WithCloseOnOverlay(false))
	if d.CloseOnOverlay {
		t.Error("expected CloseOnOverlay false")
	}
}

func TestWithCloseOnEscape(t *testing.T) {
	d := New("test", WithCloseOnEscape(false))
	if d.CloseOnEscape {
		t.Error("expected CloseOnEscape false")
	}
}

func TestWithPersistent(t *testing.T) {
	d := New("test", WithPersistent(true))
	if !d.Persistent {
		t.Error("expected Persistent true")
	}
}

func TestWithOpen(t *testing.T) {
	// WithOpen is a no-op; open/close is client-side now
	d := New("test", WithOpen(true))
	_ = d // should not panic
}

func TestWithStyled(t *testing.T) {
	d := New("test", WithStyled(true))
	if !d.IsStyled() {
		t.Error("expected styled true")
	}

	d2 := New("test", WithStyled(false))
	if d2.IsStyled() {
		t.Error("expected styled false")
	}
}

// =============================================================================
// Method Tests
// =============================================================================

// Toggle, Show, Close, ForceClose methods have been removed.
// Open/close is now handled client-side via CSS classes.

// =============================================================================
// Position Helper Tests
// =============================================================================

func TestIsLeft(t *testing.T) {
	d := New("test", WithPosition(PositionLeft))
	if !d.IsLeft() {
		t.Error("expected IsLeft true")
	}
	d.Position = PositionRight
	if d.IsLeft() {
		t.Error("expected IsLeft false for right position")
	}
}

func TestIsRight(t *testing.T) {
	d := New("test", WithPosition(PositionRight))
	if !d.IsRight() {
		t.Error("expected IsRight true")
	}
	d.Position = PositionLeft
	if d.IsRight() {
		t.Error("expected IsRight false for left position")
	}
}

func TestIsTop(t *testing.T) {
	d := New("test", WithPosition(PositionTop))
	if !d.IsTop() {
		t.Error("expected IsTop true")
	}
	d.Position = PositionBottom
	if d.IsTop() {
		t.Error("expected IsTop false for bottom position")
	}
}

func TestIsBottom(t *testing.T) {
	d := New("test", WithPosition(PositionBottom))
	if !d.IsBottom() {
		t.Error("expected IsBottom true")
	}
	d.Position = PositionTop
	if d.IsBottom() {
		t.Error("expected IsBottom false for top position")
	}
}

func TestIsHorizontal(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionLeft, true},
		{PositionRight, true},
		{PositionTop, false},
		{PositionBottom, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			d := New("test", WithPosition(tt.position))
			if d.IsHorizontal() != tt.expected {
				t.Errorf("expected IsHorizontal %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

func TestIsVertical(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionLeft, false},
		{PositionRight, false},
		{PositionTop, true},
		{PositionBottom, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			d := New("test", WithPosition(tt.position))
			if d.IsVertical() != tt.expected {
				t.Errorf("expected IsVertical %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

// =============================================================================
// CSS Class Helper Tests
// =============================================================================

func TestSizeClassHorizontal(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "w-64"},
		{SizeMd, "w-80"},
		{SizeLg, "w-96"},
		{SizeXl, "w-[32rem]"},
		{SizeFull, "w-full"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			d := New("test", WithPosition(PositionLeft), WithSize(tt.size))
			if d.SizeClass() != tt.expected {
				t.Errorf("expected SizeClass %q, got %q", tt.expected, d.SizeClass())
			}
		})
	}
}

func TestSizeClassVertical(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "h-48"},
		{SizeMd, "h-64"},
		{SizeLg, "h-96"},
		{SizeXl, "h-[32rem]"},
		{SizeFull, "h-full"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			d := New("test", WithPosition(PositionTop), WithSize(tt.size))
			if d.SizeClass() != tt.expected {
				t.Errorf("expected SizeClass %q, got %q", tt.expected, d.SizeClass())
			}
		})
	}
}

func TestPositionClass(t *testing.T) {
	tests := []struct {
		position Position
		expected string
	}{
		{PositionLeft, "left-0 top-0 h-full"},
		{PositionRight, "right-0 top-0 h-full"},
		{PositionTop, "top-0 left-0 w-full"},
		{PositionBottom, "bottom-0 left-0 w-full"},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			d := New("test", WithPosition(tt.position))
			if d.PositionClass() != tt.expected {
				t.Errorf("expected PositionClass %q, got %q", tt.expected, d.PositionClass())
			}
		})
	}
}

func TestTransformClass(t *testing.T) {
	tests := []struct {
		position Position
		expected string
	}{
		{PositionLeft, "translateX(-100%)"},
		{PositionRight, "translateX(100%)"},
		{PositionTop, "translateY(-100%)"},
		{PositionBottom, "translateY(100%)"},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			d := New("test", WithPosition(tt.position))
			if d.TransformClass() != tt.expected {
				t.Errorf("expected TransformClass %q, got %q", tt.expected, d.TransformClass())
			}
		})
	}
}

// =============================================================================
// HasTitle Tests
// =============================================================================

func TestHasTitle(t *testing.T) {
	d := New("test")
	if d.HasTitle() {
		t.Error("expected HasTitle false when no title")
	}

	d.Title = "Menu"
	if !d.HasTitle() {
		t.Error("expected HasTitle true when title set")
	}

	d.Title = ""
	if d.HasTitle() {
		t.Error("expected HasTitle false when title empty")
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
