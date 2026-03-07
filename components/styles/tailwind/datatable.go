package tailwind

import "github.com/livetemplate/lvt/components/styles"

func datatableStyles() styles.DatatableStyles {
	return styles.DatatableStyles{
		Root:            "overflow-hidden",
		FilterInput:     "w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500",
		FilterWrapper:   "mb-4",
		LoadingWrapper:  "flex items-center justify-center py-8",
		LoadingIcon:     "w-8 h-8 text-blue-600 animate-spin",
		TableWrapper:    "overflow-x-auto border border-gray-200 rounded-lg",
		Table:           "min-w-full divide-y divide-gray-200",
		Thead:           "bg-gray-50",
		Th:              "px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider",
		ThSortable:      "cursor-pointer hover:bg-gray-100",
		ThAlignCenter:   "text-center",
		ThAlignRight:    "text-right",
		ThContent:       "flex items-center gap-1",
		SortIcon:        "w-4 h-4",
		SortIconIdle:    "w-4 h-4 text-gray-300",
		Tbody:           "bg-white divide-y divide-gray-200",
		Td:              "px-4 py-3 text-sm text-gray-900",
		TdCompact:       "py-2",
		TdAlignCenter:   "text-center",
		TdAlignRight:    "text-right",
		RowStriped:      "bg-gray-50",
		RowHover:        "hover:bg-gray-100",
		RowSelected:     "bg-blue-50",
		Checkbox:        "h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500",
		EmptyCell:       "px-4 py-8 text-center text-gray-500",
		Pagination:      "flex items-center justify-between px-4 py-3 bg-white border-t border-gray-200",
		PageInfo:        "text-sm text-gray-700",
		PageInfoBold:    "font-medium",
		PageBtn:         "px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50",
		PageBtnDisabled: "opacity-50 cursor-not-allowed",
	}
}
