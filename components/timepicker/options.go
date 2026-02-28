package timepicker

// Option is a functional option for configuring time pickers.
type Option func(*TimePicker)

// DurationOption is a functional option for configuring duration pickers.
type DurationOption func(*DurationPicker)

// WithPlaceholder sets the placeholder text.
func WithPlaceholder(placeholder string) Option {
	return func(tp *TimePicker) {
		tp.Placeholder = placeholder
	}
}

// WithFormat sets the time format.
func WithFormat(format string) Option {
	return func(tp *TimePicker) {
		tp.Format = format
	}
}

// With24Hour enables 24-hour format.
func With24Hour(use24 bool) Option {
	return func(tp *TimePicker) {
		tp.Use24Hour = use24
	}
}

// WithShowSeconds enables second selection.
func WithShowSeconds(show bool) Option {
	return func(tp *TimePicker) {
		tp.ShowSeconds = show
	}
}

// WithMinuteStep sets the minute increment.
func WithMinuteStep(step int) Option {
	return func(tp *TimePicker) {
		if step > 0 && step <= 60 {
			tp.MinuteStep = step
		}
	}
}

// WithMinTime sets the minimum selectable time.
func WithMinTime(minTime string) Option {
	return func(tp *TimePicker) {
		tp.MinTime = minTime
	}
}

// WithMaxTime sets the maximum selectable time.
func WithMaxTime(maxTime string) Option {
	return func(tp *TimePicker) {
		tp.MaxTime = maxTime
	}
}

// WithTime sets the initial time.
func WithTime(hour, minute int) Option {
	return func(tp *TimePicker) {
		tp.SetTime(hour, minute)
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(tp *TimePicker) {
		tp.SetStyled(styled)
	}
}

// WithOpen sets the initial open state.
func WithOpen(open bool) Option {
	return func(tp *TimePicker) {
		tp.Open = open
	}
}

// Duration picker options

// WithDurationPlaceholder sets the placeholder text.
func WithDurationPlaceholder(placeholder string) DurationOption {
	return func(dp *DurationPicker) {
		dp.Placeholder = placeholder
	}
}

// WithDurationShowSeconds enables second selection.
func WithDurationShowSeconds(show bool) DurationOption {
	return func(dp *DurationPicker) {
		dp.ShowSeconds = show
	}
}

// WithMaxHours sets the maximum hours.
func WithMaxHours(max int) DurationOption {
	return func(dp *DurationPicker) {
		dp.MaxHours = max
	}
}

// WithDuration sets the initial duration.
func WithDuration(hours, minutes int) DurationOption {
	return func(dp *DurationPicker) {
		dp.SetDuration(hours, minutes)
	}
}

// WithDurationStyled enables Tailwind CSS styling.
func WithDurationStyled(styled bool) DurationOption {
	return func(dp *DurationPicker) {
		dp.SetStyled(styled)
	}
}
