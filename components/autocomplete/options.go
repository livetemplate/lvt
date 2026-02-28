package autocomplete

// Option is a functional option for configuring autocomplete components.
type Option func(*Autocomplete)

// WithPlaceholder sets the placeholder text.
func WithPlaceholder(placeholder string) Option {
	return func(ac *Autocomplete) {
		ac.Placeholder = placeholder
	}
}

// WithSuggestions sets the initial suggestions.
func WithSuggestions(suggestions []Suggestion) Option {
	return func(ac *Autocomplete) {
		ac.Suggestions = suggestions
	}
}

// WithMinChars sets the minimum characters before showing suggestions.
func WithMinChars(min int) Option {
	return func(ac *Autocomplete) {
		if min < 0 {
			min = 0
		}
		ac.MinChars = min
	}
}

// WithMaxSuggestions sets the maximum number of shown suggestions.
func WithMaxSuggestions(max int) Option {
	return func(ac *Autocomplete) {
		ac.MaxSuggestions = max
	}
}

// WithAllowCustom allows values not in the suggestions list.
func WithAllowCustom(allow bool) Option {
	return func(ac *Autocomplete) {
		ac.AllowCustom = allow
	}
}

// WithClearOnSelect clears the input after selection.
func WithClearOnSelect(clear bool) Option {
	return func(ac *Autocomplete) {
		ac.ClearOnSelect = clear
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) Option {
	return func(ac *Autocomplete) {
		ac.SetStyled(styled)
	}
}

// WithFilterFunc sets a custom filter function.
func WithFilterFunc(fn func(query string, suggestions []Suggestion) []Suggestion) Option {
	return func(ac *Autocomplete) {
		ac.filterFunc = fn
	}
}

// WithQuery sets the initial query.
func WithQuery(query string) Option {
	return func(ac *Autocomplete) {
		ac.Query = query
	}
}

// WithSelected sets the initially selected suggestion.
func WithSelected(s Suggestion) Option {
	return func(ac *Autocomplete) {
		ac.Selected = &s
		ac.Query = s.Label
	}
}
