package tailwind

import "github.com/livetemplate/lvt/components/styles"

func progressStyles() styles.ProgressStyles {
	return styles.ProgressStyles{
		Root:             "w-full",
		Header:           "flex justify-between mb-1",
		LabelText:        "text-sm font-medium text-gray-700",
		ValueText:        "text-sm font-medium text-gray-700",
		Track:            "w-full bg-gray-200 rounded-full overflow-hidden",
		Bar:              "rounded-full transition-all duration-300",
		BarIndeterminate: "animate-indeterminate",
		Striped:          "bg-stripes",
		Animated:         "animate-stripes",
		// Size classes
		SizeXs: "h-1",
		SizeSm: "h-2",
		SizeMd: "h-4",
		SizeLg: "h-6",
		// Color classes
		ColorPrimary: "bg-blue-500",
		ColorSuccess: "bg-green-500",
		ColorWarning: "bg-yellow-500",
		ColorDanger:  "bg-red-500",
		ColorInfo:    "bg-cyan-500",
	}
}

func circularProgressStyles() styles.CircularProgressStyles {
	return styles.CircularProgressStyles{
		Root:             "relative inline-flex items-center justify-center",
		SvgIndeterminate: "animate-spin",
		TrackCircle:      "text-gray-200",
		BarCircle:        "transition-all duration-300",
		Label:            "absolute text-sm font-medium text-gray-700",
		// Color classes
		ColorPrimary: "text-blue-500",
		ColorSuccess: "text-green-500",
		ColorWarning: "text-yellow-500",
		ColorDanger:  "text-red-500",
		ColorInfo:    "text-cyan-500",
	}
}

func spinnerStyles() styles.SpinnerStyles {
	return styles.SpinnerStyles{
		Root:       "inline-flex items-center",
		Svg:        "animate-spin",
		TrackClass: "opacity-25",
		BarClass:   "opacity-75",
		SrOnly:     "sr-only",
		// Size classes
		SizeSm: "w-4 h-4",
		SizeMd: "w-6 h-6",
		SizeLg: "w-8 h-8",
		SizeXl: "w-12 h-12",
		// Color classes
		ColorPrimary: "text-blue-500",
		ColorSuccess: "text-green-500",
		ColorWarning: "text-yellow-500",
		ColorDanger:  "text-red-500",
		ColorInfo:    "text-cyan-500",
	}
}
