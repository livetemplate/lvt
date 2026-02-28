package modal

import (
	"testing"
)

// =============================================================================
// Modal Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates modal with defaults", func(t *testing.T) {
		m := New("settings")
		if m.ID() != "settings" {
			t.Errorf("expected ID 'settings', got %q", m.ID())
		}
		if m.Namespace() != "modal" {
			t.Errorf("expected namespace 'modal', got %q", m.Namespace())
		}
		if m.Open {
			t.Error("expected closed by default")
		}
		if m.Size != SizeMd {
			t.Errorf("expected size md, got %v", m.Size)
		}
		if !m.ShowClose {
			t.Error("expected ShowClose true by default")
		}
		if !m.CloseOnOverlay {
			t.Error("expected CloseOnOverlay true by default")
		}
		if !m.CloseOnEscape {
			t.Error("expected CloseOnEscape true by default")
		}
		if !m.Centered {
			t.Error("expected Centered true by default")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		m := New("dialog",
			WithTitle("My Dialog"),
			WithSize(SizeLg),
			WithOpen(true),
			WithShowClose(false),
		)
		if m.Title != "My Dialog" {
			t.Errorf("expected title 'My Dialog', got %q", m.Title)
		}
		if m.Size != SizeLg {
			t.Errorf("expected size lg, got %v", m.Size)
		}
		if !m.Open {
			t.Error("expected open")
		}
		if m.ShowClose {
			t.Error("expected ShowClose false")
		}
	})
}

// =============================================================================
// Modal Option Tests
// =============================================================================

func TestWithOpen(t *testing.T) {
	m := New("test", WithOpen(true))
	if !m.Open {
		t.Error("expected open")
	}
}

func TestWithTitle(t *testing.T) {
	m := New("test", WithTitle("Title"))
	if m.Title != "Title" {
		t.Errorf("expected title 'Title', got %q", m.Title)
	}
}

func TestWithSize(t *testing.T) {
	sizes := []Size{SizeSm, SizeMd, SizeLg, SizeXl, SizeFull}
	for _, size := range sizes {
		t.Run(string(size), func(t *testing.T) {
			m := New("test", WithSize(size))
			if m.Size != size {
				t.Errorf("expected size %v, got %v", size, m.Size)
			}
		})
	}
}

func TestWithShowClose(t *testing.T) {
	m := New("test", WithShowClose(false))
	if m.ShowClose {
		t.Error("expected ShowClose false")
	}
}

func TestWithCloseOnOverlay(t *testing.T) {
	m := New("test", WithCloseOnOverlay(false))
	if m.CloseOnOverlay {
		t.Error("expected CloseOnOverlay false")
	}
}

func TestWithCloseOnEscape(t *testing.T) {
	m := New("test", WithCloseOnEscape(false))
	if m.CloseOnEscape {
		t.Error("expected CloseOnEscape false")
	}
}

func TestWithCentered(t *testing.T) {
	m := New("test", WithCentered(false))
	if m.Centered {
		t.Error("expected Centered false")
	}
}

func TestWithScrollable(t *testing.T) {
	m := New("test", WithScrollable(true))
	if !m.Scrollable {
		t.Error("expected Scrollable true")
	}
}

func TestWithStyled(t *testing.T) {
	m := New("test", WithStyled(true))
	if !m.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// Modal Method Tests
// =============================================================================

func TestModalShow(t *testing.T) {
	m := New("test")
	m.Show()
	if !m.Open {
		t.Error("expected open after Show")
	}
}

func TestModalHide(t *testing.T) {
	m := New("test", WithOpen(true))
	m.Hide()
	if m.Open {
		t.Error("expected closed after Hide")
	}
}

func TestModalToggle(t *testing.T) {
	m := New("test")
	m.Toggle()
	if !m.Open {
		t.Error("expected open after toggle")
	}
	m.Toggle()
	if m.Open {
		t.Error("expected closed after second toggle")
	}
}

func TestModalHasTitle(t *testing.T) {
	m := New("test")
	if m.HasTitle() {
		t.Error("expected HasTitle false")
	}
	m.Title = "Title"
	if !m.HasTitle() {
		t.Error("expected HasTitle true")
	}
}

func TestModalHasHeader(t *testing.T) {
	m := New("test", WithShowClose(false))
	if m.HasHeader() {
		t.Error("expected HasHeader false")
	}
	m.Title = "Title"
	if !m.HasHeader() {
		t.Error("expected HasHeader true with title")
	}
	m.Title = ""
	m.ShowClose = true
	if !m.HasHeader() {
		t.Error("expected HasHeader true with close button")
	}
}

func TestModalSizeClass(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "max-w-sm"},
		{SizeMd, "max-w-lg"},
		{SizeLg, "max-w-2xl"},
		{SizeXl, "max-w-4xl"},
		{SizeFull, "max-w-full mx-4"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			m := New("test", WithSize(tt.size))
			if m.SizeClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, m.SizeClass())
			}
		})
	}
}

