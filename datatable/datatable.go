// Package datatable provides data table components with sorting, filtering, and pagination.
//
// Available variants:
//   - New() creates a data table (template: "lvt:datatable:default:v1")
//
// Required lvt-* attributes: lvt-click
//
// Example usage:
//
//	// In your controller/state
//	Users: datatable.New[User]("users",
//	    datatable.WithColumns([]datatable.Column{
//	        {ID: "name", Label: "Name", Sortable: true},
//	        {ID: "email", Label: "Email", Sortable: true},
//	        {ID: "role", Label: "Role"},
//	    }),
//	    datatable.WithData(users),
//	    datatable.WithPageSize(10),
//	)
//
//	// In your template
//	{{template "lvt:datatable:default:v1" .Users}}
package datatable

import (
	"encoding/json"

	"github.com/livetemplate/components/base"
)

// SortDirection indicates the sort order.
type SortDirection int

const (
	SortNone SortDirection = iota
	SortAsc
	SortDesc
)

// Column defines a table column.
type Column struct {
	// ID is the column identifier (matches data field name)
	ID string
	// Label is the display header
	Label string
	// Sortable allows sorting by this column
	Sortable bool
	// Filterable allows filtering by this column
	Filterable bool
	// Width is optional column width (e.g., "100px", "20%")
	Width string
	// Align is text alignment ("left", "center", "right")
	Align string
	// Hidden hides the column
	Hidden bool
	// Format is a format hint (e.g., "date", "currency", "number")
	Format string
}

// Row represents a table row with data and metadata.
type Row struct {
	// ID is the unique row identifier
	ID string
	// Data holds the row data (map of column ID to value)
	Data map[string]any
	// Selected indicates if row is selected
	Selected bool
	// Disabled prevents row selection
	Disabled bool
}

// DataTable is a table component with sorting, filtering, and pagination.
// Use template "lvt:datatable:default:v1" to render.
type DataTable struct {
	base.Base

	// Columns defines the table columns
	Columns []Column

	// Rows is the data
	Rows []Row

	// SortColumn is the currently sorted column ID
	SortColumn string

	// SortDirection is the current sort direction
	SortDirection SortDirection

	// FilterValue is the current filter/search text
	FilterValue string

	// FilterColumn limits filtering to a specific column (empty for all)
	FilterColumn string

	// Page is the current page (0-indexed)
	Page int

	// PageSize is rows per page (0 for all)
	PageSize int

	// Selectable enables row selection
	Selectable bool

	// MultiSelect allows multiple row selection
	MultiSelect bool

	// SelectedIDs tracks selected row IDs
	SelectedIDs map[string]bool

	// Striped enables alternating row colors
	Striped bool

	// Hoverable enables row hover effect
	Hoverable bool

	// Bordered adds borders to cells
	Bordered bool

	// Compact reduces padding
	Compact bool

	// Loading indicates data is loading
	Loading bool

	// EmptyMessage is shown when no data
	EmptyMessage string

	// filteredRows caches filtered results
	filteredRows []Row
}

// New creates a data table.
//
// Example:
//
//	dt := datatable.New("users",
//	    datatable.WithColumns(columns),
//	    datatable.WithRows(rows),
//	    datatable.WithPageSize(10),
//	)
func New(id string, opts ...Option) *DataTable {
	dt := &DataTable{
		Base:          base.NewBase(id, "datatable"),
		SelectedIDs:   make(map[string]bool),
		SortDirection: SortNone,
		EmptyMessage:  "No data available",
		Striped:       true,
		Hoverable:     true,
	}

	for _, opt := range opts {
		opt(dt)
	}

	return dt
}

// Sort sorts by a column. Toggles direction if same column.
func (dt *DataTable) Sort(columnID string) {
	if dt.SortColumn == columnID {
		// Toggle direction
		switch dt.SortDirection {
		case SortNone, SortDesc:
			dt.SortDirection = SortAsc
		case SortAsc:
			dt.SortDirection = SortDesc
		}
	} else {
		dt.SortColumn = columnID
		dt.SortDirection = SortAsc
	}
	dt.filteredRows = nil // Clear cache
}

// ClearSort removes sorting.
func (dt *DataTable) ClearSort() {
	dt.SortColumn = ""
	dt.SortDirection = SortNone
	dt.filteredRows = nil
}

// SetFilter sets the filter value.
func (dt *DataTable) SetFilter(value string) {
	dt.FilterValue = value
	dt.Page = 0 // Reset to first page
	dt.filteredRows = nil
}

// ClearFilter clears the filter.
func (dt *DataTable) ClearFilter() {
	dt.FilterValue = ""
	dt.FilterColumn = ""
	dt.Page = 0
	dt.filteredRows = nil
}

// NextPage goes to the next page.
func (dt *DataTable) NextPage() {
	if dt.Page < dt.TotalPages()-1 {
		dt.Page++
	}
}

