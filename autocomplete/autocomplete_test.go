package autocomplete

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ac := New("test-ac")

	if ac.ID() != "test-ac" {
		t.Errorf("Expected ID 'test-ac', got '%s'", ac.ID())
	}
	if ac.Namespace() != "autocomplete" {
		t.Errorf("Expected namespace 'autocomplete', got '%s'", ac.Namespace())
	}
	if ac.Placeholder != "Type to search..." {
		t.Errorf("Expected default placeholder 'Type to search...', got '%s'", ac.Placeholder)
	}
	if ac.MinChars != 1 {
		t.Errorf("Expected MinChars 1, got %d", ac.MinChars)
	}
	if ac.MaxSuggestions != 10 {
		t.Errorf("Expected MaxSuggestions 10, got %d", ac.MaxSuggestions)
	}
	if ac.HighlightedIndex != -1 {
		t.Errorf("Expected HighlightedIndex -1, got %d", ac.HighlightedIndex)
	}
}

func TestNewMulti(t *testing.T) {
	mac := NewMulti("test-multi")

	if mac.ID() != "test-multi" {
		t.Errorf("Expected ID 'test-multi', got '%s'", mac.ID())
	}
	if !mac.ClearOnSelect {
		t.Error("Expected ClearOnSelect to be true for multi")
	}
	if mac.SelectedItems != nil && len(mac.SelectedItems) > 0 {
		t.Error("Expected SelectedItems to be empty")
	}
}

func TestWithPlaceholder(t *testing.T) {
	ac := New("test", WithPlaceholder("Search cities..."))
	if ac.Placeholder != "Search cities..." {
		t.Errorf("Expected placeholder 'Search cities...', got '%s'", ac.Placeholder)
	}
}

func TestWithSuggestions(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	ac := New("test", WithSuggestions(suggestions))

	if len(ac.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(ac.Suggestions))
	}
}

func TestWithMinChars(t *testing.T) {
	ac := New("test", WithMinChars(3))
	if ac.MinChars != 3 {
		t.Errorf("Expected MinChars 3, got %d", ac.MinChars)
	}

	// Test negative value
	ac2 := New("test2", WithMinChars(-1))
	if ac2.MinChars != 0 {
		t.Errorf("Expected MinChars 0 for negative input, got %d", ac2.MinChars)
	}
}

func TestWithMaxSuggestions(t *testing.T) {
	ac := New("test", WithMaxSuggestions(5))
	if ac.MaxSuggestions != 5 {
		t.Errorf("Expected MaxSuggestions 5, got %d", ac.MaxSuggestions)
	}
}

func TestWithAllowCustom(t *testing.T) {
	ac := New("test", WithAllowCustom(true))
	if !ac.AllowCustom {
		t.Error("Expected AllowCustom to be true")
	}
}

func TestWithClearOnSelect(t *testing.T) {
	ac := New("test", WithClearOnSelect(true))
	if !ac.ClearOnSelect {
		t.Error("Expected ClearOnSelect to be true")
	}
}

func TestWithStyled(t *testing.T) {
	ac := New("test", WithStyled(false))
	if ac.IsStyled() {
		t.Error("Expected IsStyled to be false")
	}
}

func TestWithFilterFunc(t *testing.T) {
	customFilter := func(query string, suggestions []Suggestion) []Suggestion {
		var filtered []Suggestion
		for _, s := range suggestions {
			if strings.HasPrefix(strings.ToLower(s.Label), strings.ToLower(query)) {
				filtered = append(filtered, s)
			}
		}
		return filtered
	}

	suggestions := []Suggestion{
		{Value: "1", Label: "Apple"},
		{Value: "2", Label: "Banana"},
		{Value: "3", Label: "Apricot"},
	}
	ac := New("test", WithSuggestions(suggestions), WithFilterFunc(customFilter))

	ac.SetQuery("ap")

	if len(ac.FilteredSuggestions) != 2 {
		t.Errorf("Expected 2 filtered suggestions (Apple, Apricot), got %d", len(ac.FilteredSuggestions))
	}
}

func TestWithQuery(t *testing.T) {
	ac := New("test", WithQuery("initial"))
	if ac.Query != "initial" {
		t.Errorf("Expected Query 'initial', got '%s'", ac.Query)
	}
}

func TestWithSelected(t *testing.T) {
	s := Suggestion{Value: "1", Label: "Option 1"}
	ac := New("test", WithSelected(s))

	if ac.Selected == nil {
		t.Fatal("Expected Selected to be set")
	}
	if ac.Selected.Value != "1" {
		t.Errorf("Expected Selected.Value '1', got '%s'", ac.Selected.Value)
	}
	if ac.Query != "Option 1" {
		t.Errorf("Expected Query to be set to label, got '%s'", ac.Query)
	}
}

