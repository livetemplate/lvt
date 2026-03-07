package unstyled

import "github.com/livetemplate/lvt/components/styles"

func popoverStyles() styles.PopoverStyles {
	return styles.PopoverStyles{
		Root:      "lvt-popover",
		Panel:     "lvt-popover__panel",
		Header:    "lvt-popover__header",
		Title:     "lvt-popover__title",
		CloseBtn:  "lvt-popover__close-btn",
		CloseIcon: "lvt-popover__close-icon",
		Body:      "lvt-popover__body",
		// Position classes
		PosTop:         "lvt-popover--top",
		PosTopStart:    "lvt-popover--top-start",
		PosTopEnd:      "lvt-popover--top-end",
		PosBottom:      "lvt-popover--bottom",
		PosBottomStart: "lvt-popover--bottom-start",
		PosBottomEnd:   "lvt-popover--bottom-end",
		PosLeft:        "lvt-popover--left",
		PosLeftStart:   "lvt-popover--left-start",
		PosLeftEnd:     "lvt-popover--left-end",
		PosRight:       "lvt-popover--right",
		PosRightStart:  "lvt-popover--right-start",
		PosRightEnd:    "lvt-popover--right-end",
		// Arrow classes
		ArrowTop:    "lvt-popover__arrow--top",
		ArrowBottom: "lvt-popover__arrow--bottom",
		ArrowLeft:   "lvt-popover__arrow--left",
		ArrowRight:  "lvt-popover__arrow--right",
	}
}
