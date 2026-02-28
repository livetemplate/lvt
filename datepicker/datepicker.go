// Package datepicker provides date selection components with calendar UI.
//
// Available variants:
//   - New() creates a single date picker (template: "lvt:datepicker:single:v1")
//   - NewRange() creates a date range picker (template: "lvt:datepicker:range:v1")
//   - NewInline() creates an inline calendar (template: "lvt:datepicker:inline:v1")
//
// Required lvt-* attributes: lvt-click
//
// Example usage:
//
//	// In your controller/state
//	BirthDate: datepicker.New("birthdate",
//	    datepicker.WithPlaceholder("Select date"),
//	)
//
//	// In your template
//	{{template "lvt:datepicker:single:v1" .BirthDate}}
package datepicker

import (
	"time"

	"github.com/livetemplate/components/base"
)

// DatePicker is a component for selecting a single date.
// Use template "lvt:datepicker:single:v1" to render.
type DatePicker struct {
	base.Base

	// Selected is the currently selected date (nil if none)
	Selected *time.Time

	// ViewDate is the date currently being viewed (for navigation)
	ViewDate time.Time

	// Placeholder text shown when no date is selected
	Placeholder string

	// Open indicates whether the calendar popup is visible
	Open bool

	// MinDate is the earliest selectable date (nil for no limit)
	MinDate *time.Time

	// MaxDate is the latest selectable date (nil for no limit)
	MaxDate *time.Time

	// DisabledDates are specific dates that cannot be selected
	DisabledDates []time.Time

	// DisabledWeekdays are days of the week that cannot be selected (0=Sunday, 6=Saturday)
	DisabledWeekdays []time.Weekday

	// Format for displaying the date
	Format string

	// FirstDayOfWeek (0=Sunday, 1=Monday, etc.)
	FirstDayOfWeek int
}

// RangePicker is a component for selecting a date range.
// Use template "lvt:datepicker:range:v1" to render.
type RangePicker struct {
	DatePicker

	// StartDate is the beginning of the selected range
	StartDate *time.Time

	// EndDate is the end of the selected range
	EndDate *time.Time

	// SelectingEnd indicates we're selecting the end date
	SelectingEnd bool
}

// New creates a single date picker.
//
// Example:
//
//	dp := datepicker.New("birthdate",
//	    datepicker.WithPlaceholder("Select date"),
//	    datepicker.WithFormat("Jan 2, 2006"),
//	)
func New(id string, opts ...Option) *DatePicker {
	dp := &DatePicker{
		Base:           base.NewBase(id, "datepicker"),
		ViewDate:       time.Now(),
		Placeholder:    "Select date...",
		Format:         "Jan 2, 2006",
		FirstDayOfWeek: 0, // Sunday
	}

	for _, opt := range opts {
		opt(dp)
	}

	return dp
}

// NewRange creates a date range picker.
func NewRange(id string, opts ...Option) *RangePicker {
	dp := New(id, opts...)
	return &RangePicker{
		DatePicker: *dp,
	}
}

// NewInline creates an inline calendar (always visible).
func NewInline(id string, opts ...Option) *DatePicker {
	dp := New(id, opts...)
	dp.Open = true
	return dp
}

// Toggle opens or closes the calendar.
func (dp *DatePicker) Toggle() {
	dp.Open = !dp.Open
}

// Close closes the calendar.
func (dp *DatePicker) Close() {
	dp.Open = false
}

// SelectDate selects a date.
func (dp *DatePicker) SelectDate(date time.Time) bool {
	if !dp.IsDateSelectable(date) {
		return false
	}
	dp.Selected = &date
	dp.Open = false
	return true
}

// Clear clears the selected date.
func (dp *DatePicker) Clear() {
	dp.Selected = nil
}

// PreviousMonth navigates to the previous month.
func (dp *DatePicker) PreviousMonth() {
	dp.ViewDate = dp.ViewDate.AddDate(0, -1, 0)
}

// NextMonth navigates to the next month.
func (dp *DatePicker) NextMonth() {
	dp.ViewDate = dp.ViewDate.AddDate(0, 1, 0)
}

// PreviousYear navigates to the previous year.
func (dp *DatePicker) PreviousYear() {
	dp.ViewDate = dp.ViewDate.AddDate(-1, 0, 0)
}

// NextYear navigates to the next year.
func (dp *DatePicker) NextYear() {
	dp.ViewDate = dp.ViewDate.AddDate(1, 0, 0)
}

// GoToToday navigates to and selects today.
func (dp *DatePicker) GoToToday() {
	today := time.Now()
	dp.ViewDate = today
}

