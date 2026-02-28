package datepicker

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dp := New("test-dp")

	if dp.ID() != "test-dp" {
		t.Errorf("Expected ID 'test-dp', got '%s'", dp.ID())
	}
	if dp.Namespace() != "datepicker" {
		t.Errorf("Expected namespace 'datepicker', got '%s'", dp.Namespace())
	}
	if dp.Selected != nil {
		t.Error("Expected Selected to be nil")
	}
	if dp.Placeholder != "Select date..." {
		t.Errorf("Expected default placeholder 'Select date...', got '%s'", dp.Placeholder)
	}
	if dp.Format != "Jan 2, 2006" {
		t.Errorf("Expected default format 'Jan 2, 2006', got '%s'", dp.Format)
	}
	if dp.FirstDayOfWeek != 0 {
		t.Errorf("Expected default FirstDayOfWeek 0 (Sunday), got %d", dp.FirstDayOfWeek)
	}
	if dp.Open {
		t.Error("Expected Open to be false by default")
	}
}

func TestNewRange(t *testing.T) {
	rp := NewRange("test-range")

	if rp.ID() != "test-range" {
		t.Errorf("Expected ID 'test-range', got '%s'", rp.ID())
	}
	if rp.StartDate != nil {
		t.Error("Expected StartDate to be nil")
	}
	if rp.EndDate != nil {
		t.Error("Expected EndDate to be nil")
	}
	if rp.SelectingEnd {
		t.Error("Expected SelectingEnd to be false")
	}
}

func TestNewInline(t *testing.T) {
	dp := NewInline("test-inline")

	if dp.ID() != "test-inline" {
		t.Errorf("Expected ID 'test-inline', got '%s'", dp.ID())
	}
	if !dp.Open {
		t.Error("Expected inline datepicker to have Open=true")
	}
}

func TestWithPlaceholder(t *testing.T) {
	dp := New("test", WithPlaceholder("Pick a date"))
	if dp.Placeholder != "Pick a date" {
		t.Errorf("Expected placeholder 'Pick a date', got '%s'", dp.Placeholder)
	}
}

func TestWithFormat(t *testing.T) {
	dp := New("test", WithFormat("2006-01-02"))
	if dp.Format != "2006-01-02" {
		t.Errorf("Expected format '2006-01-02', got '%s'", dp.Format)
	}
}

func TestWithSelected(t *testing.T) {
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithSelected(date))

	if dp.Selected == nil {
		t.Fatal("Expected Selected to be set")
	}
	if !sameDay(*dp.Selected, date) {
		t.Errorf("Expected selected date %v, got %v", date, *dp.Selected)
	}
	// ViewDate should also be set to selected date
	if !sameDay(dp.ViewDate, date) {
		t.Errorf("Expected ViewDate to be set to selected date")
	}
}

func TestWithMinDate(t *testing.T) {
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithMinDate(minDate))

	if dp.MinDate == nil {
		t.Fatal("Expected MinDate to be set")
	}
	if !sameDay(*dp.MinDate, minDate) {
		t.Errorf("Expected MinDate %v, got %v", minDate, *dp.MinDate)
	}
}

func TestWithMaxDate(t *testing.T) {
	maxDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithMaxDate(maxDate))

	if dp.MaxDate == nil {
		t.Fatal("Expected MaxDate to be set")
	}
	if !sameDay(*dp.MaxDate, maxDate) {
		t.Errorf("Expected MaxDate %v, got %v", maxDate, *dp.MaxDate)
	}
}

func TestWithDisabledDates(t *testing.T) {
	dates := []time.Time{
		time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC),
	}
	dp := New("test", WithDisabledDates(dates...))

	if len(dp.DisabledDates) != 2 {
		t.Errorf("Expected 2 disabled dates, got %d", len(dp.DisabledDates))
	}
}

