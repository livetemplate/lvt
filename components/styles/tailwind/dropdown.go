package tailwind

import "github.com/livetemplate/lvt/components/styles"

func dropdownStyles() styles.DropdownStyles {
	return styles.DropdownStyles{
		Root:            "relative inline-block w-full",
		TriggerBtn:      "w-full px-4 py-2 text-left bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed",
		SelectedText:    "block truncate",
		TriggerIconWrap: "absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none",
		TriggerIcon:     "w-5 h-5 text-gray-400",
		ClearBtnWrap:    "absolute inset-y-0 right-8 flex items-center pr-2 cursor-pointer",
		Dropdown:        "absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto",
		Option:          "px-4 py-2 cursor-pointer hover:bg-blue-50",
		OptionDisabled:  "opacity-50 cursor-not-allowed",
		OptionSelected:  "bg-blue-100",
	}
}
