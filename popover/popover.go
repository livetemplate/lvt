// Package popover provides rich content popover components.
//
// Available variants:
//   - New() creates a popover (template: "lvt:popover:default:v1")
//
// Required lvt-* attributes: lvt-click, lvt-click-away
//
// Example usage:
//
//	// In your controller/state
//	InfoPopover: popover.New("info",
//	    popover.WithTitle("Details"),
//	    popover.WithPosition(popover.PositionBottom),
//	)
//
//	// In your template
//	{{template "lvt:popover:default:v1" .InfoPopover}}
package popover

import (
	"github.com/livetemplate/components/base"
)

// Position defines where the popover appears relative to trigger.
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

// Trigger defines what action shows the popover.
type Trigger string

const (
	TriggerClick Trigger = "click"
	TriggerHover Trigger = "hover"
	TriggerFocus Trigger = "focus"
)

// Popover is a rich content overlay component.
// Use template "lvt:popover:default:v1" to render.
type Popover struct {
	base.Base

	// Title is the popover header
	Title string

	// Content is the popover body text
	Content string

	// Position is where the popover appears
	Position Position

	// Trigger is what action shows the popover
	Trigger Trigger

	// Open indicates whether the popover is visible
	Open bool

	// Arrow shows a pointing arrow
	Arrow bool

	// CloseOnClickAway closes popover when clicking outside
	CloseOnClickAway bool

	// ShowClose shows a close button in header
	ShowClose bool

	// Width sets the popover width
	Width string
}

// New creates a popover.
//
// Example:
//
//	p := popover.New("info",
//	    popover.WithTitle("Information"),
//	    popover.WithContent("Some details here"),
//	    popover.WithPosition(popover.PositionBottom),
//	)
func New(id string, opts ...Option) *Popover {
	p := &Popover{
		Base:             base.NewBase(id, "popover"),
		Position:         PositionBottom,
		Trigger:          TriggerClick,
		Arrow:            true,
		CloseOnClickAway: true,
		ShowClose:        false,
		Width:            "280px",
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Show opens the popover.
func (p *Popover) Show() {
	p.Open = true
}

// Hide closes the popover.
func (p *Popover) Hide() {
	p.Open = false
}

// Toggle toggles popover visibility.
func (p *Popover) Toggle() {
	p.Open = !p.Open
}

// IsTop returns true if position starts with "top".
func (p *Popover) IsTop() bool {
	return p.Position == PositionTop ||
		p.Position == PositionTopStart ||
		p.Position == PositionTopEnd
}

// IsBottom returns true if position starts with "bottom".
func (p *Popover) IsBottom() bool {
	return p.Position == PositionBottom ||
		p.Position == PositionBottomStart ||
		p.Position == PositionBottomEnd
}

// IsLeft returns true if position starts with "left".
func (p *Popover) IsLeft() bool {
	return p.Position == PositionLeft ||
		p.Position == PositionLeftStart ||
		p.Position == PositionLeftEnd
}

// IsRight returns true if position starts with "right".
func (p *Popover) IsRight() bool {
	return p.Position == PositionRight ||
		p.Position == PositionRightStart ||
		p.Position == PositionRightEnd
}

// IsClickTrigger returns true if trigger is click.
func (p *Popover) IsClickTrigger() bool {
	return p.Trigger == TriggerClick
}

// IsHoverTrigger returns true if trigger is hover.
func (p *Popover) IsHoverTrigger() bool {
	return p.Trigger == TriggerHover
}

// IsFocusTrigger returns true if trigger is focus.
func (p *Popover) IsFocusTrigger() bool {
	return p.Trigger == TriggerFocus
}

// HasTitle returns true if popover has a title.
func (p *Popover) HasTitle() bool {
	return p.Title != ""
}

// HasContent returns true if popover has content.
func (p *Popover) HasContent() bool {
	return p.Content != ""
}

// HasHeader returns true if popover should show header.
func (p *Popover) HasHeader() bool {
	return p.HasTitle() || p.ShowClose
}

// PositionClasses returns CSS classes for position.
func (p *Popover) PositionClasses() string {
	switch p.Position {
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
		return "top-full left-1/2 -translate-x-1/2 mt-2"
	}
}

// ArrowClasses returns CSS classes for the arrow.
func (p *Popover) ArrowClasses() string {
	if !p.Arrow {
		return ""
	}
	switch {
	case p.IsTop():
		return "absolute top-full left-1/2 -translate-x-1/2 border-8 border-transparent border-t-white"
	case p.IsBottom():
		return "absolute bottom-full left-1/2 -translate-x-1/2 border-8 border-transparent border-b-white"
	case p.IsLeft():
		return "absolute left-full top-1/2 -translate-y-1/2 border-8 border-transparent border-l-white"
	case p.IsRight():
		return "absolute right-full top-1/2 -translate-y-1/2 border-8 border-transparent border-r-white"
	default:
		return ""
	}
}
