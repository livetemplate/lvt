// Package autocomplete provides typeahead/autocomplete input components.
//
// Available variants:
//   - New() creates a basic autocomplete (template: "lvt:autocomplete:default:v1")
//   - NewMulti() creates a multi-select autocomplete (template: "lvt:autocomplete:multi:v1")
//
// Required lvt-* attributes: lvt-input, lvt-click, lvt-focus, lvt-blur, lvt-click-away
//
// Example usage:
//
//	// In your controller/state
//	Search: autocomplete.New("search",
//	    autocomplete.WithPlaceholder("Search cities..."),
//	    autocomplete.WithSuggestions([]autocomplete.Suggestion{
//	        {Value: "nyc", Label: "New York City"},
//	        {Value: "la", Label: "Los Angeles"},
//	    }),
//	)
//
//	// In your template
//	{{template "lvt:autocomplete:default:v1" .Search}}
package autocomplete

import (
	"strings"

	"github.com/livetemplate/components/base"
)

// Suggestion represents an autocomplete suggestion.
type Suggestion struct {
	// Value is the internal value
	Value string
	// Label is the display text
	Label string
	// Description is optional secondary text
	Description string
	// Icon is an optional icon class/name
	Icon string
	// Disabled prevents selection
	Disabled bool
	// Data holds arbitrary custom data
	Data map[string]any
}

// Autocomplete is a typeahead input component.
// Use template "lvt:autocomplete:default:v1" to render.
type Autocomplete struct {
	base.Base

	// Query is the current input text
	Query string

	// Selected is the currently selected suggestion (nil if none)
	Selected *Suggestion

	// Suggestions is the full list of available suggestions
	Suggestions []Suggestion

	// FilteredSuggestions is the currently shown suggestions
	FilteredSuggestions []Suggestion

	// Placeholder text shown when empty
	Placeholder string

	// MinChars is the minimum characters before showing suggestions (default: 1)
	MinChars int

	// MaxSuggestions limits shown suggestions (0 for unlimited)
	MaxSuggestions int

	// Open indicates whether suggestions are visible
	Open bool

	// HighlightedIndex is the currently highlighted suggestion (-1 for none)
	HighlightedIndex int

	// Loading indicates an async search is in progress
	Loading bool

	// AllowCustom allows values not in suggestions
	AllowCustom bool

	// ClearOnSelect clears input after selection (useful for multi)
	ClearOnSelect bool

	// filterFunc is a custom filter function
	filterFunc func(query string, suggestions []Suggestion) []Suggestion
}

// MultiAutocomplete allows selecting multiple suggestions.
// Use template "lvt:autocomplete:multi:v1" to render.
type MultiAutocomplete struct {
	Autocomplete

	// SelectedItems contains all selected suggestions
	SelectedItems []Suggestion
}

// New creates a basic autocomplete.
//
// Example:
//
//	ac := autocomplete.New("search",
//	    autocomplete.WithPlaceholder("Search..."),
//	    autocomplete.WithSuggestions(suggestions),
//	)
func New(id string, opts ...Option) *Autocomplete {
	ac := &Autocomplete{
		Base:             base.NewBase(id, "autocomplete"),
		Placeholder:      "Type to search...",
		MinChars:         1,
		MaxSuggestions:   10,
		HighlightedIndex: -1,
	}

	for _, opt := range opts {
		opt(ac)
	}

	return ac
}

// NewMulti creates a multi-select autocomplete.
func NewMulti(id string, opts ...Option) *MultiAutocomplete {
	ac := New(id, opts...)
	ac.ClearOnSelect = true
	return &MultiAutocomplete{
		Autocomplete: *ac,
	}
}

// SetQuery updates the search query and filters suggestions.
func (ac *Autocomplete) SetQuery(query string) {
	ac.Query = query
	ac.Filter()
	ac.HighlightedIndex = -1

	// Show suggestions if query meets minimum
	ac.Open = len(query) >= ac.MinChars && len(ac.FilteredSuggestions) > 0
}

// Filter filters suggestions based on the current query.
func (ac *Autocomplete) Filter() {
	if ac.filterFunc != nil {
		ac.FilteredSuggestions = ac.filterFunc(ac.Query, ac.Suggestions)
	} else {
		ac.FilteredSuggestions = ac.defaultFilter(ac.Query)
	}

	// Apply max limit
	if ac.MaxSuggestions > 0 && len(ac.FilteredSuggestions) > ac.MaxSuggestions {
		ac.FilteredSuggestions = ac.FilteredSuggestions[:ac.MaxSuggestions]
	}
}

