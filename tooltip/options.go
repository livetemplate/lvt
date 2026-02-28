package tooltip

// Option is a functional option for configuring tooltips.
type Option func(*Tooltip)

// WithContent sets the tooltip content.
func WithContent(content string) Option {
	return func(t *Tooltip) {
		t.Content = content
	}
}

// WithPosition sets the tooltip position.
func WithPosition(position Position) Option {
	return func(t *Tooltip) {
		t.Position = position
	}
}

// WithTrigger sets the tooltip trigger.
func WithTrigger(trigger Trigger) Option {
	return func(t *Tooltip) {
		t.Trigger = trigger
	}
}

// WithDelay sets the show delay in milliseconds.
func WithDelay(delay int) Option {
	return func(t *Tooltip) {
		t.Delay = delay
	}
}

// WithHideDelay sets the hide delay in milliseconds.
func WithHideDelay(delay int) Option {
	return func(t *Tooltip) {
		t.HideDelay = delay
	}
}

// WithArrow enables or disables the arrow.
func WithArrow(show bool) Option {
	return func(t *Tooltip) {
		t.Arrow = show
	}
}

// WithMaxWidth sets the maximum width.
func WithMaxWidth(width string) Option {
	return func(t *Tooltip) {
		t.MaxWidth = width
	}
}

// WithVisible sets the initial visibility state.
func WithVisible(visible bool) Option {
	return func(t *Tooltip) {
		t.Visible = visible
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(t *Tooltip) {
		t.SetStyled(styled)
	}
}
