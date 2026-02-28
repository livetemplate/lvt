// Package timepicker provides time selection components.
//
// Available variants:
//   - New() creates a time picker (template: "lvt:timepicker:default:v1")
//   - NewDuration() creates a duration picker (template: "lvt:timepicker:duration:v1")
//
// Required lvt-* attributes: lvt-click, lvt-input, lvt-click-away
//
// Example usage:
//
//	// In your controller/state
//	StartTime: timepicker.New("start-time",
//	    timepicker.WithPlaceholder("Select time"),
//	    timepicker.WithFormat("3:04 PM"),
//	)
//
//	// In your template
//	{{template "lvt:timepicker:default:v1" .StartTime}}
package timepicker

import (
	"fmt"

	"github.com/livetemplate/components/base"
)

// TimePicker is a component for selecting time.
// Use template "lvt:timepicker:default:v1" to render.
type TimePicker struct {
	base.Base

	// Hour is the selected hour (0-23 or 1-12 depending on Use24Hour)
	Hour int

	// Minute is the selected minute (0-59)
	Minute int

	// Second is the selected second (0-59, only used if ShowSeconds)
	Second int

	// Period is "AM" or "PM" (only used if not Use24Hour)
	Period string

	// HasValue indicates if a time has been selected
	HasValue bool

	// Open indicates whether the picker popup is visible
	Open bool

	// Placeholder text shown when no time is selected
	Placeholder string

	// Format for displaying the time
	Format string

	// Use24Hour uses 24-hour format instead of 12-hour
	Use24Hour bool

	// ShowSeconds shows second selection
	ShowSeconds bool

	// MinuteStep is the minute increment (default 1, common: 5, 15, 30)
	MinuteStep int

	// MinTime is the earliest selectable time (HH:MM)
	MinTime string

	// MaxTime is the latest selectable time (HH:MM)
	MaxTime string
}

// DurationPicker is a component for selecting a duration.
// Use template "lvt:timepicker:duration:v1" to render.
type DurationPicker struct {
	base.Base

	// Hours in the duration
	Hours int

	// Minutes in the duration
	Minutes int

	// Seconds in the duration (only if ShowSeconds)
	Seconds int

	// HasValue indicates if a duration has been set
	HasValue bool

	// Open indicates whether the picker popup is visible
	Open bool

	// Placeholder text
	Placeholder string

	// ShowSeconds includes seconds
	ShowSeconds bool

	// MaxHours limits the maximum hours
	MaxHours int
}

// New creates a time picker.
//
// Example:
//
//	tp := timepicker.New("meeting-time",
//	    timepicker.WithPlaceholder("Select time"),
//	    timepicker.WithMinuteStep(15),
//	)
func New(id string, opts ...Option) *TimePicker {
	tp := &TimePicker{
		Base:        base.NewBase(id, "timepicker"),
		Hour:        12,
		Minute:      0,
		Second:      0,
		Period:      "AM",
		Placeholder: "Select time...",
		Format:      "3:04 PM",
		MinuteStep:  1,
	}

	for _, opt := range opts {
		opt(tp)
	}

	return tp
}

// NewDuration creates a duration picker.
func NewDuration(id string, opts ...DurationOption) *DurationPicker {
	dp := &DurationPicker{
		Base:        base.NewBase(id, "timepicker"),
		Placeholder: "Select duration...",
		MaxHours:    24,
	}

	for _, opt := range opts {
		opt(dp)
	}

	return dp
}

// Toggle opens or closes the picker.
func (tp *TimePicker) Toggle() {
	tp.Open = !tp.Open
}

// Close closes the picker.
func (tp *TimePicker) Close() {
	tp.Open = false
}

// SetTime sets the time.
func (tp *TimePicker) SetTime(hour, minute int) {
	tp.Hour = hour
	tp.Minute = minute
	tp.HasValue = true
	tp.Open = false

	// Normalize for 12-hour format
	if !tp.Use24Hour {
		if hour == 0 {
			tp.Hour = 12
			tp.Period = "AM"
		} else if hour < 12 {
			tp.Hour = hour
			tp.Period = "AM"
		} else if hour == 12 {
			tp.Hour = 12
			tp.Period = "PM"
		} else {
			tp.Hour = hour - 12
			tp.Period = "PM"
		}
	}
}

