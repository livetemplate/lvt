package breadcrumbs

import (
	"testing"
)

// =============================================================================
// Breadcrumbs Constructor Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates breadcrumbs with defaults", func(t *testing.T) {
		bc := New("nav")
		if bc.ID() != "nav" {
			t.Errorf("expected ID 'nav', got %q", bc.ID())
		}
		if bc.Namespace() != "breadcrumbs" {
			t.Errorf("expected namespace 'breadcrumbs', got %q", bc.Namespace())
		}
		if bc.Separator != SeparatorChevron {
			t.Errorf("expected chevron separator, got %v", bc.Separator)
		}
		if bc.Size != SizeMd {
			t.Errorf("expected md size, got %v", bc.Size)
		}
		if bc.ShowHome {
			t.Error("expected ShowHome false by default")
		}
		if bc.HomeHref != "/" {
			t.Errorf("expected HomeHref '/', got %q", bc.HomeHref)
		}
		if bc.Collapsible {
			t.Error("expected Collapsible false by default")
		}
		if bc.MaxVisible != 3 {
			t.Errorf("expected MaxVisible 3, got %d", bc.MaxVisible)
		}
		if len(bc.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(bc.Items))
		}
	})

	t.Run("applies options", func(t *testing.T) {
		item1 := NewItem("1", WithItemLabel("Home"))
		item2 := NewItem("2", WithItemLabel("Products"))
		bc := New("trail",
			WithItems(item1, item2),
			WithSeparator(SeparatorSlash),
			WithSize(SizeLg),
			WithShowHome(true),
			WithHomeHref("/dashboard"),
			WithCollapsible(true),
			WithMaxVisible(4),
		)
		if len(bc.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(bc.Items))
		}
		if bc.Separator != SeparatorSlash {
			t.Errorf("expected slash separator, got %v", bc.Separator)
		}
		if bc.Size != SizeLg {
			t.Errorf("expected lg size, got %v", bc.Size)
		}
		if !bc.ShowHome {
			t.Error("expected ShowHome true")
		}
		if bc.HomeHref != "/dashboard" {
			t.Errorf("expected HomeHref '/dashboard', got %q", bc.HomeHref)
		}
		if !bc.Collapsible {
			t.Error("expected Collapsible true")
		}
		if bc.MaxVisible != 4 {
			t.Errorf("expected MaxVisible 4, got %d", bc.MaxVisible)
		}
	})
}

// =============================================================================
// Breadcrumbs Options Tests
// =============================================================================

func TestWithItems(t *testing.T) {
	item := NewItem("1")
	bc := New("test", WithItems(item))
	if len(bc.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(bc.Items))
	}
}

func TestWithSeparator(t *testing.T) {
	separators := []Separator{SeparatorSlash, SeparatorChevron, SeparatorArrow, SeparatorDot}
	for _, s := range separators {
		t.Run(string(s), func(t *testing.T) {
			bc := New("test", WithSeparator(s))
			if bc.Separator != s {
				t.Errorf("expected %v, got %v", s, bc.Separator)
			}
		})
	}
}

func TestWithSize(t *testing.T) {
	sizes := []Size{SizeSm, SizeMd, SizeLg}
	for _, s := range sizes {
		t.Run(string(s), func(t *testing.T) {
			bc := New("test", WithSize(s))
			if bc.Size != s {
				t.Errorf("expected %v, got %v", s, bc.Size)
			}
		})
	}
}

func TestWithShowHome(t *testing.T) {
	bc := New("test", WithShowHome(true))
	if !bc.ShowHome {
		t.Error("expected ShowHome true")
	}
}

func TestWithHomeHref(t *testing.T) {
	bc := New("test", WithHomeHref("/home"))
	if bc.HomeHref != "/home" {
		t.Errorf("expected '/home', got %q", bc.HomeHref)
	}
}

func TestWithCollapsible(t *testing.T) {
	bc := New("test", WithCollapsible(true))
	if !bc.Collapsible {
		t.Error("expected Collapsible true")
	}
}

func TestWithMaxVisible(t *testing.T) {
	bc := New("test", WithMaxVisible(5))
	if bc.MaxVisible != 5 {
		t.Errorf("expected 5, got %d", bc.MaxVisible)
	}
}

func TestWithStyled(t *testing.T) {
	bc := New("test", WithStyled(true))
	if !bc.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// Breadcrumbs Method Tests
// =============================================================================

func TestBreadcrumbsAddItem(t *testing.T) {
	bc := New("test")
	item := NewItem("1", WithItemLabel("First"))
	bc.AddItem(item)
	if len(bc.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(bc.Items))
	}
	if bc.Items[0].Label != "First" {
		t.Errorf("expected label 'First', got %q", bc.Items[0].Label)
	}
}

func TestBreadcrumbsHasItems(t *testing.T) {
	bc := New("test")
	if bc.HasItems() {
		t.Error("expected HasItems false")
	}
	bc.AddItem(NewItem("1"))
	if !bc.HasItems() {
		t.Error("expected HasItems true")
	}
}

