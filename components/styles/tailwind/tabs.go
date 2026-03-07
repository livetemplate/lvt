package tailwind

import "github.com/livetemplate/lvt/components/styles"

func tabsStyles() styles.TabsStyles {
	return styles.TabsStyles{
		Root:        "w-full",
		TabList:     "border-b border-gray-200",
		Nav:         "flex -mb-px space-x-8",
		Tab:         "py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap",
		TabActive:   "border-blue-500 text-blue-600",
		TabDefault:  "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300",
		TabDisabled: "opacity-50 cursor-not-allowed",
		TabIcon:     "mr-2",
		Badge:       "ml-2 px-2 py-0.5 text-xs rounded-full bg-gray-100 text-gray-600",
		// Vertical-specific
		VerticalRoot:    "flex",
		VerticalTabList: "border-r border-gray-200",
		VerticalNav:     "flex flex-col",
		VerticalTab:     "py-4 px-4 border-r-2 font-medium text-sm whitespace-nowrap",
		// Pills-specific
		PillsNav:     "flex gap-2",
		PillTab:      "px-4 py-2 rounded-md font-medium text-sm",
		PillActive:   "bg-blue-600 text-white",
		PillDefault:  "text-gray-500 hover:text-gray-700 hover:bg-gray-100",
		PillDisabled: "opacity-50 cursor-not-allowed",
	}
}
