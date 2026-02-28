package tagsinput

// Option is a functional option for configuring tags input.
type Option func(*TagsInput)

// WithPlaceholder sets the placeholder text.
func WithPlaceholder(placeholder string) Option {
	return func(t *TagsInput) {
		t.Placeholder = placeholder
	}
}

// WithTags sets initial tags from string values.
func WithTags(values ...string) Option {
	return func(t *TagsInput) {
		for _, v := range values {
			t.Tags = append(t.Tags, Tag{Value: v, Label: v})
		}
	}
}

// WithMaxTags limits the number of tags (0 = unlimited).
func WithMaxTags(max int) Option {
	return func(t *TagsInput) {
		t.MaxTags = max
	}
}

// WithAllowDuplicates allows duplicate tag values.
func WithAllowDuplicates(allow bool) Option {
	return func(t *TagsInput) {
		t.AllowDuplicates = allow
	}
}

// WithSeparators sets the characters that trigger tag creation.
func WithSeparators(seps ...string) Option {
	return func(t *TagsInput) {
		t.Separators = seps
	}
}

// WithSuggestions sets autocomplete suggestions.
func WithSuggestions(suggestions ...string) Option {
	return func(t *TagsInput) {
		t.Suggestions = suggestions
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) Option {
	return func(t *TagsInput) {
		t.SetStyled(styled)
	}
}
