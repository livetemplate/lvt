// Package skeleton provides loading placeholder components.
//
// Available variants:
//   - New() creates a skeleton line (template: "lvt:skeleton:default:v1")
//   - NewAvatar() creates avatar placeholder (template: "lvt:skeleton:avatar:v1")
//   - NewCard() creates card placeholder (template: "lvt:skeleton:card:v1")
//
// Example usage:
//
//	// In your controller/state
//	Loading: skeleton.New("line-1",
//	    skeleton.WithWidth("100%"),
//	    skeleton.WithHeight("16px"),
//	)
//
//	// In your template
//	{{template "lvt:skeleton:default:v1" .Loading}}
package skeleton

import (
	"github.com/livetemplate/components/base"
)

// Shape defines the skeleton shape.
type Shape string

const (
	ShapeRectangle Shape = "rectangle"
	ShapeCircle    Shape = "circle"
	ShapeRounded   Shape = "rounded"
)

// Animation defines the skeleton animation type.
type Animation string

const (
	AnimationPulse Animation = "pulse"
	AnimationWave  Animation = "wave"
	AnimationNone  Animation = "none"
)

// Skeleton is a loading placeholder component.
// Use template "lvt:skeleton:default:v1" to render.
type Skeleton struct {
	base.Base

	// Width of the skeleton
	Width string

	// Height of the skeleton
	Height string

	// Shape of the skeleton
	Shape Shape

	// Animation type
	Animation Animation

	// Lines for multi-line skeletons
	Lines int

	// LineHeight for multi-line spacing
	LineHeight string
}

// New creates a skeleton placeholder.
//
// Example:
//
//	s := skeleton.New("loading",
//	    skeleton.WithWidth("200px"),
//	    skeleton.WithHeight("20px"),
//	)
func New(id string, opts ...Option) *Skeleton {
	s := &Skeleton{
		Base:       base.NewBase(id, "skeleton"),
		Width:      "100%",
		Height:     "16px",
		Shape:      ShapeRectangle,
		Animation:  AnimationPulse,
		Lines:      1,
		LineHeight: "24px",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AvatarSkeleton is a circular avatar placeholder.
type AvatarSkeleton struct {
	base.Base

	// Size of the avatar (sm, md, lg, xl)
	Size string

	// ShowBadge shows a small circle for status
	ShowBadge bool
}

// NewAvatar creates an avatar skeleton.
//
// Example:
//
//	a := skeleton.NewAvatar("user-avatar",
//	    skeleton.WithAvatarSize("lg"),
//	)
func NewAvatar(id string, opts ...AvatarOption) *AvatarSkeleton {
	a := &AvatarSkeleton{
		Base: base.NewBase(id, "skeleton"),
		Size: "md",
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// SizeClass returns the CSS class for avatar size.
func (a *AvatarSkeleton) SizeClass() string {
	switch a.Size {
	case "sm":
		return "w-8 h-8"
	case "lg":
		return "w-16 h-16"
	case "xl":
		return "w-24 h-24"
	default: // md
		return "w-12 h-12"
	}
}

// CardSkeleton is a card-shaped placeholder.
type CardSkeleton struct {
	base.Base

	// ShowImage shows image placeholder
	ShowImage bool

	// ImageHeight for image area
	ImageHeight string

	// ShowTitle shows title placeholder
	ShowTitle bool

	// ShowDescription shows description lines
	ShowDescription bool

	// DescriptionLines number of description lines
	DescriptionLines int

	// ShowFooter shows footer area
	ShowFooter bool
}

// NewCard creates a card skeleton.
//
// Example:
//
//	c := skeleton.NewCard("card-loading",
//	    skeleton.WithCardImage(true),
//	    skeleton.WithCardDescription(true, 3),
//	)
func NewCard(id string, opts ...CardOption) *CardSkeleton {
	c := &CardSkeleton{
		Base:             base.NewBase(id, "skeleton"),
		ShowImage:        true,
		ImageHeight:      "200px",
		ShowTitle:        true,
		ShowDescription:  true,
		DescriptionLines: 3,
		ShowFooter:       false,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ShapeClass returns CSS class for the shape.
func (s *Skeleton) ShapeClass() string {
	switch s.Shape {
	case ShapeCircle:
		return "rounded-full"
	case ShapeRounded:
		return "rounded-md"
	default:
		return ""
	}
}

// AnimationClass returns CSS class for animation.
func (s *Skeleton) AnimationClass() string {
	switch s.Animation {
	case AnimationWave:
		return "animate-shimmer"
	case AnimationNone:
		return ""
	default: // pulse
		return "animate-pulse"
	}
}

// IsCircle returns true if shape is circle.
func (s *Skeleton) IsCircle() bool {
	return s.Shape == ShapeCircle
}

// IsRounded returns true if shape is rounded.
func (s *Skeleton) IsRounded() bool {
	return s.Shape == ShapeRounded
}

// IsMultiLine returns true if Lines > 1.
func (s *Skeleton) IsMultiLine() bool {
	return s.Lines > 1
}

// LineIndices returns indices for multi-line rendering.
func (s *Skeleton) LineIndices() []int {
	indices := make([]int, s.Lines)
	for i := range indices {
		indices[i] = i
	}
	return indices
}
