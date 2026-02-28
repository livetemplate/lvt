// Package menu provides navigation and context menu components.
//
// Available variants:
//   - New() creates a dropdown menu (template: "lvt:menu:default:v1")
//   - NewContext() creates a context/right-click menu (template: "lvt:menu:context:v1")
//   - NewNav() creates a navigation menu (template: "lvt:menu:nav:v1")
//
// Required lvt-* attributes: lvt-click, lvt-click-away
//
// Example usage:
//
//	// In your controller/state
//	UserMenu: menu.New("user-menu",
//	    menu.WithItems([]menu.Item{
//	        {ID: "profile", Label: "Profile", Icon: "user"},
//	        {ID: "settings", Label: "Settings", Icon: "cog"},
//	        {Type: menu.ItemTypeDivider},
//	        {ID: "logout", Label: "Logout", Icon: "logout"},
//	    }),
//	)
//
//	// In your template
//	{{template "lvt:menu:default:v1" .UserMenu}}
package menu

import (
	"github.com/livetemplate/components/base"
)

// ItemType defines the type of menu item.
type ItemType int

const (
	// ItemTypeDefault is a normal clickable item
	ItemTypeDefault ItemType = iota
	// ItemTypeDivider is a separator line
	ItemTypeDivider
	// ItemTypeHeader is a non-clickable section header
	ItemTypeHeader
	// ItemTypeSubmenu is an item with nested items
	ItemTypeSubmenu
)

// Item represents a menu item.
type Item struct {
	// ID is the item identifier
	ID string
	// Type is the item type (default, divider, header, submenu)
	Type ItemType
	// Label is the display text
	Label string
	// Icon is an optional icon class/name
	Icon string
	// Shortcut is an optional keyboard shortcut text
	Shortcut string
	// Disabled prevents interaction
	Disabled bool
	// Href makes this a link
	Href string
	// Target for link (e.g., "_blank")
	Target string
	// Items holds submenu items (when Type is ItemTypeSubmenu)
	Items []Item
	// Data holds arbitrary custom data
	Data map[string]any
	// Badge is optional badge text
	Badge string
	// BadgeColor is the badge color (e.g., "red", "blue")
	BadgeColor string
	// Active highlights the item
	Active bool
}

// Menu is a dropdown/action menu component.
// Use template "lvt:menu:default:v1" to render.
type Menu struct {
	base.Base

	// Items is the list of menu items
	Items []Item

	// Open indicates whether the menu is visible
	Open bool

	// Trigger is the button/element that opens the menu
	Trigger string

	// TriggerIcon is an optional icon for the trigger
	TriggerIcon string

	// Position is the menu position relative to trigger ("bottom-left", "bottom-right", etc.)
	Position string

	// HighlightedIndex is the currently highlighted item (-1 for none)
	HighlightedIndex int
}

// ContextMenu is a right-click context menu.
// Use template "lvt:menu:context:v1" to render.
type ContextMenu struct {
	Menu

	// X is the horizontal position
	X int
	// Y is the vertical position
	Y int
}

// NavMenu is a navigation menu with support for nested items.
// Use template "lvt:menu:nav:v1" to render.
type NavMenu struct {
	base.Base

	// Items is the list of menu items
	Items []Item

	// Orientation is "horizontal" or "vertical"
	Orientation string

	// OpenSubmenuID is the ID of the currently open submenu
	OpenSubmenuID string

	// ActiveID is the ID of the active/current item
	ActiveID string
}

