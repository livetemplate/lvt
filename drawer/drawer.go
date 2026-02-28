// Package drawer provides slide-out panel/sidebar components.
//
// Available variants:
//   - New() creates a drawer (template: "lvt:drawer:default:v1")
//
// Required lvt-* attributes: lvt-click, lvt-click-away
//
// Example usage:
//
//	// In your controller/state
//	Sidebar: drawer.New("sidebar",
//	    drawer.WithPosition("left"),
//	    drawer.WithTitle("Navigation"),
//	)
//
//	// In your template
//	{{template "lvt:drawer:default:v1" .Sidebar}}
package drawer

import (
	"github.com/livetemplate/components/base"
)

// Position defines where the drawer slides from.
type Position string

const (
	PositionLeft   Position = "left"
	PositionRight  Position = "right"
	PositionTop    Position = "top"
	PositionBottom Position = "bottom"
)

// Size defines the drawer size.
type Size string

const (
	SizeSm   Size = "sm"
	SizeMd   Size = "md"
	SizeLg   Size = "lg"
	SizeXl   Size = "xl"
	SizeFull Size = "full"
)

// Drawer is a slide-out panel component.
// Use template "lvt:drawer:default:v1" to render.
type Drawer struct {
	base.Base

	// Open indicates whether the drawer is visible
	Open bool

	// Position is where the drawer slides from
	Position Position

	// Size is the drawer size
	Size Size

	// Title is the drawer header title
	Title string

	// ShowClose shows a close button
	ShowClose bool

	// ShowOverlay shows a backdrop overlay
	ShowOverlay bool

	// CloseOnOverlay closes drawer when clicking overlay
	CloseOnOverlay bool

	// CloseOnEscape closes drawer on Escape key
	CloseOnEscape bool

	// Persistent prevents closing (no overlay click, no escape)
	Persistent bool
}

// New creates a drawer.
//
// Example:
//
//	d := drawer.New("sidebar",
//	    drawer.WithPosition(drawer.PositionLeft),
//	    drawer.WithSize(drawer.SizeMd),
//	)
func New(id string, opts ...Option) *Drawer {
	d := &Drawer{
		Base:           base.NewBase(id, "drawer"),
		Position:       PositionLeft,
		Size:           SizeMd,
		ShowClose:      true,
		ShowOverlay:    true,
		CloseOnOverlay: true,
		CloseOnEscape:  true,
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Toggle opens or closes the drawer.
func (d *Drawer) Toggle() {
	d.Open = !d.Open
}

// Show opens the drawer.
func (d *Drawer) Show() {
	d.Open = true
}

// Close closes the drawer.
func (d *Drawer) Close() {
	if !d.Persistent {
		d.Open = false
	}
}

// ForceClose closes the drawer even if persistent.
func (d *Drawer) ForceClose() {
	d.Open = false
}

// IsLeft returns true if position is left.
func (d *Drawer) IsLeft() bool {
	return d.Position == PositionLeft
}

// IsRight returns true if position is right.
func (d *Drawer) IsRight() bool {
	return d.Position == PositionRight
}

// IsTop returns true if position is top.
func (d *Drawer) IsTop() bool {
	return d.Position == PositionTop
}

// IsBottom returns true if position is bottom.
func (d *Drawer) IsBottom() bool {
	return d.Position == PositionBottom
}

// IsHorizontal returns true if position is left or right.
func (d *Drawer) IsHorizontal() bool {
	return d.IsLeft() || d.IsRight()
}

// IsVertical returns true if position is top or bottom.
func (d *Drawer) IsVertical() bool {
	return d.IsTop() || d.IsBottom()
}

// SizeClass returns CSS classes for the size.
func (d *Drawer) SizeClass() string {
	if d.IsHorizontal() {
		switch d.Size {
		case SizeSm:
			return "w-64"
		case SizeLg:
			return "w-96"
		case SizeXl:
			return "w-[32rem]"
		case SizeFull:
			return "w-full"
		default: // md
			return "w-80"
		}
	}
	// Vertical
	switch d.Size {
	case SizeSm:
		return "h-48"
	case SizeLg:
		return "h-96"
	case SizeXl:
		return "h-[32rem]"
	case SizeFull:
		return "h-full"
	default: // md
		return "h-64"
	}
}

// PositionClass returns CSS classes for position.
func (d *Drawer) PositionClass() string {
	switch d.Position {
	case PositionRight:
		return "right-0 top-0 h-full"
	case PositionTop:
		return "top-0 left-0 w-full"
	case PositionBottom:
		return "bottom-0 left-0 w-full"
	default: // left
		return "left-0 top-0 h-full"
	}
}

// TransformClass returns CSS transform classes for animation.
func (d *Drawer) TransformClass() string {
	if d.Open {
		return "translate-x-0 translate-y-0"
	}
	switch d.Position {
	case PositionRight:
		return "translate-x-full"
	case PositionTop:
		return "-translate-y-full"
	case PositionBottom:
		return "translate-y-full"
	default: // left
		return "-translate-x-full"
	}
}

// HasTitle returns true if drawer has a title.
func (d *Drawer) HasTitle() bool {
	return d.Title != ""
}