// =============================================================================
// ConfirmModal Tests
// =============================================================================

func TestNewConfirm(t *testing.T) {
	t.Run("creates confirm with defaults", func(t *testing.T) {
		c := NewConfirm("delete")
		if c.ID() != "delete" {
			t.Errorf("expected ID 'delete', got %q", c.ID())
		}
		if c.ConfirmText != "Confirm" {
			t.Errorf("expected ConfirmText 'Confirm', got %q", c.ConfirmText)
		}
		if c.CancelText != "Cancel" {
			t.Errorf("expected CancelText 'Cancel', got %q", c.CancelText)
		}
		if c.Icon != "warning" {
			t.Errorf("expected Icon 'warning', got %q", c.Icon)
		}
	})

	t.Run("applies options", func(t *testing.T) {
		c := NewConfirm("remove",
			WithConfirmTitle("Remove Item"),
			WithConfirmMessage("Are you sure?"),
			WithConfirmDestructive(true),
		)
		if c.Title != "Remove Item" {
			t.Errorf("expected title 'Remove Item', got %q", c.Title)
		}
		if c.Message != "Are you sure?" {
			t.Errorf("expected message 'Are you sure?', got %q", c.Message)
		}
		if !c.Destructive {
			t.Error("expected destructive")
		}
	})
}

func TestConfirmModalShow(t *testing.T) {
	c := NewConfirm("test")
	c.Show()
	if !c.Open {
		t.Error("expected open")
	}
}

func TestConfirmModalHide(t *testing.T) {
	c := NewConfirm("test", WithConfirmOpen(true))
	c.Hide()
	if c.Open {
		t.Error("expected closed")
	}
}

func TestConfirmModalHasTitle(t *testing.T) {
	c := NewConfirm("test")
	if c.HasTitle() {
		t.Error("expected HasTitle false")
	}
	c.Title = "Title"
	if !c.HasTitle() {
		t.Error("expected HasTitle true")
	}
}

func TestConfirmModalHasMessage(t *testing.T) {
	c := NewConfirm("test")
	if c.HasMessage() {
		t.Error("expected HasMessage false")
	}
	c.Message = "Message"
	if !c.HasMessage() {
		t.Error("expected HasMessage true")
	}
}

func TestConfirmButtonClass(t *testing.T) {
	c := NewConfirm("test")
	if c.ConfirmButtonClass() != "bg-blue-600 hover:bg-blue-700 text-white" {
		t.Errorf("expected blue button class, got %q", c.ConfirmButtonClass())
	}

	c.Destructive = true
	if c.ConfirmButtonClass() != "bg-red-600 hover:bg-red-700 text-white" {
		t.Errorf("expected red button class, got %q", c.ConfirmButtonClass())
	}
}

func TestConfirmIconClass(t *testing.T) {
	tests := []struct {
		icon        string
		destructive bool
		expected    string
	}{
		{"warning", false, "text-yellow-500"},
		{"warning", true, "text-red-500"},
		{"info", false, "text-blue-500"},
		{"success", false, "text-green-500"},
		{"error", false, "text-red-500"},
	}

	for _, tt := range tests {
		t.Run(tt.icon, func(t *testing.T) {
			c := NewConfirm("test", WithConfirmIcon(tt.icon), WithConfirmDestructive(tt.destructive))
			if c.IconClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, c.IconClass())
			}
		})
	}
}

// =============================================================================
// SheetModal Tests
// =============================================================================

func TestNewSheet(t *testing.T) {
	t.Run("creates sheet with defaults", func(t *testing.T) {
		s := NewSheet("filters")
		if s.ID() != "filters" {
			t.Errorf("expected ID 'filters', got %q", s.ID())
		}
		if s.Position != "right" {
			t.Errorf("expected position 'right', got %q", s.Position)
		}
		if s.Size != SizeMd {
			t.Errorf("expected size md, got %v", s.Size)
		}
		if !s.ShowClose {
			t.Error("expected ShowClose true")
		}
		if !s.CloseOnOverlay {
			t.Error("expected CloseOnOverlay true")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		s := NewSheet("menu",
			WithSheetTitle("Menu"),
			WithSheetPosition("left"),
			WithSheetSize(SizeLg),
		)
		if s.Title != "Menu" {
			t.Errorf("expected title 'Menu', got %q", s.Title)
		}
		if s.Position != "left" {
			t.Errorf("expected position 'left', got %q", s.Position)
		}
		if s.Size != SizeLg {
			t.Errorf("expected size lg, got %v", s.Size)
		}
	})
}

func TestSheetModalShow(t *testing.T) {
	s := NewSheet("test")
	s.Show()
	if !s.Open {
		t.Error("expected open")
	}
}

func TestSheetModalHide(t *testing.T) {
	s := NewSheet("test", WithSheetOpen(true))
	s.Hide()
	if s.Open {
		t.Error("expected closed")
	}
}

func TestSheetModalToggle(t *testing.T) {
	s := NewSheet("test")
	s.Toggle()
	if !s.Open {
		t.Error("expected open after toggle")
	}
	s.Toggle()
	if s.Open {
		t.Error("expected closed after second toggle")
	}
}

func TestSheetIsLeft(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("left"))
	if !s.IsLeft() {
		t.Error("expected IsLeft true")
	}
}

