package progress

// Option is a functional option for configuring progress bars.
type Option func(*Progress)

// WithValue sets the current progress value.
func WithValue(value float64) Option {
	return func(p *Progress) {
		p.Value = value
	}
}

// WithMax sets the maximum value.
func WithMax(max float64) Option {
	return func(p *Progress) {
		p.Max = max
	}
}

// WithSize sets the progress bar size.
func WithSize(size Size) Option {
	return func(p *Progress) {
		p.Size = size
	}
}

// WithColor sets the progress bar color.
func WithColor(color Color) Option {
	return func(p *Progress) {
		p.Color = color
	}
}

// WithShowLabel shows the percentage label.
func WithShowLabel(show bool) Option {
	return func(p *Progress) {
		p.ShowLabel = show
	}
}

// WithLabel sets a custom label.
func WithLabel(label string) Option {
	return func(p *Progress) {
		p.Label = label
	}
}

// WithStriped enables striped pattern.
func WithStriped(striped bool) Option {
	return func(p *Progress) {
		p.Striped = striped
	}
}

// WithAnimated enables stripe animation.
func WithAnimated(animated bool) Option {
	return func(p *Progress) {
		p.Animated = animated
	}
}

// WithIndeterminate enables indeterminate mode.
func WithIndeterminate(indeterminate bool) Option {
	return func(p *Progress) {
		p.Indeterminate = indeterminate
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(p *Progress) {
		p.SetStyled(styled)
	}
}

// CircularOption is a functional option for circular progress.
type CircularOption func(*CircularProgress)

// WithCircularValue sets the current progress value.
func WithCircularValue(value float64) CircularOption {
	return func(c *CircularProgress) {
		c.Value = value
	}
}

// WithCircularMax sets the maximum value.
func WithCircularMax(max float64) CircularOption {
	return func(c *CircularProgress) {
		c.Max = max
	}
}

// WithCircularSize sets the size in pixels.
func WithCircularSize(size int) CircularOption {
	return func(c *CircularProgress) {
		c.Size = size
	}
}

// WithCircularStrokeWidth sets the stroke width.
func WithCircularStrokeWidth(width int) CircularOption {
	return func(c *CircularProgress) {
		c.StrokeWidth = width
	}
}

// WithCircularColor sets the color.
func WithCircularColor(color Color) CircularOption {
	return func(c *CircularProgress) {
		c.Color = color
	}
}

// WithCircularShowLabel shows percentage in center.
func WithCircularShowLabel(show bool) CircularOption {
	return func(c *CircularProgress) {
		c.ShowLabel = show
	}
}

// WithCircularLabel sets a custom label.
func WithCircularLabel(label string) CircularOption {
	return func(c *CircularProgress) {
		c.Label = label
	}
}

// WithCircularIndeterminate enables spinning mode.
func WithCircularIndeterminate(indeterminate bool) CircularOption {
	return func(c *CircularProgress) {
		c.Indeterminate = indeterminate
	}
}

// WithCircularStyled enables Tailwind CSS styling.
func WithCircularStyled(styled bool) CircularOption {
	return func(c *CircularProgress) {
		c.SetStyled(styled)
	}
}

// SpinnerOption is a functional option for spinners.
type SpinnerOption func(*Spinner)

// WithSpinnerSize sets the spinner size.
func WithSpinnerSize(size string) SpinnerOption {
	return func(s *Spinner) {
		s.Size = size
	}
}

// WithSpinnerColor sets the spinner color.
func WithSpinnerColor(color Color) SpinnerOption {
	return func(s *Spinner) {
		s.Color = color
	}
}

// WithSpinnerLabel sets the accessibility label.
func WithSpinnerLabel(label string) SpinnerOption {
	return func(s *Spinner) {
		s.Label = label
	}
}

// WithSpinnerStyled enables Tailwind CSS styling.
func WithSpinnerStyled(styled bool) SpinnerOption {
	return func(s *Spinner) {
		s.SetStyled(styled)
	}
}
