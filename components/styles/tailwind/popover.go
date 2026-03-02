package tailwind

import "github.com/livetemplate/lvt/components/styles"

func popoverStyles() styles.PopoverStyles {
	return styles.PopoverStyles{
		Root:      "relative inline-block",
		Panel:     "absolute z-50 bg-white rounded-lg shadow-lg border border-gray-200",
		Header:    "flex items-center justify-between px-4 py-2 border-b border-gray-100",
		Title:     "text-sm font-semibold text-gray-900",
		CloseBtn:  "p-1 text-gray-400 hover:text-gray-600 rounded",
		CloseIcon: "w-4 h-4",
		Body:      "px-4 py-3 text-sm text-gray-600",
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
		ArrowTop:    "absolute top-full left-1/2 -translate-x-1/2 border-8 border-transparent border-t-white",
		ArrowBottom: "absolute bottom-full left-1/2 -translate-x-1/2 border-8 border-transparent border-b-white",
		ArrowLeft:   "absolute left-full top-1/2 -translate-y-1/2 border-8 border-transparent border-l-white",
		ArrowRight:  "absolute right-full top-1/2 -translate-y-1/2 border-8 border-transparent border-r-white",
	}
}