func TestBreadcrumbsItemCount(t *testing.T) {
	bc := New("test")
	if bc.ItemCount() != 0 {
		t.Errorf("expected 0, got %d", bc.ItemCount())
	}
	bc.AddItem(NewItem("1"))
	bc.AddItem(NewItem("2"))
	if bc.ItemCount() != 2 {
		t.Errorf("expected 2, got %d", bc.ItemCount())
	}
}

func TestBreadcrumbsLastItem(t *testing.T) {
	bc := New("test")
	if bc.LastItem() != nil {
		t.Error("expected nil for empty breadcrumbs")
	}
	bc.AddItem(NewItem("1", WithItemLabel("First")))
	bc.AddItem(NewItem("2", WithItemLabel("Last")))
	last := bc.LastItem()
	if last == nil {
		t.Fatal("expected last item")
	}
	if last.Label != "Last" {
		t.Errorf("expected label 'Last', got %q", last.Label)
	}
}

func TestBreadcrumbsIsCollapsed(t *testing.T) {
	t.Run("not collapsible", func(t *testing.T) {
		bc := New("test",
			WithItems(NewItem("1"), NewItem("2"), NewItem("3"), NewItem("4")),
		)
		if bc.IsCollapsed() {
			t.Error("expected not collapsed when Collapsible is false")
		}
	})

	t.Run("collapsible but few items", func(t *testing.T) {
		bc := New("test",
			WithCollapsible(true),
			WithMaxVisible(3),
			WithItems(NewItem("1"), NewItem("2")),
		)
		if bc.IsCollapsed() {
			t.Error("expected not collapsed with fewer items than MaxVisible")
		}
	})

	t.Run("collapsible with many items", func(t *testing.T) {
		bc := New("test",
			WithCollapsible(true),
			WithMaxVisible(3),
			WithItems(NewItem("1"), NewItem("2"), NewItem("3"), NewItem("4")),
		)
		if !bc.IsCollapsed() {
			t.Error("expected collapsed with more items than MaxVisible")
		}
	})
}

func TestBreadcrumbsVisibleItems(t *testing.T) {
	t.Run("not collapsed", func(t *testing.T) {
		bc := New("test",
			WithItems(NewItem("1"), NewItem("2")),
		)
		visible := bc.VisibleItems()
		if len(visible) != 2 {
			t.Errorf("expected 2, got %d", len(visible))
		}
	})

	t.Run("collapsed", func(t *testing.T) {
		bc := New("test",
			WithCollapsible(true),
			WithMaxVisible(3),
			WithItems(
				NewItem("1", WithItemLabel("First")),
				NewItem("2", WithItemLabel("Second")),
				NewItem("3", WithItemLabel("Third")),
				NewItem("4", WithItemLabel("Fourth")),
				NewItem("5", WithItemLabel("Fifth")),
			),
		)
		visible := bc.VisibleItems()
		if len(visible) != 3 {
			t.Errorf("expected 3 visible items, got %d", len(visible))
		}
		if visible[0].Label != "First" {
			t.Errorf("expected first item 'First', got %q", visible[0].Label)
		}
		if visible[1].Label != "Fourth" {
			t.Errorf("expected second visible 'Fourth', got %q", visible[1].Label)
		}
		if visible[2].Label != "Fifth" {
			t.Errorf("expected third visible 'Fifth', got %q", visible[2].Label)
		}
	})
}

func TestBreadcrumbsHiddenCount(t *testing.T) {
	t.Run("not collapsed", func(t *testing.T) {
		bc := New("test", WithItems(NewItem("1"), NewItem("2")))
		if bc.HiddenCount() != 0 {
			t.Errorf("expected 0, got %d", bc.HiddenCount())
		}
	})

	t.Run("collapsed", func(t *testing.T) {
		bc := New("test",
			WithCollapsible(true),
			WithMaxVisible(3),
			WithItems(NewItem("1"), NewItem("2"), NewItem("3"), NewItem("4"), NewItem("5")),
		)
		if bc.HiddenCount() != 2 {
			t.Errorf("expected 2 hidden, got %d", bc.HiddenCount())
		}
	})
}

func TestBreadcrumbsSeparatorSymbol(t *testing.T) {
	tests := []struct {
		separator Separator
		expected  string
	}{
		{SeparatorSlash, "/"},
		{SeparatorChevron, ""},
		{SeparatorArrow, "→"},
		{SeparatorDot, "•"},
	}

	for _, tt := range tests {
		t.Run(string(tt.separator), func(t *testing.T) {
			bc := New("test", WithSeparator(tt.separator))
			if bc.SeparatorSymbol() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, bc.SeparatorSymbol())
			}
		})
	}
}

func TestBreadcrumbsIsChevronSeparator(t *testing.T) {
	bc := New("test", WithSeparator(SeparatorChevron))
	if !bc.IsChevronSeparator() {
		t.Error("expected IsChevronSeparator true")
	}
	bc.Separator = SeparatorSlash
	if bc.IsChevronSeparator() {
		t.Error("expected IsChevronSeparator false")
	}
}

