package dropdown

// Option is a functional option for configuring dropdowns.
type Option func(*Dropdown)

// WithPlaceholder sets the placeholder text shown when nothing is selected.
func WithPlaceholder(placeholder string) Option {
	return func(d *Dropdown) {
		d.Placeholder = placeholder
	}
}

// WithSelected pre-selects an item by value.
func WithSelected(value string) Option {
	return func(d *Dropdown) {
		for i := range d.Options {
			if d.Options[i].Value == value {
				d.Selected = &d.Options[i]
				return
			}
		}
	}
}

// WithDisabled disables the dropdown.
func WithDisabled(disabled bool) Option {
	return func(d *Dropdown) {
		d.Disabled = disabled
	}
}

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(d *Dropdown) {
		d.Open = open
	}
}

// WithStyled enables Tailwind CSS styling for the component.
// When false, renders semantic HTML without styling classes.
func WithStyled(styled bool) Option {
	return func(d *Dropdown) {
		d.SetStyled(styled)
	}
}

// SearchableOption is a functional option for configuring searchable dropdowns.
type SearchableOption func(*Searchable)

// WithMinChars sets the minimum characters required before filtering starts.
func WithMinChars(minChars int) SearchableOption {
	return func(s *Searchable) {
		s.MinChars = minChars
	}
}

// WithQuery sets the initial search query.
func WithQuery(query string) SearchableOption {
	return func(s *Searchable) {
		s.Query = query
	}
}

// MultiOption is a functional option for configuring multi-select dropdowns.
type MultiOption func(*Multi)

// WithMaxSelections limits how many items can be selected (0 = unlimited).
func WithMaxSelections(max int) MultiOption {
	return func(m *Multi) {
		m.MaxSelections = max
	}
}

// WithSelectedValues pre-selects multiple items by their values.
func WithSelectedValues(values []string) MultiOption {
	return func(m *Multi) {
		valueSet := make(map[string]bool)
		for _, v := range values {
			valueSet[v] = true
		}

		m.SelectedItems = make([]Item, 0)
		for _, opt := range m.Options {
			if valueSet[opt.Value] {
				m.SelectedItems = append(m.SelectedItems, opt)
			}
		}
	}
}
