// Package progress provides progress bar and spinner components.
//
// Available variants:
//   - New() creates a progress bar (template: "lvt:progress:default:v1")
//   - NewCircular() creates a circular progress (template: "lvt:progress:circular:v1")
//   - NewSpinner() creates a spinner (template: "lvt:progress:spinner:v1")
//
// Example usage:
//
//	// In your controller/state
//	UploadProgress: progress.New("upload",
//	    progress.WithValue(75),
//	    progress.WithShowLabel(true),
//	)
//
//	// In your template
//	{{template "lvt:progress:default:v1" .UploadProgress}}
package progress

import (
	"fmt"

	"github.com/livetemplate/components/base"
)

// Size defines the progress bar size.
type Size string

const (
	SizeXs Size = "xs"
	SizeSm Size = "sm"
	SizeMd Size = "md"
	SizeLg Size = "lg"
)

// Color defines the progress color.
type Color string

const (
	ColorPrimary Color = "primary"
	ColorSuccess Color = "success"
	ColorWarning Color = "warning"
	ColorDanger  Color = "danger"
	ColorInfo    Color = "info"
)

// Progress is a linear progress bar component.
// Use template "lvt:progress:default:v1" to render.
type Progress struct {
	base.Base

	// Value is the current progress (0-100)
	Value float64

	// Max is the maximum value (default 100)
	Max float64

	// Size of the progress bar
	Size Size

	// Color of the progress bar
	Color Color

	// ShowLabel shows percentage label
	ShowLabel bool

	// Label is custom label text (overrides percentage)
	Label string

	// Striped shows striped pattern
	Striped bool

	// Animated animates the stripes
	Animated bool

	// Indeterminate shows indeterminate animation
	Indeterminate bool
}

