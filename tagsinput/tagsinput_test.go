package tagsinput

import (
	"html/template"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ti := New("skills")

	if ti.ID() != "skills" {
		t.Errorf("expected ID 'skills', got '%s'", ti.ID())
	}
	if ti.Count() != 0 {
		t.Error("expected 0 tags initially")
	}
	if ti.Placeholder != "Add tag..." {
		t.Errorf("expected default placeholder, got '%s'", ti.Placeholder)
	}
}

func TestNewWithOptions(t *testing.T) {
	ti := New("test",
		WithPlaceholder("Add skill..."),
		WithTags("go", "python"),
		WithMaxTags(5),
		WithStyled(false),
	)

	if ti.Placeholder != "Add skill..." {
		t.Errorf("expected placeholder 'Add skill...', got '%s'", ti.Placeholder)
	}
	if ti.Count() != 2 {
		t.Errorf("expected 2 tags, got %d", ti.Count())
	}
	if ti.MaxTags != 5 {
		t.Errorf("expected MaxTags 5, got %d", ti.MaxTags)
	}
	if ti.IsStyled() {
		t.Error("expected IsStyled to be false")
	}
}

func TestTagsInput_AddTag(t *testing.T) {
	ti := New("test")

	if !ti.AddTag("go") {
		t.Error("expected AddTag to return true")
	}
	if ti.Count() != 1 {
		t.Error("expected 1 tag after AddTag")
	}
	if ti.Tags[0].Value != "go" {
		t.Errorf("expected tag value 'go', got '%s'", ti.Tags[0].Value)
	}
}

func TestTagsInput_AddTagTrimsSpace(t *testing.T) {
	ti := New("test")

	ti.AddTag("  go  ")

	if ti.Tags[0].Value != "go" {
		t.Errorf("expected trimmed tag value 'go', got '%s'", ti.Tags[0].Value)
	}
}

func TestTagsInput_AddTagEmpty(t *testing.T) {
	ti := New("test")

	if ti.AddTag("") {
		t.Error("expected AddTag to return false for empty value")
	}
	if ti.AddTag("   ") {
		t.Error("expected AddTag to return false for whitespace-only value")
	}
}

func TestTagsInput_AddTagMaxTags(t *testing.T) {
	ti := New("test", WithMaxTags(2))

	ti.AddTag("a")
	ti.AddTag("b")
	result := ti.AddTag("c")

	if result {
		t.Error("expected AddTag to return false when max reached")
	}
	if ti.Count() != 2 {
		t.Errorf("expected 2 tags (max), got %d", ti.Count())
	}
}

func TestTagsInput_AddTagNoDuplicates(t *testing.T) {
	ti := New("test")

	ti.AddTag("go")
	result := ti.AddTag("go")

	if result {
		t.Error("expected AddTag to return false for duplicate")
	}
	if ti.Count() != 1 {
		t.Error("expected 1 tag (no duplicates)")
	}
}

func TestTagsInput_AddTagAllowDuplicates(t *testing.T) {
	ti := New("test", WithAllowDuplicates(true))

	ti.AddTag("go")
	ti.AddTag("go")

	if ti.Count() != 2 {
		t.Error("expected 2 tags (duplicates allowed)")
	}
}

func TestTagsInput_RemoveTag(t *testing.T) {
	ti := New("test", WithTags("go", "python", "rust"))

	ti.RemoveTag("python")

	if ti.Count() != 2 {
		t.Error("expected 2 tags after remove")
	}
	if ti.HasTag("python") {
		t.Error("expected python to be removed")
	}
}

func TestTagsInput_RemoveTagAt(t *testing.T) {
	ti := New("test", WithTags("a", "b", "c"))

	ti.RemoveTagAt(1)

	if ti.Count() != 2 {
		t.Error("expected 2 tags after RemoveTagAt")
	}
	if ti.HasTag("b") {
		t.Error("expected 'b' to be removed")
	}
}

func TestTagsInput_RemoveTagAtInvalid(t *testing.T) {
	ti := New("test", WithTags("a"))

	ti.RemoveTagAt(-1)
	ti.RemoveTagAt(10)

	if ti.Count() != 1 {
		t.Error("expected tags to remain unchanged for invalid index")
	}
}

func TestTagsInput_RemoveLast(t *testing.T) {
	ti := New("test", WithTags("a", "b", "c"))

	ti.RemoveLast()

	if ti.Count() != 2 {
		t.Error("expected 2 tags after RemoveLast")
	}
	if ti.HasTag("c") {
		t.Error("expected 'c' to be removed")
	}
}

func TestTagsInput_RemoveLastEmpty(t *testing.T) {
	ti := New("test")

	ti.RemoveLast() // Should not panic

	if ti.Count() != 0 {
		t.Error("expected 0 tags")
	}
}

func TestTagsInput_HasTag(t *testing.T) {
	ti := New("test", WithTags("go", "python"))

	if !ti.HasTag("go") {
		t.Error("expected HasTag('go') to be true")
	}
	if ti.HasTag("rust") {
		t.Error("expected HasTag('rust') to be false")
	}
}

func TestTagsInput_Clear(t *testing.T) {
	ti := New("test", WithTags("a", "b", "c"))
	ti.Input = "test"

	ti.Clear()

	if ti.Count() != 0 {
		t.Error("expected 0 tags after Clear")
	}
	if ti.Input != "" {
		t.Error("expected empty input after Clear")
	}
}

