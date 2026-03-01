package tailwind

import "github.com/livetemplate/lvt/components/styles"

func timepickerStyles() styles.TimepickerStyles {
	return styles.TimepickerStyles{
		Root:          "relative inline-block",
		TriggerBtn:    "w-full px-4 py-2 text-left bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500",
		TriggerText:   "flex items-center justify-between",
		Placeholder:   "text-gray-400",
		TriggerIcon:   "w-5 h-5 text-gray-400",
		Dropdown:      "absolute z-10 mt-1 p-4 bg-white border border-gray-200 rounded-lg shadow-lg",
		SpinnerLayout: "flex items-center justify-center gap-2",
		SpinnerCol:    "flex flex-col items-center",
		SpinnerBtn:    "p-1 hover:bg-gray-100 rounded",
		SpinnerIcon:   "w-5 h-5",
		SpinnerInput:  "w-12 text-center text-lg font-semibold border border-gray-300 rounded py-1",
		Separator:     "text-lg font-semibold",
		PeriodCol:     "flex flex-col items-center ml-2",
		PeriodActive:  "px-2 py-1 text-sm font-medium rounded bg-blue-600 text-white",
		PeriodDefault: "px-2 py-1 text-sm font-medium rounded bg-gray-100 text-gray-700 hover:bg-gray-200",
		Footer:        "flex justify-between mt-4 pt-4 border-t border-gray-200",
		NowBtn:        "text-sm text-blue-600 hover:text-blue-700",
		ClearBtn:      "text-sm text-gray-500 hover:text-gray-700",
	}
}
