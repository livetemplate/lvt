package unstyled

import "github.com/livetemplate/lvt/components/styles"

func timelineStyles() styles.TimelineStyles {
	return styles.TimelineStyles{
		VerticalRoot:   "lvt-timeline--vertical",
		HorizontalRoot: "lvt-timeline--horizontal",
		Connector:      "lvt-timeline__connector",
	}
}

func timelineItemStyles() styles.TimelineItemStyles {
	return styles.TimelineItemStyles{
		VerticalItem:      "lvt-timeline-item--vertical",
		HorizontalItem:    "lvt-timeline-item--horizontal",
		IndicatorVertical: "lvt-timeline-item__indicator--vertical",
		IndicatorHoriz:    "lvt-timeline-item__indicator--horizontal",
		IndicatorIcon:     "lvt-timeline-item__indicator-icon",
		IndicatorDot:      "lvt-timeline-item__indicator-dot",
		CheckIcon:         "lvt-timeline-item__check-icon",
		ContentVertical:   "lvt-timeline-item__content--vertical",
		ContentHoriz:      "lvt-timeline-item__content--horizontal",
		Time:              "lvt-timeline-item__time",
		TimeVertical:      "lvt-timeline-item__time--vertical",
		TimeHoriz:         "lvt-timeline-item__time--horizontal",
		Title:             "lvt-timeline-item__title",
		TitleHoriz:        "lvt-timeline-item__title--horizontal",
		Description:       "lvt-timeline-item__description",
		HorizConnector:    "lvt-timeline-item__connector--horizontal",
		// Indicator color classes
		ColorGray:   "lvt-timeline-item--gray",
		ColorBlue:   "lvt-timeline-item--blue",
		ColorGreen:  "lvt-timeline-item--green",
		ColorYellow: "lvt-timeline-item--yellow",
		ColorRed:    "lvt-timeline-item--red",
		ColorPurple: "lvt-timeline-item--purple",
		// Status classes
		StatusDefault:  "lvt-timeline-item--status-default",
		StatusPending:  "lvt-timeline-item--status-pending",
		StatusActive:   "lvt-timeline-item--status-active",
		StatusComplete: "lvt-timeline-item--status-complete",
		StatusError:    "lvt-timeline-item--status-error",
		// Ring classes
		RingBlue:   "lvt-timeline-item--ring-blue",
		RingGreen:  "lvt-timeline-item--ring-green",
		RingYellow: "lvt-timeline-item--ring-yellow",
		RingRed:    "lvt-timeline-item--ring-red",
		RingPurple: "lvt-timeline-item--ring-purple",
		RingGray:   "lvt-timeline-item--ring-gray",
	}
}
