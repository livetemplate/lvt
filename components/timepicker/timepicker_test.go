package timepicker

import (
	"testing"
)

func TestNew(t *testing.T) {
	tp := New("test-tp")

	if tp.ID() != "test-tp" {
		t.Errorf("Expected ID 'test-tp', got '%s'", tp.ID())
	}
	if tp.Namespace() != "timepicker" {
		t.Errorf("Expected namespace 'timepicker', got '%s'", tp.Namespace())
	}
	if tp.Hour != 12 {
		t.Errorf("Expected default Hour 12, got %d", tp.Hour)
	}
	if tp.Minute != 0 {
		t.Errorf("Expected default Minute 0, got %d", tp.Minute)
	}
	if tp.Period != "AM" {
		t.Errorf("Expected default Period 'AM', got '%s'", tp.Period)
	}
	if tp.MinuteStep != 1 {
		t.Errorf("Expected default MinuteStep 1, got %d", tp.MinuteStep)
	}
	if tp.HasValue {
		t.Error("Expected HasValue to be false")
	}
}

func TestNewDuration(t *testing.T) {
	dp := NewDuration("test-dur")

	if dp.ID() != "test-dur" {
		t.Errorf("Expected ID 'test-dur', got '%s'", dp.ID())
	}
	if dp.MaxHours != 24 {
		t.Errorf("Expected default MaxHours 24, got %d", dp.MaxHours)
	}
}

func TestWithPlaceholder(t *testing.T) {
	tp := New("test", WithPlaceholder("Pick time"))
	if tp.Placeholder != "Pick time" {
		t.Errorf("Expected placeholder 'Pick time', got '%s'", tp.Placeholder)
	}
}

func TestWith24Hour(t *testing.T) {
	tp := New("test", With24Hour(true))
	if !tp.Use24Hour {
		t.Error("Expected Use24Hour to be true")
	}
}

func TestWithShowSeconds(t *testing.T) {
	tp := New("test", WithShowSeconds(true))
	if !tp.ShowSeconds {
		t.Error("Expected ShowSeconds to be true")
	}
}

func TestWithMinuteStep(t *testing.T) {
	tp := New("test", WithMinuteStep(15))
	if tp.MinuteStep != 15 {
		t.Errorf("Expected MinuteStep 15, got %d", tp.MinuteStep)
	}

	// Invalid step ignored
	tp2 := New("test2", WithMinuteStep(0))
	if tp2.MinuteStep != 1 {
		t.Error("Expected invalid MinuteStep to be ignored")
	}
}

func TestWithTime(t *testing.T) {
	tp := New("test", WithTime(14, 30))

	if !tp.HasValue {
		t.Error("Expected HasValue to be true")
	}
	// 14:30 in 12-hour is 2:30 PM
	if tp.Hour != 2 || tp.Period != "PM" {
		t.Errorf("Expected 2 PM, got %d %s", tp.Hour, tp.Period)
	}
}

func TestWithStyled(t *testing.T) {
	tp := New("test", WithStyled(false))
	if tp.IsStyled() {
		t.Error("Expected IsStyled to be false")
	}
}

func TestToggle(t *testing.T) {
	tp := New("test")

	tp.Toggle()
	if !tp.Open {
		t.Error("Expected Open to be true after toggle")
	}

	tp.Toggle()
	if tp.Open {
		t.Error("Expected Open to be false after second toggle")
	}
}

func TestClose(t *testing.T) {
	tp := New("test", WithOpen(true))
	tp.Close()
	if tp.Open {
		t.Error("Expected Open to be false after Close")
	}
}

func TestSetTime(t *testing.T) {
	tp := New("test")
	tp.SetTime(14, 30)

	if !tp.HasValue {
		t.Error("Expected HasValue to be true")
	}
	if tp.Hour != 2 || tp.Minute != 30 || tp.Period != "PM" {
		t.Errorf("Expected 2:30 PM, got %d:%d %s", tp.Hour, tp.Minute, tp.Period)
	}
}

