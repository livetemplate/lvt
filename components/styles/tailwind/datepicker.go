package tailwind

import "github.com/livetemplate/lvt/components/styles"

func datepickerStyles() styles.DatepickerStyles {
	return styles.DatepickerStyles{
		Root:         "relative inline-block",
		TriggerBtn:   "w-full px-4 py-2 text-left bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500",
		TriggerText:  "flex items-center justify-between",
		Placeholder:  "text-gray-400",
		TriggerIcon:  "w-5 h-5 text-gray-400",
		Dropdown:     "absolute z-10 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg p-4",
		Calendar:     "w-64",
		NavWrapper:   "flex items-center justify-between mb-4",
		NavBtn:       "p-1 hover:bg-gray-100 rounded",
		NavIcon:      "w-5 h-5",
		MonthLabel:   "font-semibold text-gray-900",
		WeekdayGrid:  "grid grid-cols-7 gap-1 mb-2",
		Weekday:      "text-center text-xs font-medium text-gray-500 py-1",
		WeekGrid:     "grid grid-cols-7 gap-1",
		DayBtn:       "w-8 h-8 text-sm rounded-full flex items-center justify-center",
		DayOutMonth:  "text-gray-300",
		DayDisabled:  "text-gray-300 cursor-not-allowed",
		DaySelected:  "bg-blue-600 text-white",
		DayToday:     "border border-blue-600 text-blue-600",
		DayDefault:   "text-gray-700 hover:bg-gray-100",
		Footer:       "flex justify-between mt-4 pt-4 border-t border-gray-200",
		TodayBtn:     "text-sm text-blue-600 hover:text-blue-700",
		ClearBtn:     "text-sm text-gray-500 hover:text-gray-700",
	}
}
