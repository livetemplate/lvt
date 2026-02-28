package menu

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New("test-menu")

	if m.ID() != "test-menu" {
		t.Errorf("Expected ID 'test-menu', got '%s'", m.ID())
	}
	if m.Namespace() != "menu" {
		t.Errorf("Expected namespace 'menu', got '%s'", m.Namespace())
	}
	if m.Position != "bottom-left" {
		t.Errorf("Expected default position 'bottom-left', got '%s'", m.Position)
	}
	if m.HighlightedIndex != -1 {
		t.Errorf("Expected HighlightedIndex -1, got %d", m.HighlightedIndex)
	}
	if m.Open {
		t.Error("Expected Open to be false")
	}
}

func TestNewContext(t *testing.T) {
	cm := NewContext("ctx-menu")

	if cm.ID() != "ctx-menu" {
		t.Errorf("Expected ID 'ctx-menu', got '%s'", cm.ID())
	}
	if cm.X != 0 || cm.Y != 0 {
		t.Error("Expected X and Y to be 0")
	}
}

func TestNewNav(t *testing.T) {
	nm := NewNav("nav-menu")

	if nm.ID() != "nav-menu" {
		t.Errorf("Expected ID 'nav-menu', got '%s'", nm.ID())
	}
	if nm.Orientation != "horizontal" {
		t.Errorf("Expected default orientation 'horizontal', got '%s'", nm.Orientation)
	}
}

func TestWithItems(t *testing.T) {
	items := []Item{
		{ID: "1", Label: "One"},
		{ID: "2", Label: "Two"},
	}
	m := New("test", WithItems(items))

	if len(m.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(m.Items))
	}
}

func TestWithTrigger(t *testing.T) {
	m := New("test", WithTrigger("Actions"))
	if m.Trigger != "Actions" {
		t.Errorf("Expected Trigger 'Actions', got '%s'", m.Trigger)
	}
}

func TestWithTriggerIcon(t *testing.T) {
	m := New("test", WithTriggerIcon("⚙️"))
	if m.TriggerIcon != "⚙️" {
		t.Errorf("Expected TriggerIcon '⚙️', got '%s'", m.TriggerIcon)
	}
}

func TestWithPosition(t *testing.T) {
	m := New("test", WithPosition("bottom-right"))
	if m.Position != "bottom-right" {
		t.Errorf("Expected Position 'bottom-right', got '%s'", m.Position)
	}
}

func TestWithOpen(t *testing.T) {
	m := New("test", WithOpen(true))
	if !m.Open {
		t.Error("Expected Open to be true")
	}
}

func TestWithStyled(t *testing.T) {
	m := New("test", WithStyled(false))
	if m.IsStyled() {
		t.Error("Expected IsStyled to be false")
	}
}

func TestToggle(t *testing.T) {
	m := New("test")

	m.Toggle()
	if !m.Open {
		t.Error("Expected Open to be true after toggle")
	}

	m.HighlightedIndex = 2
	m.Toggle()
	if m.Open {
		t.Error("Expected Open to be false after second toggle")
	}
	if m.HighlightedIndex != -1 {
		t.Error("Expected HighlightedIndex to reset to -1 on close")
	}
}

func TestShow(t *testing.T) {
	m := New("test")
	m.Show()
	if !m.Open {
		t.Error("Expected Open to be true after Show")
	}
}

func TestClose(t *testing.T) {
	m := New("test", WithOpen(true))
	m.HighlightedIndex = 2
	m.Close()

	if m.Open {
		t.Error("Expected Open to be false after Close")
	}
	if m.HighlightedIndex != -1 {
		t.Error("Expected HighlightedIndex to reset to -1")
	}
}

func TestSelectIndex(t *testing.T) {
	items := []Item{
		{ID: "edit", Label: "Edit"},
		{ID: "delete", Label: "Delete"},
	}
	m := New("test", WithItems(items), WithOpen(true))

	id := m.SelectIndex(0)

	if id != "edit" {
		t.Errorf("Expected 'edit', got '%s'", id)
	}
	if m.Open {
		t.Error("Expected Open to be false after selection")
	}
}

func TestSelectIndexOutOfBounds(t *testing.T) {
	m := New("test")

	if m.SelectIndex(0) != "" {
		t.Error("Expected empty string for out of bounds")
	}
	if m.SelectIndex(-1) != "" {
		t.Error("Expected empty string for negative index")
	}
}

func TestSelectIndexDisabled(t *testing.T) {
	items := []Item{
		{ID: "edit", Label: "Edit", Disabled: true},
	}
	m := New("test", WithItems(items))

	if m.SelectIndex(0) != "" {
		t.Error("Expected empty string for disabled item")
	}
}

