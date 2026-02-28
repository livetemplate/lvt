// Package accordion provides collapsible content sections for organizing information.
//
// Available variants:
//   - New() creates a basic accordion (template: "lvt:accordion:default:v1")
//   - NewSingle() creates single-open accordion (template: "lvt:accordion:single:v1")
//
// Required lvt-* attributes: lvt-click
//
// Example usage:
//
//	// In your controller/state
//	FAQ: accordion.New("faq", []accordion.Item{
//	    {ID: "q1", Title: "What is LiveTemplate?", Content: "..."},
//	    {ID: "q2", Title: "How do I get started?", Content: "..."},
//	})
//
//	// In your template
//	{{template "lvt:accordion:default:v1" .FAQ}}
package accordion

import (
	"github.com/livetemplate/components/base"
)

// Item represents a single accordion section.
type Item struct {
	ID       string // Unique identifier for this item
	Title    string // Header text shown (always visible)
	Content  string // Body content (shown when expanded)
	Disabled bool   // Whether this item can be toggled
	Icon     string // Optional icon (HTML or class name)
}

// Accordion is the base component for collapsible sections.
// Use template "lvt:accordion:default:v1" to render.
type Accordion struct {
	base.Base

	// Items is the list of accordion sections
	Items []Item

	// OpenIDs contains the IDs of currently open items
	OpenIDs map[string]bool

	// AllowMultiple allows multiple items to be open at once
	// When false, opening one item closes others (single-open mode)
	AllowMultiple bool
}

// New creates an accordion that allows multiple items open.
//
// Example:
//
//	faq := accordion.New("faq", []accordion.Item{
//	    {ID: "q1", Title: "Question 1", Content: "Answer 1"},
//	    {ID: "q2", Title: "Question 2", Content: "Answer 2"},
//	},
//	    accordion.WithOpen("q1"),
//	)
func New(id string, items []Item, opts ...Option) *Accordion {
	a := &Accordion{
		Base:          base.NewBase(id, "accordion"),
		Items:         items,
		OpenIDs:       make(map[string]bool),
		AllowMultiple: true,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// NewSingle creates an accordion where only one item can be open at a time.
//
// Example:
//
//	nav := accordion.NewSingle("nav", sections,
//	    accordion.WithOpen("section1"),
//	)
func NewSingle(id string, items []Item, opts ...Option) *Accordion {
	a := New(id, items, opts...)
	a.AllowMultiple = false
	return a
}

// Toggle toggles the open state of an item.
func (a *Accordion) Toggle(itemID string) {
	// Check if item exists and is not disabled
	for _, item := range a.Items {
		if item.ID == itemID {
			if item.Disabled {
				return
			}
			break
		}
	}

	if a.IsOpen(itemID) {
		// Close the item
		delete(a.OpenIDs, itemID)
	} else {
		// Open the item
		if !a.AllowMultiple {
			// Close all others first
			a.OpenIDs = make(map[string]bool)
		}
		a.OpenIDs[itemID] = true
	}
}

// Open opens an item.
func (a *Accordion) Open(itemID string) {
	// Check if item exists and is not disabled
	for _, item := range a.Items {
		if item.ID == itemID {
			if item.Disabled {
				return
			}
			break
		}
	}

	if !a.AllowMultiple {
		a.OpenIDs = make(map[string]bool)
	}
	a.OpenIDs[itemID] = true
}

// Close closes an item.
func (a *Accordion) Close(itemID string) {
	delete(a.OpenIDs, itemID)
}

// IsOpen checks if an item is currently open.
func (a *Accordion) IsOpen(itemID string) bool {
	return a.OpenIDs[itemID]
}

// OpenAll opens all non-disabled items.
func (a *Accordion) OpenAll() {
	for _, item := range a.Items {
		if !item.Disabled {
			a.OpenIDs[item.ID] = true
		}
	}
}

// CloseAll closes all items.
func (a *Accordion) CloseAll() {
	a.OpenIDs = make(map[string]bool)
}

// OpenCount returns the number of open items.
func (a *Accordion) OpenCount() int {
	return len(a.OpenIDs)
}

// ItemCount returns the total number of items.
func (a *Accordion) ItemCount() int {
	return len(a.Items)
}

// GetItem returns an item by ID, or nil if not found.
func (a *Accordion) GetItem(itemID string) *Item {
	for i := range a.Items {
		if a.Items[i].ID == itemID {
			return &a.Items[i]
		}
	}
	return nil
}

// AddItem adds a new item to the accordion.
func (a *Accordion) AddItem(item Item) {
	a.Items = append(a.Items, item)
}

// RemoveItem removes an item by ID.
func (a *Accordion) RemoveItem(itemID string) {
	for i, item := range a.Items {
		if item.ID == itemID {
			a.Items = append(a.Items[:i], a.Items[i+1:]...)
			delete(a.OpenIDs, itemID)
			return
		}
	}
}
