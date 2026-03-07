package tailwind

import "github.com/livetemplate/lvt/components/styles"

func tagsInputStyles() styles.TagsInputStyles {
	return styles.TagsInputStyles{
		Root:          "relative",
		Wrapper:       "flex flex-wrap gap-2 p-2 border border-gray-300 rounded-md bg-white focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-blue-500",
		Tag:           "inline-flex items-center gap-1 px-2 py-1 text-sm bg-blue-100 text-blue-800 rounded-md",
		TagRemoveBtn:  "inline-flex items-center justify-center w-4 h-4 text-blue-600 hover:text-blue-800 hover:bg-blue-200 rounded-full",
		TagRemoveIcon: "w-3 h-3",
		Input:         "flex-1 min-w-[120px] outline-none bg-transparent text-sm",
		Dropdown:      "absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-48 overflow-auto",
		Suggestion:    "px-3 py-2 text-sm cursor-pointer hover:bg-blue-50",
		Counter:       "mt-1 text-xs text-gray-500",
	}
}
