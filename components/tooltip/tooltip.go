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
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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

// Styles returns the resolved TooltipStyles for this component.
func (t *Tooltip) Styles() styles.TooltipStyles {
	if s, ok := t.StyleData().(styles.TooltipStyles); ok {
		return s
	}
	adapter := styles.ForStyled(t.IsStyled())
	if adapter == nil {
		return styles.TooltipStyles{}
	}
	s := adapter.TooltipStyles()
	t.SetStyleData(s)
	return s
}

// PositionClasses returns CSS classes for position.
func (t *Tooltip) PositionClasses() string {
	s := t.Styles()
	switch t.Position {
	case PositionTop:
		return s.PosTop
	case PositionTopStart:
		return s.PosTopStart
	case PositionTopEnd:
		return s.PosTopEnd
	case PositionBottom:
		return s.PosBottom
	case PositionBottomStart:
		return s.PosBottomStart
	case PositionBottomEnd:
		return s.PosBottomEnd
	case PositionLeft:
		return s.PosLeft
	case PositionLeftStart:
		return s.PosLeftStart
	case PositionLeftEnd:
		return s.PosLeftEnd
	case PositionRight:
		return s.PosRight
	case PositionRightStart:
		return s.PosRightStart
	case PositionRightEnd:
		return s.PosRightEnd
	default:
		return s.PosTop
	}
}

// ArrowClasses returns CSS classes for the arrow.
func (t *Tooltip) ArrowClasses() string {
	if !t.Arrow {
		return ""
	}
	s := t.Styles()
	switch {
	case t.IsTop():
		return s.ArrowTop
	case t.IsBottom():
		return s.ArrowBottom
	case t.IsLeft():
		return s.ArrowLeft
	case t.IsRight():
		return s.ArrowRight
	default:
		return ""
	}
}

// HasContent returns true if tooltip has content.
func (t *Tooltip) HasContent() bool {
	return t.Content != ""
}
