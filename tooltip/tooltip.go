// Package tooltip provides tooltip/hint components.
//
// Available variants:
//   - New() creates a tooltip (template: "lvt:tooltip:default:v1")
//
// Required lvt-* attributes: lvt-mouseenter, lvt-mouseleave, lvt-focus, lvt-blur
//
// Example usage:
//
//	// In your controller/state
//	HelpTip: tooltip.New("help-tip",
//	    tooltip.WithContent("Click to learn more"),
//	    tooltip.WithPosition(tooltip.PositionTop),
//	)
//
//	// In your template
//	{{template "lvt:tooltip:default:v1" .HelpTip}}
package tooltip

import (
	"github.com/livetemplate/components/base"
)

// Position defines where the tooltip appears relative to trigger.
type Position string

const (
	PositionTop         Position = "top"
	PositionTopStart    Position = "top-start"
	PositionTopEnd      Position = "top-end"
	PositionBottom      Position = "bottom"
	PositionBottomStart Position = "bottom-start"
	PositionBottomEnd   Position = "bottom-end"
	PositionLeft        Position = "left"
	PositionLeftStart   Position = "left-start"
	PositionLeftEnd     Position = "left-end"
	PositionRight       Position = "right"
	PositionRightStart  Position = "right-start"
	PositionRightEnd    Position = "right-end"
)

// Trigger defines what action shows the tooltip.
type Trigger string

const (
	TriggerHover Trigger = "hover"
	TriggerFocus Trigger = "focus"
	TriggerClick Trigger = "click"
)

// Tooltip is a contextual hint component.
// Use template "lvt:tooltip:default:v1" to render.
type Tooltip struct {
	base.Base

	// Content is the tooltip text
	Content string

	// Position is where the tooltip appears
	Position Position

	// Trigger is what action shows the tooltip
	Trigger Trigger

	// Visible indicates whether the tooltip is shown
	Visible bool

	// Delay in milliseconds before showing (for hover)
	Delay int

	// HideDelay in milliseconds before hiding (for hover)
	HideDelay int

	// Arrow shows a pointing arrow
	Arrow bool

	// MaxWidth limits tooltip width
	MaxWidth string
}

// New creates a tooltip.
//
// Example:
//
//	t := tooltip.New("help",
//	    tooltip.WithContent("Help text"),
//	    tooltip.WithPosition(tooltip.PositionTop),
//	)
func New(id string, opts ...Option) *Tooltip {
	t := &Tooltip{
		Base:      base.NewBase(id, "tooltip"),
		Position:  PositionTop,
		Trigger:   TriggerHover,
		Delay:     0,
		HideDelay: 0,
		Arrow:     true,
		MaxWidth:  "200px",
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Show makes the tooltip visible.
func (t *Tooltip) Show() {
	t.Visible = true
}

// Hide makes the tooltip invisible.
func (t *Tooltip) Hide() {
	t.Visible = false
}

// Toggle toggles tooltip visibility.
func (t *Tooltip) Toggle() {
	t.Visible = !t.Visible
}

// IsTop returns true if position starts with "top".
func (t *Tooltip) IsTop() bool {
	return t.Position == PositionTop ||
		t.Position == PositionTopStart ||
		t.Position == PositionTopEnd
}

// IsBottom returns true if position starts with "bottom".
func (t *Tooltip) IsBottom() bool {
	return t.Position == PositionBottom ||
		t.Position == PositionBottomStart ||
		t.Position == PositionBottomEnd
}

// IsLeft returns true if position starts with "left".
func (t *Tooltip) IsLeft() bool {
	return t.Position == PositionLeft ||
		t.Position == PositionLeftStart ||
		t.Position == PositionLeftEnd
}

// IsRight returns true if position starts with "right".
func (t *Tooltip) IsRight() bool {
	return t.Position == PositionRight ||
		t.Position == PositionRightStart ||
		t.Position == PositionRightEnd
}

// IsHoverTrigger returns true if trigger is hover.
func (t *Tooltip) IsHoverTrigger() bool {
	return t.Trigger == TriggerHover
}

// IsFocusTrigger returns true if trigger is focus.
func (t *Tooltip) IsFocusTrigger() bool {
	return t.Trigger == TriggerFocus
}

// IsClickTrigger returns true if trigger is click.
func (t *Tooltip) IsClickTrigger() bool {
	return t.Trigger == TriggerClick
}

// PositionClasses returns CSS classes for position.
func (t *Tooltip) PositionClasses() string {
	switch t.Position {
	case PositionTop:
		return "bottom-full left-1/2 -translate-x-1/2 mb-2"
	case PositionTopStart:
		return "bottom-full left-0 mb-2"
	case PositionTopEnd:
		return "bottom-full right-0 mb-2"
	case PositionBottom:
		return "top-full left-1/2 -translate-x-1/2 mt-2"
	case PositionBottomStart:
		return "top-full left-0 mt-2"
	case PositionBottomEnd:
		return "top-full right-0 mt-2"
	case PositionLeft:
		return "right-full top-1/2 -translate-y-1/2 mr-2"
	case PositionLeftStart:
		return "right-full top-0 mr-2"
	case PositionLeftEnd:
		return "right-full bottom-0 mr-2"
	case PositionRight:
		return "left-full top-1/2 -translate-y-1/2 ml-2"
	case PositionRightStart:
		return "left-full top-0 ml-2"
	case PositionRightEnd:
		return "left-full bottom-0 ml-2"
	default:
		return "bottom-full left-1/2 -translate-x-1/2 mb-2"
	}
}

// ArrowClasses returns CSS classes for the arrow.
func (t *Tooltip) ArrowClasses() string {
	if !t.Arrow {
		return ""
	}
	switch {
	case t.IsTop():
		return "absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-900"
	case t.IsBottom():
		return "absolute bottom-full left-1/2 -translate-x-1/2 border-4 border-transparent border-b-gray-900"
	case t.IsLeft():
		return "absolute left-full top-1/2 -translate-y-1/2 border-4 border-transparent border-l-gray-900"
	case t.IsRight():
		return "absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-gray-900"
	default:
		return ""
	}
}

// HasContent returns true if tooltip has content.
func (t *Tooltip) HasContent() bool {
	return t.Content != ""
}
