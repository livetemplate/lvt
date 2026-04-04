package datepicker

import "time"

// Option is a functional option for configuring date pickers.
type Option func(*DatePicker)

// WithPlaceholder sets the placeholder text.
func WithPlaceholder(placeholder string) Option {
	return func(dp *DatePicker) {
		dp.Placeholder = placeholder
	}
}

// WithFormat sets the date display format.
func WithFormat(format string) Option {
	return func(dp *DatePicker) {
		dp.Format = format
	}
}

// WithSelected sets the initially selected date.
func WithSelected(date time.Time) Option {
	return func(dp *DatePicker) {
		dp.Selected = &date
		dp.ViewDate = date
	}
}

// WithMinDate sets the minimum selectable date.
func WithMinDate(date time.Time) Option {
	return func(dp *DatePicker) {
		dp.MinDate = &date
	}
}

// WithMaxDate sets the maximum selectable date.
func WithMaxDate(date time.Time) Option {
	return func(dp *DatePicker) {
		dp.MaxDate = &date
	}
}

// WithDisabledDates sets specific dates that cannot be selected.
func WithDisabledDates(dates ...time.Time) Option {
	return func(dp *DatePicker) {
		dp.DisabledDates = dates
	}
}

// WithDisabledWeekdays sets days of the week that cannot be selected.
func WithDisabledWeekdays(weekdays ...time.Weekday) Option {
	return func(dp *DatePicker) {
		dp.DisabledWeekdays = weekdays
	}
}

// WithFirstDayOfWeek sets the first day of the week (0=Sunday, 1=Monday, etc.).
func WithFirstDayOfWeek(day int) Option {
	return func(dp *DatePicker) {
		dp.FirstDayOfWeek = day % 7
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) Option {
	return func(dp *DatePicker) {
		dp.SetStyled(styled)
	}
}

// WithOpen is a no-op, retained for backward compatibility.
// Open/close state is now managed client-side via CSS classes.
func WithOpen(_ bool) Option {
	return func(_ *DatePicker) {}
}