// New creates a dropdown menu.
//
// Example:
//
//	m := menu.New("actions",
//	    menu.WithTrigger("Actions"),
//	    menu.WithItems(items),
//	)
func New(id string, opts ...Option) *Menu {
	m := &Menu{
		Base:             base.NewBase(id, "menu"),
		Position:         "bottom-left",
		HighlightedIndex: -1,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// NewContext creates a context menu.
func NewContext(id string, opts ...Option) *ContextMenu {
	m := New(id, opts...)
	return &ContextMenu{
		Menu: *m,
	}
}

// NewNav creates a navigation menu.
func NewNav(id string, opts ...NavOption) *NavMenu {
	nm := &NavMenu{
		Base:        base.NewBase(id, "menu"),
		Orientation: "horizontal",
	}

	for _, opt := range opts {
		opt(nm)
	}

	return nm
}

// Toggle opens or closes the menu.
func (m *Menu) Toggle() {
	m.Open = !m.Open
	if !m.Open {
		m.HighlightedIndex = -1
	}
}

// Show opens the menu.
func (m *Menu) Show() {
	m.Open = true
}

// Close closes the menu.
func (m *Menu) Close() {
	m.Open = false
	m.HighlightedIndex = -1
}

// SelectIndex selects the item at the given index.
func (m *Menu) SelectIndex(index int) string {
	clickableItems := m.ClickableItems()
	if index < 0 || index >= len(clickableItems) {
		return ""
	}

	item := clickableItems[index]
	if item.Disabled {
		return ""
	}

	m.Close()
	return item.ID
}

// ClickableItems returns items that can be clicked (excludes dividers, headers, disabled).
func (m *Menu) ClickableItems() []Item {
	var clickable []Item
	for _, item := range m.Items {
		if item.Type == ItemTypeDefault && !item.Disabled {
			clickable = append(clickable, item)
		}
	}
	return clickable
}

// HighlightNext moves highlight to the next clickable item.
func (m *Menu) HighlightNext() {
	clickable := m.ClickableItems()
	if len(clickable) == 0 {
		return
	}

	m.HighlightedIndex++
	if m.HighlightedIndex >= len(clickable) {
		m.HighlightedIndex = 0
	}
}

// HighlightPrevious moves highlight to the previous clickable item.
func (m *Menu) HighlightPrevious() {
	clickable := m.ClickableItems()
	if len(clickable) == 0 {
		return
	}

	m.HighlightedIndex--
	if m.HighlightedIndex < 0 {
		m.HighlightedIndex = len(clickable) - 1
	}
}

// IsHighlighted checks if an item ID is highlighted.
func (m *Menu) IsHighlighted(id string) bool {
	clickable := m.ClickableItems()
	if m.HighlightedIndex < 0 || m.HighlightedIndex >= len(clickable) {
		return false
	}
	return clickable[m.HighlightedIndex].ID == id
}

// GetItem returns an item by ID.
func (m *Menu) GetItem(id string) *Item {
	for i := range m.Items {
		if m.Items[i].ID == id {
			return &m.Items[i]
		}
	}
	return nil
}

// SetItemDisabled enables or disables an item.
func (m *Menu) SetItemDisabled(id string, disabled bool) {
	if item := m.GetItem(id); item != nil {
		item.Disabled = disabled
	}
}

// ContextMenu methods

// ShowAt shows the context menu at the specified position.
func (cm *ContextMenu) ShowAt(x, y int) {
	cm.X = x
	cm.Y = y
	cm.Open = true
}

// NavMenu methods

// ToggleSubmenu opens or closes a submenu.
func (nm *NavMenu) ToggleSubmenu(id string) {
	if nm.OpenSubmenuID == id {
		nm.OpenSubmenuID = ""
	} else {
		nm.OpenSubmenuID = id
	}
}

// OpenSubmenu opens a specific submenu.
func (nm *NavMenu) OpenSubmenu(id string) {
	nm.OpenSubmenuID = id
}

// CloseSubmenu closes the open submenu.
func (nm *NavMenu) CloseSubmenu() {
	nm.OpenSubmenuID = ""
}

// IsSubmenuOpen checks if a submenu is open.
func (nm *NavMenu) IsSubmenuOpen(id string) bool {
	return nm.OpenSubmenuID == id
}

// SetActive sets the active item.
func (nm *NavMenu) SetActive(id string) {
	nm.ActiveID = id
	// Also update Active flag on items
	for i := range nm.Items {
		nm.Items[i].Active = nm.Items[i].ID == id
		// Check nested items
		for j := range nm.Items[i].Items {
			nm.Items[i].Items[j].Active = nm.Items[i].Items[j].ID == id
		}
	}
}

// IsActive checks if an item is active.
func (nm *NavMenu) IsActive(id string) bool {
	return nm.ActiveID == id
}

// GetItem returns an item by ID (searches nested items too).
func (nm *NavMenu) GetItem(id string) *Item {
	for i := range nm.Items {
		if nm.Items[i].ID == id {
			return &nm.Items[i]
		}
		for j := range nm.Items[i].Items {
			if nm.Items[i].Items[j].ID == id {
				return &nm.Items[i].Items[j]
			}
		}
	}
	return nil
}

// Helper functions for templates

// IsDivider checks if item is a divider.
func (i Item) IsDivider() bool {
	return i.Type == ItemTypeDivider
}

// IsHeader checks if item is a header.
func (i Item) IsHeader() bool {
	return i.Type == ItemTypeHeader
}

// IsSubmenu checks if item has submenu items.
func (i Item) IsSubmenu() bool {
	return i.Type == ItemTypeSubmenu || len(i.Items) > 0
}

// IsLink checks if item is a link.
func (i Item) IsLink() bool {
	return i.Href != ""
}

// HasBadge checks if item has a badge.
func (i Item) HasBadge() bool {
	return i.Badge != ""
}

// HasIcon checks if item has an icon.
func (i Item) HasIcon() bool {
	return i.Icon != ""
}

// HasShortcut checks if item has a keyboard shortcut.
func (i Item) HasShortcut() bool {
	return i.Shortcut != ""
}
