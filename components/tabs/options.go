package tabs

// Option is a functional option for configuring tabs.
type Option func(*Tabs)

// WithActive sets the initially active tab by ID.
func WithActive(tabID string) Option {
	return func(t *Tabs) {
		t.SetActive(tabID)
	}
}

// WithStyled enables Tailwind CSS styling for the component.
// When false, renders semantic HTML without styling classes.
func WithStyled(styled bool) Option {
	return func(t *Tabs) {
		t.SetStyled(styled)
	}
}
