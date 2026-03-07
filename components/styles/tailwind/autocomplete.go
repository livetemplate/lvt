package tailwind

import "github.com/livetemplate/lvt/components/styles"

func autocompleteStyles() styles.AutocompleteStyles {
	return styles.AutocompleteStyles{
		Root:           "relative",
		InputWrapper:   "relative",
		Input:          "w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500",
		InputLoading:   "pr-10",
		LoadingWrapper: "absolute inset-y-0 right-0 flex items-center pr-3",
		LoadingIcon:    "w-5 h-5 text-gray-400 animate-spin",
		ClearBtn:          "absolute inset-y-0 right-0 flex items-center pr-3",
		ClearIcon:         "w-5 h-5 text-gray-400 hover:text-gray-600",
		SelectedBadge:     "inline-flex items-center gap-1 px-2 py-1 rounded-md bg-blue-100 text-blue-800 text-sm",
		SelectedBadgeBtn:  "text-blue-600 hover:text-blue-800",
		Dropdown:          "absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-auto",
		Option:         "px-4 py-2 cursor-pointer text-gray-900 hover:bg-gray-100",
		OptionDisabled: "text-gray-400 cursor-not-allowed",
		OptionActive:   "bg-blue-600 text-white",
		OptionLabel:    "font-medium",
		OptionDesc:     "text-sm text-gray-500",
		OptionDescAlt:  "text-sm text-blue-200",
		OptionIcon:     "mr-2",
		OptionLayout:   "flex items-center",
		Empty:          "px-4 py-2 text-gray-500 text-center",
	}
}
