package unstyled

import "github.com/livetemplate/lvt/components/styles"

func ratingStyles() styles.RatingStyles {
	return styles.RatingStyles{
		Root:          "lvt-rating",
		Label:         "lvt-rating__label",
		StarsWrapper:  "lvt-rating__stars",
		StarBtn:       "lvt-rating__star",
		StarReadonly:  "lvt-rating__star--readonly",
		HalfStarOuter: "lvt-rating__half-star",
		HalfStarInner: "lvt-rating__half-star-inner",
		ValueText:     "lvt-rating__value",
		CountText:     "lvt-rating__count",
		// Size classes
		SizeSm: "lvt-rating--sm",
		SizeMd: "lvt-rating--md",
		SizeLg: "lvt-rating--lg",
		SizeXl: "lvt-rating--xl",
		// Color classes
		ColorYellow: "lvt-rating--yellow",
		ColorRed:    "lvt-rating--red",
		ColorBlue:   "lvt-rating--blue",
		ColorGreen:  "lvt-rating--green",
		// Empty color classes
		EmptyDefault: "lvt-rating__star--empty",
		EmptyYellow:  "lvt-rating__star--empty-yellow",
		EmptyRed:     "lvt-rating__star--empty-red",
		EmptyBlue:    "lvt-rating__star--empty-blue",
		EmptyGreen:   "lvt-rating__star--empty-green",
	}
}
