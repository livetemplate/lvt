package tailwind

import "github.com/livetemplate/lvt/components/styles"

func accordionStyles() styles.AccordionStyles {
	return styles.AccordionStyles{
		Root:         "divide-y divide-gray-200 border border-gray-200 rounded-lg",
		Item:         "",
		ItemDisabled: "opacity-50",
		Header:       "flex items-center justify-between w-full px-4 py-4 text-left font-medium text-gray-900 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500",
		HeaderIcon:   "mr-3",
		ChevronIcon:  "w-5 h-5 text-gray-500 transition-transform duration-200",
		ChevronOpen:  "rotate-180",
		Content:      "px-4 pb-4 text-gray-600",
	}
}
