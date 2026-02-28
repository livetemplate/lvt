package popover

import (
	"testing"
)

// =============================================================================
// Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates popover with defaults", func(t *testing.T) {
		p := New("info")
		if p.ID() != "info" {
			t.Errorf("expected ID 'info', got %q", p.ID())
		}
		if p.Namespace() != "popover" {
			t.Errorf("expected namespace 'popover', got %q", p.Namespace())
		}
		if p.Position != PositionBottom {
			t.Errorf("expected position bottom, got %v", p.Position)
		}
		if p.Trigger != TriggerClick {
			t.Errorf("expected trigger click, got %v", p.Trigger)
		}
		if p.Open {
			t.Error("expected popover to be closed by default")
		}
		if !p.Arrow {
			t.Error("expected arrow true by default")
		}
		if !p.CloseOnClickAway {
			t.Error("expected CloseOnClickAway true by default")
		}
		if p.ShowClose {
			t.Error("expected ShowClose false by default")
		}
		if p.Width != "280px" {
			t.Errorf("expected width '280px', got %q", p.Width)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		p := New("details",
			WithTitle("Details"),
			WithContent("More info here"),
			WithPosition(PositionTop),
			WithTrigger(TriggerHover),
			WithShowClose(true),
		)
		if p.Title != "Details" {
			t.Errorf("expected title 'Details', got %q", p.Title)
		}
		if p.Content != "More info here" {
			t.Errorf("expected content 'More info here', got %q", p.Content)
		}
		if p.Position != PositionTop {
			t.Errorf("expected position top, got %v", p.Position)
		}
		if p.Trigger != TriggerHover {
			t.Errorf("expected trigger hover, got %v", p.Trigger)
		}
		if !p.ShowClose {
			t.Error("expected ShowClose true")
		}
	})
}

// =============================================================================
// Option Tests
// =============================================================================

func TestWithTitle(t *testing.T) {
	p := New("test", WithTitle("Header"))
	if p.Title != "Header" {
		t.Errorf("expected title 'Header', got %q", p.Title)
	}
}

func TestWithContent(t *testing.T) {
	p := New("test", WithContent("Body text"))
	if p.Content != "Body text" {
		t.Errorf("expected content 'Body text', got %q", p.Content)
	}
}

func TestWithPosition(t *testing.T) {
	positions := []Position{
		PositionTop, PositionTopStart, PositionTopEnd,
		PositionBottom, PositionBottomStart, PositionBottomEnd,
		PositionLeft, PositionLeftStart, PositionLeftEnd,
		PositionRight, PositionRightStart, PositionRightEnd,
	}

	for _, pos := range positions {
		t.Run(string(pos), func(t *testing.T) {
			p := New("test", WithPosition(pos))
			if p.Position != pos {
				t.Errorf("expected position %v, got %v", pos, p.Position)
			}
		})
	}
}

func TestWithTrigger(t *testing.T) {
	triggers := []Trigger{TriggerClick, TriggerHover, TriggerFocus}

	for _, trigger := range triggers {
		t.Run(string(trigger), func(t *testing.T) {
			p := New("test", WithTrigger(trigger))
			if p.Trigger != trigger {
				t.Errorf("expected trigger %v, got %v", trigger, p.Trigger)
			}
		})
	}
}

func TestWithArrow(t *testing.T) {
	p := New("test", WithArrow(false))
	if p.Arrow {
		t.Error("expected arrow false")
	}
}

func TestWithCloseOnClickAway(t *testing.T) {
	p := New("test", WithCloseOnClickAway(false))
	if p.CloseOnClickAway {
		t.Error("expected CloseOnClickAway false")
	}
}

func TestWithShowClose(t *testing.T) {
	p := New("test", WithShowClose(true))
	if !p.ShowClose {
		t.Error("expected ShowClose true")
	}
}

func TestWithWidth(t *testing.T) {
	p := New("test", WithWidth("400px"))
	if p.Width != "400px" {
		t.Errorf("expected width '400px', got %q", p.Width)
	}
}

func TestWithOpen(t *testing.T) {
	p := New("test", WithOpen(true))
	if !p.Open {
		t.Error("expected Open true")
	}
}

func TestWithStyled(t *testing.T) {
	p := New("test", WithStyled(true))
	if !p.IsStyled() {
		t.Error("expected styled true")
	}

	p2 := New("test", WithStyled(false))
	if p2.IsStyled() {
		t.Error("expected styled false")
	}
}

// =============================================================================
// Method Tests
// =============================================================================

func TestShow(t *testing.T) {
	p := New("test")
	p.Show()
	if !p.Open {
		t.Error("expected open after Show")
	}
}

func TestHide(t *testing.T) {
	p := New("test", WithOpen(true))
	p.Hide()
	if p.Open {
		t.Error("expected closed after Hide")
	}
}

func TestToggle(t *testing.T) {
	p := New("test")
	if p.Open {
		t.Error("expected initially closed")
	}

	p.Toggle()
	if !p.Open {
		t.Error("expected open after toggle")
	}

	p.Toggle()
	if p.Open {
		t.Error("expected closed after second toggle")
	}
}

// =============================================================================
// Position Helper Tests
// =============================================================================