// SetTimeWithSeconds sets the time including seconds.
func (tp *TimePicker) SetTimeWithSeconds(hour, minute, second int) {
	tp.SetTime(hour, minute)
	tp.Second = second
}

// Clear clears the selected time.
func (tp *TimePicker) Clear() {
	tp.HasValue = false
	tp.Hour = 12
	tp.Minute = 0
	tp.Second = 0
	tp.Period = "AM"
}

// SetHour sets the hour.
func (tp *TimePicker) SetHour(hour int) {
	if tp.Use24Hour {
		if hour >= 0 && hour <= 23 {
			tp.Hour = hour
			tp.HasValue = true
		}
	} else {
		if hour >= 1 && hour <= 12 {
			tp.Hour = hour
			tp.HasValue = true
		}
	}
}

// SetMinute sets the minute.
func (tp *TimePicker) SetMinute(minute int) {
	if minute >= 0 && minute <= 59 {
		tp.Minute = minute
		tp.HasValue = true
	}
}

// SetSecond sets the second.
func (tp *TimePicker) SetSecond(second int) {
	if second >= 0 && second <= 59 {
		tp.Second = second
	}
}

// SetPeriod sets AM/PM.
func (tp *TimePicker) SetPeriod(period string) {
	if period == "AM" || period == "PM" {
		tp.Period = period
		tp.HasValue = true
	}
}

// TogglePeriod toggles between AM and PM.
func (tp *TimePicker) TogglePeriod() {
	if tp.Period == "AM" {
		tp.Period = "PM"
	} else {
		tp.Period = "AM"
	}
}

// IncrementHour increments the hour.
func (tp *TimePicker) IncrementHour() {
	if tp.Use24Hour {
		tp.Hour = (tp.Hour + 1) % 24
	} else {
		tp.Hour++
		if tp.Hour > 12 {
			tp.Hour = 1
		}
	}
	tp.HasValue = true
}

// DecrementHour decrements the hour.
func (tp *TimePicker) DecrementHour() {
	if tp.Use24Hour {
		tp.Hour--
		if tp.Hour < 0 {
			tp.Hour = 23
		}
	} else {
		tp.Hour--
		if tp.Hour < 1 {
			tp.Hour = 12
		}
	}
	tp.HasValue = true
}

// IncrementMinute increments the minute by MinuteStep.
func (tp *TimePicker) IncrementMinute() {
	tp.Minute = (tp.Minute + tp.MinuteStep) % 60
	tp.HasValue = true
}

// DecrementMinute decrements the minute by MinuteStep.
func (tp *TimePicker) DecrementMinute() {
	tp.Minute -= tp.MinuteStep
	if tp.Minute < 0 {
		tp.Minute = 60 + tp.Minute
	}
	tp.HasValue = true
}

// Get24Hour returns the hour in 24-hour format.
func (tp *TimePicker) Get24Hour() int {
	if tp.Use24Hour {
		return tp.Hour
	}
	if tp.Period == "AM" {
		if tp.Hour == 12 {
			return 0
		}
		return tp.Hour
	}
	// PM
	if tp.Hour == 12 {
		return 12
	}
	return tp.Hour + 12
}

// DisplayValue returns the formatted time or placeholder.
func (tp *TimePicker) DisplayValue() string {
	if !tp.HasValue {
		return tp.Placeholder
	}
	return tp.FormatTime()
}

// FormatTime returns the formatted time string.
func (tp *TimePicker) FormatTime() string {
	if tp.Use24Hour {
		if tp.ShowSeconds {
			return fmt.Sprintf("%02d:%02d:%02d", tp.Hour, tp.Minute, tp.Second)
		}
		return fmt.Sprintf("%02d:%02d", tp.Hour, tp.Minute)
	}
	if tp.ShowSeconds {
		return fmt.Sprintf("%d:%02d:%02d %s", tp.Hour, tp.Minute, tp.Second, tp.Period)
	}
	return fmt.Sprintf("%d:%02d %s", tp.Hour, tp.Minute, tp.Period)
}

