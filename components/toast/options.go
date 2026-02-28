package toast

// ContainerOption is a functional option for configuring toast containers.
type ContainerOption func(*Container)

// WithPosition sets the position of the toast container.
func WithPosition(pos Position) ContainerOption {
	return func(c *Container) {
		c.Position = pos
	}
}

// WithMaxVisible limits how many toasts are shown at once.
func WithMaxVisible(max int) ContainerOption {
	return func(c *Container) {
		c.MaxVisible = max
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) ContainerOption {
	return func(c *Container) {
		c.SetStyled(styled)
	}
}

// MessageOption is a functional option for configuring toast messages.
type MessageOption func(*Message)

// WithTitle sets the toast title.
func WithTitle(title string) MessageOption {
	return func(m *Message) {
		m.Title = title
	}
}

// WithBody sets the toast body text.
func WithBody(body string) MessageOption {
	return func(m *Message) {
		m.Body = body
	}
}

// WithType sets the toast type (info, success, warning, error).
func WithType(t Type) MessageOption {
	return func(m *Message) {
		m.Type = t
	}
}

// WithDismissible sets whether the toast can be dismissed.
func WithDismissible(dismissible bool) MessageOption {
	return func(m *Message) {
		m.Dismissible = dismissible
	}
}

// WithIcon sets a custom icon for the toast.
func WithIcon(icon string) MessageOption {
	return func(m *Message) {
		m.Icon = icon
	}
}

// NewMessage creates a new toast message with options.
func NewMessage(opts ...MessageOption) Message {
	m := Message{
		Type:        Info,
		Dismissible: true,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}
