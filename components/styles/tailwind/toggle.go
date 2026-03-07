package tailwind

import "github.com/livetemplate/lvt/components/styles"

func toggleStyles() styles.ToggleStyles {
	return styles.ToggleStyles{
		UseCustomTrack: true,
		Label:          "inline-flex items-center gap-3 cursor-pointer",
		LabelDisabled:  "inline-flex items-center gap-3 cursor-not-allowed opacity-60",
		Input:          "sr-only peer",
		LabelText:      "text-sm font-medium text-gray-700",
		Description:    "text-sm text-gray-500",
		Track:          "relative rounded-full transition-colors duration-200 ease-in-out",
		Knob:           "absolute top-0.5 bg-white rounded-full shadow transform transition-transform duration-200 ease-in-out",
		// Track size classes
		TrackSm: "w-8 h-4",
		TrackMd: "w-11 h-6",
		TrackLg: "w-14 h-8",
		// Knob size classes
		KnobSm: "w-3 h-3",
		KnobMd: "w-5 h-5",
		KnobLg: "w-6 h-6",
		// Track color classes
		TrackChecked:           "bg-blue-600",
		TrackUnchecked:         "bg-gray-200",
		TrackCheckedDisabled:   "bg-blue-300",
		TrackUncheckedDisabled: "bg-gray-200",
		// Knob translate classes
		KnobUnchecked: "translate-x-0.5",
		KnobSmChecked: "translate-x-4",
		KnobMdChecked: "translate-x-5",
		KnobLgChecked: "translate-x-7",
	}
}

func checkboxStyles() styles.CheckboxStyles {
	return styles.CheckboxStyles{
		UseCustomCheckbox: true,
		Label:             "inline-flex items-start gap-3 cursor-pointer",
		LabelDisabled:     "inline-flex items-start gap-3 cursor-not-allowed opacity-60",
		Input:             "sr-only peer",
		CheckboxWrap:      "relative flex items-center justify-center",
		CheckboxBox:       "w-5 h-5 rounded border-2 transition-colors flex items-center justify-center",
		CheckIcon:         "w-3 h-3 text-white",
		LabelTextWrap:     "flex-1",
		LabelText:         "text-sm font-medium text-gray-700",
		Description:       "text-sm text-gray-500",
		// State classes
		StateChecked:           "bg-blue-600 border-blue-600",
		StateUnchecked:         "bg-white border-gray-300",
		StateCheckedDisabled:   "bg-blue-300 border-blue-300",
		StateUncheckedDisabled: "bg-gray-100 border-gray-200",
	}
}
