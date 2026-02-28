package tooltip

import (
	"testing"
)

// =============================================================================
// Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates tooltip with defaults", func(t *testing.T) {
		tip := New("help")
		if tip.ID() != "help" {
			t.Errorf("expected ID 'help', got %q", tip.ID())
		}
		if tip.Namespace() != "tooltip" {
			t.Errorf("expected namespace 'tooltip', got %q", tip.Namespace())
		}
		if tip.Position != PositionTop {
			t.Errorf("expected position top, got %v", tip.Position)
		}
		if tip.Trigger != TriggerHover {
			t.Errorf("expected trigger hover, got %v", tip.Trigger)
		}
		if tip.Visible {
			t.Error("expected tooltip to be hidden by default")
		}
		if !tip.Arrow {
			t.Error("expected arrow true by default")
		}
		if tip.MaxWidth != "200px" {
			t.Errorf("expected max width '200px', got %q", tip.MaxWidth)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		tip := New("info",
			WithContent("Information"),
			WithPosition(PositionBottom),
			WithTrigger(TriggerClick),
			WithArrow(false),
		)
		if tip.Content != "Information" {
			t.Errorf("expected content 'Information', got %q", tip.Content)
		}
		if tip.Position != PositionBottom {
			t.Errorf("expected position bottom, got %v", tip.Position)
		}
		if tip.Trigger != TriggerClick {
			t.Errorf("expected trigger click, got %v", tip.Trigger)
		}
		if tip.Arrow {
			t.Error("expected arrow false")
		}
	})
}

// =============================================================================
// Option Tests
// =============================================================================

func TestWithContent(t *testing.T) {
	tip := New("test", WithContent("Help text"))
	if tip.Content != "Help text" {
		t.Errorf("expected content 'Help text', got %q", tip.Content)
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
			tip := New("test", WithPosition(pos))
			if tip.Position != pos {
				t.Errorf("expected position %v, got %v", pos, tip.Position)
			}
		})
	}
}

func TestWithTrigger(t *testing.T) {
	triggers := []Trigger{TriggerHover, TriggerFocus, TriggerClick}

	for _, trigger := range triggers {
		t.Run(string(trigger), func(t *testing.T) {
			tip := New("test", WithTrigger(trigger))
			if tip.Trigger != trigger {
				t.Errorf("expected trigger %v, got %v", trigger, tip.Trigger)
			}
		})
	}
}

func TestWithDelay(t *testing.T) {
	tip := New("test", WithDelay(300))
	if tip.Delay != 300 {
		t.Errorf("expected delay 300, got %d", tip.Delay)
	}
}

func TestWithHideDelay(t *testing.T) {
	tip := New("test", WithHideDelay(100))
	if tip.HideDelay != 100 {
		t.Errorf("expected hide delay 100, got %d", tip.HideDelay)
	}
}

func TestWithArrow(t *testing.T) {
	tip := New("test", WithArrow(false))
	if tip.Arrow {
		t.Error("expected arrow false")
	}
}

func TestWithMaxWidth(t *testing.T) {
	tip := New("test", WithMaxWidth("300px"))
	if tip.MaxWidth != "300px" {
		t.Errorf("expected max width '300px', got %q", tip.MaxWidth)
	}
}

func TestWithVisible(t *testing.T) {
	tip := New("test", WithVisible(true))
	if !tip.Visible {
		t.Error("expected visible true")
	}
}

func TestWithStyled(t *testing.T) {
	tip := New("test", WithStyled(true))
	if !tip.IsStyled() {
		t.Error("expected styled true")
	}

	tip2 := New("test", WithStyled(false))
	if tip2.IsStyled() {
		t.Error("expected styled false")
	}
}

// =============================================================================
// Method Tests
// =============================================================================

func TestShow(t *testing.T) {
	tip := New("test")
	tip.Show()
	if !tip.Visible {
		t.Error("expected visible after Show")
	}
}

func TestHide(t *testing.T) {
	tip := New("test", WithVisible(true))
	tip.Hide()
	if tip.Visible {
		t.Error("expected hidden after Hide")
	}
}