func TestIsTop(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionTop, true},
		{PositionTopStart, true},
		{PositionTopEnd, true},
		{PositionBottom, false},
		{PositionLeft, false},
		{PositionRight, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			p := New("test", WithPosition(tt.position))
			if p.IsTop() != tt.expected {
				t.Errorf("expected IsTop %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

func TestIsBottom(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionBottom, true},
		{PositionBottomStart, true},
		{PositionBottomEnd, true},
		{PositionTop, false},
		{PositionLeft, false},
		{PositionRight, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			p := New("test", WithPosition(tt.position))
			if p.IsBottom() != tt.expected {
				t.Errorf("expected IsBottom %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

func TestIsLeft(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionLeft, true},
		{PositionLeftStart, true},
		{PositionLeftEnd, true},
		{PositionTop, false},
		{PositionBottom, false},
		{PositionRight, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			p := New("test", WithPosition(tt.position))
			if p.IsLeft() != tt.expected {
				t.Errorf("expected IsLeft %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

func TestIsRight(t *testing.T) {
	tests := []struct {
		position Position
		expected bool
	}{
		{PositionRight, true},
		{PositionRightStart, true},
		{PositionRightEnd, true},
		{PositionTop, false},
		{PositionBottom, false},
		{PositionLeft, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			p := New("test", WithPosition(tt.position))
			if p.IsRight() != tt.expected {
				t.Errorf("expected IsRight %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

// =============================================================================
// Trigger Helper Tests
// =============================================================================

func TestIsClickTrigger(t *testing.T) {
	p := New("test", WithTrigger(TriggerClick))
	if !p.IsClickTrigger() {
		t.Error("expected IsClickTrigger true")
	}
	p.Trigger = TriggerHover
	if p.IsClickTrigger() {
		t.Error("expected IsClickTrigger false for hover")
	}
}

func TestIsHoverTrigger(t *testing.T) {
	p := New("test", WithTrigger(TriggerHover))
	if !p.IsHoverTrigger() {
		t.Error("expected IsHoverTrigger true")
	}
	p.Trigger = TriggerClick
	if p.IsHoverTrigger() {
		t.Error("expected IsHoverTrigger false for click")
	}
}

func TestIsFocusTrigger(t *testing.T) {
	p := New("test", WithTrigger(TriggerFocus))
	if !p.IsFocusTrigger() {
		t.Error("expected IsFocusTrigger true")
	}
	p.Trigger = TriggerClick
	if p.IsFocusTrigger() {
		t.Error("expected IsFocusTrigger false for click")
	}
}

// =============================================================================
// Content Helper Tests
// =============================================================================

func TestHasTitle(t *testing.T) {
	p := New("test")
	if p.HasTitle() {
		t.Error("expected HasTitle false when no title")
	}

	p.Title = "Header"
	if !p.HasTitle() {
		t.Error("expected HasTitle true when title set")
	}
}

func TestHasContent(t *testing.T) {
	p := New("test")
	if p.HasContent() {
		t.Error("expected HasContent false when no content")
	}

	p.Content = "Body"
	if !p.HasContent() {
		t.Error("expected HasContent true when content set")
	}
}

func TestHasHeader(t *testing.T) {
	p := New("test")
	if p.HasHeader() {
		t.Error("expected HasHeader false when no title or close")
	}

	p.Title = "Title"
	if !p.HasHeader() {
		t.Error("expected HasHeader true when title set")
	}

	p.Title = ""
	p.ShowClose = true
	if !p.HasHeader() {
		t.Error("expected HasHeader true when ShowClose set")
	}
}

// =============================================================================
// CSS Class Tests
// =============================================================================

func TestPositionClasses(t *testing.T) {
	tests := []struct {
		position Position
	}{
		{PositionTop},
		{PositionTopStart},
		{PositionTopEnd},
		{PositionBottom},
		{PositionBottomStart},
		{PositionBottomEnd},
		{PositionLeft},
		{PositionLeftStart},
		{PositionLeftEnd},
		{PositionRight},
		{PositionRightStart},
		{PositionRightEnd},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			p := New("test", WithPosition(tt.position))
			classes := p.PositionClasses()
			if classes == "" {
				t.Error("expected non-empty position classes")
			}
		})
	}
}

func TestArrowClasses(t *testing.T) {
	t.Run("no arrow when disabled", func(t *testing.T) {
		p := New("test", WithArrow(false))
		if p.ArrowClasses() != "" {
			t.Error("expected empty arrow classes when arrow disabled")
		}
	})

	t.Run("arrow classes for top", func(t *testing.T) {
		p := New("test", WithPosition(PositionTop), WithArrow(true))
		classes := p.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for top position")
		}
	})

	t.Run("arrow classes for bottom", func(t *testing.T) {
		p := New("test", WithPosition(PositionBottom), WithArrow(true))
		classes := p.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for bottom position")
		}
	})

	t.Run("arrow classes for left", func(t *testing.T) {
		p := New("test", WithPosition(PositionLeft), WithArrow(true))
		classes := p.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for left position")
		}
	})

	t.Run("arrow classes for right", func(t *testing.T) {
		p := New("test", WithPosition(PositionRight), WithArrow(true))
		classes := p.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for right position")
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
