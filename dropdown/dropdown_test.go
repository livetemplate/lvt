package dropdown

import (
	"html/template"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	options := []Item{
		{Value: "us", Label: "United States"},
		{Value: "ca", Label: "Canada"},
		{Value: "mx", Label: "Mexico"},
	}

	d := New("country", options)

	if d.ID() != "country" {
		t.Errorf("expected ID 'country', got '%s'", d.ID())
	}
	if len(d.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(d.Options))
	}
	if d.Placeholder != "Select..." {
		t.Errorf("expected placeholder 'Select...', got '%s'", d.Placeholder)
	}
	if d.Selected != nil {
		t.Error("expected Selected to be nil")
	}
	if d.Open {
		t.Error("expected Open to be false")
	}
	if d.Disabled {
		t.Error("expected Disabled to be false")
	}
}

func TestNewWithOptions(t *testing.T) {
	options := []Item{
		{Value: "us", Label: "United States"},
		{Value: "ca", Label: "Canada"},
	}

	d := New("country", options,
		WithPlaceholder("Choose a country"),
		WithSelected("ca"),
		WithDisabled(true),
		WithOpen(true),
	)

	if d.Placeholder != "Choose a country" {
		t.Errorf("expected placeholder 'Choose a country', got '%s'", d.Placeholder)
	}
	if d.Selected == nil || d.Selected.Value != "ca" {
		t.Error("expected 'ca' to be selected")
	}
	if !d.Disabled {
		t.Error("expected Disabled to be true")
	}
	if !d.Open {
		t.Error("expected Open to be true")
	}
}

func TestDropdown_Toggle(t *testing.T) {
	d := New("test", nil)

	if d.Open {
		t.Error("expected initially closed")
	}

	d.Toggle()
	if !d.Open {
		t.Error("expected open after toggle")
	}

	d.Toggle()
	if d.Open {
		t.Error("expected closed after second toggle")
	}
}

func TestDropdown_Close(t *testing.T) {
	d := New("test", nil, WithOpen(true))

	d.Close()
	if d.Open {
		t.Error("expected closed after Close()")
	}
}

func TestDropdown_Select(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
	}

	d := New("test", options, WithOpen(true))

	d.Select("b")

	if d.Selected == nil || d.Selected.Value != "b" {
		t.Error("expected 'b' to be selected")
	}
	if d.Open {
		t.Error("expected dropdown to close after selection")
	}
}

func TestDropdown_SelectNonExistent(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Option A"},
	}

	d := New("test", options, WithOpen(true))
	d.Select("nonexistent")

	if d.Selected != nil {
		t.Error("expected Selected to remain nil for nonexistent value")
	}
}

func TestDropdown_Clear(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Option A"},
	}

	d := New("test", options, WithSelected("a"))
	d.Clear()

	if d.Selected != nil {
		t.Error("expected Selected to be nil after Clear()")
	}
}

func TestDropdown_Value(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Option A"},
	}

	d := New("test", options)

	if d.Value() != "" {
		t.Error("expected empty value when nothing selected")
	}

	d.Select("a")
	if d.Value() != "a" {
		t.Errorf("expected value 'a', got '%s'", d.Value())
	}
}

func TestNewSearchable(t *testing.T) {
	options := []Item{
		{Value: "us", Label: "United States"},
		{Value: "ca", Label: "Canada"},
	}

	s := NewSearchable("country", options)

	if s.ID() != "country" {
		t.Errorf("expected ID 'country', got '%s'", s.ID())
	}
	if s.Placeholder != "Search..." {
		t.Errorf("expected placeholder 'Search...', got '%s'", s.Placeholder)
	}
	if s.MinChars != 1 {
		t.Errorf("expected MinChars 1, got %d", s.MinChars)
	}
	if s.Query != "" {
		t.Error("expected empty query")
	}
}

func TestSearchable_Search(t *testing.T) {
	options := []Item{
		{Value: "us", Label: "United States"},
		{Value: "ca", Label: "Canada"},
		{Value: "mx", Label: "Mexico"},
	}

	s := NewSearchable("country", options)

	s.Search("can")

	if s.Query != "can" {
		t.Errorf("expected query 'can', got '%s'", s.Query)
	}
	if !s.Open {
		t.Error("expected Open to be true after search")
	}
	if len(s.FilteredOptions) != 1 {
		t.Errorf("expected 1 filtered option, got %d", len(s.FilteredOptions))
	}
	if s.FilteredOptions[0].Value != "ca" {
		t.Error("expected 'ca' in filtered options")
	}
}

func TestSearchable_SearchCaseInsensitive(t *testing.T) {
	options := []Item{
		{Value: "us", Label: "United States"},
		{Value: "ca", Label: "Canada"},
	}

	s := NewSearchable("country", options)
	s.Search("UNITED")

	if len(s.FilteredOptions) != 1 || s.FilteredOptions[0].Value != "us" {
		t.Error("expected case-insensitive match for 'United States'")
	}
}