func TestToggle(t *testing.T) {
	tip := New("test")
	if tip.Visible {
		t.Error("expected initially hidden")
	}

	tip.Toggle()
	if !tip.Visible {
		t.Error("expected visible after toggle")
	}

	tip.Toggle()
	if tip.Visible {
		t.Error("expected hidden after second toggle")
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
			tip := New("test", WithPosition(tt.position))
			if tip.IsTop() != tt.expected {
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
			tip := New("test", WithPosition(tt.position))
			if tip.IsBottom() != tt.expected {
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
			tip := New("test", WithPosition(tt.position))
			if tip.IsLeft() != tt.expected {
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
			tip := New("test", WithPosition(tt.position))
			if tip.IsRight() != tt.expected {
				t.Errorf("expected IsRight %v for position %v", tt.expected, tt.position)
			}
		})
	}
}

// =============================================================================
// Trigger Helper Tests
// =============================================================================

func TestIsHoverTrigger(t *testing.T) {
	tip := New("test", WithTrigger(TriggerHover))
	if !tip.IsHoverTrigger() {
		t.Error("expected IsHoverTrigger true")
	}
	tip.Trigger = TriggerClick
	if tip.IsHoverTrigger() {
		t.Error("expected IsHoverTrigger false for click")
	}
}

func TestIsFocusTrigger(t *testing.T) {
	tip := New("test", WithTrigger(TriggerFocus))
	if !tip.IsFocusTrigger() {
		t.Error("expected IsFocusTrigger true")
	}
	tip.Trigger = TriggerHover
	if tip.IsFocusTrigger() {
		t.Error("expected IsFocusTrigger false for hover")
	}
}

func TestIsClickTrigger(t *testing.T) {
	tip := New("test", WithTrigger(TriggerClick))
	if !tip.IsClickTrigger() {
		t.Error("expected IsClickTrigger true")
	}
	tip.Trigger = TriggerHover
	if tip.IsClickTrigger() {
		t.Error("expected IsClickTrigger false for hover")
	}
}

// =============================================================================
// CSS Class Tests
// =============================================================================

func TestPositionClasses(t *testing.T) {
	tests := []struct {
		position Position
		contains string
	}{
		{PositionTop, "bottom-full"},
		{PositionTopStart, "bottom-full left-0"},
		{PositionTopEnd, "bottom-full right-0"},
		{PositionBottom, "top-full"},
		{PositionBottomStart, "top-full left-0"},
		{PositionBottomEnd, "top-full right-0"},
		{PositionLeft, "right-full"},
		{PositionLeftStart, "right-full top-0"},
		{PositionLeftEnd, "right-full bottom-0"},
		{PositionRight, "left-full"},
		{PositionRightStart, "left-full top-0"},
		{PositionRightEnd, "left-full bottom-0"},
	}

	for _, tt := range tests {
		t.Run(string(tt.position), func(t *testing.T) {
			tip := New("test", WithPosition(tt.position))
			classes := tip.PositionClasses()
			if classes == "" {
				t.Error("expected non-empty position classes")
			}
		})
	}
}

func TestArrowClasses(t *testing.T) {
	t.Run("no arrow when disabled", func(t *testing.T) {
		tip := New("test", WithArrow(false))
		if tip.ArrowClasses() != "" {
			t.Error("expected empty arrow classes when arrow disabled")
		}
	})

	t.Run("arrow classes for top", func(t *testing.T) {
		tip := New("test", WithPosition(PositionTop), WithArrow(true))
		classes := tip.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for top position")
		}
	})

	t.Run("arrow classes for bottom", func(t *testing.T) {
		tip := New("test", WithPosition(PositionBottom), WithArrow(true))
		classes := tip.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for bottom position")
		}
	})

	t.Run("arrow classes for left", func(t *testing.T) {
		tip := New("test", WithPosition(PositionLeft), WithArrow(true))
		classes := tip.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for left position")
		}
	})

	t.Run("arrow classes for right", func(t *testing.T) {
		tip := New("test", WithPosition(PositionRight), WithArrow(true))
		classes := tip.ArrowClasses()
		if classes == "" {
			t.Error("expected arrow classes for right position")
		}
	})
}

// =============================================================================
// HasContent Tests
// =============================================================================

func TestHasContent(t *testing.T) {
	tip := New("test")
	if tip.HasContent() {
		t.Error("expected HasContent false when no content")
	}

	tip.Content = "Help"
	if !tip.HasContent() {
		t.Error("expected HasContent true when content set")
	}

	tip.Content = ""
	if tip.HasContent() {
		t.Error("expected HasContent false when content empty")
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