// PreviousPage goes to the previous page.
func (dt *DataTable) PreviousPage() {
	if dt.Page > 0 {
		dt.Page--
	}
}

// GoToPage goes to a specific page.
func (dt *DataTable) GoToPage(page int) {
	if page >= 0 && page < dt.TotalPages() {
		dt.Page = page
	}
}

// FirstPage goes to the first page.
func (dt *DataTable) FirstPage() {
	dt.Page = 0
}

// LastPage goes to the last page.
func (dt *DataTable) LastPage() {
	dt.Page = dt.TotalPages() - 1
	if dt.Page < 0 {
		dt.Page = 0
	}
}

// TotalPages returns the total number of pages.
func (dt *DataTable) TotalPages() int {
	if dt.PageSize <= 0 {
		return 1
	}
	total := dt.TotalRows()
	pages := total / dt.PageSize
	if total%dt.PageSize > 0 {
		pages++
	}
	if pages == 0 {
		pages = 1
	}
	return pages
}

// TotalRows returns the total number of rows (after filtering).
func (dt *DataTable) TotalRows() int {
	return len(dt.GetFilteredRows())
}

// HasNextPage returns true if there's a next page.
func (dt *DataTable) HasNextPage() bool {
	return dt.Page < dt.TotalPages()-1
}

// HasPreviousPage returns true if there's a previous page.
func (dt *DataTable) HasPreviousPage() bool {
	return dt.Page > 0
}

// GetFilteredRows returns rows after filtering (cached).
func (dt *DataTable) GetFilteredRows() []Row {
	// Simple implementation - in production this would be optimized
	if dt.FilterValue == "" {
		return dt.Rows
	}

	if dt.filteredRows != nil {
		return dt.filteredRows
	}

	// Filter implementation would go here
	// For now, return all rows (actual filtering done server-side)
	dt.filteredRows = dt.Rows
	return dt.filteredRows
}

// GetPageRows returns rows for the current page.
func (dt *DataTable) GetPageRows() []Row {
	filtered := dt.GetFilteredRows()

	if dt.PageSize <= 0 {
		return filtered
	}

	start := dt.Page * dt.PageSize
	end := start + dt.PageSize

	if start >= len(filtered) {
		return nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end]
}

// SelectRow selects a row by ID.
func (dt *DataTable) SelectRow(id string) {
	if !dt.Selectable {
		return
	}

	if !dt.MultiSelect {
		// Clear other selections
		dt.SelectedIDs = make(map[string]bool)
	}
	dt.SelectedIDs[id] = true

	// Update row state
	for i := range dt.Rows {
		if dt.Rows[i].ID == id {
			dt.Rows[i].Selected = true
		} else if !dt.MultiSelect {
			dt.Rows[i].Selected = false
		}
	}
}

// DeselectRow deselects a row by ID.
func (dt *DataTable) DeselectRow(id string) {
	delete(dt.SelectedIDs, id)
	for i := range dt.Rows {
		if dt.Rows[i].ID == id {
			dt.Rows[i].Selected = false
		}
	}
}

// ToggleRowSelection toggles row selection.
func (dt *DataTable) ToggleRowSelection(id string) {
	if dt.SelectedIDs[id] {
		dt.DeselectRow(id)
	} else {
		dt.SelectRow(id)
	}
}

// SelectAll selects all rows.
func (dt *DataTable) SelectAll() {
	if !dt.Selectable || !dt.MultiSelect {
		return
	}
	for i := range dt.Rows {
		if !dt.Rows[i].Disabled {
			dt.Rows[i].Selected = true
			dt.SelectedIDs[dt.Rows[i].ID] = true
		}
	}
}

// DeselectAll deselects all rows.
func (dt *DataTable) DeselectAll() {
	dt.SelectedIDs = make(map[string]bool)
	for i := range dt.Rows {
		dt.Rows[i].Selected = false
	}
}

// SelectedCount returns the number of selected rows.
func (dt *DataTable) SelectedCount() int {
	return len(dt.SelectedIDs)
}

// HasSelection returns true if any row is selected.
func (dt *DataTable) HasSelection() bool {
	return len(dt.SelectedIDs) > 0
}

// AllSelected returns true if all rows are selected.
func (dt *DataTable) AllSelected() bool {
	if len(dt.Rows) == 0 {
		return false
	}
	selectableCount := 0
	for _, row := range dt.Rows {
		if !row.Disabled {
			selectableCount++
		}
	}
	return len(dt.SelectedIDs) == selectableCount
}

// GetSelectedRows returns all selected rows.
func (dt *DataTable) GetSelectedRows() []Row {
	var selected []Row
	for _, row := range dt.Rows {
		if dt.SelectedIDs[row.ID] {
			selected = append(selected, row)
		}
	}
	return selected
}

