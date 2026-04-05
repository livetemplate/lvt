// Package popover provides rich content popover components.
//
// Available variants:
//   - New() creates a popover (template: "lvt:popover:default:v1")
//
// Open/close is handled client-side via CSS classes and onclick/onmouseenter/onmouseleave handlers.
// Server actions handle data operations only.
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
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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

// Styles returns the resolved PopoverStyles for this component.
func (p *Popover) Styles() styles.PopoverStyles {
	if s, ok := p.StyleData().(styles.PopoverStyles); ok {
		return s
	}
	adapter := styles.ForStyled(p.IsStyled())
	if adapter == nil {
		return styles.PopoverStyles{}
	}
	s := adapter.PopoverStyles()
	p.SetStyleData(s)
	return s
}

// PositionClasses returns CSS classes for position.
func (p *Popover) PositionClasses() string {
	s := p.Styles()
	switch p.Position {
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
		return s.PosBottom
	}
}

// ArrowClasses returns CSS classes for the arrow.
func (p *Popover) ArrowClasses() string {
	if !p.Arrow {
		return ""
	}
	s := p.Styles()
	switch {
	case p.IsTop():
		return s.ArrowTop
	case p.IsBottom():
		return s.ArrowBottom
	case p.IsLeft():
		return s.ArrowLeft
	case p.IsRight():
		return s.ArrowRight
	default:
		return ""
	}
}
