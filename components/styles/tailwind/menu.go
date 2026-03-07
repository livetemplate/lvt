package tailwind

import "github.com/livetemplate/lvt/components/styles"

func menuStyles() styles.MenuStyles {
	return styles.MenuStyles{
		Root:            "relative inline-block text-left",
		TriggerBtn:      "inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
		TriggerIcon:     "w-4 h-4",
		TriggerIconOpen: "rotate-180",
		Panel:           "absolute z-10 mt-2 w-56 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none",
		PanelInner:      "py-1",
		Divider:         "my-1 border-gray-200",
		SectionHeader:   "px-4 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider",
		Item:            "flex items-center justify-between px-4 py-2 text-sm",
		ItemDisabled:    "text-gray-400 cursor-not-allowed",
		ItemActive:      "bg-gray-100 text-gray-900",
		ItemDefault:     "text-gray-700 hover:bg-gray-100",
		ItemIcon:        "w-5 h-5",
		Badge:           "px-2 py-0.5 text-xs rounded-full",
		BadgeRed:        "bg-red-100 text-red-700",
		BadgeBlue:       "bg-blue-100 text-blue-700",
		BadgeGreen:      "bg-green-100 text-green-700",
		BadgeDefault:    "bg-gray-100 text-gray-700",
		Shortcut:        "text-xs text-gray-400",
		// Position variants
		PositionBottomLeft:  "left-0",
		PositionBottomRight: "right-0",
		PositionTopLeft:     "bottom-full mb-2 left-0",
		PositionTopRight:    "bottom-full mb-2 right-0",
	}
}
