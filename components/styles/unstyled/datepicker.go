package unstyled

import "github.com/livetemplate/lvt/components/styles"

func datepickerStyles() styles.DatepickerStyles {
	return styles.DatepickerStyles{
		Root:        "lvt-datepicker",
		TriggerBtn:  "lvt-datepicker__trigger",
		TriggerText: "lvt-datepicker__trigger-text",
		Placeholder: "lvt-datepicker__placeholder",
		TriggerIcon: "lvt-datepicker__trigger-icon",
		Dropdown:    "lvt-datepicker__dropdown",
		Calendar:    "lvt-datepicker__calendar",
		NavWrapper:  "lvt-datepicker__nav",
		NavBtn:      "lvt-datepicker__nav-btn",
		NavIcon:     "lvt-datepicker__nav-icon",
		MonthLabel:  "lvt-datepicker__month-label",
		WeekdayGrid: "lvt-datepicker__weekday-grid",
		Weekday:     "lvt-datepicker__weekday",
		WeekGrid:    "lvt-datepicker__week-grid",
		DayBtn:      "lvt-datepicker__day",
		DayOutMonth: "lvt-datepicker__day--out-month",
		DayDisabled: "lvt-datepicker__day--disabled",
		DaySelected: "lvt-datepicker__day--selected",
		DayToday:    "lvt-datepicker__day--today",
		DayDefault:  "lvt-datepicker__day--default",
		Footer:      "lvt-datepicker__footer",
		TodayBtn:    "lvt-datepicker__today-btn",
		ClearBtn:    "lvt-datepicker__clear-btn",
	}
}