func TestClickableItems(t *testing.T) {
	items := []Item{
		{ID: "edit", Label: "Edit"},
		{Type: ItemTypeDivider},
		{ID: "delete", Label: "Delete"},
		{Type: ItemTypeHeader, Label: "Section"},
		{ID: "disabled", Label: "Disabled", Disabled: true},
	}
	m := New("test", WithItems(items))

	clickable := m.ClickableItems()

	if len(clickable) != 2 {
		t.Errorf("Expected 2 clickable items, got %d", len(clickable))
	}
}

func TestHighlightNext(t *testing.T) {
	items := []Item{
		{ID: "1", Label: "One"},
		{ID: "2", Label: "Two"},
	}
	m := New("test", WithItems(items))

	m.HighlightNext()
	if m.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0, got %d", m.HighlightedIndex)
	}

	m.HighlightNext()
	if m.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1, got %d", m.HighlightedIndex)
	}

	// Should wrap
	m.HighlightNext()
	if m.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0 (wrapped), got %d", m.HighlightedIndex)
	}
}

func TestHighlightPrevious(t *testing.T) {
	items := []Item{
		{ID: "1", Label: "One"},
		{ID: "2", Label: "Two"},
	}
	m := New("test", WithItems(items))

	m.HighlightPrevious()
	if m.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1, got %d", m.HighlightedIndex)
	}

	m.HighlightPrevious()
	if m.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0, got %d", m.HighlightedIndex)
	}

	// Should wrap
	m.HighlightPrevious()
	if m.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1 (wrapped), got %d", m.HighlightedIndex)
	}
}

func TestIsHighlighted(t *testing.T) {
	items := []Item{
		{ID: "1", Label: "One"},
		{ID: "2", Label: "Two"},
	}
	m := New("test", WithItems(items))
	m.HighlightedIndex = 0

	if !m.IsHighlighted("1") {
		t.Error("Expected '1' to be highlighted")
	}
	if m.IsHighlighted("2") {
		t.Error("Expected '2' to not be highlighted")
	}
}

func TestGetItem(t *testing.T) {
	items := []Item{
		{ID: "edit", Label: "Edit"},
		{ID: "delete", Label: "Delete"},
	}
	m := New("test", WithItems(items))

	item := m.GetItem("edit")
	if item == nil {
		t.Fatal("Expected to find item 'edit'")
	}
	if item.Label != "Edit" {
		t.Errorf("Expected Label 'Edit', got '%s'", item.Label)
	}

	if m.GetItem("notfound") != nil {
		t.Error("Expected nil for not found item")
	}
}

func TestSetItemDisabled(t *testing.T) {
	items := []Item{
		{ID: "edit", Label: "Edit"},
	}
	m := New("test", WithItems(items))

	m.SetItemDisabled("edit", true)

	if !m.Items[0].Disabled {
		t.Error("Expected item to be disabled")
	}

	m.SetItemDisabled("edit", false)
	if m.Items[0].Disabled {
		t.Error("Expected item to be enabled")
	}
}

// ContextMenu tests

func TestShowAt(t *testing.T) {
	cm := NewContext("test")

	cm.ShowAt(100, 200)

	if cm.X != 100 {
		t.Errorf("Expected X 100, got %d", cm.X)
	}
	if cm.Y != 200 {
		t.Errorf("Expected Y 200, got %d", cm.Y)
	}
	if !cm.Open {
		t.Error("Expected Open to be true")
	}
}

// NavMenu tests

func TestWithNavItems(t *testing.T) {
	items := []Item{
		{ID: "home", Label: "Home"},
		{ID: "about", Label: "About"},
	}
	nm := NewNav("test", WithNavItems(items))

	if len(nm.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(nm.Items))
	}
}

func TestWithOrientation(t *testing.T) {
	nm := NewNav("test", WithOrientation("vertical"))
	if nm.Orientation != "vertical" {
		t.Errorf("Expected orientation 'vertical', got '%s'", nm.Orientation)
	}
}

func TestWithActiveID(t *testing.T) {
	items := []Item{
		{ID: "home", Label: "Home"},
		{ID: "about", Label: "About"},
	}
	nm := NewNav("test", WithNavItems(items), WithActiveID("home"))

	if nm.ActiveID != "home" {
		t.Errorf("Expected ActiveID 'home', got '%s'", nm.ActiveID)
	}
	if !nm.Items[0].Active {
		t.Error("Expected first item to be active")
	}
}

func TestToggleSubmenu(t *testing.T) {
	nm := NewNav("test")

	nm.ToggleSubmenu("products")
	if nm.OpenSubmenuID != "products" {
		t.Errorf("Expected OpenSubmenuID 'products', got '%s'", nm.OpenSubmenuID)
	}

	nm.ToggleSubmenu("products")
	if nm.OpenSubmenuID != "" {
		t.Errorf("Expected OpenSubmenuID to be empty, got '%s'", nm.OpenSubmenuID)
	}
}

