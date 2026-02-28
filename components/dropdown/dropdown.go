// Package dropdown provides dropdown/select components with single and multi-select support.
//
// Available variants:
//   - New() creates a basic dropdown (template: "lvt:dropdown:default:v1")
//   - NewSearchable() creates a searchable dropdown (template: "lvt:dropdown:searchable:v1")
//   - NewMulti() creates a multi-select dropdown (template: "lvt:dropdown:multi:v1")
//
// Required lvt-* attributes: lvt-click, lvt-click-away
// Optional: lvt-debounce (for searchable), lvt-focus-trap
//
// Example usage:
//
//	// In your controller/state
//	CountrySelect: dropdown.NewSearchable("country", countries,
//	    dropdown.WithPlaceholder("Select country"),
//	)
//
//	// In your template
//	{{template "lvt:dropdown:searchable:v1" .CountrySelect}}
package dropdown

import (
	"github.com/livetemplate/components/base"
)

// Item represents a single option in the dropdown.
type Item struct {
	Value    string // The value sent to the server when selected
	Label    string // The display text shown to users
	Disabled bool   // Whether this option is disabled
	Group    string // Optional group/category for grouped dropdowns
}

// Dropdown is a basic single-select dropdown component.
// Use template "lvt:dropdown:default:v1" to render.
type Dropdown struct {
	base.Base

	// Options is the list of selectable items
	Options []Item

	// Selected is the currently selected item (nil if none)
	Selected *Item

	// Placeholder is shown when nothing is selected
	Placeholder string

	// Open indicates whether the dropdown menu is currently visible
	Open bool

	// Disabled prevents user interaction
	Disabled bool
}