func TestSearchable_SearchBelowMinChars(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
	}

	s := NewSearchable("test", options)
	s.MinChars = 3

	s.Search("al")

	if s.FilteredOptions != nil {
		t.Error("expected FilteredOptions to be nil when below MinChars")
	}
}

func TestSearchable_VisibleOptions(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	s := NewSearchable("test", options)

	// No query - should return all options
	visible := s.VisibleOptions()
	if len(visible) != 2 {
		t.Errorf("expected 2 visible options without query, got %d", len(visible))
	}

	// With query - should return filtered
	s.Search("alp")
	visible = s.VisibleOptions()
	if len(visible) != 1 || visible[0].Value != "a" {
		t.Error("expected filtered options after search")
	}
}

func TestSearchable_ClearSearch(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
	}

	s := NewSearchable("test", options)
	s.Search("alp")
	s.ClearSearch()

	if s.Query != "" {
		t.Error("expected empty query after ClearSearch")
	}
	if s.FilteredOptions != nil {
		t.Error("expected FilteredOptions to be nil after ClearSearch")
	}
}

func TestNewMulti(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	m := NewMulti("tags", options)

	if m.ID() != "tags" {
		t.Errorf("expected ID 'tags', got '%s'", m.ID())
	}
	if len(m.SelectedItems) != 0 {
		t.Error("expected empty SelectedItems")
	}
	if m.MaxSelections != 0 {
		t.Error("expected MaxSelections 0 (unlimited)")
	}
}

func TestMulti_ToggleItem(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	m := NewMulti("tags", options)

	// Select first item
	m.ToggleItem("a")
	if len(m.SelectedItems) != 1 || m.SelectedItems[0].Value != "a" {
		t.Error("expected 'a' to be selected")
	}

	// Select second item
	m.ToggleItem("b")
	if len(m.SelectedItems) != 2 {
		t.Error("expected both items to be selected")
	}

	// Deselect first item
	m.ToggleItem("a")
	if len(m.SelectedItems) != 1 || m.SelectedItems[0].Value != "b" {
		t.Error("expected only 'b' to remain selected")
	}
}

func TestMulti_ToggleItemMaxSelections(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
		{Value: "c", Label: "Charlie"},
	}

	m := NewMulti("tags", options)
	m.MaxSelections = 2

	m.ToggleItem("a")
	m.ToggleItem("b")
	m.ToggleItem("c") // Should be ignored due to max selections

	if len(m.SelectedItems) != 2 {
		t.Errorf("expected 2 selected items (max), got %d", len(m.SelectedItems))
	}
	if m.IsSelected("c") {
		t.Error("expected 'c' to not be selected due to max selections")
	}
}

func TestMulti_IsSelected(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	m := NewMulti("tags", options)
	m.ToggleItem("a")

	if !m.IsSelected("a") {
		t.Error("expected 'a' to be selected")
	}
	if m.IsSelected("b") {
		t.Error("expected 'b' to not be selected")
	}
}

func TestMulti_Values(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	m := NewMulti("tags", options)
	m.ToggleItem("a")
	m.ToggleItem("b")

	values := m.Values()
	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}
}

func TestMulti_ClearAll(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
	}

	m := NewMulti("tags", options)
	m.ToggleItem("a")
	m.ToggleItem("b")
	m.ClearAll()

	if len(m.SelectedItems) != 0 {
		t.Error("expected empty SelectedItems after ClearAll")
	}
}

func TestMulti_SelectAll(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta", Disabled: true},
		{Value: "c", Label: "Charlie"},
	}

	m := NewMulti("tags", options)
	m.SelectAll()

	// Should select all except disabled
	if len(m.SelectedItems) != 2 {
		t.Errorf("expected 2 selected items, got %d", len(m.SelectedItems))
	}
	if m.IsSelected("b") {
		t.Error("expected disabled item 'b' to not be selected")
	}
}

func TestMulti_SelectAllWithMax(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
		{Value: "c", Label: "Charlie"},
	}

	m := NewMulti("tags", options)
	m.MaxSelections = 2
	m.SelectAll()

	if len(m.SelectedItems) != 2 {
		t.Errorf("expected 2 selected items (max), got %d", len(m.SelectedItems))
	}
}

func TestMulti_DisplayText(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
		{Value: "c", Label: "Charlie"},
	}

	m := NewMulti("tags", options, WithPlaceholder("Select tags"))

	// No selection
	if m.DisplayText() != "Select tags" {
		t.Errorf("expected placeholder, got '%s'", m.DisplayText())
	}

	// One selection
	m.ToggleItem("a")
	if m.DisplayText() != "Alpha" {
		t.Errorf("expected 'Alpha', got '%s'", m.DisplayText())
	}

	// Multiple selections
	m.ToggleItem("b")
	m.ToggleItem("c")
	text := m.DisplayText()
	if !strings.Contains(text, "Alpha") || !strings.Contains(text, "+ 2 more") {
		t.Errorf("expected 'Alpha + 2 more', got '%s'", text)
	}
}

