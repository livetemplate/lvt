// Package tagsinput provides a tag/chip input component for entering multiple values.
//
// Available variants:
//   - New() creates a basic tags input (template: "lvt:tagsinput:default:v1")
//
// Required lvt-* attributes: lvt-input, lvt-keydown, lvt-click
//
// Example usage:
//
//	// In your controller/state
//	Tags: tagsinput.New("skills",
//	    tagsinput.WithPlaceholder("Add skills..."),
//	    tagsinput.WithTags("go", "typescript", "python"),
//	)
//
//	// In your template
//	{{template "lvt:tagsinput:default:v1" .Tags}}
package tagsinput

import (
	"github.com/livetemplate/components/base"
)

// Tag represents a single tag/chip.
type Tag struct {
	Value string // The tag value
	Label string // Display text (defaults to Value if empty)
}

// TagsInput is a component for entering multiple tags.
// Use template "lvt:tagsinput:default:v1" to render.
type TagsInput struct {
	base.Base

	// Tags is the list of current tags
	Tags []Tag

	// Input is the current input value (before becoming a tag)
	Input string

	// Placeholder text shown when empty
	Placeholder string

	// MaxTags limits the number of tags (0 = unlimited)
	MaxTags int

	// AllowDuplicates allows the same tag value multiple times
	AllowDuplicates bool

	// Separator characters that trigger tag creation (default: comma, Enter)
	Separators []string

	// Suggestions for autocomplete (optional)
	Suggestions []string

	// ShowSuggestions controls visibility of suggestion dropdown
	ShowSuggestions bool
}

// New creates a new tags input.
//
// Example:
//
//	tags := tagsinput.New("skills",
//	    tagsinput.WithPlaceholder("Add skills..."),
//	    tagsinput.WithTags("go", "python"),
//	    tagsinput.WithMaxTags(10),
//	)
func New(id string, opts ...Option) *TagsInput {
	t := &TagsInput{
		Base:        base.NewBase(id, "tagsinput"),
		Tags:        make([]Tag, 0),
		Placeholder: "Add tag...",
		Separators:  []string{","},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// AddTag adds a new tag if valid.
func (t *TagsInput) AddTag(value string) bool {
	value = trimSpace(value)
	if value == "" {
		return false
	}

	// Check max tags
	if t.MaxTags > 0 && len(t.Tags) >= t.MaxTags {
		return false
	}

	// Check duplicates
	if !t.AllowDuplicates && t.HasTag(value) {
		return false
	}

	t.Tags = append(t.Tags, Tag{Value: value, Label: value})
	t.Input = ""
	return true
}

// RemoveTag removes a tag by value.
func (t *TagsInput) RemoveTag(value string) {
	for i, tag := range t.Tags {
		if tag.Value == value {
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
			return
		}
	}
}

// RemoveTagAt removes a tag by index.
func (t *TagsInput) RemoveTagAt(index int) {
	if index >= 0 && index < len(t.Tags) {
		t.Tags = append(t.Tags[:index], t.Tags[index+1:]...)
	}
}

// RemoveLast removes the last tag (for backspace behavior).
func (t *TagsInput) RemoveLast() {
	if len(t.Tags) > 0 {
		t.Tags = t.Tags[:len(t.Tags)-1]
	}
}

// HasTag checks if a tag value already exists.
func (t *TagsInput) HasTag(value string) bool {
	for _, tag := range t.Tags {
		if tag.Value == value {
			return true
		}
	}
	return false
}

// Clear removes all tags.
func (t *TagsInput) Clear() {
	t.Tags = make([]Tag, 0)
	t.Input = ""
}

// SetInput updates the current input value.
func (t *TagsInput) SetInput(value string) {
	t.Input = value

	// Check for separator characters
	for _, sep := range t.Separators {
		if contains(value, sep) {
			// Split by separator and add tags
			parts := split(value, sep)
			for _, part := range parts {
				t.AddTag(part)
			}
			return
		}
	}
}

// Values returns all tag values as a slice.
func (t *TagsInput) Values() []string {
	values := make([]string, len(t.Tags))
	for i, tag := range t.Tags {
		values[i] = tag.Value
	}
	return values
}

// Count returns the number of tags.
func (t *TagsInput) Count() int {
	return len(t.Tags)
}

// IsEmpty returns true if there are no tags.
func (t *TagsInput) IsEmpty() bool {
	return len(t.Tags) == 0
}

// CanAddMore returns true if more tags can be added.
func (t *TagsInput) CanAddMore() bool {
	return t.MaxTags == 0 || len(t.Tags) < t.MaxTags
}

// FilteredSuggestions returns suggestions that match the current input
// and aren't already tags.
func (t *TagsInput) FilteredSuggestions() []string {
	if t.Input == "" || len(t.Suggestions) == 0 {
		return nil
	}

	inputLower := toLower(t.Input)
	filtered := make([]string, 0)

	for _, s := range t.Suggestions {
		// Skip if already a tag
		if t.HasTag(s) {
			continue
		}
		// Match if input is prefix
		if len(inputLower) <= len(s) && toLower(s[:len(inputLower)]) == inputLower {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// Helper functions
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
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

func split(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}

	var result []string
	for {
		idx := -1
		for i := 0; i <= len(s)-len(sep); i++ {
			if s[i:i+len(sep)] == sep {
				idx = i
				break
			}
		}
		if idx < 0 {
			result = append(result, s)
			break
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	return result
}

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
