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
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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

// Styles returns the resolved AvatarSkeletonStyles for this component.
func (a *AvatarSkeleton) Styles() styles.AvatarSkeletonStyles {
	if st, ok := a.StyleData().(styles.AvatarSkeletonStyles); ok {
		return st
	}
	adapter := styles.ForStyled(a.IsStyled())
	if adapter == nil {
		return styles.AvatarSkeletonStyles{}
	}
	st := adapter.AvatarSkeletonStyles()
	a.SetStyleData(st)
	return st
}

// SizeClass returns the CSS class for avatar size.
func (a *AvatarSkeleton) SizeClass() string {
	st := a.Styles()
	switch a.Size {
	case "sm":
		return st.SizeSm
	case "lg":
		return st.SizeLg
	case "xl":
		return st.SizeXl
	default: // md
		return st.SizeMd
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

// Styles returns the resolved CardSkeletonStyles for this component.
func (c *CardSkeleton) Styles() styles.CardSkeletonStyles {
	if st, ok := c.StyleData().(styles.CardSkeletonStyles); ok {
		return st
	}
	adapter := styles.ForStyled(c.IsStyled())
	if adapter == nil {
		return styles.CardSkeletonStyles{}
	}
	st := adapter.CardSkeletonStyles()
	c.SetStyleData(st)
	return st
}

// DescLineIndices returns indices for description line rendering.
func (c *CardSkeleton) DescLineIndices() []int {
	indices := make([]int, c.DescriptionLines)
	for i := range indices {
		indices[i] = i + 1
	}
	return indices
}

// IsLastDescLine returns true if this is the last description line.
func (c *CardSkeleton) IsLastDescLine(i int) bool {
	return i == c.DescriptionLines
}

// Styles returns the resolved SkeletonStyles for this component.
func (s *Skeleton) Styles() styles.SkeletonStyles {
	if st, ok := s.StyleData().(styles.SkeletonStyles); ok {
		return st
	}
	adapter := styles.ForStyled(s.IsStyled())
	if adapter == nil {
		return styles.SkeletonStyles{}
	}
	st := adapter.SkeletonStyles()
	s.SetStyleData(st)
	return st
}

// ShapeClass returns CSS class for the shape.
func (s *Skeleton) ShapeClass() string {
	st := s.Styles()
	switch s.Shape {
	case ShapeCircle:
		return st.ShapeCircle
	case ShapeRounded:
		return st.ShapeRounded
	default:
		return ""
	}
}

// AnimationClass returns CSS class for animation.
func (s *Skeleton) AnimationClass() string {
	st := s.Styles()
	switch s.Animation {
	case AnimationWave:
		return st.AnimationWave
	case AnimationNone:
		return ""
	default: // pulse
		return st.AnimationPulse
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
