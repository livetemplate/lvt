package tailwind

import "github.com/livetemplate/lvt/components/styles"

func ratingStyles() styles.RatingStyles {
	return styles.RatingStyles{
		Root:          "inline-flex items-center gap-2",
		Label:         "text-sm text-gray-600",
		StarsWrapper:  "inline-flex",
		StarBtn:       "cursor-pointer transition-colors",
		StarReadonly:   "cursor-default",
		HalfStarOuter: "relative",
		HalfStarInner: "absolute overflow-hidden",
		ValueText:     "text-sm font-medium text-gray-700",
		CountText:     "text-sm text-gray-500",
		// Size classes
		SizeSm: "text-lg",
		SizeMd: "text-2xl",
		SizeLg: "text-3xl",
		SizeXl: "text-4xl",
		// Color classes
		ColorYellow: "text-yellow-400",
		ColorRed:    "text-red-500",
		ColorBlue:   "text-blue-500",
		ColorGreen:  "text-green-500",
		// Empty color classes
		EmptyDefault: "text-gray-300",
		EmptyYellow:  "text-yellow-200",
		EmptyRed:     "text-red-200",
		EmptyBlue:    "text-blue-200",
		EmptyGreen:   "text-green-200",
	}
}