func TestSheetIsRight(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("right"))
	if !s.IsRight() {
		t.Error("expected IsRight true")
	}
}

func TestSheetIsTop(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("top"))
	if !s.IsTop() {
		t.Error("expected IsTop true")
	}
}

func TestSheetIsBottom(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("bottom"))
	if !s.IsBottom() {
		t.Error("expected IsBottom true")
	}
}

func TestSheetIsHorizontal(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("left"))
	if !s.IsHorizontal() {
		t.Error("expected IsHorizontal true for left")
	}
	s.Position = "right"
	if !s.IsHorizontal() {
		t.Error("expected IsHorizontal true for right")
	}
	s.Position = "top"
	if s.IsHorizontal() {
		t.Error("expected IsHorizontal false for top")
	}
}

func TestSheetIsVertical(t *testing.T) {
	s := NewSheet("test", WithSheetPosition("top"))
	if !s.IsVertical() {
		t.Error("expected IsVertical true for top")
	}
	s.Position = "bottom"
	if !s.IsVertical() {
		t.Error("expected IsVertical true for bottom")
	}
	s.Position = "left"
	if s.IsVertical() {
		t.Error("expected IsVertical false for left")
	}
}

func TestSheetPositionClass(t *testing.T) {
	tests := []struct {
		position string
		expected string
	}{
		{"left", "left-0 top-0 h-full"},
		{"right", "right-0 top-0 h-full"},
		{"top", "top-0 left-0 w-full"},
		{"bottom", "bottom-0 left-0 w-full"},
	}

	for _, tt := range tests {
		t.Run(tt.position, func(t *testing.T) {
			s := NewSheet("test", WithSheetPosition(tt.position))
			if s.PositionClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, s.PositionClass())
			}
		})
	}
}

func TestSheetSizeClass(t *testing.T) {
	t.Run("horizontal sizes", func(t *testing.T) {
		tests := []struct {
			size     Size
			expected string
		}{
			{SizeSm, "w-64"},
			{SizeMd, "w-80"},
			{SizeLg, "w-96"},
		}
		for _, tt := range tests {
			s := NewSheet("test", WithSheetPosition("right"), WithSheetSize(tt.size))
			if s.SizeClass() != tt.expected {
				t.Errorf("size %v: expected %q, got %q", tt.size, tt.expected, s.SizeClass())
			}
		}
	})

	t.Run("vertical sizes", func(t *testing.T) {
		tests := []struct {
			size     Size
			expected string
		}{
			{SizeSm, "h-48"},
			{SizeMd, "h-64"},
			{SizeLg, "h-96"},
		}
		for _, tt := range tests {
			s := NewSheet("test", WithSheetPosition("top"), WithSheetSize(tt.size))
			if s.SizeClass() != tt.expected {
				t.Errorf("size %v: expected %q, got %q", tt.size, tt.expected, s.SizeClass())
			}
		}
	})
}

func TestSheetTransformClass(t *testing.T) {
	t.Run("open", func(t *testing.T) {
		s := NewSheet("test", WithSheetOpen(true))
		if s.TransformClass() != "translate-x-0 translate-y-0" {
			t.Errorf("expected no transform when open, got %q", s.TransformClass())
		}
	})

	t.Run("closed positions", func(t *testing.T) {
		tests := []struct {
			position string
			expected string
		}{
			{"left", "-translate-x-full"},
			{"right", "translate-x-full"},
			{"top", "-translate-y-full"},
			{"bottom", "translate-y-full"},
		}
		for _, tt := range tests {
			s := NewSheet("test", WithSheetPosition(tt.position))
			if s.TransformClass() != tt.expected {
				t.Errorf("position %s: expected %q, got %q", tt.position, tt.expected, s.TransformClass())
			}
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
