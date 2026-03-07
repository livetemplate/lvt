package tailwind

import "github.com/livetemplate/lvt/components/styles"

func tooltipStyles() styles.TooltipStyles {
	return styles.TooltipStyles{
		Root:  "relative inline-block",
		Panel: "absolute z-50 px-2 py-1 text-sm text-white bg-gray-900 rounded shadow-lg whitespace-nowrap",
		// Position classes
		PosTop:         "bottom-full left-1/2 -translate-x-1/2 mb-2",
		PosTopStart:    "bottom-full left-0 mb-2",
		PosTopEnd:      "bottom-full right-0 mb-2",
		PosBottom:      "top-full left-1/2 -translate-x-1/2 mt-2",
		PosBottomStart: "top-full left-0 mt-2",
		PosBottomEnd:   "top-full right-0 mt-2",
		PosLeft:        "right-full top-1/2 -translate-y-1/2 mr-2",
		PosLeftStart:   "right-full top-0 mr-2",
		PosLeftEnd:     "right-full bottom-0 mr-2",
		PosRight:       "left-full top-1/2 -translate-y-1/2 ml-2",
		PosRightStart:  "left-full top-0 ml-2",
		PosRightEnd:    "left-full bottom-0 ml-2",
		// Arrow classes
		ArrowTop:    "absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-900",
		ArrowBottom: "absolute bottom-full left-1/2 -translate-x-1/2 border-4 border-transparent border-b-gray-900",
		ArrowLeft:   "absolute left-full top-1/2 -translate-y-1/2 border-4 border-transparent border-l-gray-900",
		ArrowRight:  "absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-gray-900",
	}
}
