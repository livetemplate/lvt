package unstyled

import "github.com/livetemplate/lvt/components/styles"

func timepickerStyles() styles.TimepickerStyles {
	return styles.TimepickerStyles{
		Root:          "lvt-timepicker",
		TriggerBtn:    "lvt-timepicker__trigger",
		TriggerText:   "lvt-timepicker__trigger-text",
		Placeholder:   "lvt-timepicker__placeholder",
		TriggerIcon:   "lvt-timepicker__trigger-icon",
		Dropdown:      "lvt-timepicker__dropdown",
		SpinnerLayout: "lvt-timepicker__spinners",
		SpinnerCol:    "lvt-timepicker__spinner-col",
		SpinnerBtn:    "lvt-timepicker__spinner-btn",
		SpinnerIcon:   "lvt-timepicker__spinner-icon",
		SpinnerInput:  "lvt-timepicker__spinner-input",
		Separator:     "lvt-timepicker__separator",
		PeriodCol:     "lvt-timepicker__period-col",
		PeriodActive:  "lvt-timepicker__period--active",
		PeriodDefault: "lvt-timepicker__period--default",
		Footer:        "lvt-timepicker__footer",
		NowBtn:        "lvt-timepicker__now-btn",
		ClearBtn:      "lvt-timepicker__clear-btn",
	}
}
