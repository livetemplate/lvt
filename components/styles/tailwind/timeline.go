package tailwind

import "github.com/livetemplate/lvt/components/styles"

func timelineStyles() styles.TimelineStyles {
	return styles.TimelineStyles{
		VerticalRoot:   "relative",
		HorizontalRoot: "flex flex-row space-x-4 overflow-x-auto",
		Connector:      "absolute left-4 top-0 bottom-0 w-0.5 bg-gray-200",
	}
}

func timelineItemStyles() styles.TimelineItemStyles {
	return styles.TimelineItemStyles{
		VerticalItem:      "relative pl-10 pb-8 last:pb-0",
		HorizontalItem:    "flex-shrink-0 flex flex-col items-center",
		IndicatorVertical: "absolute left-0 w-8 h-8 rounded-full flex items-center justify-center z-10",
		IndicatorHoriz:    "w-10 h-10 rounded-full flex items-center justify-center mb-3",
		IndicatorIcon:     "text-white text-sm",
		IndicatorDot:      "w-2.5 h-2.5 rounded-full bg-white",
		CheckIcon:         "w-4 h-4 text-white",
		ContentVertical:   "flex-1",
		ContentHoriz:      "text-center max-w-32",
		Time:              "text-sm text-gray-500",
		TimeVertical:      "mb-1 block",
		TimeHoriz:         "text-xs block",
		Title:             "text-base font-semibold text-gray-900",
		TitleHoriz:        "text-sm font-medium text-gray-900",
		Description:       "mt-1 text-sm text-gray-600",
		HorizConnector:    "absolute top-5 left-full w-4 h-0.5 bg-gray-200",
		// Indicator color classes
		ColorGray:   "bg-gray-400",
		ColorBlue:   "bg-blue-500",
		ColorGreen:  "bg-green-500",
		ColorYellow: "bg-yellow-500",
		ColorRed:    "bg-red-500",
		ColorPurple: "bg-purple-500",
		// Status classes
		StatusDefault:  "bg-gray-400 text-white",
		StatusPending:  "bg-gray-200 text-gray-500",
		StatusActive:   "bg-blue-500 text-white ring-4 ring-blue-100",
		StatusComplete: "bg-green-500 text-white",
		StatusError:    "bg-red-500 text-white",
		// Ring classes
		RingBlue:   "ring-4 ring-blue-100",
		RingGreen:  "ring-4 ring-green-100",
		RingYellow: "ring-4 ring-yellow-100",
		RingRed:    "ring-4 ring-red-100",
		RingPurple: "ring-4 ring-purple-100",
		RingGray:   "ring-4 ring-gray-100",
	}
}
