// Package toggle provides toggle/switch components.
//
// Available variants:
//   - New() creates a toggle switch (template: "lvt:toggle:default:v1")
//   - NewCheckbox() creates a styled checkbox (template: "lvt:toggle:checkbox:v1")
//
// Required lvt-* attributes: lvt-on:click, lvt-on:change
//
// Example usage:
//
//	// In your controller/state
//	DarkMode: toggle.New("dark-mode",
//	    toggle.WithLabel("Dark Mode"),
//	    toggle.WithChecked(true),
//	)
//
//	// In your template
//	{{template "lvt:toggle:default:v1" .DarkMode}}
package toggle

import (
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
)

// Size defines the toggle size.
type Size string

const (
	SizeSm Size = "sm"
	SizeMd Size = "md"
	SizeLg Size = "lg"
)

// Toggle is a switch/toggle component.
// Use template "lvt:toggle:default:v1" to render.
type Toggle struct {
	base.Base

	// Checked indicates whether the toggle is on
	Checked bool

	// Disabled prevents interaction
	Disabled bool

	// Label is the toggle label text
	Label string

	// LabelPosition is where the label appears (left, right)
	LabelPosition string

	// Size of the toggle
	Size Size

	// Name for form submission
	Name string

	// Value for form submission when checked
	Value string

	// Required for form validation
	Required bool

	// Description is helper text below the toggle
	Description string
}

// New creates a toggle switch.
//
// Example:
//
//	t := toggle.New("notifications",
//	    toggle.WithLabel("Email Notifications"),
//	    toggle.WithChecked(user.EmailNotifications),
//	)
func New(id string, opts ...Option) *Toggle {
	t := &Toggle{
		Base:          base.NewBase(id, "toggle"),
		LabelPosition: "right",
		Size:          SizeMd,
		Value:         "on",
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Toggle toggles the checked state.
func (t *Toggle) Toggle() {
	if !t.Disabled {
		t.Checked = !t.Checked
	}
}

// Check sets checked to true.
func (t *Toggle) Check() {
	if !t.Disabled {
		t.Checked = true
	}
}

// Uncheck sets checked to false.
func (t *Toggle) Uncheck() {
	if !t.Disabled {
		t.Checked = false
	}
}

// SetChecked sets the checked state.
func (t *Toggle) SetChecked(checked bool) {
	if !t.Disabled {
		t.Checked = checked
	}
}

// IsOn returns true if toggle is checked.
func (t *Toggle) IsOn() bool {
	return t.Checked
}

// IsOff returns true if toggle is not checked.
func (t *Toggle) IsOff() bool {
	return !t.Checked
}

// HasLabel returns true if toggle has a label.
func (t *Toggle) HasLabel() bool {
	return t.Label != ""
}

// HasDescription returns true if toggle has a description.
func (t *Toggle) HasDescription() bool {
	return t.Description != ""
}

// IsLabelLeft returns true if label is on left.
func (t *Toggle) IsLabelLeft() bool {
	return t.LabelPosition == "left"
}

// IsLabelRight returns true if label is on right.
func (t *Toggle) IsLabelRight() bool {
	return t.LabelPosition == "right"
}

// Styles returns the resolved ToggleStyles for this component.
func (t *Toggle) Styles() styles.ToggleStyles {
	if s, ok := t.StyleData().(styles.ToggleStyles); ok {
		return s
	}
	adapter := styles.ForStyled(t.IsStyled())
	if adapter == nil {
		return styles.ToggleStyles{}
	}
	s := adapter.ToggleStyles()
	t.SetStyleData(s)
	return s
}

// SizeClasses returns CSS classes for the toggle track.
func (t *Toggle) SizeClasses() string {
	s := t.Styles()
	switch t.Size {
	case SizeSm:
		return s.TrackSm
	case SizeLg:
		return s.TrackLg
	default:
		return s.TrackMd
	}
}

// KnobSizeClasses returns CSS classes for the toggle knob.
func (t *Toggle) KnobSizeClasses() string {
	s := t.Styles()
	switch t.Size {
	case SizeSm:
		return s.KnobSm
	case SizeLg:
		return s.KnobLg
	default:
		return s.KnobMd
	}
}

// KnobTranslateClass returns CSS class for knob position.
func (t *Toggle) KnobTranslateClass() string {
	s := t.Styles()
	if !t.Checked {
		return s.KnobUnchecked
	}
	switch t.Size {
	case SizeSm:
		return s.KnobSmChecked
	case SizeLg:
		return s.KnobLgChecked
	default:
		return s.KnobMdChecked
	}
}

// TrackColorClass returns CSS class for track color.
func (t *Toggle) TrackColorClass() string {
	s := t.Styles()
	if t.Disabled {
		if t.Checked {
			return s.TrackCheckedDisabled
		}
		return s.TrackUncheckedDisabled
	}
	if t.Checked {
		return s.TrackChecked
	}
	return s.TrackUnchecked
}

// Checkbox is a styled checkbox component.
type Checkbox struct {
	base.Base

	// Checked indicates whether the checkbox is checked
	Checked bool

	// Indeterminate shows the indeterminate state
	Indeterminate bool

	// Disabled prevents interaction
	Disabled bool

	// Label is the checkbox label text
	Label string

	// Name for form submission
	Name string

	// Value for form submission when checked
	Value string

	// Required for form validation
	Required bool

	// Description is helper text
	Description string
}

// NewCheckbox creates a styled checkbox.
//
// Example:
//
//	c := toggle.NewCheckbox("terms",
//	    toggle.WithCheckboxLabel("I agree to the terms"),
//	    toggle.WithCheckboxRequired(true),
//	)
func NewCheckbox(id string, opts ...CheckboxOption) *Checkbox {
	c := &Checkbox{
		Base:  base.NewBase(id, "toggle"),
		Value: "on",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Toggle toggles the checked state.
func (c *Checkbox) Toggle() {
	if !c.Disabled {
		c.Checked = !c.Checked
		c.Indeterminate = false
	}
}

// Check sets checked to true.
func (c *Checkbox) Check() {
	if !c.Disabled {
		c.Checked = true
		c.Indeterminate = false
	}
}

// Uncheck sets checked to false.
func (c *Checkbox) Uncheck() {
	if !c.Disabled {
		c.Checked = false
		c.Indeterminate = false
	}
}

// SetIndeterminate sets the indeterminate state.
func (c *Checkbox) SetIndeterminate(indeterminate bool) {
	c.Indeterminate = indeterminate
	if indeterminate {
		c.Checked = false
	}
}

// HasLabel returns true if checkbox has a label.
func (c *Checkbox) HasLabel() bool {
	return c.Label != ""
}

// HasDescription returns true if checkbox has a description.
func (c *Checkbox) HasDescription() bool {
	return c.Description != ""
}

// Styles returns the resolved CheckboxStyles for this component.
func (c *Checkbox) Styles() styles.CheckboxStyles {
	if s, ok := c.StyleData().(styles.CheckboxStyles); ok {
		return s
	}
	adapter := styles.ForStyled(c.IsStyled())
	if adapter == nil {
		return styles.CheckboxStyles{}
	}
	s := adapter.CheckboxStyles()
	c.SetStyleData(s)
	return s
}

// CheckboxStateClass returns CSS class for checkbox state.
func (c *Checkbox) CheckboxStateClass() string {
	s := c.Styles()
	if c.Disabled {
		if c.Checked {
			return s.StateCheckedDisabled
		}
		return s.StateUncheckedDisabled
	}
	if c.Checked {
		return s.StateChecked
	}
	return s.StateUnchecked
}