func TestSetTimeMidnight(t *testing.T) {
	tp := New("test")
	tp.SetTime(0, 0)

	if tp.Hour != 12 || tp.Period != "AM" {
		t.Errorf("Expected 12:00 AM for midnight, got %d:%d %s", tp.Hour, tp.Minute, tp.Period)
	}
}

func TestSetTimeNoon(t *testing.T) {
	tp := New("test")
	tp.SetTime(12, 0)

	if tp.Hour != 12 || tp.Period != "PM" {
		t.Errorf("Expected 12:00 PM for noon, got %d:%d %s", tp.Hour, tp.Minute, tp.Period)
	}
}

func TestSetTime24Hour(t *testing.T) {
	tp := New("test", With24Hour(true))
	tp.SetTime(14, 30)

	if tp.Hour != 14 {
		t.Errorf("Expected Hour 14 in 24-hour mode, got %d", tp.Hour)
	}
}

func TestClear(t *testing.T) {
	tp := New("test", WithTime(14, 30))
	tp.Clear()

	if tp.HasValue {
		t.Error("Expected HasValue to be false")
	}
	if tp.Hour != 12 || tp.Minute != 0 || tp.Period != "AM" {
		t.Error("Expected time to be reset to defaults")
	}
}

func TestSetHour(t *testing.T) {
	tp := New("test")
	tp.SetHour(5)

	if tp.Hour != 5 {
		t.Errorf("Expected Hour 5, got %d", tp.Hour)
	}
	if !tp.HasValue {
		t.Error("Expected HasValue to be true")
	}
}

func TestSetHour24(t *testing.T) {
	tp := New("test", With24Hour(true))
	tp.SetHour(23)

	if tp.Hour != 23 {
		t.Errorf("Expected Hour 23, got %d", tp.Hour)
	}
}

func TestSetMinute(t *testing.T) {
	tp := New("test")
	tp.SetMinute(45)

	if tp.Minute != 45 {
		t.Errorf("Expected Minute 45, got %d", tp.Minute)
	}
}

func TestSetSecond(t *testing.T) {
	tp := New("test")
	tp.SetSecond(30)

	if tp.Second != 30 {
		t.Errorf("Expected Second 30, got %d", tp.Second)
	}
}

func TestSetPeriod(t *testing.T) {
	tp := New("test")
	tp.SetPeriod("PM")

	if tp.Period != "PM" {
		t.Errorf("Expected Period 'PM', got '%s'", tp.Period)
	}
}

func TestTogglePeriod(t *testing.T) {
	tp := New("test")

	tp.TogglePeriod()
	if tp.Period != "PM" {
		t.Errorf("Expected Period 'PM', got '%s'", tp.Period)
	}

	tp.TogglePeriod()
	if tp.Period != "AM" {
		t.Errorf("Expected Period 'AM', got '%s'", tp.Period)
	}
}

func TestIncrementHour(t *testing.T) {
	tp := New("test")
	tp.Hour = 11

	tp.IncrementHour()
	if tp.Hour != 12 {
		t.Errorf("Expected Hour 12, got %d", tp.Hour)
	}

	tp.IncrementHour()
	if tp.Hour != 1 {
		t.Errorf("Expected Hour 1 (wrap), got %d", tp.Hour)
	}
}

func TestIncrementHour24(t *testing.T) {
	tp := New("test", With24Hour(true))
	tp.Hour = 23

	tp.IncrementHour()
	if tp.Hour != 0 {
		t.Errorf("Expected Hour 0 (wrap), got %d", tp.Hour)
	}
}

func TestDecrementHour(t *testing.T) {
	tp := New("test")
	tp.Hour = 1

	tp.DecrementHour()
	if tp.Hour != 12 {
		t.Errorf("Expected Hour 12 (wrap), got %d", tp.Hour)
	}
}