func TestWithDisabledWeekdays(t *testing.T) {
	dp := New("test", WithDisabledWeekdays(time.Saturday, time.Sunday))

	if len(dp.DisabledWeekdays) != 2 {
		t.Errorf("Expected 2 disabled weekdays, got %d", len(dp.DisabledWeekdays))
	}
}

func TestWithFirstDayOfWeek(t *testing.T) {
	dp := New("test", WithFirstDayOfWeek(1)) // Monday
	if dp.FirstDayOfWeek != 1 {
		t.Errorf("Expected FirstDayOfWeek 1, got %d", dp.FirstDayOfWeek)
	}

	// Test wrapping
	dp2 := New("test2", WithFirstDayOfWeek(8))
	if dp2.FirstDayOfWeek != 1 { // 8 % 7 = 1
		t.Errorf("Expected FirstDayOfWeek 1 (wrapped), got %d", dp2.FirstDayOfWeek)
	}
}

func TestWithStyled(t *testing.T) {
	dp := New("test", WithStyled(true))
	if !dp.IsStyled() {
		t.Error("Expected IsStyled to be true")
	}
}

func TestWithOpen(t *testing.T) {
	dp := New("test", WithOpen(true))
	if !dp.Open {
		t.Error("Expected Open to be true")
	}
}

func TestToggle(t *testing.T) {
	dp := New("test")

	if dp.Open {
		t.Error("Expected initial Open to be false")
	}

	dp.Toggle()
	if !dp.Open {
		t.Error("Expected Open to be true after toggle")
	}

	dp.Toggle()
	if dp.Open {
		t.Error("Expected Open to be false after second toggle")
	}
}

func TestClose(t *testing.T) {
	dp := New("test", WithOpen(true))
	dp.Close()
	if dp.Open {
		t.Error("Expected Open to be false after Close")
	}
}

func TestSelectDate(t *testing.T) {
	dp := New("test", WithOpen(true))
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	result := dp.SelectDate(date)

	if !result {
		t.Error("Expected SelectDate to return true")
	}
	if dp.Selected == nil {
		t.Fatal("Expected Selected to be set")
	}
	if !sameDay(*dp.Selected, date) {
		t.Errorf("Expected selected date %v, got %v", date, *dp.Selected)
	}
	if dp.Open {
		t.Error("Expected Open to be false after selection")
	}
}

func TestSelectDateDisabled(t *testing.T) {
	disabledDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithDisabledDates(disabledDate))

	result := dp.SelectDate(disabledDate)

	if result {
		t.Error("Expected SelectDate to return false for disabled date")
	}
	if dp.Selected != nil {
		t.Error("Expected Selected to remain nil")
	}
}

func TestClear(t *testing.T) {
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithSelected(date))

	dp.Clear()
	if dp.Selected != nil {
		t.Error("Expected Selected to be nil after Clear")
	}
}

func TestPreviousMonth(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	dp.PreviousMonth()

	if dp.ViewDate.Month() != time.May {
		t.Errorf("Expected month May, got %v", dp.ViewDate.Month())
	}
	if dp.ViewDate.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", dp.ViewDate.Year())
	}
}

func TestPreviousMonthYearRollover(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	dp.PreviousMonth()

	if dp.ViewDate.Month() != time.December {
		t.Errorf("Expected month December, got %v", dp.ViewDate.Month())
	}
	if dp.ViewDate.Year() != 2023 {
		t.Errorf("Expected year 2023, got %d", dp.ViewDate.Year())
	}
}

func TestNextMonth(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	dp.NextMonth()

	if dp.ViewDate.Month() != time.July {
		t.Errorf("Expected month July, got %v", dp.ViewDate.Month())
	}
}

func TestNextMonthYearRollover(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)

	dp.NextMonth()

	if dp.ViewDate.Month() != time.January {
		t.Errorf("Expected month January, got %v", dp.ViewDate.Month())
	}
	if dp.ViewDate.Year() != 2025 {
		t.Errorf("Expected year 2025, got %d", dp.ViewDate.Year())
	}
}

