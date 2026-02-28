package datatable

// Option is a functional option for configuring data tables.
type Option func(*DataTable)

// WithColumns sets the table columns.
func WithColumns(columns []Column) Option {
	return func(dt *DataTable) {
		dt.Columns = columns
	}
}

// WithRows sets the table data.
func WithRows(rows []Row) Option {
	return func(dt *DataTable) {
		dt.Rows = rows
	}
}

// WithPageSize sets the number of rows per page.
func WithPageSize(size int) Option {
	return func(dt *DataTable) {
		dt.PageSize = size
	}
}

// WithSelectable enables row selection.
func WithSelectable(selectable bool) Option {
	return func(dt *DataTable) {
		dt.Selectable = selectable
	}
}

// WithMultiSelect enables multiple row selection.
func WithMultiSelect(multiSelect bool) Option {
	return func(dt *DataTable) {
		dt.MultiSelect = multiSelect
		if multiSelect {
			dt.Selectable = true
		}
	}
}

// WithStriped enables alternating row colors.
func WithStriped(striped bool) Option {
	return func(dt *DataTable) {
		dt.Striped = striped
	}
}

// WithHoverable enables row hover effect.
func WithHoverable(hoverable bool) Option {
	return func(dt *DataTable) {
		dt.Hoverable = hoverable
	}
}

// WithBordered adds borders to cells.
func WithBordered(bordered bool) Option {
	return func(dt *DataTable) {
		dt.Bordered = bordered
	}
}

// WithCompact reduces padding.
func WithCompact(compact bool) Option {
	return func(dt *DataTable) {
		dt.Compact = compact
	}
}

// WithEmptyMessage sets the message shown when no data.
func WithEmptyMessage(message string) Option {
	return func(dt *DataTable) {
		dt.EmptyMessage = message
	}
}

// WithStyled enables Tailwind CSS styling for the component.
func WithStyled(styled bool) Option {
	return func(dt *DataTable) {
		dt.SetStyled(styled)
	}
}

// WithSort sets initial sorting.
func WithSort(columnID string, direction SortDirection) Option {
	return func(dt *DataTable) {
		dt.SortColumn = columnID
		dt.SortDirection = direction
	}
}

// WithFilter sets initial filter.
func WithFilter(value string) Option {
	return func(dt *DataTable) {
		dt.FilterValue = value
	}
}

// WithLoading sets initial loading state.
func WithLoading(loading bool) Option {
	return func(dt *DataTable) {
		dt.Loading = loading
	}
}