func TestBreadcrumbsSizeClass(t *testing.T) {
	tests := []struct {
		size     Size
		expected string
	}{
		{SizeSm, "text-sm"},
		{SizeMd, "text-base"},
		{SizeLg, "text-lg"},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			bc := New("test", WithSize(tt.size))
			if bc.SizeClass() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, bc.SizeClass())
			}
		})
	}
}

// =============================================================================
// BreadcrumbItem Constructor Tests
// =============================================================================

func TestNewItem(t *testing.T) {
	t.Run("creates item with defaults", func(t *testing.T) {
		item := NewItem("home")
		if item.ID() != "home" {
			t.Errorf("expected ID 'home', got %q", item.ID())
		}
		if item.Namespace() != "breadcrumb-item" {
			t.Errorf("expected namespace 'breadcrumb-item', got %q", item.Namespace())
		}
		if item.Current {
			t.Error("expected Current false by default")
		}
		if item.Disabled {
			t.Error("expected Disabled false by default")
		}
	})

	t.Run("applies options", func(t *testing.T) {
		item := NewItem("product",
			WithItemLabel("Product"),
			WithItemHref("/product"),
			WithItemIcon("📦"),
			WithItemCurrent(true),
			WithItemDisabled(true),
		)
		if item.Label != "Product" {
			t.Errorf("expected label 'Product', got %q", item.Label)
		}
		if item.Href != "/product" {
			t.Errorf("expected href '/product', got %q", item.Href)
		}
		if item.Icon != "📦" {
			t.Errorf("expected icon '📦', got %q", item.Icon)
		}
		if !item.Current {
			t.Error("expected Current true")
		}
		if !item.Disabled {
			t.Error("expected Disabled true")
		}
	})
}

// =============================================================================
// BreadcrumbItem Options Tests
// =============================================================================

func TestWithItemLabel(t *testing.T) {
	item := NewItem("test", WithItemLabel("Label"))
	if item.Label != "Label" {
		t.Errorf("expected 'Label', got %q", item.Label)
	}
}

func TestWithItemHref(t *testing.T) {
	item := NewItem("test", WithItemHref("/path"))
	if item.Href != "/path" {
		t.Errorf("expected '/path', got %q", item.Href)
	}
}

func TestWithItemIcon(t *testing.T) {
	item := NewItem("test", WithItemIcon("star"))
	if item.Icon != "star" {
		t.Errorf("expected 'star', got %q", item.Icon)
	}
}

func TestWithItemCurrent(t *testing.T) {
	item := NewItem("test", WithItemCurrent(true))
	if !item.Current {
		t.Error("expected Current true")
	}
}

func TestWithItemDisabled(t *testing.T) {
	item := NewItem("test", WithItemDisabled(true))
	if !item.Disabled {
		t.Error("expected Disabled true")
	}
}

func TestWithItemStyled(t *testing.T) {
	item := NewItem("test", WithItemStyled(true))
	if !item.IsStyled() {
		t.Error("expected styled")
	}
}

// =============================================================================
// BreadcrumbItem Method Tests
// =============================================================================

func TestItemHasLabel(t *testing.T) {
	item := NewItem("test")
	if item.HasLabel() {
		t.Error("expected HasLabel false")
	}
	item.Label = "Label"
	if !item.HasLabel() {
		t.Error("expected HasLabel true")
	}
}

func TestItemHasHref(t *testing.T) {
	item := NewItem("test")
	if item.HasHref() {
		t.Error("expected HasHref false")
	}
	item.Href = "/path"
	if !item.HasHref() {
		t.Error("expected HasHref true")
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

func TestItemIsClickable(t *testing.T) {
	t.Run("clickable with href", func(t *testing.T) {
		item := NewItem("test", WithItemHref("/path"))
		if !item.IsClickable() {
			t.Error("expected IsClickable true")
		}
	})

	t.Run("not clickable without href", func(t *testing.T) {
		item := NewItem("test")
		if item.IsClickable() {
			t.Error("expected IsClickable false without href")
		}
	})

	t.Run("not clickable when current", func(t *testing.T) {
		item := NewItem("test", WithItemHref("/path"), WithItemCurrent(true))
		if item.IsClickable() {
			t.Error("expected IsClickable false when current")
		}
	})

	t.Run("not clickable when disabled", func(t *testing.T) {
		item := NewItem("test", WithItemHref("/path"), WithItemDisabled(true))
		if item.IsClickable() {
			t.Error("expected IsClickable false when disabled")
		}
	})
}

func TestItemLinkClass(t *testing.T) {
	t.Run("current", func(t *testing.T) {
		item := NewItem("test", WithItemCurrent(true))
		if item.LinkClass() != "text-gray-700 font-medium" {
			t.Errorf("expected current class, got %q", item.LinkClass())
		}
	})

	t.Run("disabled", func(t *testing.T) {
		item := NewItem("test", WithItemDisabled(true))
		if item.LinkClass() != "text-gray-400 cursor-not-allowed" {
			t.Errorf("expected disabled class, got %q", item.LinkClass())
		}
	})

	t.Run("default", func(t *testing.T) {
		item := NewItem("test")
		if item.LinkClass() != "text-gray-500 hover:text-gray-700" {
			t.Errorf("expected default class, got %q", item.LinkClass())
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