func TestPreviousYear(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	dp.PreviousYear()

	if dp.ViewDate.Year() != 2023 {
		t.Errorf("Expected year 2023, got %d", dp.ViewDate.Year())
	}
	if dp.ViewDate.Month() != time.June {
		t.Errorf("Expected month June, got %v", dp.ViewDate.Month())
	}
}

func TestNextYear(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	dp.NextYear()

	if dp.ViewDate.Year() != 2025 {
		t.Errorf("Expected year 2025, got %d", dp.ViewDate.Year())
	}
}

func TestGoToToday(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	dp.GoToToday()

	now := time.Now()
	if dp.ViewDate.Year() != now.Year() || dp.ViewDate.Month() != now.Month() {
		t.Errorf("Expected ViewDate to be today's month/year")
	}
}

func TestIsDateSelectableMinDate(t *testing.T) {
	minDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithMinDate(minDate))

	// Before min - not selectable
	if dp.IsDateSelectable(time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected date before MinDate to not be selectable")
	}

	// On min - selectable
	if !dp.IsDateSelectable(minDate) {
		t.Error("Expected MinDate itself to be selectable")
	}

	// After min - selectable
	if !dp.IsDateSelectable(time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected date after MinDate to be selectable")
	}
}

func TestIsDateSelectableMaxDate(t *testing.T) {
	maxDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithMaxDate(maxDate))

	// After max - not selectable
	if dp.IsDateSelectable(time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected date after MaxDate to not be selectable")
	}

	// On max - selectable
	if !dp.IsDateSelectable(maxDate) {
		t.Error("Expected MaxDate itself to be selectable")
	}

	// Before max - selectable
	if !dp.IsDateSelectable(time.Date(2024, 6, 19, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected date before MaxDate to be selectable")
	}
}

func TestIsDateSelectableDisabledWeekdays(t *testing.T) {
	dp := New("test", WithDisabledWeekdays(time.Saturday, time.Sunday))

	// Saturday - not selectable (June 15, 2024 is Saturday)
	if dp.IsDateSelectable(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected Saturday to not be selectable")
	}

	// Sunday - not selectable (June 16, 2024 is Sunday)
	if dp.IsDateSelectable(time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected Sunday to not be selectable")
	}

	// Monday - selectable (June 17, 2024 is Monday)
	if !dp.IsDateSelectable(time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected Monday to be selectable")
	}
}

func TestIsDateSelectableDisabledDates(t *testing.T) {
	disabledDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithDisabledDates(disabledDate))

	if dp.IsDateSelectable(disabledDate) {
		t.Error("Expected disabled date to not be selectable")
	}

	if !dp.IsDateSelectable(time.Date(2024, 6, 11, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected non-disabled date to be selectable")
	}
}

func TestIsSelected(t *testing.T) {
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithSelected(date))

	if !dp.IsSelected(date) {
		t.Error("Expected date to be selected")
	}

	otherDate := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)
	if dp.IsSelected(otherDate) {
		t.Error("Expected other date to not be selected")
	}
}

func TestIsSelectedWithNil(t *testing.T) {
	dp := New("test")

	if dp.IsSelected(time.Now()) {
		t.Error("Expected no date to be selected when Selected is nil")
	}
}

func TestIsToday(t *testing.T) {
	dp := New("test")
	today := time.Now()

	if !dp.IsToday(today) {
		t.Error("Expected today to be recognized as today")
	}

	yesterday := today.AddDate(0, 0, -1)
	if dp.IsToday(yesterday) {
		t.Error("Expected yesterday to not be today")
	}
}

func TestDisplayValue(t *testing.T) {
	dp := New("test", WithPlaceholder("Pick date"))

	// No date selected - show placeholder
	if dp.DisplayValue() != "Pick date" {
		t.Errorf("Expected placeholder 'Pick date', got '%s'", dp.DisplayValue())
	}

	// Date selected - show formatted date
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	dp.SelectDate(date)
	if dp.DisplayValue() != "Jun 15, 2024" {
		t.Errorf("Expected 'Jun 15, 2024', got '%s'", dp.DisplayValue())
	}
}