func TestSetQuery(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "Apple"},
		{Value: "2", Label: "Banana"},
	}
	ac := New("test", WithSuggestions(suggestions))

	ac.SetQuery("app")

	if ac.Query != "app" {
		t.Errorf("Expected Query 'app', got '%s'", ac.Query)
	}
	if len(ac.FilteredSuggestions) != 1 {
		t.Errorf("Expected 1 filtered suggestion, got %d", len(ac.FilteredSuggestions))
	}
	if !ac.Open {
		t.Error("Expected Open to be true when matching suggestions found")
	}
	if ac.HighlightedIndex != -1 {
		t.Errorf("Expected HighlightedIndex reset to -1, got %d", ac.HighlightedIndex)
	}
}

func TestSetQueryBelowMinChars(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "Apple"},
	}
	ac := New("test", WithSuggestions(suggestions), WithMinChars(3))

	ac.SetQuery("ap") // Only 2 chars

	if ac.Open {
		t.Error("Expected Open to be false when query below MinChars")
	}
}

func TestFilter(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "New York"},
		{Value: "2", Label: "Los Angeles"},
		{Value: "3", Label: "New Orleans"},
	}
	ac := New("test", WithSuggestions(suggestions))

	ac.Query = "new"
	ac.Filter()

	if len(ac.FilteredSuggestions) != 2 {
		t.Errorf("Expected 2 filtered suggestions, got %d", len(ac.FilteredSuggestions))
	}
}

func TestFilterMaxSuggestions(t *testing.T) {
	suggestions := make([]Suggestion, 20)
	for i := 0; i < 20; i++ {
		suggestions[i] = Suggestion{Value: string(rune('a' + i)), Label: "Item"}
	}
	ac := New("test", WithSuggestions(suggestions), WithMaxSuggestions(5))

	ac.Query = "Item"
	ac.Filter()

	if len(ac.FilteredSuggestions) != 5 {
		t.Errorf("Expected 5 filtered suggestions (MaxSuggestions), got %d", len(ac.FilteredSuggestions))
	}
}

func TestFilterByDescription(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "City", Description: "New York"},
		{Value: "2", Label: "City", Description: "Los Angeles"},
	}
	ac := New("test", WithSuggestions(suggestions))

	ac.Query = "york"
	ac.Filter()

	if len(ac.FilteredSuggestions) != 1 {
		t.Errorf("Expected 1 filtered suggestion (by description), got %d", len(ac.FilteredSuggestions))
	}
}

func TestSelectIndex(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.SetQuery("o")

	result := ac.SelectIndex(0)

	if !result {
		t.Error("Expected SelectIndex to return true")
	}
	if ac.Selected == nil {
		t.Fatal("Expected Selected to be set")
	}
	if ac.Selected.Value != "1" {
		t.Errorf("Expected Selected.Value '1', got '%s'", ac.Selected.Value)
	}
	if ac.Open {
		t.Error("Expected Open to be false after selection")
	}
}

func TestSelectIndexOutOfBounds(t *testing.T) {
	ac := New("test")

	if ac.SelectIndex(0) {
		t.Error("Expected SelectIndex to return false for out of bounds")
	}
	if ac.SelectIndex(-1) {
		t.Error("Expected SelectIndex to return false for negative index")
	}
}

func TestSelectIndexDisabled(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One", Disabled: true},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.Filter()

	if ac.SelectIndex(0) {
		t.Error("Expected SelectIndex to return false for disabled suggestion")
	}
}

func TestSelect(t *testing.T) {
	ac := New("test")
	s := Suggestion{Value: "1", Label: "One"}

	result := ac.Select(s)

	if !result {
		t.Error("Expected Select to return true")
	}
	if ac.Selected == nil || ac.Selected.Value != "1" {
		t.Error("Expected suggestion to be selected")
	}
	if ac.Query != "One" {
		t.Errorf("Expected Query to be 'One', got '%s'", ac.Query)
	}
}

func TestSelectClearOnSelect(t *testing.T) {
	ac := New("test", WithClearOnSelect(true))
	s := Suggestion{Value: "1", Label: "One"}

	ac.Select(s)

	if ac.Query != "" {
		t.Errorf("Expected Query to be empty with ClearOnSelect, got '%s'", ac.Query)
	}
}

func TestSelectHighlighted(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.Filter()
	ac.HighlightedIndex = 1

	result := ac.SelectHighlighted()

	if !result {
		t.Error("Expected SelectHighlighted to return true")
	}
	if ac.Selected == nil || ac.Selected.Value != "2" {
		t.Error("Expected second suggestion to be selected")
	}
}

func TestSelectHighlightedNoHighlight(t *testing.T) {
	ac := New("test")

	if ac.SelectHighlighted() {
		t.Error("Expected SelectHighlighted to return false when nothing highlighted")
	}
}

