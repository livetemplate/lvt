package accordion

// Option is a functional option for configuring accordions.
type Option func(*Accordion)

// WithOpen sets initially open item(s) by ID.
func WithOpen(itemIDs ...string) Option {
	return func(a *Accordion) {
		for _, id := range itemIDs {
			a.Open(id)
		}
	}
}

// WithAllOpen opens all items initially.
func WithAllOpen() Option {
	return func(a *Accordion) {
		a.OpenAll()
	}
}

// WithStyled enables Tailwind CSS styling for the component.
// When false, renders semantic HTML without styling classes.
func WithStyled(styled bool) Option {
	return func(a *Accordion) {
		a.SetStyled(styled)
	}
}