func TestDisplayValueCustomFormat(t *testing.T) {
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	dp := New("test", WithFormat("2006-01-02"), WithSelected(date))

	if dp.DisplayValue() != "2024-06-15" {
		t.Errorf("Expected '2024-06-15', got '%s'", dp.DisplayValue())
	}
}

func TestViewMonth(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	if dp.ViewMonth() != "June" {
		t.Errorf("Expected 'June', got '%s'", dp.ViewMonth())
	}
}

func TestViewYear(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	if dp.ViewYear() != 2024 {
		t.Errorf("Expected 2024, got %d", dp.ViewYear())
	}
}

func TestWeekdayNames(t *testing.T) {
	dp := New("test") // Sunday first
	names := dp.WeekdayNames()

	expected := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	if len(names) != 7 {
		t.Fatalf("Expected 7 weekday names, got %d", len(names))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected weekday %d to be '%s', got '%s'", i, expected[i], name)
		}
	}
}

func TestWeekdayNamesMondayFirst(t *testing.T) {
	dp := New("test", WithFirstDayOfWeek(1)) // Monday first
	names := dp.WeekdayNames()

	expected := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected weekday %d to be '%s', got '%s'", i, expected[i], name)
		}
	}
}

func TestCalendarWeeks(t *testing.T) {
	dp := New("test")
	dp.ViewDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC) // June 2024

	weeks := dp.CalendarWeeks()

	// June 2024 should have 5 or 6 weeks displayed
	if len(weeks) < 4 || len(weeks) > 6 {
		t.Errorf("Expected 4-6 weeks, got %d", len(weeks))
	}

	// Each week should have 7 days
	for i, week := range weeks {
		if len(week) != 7 {
			t.Errorf("Week %d should have 7 days, got %d", i, len(week))
		}
	}

	// Find June 1st - it should be a Saturday (index 6 for Sunday-first)
	var june1Found bool
	for _, week := range weeks {
		for _, day := range week {
			if day.Day == 1 && day.InMonth {
				june1Found = true
				// June 1, 2024 is a Saturday
				if day.Date.Weekday() != time.Saturday {
					t.Errorf("Expected June 1 to be Saturday, got %v", day.Date.Weekday())
				}
			}
		}
	}
	if !june1Found {
		t.Error("June 1st not found in calendar")
	}
}

func TestCalendarDayDateString(t *testing.T) {
	day := CalendarDay{
		Date: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
	}

	if day.DateString() != "2024-06-15" {
		t.Errorf("Expected '2024-06-15', got '%s'", day.DateString())
	}
}

// RangePicker tests

func TestSelectRangeDateStart(t *testing.T) {
	rp := NewRange("test")
	date := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	result := rp.SelectRangeDate(date)

	if !result {
		t.Error("Expected SelectRangeDate to return true")
	}
	if rp.StartDate == nil {
		t.Fatal("Expected StartDate to be set")
	}
	if !sameDay(*rp.StartDate, date) {
		t.Errorf("Expected StartDate %v, got %v", date, *rp.StartDate)
	}
	if !rp.SelectingEnd {
		t.Error("Expected SelectingEnd to be true after selecting start")
	}
	if rp.EndDate != nil {
		t.Error("Expected EndDate to be nil")
	}
}