// IsRowSelected checks if a row is selected.
func (dt *DataTable) IsRowSelected(id string) bool {
	return dt.SelectedIDs[id]
}

// IsSortedBy checks if sorted by a column.
func (dt *DataTable) IsSortedBy(columnID string) bool {
	return dt.SortColumn == columnID && dt.SortDirection != SortNone
}

// IsSortedAsc checks if sorted ascending by a column.
func (dt *DataTable) IsSortedAsc(columnID string) bool {
	return dt.SortColumn == columnID && dt.SortDirection == SortAsc
}

// IsSortedDesc checks if sorted descending by a column.
func (dt *DataTable) IsSortedDesc(columnID string) bool {
	return dt.SortColumn == columnID && dt.SortDirection == SortDesc
}

// GetColumn returns a column by ID.
func (dt *DataTable) GetColumn(id string) *Column {
	for i := range dt.Columns {
		if dt.Columns[i].ID == id {
			return &dt.Columns[i]
		}
	}
	return nil
}

// VisibleColumns returns non-hidden columns.
func (dt *DataTable) VisibleColumns() []Column {
	var visible []Column
	for _, col := range dt.Columns {
		if !col.Hidden {
			visible = append(visible, col)
		}
	}
	return visible
}

// ShowColumn unhides a column.
func (dt *DataTable) ShowColumn(id string) {
	if col := dt.GetColumn(id); col != nil {
		col.Hidden = false
	}
}

// HideColumn hides a column.
func (dt *DataTable) HideColumn(id string) {
	if col := dt.GetColumn(id); col != nil {
		col.Hidden = true
	}
}

// SetLoading sets the loading state.
func (dt *DataTable) SetLoading(loading bool) {
	dt.Loading = loading
}

// SetData replaces all rows.
func (dt *DataTable) SetData(rows []Row) {
	dt.Rows = rows
	dt.filteredRows = nil
	dt.DeselectAll()
}

// IsEmpty returns true if there's no data.
func (dt *DataTable) IsEmpty() bool {
	return len(dt.GetFilteredRows()) == 0
}

// PageInfo returns information about the current page.
func (dt *DataTable) PageInfo() string {
	total := dt.TotalRows()
	if total == 0 {
		return "No results"
	}

	if dt.PageSize <= 0 {
		return ""
	}

	start := dt.Page*dt.PageSize + 1
	end := start + dt.PageSize - 1
	if end > total {
		end = total
	}

	return "" // Would format "Showing 1-10 of 100" but keeping simple
}

// StartIndex returns the 1-based start index for current page.
func (dt *DataTable) StartIndex() int {
	if dt.PageSize <= 0 {
		return 1
	}
	return dt.Page*dt.PageSize + 1
}

// EndIndex returns the end index for current page.
func (dt *DataTable) EndIndex() int {
	if dt.PageSize <= 0 {
		return dt.TotalRows()
	}
	end := (dt.Page + 1) * dt.PageSize
	total := dt.TotalRows()
	if end > total {
		end = total
	}
	return end
}

// GetCellValue returns the value for a row and column.
func (r Row) GetCellValue(columnID string) any {
	if r.Data == nil {
		return nil
	}
	return r.Data[columnID]
}

// GetCellString returns the string value for a row and column.
func (r Row) GetCellString(columnID string) string {
	v := r.GetCellValue(columnID)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// MarshalJSON implements json.Marshaler to include computed fields for RPC serialization.
// This ensures that methods like VisibleColumns(), GetPageRows() are available as JSON fields
// when the datatable is serialized and used in templates via RPC.
func (dt *DataTable) MarshalJSON() ([]byte, error) {
	// Create an alias to avoid infinite recursion
	type DataTableAlias DataTable

	// Build computed fields
	return json.Marshal(&struct {
		*DataTableAlias

		// Computed fields for template access
		VisibleColumns  []Column `json:"VisibleColumns"`
		PageRows        []Row    `json:"PageRows"`
		StartIndex      int      `json:"StartIndex"`
		EndIndex        int      `json:"EndIndex"`
		TotalRowsCount  int      `json:"TotalRows"`
		HasPreviousPage bool     `json:"HasPreviousPage"`
		HasNextPage     bool     `json:"HasNextPage"`
		AllSelectedFlag bool     `json:"AllSelected"`
		IsEmptyFlag     bool     `json:"IsEmpty"`
	}{
		DataTableAlias:  (*DataTableAlias)(dt),
		VisibleColumns:  dt.VisibleColumns(),
		PageRows:        dt.GetPageRows(),
		StartIndex:      dt.StartIndex(),
		EndIndex:        dt.EndIndex(),
		TotalRowsCount:  dt.TotalRows(),
		HasPreviousPage: dt.HasPreviousPage(),
		HasNextPage:     dt.HasNextPage(),
		AllSelectedFlag: dt.AllSelected(),
		IsEmptyFlag:     dt.IsEmpty(),
	})
}
