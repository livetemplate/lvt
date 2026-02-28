package popover

// Option is a functional option for configuring popovers.
type Option func(*Popover)

// WithTitle sets the popover title.
func WithTitle(title string) Option {
	return func(p *Popover) {
		p.Title = title
	}
}

// WithContent sets the popover content.
func WithContent(content string) Option {
	return func(p *Popover) {
		p.Content = content
	}
}

// WithPosition sets the popover position.
func WithPosition(position Position) Option {
	return func(p *Popover) {
		p.Position = position
	}
}

// WithTrigger sets the popover trigger.
func WithTrigger(trigger Trigger) Option {
	return func(p *Popover) {
		p.Trigger = trigger
	}
}

// WithArrow enables or disables the arrow.
func WithArrow(show bool) Option {
	return func(p *Popover) {
		p.Arrow = show
	}
}

// WithCloseOnClickAway enables closing on click outside.
func WithCloseOnClickAway(close bool) Option {
	return func(p *Popover) {
		p.CloseOnClickAway = close
	}
}

// WithShowClose shows a close button.
func WithShowClose(show bool) Option {
	return func(p *Popover) {
		p.ShowClose = show
	}
}

// WithWidth sets the popover width.
func WithWidth(width string) Option {
	return func(p *Popover) {
		p.Width = width
	}
}

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(p *Popover) {
		p.Open = open
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(p *Popover) {
		p.SetStyled(styled)
	}
}
