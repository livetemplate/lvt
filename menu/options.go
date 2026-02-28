package menu

// Option is a functional option for configuring menu components.
type Option func(*Menu)

// NavOption is a functional option for configuring navigation menus.
type NavOption func(*NavMenu)

// WithItems sets the menu items.
func WithItems(items []Item) Option {
	return func(m *Menu) {
		m.Items = items
	}
}

// WithTrigger sets the trigger button text.
func WithTrigger(trigger string) Option {
	return func(m *Menu) {
		m.Trigger = trigger
	}
}

// WithTriggerIcon sets the trigger button icon.
func WithTriggerIcon(icon string) Option {
	return func(m *Menu) {
		m.TriggerIcon = icon
	}
}

// WithPosition sets the menu position relative to trigger.
// Options: "bottom-left", "bottom-right", "top-left", "top-right"
func WithPosition(position string) Option {
	return func(m *Menu) {
		m.Position = position
	}
}

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(m *Menu) {
		m.Open = open
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) Option {
	return func(m *Menu) {
		m.SetStyled(styled)
	}
}

// NavMenu options

// WithNavItems sets the navigation items.
func WithNavItems(items []Item) NavOption {
	return func(nm *NavMenu) {
		nm.Items = items
	}
}

// WithOrientation sets the menu orientation.
// Options: "horizontal", "vertical"
func WithOrientation(orientation string) NavOption {
	return func(nm *NavMenu) {
		nm.Orientation = orientation
	}
}

// WithActiveID sets the active item ID.
func WithActiveID(id string) NavOption {
	return func(nm *NavMenu) {
		nm.SetActive(id)
	}
}

// WithNavStyled enables Tailwind CSS styling for navigation menu.
func WithNavStyled(styled bool) NavOption {
	return func(nm *NavMenu) {
		nm.SetStyled(styled)
	}
}
