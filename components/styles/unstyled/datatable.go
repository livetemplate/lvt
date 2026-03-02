package unstyled

import "github.com/livetemplate/lvt/components/styles"

func datatableStyles() styles.DatatableStyles {
	return styles.DatatableStyles{
		Root:            "lvt-datatable",
		FilterInput:     "lvt-datatable__filter-input",
		FilterWrapper:   "lvt-datatable__filter-wrapper",
		LoadingWrapper:  "lvt-datatable__loading-wrapper",
		LoadingIcon:     "lvt-datatable__loading-icon",
		TableWrapper:    "lvt-datatable__table-wrapper",
		Table:           "lvt-datatable__table",
		Thead:           "lvt-datatable__thead",
		Th:              "lvt-datatable__th",
		ThSortable:      "lvt-datatable__th--sortable",
		ThAlignCenter:   "lvt-datatable__th--center",
		ThAlignRight:    "lvt-datatable__th--right",
		ThContent:       "lvt-datatable__th-content",
		SortIcon:        "lvt-datatable__sort-icon",
		SortIconIdle:    "lvt-datatable__sort-icon--idle",
		Tbody:           "lvt-datatable__tbody",
		Td:              "lvt-datatable__td",
		TdCompact:       "lvt-datatable__td--compact",
		TdAlignCenter:   "lvt-datatable__td--center",
		TdAlignRight:    "lvt-datatable__td--right",
		RowStriped:      "lvt-datatable__row--striped",
		RowHover:        "lvt-datatable__row--hoverable",
		RowSelected:     "lvt-datatable__row--selected",
		Checkbox:        "lvt-datatable__checkbox",
		EmptyCell:       "lvt-datatable__empty",
		Pagination:      "lvt-datatable__pagination",
		PageInfo:        "lvt-datatable__page-info",
		PageInfoBold:    "lvt-datatable__page-info-bold",
		PageBtn:         "lvt-datatable__page-btn",
		PageBtnDisabled: "lvt-datatable__page-btn--disabled",
	}
}