func TestTagsInput_SetInput(t *testing.T) {
	ti := New("test")

	ti.SetInput("go")

	if ti.Input != "go" {
		t.Errorf("expected input 'go', got '%s'", ti.Input)
	}
}

func TestTagsInput_SetInputWithSeparator(t *testing.T) {
	ti := New("test", WithSeparators(","))

	ti.SetInput("go,python,rust")

	// Should have added tags
	if ti.Count() != 3 {
		t.Errorf("expected 3 tags after separator input, got %d", ti.Count())
	}
	if ti.Input != "" {
		t.Error("expected input to be cleared after adding tags")
	}
}

func TestTagsInput_Values(t *testing.T) {
	ti := New("test", WithTags("go", "python"))

	values := ti.Values()

	if len(values) != 2 {
		t.Error("expected 2 values")
	}
	if values[0] != "go" || values[1] != "python" {
		t.Error("expected values to match tags")
	}
}

func TestTagsInput_IsEmpty(t *testing.T) {
	ti := New("test")

	if !ti.IsEmpty() {
		t.Error("expected IsEmpty to be true initially")
	}

	ti.AddTag("go")

	if ti.IsEmpty() {
		t.Error("expected IsEmpty to be false after AddTag")
	}
}

func TestTagsInput_CanAddMore(t *testing.T) {
	ti := New("test", WithMaxTags(2))

	if !ti.CanAddMore() {
		t.Error("expected CanAddMore to be true initially")
	}

	ti.AddTag("a")
	ti.AddTag("b")

	if ti.CanAddMore() {
		t.Error("expected CanAddMore to be false at max")
	}
}

func TestTagsInput_CanAddMoreUnlimited(t *testing.T) {
	ti := New("test") // MaxTags = 0 (unlimited)

	for i := 0; i < 100; i++ {
		if !ti.CanAddMore() {
			t.Error("expected CanAddMore to always be true when unlimited")
		}
		ti.AddTag("tag" + string(rune('0'+i%10)))
	}
}

func TestTagsInput_FilteredSuggestions(t *testing.T) {
	ti := New("test", WithSuggestions("golang", "python", "rust", "go"))

	ti.Input = "go"
	suggestions := ti.FilteredSuggestions()

	if len(suggestions) != 2 { // "golang" and "go"
		t.Errorf("expected 2 suggestions, got %d", len(suggestions))
	}
}

func TestTagsInput_FilteredSuggestionsExcludesExisting(t *testing.T) {
	ti := New("test",
		WithSuggestions("golang", "python", "rust"),
		WithTags("golang"),
	)

	ti.Input = "go"
	suggestions := ti.FilteredSuggestions()

	// Should not include "golang" since it's already a tag
	for _, s := range suggestions {
		if s == "golang" {
			t.Error("filtered suggestions should exclude existing tags")
		}
	}
}

func TestTagsInput_FilteredSuggestionsEmpty(t *testing.T) {
	ti := New("test", WithSuggestions("golang", "python"))

	ti.Input = ""
	suggestions := ti.FilteredSuggestions()

	if suggestions != nil {
		t.Error("expected nil suggestions for empty input")
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()

	if ts == nil {
		t.Fatal("expected non-nil TemplateSet")
	}
	if ts.Namespace != "tagsinput" {
		t.Errorf("expected namespace 'tagsinput', got '%s'", ts.Namespace)
	}
}

func TestTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	ti := New("skills",
		WithTags("go", "python"),
		WithPlaceholder("Add skill..."),
		WithStyled(true),
	)

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:tagsinput:default:v1", ti)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `data-tagsinput="skills"`) {
		t.Error("expected data-tagsinput attribute")
	}
	if !strings.Contains(html, "go") {
		t.Error("expected 'go' tag in output")
	}
	if !strings.Contains(html, "python") {
		t.Error("expected 'python' tag in output")
	}
	if !strings.Contains(html, `lvt-input="input_tag_skills"`) {
		t.Error("expected lvt-input attribute")
	}
	if !strings.Contains(html, `lvt-click="remove_tag_skills"`) {
		t.Error("expected remove button with lvt-click")
	}
}

func TestUnstyledTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	ti := New("test", WithTags("go"), WithStyled(false))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:tagsinput:default:v1", ti)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Unstyled version should not have Tailwind classes
	if strings.Contains(html, "flex-wrap") {
		t.Error("unstyled template should not have Tailwind classes")
	}
}

func TestMaxTagsDisplay(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	ti := New("test", WithMaxTags(5), WithTags("a", "b"), WithStyled(true))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:tagsinput:default:v1", ti)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "2/5") {
		t.Error("expected '2/5 tags' counter in output")
	}
}

// Helper function tests
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\ntest\r\n", "test"},
		{"no trim", "no trim"},
		{"", ""},
	}

	for _, tc := range tests {
		result := trimSpace(tc.input)
		if result != tc.expected {
			t.Errorf("trimSpace(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		input    string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"a", ",", []string{"a"}},
		{"", ",", []string{""}},
	}

	for _, tc := range tests {
		result := split(tc.input, tc.sep)
		if len(result) != len(tc.expected) {
			t.Errorf("split(%q, %q) length = %d, expected %d", tc.input, tc.sep, len(result), len(tc.expected))
		}
	}
}