// New creates a basic single-select dropdown.
//
// Example:
//
//	countries := []dropdown.Item{
//	    {Value: "us", Label: "United States"},
//	    {Value: "ca", Label: "Canada"},
//	    {Value: "mx", Label: "Mexico"},
//	}
//	d := dropdown.New("country", countries,
//	    dropdown.WithPlaceholder("Select country"),
//	    dropdown.WithSelected("us"),
//	)
func New(id string, options []Item, opts ...Option) *Dropdown {
	d := &Dropdown{
		Base:        base.NewBase(id, "dropdown"),
		Options:     options,
		Placeholder: "Select...",
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Toggle opens or closes the dropdown.
func (d *Dropdown) Toggle() {
	d.Open = !d.Open
}

// Close closes the dropdown.
func (d *Dropdown) Close() {
	d.Open = false
}

// Select selects an item by value.
func (d *Dropdown) Select(value string) {
	for i := range d.Options {
		if d.Options[i].Value == value {
			d.Selected = &d.Options[i]
			d.Open = false
			return
		}
	}
}

// Clear clears the selection.
func (d *Dropdown) Clear() {
	d.Selected = nil
}

// Value returns the currently selected value, or empty string if none.
func (d *Dropdown) Value() string {
	if d.Selected != nil {
		return d.Selected.Value
	}
	return ""
}

// Searchable is a dropdown with search/filter capability.
// Use template "lvt:dropdown:searchable:v1" to render.
type Searchable struct {
	Dropdown

	// Query is the current search query
	Query string

	// FilteredOptions is the list of options matching the current query
	// If nil, all options are shown
	FilteredOptions []Item

	// MinChars is the minimum characters required before filtering starts
	MinChars int
}

// NewSearchable creates a searchable dropdown.
//
// Example:
//
//	d := dropdown.NewSearchable("country", countries,
//	    dropdown.WithPlaceholder("Search countries..."),
//	    dropdown.WithMinChars(2),
//	)
func NewSearchable(id string, options []Item, opts ...Option) *Searchable {
	s := &Searchable{
		Dropdown: Dropdown{
			Base:        base.NewBase(id, "dropdown"),
			Options:     options,
			Placeholder: "Search...",
		},
		MinChars: 1,
	}

	for _, opt := range opts {
		opt(&s.Dropdown)
	}

	return s
}

// Search filters options based on the query.
func (s *Searchable) Search(query string) {
	s.Query = query
	s.Open = true

	if len(query) < s.MinChars {
		s.FilteredOptions = nil
		return
	}

	s.FilteredOptions = make([]Item, 0)
	queryLower := toLower(query)
	for _, opt := range s.Options {
		if contains(toLower(opt.Label), queryLower) {
			s.FilteredOptions = append(s.FilteredOptions, opt)
		}
	}
}

// VisibleOptions returns the options to display (filtered if searching, all otherwise).
func (s *Searchable) VisibleOptions() []Item {
	if s.Query != "" && len(s.Query) >= s.MinChars {
		return s.FilteredOptions
	}
	return s.Options
}

// ClearSearch clears the search query and shows all options.
func (s *Searchable) ClearSearch() {
	s.Query = ""
	s.FilteredOptions = nil
}

// Multi is a multi-select dropdown with checkboxes.
// Use template "lvt:dropdown:multi:v1" to render.
type Multi struct {
	Dropdown

	// SelectedItems contains all selected items
	SelectedItems []Item

	// MaxSelections limits how many items can be selected (0 = unlimited)
	MaxSelections int
}

// NewMulti creates a multi-select dropdown.
//
// Example:
//
//	d := dropdown.NewMulti("tags", tagOptions,
//	    dropdown.WithPlaceholder("Select tags..."),
//	    dropdown.WithMaxSelections(5),
//	)
func NewMulti(id string, options []Item, opts ...Option) *Multi {
	m := &Multi{
		Dropdown: Dropdown{
			Base:        base.NewBase(id, "dropdown"),
			Options:     options,
			Placeholder: "Select...",
		},
		SelectedItems: make([]Item, 0),
	}

	for _, opt := range opts {
		opt(&m.Dropdown)
	}

	return m
}

// ToggleItem toggles selection of an item by value.
func (m *Multi) ToggleItem(value string) {
	// Check if already selected
	for i, item := range m.SelectedItems {
		if item.Value == value {
			// Remove from selection
			m.SelectedItems = append(m.SelectedItems[:i], m.SelectedItems[i+1:]...)
			return
		}
	}

	// Check max selections
	if m.MaxSelections > 0 && len(m.SelectedItems) >= m.MaxSelections {
		return
	}

	// Add to selection
	for _, opt := range m.Options {
		if opt.Value == value {
			m.SelectedItems = append(m.SelectedItems, opt)
			return
		}
	}
}

// IsSelected checks if an item is currently selected.
func (m *Multi) IsSelected(value string) bool {
	for _, item := range m.SelectedItems {
		if item.Value == value {
			return true
		}
	}
	return false
}

// Values returns all selected values.
func (m *Multi) Values() []string {
	values := make([]string, len(m.SelectedItems))
	for i, item := range m.SelectedItems {
		values[i] = item.Value
	}
	return values
}

// ClearAll clears all selections.
func (m *Multi) ClearAll() {
	m.SelectedItems = make([]Item, 0)
}

// SelectAll selects all non-disabled options.
func (m *Multi) SelectAll() {
	m.SelectedItems = make([]Item, 0)
	for _, opt := range m.Options {
		if !opt.Disabled {
			if m.MaxSelections > 0 && len(m.SelectedItems) >= m.MaxSelections {
				break
			}
			m.SelectedItems = append(m.SelectedItems, opt)
		}
	}
}

// DisplayText returns a summary of selected items for display.
func (m *Multi) DisplayText() string {
	count := len(m.SelectedItems)
	switch count {
	case 0:
		return m.Placeholder
	case 1:
		return m.SelectedItems[0].Label
	default:
		return m.SelectedItems[0].Label + " + " + itoa(count-1) + " more"
	}
}

// Helper functions to avoid importing strings/strconv
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
