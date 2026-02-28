// Package tabs provides tab navigation components for organizing content.
//
// Available variants:
//   - New() creates horizontal tabs (template: "lvt:tabs:horizontal:v1")
//   - NewVertical() creates vertical tabs (template: "lvt:tabs:vertical:v1")
//   - NewPills() creates pill-style tabs (template: "lvt:tabs:pills:v1")
//
// Required lvt-* attributes: lvt-click
//
// Example usage:
//
//	// In your controller/state
//	SettingsTabs: tabs.New("settings", []tabs.Tab{
//	    {ID: "general", Label: "General"},
//	    {ID: "security", Label: "Security"},
//	    {ID: "notifications", Label: "Notifications"},
//	},
//	    tabs.WithActive("general"),
//	)
//
//	// In your template
//	{{template "lvt:tabs:horizontal:v1" .SettingsTabs}}
//	{{if eq .SettingsTabs.ActiveID "general"}}
//	  <!-- General settings content -->
//	{{end}}
package tabs

import (
	"github.com/livetemplate/components/base"
)

// Tab represents a single tab item.
type Tab struct {
	ID       string // Unique identifier for this tab
	Label    string // Display text shown in the tab
	Icon     string // Optional icon (HTML or class name)
	Disabled bool   // Whether this tab is disabled
	Badge    string // Optional badge text (e.g., count)
}

// Tabs is the base component for tab navigation.
// Use template "lvt:tabs:horizontal:v1" to render.
type Tabs struct {
	base.Base

	// Items is the list of tabs
	Items []Tab

	// ActiveID is the ID of the currently active tab
	ActiveID string
}

// New creates horizontal tabs.
//
// Example:
//
//	tabs := tabs.New("settings", []tabs.Tab{
//	    {ID: "general", Label: "General"},
//	    {ID: "security", Label: "Security"},
//	    {ID: "notifications", Label: "Notifications"},
//	},
//	    tabs.WithActive("general"),
//	)
func New(id string, items []Tab, opts ...Option) *Tabs {
	t := &Tabs{
		Base:  base.NewBase(id, "tabs"),
		Items: items,
	}

	// Default to first tab if items exist
	if len(items) > 0 {
		t.ActiveID = items[0].ID
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// NewVertical creates vertical tabs (sidebar style).
//
// Example:
//
//	tabs := tabs.NewVertical("nav", navItems,
//	    tabs.WithActive("dashboard"),
//	)
func NewVertical(id string, items []Tab, opts ...Option) *Tabs {
	return New(id, items, opts...)
}

// NewPills creates pill-style tabs.
//
// Example:
//
//	tabs := tabs.NewPills("filter", filterItems,
//	    tabs.WithActive("all"),
//	)
func NewPills(id string, items []Tab, opts ...Option) *Tabs {
	return New(id, items, opts...)
}

// SetActive sets the active tab by ID.
func (t *Tabs) SetActive(tabID string) {
	// Only set if it's a valid, non-disabled tab
	for _, tab := range t.Items {
		if tab.ID == tabID && !tab.Disabled {
			t.ActiveID = tabID
			return
		}
	}
}

// ActiveTab returns the currently active tab, or nil if none.
func (t *Tabs) ActiveTab() *Tab {
	for i := range t.Items {
		if t.Items[i].ID == t.ActiveID {
			return &t.Items[i]
		}
	}
	return nil
}

// IsActive checks if a tab is currently active.
func (t *Tabs) IsActive(tabID string) bool {
	return t.ActiveID == tabID
}

// Next activates the next non-disabled tab (wraps around).
func (t *Tabs) Next() {
	if len(t.Items) == 0 {
		return
	}

	currentIdx := t.findActiveIndex()
	for i := 1; i <= len(t.Items); i++ {
		nextIdx := (currentIdx + i) % len(t.Items)
		if !t.Items[nextIdx].Disabled {
			t.ActiveID = t.Items[nextIdx].ID
			return
		}
	}
}

// Previous activates the previous non-disabled tab (wraps around).
func (t *Tabs) Previous() {
	if len(t.Items) == 0 {
		return
	}

	currentIdx := t.findActiveIndex()
	for i := 1; i <= len(t.Items); i++ {
		prevIdx := (currentIdx - i + len(t.Items)) % len(t.Items)
		if !t.Items[prevIdx].Disabled {
			t.ActiveID = t.Items[prevIdx].ID
			return
		}
	}
}

// findActiveIndex returns the index of the active tab.
func (t *Tabs) findActiveIndex() int {
	for i, tab := range t.Items {
		if tab.ID == t.ActiveID {
			return i
		}
	}
	return 0
}

// AddTab adds a new tab to the list.
func (t *Tabs) AddTab(tab Tab) {
	t.Items = append(t.Items, tab)
	// If this is the first tab, make it active
	if len(t.Items) == 1 {
		t.ActiveID = tab.ID
	}
}

// RemoveTab removes a tab by ID.
func (t *Tabs) RemoveTab(tabID string) {
	for i, tab := range t.Items {
		if tab.ID == tabID {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			// If we removed the active tab, activate the first available
			if t.ActiveID == tabID && len(t.Items) > 0 {
				t.SetActive(t.Items[0].ID)
			}
			return
		}
	}
}

// TabCount returns the number of tabs.
func (t *Tabs) TabCount() int {
	return len(t.Items)
}

// EnabledTabCount returns the number of non-disabled tabs.
func (t *Tabs) EnabledTabCount() int {
	count := 0
	for _, tab := range t.Items {
		if !tab.Disabled {
			count++
		}
	}
	return count
}