func TestWithSelectedValues(t *testing.T) {
	options := []Item{
		{Value: "a", Label: "Alpha"},
		{Value: "b", Label: "Beta"},
		{Value: "c", Label: "Charlie"},
	}

	m := NewMulti("tags", options)
	WithSelectedValues([]string{"a", "c"})(m)

	if len(m.SelectedItems) != 2 {
		t.Errorf("expected 2 selected items, got %d", len(m.SelectedItems))
	}
	if !m.IsSelected("a") || !m.IsSelected("c") {
		t.Error("expected 'a' and 'c' to be selected")
	}
	if m.IsSelected("b") {
		t.Error("expected 'b' to not be selected")
	}
}

func TestWithMinChars(t *testing.T) {
	options := []Item{{Value: "a", Label: "Alpha"}}

	s := NewSearchable("test", options)
	WithMinChars(3)(s)

	if s.MinChars != 3 {
		t.Errorf("expected MinChars 3, got %d", s.MinChars)
	}
}

func TestWithQuery(t *testing.T) {
	options := []Item{{Value: "a", Label: "Alpha"}}

	s := NewSearchable("test", options)
	WithQuery("initial")(s)

	if s.Query != "initial" {
		t.Errorf("expected query 'initial', got '%s'", s.Query)
	}
}

func TestWithMaxSelections(t *testing.T) {
	options := []Item{{Value: "a", Label: "Alpha"}}

	m := NewMulti("test", options)
	WithMaxSelections(5)(m)

	if m.MaxSelections != 5 {
		t.Errorf("expected MaxSelections 5, got %d", m.MaxSelections)
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()

	if ts == nil {
		t.Fatal("expected non-nil TemplateSet")
	}
	if ts.Namespace != "dropdown" {
		t.Errorf("expected namespace 'dropdown', got '%s'", ts.Namespace)
	}
}

func TestTemplateRendering(t *testing.T) {
	// Parse the templates
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	// Test default template
	t.Run("default template", func(t *testing.T) {
		options := []Item{
			{Value: "a", Label: "Alpha"},
			{Value: "b", Label: "Beta"},
		}
		d := New("test", options, WithPlaceholder("Select..."), WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:dropdown:default:v1", d)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-dropdown="test"`) {
			t.Error("expected data-dropdown attribute")
		}
		if !strings.Contains(html, "Select...") {
			t.Error("expected placeholder text")
		}
		if !strings.Contains(html, `lvt-click="toggle_test"`) {
			t.Error("expected lvt-click attribute")
		}
	})

	// Test searchable template
	t.Run("searchable template", func(t *testing.T) {
		options := []Item{
			{Value: "a", Label: "Alpha"},
		}
		s := NewSearchable("search-test", options, WithStyled(true))
		s.Open = true

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:dropdown:searchable:v1", s)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `lvt-input="search_search-test"`) {
			t.Error("expected lvt-input attribute")
		}
		if !strings.Contains(html, `lvt-debounce="150"`) {
			t.Error("expected lvt-debounce attribute")
		}
	})

	// Test multi template
	t.Run("multi template", func(t *testing.T) {
		options := []Item{
			{Value: "a", Label: "Alpha"},
			{Value: "b", Label: "Beta"},
		}
		m := NewMulti("multi-test", options, WithStyled(true))
		m.Open = true
		m.ToggleItem("a")

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:dropdown:multi:v1", m)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `aria-multiselectable="true"`) {
			t.Error("expected aria-multiselectable attribute")
		}
		if !strings.Contains(html, `lvt-change="toggle_item_multi-test"`) {
			t.Error("expected lvt-change attribute")
		}
		if !strings.Contains(html, "1 selected") {
			t.Error("expected selected count")
		}
	})
}

func TestUnstyledTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	options := []Item{
		{Value: "a", Label: "Alpha"},
	}

	// Default template unstyled
	t.Run("default unstyled", func(t *testing.T) {
		d := New("test", options, WithStyled(false))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:dropdown:default:v1", d)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		// Unstyled version should not have Tailwind classes
		if strings.Contains(html, "class=\"relative") {
			t.Error("unstyled template should not have Tailwind classes")
		}
	})
}

// Helper tests
func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"already lower", "already lower"},
		{"MiXeD", "mixed"},
		{"", ""},
	}

	for _, tc := range tests {
		result := toLower(tc.input)
		if result != tc.expected {
			t.Errorf("toLower(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"hello", "", true},
		{"", "hello", false},
		{"abc", "abcd", false},
	}

	for _, tc := range tests {
		result := contains(tc.s, tc.substr)
		if result != tc.expected {
			t.Errorf("contains(%q, %q) = %v, expected %v", tc.s, tc.substr, result, tc.expected)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{100, "100"},
		{-5, "-5"},
		{-123, "-123"},
	}

	for _, tc := range tests {
		result := itoa(tc.input)
		if result != tc.expected {
			t.Errorf("itoa(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