// IsDateSelectable checks if a date can be selected.
func (dp *DatePicker) IsDateSelectable(date time.Time) bool {
	// Check min date
	if dp.MinDate != nil && date.Before(*dp.MinDate) {
		return false
	}

	// Check max date
	if dp.MaxDate != nil && date.After(*dp.MaxDate) {
		return false
	}

	// Check disabled weekdays
	for _, wd := range dp.DisabledWeekdays {
		if date.Weekday() == wd {
			return false
		}
	}

	// Check disabled dates
	for _, dd := range dp.DisabledDates {
		if sameDay(date, dd) {
			return false
		}
	}

	return true
}

// IsSelected checks if a date is the selected date.
func (dp *DatePicker) IsSelected(date time.Time) bool {
	if dp.Selected == nil {
		return false
	}
	return sameDay(*dp.Selected, date)
}

// IsToday checks if a date is today.
func (dp *DatePicker) IsToday(date time.Time) bool {
	return sameDay(date, time.Now())
}

// DisplayValue returns the formatted selected date or placeholder.
func (dp *DatePicker) DisplayValue() string {
	if dp.Selected == nil {
		return dp.Placeholder
	}
	return dp.Selected.Format(dp.Format)
}

// ViewMonth returns the month being viewed.
func (dp *DatePicker) ViewMonth() string {
	return dp.ViewDate.Format("January")
}

// ViewYear returns the year being viewed.
func (dp *DatePicker) ViewYear() int {
	return dp.ViewDate.Year()
}

// CalendarWeeks returns the weeks for the current view month.
func (dp *DatePicker) CalendarWeeks() [][]CalendarDay {
	year, month, _ := dp.ViewDate.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, dp.ViewDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Find start day (might be in previous month)
	startOffset := int(firstOfMonth.Weekday()) - dp.FirstDayOfWeek
	if startOffset < 0 {
		startOffset += 7
	}
	startDate := firstOfMonth.AddDate(0, 0, -startOffset)

	var weeks [][]CalendarDay

	current := startDate
	for {
		var week []CalendarDay
		for i := 0; i < 7; i++ {
			day := CalendarDay{
				Date:        current,
				Day:         current.Day(),
				InMonth:     current.Month() == month,
				IsToday:     dp.IsToday(current),
				IsSelected:  dp.IsSelected(current),
				IsDisabled:  !dp.IsDateSelectable(current),
			}
			week = append(week, day)
			current = current.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)

		// Stop if we've passed the last day of the month and completed the week
		if current.After(lastOfMonth) && len(weeks) >= 4 {
			break
		}
		// Safety limit
		if len(weeks) >= 6 {
			break
		}
	}

	return weeks
}

// WeekdayNames returns the names of weekdays starting from FirstDayOfWeek.
func (dp *DatePicker) WeekdayNames() []string {
	names := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	result := make([]string, 7)
	for i := 0; i < 7; i++ {
		result[i] = names[(dp.FirstDayOfWeek+i)%7]
	}
	return result
}

// CalendarDay represents a day in the calendar view.
type CalendarDay struct {
	Date       time.Time
	Day        int
	InMonth    bool
	IsToday    bool
	IsSelected bool
	IsDisabled bool
}

// DateString returns the date as a string (for lvt-data attributes).
func (d CalendarDay) DateString() string {
	return d.Date.Format("2006-01-02")
}

// RangePicker methods

// SelectRangeDate handles date selection for range picker.
func (rp *RangePicker) SelectRangeDate(date time.Time) bool {
	if !rp.IsDateSelectable(date) {
		return false
	}

	if !rp.SelectingEnd || rp.StartDate == nil {
		// Selecting start date
		rp.StartDate = &date
		rp.EndDate = nil
		rp.SelectingEnd = true
	} else {
		// Selecting end date
		if date.Before(*rp.StartDate) {
			// Swap if end is before start
			rp.EndDate = rp.StartDate
			rp.StartDate = &date
		} else {
			rp.EndDate = &date
		}
		rp.SelectingEnd = false
		rp.Open = false
	}
	return true
}

// ClearRange clears both dates.
func (rp *RangePicker) ClearRange() {
	rp.StartDate = nil
	rp.EndDate = nil
	rp.SelectingEnd = false
}

// IsInRange checks if a date is within the selected range.
func (rp *RangePicker) IsInRange(date time.Time) bool {
	if rp.StartDate == nil || rp.EndDate == nil {
		return false
	}
	return (date.After(*rp.StartDate) || sameDay(date, *rp.StartDate)) &&
		(date.Before(*rp.EndDate) || sameDay(date, *rp.EndDate))
}

// DisplayRangeValue returns the formatted range or placeholder.
func (rp *RangePicker) DisplayRangeValue() string {
	if rp.StartDate == nil {
		return rp.Placeholder
	}
	start := rp.StartDate.Format(rp.Format)
	if rp.EndDate == nil {
		return start + " - ..."
	}
	return start + " - " + rp.EndDate.Format(rp.Format)
}

// Helper functions
func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
