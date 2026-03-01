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

	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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

// Styles returns the resolved ProgressStyles for this component.
func (p *Progress) Styles() styles.ProgressStyles {
	if s, ok := p.StyleData().(styles.ProgressStyles); ok {
		return s
	}
	adapter := styles.ForStyled(p.IsStyled())
	if adapter == nil {
		return styles.ProgressStyles{}
	}
	s := adapter.ProgressStyles()
	p.SetStyleData(s)
	return s
}

// SizeClass returns CSS class for size.
func (p *Progress) SizeClass() string {
	st := p.Styles()
	switch p.Size {
	case SizeXs:
		return st.SizeXs
	case SizeSm:
		return st.SizeSm
	case SizeLg:
		return st.SizeLg
	default: // md
		return st.SizeMd
	}
}

// ColorClass returns CSS class for color.
func (p *Progress) ColorClass() string {
	st := p.Styles()
	switch p.Color {
	case ColorSuccess:
		return st.ColorSuccess
	case ColorWarning:
		return st.ColorWarning
	case ColorDanger:
		return st.ColorDanger
	case ColorInfo:
		return st.ColorInfo
	default: // primary
		return st.ColorPrimary
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

// Styles returns the resolved CircularProgressStyles for this component.
func (c *CircularProgress) Styles() styles.CircularProgressStyles {
	if s, ok := c.StyleData().(styles.CircularProgressStyles); ok {
		return s
	}
	adapter := styles.ForStyled(c.IsStyled())
	if adapter == nil {
		return styles.CircularProgressStyles{}
	}
	s := adapter.CircularProgressStyles()
	c.SetStyleData(s)
	return s
}

// ColorClass returns CSS class for color.
func (c *CircularProgress) ColorClass() string {
	st := c.Styles()
	switch c.Color {
	case ColorSuccess:
		return st.ColorSuccess
	case ColorWarning:
		return st.ColorWarning
	case ColorDanger:
		return st.ColorDanger
	case ColorInfo:
		return st.ColorInfo
	default: // primary
		return st.ColorPrimary
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

// Styles returns the resolved SpinnerStyles for this component.
func (s *Spinner) Styles() styles.SpinnerStyles {
	if st, ok := s.StyleData().(styles.SpinnerStyles); ok {
		return st
	}
	adapter := styles.ForStyled(s.IsStyled())
	if adapter == nil {
		return styles.SpinnerStyles{}
	}
	st := adapter.SpinnerStyles()
	s.SetStyleData(st)
	return st
}

// SizeClass returns CSS class for spinner size.
func (s *Spinner) SizeClass() string {
	st := s.Styles()
	switch s.Size {
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

// ColorClass returns CSS class for spinner color.
func (s *Spinner) ColorClass() string {
	st := s.Styles()
	switch s.Color {
	case ColorSuccess:
		return st.ColorSuccess
	case ColorWarning:
		return st.ColorWarning
	case ColorDanger:
		return st.ColorDanger
	case ColorInfo:
		return st.ColorInfo
	default: // primary
		return st.ColorPrimary
	}
}
