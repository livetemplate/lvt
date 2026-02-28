package modal

// Option is a functional option for configuring modals.
type Option func(*Modal)

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(m *Modal) {
		m.Open = open
	}
}

// WithTitle sets the modal title.
func WithTitle(title string) Option {
	return func(m *Modal) {
		m.Title = title
	}
}

// WithSize sets the modal size.
func WithSize(size Size) Option {
	return func(m *Modal) {
		m.Size = size
	}
}

// WithShowClose shows or hides the close button.
func WithShowClose(show bool) Option {
	return func(m *Modal) {
		m.ShowClose = show
	}
}

// WithCloseOnOverlay enables closing on overlay click.
func WithCloseOnOverlay(close bool) Option {
	return func(m *Modal) {
		m.CloseOnOverlay = close
	}
}

// WithCloseOnEscape enables closing on Escape key.
func WithCloseOnEscape(close bool) Option {
	return func(m *Modal) {
		m.CloseOnEscape = close
	}
}

// WithCentered enables vertical centering.
func WithCentered(centered bool) Option {
	return func(m *Modal) {
		m.Centered = centered
	}
}

// WithScrollable enables scrollable modal body.
func WithScrollable(scrollable bool) Option {
	return func(m *Modal) {
		m.Scrollable = scrollable
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(m *Modal) {
		m.SetStyled(styled)
	}
}

// ConfirmOption is a functional option for confirm dialogs.
type ConfirmOption func(*ConfirmModal)

// WithConfirmOpen sets the initial open state.
func WithConfirmOpen(open bool) ConfirmOption {
	return func(c *ConfirmModal) {
		c.Open = open
	}
}

// WithConfirmTitle sets the dialog title.
func WithConfirmTitle(title string) ConfirmOption {
	return func(c *ConfirmModal) {
		c.Title = title
	}
}

// WithConfirmMessage sets the confirmation message.
func WithConfirmMessage(message string) ConfirmOption {
	return func(c *ConfirmModal) {
		c.Message = message
	}
}

// WithConfirmText sets the confirm button text.
func WithConfirmText(text string) ConfirmOption {
	return func(c *ConfirmModal) {
		c.ConfirmText = text
	}
}

// WithCancelText sets the cancel button text.
func WithCancelText(text string) ConfirmOption {
	return func(c *ConfirmModal) {
		c.CancelText = text
	}
}

// WithConfirmDestructive marks the action as destructive.
func WithConfirmDestructive(destructive bool) ConfirmOption {
	return func(c *ConfirmModal) {
		c.Destructive = destructive
	}
}

// WithConfirmIcon sets the dialog icon.
func WithConfirmIcon(icon string) ConfirmOption {
	return func(c *ConfirmModal) {
		c.Icon = icon
	}
}

// WithConfirmStyled enables Tailwind CSS styling.
func WithConfirmStyled(styled bool) ConfirmOption {
	return func(c *ConfirmModal) {
		c.SetStyled(styled)
	}
}

// SheetOption is a functional option for sheet modals.
type SheetOption func(*SheetModal)

// WithSheetOpen sets the initial open state.
func WithSheetOpen(open bool) SheetOption {
	return func(s *SheetModal) {
		s.Open = open
	}
}

// WithSheetTitle sets the sheet title.
func WithSheetTitle(title string) SheetOption {
	return func(s *SheetModal) {
		s.Title = title
	}
}

// WithSheetPosition sets the sheet position.
func WithSheetPosition(position string) SheetOption {
	return func(s *SheetModal) {
		s.Position = position
	}
}

// WithSheetSize sets the sheet size.
func WithSheetSize(size Size) SheetOption {
	return func(s *SheetModal) {
		s.Size = size
	}
}

// WithSheetShowClose shows or hides the close button.
func WithSheetShowClose(show bool) SheetOption {
	return func(s *SheetModal) {
		s.ShowClose = show
	}
}

// WithSheetCloseOnOverlay enables closing on overlay click.
func WithSheetCloseOnOverlay(close bool) SheetOption {
	return func(s *SheetModal) {
		s.CloseOnOverlay = close
	}
}

// WithSheetStyled enables Tailwind CSS styling.
func WithSheetStyled(styled bool) SheetOption {
	return func(s *SheetModal) {
		s.SetStyled(styled)
	}
}