func TestDecrementHour24(t *testing.T) {
	tp := New("test", With24Hour(true))
	tp.Hour = 0

	tp.DecrementHour()
	if tp.Hour != 23 {
		t.Errorf("Expected Hour 23 (wrap), got %d", tp.Hour)
	}
}

func TestIncrementMinute(t *testing.T) {
	tp := New("test", WithMinuteStep(15))
	tp.Minute = 45

	tp.IncrementMinute()
	if tp.Minute != 0 {
		t.Errorf("Expected Minute 0 (wrap), got %d", tp.Minute)
	}
}

func TestDecrementMinute(t *testing.T) {
	tp := New("test", WithMinuteStep(15))
	tp.Minute = 0

	tp.DecrementMinute()
	if tp.Minute != 45 {
		t.Errorf("Expected Minute 45 (wrap), got %d", tp.Minute)
	}
}

func TestGet24Hour(t *testing.T) {
	tp := New("test")

	// 12 AM = 0
	tp.Hour = 12
	tp.Period = "AM"
	if tp.Get24Hour() != 0 {
		t.Errorf("Expected 0 for 12 AM, got %d", tp.Get24Hour())
	}

	// 12 PM = 12
	tp.Period = "PM"
	if tp.Get24Hour() != 12 {
		t.Errorf("Expected 12 for 12 PM, got %d", tp.Get24Hour())
	}

	// 5 AM = 5
	tp.Hour = 5
	tp.Period = "AM"
	if tp.Get24Hour() != 5 {
		t.Errorf("Expected 5 for 5 AM, got %d", tp.Get24Hour())
	}

	// 5 PM = 17
	tp.Period = "PM"
	if tp.Get24Hour() != 17 {
		t.Errorf("Expected 17 for 5 PM, got %d", tp.Get24Hour())
	}
}

func TestDisplayValue(t *testing.T) {
	tp := New("test", WithPlaceholder("Pick time"))

	if tp.DisplayValue() != "Pick time" {
		t.Errorf("Expected placeholder, got '%s'", tp.DisplayValue())
	}

	tp.SetTime(14, 30)
	if tp.DisplayValue() == "Pick time" {
		t.Error("Expected formatted time, not placeholder")
	}
}

func TestFormatTime(t *testing.T) {
	tp := New("test")
	tp.Hour = 2
	tp.Minute = 30
	tp.Period = "PM"
	tp.HasValue = true

	if tp.FormatTime() != "2:30 PM" {
		t.Errorf("Expected '2:30 PM', got '%s'", tp.FormatTime())
	}
}

func TestFormatTime24(t *testing.T) {
	tp := New("test", With24Hour(true))
	tp.Hour = 14
	tp.Minute = 30
	tp.HasValue = true

	if tp.FormatTime() != "14:30" {
		t.Errorf("Expected '14:30', got '%s'", tp.FormatTime())
	}
}

func TestFormatTimeWithSeconds(t *testing.T) {
	tp := New("test", WithShowSeconds(true))
	tp.Hour = 2
	tp.Minute = 30
	tp.Second = 45
	tp.Period = "PM"
	tp.HasValue = true

	if tp.FormatTime() != "2:30:45 PM" {
		t.Errorf("Expected '2:30:45 PM', got '%s'", tp.FormatTime())
	}
}

func TestHourOptions(t *testing.T) {
	tp := New("test")
	hours := tp.HourOptions()

	if len(hours) != 12 {
		t.Errorf("Expected 12 hour options, got %d", len(hours))
	}
	if hours[0] != 1 || hours[11] != 12 {
		t.Error("Expected hours 1-12")
	}
}

func TestHourOptions24(t *testing.T) {
	tp := New("test", With24Hour(true))
	hours := tp.HourOptions()

	if len(hours) != 24 {
		t.Errorf("Expected 24 hour options, got %d", len(hours))
	}
	if hours[0] != 0 || hours[23] != 23 {
		t.Error("Expected hours 0-23")
	}
}