// New creates a progress bar.
//
// Example:
//
//	p := progress.New("download",
//	    progress.WithValue(50),
//	    progress.WithShowLabel(true),
//	)
func New(id string, opts ...Option) *Progress {
	p := &Progress{
		Base:  base.NewBase(id, "progress"),
		Value: 0,
		Max:   100,
		Size:  SizeMd,
		Color: ColorPrimary,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Percentage returns the progress percentage (0-100).
func (p *Progress) Percentage() float64 {
	if p.Max <= 0 {
		return 0
	}
	pct := (p.Value / p.Max) * 100
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}

// PercentageStr returns the percentage as a formatted string.
func (p *Progress) PercentageStr() string {
	return fmt.Sprintf("%.0f%%", p.Percentage())
}

// DisplayLabel returns the label to display.
func (p *Progress) DisplayLabel() string {
	if p.Label != "" {
		return p.Label
	}
	return p.PercentageStr()
}

// SetValue sets the progress value.
func (p *Progress) SetValue(value float64) {
	p.Value = value
}

// Increment increases the value by amount.
func (p *Progress) Increment(amount float64) {
	p.Value += amount
	if p.Value > p.Max {
		p.Value = p.Max
	}
}

// Decrement decreases the value by amount.
func (p *Progress) Decrement(amount float64) {
	p.Value -= amount
	if p.Value < 0 {
		p.Value = 0
	}
}

// Reset sets value to 0.
func (p *Progress) Reset() {
	p.Value = 0
}

// Complete sets value to max.
func (p *Progress) Complete() {
	p.Value = p.Max
}

// IsComplete returns true if value equals max.
func (p *Progress) IsComplete() bool {
	return p.Value >= p.Max
}

// SizeClass returns CSS class for size.
func (p *Progress) SizeClass() string {
	switch p.Size {
	case SizeXs:
		return "h-1"
	case SizeSm:
		return "h-2"
	case SizeLg:
		return "h-6"
	default: // md
		return "h-4"
	}
}

// ColorClass returns CSS class for color.
func (p *Progress) ColorClass() string {
	switch p.Color {
	case ColorSuccess:
		return "bg-green-500"
	case ColorWarning:
		return "bg-yellow-500"
	case ColorDanger:
		return "bg-red-500"
	case ColorInfo:
		return "bg-cyan-500"
	default: // primary
		return "bg-blue-500"
	}
}

// CircularProgress is a circular progress indicator.
type CircularProgress struct {
	base.Base

	// Value is the current progress (0-100)
	Value float64

	// Max is the maximum value (default 100)
	Max float64

	// Size in pixels
	Size int

	// StrokeWidth of the circle
	StrokeWidth int

	// Color of the progress
	Color Color

	// ShowLabel shows percentage in center
	ShowLabel bool

	// Label is custom label text
	Label string

	// Indeterminate shows spinning animation
	Indeterminate bool
}

// NewCircular creates a circular progress indicator.
//
// Example:
//
//	c := progress.NewCircular("loading",
//	    progress.WithCircularValue(75),
//	    progress.WithCircularSize(80),
//	)
func NewCircular(id string, opts ...CircularOption) *CircularProgress {
	c := &CircularProgress{
		Base:        base.NewBase(id, "progress"),
		Value:       0,
		Max:         100,
		Size:        48,
		StrokeWidth: 4,
		Color:       ColorPrimary,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Percentage returns the progress percentage.
func (c *CircularProgress) Percentage() float64 {
	if c.Max <= 0 {
		return 0
	}
	pct := (c.Value / c.Max) * 100
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}

// PercentageStr returns the percentage as formatted string.
func (c *CircularProgress) PercentageStr() string {
	return fmt.Sprintf("%.0f%%", c.Percentage())
}

// DisplayLabel returns the label to display.
func (c *CircularProgress) DisplayLabel() string {
	if c.Label != "" {
		return c.Label
	}
	return c.PercentageStr()
}

// Radius returns the circle radius.
func (c *CircularProgress) Radius() int {
	return (c.Size - c.StrokeWidth) / 2
}

// Circumference returns the circle circumference.
func (c *CircularProgress) Circumference() float64 {
	return 2 * 3.14159 * float64(c.Radius())
}

// DashOffset returns the stroke-dashoffset for progress.
func (c *CircularProgress) DashOffset() float64 {
	return c.Circumference() * (1 - c.Percentage()/100)
}

// Center returns the center point of the SVG.
func (c *CircularProgress) Center() int {
	return c.Size / 2
}

// ColorClass returns CSS class for color.
func (c *CircularProgress) ColorClass() string {
	switch c.Color {
	case ColorSuccess:
		return "text-green-500"
	case ColorWarning:
		return "text-yellow-500"
	case ColorDanger:
		return "text-red-500"
	case ColorInfo:
		return "text-cyan-500"
	default: // primary
		return "text-blue-500"
	}
}

// Spinner is a loading spinner component.
type Spinner struct {
	base.Base

	// Size of the spinner (sm, md, lg)
	Size string

	// Color of the spinner
	Color Color

	// Label for accessibility
	Label string
}

// NewSpinner creates a spinner.
//
// Example:
//
//	s := progress.NewSpinner("loading",
//	    progress.WithSpinnerSize("lg"),
//	)
func NewSpinner(id string, opts ...SpinnerOption) *Spinner {
	s := &Spinner{
		Base:  base.NewBase(id, "progress"),
		Size:  "md",
		Color: ColorPrimary,
		Label: "Loading...",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SizeClass returns CSS class for spinner size.
func (s *Spinner) SizeClass() string {
	switch s.Size {
	case "sm":
		return "w-4 h-4"
	case "lg":
		return "w-8 h-8"
	case "xl":
		return "w-12 h-12"
	default: // md
		return "w-6 h-6"
	}
}

// ColorClass returns CSS class for spinner color.
func (s *Spinner) ColorClass() string {
	switch s.Color {
	case ColorSuccess:
		return "text-green-500"
	case ColorWarning:
		return "text-yellow-500"
	case ColorDanger:
		return "text-red-500"
	case ColorInfo:
		return "text-cyan-500"
	default: // primary
		return "text-blue-500"
	}
}
