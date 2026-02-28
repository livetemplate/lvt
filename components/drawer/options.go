package drawer

// Option is a functional option for configuring drawers.
type Option func(*Drawer)

// WithPosition sets the drawer position.
func WithPosition(position Position) Option {
	return func(d *Drawer) {
		d.Position = position
	}
}

// WithSize sets the drawer size.
func WithSize(size Size) Option {
	return func(d *Drawer) {
		d.Size = size
	}
}

// WithTitle sets the drawer title.
func WithTitle(title string) Option {
	return func(d *Drawer) {
		d.Title = title
	}
}

// WithShowClose shows or hides the close button.
func WithShowClose(show bool) Option {
	return func(d *Drawer) {
		d.ShowClose = show
	}
}

// WithShowOverlay shows or hides the backdrop overlay.
func WithShowOverlay(show bool) Option {
	return func(d *Drawer) {
		d.ShowOverlay = show
	}
}

// WithCloseOnOverlay enables closing when clicking overlay.
func WithCloseOnOverlay(close bool) Option {
	return func(d *Drawer) {
		d.CloseOnOverlay = close
	}
}

// WithCloseOnEscape enables closing on Escape key.
func WithCloseOnEscape(close bool) Option {
	return func(d *Drawer) {
		d.CloseOnEscape = close
	}
}

// WithPersistent makes the drawer persistent (can't be closed by user).
func WithPersistent(persistent bool) Option {
	return func(d *Drawer) {
		d.Persistent = persistent
	}
}

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(d *Drawer) {
		d.Open = open
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(d *Drawer) {
		d.SetStyled(styled)
	}
}
