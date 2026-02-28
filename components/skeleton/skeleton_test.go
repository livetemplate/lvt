package skeleton

import (
	"testing"
)

// =============================================================================
// Skeleton Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates skeleton with defaults", func(t *testing.T) {
		s := New("loading")
		if s.ID() != "loading" {
			t.Errorf("expected ID 'loading', got %q", s.ID())
		}
		if s.Namespace() != "skeleton" {
			t.Errorf("expected namespace 'skeleton', got %q", s.Namespace())
		}
		if s.Width != "100%" {
			t.Errorf("expected width '100%%', got %q", s.Width)
		}
		if s.Height != "16px" {
			t.Errorf("expected height '16px', got %q", s.Height)
		}
		if s.Shape != ShapeRectangle {
			t.Errorf("expected shape rectangle, got %v", s.Shape)
		}
		if s.Animation != AnimationPulse {
			t.Errorf("expected animation pulse, got %v", s.Animation)
		}
		if s.Lines != 1 {
			t.Errorf("expected 1 line, got %d", s.Lines)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		s := New("text",
			WithWidth("200px"),
			WithHeight("20px"),
			WithShape(ShapeRounded),
			WithAnimation(AnimationWave),
			WithLines(3),
		)
		if s.Width != "200px" {
			t.Errorf("expected width '200px', got %q", s.Width)
		}
		if s.Height != "20px" {
			t.Errorf("expected height '20px', got %q", s.Height)
		}
		if s.Shape != ShapeRounded {
			t.Errorf("expected shape rounded, got %v", s.Shape)
		}
		if s.Animation != AnimationWave {
			t.Errorf("expected animation wave, got %v", s.Animation)
		}
		if s.Lines != 3 {
			t.Errorf("expected 3 lines, got %d", s.Lines)
		}
	})
}

// =============================================================================
// Skeleton Option Tests
// =============================================================================

func TestWithWidth(t *testing.T) {
	s := New("test", WithWidth("50%"))
	if s.Width != "50%" {
		t.Errorf("expected width '50%%', got %q", s.Width)
	}
}

func TestWithHeight(t *testing.T) {
	s := New("test", WithHeight("24px"))
	if s.Height != "24px" {
		t.Errorf("expected height '24px', got %q", s.Height)
	}
}

func TestWithShape(t *testing.T) {
	shapes := []Shape{ShapeRectangle, ShapeCircle, ShapeRounded}

	for _, shape := range shapes {
		t.Run(string(shape), func(t *testing.T) {
			s := New("test", WithShape(shape))
			if s.Shape != shape {
				t.Errorf("expected shape %v, got %v", shape, s.Shape)
			}
		})
	}
}

func TestWithAnimation(t *testing.T) {
	animations := []Animation{AnimationPulse, AnimationWave, AnimationNone}

	for _, anim := range animations {
		t.Run(string(anim), func(t *testing.T) {
			s := New("test", WithAnimation(anim))
			if s.Animation != anim {
				t.Errorf("expected animation %v, got %v", anim, s.Animation)
			}
		})
	}
}

func TestWithLines(t *testing.T) {
	s := New("test", WithLines(5))
	if s.Lines != 5 {
		t.Errorf("expected 5 lines, got %d", s.Lines)
	}
}

func TestWithLineHeight(t *testing.T) {
	s := New("test", WithLineHeight("32px"))
	if s.LineHeight != "32px" {
		t.Errorf("expected line height '32px', got %q", s.LineHeight)
	}
}

func TestWithStyled(t *testing.T) {
	s := New("test", WithStyled(true))
	if !s.IsStyled() {
		t.Error("expected styled true")
	}

	s2 := New("test", WithStyled(false))
	if s2.IsStyled() {
		t.Error("expected styled false")
	}
}

// =============================================================================
// Skeleton Method Tests
// =============================================================================

func TestShapeClass(t *testing.T) {
	tests := []struct {
		shape    Shape
		expected string
	}{
		{ShapeRectangle, ""},
		{ShapeCircle, "rounded-full"},
		{ShapeRounded, "rounded-md"},
	}

	for _, tt := range tests {
		t.Run(string(tt.shape), func(t *testing.T) {
			s := New("test", WithShape(tt.shape))
			if s.ShapeClass() != tt.expected {
				t.Errorf("expected ShapeClass %q, got %q", tt.expected, s.ShapeClass())
			}
		})
	}
}

func TestAnimationClass(t *testing.T) {
	tests := []struct {
		animation Animation
		expected  string
	}{
		{AnimationPulse, "animate-pulse"},
		{AnimationWave, "animate-shimmer"},
		{AnimationNone, ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.animation), func(t *testing.T) {
			s := New("test", WithAnimation(tt.animation))
			if s.AnimationClass() != tt.expected {
				t.Errorf("expected AnimationClass %q, got %q", tt.expected, s.AnimationClass())
			}
		})
	}
}