func TestSelectHighlightedAllowCustom(t *testing.T) {
	ac := New("test", WithAllowCustom(true), WithQuery("custom value"))

	result := ac.SelectHighlighted()

	if !result {
		t.Error("Expected SelectHighlighted to return true with AllowCustom")
	}
	if ac.Selected == nil || ac.Selected.Value != "custom value" {
		t.Error("Expected custom value to be selected")
	}
}

func TestClear(t *testing.T) {
	s := Suggestion{Value: "1", Label: "One"}
	ac := New("test", WithSelected(s))

	ac.Clear()

	if ac.Selected != nil {
		t.Error("Expected Selected to be nil")
	}
	if ac.Query != "" {
		t.Error("Expected Query to be empty")
	}
	if ac.Open {
		t.Error("Expected Open to be false")
	}
}

func TestFocus(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
	}
	ac := New("test", WithSuggestions(suggestions), WithMinChars(1))
	ac.Query = "o"

	ac.Focus()

	if !ac.Open {
		t.Error("Expected Open to be true after Focus with matching suggestions")
	}
}

func TestFocusBelowMinChars(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
	}
	ac := New("test", WithSuggestions(suggestions), WithMinChars(3))
	ac.Query = "o"

	ac.Focus()

	if ac.Open {
		t.Error("Expected Open to be false when query below MinChars")
	}
}

func TestBlur(t *testing.T) {
	ac := New("test")
	ac.Open = true
	ac.HighlightedIndex = 2

	ac.Blur()

	if ac.Open {
		t.Error("Expected Open to be false after Blur")
	}
	if ac.HighlightedIndex != -1 {
		t.Errorf("Expected HighlightedIndex -1, got %d", ac.HighlightedIndex)
	}
}

func TestHighlightNext(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.Filter()

	ac.HighlightNext()
	if ac.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0, got %d", ac.HighlightedIndex)
	}

	ac.HighlightNext()
	if ac.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1, got %d", ac.HighlightedIndex)
	}

	// Should wrap around
	ac.HighlightNext()
	if ac.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0 (wrapped), got %d", ac.HighlightedIndex)
	}
}

func TestHighlightNextSkipsDisabled(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two", Disabled: true},
		{Value: "3", Label: "Three"},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.Filter()
	ac.HighlightedIndex = 0

	ac.HighlightNext()

	if ac.HighlightedIndex != 2 {
		t.Errorf("Expected HighlightedIndex 2 (skipped disabled), got %d", ac.HighlightedIndex)
	}
}

func TestHighlightPrevious(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	ac := New("test", WithSuggestions(suggestions))
	ac.Filter()

	// Should start at end when going previous from -1
	ac.HighlightPrevious()
	if ac.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1, got %d", ac.HighlightedIndex)
	}

	ac.HighlightPrevious()
	if ac.HighlightedIndex != 0 {
		t.Errorf("Expected HighlightedIndex 0, got %d", ac.HighlightedIndex)
	}

	// Should wrap to end
	ac.HighlightPrevious()
	if ac.HighlightedIndex != 1 {
		t.Errorf("Expected HighlightedIndex 1 (wrapped), got %d", ac.HighlightedIndex)
	}
}

func TestIsHighlighted(t *testing.T) {
	ac := New("test")
	ac.HighlightedIndex = 2

	if !ac.IsHighlighted(2) {
		t.Error("Expected index 2 to be highlighted")
	}
	if ac.IsHighlighted(1) {
		t.Error("Expected index 1 to not be highlighted")
	}
}

func TestHasSelection(t *testing.T) {
	ac := New("test")

	if ac.HasSelection() {
		t.Error("Expected HasSelection to be false initially")
	}

	ac.Selected = &Suggestion{Value: "1", Label: "One"}
	if !ac.HasSelection() {
		t.Error("Expected HasSelection to be true")
	}
}

func TestDisplayValue(t *testing.T) {
	ac := New("test")

	// No selection - return query
	ac.Query = "search"
	if ac.DisplayValue() != "search" {
		t.Errorf("Expected DisplayValue 'search', got '%s'", ac.DisplayValue())
	}

	// With selection - return label
	ac.Selected = &Suggestion{Value: "1", Label: "Selected Item"}
	if ac.DisplayValue() != "Selected Item" {
		t.Errorf("Expected DisplayValue 'Selected Item', got '%s'", ac.DisplayValue())
	}
}

func TestSetSuggestions(t *testing.T) {
	ac := New("test")
	ac.Query = "new"

	newSuggestions := []Suggestion{
		{Value: "1", Label: "New York"},
		{Value: "2", Label: "Old City"},
	}
	ac.SetSuggestions(newSuggestions)

	if len(ac.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(ac.Suggestions))
	}
	// Should auto-filter
	if len(ac.FilteredSuggestions) != 1 {
		t.Errorf("Expected 1 filtered suggestion, got %d", len(ac.FilteredSuggestions))
	}
}

