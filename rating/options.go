package rating

// Option is a functional option for configuring ratings.
type Option func(*Rating)

// WithValue sets the initial rating value.
func WithValue(value float64) Option {
	return func(r *Rating) {
		r.SetValue(value)
	}
}

// WithMaxStars sets the maximum number of stars.
func WithMaxStars(max int) Option {
	return func(r *Rating) {
		if max > 0 {
			r.MaxStars = max
		}
	}
}

// WithAllowHalf enables half-star ratings.
func WithAllowHalf(allow bool) Option {
	return func(r *Rating) {
		r.AllowHalf = allow
	}
}

// WithAllowClear enables clearing by clicking the current value.
func WithAllowClear(allow bool) Option {
	return func(r *Rating) {
		r.AllowClear = allow
	}
}

// WithReadonly makes the rating read-only.
func WithReadonly(readonly bool) Option {
	return func(r *Rating) {
		r.Readonly = readonly
	}
}

// WithSize sets the star size.
func WithSize(size string) Option {
	return func(r *Rating) {
		r.Size = size
	}
}

// WithColor sets the active star color.
func WithColor(color string) Option {
	return func(r *Rating) {
		r.Color = color
	}
}

// WithShowValue displays the numeric value.
func WithShowValue(show bool) Option {
	return func(r *Rating) {
		r.ShowValue = show
	}
}

// WithShowCount displays the rating count.
func WithShowCount(show bool) Option {
	return func(r *Rating) {
		r.ShowCount = show
	}
}

// WithCount sets the number of ratings.
func WithCount(count int) Option {
	return func(r *Rating) {
		r.Count = count
	}
}

// WithLabel sets the label text.
func WithLabel(label string) Option {
	return func(r *Rating) {
		r.Label = label
	}
}

// WithCharacter sets the rating character.
func WithCharacter(char string) Option {
	return func(r *Rating) {
		r.Character = char
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(r *Rating) {
		r.SetStyled(styled)
	}
}
