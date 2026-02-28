// Package rating provides star rating components.
//
// Available variants:
//   - New() creates a star rating (template: "lvt:rating:default:v1")
//   - NewReadonly() creates a read-only rating display (template: "lvt:rating:readonly:v1")
//
// Required lvt-* attributes: lvt-click, lvt-mouseover, lvt-mouseleave
//
// Example usage:
//
//	// In your controller/state
//	ProductRating: rating.New("product-rating",
//	    rating.WithMaxStars(5),
//	    rating.WithValue(4),
//	)
//
//	// In your template
//	{{template "lvt:rating:default:v1" .ProductRating}}
package rating

import (
	"math"

	"github.com/livetemplate/components/base"
)

// Rating is a star rating component.
// Use template "lvt:rating:default:v1" to render.
type Rating struct {
	base.Base

	// Value is the current rating (can be fractional for half stars)
	Value float64

	// MaxStars is the maximum number of stars (default 5)
	MaxStars int

	// AllowHalf enables half-star ratings
	AllowHalf bool

	// AllowClear allows clearing by clicking the current value
	AllowClear bool

	// Readonly prevents user interaction
	Readonly bool

	// Size is the star size ("sm", "md", "lg", "xl")
	Size string

	// Color is the active star color
	Color string

	// EmptyColor is the inactive star color
	EmptyColor string

	// HoverValue is the value being hovered (-1 if not hovering)
	HoverValue float64

	// ShowValue displays the numeric value
	ShowValue bool

	// ShowCount displays the rating count
	ShowCount bool

	// Count is the number of ratings (for display)
	Count int

	// Label is optional label text
	Label string

	// Character is the rating character (default star)
	Character string
}

// New creates a star rating.
//
// Example:
//
//	r := rating.New("product",
//	    rating.WithMaxStars(5),
//	    rating.WithAllowHalf(true),
//	)
func New(id string, opts ...Option) *Rating {
	r := &Rating{
		Base:       base.NewBase(id, "rating"),
		MaxStars:   5,
		Size:       "md",
		Color:      "yellow",
		EmptyColor: "gray",
		HoverValue: -1,
		Character:  "★",
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// NewReadonly creates a read-only rating display.
func NewReadonly(id string, value float64, opts ...Option) *Rating {
	r := New(id, opts...)
	r.Value = value
	r.Readonly = true
	return r
}

// SetValue sets the rating value.
func (r *Rating) SetValue(value float64) {
	if value < 0 {
		value = 0
	}
	if value > float64(r.MaxStars) {
		value = float64(r.MaxStars)
	}
	if r.AllowHalf {
		// Round to nearest 0.5
		r.Value = math.Round(value*2) / 2
	} else {
		r.Value = math.Round(value)
	}
}

// Clear clears the rating.
func (r *Rating) Clear() {
	r.Value = 0
}

// Click handles clicking on a star.
func (r *Rating) Click(starIndex int) {
	if r.Readonly {
		return
	}

	newValue := float64(starIndex)

	// If clicking current value and AllowClear, clear rating
	if r.AllowClear && r.Value == newValue {
		r.Clear()
		return
	}

	r.SetValue(newValue)
}

// ClickHalf handles clicking on a half star.
func (r *Rating) ClickHalf(starIndex int, isFirstHalf bool) {
	if r.Readonly {
		return
	}

	var newValue float64
	if r.AllowHalf && isFirstHalf {
		newValue = float64(starIndex) - 0.5
	} else {
		newValue = float64(starIndex)
	}

	// If clicking current value and AllowClear, clear rating
	if r.AllowClear && r.Value == newValue {
		r.Clear()
		return
	}

	r.SetValue(newValue)
}

// Hover sets the hover value.
func (r *Rating) Hover(starIndex int) {
	if r.Readonly {
		return
	}
	r.HoverValue = float64(starIndex)
}

// HoverHalf sets the hover value for half stars.
func (r *Rating) HoverHalf(starIndex int, isFirstHalf bool) {
	if r.Readonly {
		return
	}
	if r.AllowHalf && isFirstHalf {
		r.HoverValue = float64(starIndex) - 0.5
	} else {
		r.HoverValue = float64(starIndex)
	}
}

// Leave clears the hover value.
func (r *Rating) Leave() {
	r.HoverValue = -1
}

// DisplayValue returns the value to display (hover or actual).
func (r *Rating) DisplayValue() float64 {
	if r.HoverValue >= 0 {
		return r.HoverValue
	}
	return r.Value
}

// IsStarFull returns true if a star is fully filled.
func (r *Rating) IsStarFull(starIndex int) bool {
	return r.DisplayValue() >= float64(starIndex)
}

// IsStarHalf returns true if a star is half filled.
func (r *Rating) IsStarHalf(starIndex int) bool {
	val := r.DisplayValue()
	return val >= float64(starIndex)-0.5 && val < float64(starIndex)
}

// IsStarEmpty returns true if a star is empty.
func (r *Rating) IsStarEmpty(starIndex int) bool {
	return r.DisplayValue() < float64(starIndex)-0.5
}

// StarState returns "full", "half", or "empty" for a star.
func (r *Rating) StarState(starIndex int) string {
	if r.IsStarFull(starIndex) {
		return "full"
	}
	if r.IsStarHalf(starIndex) {
		return "half"
	}
	return "empty"
}

// Stars returns the star indices (1 to MaxStars).
func (r *Rating) Stars() []int {
	stars := make([]int, r.MaxStars)
	for i := 0; i < r.MaxStars; i++ {
		stars[i] = i + 1
	}
	return stars
}

// Percentage returns the rating as a percentage.
func (r *Rating) Percentage() float64 {
	if r.MaxStars == 0 {
		return 0
	}
	return (r.Value / float64(r.MaxStars)) * 100
}

// FormatValue returns the formatted value string.
func (r *Rating) FormatValue() string {
	if r.AllowHalf && r.Value != math.Floor(r.Value) {
		return formatFloat(r.Value, 1)
	}
	return formatFloat(r.Value, 0)
}

// FormatValueWithMax returns the value formatted as "4/5".
func (r *Rating) FormatValueWithMax() string {
	return r.FormatValue() + "/" + formatFloat(float64(r.MaxStars), 0)
}

// SizeClass returns CSS size classes.
func (r *Rating) SizeClass() string {
	switch r.Size {
	case "sm":
		return "text-lg"
	case "lg":
		return "text-3xl"
	case "xl":
		return "text-4xl"
	default: // md
		return "text-2xl"
	}
}

// ColorClass returns CSS color classes for filled stars.
func (r *Rating) ColorClass() string {
	switch r.Color {
	case "red":
		return "text-red-500"
	case "blue":
		return "text-blue-500"
	case "green":
		return "text-green-500"
	default: // yellow
		return "text-yellow-400"
	}
}

// EmptyColorClass returns CSS color classes for empty stars.
func (r *Rating) EmptyColorClass() string {
	return "text-gray-300"
}

// Helper function
func formatFloat(f float64, decimals int) string {
	if decimals == 0 {
		return string(rune('0' + int(f)))
	}
	whole := int(f)
	frac := int((f - float64(whole)) * 10)
	if whole >= 10 {
		return string(rune('0'+whole/10)) + string(rune('0'+whole%10)) + "." + string(rune('0'+frac))
	}
	return string(rune('0'+whole)) + "." + string(rune('0'+frac))
}