func TestOpenSubmenu(t *testing.T) {
	nm := NewNav("test")
	nm.OpenSubmenu("products")

	if nm.OpenSubmenuID != "products" {
		t.Errorf("Expected OpenSubmenuID 'products', got '%s'", nm.OpenSubmenuID)
	}
}

func TestCloseSubmenu(t *testing.T) {
	nm := NewNav("test")
	nm.OpenSubmenuID = "products"
	nm.CloseSubmenu()

	if nm.OpenSubmenuID != "" {
		t.Error("Expected OpenSubmenuID to be empty")
	}
}

func TestIsSubmenuOpen(t *testing.T) {
	nm := NewNav("test")
	nm.OpenSubmenuID = "products"

	if !nm.IsSubmenuOpen("products") {
		t.Error("Expected 'products' submenu to be open")
	}
	if nm.IsSubmenuOpen("services") {
		t.Error("Expected 'services' submenu to not be open")
	}
}

func TestSetActive(t *testing.T) {
	items := []Item{
		{ID: "home", Label: "Home"},
		{
			ID:    "products",
			Label: "Products",
			Items: []Item{
				{ID: "product1", Label: "Product 1"},
			},
		},
	}
	nm := NewNav("test", WithNavItems(items))

	nm.SetActive("product1")

	if nm.ActiveID != "product1" {
		t.Errorf("Expected ActiveID 'product1', got '%s'", nm.ActiveID)
	}
	if !nm.Items[1].Items[0].Active {
		t.Error("Expected nested item to be active")
	}
}

func TestNavMenuGetItem(t *testing.T) {
	items := []Item{
		{ID: "home", Label: "Home"},
		{
			ID:    "products",
			Label: "Products",
			Items: []Item{
				{ID: "product1", Label: "Product 1"},
			},
		},
	}
	nm := NewNav("test", WithNavItems(items))

	// Top level
	item := nm.GetItem("home")
	if item == nil || item.ID != "home" {
		t.Error("Expected to find 'home' item")
	}

	// Nested
	nested := nm.GetItem("product1")
	if nested == nil || nested.ID != "product1" {
		t.Error("Expected to find nested 'product1' item")
	}

	// Not found
	if nm.GetItem("notfound") != nil {
		t.Error("Expected nil for not found")
	}
}

// Item helper tests

func TestItemIsDivider(t *testing.T) {
	divider := Item{Type: ItemTypeDivider}
	normal := Item{Type: ItemTypeDefault}

	if !divider.IsDivider() {
		t.Error("Expected divider to be divider")
	}
	if normal.IsDivider() {
		t.Error("Expected normal item to not be divider")
	}
}

func TestItemIsHeader(t *testing.T) {
	header := Item{Type: ItemTypeHeader}
	normal := Item{Type: ItemTypeDefault}

	if !header.IsHeader() {
		t.Error("Expected header to be header")
	}
	if normal.IsHeader() {
		t.Error("Expected normal item to not be header")
	}
}

func TestItemIsSubmenu(t *testing.T) {
	submenu1 := Item{Type: ItemTypeSubmenu}
	submenu2 := Item{Items: []Item{{ID: "child"}}}
	normal := Item{Type: ItemTypeDefault}

	if !submenu1.IsSubmenu() {
		t.Error("Expected submenu type to be submenu")
	}
	if !submenu2.IsSubmenu() {
		t.Error("Expected item with children to be submenu")
	}
	if normal.IsSubmenu() {
		t.Error("Expected normal item to not be submenu")
	}
}

func TestItemIsLink(t *testing.T) {
	link := Item{Href: "/page"}
	normal := Item{}

	if !link.IsLink() {
		t.Error("Expected item with href to be link")
	}
	if normal.IsLink() {
		t.Error("Expected normal item to not be link")
	}
}

func TestItemHasBadge(t *testing.T) {
	withBadge := Item{Badge: "New"}
	without := Item{}

	if !withBadge.HasBadge() {
		t.Error("Expected item with badge to have badge")
	}
	if without.HasBadge() {
		t.Error("Expected item without badge to not have badge")
	}
}

func TestItemHasIcon(t *testing.T) {
	withIcon := Item{Icon: "🏠"}
	without := Item{}

	if !withIcon.HasIcon() {
		t.Error("Expected item with icon to have icon")
	}
	if without.HasIcon() {
		t.Error("Expected item without icon to not have icon")
	}
}

func TestItemHasShortcut(t *testing.T) {
	with := Item{Shortcut: "⌘K"}
	without := Item{}

	if !with.HasShortcut() {
		t.Error("Expected item with shortcut to have shortcut")
	}
	if without.HasShortcut() {
		t.Error("Expected item without shortcut to not have shortcut")
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
