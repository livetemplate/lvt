package toggle

// Option is a functional option for configuring toggles.
type Option func(*Toggle)

// WithChecked sets the initial checked state.
func WithChecked(checked bool) Option {
	return func(t *Toggle) {
		t.Checked = checked
	}
}

// WithDisabled sets the disabled state.
func WithDisabled(disabled bool) Option {
	return func(t *Toggle) {
		t.Disabled = disabled
	}
}

// WithLabel sets the toggle label.
func WithLabel(label string) Option {
	return func(t *Toggle) {
		t.Label = label
	}
}

// WithLabelPosition sets where the label appears.
func WithLabelPosition(position string) Option {
	return func(t *Toggle) {
		t.LabelPosition = position
	}
}

// WithSize sets the toggle size.
func WithSize(size Size) Option {
	return func(t *Toggle) {
		t.Size = size
	}
}

// WithName sets the form field name.
func WithName(name string) Option {
	return func(t *Toggle) {
		t.Name = name
	}
}

// WithValue sets the value when checked.
func WithValue(value string) Option {
	return func(t *Toggle) {
		t.Value = value
	}
}

// WithRequired sets the required state.
func WithRequired(required bool) Option {
	return func(t *Toggle) {
		t.Required = required
	}
}

// WithDescription sets the helper text.
func WithDescription(description string) Option {
	return func(t *Toggle) {
		t.Description = description
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(t *Toggle) {
		t.SetStyled(styled)
	}
}

// CheckboxOption is a functional option for checkboxes.
type CheckboxOption func(*Checkbox)

// WithCheckboxChecked sets the initial checked state.
func WithCheckboxChecked(checked bool) CheckboxOption {
	return func(c *Checkbox) {
		c.Checked = checked
	}
}

// WithCheckboxIndeterminate sets the indeterminate state.
func WithCheckboxIndeterminate(indeterminate bool) CheckboxOption {
	return func(c *Checkbox) {
		c.Indeterminate = indeterminate
	}
}

// WithCheckboxDisabled sets the disabled state.
func WithCheckboxDisabled(disabled bool) CheckboxOption {
	return func(c *Checkbox) {
		c.Disabled = disabled
	}
}

// WithCheckboxLabel sets the checkbox label.
func WithCheckboxLabel(label string) CheckboxOption {
	return func(c *Checkbox) {
		c.Label = label
	}
}

// WithCheckboxName sets the form field name.
func WithCheckboxName(name string) CheckboxOption {
	return func(c *Checkbox) {
		c.Name = name
	}
}

// WithCheckboxValue sets the value when checked.
func WithCheckboxValue(value string) CheckboxOption {
	return func(c *Checkbox) {
		c.Value = value
	}
}

// WithCheckboxRequired sets the required state.
func WithCheckboxRequired(required bool) CheckboxOption {
	return func(c *Checkbox) {
		c.Required = required
	}
}

// WithCheckboxDescription sets the helper text.
func WithCheckboxDescription(description string) CheckboxOption {
	return func(c *Checkbox) {
		c.Description = description
	}
}

// WithCheckboxStyled enables Tailwind CSS styling.
func WithCheckboxStyled(styled bool) CheckboxOption {
	return func(c *Checkbox) {
		c.SetStyled(styled)
	}
}
