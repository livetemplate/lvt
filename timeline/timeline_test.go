package timeline

import (
	"testing"
)

// =============================================================================
// Timeline Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates timeline with defaults", func(t *testing.T) {
		tl := New("history")
		if tl.ID() != "history" {
			t.Errorf("expected ID 'history', got %q", tl.ID())
		}
		if tl.Namespace() != "timeline" {
			t.Errorf("expected namespace 'timeline', got %q", tl.Namespace())
		}
		if tl.Orientation != OrientationVertical {
			t.Errorf("expected vertical orientation, got %v", tl.Orientation)
		}
		if tl.Position != PositionLeft {
			t.Errorf("expected left position, got %v", tl.Position)
		}
		if !tl.ShowConnectors {
			t.Error("expected ShowConnectors true by default")
		}
		if tl.Reverse {
			t.Error("expected Reverse false by default")
		}
		if len(tl.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(tl.Items))
		}
	})

	t.Run("applies options", func(t *testing.T) {
		item1 := NewItem("1", WithItemTitle("First"))
		item2 := NewItem("2", WithItemTitle("Second"))
		tl := New("events",
			WithItems(item1, item2),
			WithOrientation(OrientationHorizontal),
			WithPosition(PositionAlternate),
			WithShowConnectors(false),
			WithReverse(true),
		)
		if len(tl.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(tl.Items))
		}
		if tl.Orientation != OrientationHorizontal {
			t.Errorf("expected horizontal, got %v", tl.Orientation)
		}
		if tl.Position != PositionAlternate {
			t.Errorf("expected alternate, got %v", tl.Position)
		}
		if tl.ShowConnectors {
			t.Error("expected ShowConnectors false")
		}
		if !tl.Reverse {
			t.Error("expected Reverse true")
		}
	})
}

// =============================================================================
// Timeline Options Tests
// =============================================================================

func TestWithItems(t *testing.T) {
	item := NewItem("1")
	tl := New("test", WithItems(item))
	if len(tl.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(tl.Items))
	}
}

func TestWithOrientation(t *testing.T) {
	orientations := []Orientation{OrientationVertical, OrientationHorizontal}
	for _, o := range orientations {
		t.Run(string(o), func(t *testing.T) {
			tl := New("test", WithOrientation(o))
			if tl.Orientation != o {
				t.Errorf("expected %v, got %v", o, tl.Orientation)
			}
		})
	}
}

func TestWithPosition(t *testing.T) {
	positions := []Position{PositionLeft, PositionRight, PositionAlternate}
	for _, p := range positions {
		t.Run(string(p), func(t *testing.T) {
			tl := New("test", WithPosition(p))
			if tl.Position != p {
				t.Errorf("expected %v, got %v", p, tl.Position)
			}
		})
	}
}

func TestWithShowConnectors(t *testing.T) {
	tl := New("test", WithShowConnectors(false))
	if tl.ShowConnectors {
		t.Error("expected ShowConnectors false")
	}
}

func TestWithReverse(t *testing.T) {
	tl := New("test", WithReverse(true))
	if !tl.Reverse {
		t.Error("expected Reverse true")
	}
}

func TestWithStyled(t *testing.T) {
	tl := New("test", WithStyled(true))
	if !tl.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// Timeline Method Tests
// =============================================================================

func TestTimelineAddItem(t *testing.T) {
	tl := New("test")
	item := NewItem("1", WithItemTitle("First"))
	tl.AddItem(item)
	if len(tl.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(tl.Items))
	}
	if tl.Items[0].Title != "First" {
		t.Errorf("expected title 'First', got %q", tl.Items[0].Title)
	}
}

func TestTimelineRemoveItem(t *testing.T) {
	item1 := NewItem("1")
	item2 := NewItem("2")
	tl := New("test", WithItems(item1, item2))
	tl.RemoveItem("1")
	if len(tl.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(tl.Items))
	}
	if tl.Items[0].ID() != "2" {
		t.Errorf("expected item '2', got %q", tl.Items[0].ID())
	}
}

func TestTimelineGetItem(t *testing.T) {
	item := NewItem("1", WithItemTitle("Found"))
	tl := New("test", WithItems(item))

	found := tl.GetItem("1")
	if found == nil {
		t.Fatal("expected to find item")
	}
	if found.Title != "Found" {
		t.Errorf("expected title 'Found', got %q", found.Title)
	}

	notFound := tl.GetItem("999")
	if notFound != nil {
		t.Error("expected nil for non-existent item")
	}
}

func TestTimelineHasItems(t *testing.T) {
	tl := New("test")
	if tl.HasItems() {
		t.Error("expected HasItems false")
	}
	tl.AddItem(NewItem("1"))
	if !tl.HasItems() {
		t.Error("expected HasItems true")
	}
}

func TestTimelineItemCount(t *testing.T) {
	tl := New("test")
	if tl.ItemCount() != 0 {
		t.Errorf("expected 0, got %d", tl.ItemCount())
	}
	tl.AddItem(NewItem("1"))
	tl.AddItem(NewItem("2"))
	if tl.ItemCount() != 2 {
		t.Errorf("expected 2, got %d", tl.ItemCount())
	}
}