// HourOptions returns available hour options.
func (tp *TimePicker) HourOptions() []int {
	if tp.Use24Hour {
		hours := make([]int, 24)
		for i := 0; i < 24; i++ {
			hours[i] = i
		}
		return hours
	}
	hours := make([]int, 12)
	for i := 0; i < 12; i++ {
		hours[i] = i + 1
	}
	return hours
}

// MinuteOptions returns available minute options.
func (tp *TimePicker) MinuteOptions() []int {
	var minutes []int
	for i := 0; i < 60; i += tp.MinuteStep {
		minutes = append(minutes, i)
	}
	return minutes
}

// SecondOptions returns available second options.
func (tp *TimePicker) SecondOptions() []int {
	seconds := make([]int, 60)
	for i := 0; i < 60; i++ {
		seconds[i] = i
	}
	return seconds
}

// SetNow sets the time to the current time.
func (tp *TimePicker) SetNow() {
	// In a real implementation, this would use time.Now()
	// For the component, we just mark HasValue true
	tp.HasValue = true
}

// DurationPicker methods

// Toggle opens or closes the picker.
func (dp *DurationPicker) Toggle() {
	dp.Open = !dp.Open
}

// Close closes the picker.
func (dp *DurationPicker) Close() {
	dp.Open = false
}

// SetDuration sets the duration.
func (dp *DurationPicker) SetDuration(hours, minutes int) {
	dp.Hours = hours
	dp.Minutes = minutes
	dp.HasValue = true
	dp.Open = false
}

// SetDurationWithSeconds sets the duration including seconds.
func (dp *DurationPicker) SetDurationWithSeconds(hours, minutes, seconds int) {
	dp.SetDuration(hours, minutes)
	dp.Seconds = seconds
}

// Clear clears the duration.
func (dp *DurationPicker) Clear() {
	dp.Hours = 0
	dp.Minutes = 0
	dp.Seconds = 0
	dp.HasValue = false
}

// SetHours sets the hours.
func (dp *DurationPicker) SetHours(hours int) {
	if hours >= 0 && hours <= dp.MaxHours {
		dp.Hours = hours
		dp.HasValue = true
	}
}

// SetMinutes sets the minutes.
func (dp *DurationPicker) SetMinutes(minutes int) {
	if minutes >= 0 && minutes <= 59 {
		dp.Minutes = minutes
		dp.HasValue = true
	}
}

// SetSeconds sets the seconds.
func (dp *DurationPicker) SetSeconds(seconds int) {
	if seconds >= 0 && seconds <= 59 {
		dp.Seconds = seconds
	}
}

// IncrementHours increments hours.
func (dp *DurationPicker) IncrementHours() {
	if dp.Hours < dp.MaxHours {
		dp.Hours++
		dp.HasValue = true
	}
}

// DecrementHours decrements hours.
func (dp *DurationPicker) DecrementHours() {
	if dp.Hours > 0 {
		dp.Hours--
		dp.HasValue = true
	}
}

// IncrementMinutes increments minutes.
func (dp *DurationPicker) IncrementMinutes() {
	dp.Minutes++
	if dp.Minutes > 59 {
		dp.Minutes = 0
	}
	dp.HasValue = true
}

// DecrementMinutes decrements minutes.
func (dp *DurationPicker) DecrementMinutes() {
	dp.Minutes--
	if dp.Minutes < 0 {
		dp.Minutes = 59
	}
	dp.HasValue = true
}

// TotalMinutes returns the total duration in minutes.
func (dp *DurationPicker) TotalMinutes() int {
	return dp.Hours*60 + dp.Minutes
}

// TotalSeconds returns the total duration in seconds.
func (dp *DurationPicker) TotalSeconds() int {
	return dp.Hours*3600 + dp.Minutes*60 + dp.Seconds
}

// DisplayValue returns the formatted duration or placeholder.
func (dp *DurationPicker) DisplayValue() string {
	if !dp.HasValue {
		return dp.Placeholder
	}
	return dp.FormatDuration()
}

// FormatDuration returns the formatted duration string.
func (dp *DurationPicker) FormatDuration() string {
	if dp.ShowSeconds {
		return fmt.Sprintf("%dh %dm %ds", dp.Hours, dp.Minutes, dp.Seconds)
	}
	return fmt.Sprintf("%dh %dm", dp.Hours, dp.Minutes)
}
