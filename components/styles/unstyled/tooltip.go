package unstyled

import "github.com/livetemplate/lvt/components/styles"

func tooltipStyles() styles.TooltipStyles {
	return styles.TooltipStyles{
		Root:  "lvt-tooltip",
		Panel: "lvt-tooltip__panel",
		// Position classes
		PosTop:         "lvt-tooltip--top",
		PosTopStart:    "lvt-tooltip--top-start",
		PosTopEnd:      "lvt-tooltip--top-end",
		PosBottom:      "lvt-tooltip--bottom",
		PosBottomStart: "lvt-tooltip--bottom-start",
		PosBottomEnd:   "lvt-tooltip--bottom-end",
		PosLeft:        "lvt-tooltip--left",
		PosLeftStart:   "lvt-tooltip--left-start",
		PosLeftEnd:     "lvt-tooltip--left-end",
		PosRight:       "lvt-tooltip--right",
		PosRightStart:  "lvt-tooltip--right-start",
		PosRightEnd:    "lvt-tooltip--right-end",
		// Arrow classes
		ArrowTop:    "lvt-tooltip__arrow--top",
		ArrowBottom: "lvt-tooltip__arrow--bottom",
		ArrowLeft:   "lvt-tooltip__arrow--left",
		ArrowRight:  "lvt-tooltip__arrow--right",
	}
}
