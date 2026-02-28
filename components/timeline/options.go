package timeline

// =============================================================================
// Timeline Options
// =============================================================================

// WithItems sets the timeline items.
func WithItems(items ...*TimelineItem) Option {
	return func(t *Timeline) {
		t.Items = items
	}
}

// WithOrientation sets the timeline orientation.
func WithOrientation(o Orientation) Option {
	return func(t *Timeline) {
		t.Orientation = o
	}
}

// WithPosition sets the content position.
func WithPosition(p Position) Option {
	return func(t *Timeline) {
		t.Position = p
	}
}

// WithShowConnectors enables/disables connectors.
func WithShowConnectors(show bool) Option {
	return func(t *Timeline) {
		t.ShowConnectors = show
	}
}

// WithReverse reverses the item order.
func WithReverse(reverse bool) Option {
	return func(t *Timeline) {
		t.Reverse = reverse
	}
}

// WithStyled sets the styled mode.
func WithStyled(styled bool) Option {
	return func(t *Timeline) {
		t.SetStyled(styled)
	}
}

// =============================================================================
// TimelineItem Options
// =============================================================================

// WithItemTitle sets the item title.
func WithItemTitle(title string) ItemOption {
	return func(i *TimelineItem) {
		i.Title = title
	}
}

// WithItemDescription sets the item description.
func WithItemDescription(desc string) ItemOption {
	return func(i *TimelineItem) {
		i.Description = desc
	}
}

// WithItemTime sets the item time/date.
func WithItemTime(time string) ItemOption {
	return func(i *TimelineItem) {
		i.Time = time
	}
}

// WithItemIcon sets the item icon.
func WithItemIcon(icon string) ItemOption {
	return func(i *TimelineItem) {
		i.Icon = icon
	}
}

// WithItemStatus sets the item status.
func WithItemStatus(status Status) ItemOption {
	return func(i *TimelineItem) {
		i.Status = status
	}
}

// WithItemColor sets the item color.
func WithItemColor(color Color) ItemOption {
	return func(i *TimelineItem) {
		i.Color = color
	}
}

// WithItemActive sets the item as active.
func WithItemActive(active bool) ItemOption {
	return func(i *TimelineItem) {
		i.Active = active
	}
}

// WithItemCompleted sets the item as completed.
func WithItemCompleted(completed bool) ItemOption {
	return func(i *TimelineItem) {
		i.Completed = completed
	}
}

// WithItemStyled sets the item styled mode.
func WithItemStyled(styled bool) ItemOption {
	return func(i *TimelineItem) {
		i.SetStyled(styled)
	}
}