func TestTimelineIsVertical(t *testing.T) {
	tl := New("test")
	if !tl.IsVertical() {
		t.Error("expected IsVertical true by default")
	}
	tl.Orientation = OrientationHorizontal
	if tl.IsVertical() {
		t.Error("expected IsVertical false")
	}
}

func TestTimelineIsHorizontal(t *testing.T) {
	tl := New("test", WithOrientation(OrientationHorizontal))
	if !tl.IsHorizontal() {
		t.Error("expected IsHorizontal true")
	}
}

func TestTimelineIsAlternate(t *testing.T) {
	tl := New("test", WithPosition(PositionAlternate))
	if !tl.IsAlternate() {
		t.Error("expected IsAlternate true")
	}
}

func TestTimelineOrientationClass(t *testing.T) {
	t.Run("vertical", func(t *testing.T) {
		tl := New("test")
		if tl.OrientationClass() != "flex flex-col" {
			t.Errorf("expected 'flex flex-col', got %q", tl.OrientationClass())
		}
	})

	t.Run("horizontal", func(t *testing.T) {
		tl := New("test", WithOrientation(OrientationHorizontal))
		if tl.OrientationClass() != "flex flex-row" {
			t.Errorf("expected 'flex flex-row', got %q", tl.OrientationClass())
		}
	})
}

// =============================================================================
// TimelineItem Constructor Tests
// =============================================================================

