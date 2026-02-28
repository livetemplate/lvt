// Package toggle provides toggle/switch components.
//
// Available variants:
//   - New() creates a toggle switch (template: "lvt:toggle:default:v1")
//   - NewCheckbox() creates a styled checkbox (template: "lvt:toggle:checkbox:v1")
//
// Required lvt-* attributes: lvt-click, lvt-change
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
	"github.com/livetemplate/components/base"
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

// SizeClasses returns CSS classes for the toggle track.
func (t *Toggle) SizeClasses() string {
	switch t.Size {
	case SizeSm:
		return "w-8 h-4"
	case SizeLg:
		return "w-14 h-8"
	default: // md
		return "w-11 h-6"
	}
}

// KnobSizeClasses returns CSS classes for the toggle knob.
func (t *Toggle) KnobSizeClasses() string {
	switch t.Size {
	case SizeSm:
		return "w-3 h-3"
	case SizeLg:
		return "w-6 h-6"
	default: // md
		return "w-5 h-5"
	}
}

// KnobTranslateClass returns CSS class for knob position.
func (t *Toggle) KnobTranslateClass() string {
	if !t.Checked {
		return "translate-x-0.5"
	}
	switch t.Size {
	case SizeSm:
		return "translate-x-4"
	case SizeLg:
		return "translate-x-7"
	default: // md
		return "translate-x-5"
	}
}

// TrackColorClass returns CSS class for track color.
func (t *Toggle) TrackColorClass() string {
	if t.Disabled {
		if t.Checked {
			return "bg-blue-300"
		}
		return "bg-gray-200"
	}
	if t.Checked {
		return "bg-blue-600"
	}
	return "bg-gray-200"
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

// CheckboxStateClass returns CSS class for checkbox state.
func (c *Checkbox) CheckboxStateClass() string {
	if c.Disabled {
		if c.Checked {
			return "bg-blue-300 border-blue-300"
		}
		return "bg-gray-100 border-gray-200"
	}
	if c.Checked {
		return "bg-blue-600 border-blue-600"
	}
	return "bg-white border-gray-300"
}