func TestMinuteOptions(t *testing.T) {
	tp := New("test", WithMinuteStep(15))
	minutes := tp.MinuteOptions()

	if len(minutes) != 4 {
		t.Errorf("Expected 4 minute options (0, 15, 30, 45), got %d", len(minutes))
	}
}

// Duration picker tests

func TestDurationToggle(t *testing.T) {
	dp := NewDuration("test")

	dp.Toggle()
	if !dp.Open {
		t.Error("Expected Open to be true")
	}

	dp.Toggle()
	if dp.Open {
		t.Error("Expected Open to be false")
	}
}

func TestSetDuration(t *testing.T) {
	dp := NewDuration("test")
	dp.SetDuration(2, 30)

	if dp.Hours != 2 || dp.Minutes != 30 {
		t.Errorf("Expected 2h 30m, got %dh %dm", dp.Hours, dp.Minutes)
	}
	if !dp.HasValue {
		t.Error("Expected HasValue to be true")
	}
}

func TestDurationClear(t *testing.T) {
	dp := NewDuration("test")
	dp.SetDuration(2, 30)
	dp.Clear()

	if dp.HasValue {
		t.Error("Expected HasValue to be false")
	}
	if dp.Hours != 0 || dp.Minutes != 0 {
		t.Error("Expected duration to be reset")
	}
}

func TestIncrementDecrementHours(t *testing.T) {
	dp := NewDuration("test", WithMaxHours(5))
	dp.Hours = 4

	dp.IncrementHours()
	if dp.Hours != 5 {
		t.Errorf("Expected Hours 5, got %d", dp.Hours)
	}

	dp.IncrementHours() // Should not exceed max
	if dp.Hours != 5 {
		t.Error("Expected Hours to stay at max")
	}

	dp.DecrementHours()
	if dp.Hours != 4 {
		t.Errorf("Expected Hours 4, got %d", dp.Hours)
	}
}

func TestIncrementDecrementMinutes(t *testing.T) {
	dp := NewDuration("test")
	dp.Minutes = 59

	dp.IncrementMinutes()
	if dp.Minutes != 0 {
		t.Errorf("Expected Minutes 0 (wrap), got %d", dp.Minutes)
	}

	dp.DecrementMinutes()
	if dp.Minutes != 59 {
		t.Errorf("Expected Minutes 59 (wrap), got %d", dp.Minutes)
	}
}

func TestTotalMinutes(t *testing.T) {
	dp := NewDuration("test")
	dp.Hours = 2
	dp.Minutes = 30

	if dp.TotalMinutes() != 150 {
		t.Errorf("Expected 150 minutes, got %d", dp.TotalMinutes())
	}
}

func TestTotalSeconds(t *testing.T) {
	dp := NewDuration("test", WithDurationShowSeconds(true))
	dp.Hours = 1
	dp.Minutes = 30
	dp.Seconds = 45

	expected := 1*3600 + 30*60 + 45
	if dp.TotalSeconds() != expected {
		t.Errorf("Expected %d seconds, got %d", expected, dp.TotalSeconds())
	}
}

func TestDurationDisplayValue(t *testing.T) {
	dp := NewDuration("test", WithDurationPlaceholder("Set duration"))

	if dp.DisplayValue() != "Set duration" {
		t.Errorf("Expected placeholder, got '%s'", dp.DisplayValue())
	}

	dp.SetDuration(2, 30)
	if dp.DisplayValue() != "2h 30m" {
		t.Errorf("Expected '2h 30m', got '%s'", dp.DisplayValue())
	}
}

func TestDurationFormatWithSeconds(t *testing.T) {
	dp := NewDuration("test", WithDurationShowSeconds(true))
	dp.Hours = 1
	dp.Minutes = 30
	dp.Seconds = 45
	dp.HasValue = true

	if dp.FormatDuration() != "1h 30m 45s" {
		t.Errorf("Expected '1h 30m 45s', got '%s'", dp.FormatDuration())
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