// defaultFilter performs case-insensitive substring matching.
func (ac *Autocomplete) defaultFilter(query string) []Suggestion {
	if query == "" {
		return ac.Suggestions
	}

	queryLower := strings.ToLower(query)
	var filtered []Suggestion

	for _, s := range ac.Suggestions {
		if strings.Contains(strings.ToLower(s.Label), queryLower) ||
			strings.Contains(strings.ToLower(s.Value), queryLower) ||
			strings.Contains(strings.ToLower(s.Description), queryLower) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// SelectIndex selects the suggestion at the given index.
func (ac *Autocomplete) SelectIndex(index int) bool {
	if index < 0 || index >= len(ac.FilteredSuggestions) {
		return false
	}

	s := ac.FilteredSuggestions[index]
	if s.Disabled {
		return false
	}

	return ac.Select(s)
}

// Select selects a suggestion.
func (ac *Autocomplete) Select(s Suggestion) bool {
	if s.Disabled {
		return false
	}

	ac.Selected = &s
	ac.Query = s.Label
	ac.Open = false
	ac.HighlightedIndex = -1

	if ac.ClearOnSelect {
		ac.Query = ""
	}

	return true
}

// SelectHighlighted selects the currently highlighted suggestion.
func (ac *Autocomplete) SelectHighlighted() bool {
	if ac.HighlightedIndex < 0 || ac.HighlightedIndex >= len(ac.FilteredSuggestions) {
		// If custom values allowed and we have a query, accept it
		if ac.AllowCustom && ac.Query != "" {
			ac.Selected = &Suggestion{
				Value: ac.Query,
				Label: ac.Query,
			}
			ac.Open = false
			return true
		}
		return false
	}

	return ac.SelectIndex(ac.HighlightedIndex)
}

// Clear clears the selection and query.
func (ac *Autocomplete) Clear() {
	ac.Selected = nil
	ac.Query = ""
	ac.Open = false
	ac.HighlightedIndex = -1
	ac.FilteredSuggestions = nil
}

// Focus opens suggestions if query meets minimum.
func (ac *Autocomplete) Focus() {
	ac.Filter()
	if len(ac.Query) >= ac.MinChars && len(ac.FilteredSuggestions) > 0 {
		ac.Open = true
	}
}

// Blur closes suggestions.
func (ac *Autocomplete) Blur() {
	ac.Open = false
	ac.HighlightedIndex = -1
}

// HighlightNext moves highlight to the next suggestion.
func (ac *Autocomplete) HighlightNext() {
	if len(ac.FilteredSuggestions) == 0 {
		return
	}

	ac.HighlightedIndex++
	if ac.HighlightedIndex >= len(ac.FilteredSuggestions) {
		ac.HighlightedIndex = 0
	}

	// Skip disabled
	for i := 0; i < len(ac.FilteredSuggestions); i++ {
		if !ac.FilteredSuggestions[ac.HighlightedIndex].Disabled {
			return
		}
		ac.HighlightedIndex = (ac.HighlightedIndex + 1) % len(ac.FilteredSuggestions)
	}
}

// HighlightPrevious moves highlight to the previous suggestion.
func (ac *Autocomplete) HighlightPrevious() {
	if len(ac.FilteredSuggestions) == 0 {
		return
	}

	ac.HighlightedIndex--
	if ac.HighlightedIndex < 0 {
		ac.HighlightedIndex = len(ac.FilteredSuggestions) - 1
	}

	// Skip disabled
	for i := 0; i < len(ac.FilteredSuggestions); i++ {
		if !ac.FilteredSuggestions[ac.HighlightedIndex].Disabled {
			return
		}
		ac.HighlightedIndex--
		if ac.HighlightedIndex < 0 {
			ac.HighlightedIndex = len(ac.FilteredSuggestions) - 1
		}
	}
}

// IsHighlighted checks if a suggestion index is highlighted.
func (ac *Autocomplete) IsHighlighted(index int) bool {
	return ac.HighlightedIndex == index
}

// HasSelection returns true if a suggestion is selected.
func (ac *Autocomplete) HasSelection() bool {
	return ac.Selected != nil
}

// DisplayValue returns the display text (selected label or query).
func (ac *Autocomplete) DisplayValue() string {
	if ac.Selected != nil {
		return ac.Selected.Label
	}
	return ac.Query
}

// SetSuggestions updates the suggestions and re-filters.
func (ac *Autocomplete) SetSuggestions(suggestions []Suggestion) {
	ac.Suggestions = suggestions
	ac.Filter()
}

// SetLoading sets the loading state.
func (ac *Autocomplete) SetLoading(loading bool) {
	ac.Loading = loading
}

// MultiAutocomplete methods

// SelectMulti adds a suggestion to selected items.
func (mac *MultiAutocomplete) SelectMulti(s Suggestion) bool {
	if s.Disabled {
		return false
	}

	// Check if already selected
	for _, item := range mac.SelectedItems {
		if item.Value == s.Value {
			return false
		}
	}

	mac.SelectedItems = append(mac.SelectedItems, s)
	mac.Query = ""
	mac.Open = false
	mac.HighlightedIndex = -1

	return true
}

// RemoveSelected removes a suggestion from selected items by value.
func (mac *MultiAutocomplete) RemoveSelected(value string) bool {
	for i, item := range mac.SelectedItems {
		if item.Value == value {
			mac.SelectedItems = append(mac.SelectedItems[:i], mac.SelectedItems[i+1:]...)
			return true
		}
	}
	return false
}

// ClearMulti clears all selected items.
func (mac *MultiAutocomplete) ClearMulti() {
	mac.SelectedItems = nil
	mac.Query = ""
	mac.Open = false
	mac.HighlightedIndex = -1
}

// IsSelectedMulti checks if a value is in selected items.
func (mac *MultiAutocomplete) IsSelectedMulti(value string) bool {
	for _, item := range mac.SelectedItems {
		if item.Value == value {
			return true
		}
	}
	return false
}

// HasSelectedItems returns true if any items are selected.
func (mac *MultiAutocomplete) HasSelectedItems() bool {
	return len(mac.SelectedItems) > 0
}

// SelectedValues returns the values of all selected items.
func (mac *MultiAutocomplete) SelectedValues() []string {
	values := make([]string, len(mac.SelectedItems))
	for i, item := range mac.SelectedItems {
		values[i] = item.Value
	}
	return values
}

// FilteredExcludingSelected returns filtered suggestions excluding already selected items.
func (mac *MultiAutocomplete) FilteredExcludingSelected() []Suggestion {
	var filtered []Suggestion
	for _, s := range mac.FilteredSuggestions {
		if !mac.IsSelectedMulti(s.Value) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