func TestSetLoading(t *testing.T) {
	ac := New("test")

	ac.SetLoading(true)
	if !ac.Loading {
		t.Error("Expected Loading to be true")
	}

	ac.SetLoading(false)
	if ac.Loading {
		t.Error("Expected Loading to be false")
	}
}

// MultiAutocomplete tests

func TestSelectMulti(t *testing.T) {
	mac := NewMulti("test")
	s := Suggestion{Value: "1", Label: "One"}

	result := mac.SelectMulti(s)

	if !result {
		t.Error("Expected SelectMulti to return true")
	}
	if len(mac.SelectedItems) != 1 {
		t.Errorf("Expected 1 selected item, got %d", len(mac.SelectedItems))
	}
	if mac.Query != "" {
		t.Error("Expected Query to be cleared")
	}
}

func TestSelectMultiDuplicate(t *testing.T) {
	mac := NewMulti("test")
	s := Suggestion{Value: "1", Label: "One"}

	mac.SelectMulti(s)
	result := mac.SelectMulti(s)

	if result {
		t.Error("Expected SelectMulti to return false for duplicate")
	}
	if len(mac.SelectedItems) != 1 {
		t.Errorf("Expected 1 selected item, got %d", len(mac.SelectedItems))
	}
}

func TestSelectMultiDisabled(t *testing.T) {
	mac := NewMulti("test")
	s := Suggestion{Value: "1", Label: "One", Disabled: true}

	if mac.SelectMulti(s) {
		t.Error("Expected SelectMulti to return false for disabled")
	}
}

func TestRemoveSelected(t *testing.T) {
	mac := NewMulti("test")
	mac.SelectedItems = []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}

	result := mac.RemoveSelected("1")

	if !result {
		t.Error("Expected RemoveSelected to return true")
	}
	if len(mac.SelectedItems) != 1 {
		t.Errorf("Expected 1 selected item, got %d", len(mac.SelectedItems))
	}
	if mac.SelectedItems[0].Value != "2" {
		t.Error("Expected remaining item to be '2'")
	}
}

func TestRemoveSelectedNotFound(t *testing.T) {
	mac := NewMulti("test")
	mac.SelectedItems = []Suggestion{
		{Value: "1", Label: "One"},
	}

	if mac.RemoveSelected("99") {
		t.Error("Expected RemoveSelected to return false for not found")
	}
}

func TestClearMulti(t *testing.T) {
	mac := NewMulti("test")
	mac.SelectedItems = []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}
	mac.Query = "search"
	mac.Open = true

	mac.ClearMulti()

	if mac.SelectedItems != nil {
		t.Error("Expected SelectedItems to be nil")
	}
	if mac.Query != "" {
		t.Error("Expected Query to be empty")
	}
	if mac.Open {
		t.Error("Expected Open to be false")
	}
}

func TestIsSelectedMulti(t *testing.T) {
	mac := NewMulti("test")
	mac.SelectedItems = []Suggestion{
		{Value: "1", Label: "One"},
	}

	if !mac.IsSelectedMulti("1") {
		t.Error("Expected '1' to be selected")
	}
	if mac.IsSelectedMulti("2") {
		t.Error("Expected '2' to not be selected")
	}
}

func TestHasSelectedItems(t *testing.T) {
	mac := NewMulti("test")

	if mac.HasSelectedItems() {
		t.Error("Expected HasSelectedItems to be false initially")
	}

	mac.SelectedItems = []Suggestion{{Value: "1", Label: "One"}}
	if !mac.HasSelectedItems() {
		t.Error("Expected HasSelectedItems to be true")
	}
}

func TestSelectedValues(t *testing.T) {
	mac := NewMulti("test")
	mac.SelectedItems = []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
	}

	values := mac.SelectedValues()

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
	if values[0] != "1" || values[1] != "2" {
		t.Errorf("Expected ['1', '2'], got %v", values)
	}
}

func TestFilteredExcludingSelected(t *testing.T) {
	suggestions := []Suggestion{
		{Value: "1", Label: "One"},
		{Value: "2", Label: "Two"},
		{Value: "3", Label: "Three"},
	}
	mac := NewMulti("test", WithSuggestions(suggestions))
	mac.Filter()
	mac.SelectedItems = []Suggestion{
		{Value: "2", Label: "Two"},
	}

	filtered := mac.FilteredExcludingSelected()

	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered suggestions, got %d", len(filtered))
	}
	for _, s := range filtered {
		if s.Value == "2" {
			t.Error("Expected '2' to be excluded from filtered")
		}
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
