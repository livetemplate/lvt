package breadcrumbs

// =============================================================================
// Breadcrumbs Options
// =============================================================================

// WithItems sets the breadcrumb items.
func WithItems(items ...*BreadcrumbItem) Option {
	return func(bc *Breadcrumbs) {
		bc.Items = items
	}
}

// WithSeparator sets the separator style.
func WithSeparator(s Separator) Option {
	return func(bc *Breadcrumbs) {
		bc.Separator = s
	}
}

// WithSize sets the text size.
func WithSize(s Size) Option {
	return func(bc *Breadcrumbs) {
		bc.Size = s
	}
}

// WithShowHome enables the home icon.
func WithShowHome(show bool) Option {
	return func(bc *Breadcrumbs) {
		bc.ShowHome = show
	}
}

// WithHomeHref sets the home link URL.
func WithHomeHref(href string) Option {
	return func(bc *Breadcrumbs) {
		bc.HomeHref = href
	}
}

// WithCollapsible enables collapsing.
func WithCollapsible(collapsible bool) Option {
	return func(bc *Breadcrumbs) {
		bc.Collapsible = collapsible
	}
}

// WithMaxVisible sets max visible items when collapsed.
func WithMaxVisible(max int) Option {
	return func(bc *Breadcrumbs) {
		bc.MaxVisible = max
	}
}

// WithStyled sets the styled mode.
func WithStyled(styled bool) Option {
	return func(bc *Breadcrumbs) {
		bc.SetStyled(styled)
	}
}

// =============================================================================
// BreadcrumbItem Options
// =============================================================================

// WithItemLabel sets the item label.
func WithItemLabel(label string) ItemOption {
	return func(i *BreadcrumbItem) {
		i.Label = label
	}
}

// WithItemHref sets the item href.
func WithItemHref(href string) ItemOption {
	return func(i *BreadcrumbItem) {
		i.Href = href
	}
}

// WithItemIcon sets the item icon.
func WithItemIcon(icon string) ItemOption {
	return func(i *BreadcrumbItem) {
		i.Icon = icon
	}
}

// WithItemCurrent marks the item as current.
func WithItemCurrent(current bool) ItemOption {
	return func(i *BreadcrumbItem) {
		i.Current = current
	}
}

// WithItemDisabled disables the item.
func WithItemDisabled(disabled bool) ItemOption {
	return func(i *BreadcrumbItem) {
		i.Disabled = disabled
	}
}

// WithItemStyled sets the item styled mode.
func WithItemStyled(styled bool) ItemOption {
	return func(i *BreadcrumbItem) {
		i.SetStyled(styled)
	}
}