func TestIsCircle(t *testing.T) {
	s := New("test", WithShape(ShapeCircle))
	if !s.IsCircle() {
		t.Error("expected IsCircle true")
	}
	s.Shape = ShapeRectangle
	if s.IsCircle() {
		t.Error("expected IsCircle false for rectangle")
	}
}

func TestIsRounded(t *testing.T) {
	s := New("test", WithShape(ShapeRounded))
	if !s.IsRounded() {
		t.Error("expected IsRounded true")
	}
	s.Shape = ShapeRectangle
	if s.IsRounded() {
		t.Error("expected IsRounded false for rectangle")
	}
}

func TestIsMultiLine(t *testing.T) {
	s := New("test")
	if s.IsMultiLine() {
		t.Error("expected IsMultiLine false with 1 line")
	}

	s.Lines = 3
	if !s.IsMultiLine() {
		t.Error("expected IsMultiLine true with 3 lines")
	}
}

func TestLineIndices(t *testing.T) {
	s := New("test", WithLines(4))
	indices := s.LineIndices()
	if len(indices) != 4 {
		t.Errorf("expected 4 indices, got %d", len(indices))
	}
	for i, idx := range indices {
		if idx != i {
			t.Errorf("expected index %d at position %d, got %d", i, i, idx)
		}
	}
}

// =============================================================================
// Avatar Skeleton Tests
// =============================================================================

func TestNewAvatar(t *testing.T) {
	t.Run("creates avatar with defaults", func(t *testing.T) {
		a := NewAvatar("user")
		if a.ID() != "user" {
			t.Errorf("expected ID 'user', got %q", a.ID())
		}
		if a.Size != "md" {
			t.Errorf("expected size 'md', got %q", a.Size)
		}
		if a.ShowBadge {
			t.Error("expected ShowBadge false by default")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		a := NewAvatar("profile",
			WithAvatarSize("lg"),
			WithAvatarBadge(true),
		)
		if a.Size != "lg" {
			t.Errorf("expected size 'lg', got %q", a.Size)
		}
		if !a.ShowBadge {
			t.Error("expected ShowBadge true")
		}
	})
}

func TestAvatarSizeClass(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"sm", "w-8 h-8"},
		{"md", "w-12 h-12"},
		{"lg", "w-16 h-16"},
		{"xl", "w-24 h-24"},
	}

	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			a := NewAvatar("test", WithAvatarSize(tt.size))
			if a.SizeClass() != tt.expected {
				t.Errorf("expected SizeClass %q, got %q", tt.expected, a.SizeClass())
			}
		})
	}
}

func TestWithAvatarStyled(t *testing.T) {
	a := NewAvatar("test", WithAvatarStyled(true))
	if !a.IsStyled() {
		t.Error("expected avatar styled true")
	}
}

// =============================================================================
// Card Skeleton Tests
// =============================================================================

func TestNewCard(t *testing.T) {
	t.Run("creates card with defaults", func(t *testing.T) {
		c := NewCard("content")
		if c.ID() != "content" {
			t.Errorf("expected ID 'content', got %q", c.ID())
		}
		if !c.ShowImage {
			t.Error("expected ShowImage true by default")
		}
		if c.ImageHeight != "200px" {
			t.Errorf("expected ImageHeight '200px', got %q", c.ImageHeight)
		}
		if !c.ShowTitle {
			t.Error("expected ShowTitle true by default")
		}
		if !c.ShowDescription {
			t.Error("expected ShowDescription true by default")
		}
		if c.DescriptionLines != 3 {
			t.Errorf("expected 3 description lines, got %d", c.DescriptionLines)
		}
		if c.ShowFooter {
			t.Error("expected ShowFooter false by default")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		c := NewCard("post",
			WithCardImage(false),
			WithCardTitle(false),
			WithCardDescription(true, 5),
			WithCardFooter(true),
		)
		if c.ShowImage {
			t.Error("expected ShowImage false")
		}
		if c.ShowTitle {
			t.Error("expected ShowTitle false")
		}
		if !c.ShowDescription {
			t.Error("expected ShowDescription true")
		}
		if c.DescriptionLines != 5 {
			t.Errorf("expected 5 description lines, got %d", c.DescriptionLines)
		}
		if !c.ShowFooter {
			t.Error("expected ShowFooter true")
		}
	})
}

func TestWithCardImageHeight(t *testing.T) {
	c := NewCard("test", WithCardImageHeight("300px"))
	if c.ImageHeight != "300px" {
		t.Errorf("expected ImageHeight '300px', got %q", c.ImageHeight)
	}
}

func TestWithCardStyled(t *testing.T) {
	c := NewCard("test", WithCardStyled(true))
	if !c.IsStyled() {
		t.Error("expected card styled true")
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