func TestNewItem(t *testing.T) {
	t.Run("creates item with defaults", func(t *testing.T) {
		item := NewItem("step1")
		if item.ID() != "step1" {
			t.Errorf("expected ID 'step1', got %q", item.ID())
		}
		if item.Namespace() != "timeline-item" {
			t.Errorf("expected namespace 'timeline-item', got %q", item.Namespace())
		}
		if item.Status != StatusDefault {
			t.Errorf("expected default status, got %v", item.Status)
		}
		if item.Color != ColorGray {
			t.Errorf("expected gray color, got %v", item.Color)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		item := NewItem("event",
			WithItemTitle("Meeting"),
			WithItemDescription("Team sync"),
			WithItemTime("10:00 AM"),
			WithItemIcon("📅"),
			WithItemStatus(StatusComplete),
			WithItemColor(ColorGreen),
			WithItemActive(true),
			WithItemCompleted(true),
		)
		if item.Title != "Meeting" {
			t.Errorf("expected title 'Meeting', got %q", item.Title)
		}
		if item.Description != "Team sync" {
			t.Errorf("expected description 'Team sync', got %q", item.Description)
		}
		if item.Time != "10:00 AM" {
			t.Errorf("expected time '10:00 AM', got %q", item.Time)
		}
		if item.Icon != "📅" {
			t.Errorf("expected icon '📅', got %q", item.Icon)
		}
		if item.Status != StatusComplete {
			t.Errorf("expected complete status, got %v", item.Status)
		}
		if item.Color != ColorGreen {
			t.Errorf("expected green color, got %v", item.Color)
		}
		if !item.Active {
			t.Error("expected Active true")
		}
		if !item.Completed {
			t.Error("expected Completed true")
		}
	})
}

// =============================================================================
// TimelineItem Options Tests
// =============================================================================

func TestWithItemTitle(t *testing.T) {
	item := NewItem("test", WithItemTitle("Title"))
	if item.Title != "Title" {
		t.Errorf("expected 'Title', got %q", item.Title)
	}
}

func TestWithItemDescription(t *testing.T) {
	item := NewItem("test", WithItemDescription("Desc"))
	if item.Description != "Desc" {
		t.Errorf("expected 'Desc', got %q", item.Description)
	}
}

func TestWithItemTime(t *testing.T) {
	item := NewItem("test", WithItemTime("2024-01-01"))
	if item.Time != "2024-01-01" {
		t.Errorf("expected '2024-01-01', got %q", item.Time)
	}
}

func TestWithItemIcon(t *testing.T) {
	item := NewItem("test", WithItemIcon("star"))
	if item.Icon != "star" {
		t.Errorf("expected 'star', got %q", item.Icon)
	}
}

func TestWithItemStatus(t *testing.T) {
	statuses := []Status{StatusDefault, StatusPending, StatusActive, StatusComplete, StatusError}
	for _, s := range statuses {
		t.Run(string(s), func(t *testing.T) {
			item := NewItem("test", WithItemStatus(s))
			if item.Status != s {
				t.Errorf("expected %v, got %v", s, item.Status)
			}
		})
	}
}

func TestWithItemColor(t *testing.T) {
	colors := []Color{ColorGray, ColorBlue, ColorGreen, ColorYellow, ColorRed, ColorPurple}
	for _, c := range colors {
		t.Run(string(c), func(t *testing.T) {
			item := NewItem("test", WithItemColor(c))
			if item.Color != c {
				t.Errorf("expected %v, got %v", c, item.Color)
			}
		})
	}
}

func TestWithItemActive(t *testing.T) {
	item := NewItem("test", WithItemActive(true))
	if !item.Active {
		t.Error("expected Active true")
	}
}

func TestWithItemCompleted(t *testing.T) {
	item := NewItem("test", WithItemCompleted(true))
	if !item.Completed {
		t.Error("expected Completed true")
	}
}

func TestWithItemStyled(t *testing.T) {
	item := NewItem("test", WithItemStyled(true))
	if !item.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// TimelineItem Method Tests
// =============================================================================

func TestItemHasTitle(t *testing.T) {
	item := NewItem("test")
	if item.HasTitle() {
		t.Error("expected HasTitle false")
	}
	item.Title = "Title"
	if !item.HasTitle() {
		t.Error("expected HasTitle true")
	}
}

func TestItemHasDescription(t *testing.T) {
	item := NewItem("test")
	if item.HasDescription() {
		t.Error("expected HasDescription false")
	}
	item.Description = "Desc"
	if !item.HasDescription() {
		t.Error("expected HasDescription true")
	}
}

func TestItemHasTime(t *testing.T) {
	item := NewItem("test")
	if item.HasTime() {
		t.Error("expected HasTime false")
	}
	item.Time = "Now"
	if !item.HasTime() {
		t.Error("expected HasTime true")
	}
}

func TestItemHasIcon(t *testing.T) {
	item := NewItem("test")
	if item.HasIcon() {
		t.Error("expected HasIcon false")
	}
	item.Icon = "star"
	if !item.HasIcon() {
		t.Error("expected HasIcon true")
	}
}

func TestItemIsPending(t *testing.T) {
	item := NewItem("test", WithItemStatus(StatusPending))
	if !item.IsPending() {
		t.Error("expected IsPending true")
	}
}

func TestItemIsActive(t *testing.T) {
	t.Run("via status", func(t *testing.T) {
		item := NewItem("test", WithItemStatus(StatusActive))
		if !item.IsActive() {
			t.Error("expected IsActive true")
		}
	})

	t.Run("via Active field", func(t *testing.T) {
		item := NewItem("test", WithItemActive(true))
		if !item.IsActive() {
			t.Error("expected IsActive true")
		}
	})
}

func TestItemIsComplete(t *testing.T) {
	t.Run("via status", func(t *testing.T) {
		item := NewItem("test", WithItemStatus(StatusComplete))
		if !item.IsComplete() {
			t.Error("expected IsComplete true")
		}
	})

	t.Run("via Completed field", func(t *testing.T) {
		item := NewItem("test", WithItemCompleted(true))
		if !item.IsComplete() {
			t.Error("expected IsComplete true")
		}
	})
}

func TestItemIsError(t *testing.T) {
	item := NewItem("test", WithItemStatus(StatusError))
	if !item.IsError() {
		t.Error("expected IsError true")
	}
}

func TestItemIndicatorClass(t *testing.T) {
	tests := []struct {
		color    Color
		expected string
	}{
		{ColorGray, "bg-gray-400"},
		{ColorBlue, "bg-blue-500"},
		{ColorGreen, "bg-green-500"},
		{ColorYellow, "bg-yellow-500"},
		{ColorRed, "bg-red-500"},
		{ColorPurple, "bg-purple-500"},
	}

	for _, tt := range tests {
		t.Run(string(tt.color), func(t *testing.T) {
			item := NewItem("test", WithItemColor(tt.color))
			if item.IndicatorClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, item.IndicatorClass())
			}
		})
	}
}

func TestItemStatusClass(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusDefault, "bg-gray-400 text-white"},
		{StatusPending, "bg-gray-200 text-gray-500"},
		{StatusActive, "bg-blue-500 text-white ring-4 ring-blue-100"},
		{StatusComplete, "bg-green-500 text-white"},
		{StatusError, "bg-red-500 text-white"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			item := NewItem("test", WithItemStatus(tt.status))
			if item.StatusClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, item.StatusClass())
			}
		})
	}
}

func TestItemRingClass(t *testing.T) {
	t.Run("not active", func(t *testing.T) {
		item := NewItem("test")
		if item.RingClass() != "" {
			t.Errorf("expected empty, got %q", item.RingClass())
		}
	})

	t.Run("active with colors", func(t *testing.T) {
		tests := []struct {
			color    Color
			expected string
		}{
			{ColorGray, "ring-4 ring-gray-100"},
			{ColorBlue, "ring-4 ring-blue-100"},
			{ColorGreen, "ring-4 ring-green-100"},
			{ColorYellow, "ring-4 ring-yellow-100"},
			{ColorRed, "ring-4 ring-red-100"},
			{ColorPurple, "ring-4 ring-purple-100"},
		}

		for _, tt := range tests {
			t.Run(string(tt.color), func(t *testing.T) {
				item := NewItem("test", WithItemActive(true), WithItemColor(tt.color))
				if item.RingClass() != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, item.RingClass())
				}
			})
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