func TestSelectRangeDateEnd(t *testing.T) {
	rp := NewRange("test", WithOpen(true))
	startDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)

	rp.SelectRangeDate(startDate)
	rp.SelectRangeDate(endDate)

	if rp.StartDate == nil || rp.EndDate == nil {
		t.Fatal("Expected both StartDate and EndDate to be set")
	}
	if !sameDay(*rp.StartDate, startDate) {
		t.Errorf("Expected StartDate %v, got %v", startDate, *rp.StartDate)
	}
	if !sameDay(*rp.EndDate, endDate) {
		t.Errorf("Expected EndDate %v, got %v", endDate, *rp.EndDate)
	}
	if rp.SelectingEnd {
		t.Error("Expected SelectingEnd to be false after complete selection")
	}
	if rp.Open {
		t.Error("Expected Open to be false after completing range selection")
	}
}

func TestSelectRangeDateEndBeforeStart(t *testing.T) {
	rp := NewRange("test")
	startDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC) // Before start

	rp.SelectRangeDate(startDate)
	rp.SelectRangeDate(endDate)

	// Should swap dates
	if !sameDay(*rp.StartDate, endDate) {
		t.Errorf("Expected StartDate to be swapped to %v, got %v", endDate, *rp.StartDate)
	}
	if !sameDay(*rp.EndDate, startDate) {
		t.Errorf("Expected EndDate to be swapped to %v, got %v", startDate, *rp.EndDate)
	}
}

func TestClearRange(t *testing.T) {
	rp := NewRange("test")
	startDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)

	rp.SelectRangeDate(startDate)
	rp.SelectRangeDate(endDate)
	rp.ClearRange()

	if rp.StartDate != nil {
		t.Error("Expected StartDate to be nil after ClearRange")
	}
	if rp.EndDate != nil {
		t.Error("Expected EndDate to be nil after ClearRange")
	}
	if rp.SelectingEnd {
		t.Error("Expected SelectingEnd to be false after ClearRange")
	}
}

func TestIsInRange(t *testing.T) {
	rp := NewRange("test")
	startDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)

	rp.SelectRangeDate(startDate)
	rp.SelectRangeDate(endDate)

	// Start date - in range
	if !rp.IsInRange(startDate) {
		t.Error("Expected start date to be in range")
	}

	// End date - in range
	if !rp.IsInRange(endDate) {
		t.Error("Expected end date to be in range")
	}

	// Middle date - in range
	middleDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !rp.IsInRange(middleDate) {
		t.Error("Expected middle date to be in range")
	}

	// Before start - not in range
	beforeDate := time.Date(2024, 6, 5, 0, 0, 0, 0, time.UTC)
	if rp.IsInRange(beforeDate) {
		t.Error("Expected date before start to not be in range")
	}

	// After end - not in range
	afterDate := time.Date(2024, 6, 25, 0, 0, 0, 0, time.UTC)
	if rp.IsInRange(afterDate) {
		t.Error("Expected date after end to not be in range")
	}
}

func TestIsInRangeIncomplete(t *testing.T) {
	rp := NewRange("test")

	// No range selected - nothing in range
	if rp.IsInRange(time.Now()) {
		t.Error("Expected nothing to be in range when no dates selected")
	}

	// Only start date - nothing in range
	rp.SelectRangeDate(time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC))
	if rp.IsInRange(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected nothing to be in range when only start date selected")
	}
}

func TestDisplayRangeValue(t *testing.T) {
	rp := NewRange("test", WithPlaceholder("Select range"))

	// No dates - placeholder
	if rp.DisplayRangeValue() != "Select range" {
		t.Errorf("Expected placeholder, got '%s'", rp.DisplayRangeValue())
	}

	// Start date only
	startDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	rp.SelectRangeDate(startDate)
	if rp.DisplayRangeValue() != "Jun 10, 2024 - ..." {
		t.Errorf("Expected 'Jun 10, 2024 - ...', got '%s'", rp.DisplayRangeValue())
	}

	// Both dates
	endDate := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
	rp.SelectRangeDate(endDate)
	if rp.DisplayRangeValue() != "Jun 10, 2024 - Jun 20, 2024" {
		t.Errorf("Expected 'Jun 10, 2024 - Jun 20, 2024', got '%s'", rp.DisplayRangeValue())
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
