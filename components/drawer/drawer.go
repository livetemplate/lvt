// Package drawer provides slide-out panel/sidebar components.
//
// Available variants:
//   - New() creates a drawer (template: "lvt:drawer:default:v1")
//
// Open/close is handled client-side via CSS classes.
// The consuming page triggers the drawer by adding/removing the "open" class
// on the element with data-drawer attribute.
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
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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
// Open/close is handled client-side via CSS classes.
// Use template "lvt:drawer:default:v1" to render.
type Drawer struct {
	base.Base

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

// Styles returns the resolved DrawerStyles for this component.
func (d *Drawer) Styles() styles.DrawerStyles {
	if s, ok := d.StyleData().(styles.DrawerStyles); ok {
		return s
	}
	adapter := styles.ForStyled(d.IsStyled())
	if adapter == nil {
		return styles.DrawerStyles{}
	}
	s := adapter.DrawerStyles()
	d.SetStyleData(s)
	return s
}

// SizeClass returns CSS classes for the size.
func (d *Drawer) SizeClass() string {
	s := d.Styles()
	if d.IsHorizontal() {
		switch d.Size {
		case SizeSm:
			return s.SizeSmH
		case SizeLg:
			return s.SizeLgH
		case SizeXl:
			return s.SizeXlH
		case SizeFull:
			return s.SizeFullH
		default: // md
			return s.SizeMdH
		}
	}
	// Vertical
	switch d.Size {
	case SizeSm:
		return s.SizeSmV
	case SizeLg:
		return s.SizeLgV
	case SizeXl:
		return s.SizeXlV
	case SizeFull:
		return s.SizeFullV
	default: // md
		return s.SizeMdV
	}
}

// PositionClass returns CSS classes for position.
func (d *Drawer) PositionClass() string {
	s := d.Styles()
	switch d.Position {
	case PositionRight:
		return s.PositionRight
	case PositionTop:
		return s.PositionTop
	case PositionBottom:
		return s.PositionBottom
	default: // left
		return s.PositionLeft
	}
}

// TransformClass returns the CSS transform value for the closed state.
// This is used as a CSS custom property; the open state uses translate(0,0).
func (d *Drawer) TransformClass() string {
	switch d.Position {
	case PositionRight:
		return "translateX(100%)"
	case PositionTop:
		return "translateY(-100%)"
	case PositionBottom:
		return "translateY(100%)"
	default: // left
		return "translateX(-100%)"
	}
}

// HasTitle returns true if drawer has a title.
func (d *Drawer) HasTitle() bool {
	return d.Title != ""
}
