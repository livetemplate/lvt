package datatable

import (
	"testing"
)

func TestNew(t *testing.T) {
	dt := New("test-table")

	if dt.ID() != "test-table" {
		t.Errorf("Expected ID 'test-table', got '%s'", dt.ID())
	}
	if dt.Namespace() != "datatable" {
		t.Errorf("Expected namespace 'datatable', got '%s'", dt.Namespace())
	}
	if dt.SortDirection != SortNone {
		t.Errorf("Expected SortDirection SortNone, got %v", dt.SortDirection)
	}
	if dt.EmptyMessage != "No data available" {
		t.Errorf("Expected default EmptyMessage, got '%s'", dt.EmptyMessage)
	}
	if !dt.Striped {
		t.Error("Expected Striped to be true by default")
	}
	if !dt.Hoverable {
		t.Error("Expected Hoverable to be true by default")
	}
}

func TestWithColumns(t *testing.T) {
	columns := []Column{
		{ID: "name", Label: "Name"},
		{ID: "email", Label: "Email"},
	}
	dt := New("test", WithColumns(columns))

	if len(dt.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(dt.Columns))
	}
}

func TestWithRows(t *testing.T) {
	rows := []Row{
		{ID: "1", Data: map[string]any{"name": "Alice"}},
		{ID: "2", Data: map[string]any{"name": "Bob"}},
	}
	dt := New("test", WithRows(rows))

	if len(dt.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(dt.Rows))
	}
}

func TestWithPageSize(t *testing.T) {
	dt := New("test", WithPageSize(10))
	if dt.PageSize != 10 {
		t.Errorf("Expected PageSize 10, got %d", dt.PageSize)
	}
}

func TestWithSelectable(t *testing.T) {
	dt := New("test", WithSelectable(true))
	if !dt.Selectable {
		t.Error("Expected Selectable to be true")
	}
}

func TestWithMultiSelect(t *testing.T) {
	dt := New("test", WithMultiSelect(true))
	if !dt.MultiSelect {
		t.Error("Expected MultiSelect to be true")
	}
	if !dt.Selectable {
		t.Error("Expected Selectable to be enabled with MultiSelect")
	}
}

func TestWithStriped(t *testing.T) {
	dt := New("test", WithStriped(false))
	if dt.Striped {
		t.Error("Expected Striped to be false")
	}
}

func TestWithBordered(t *testing.T) {
	dt := New("test", WithBordered(true))
	if !dt.Bordered {
		t.Error("Expected Bordered to be true")
	}
}

func TestWithCompact(t *testing.T) {
	dt := New("test", WithCompact(true))
	if !dt.Compact {
		t.Error("Expected Compact to be true")
	}
}

func TestWithEmptyMessage(t *testing.T) {
	dt := New("test", WithEmptyMessage("No users found"))
	if dt.EmptyMessage != "No users found" {
		t.Errorf("Expected EmptyMessage 'No users found', got '%s'", dt.EmptyMessage)
	}
}

func TestWithStyled(t *testing.T) {
	dt := New("test", WithStyled(false))
	if dt.IsStyled() {
		t.Error("Expected IsStyled to be false")
	}
}

func TestWithSort(t *testing.T) {
	dt := New("test", WithSort("name", SortAsc))
	if dt.SortColumn != "name" {
		t.Errorf("Expected SortColumn 'name', got '%s'", dt.SortColumn)
	}
	if dt.SortDirection != SortAsc {
		t.Errorf("Expected SortDirection SortAsc, got %v", dt.SortDirection)
	}
}

func TestWithFilter(t *testing.T) {
	dt := New("test", WithFilter("search"))
	if dt.FilterValue != "search" {
		t.Errorf("Expected FilterValue 'search', got '%s'", dt.FilterValue)
	}
}

func TestWithLoading(t *testing.T) {
	dt := New("test", WithLoading(true))
	if !dt.Loading {
		t.Error("Expected Loading to be true")
	}
}

func TestSort(t *testing.T) {
	dt := New("test")

	// First sort - ascending
	dt.Sort("name")
	if dt.SortColumn != "name" {
		t.Errorf("Expected SortColumn 'name', got '%s'", dt.SortColumn)
	}
	if dt.SortDirection != SortAsc {
		t.Errorf("Expected SortDirection SortAsc, got %v", dt.SortDirection)
	}

	// Same column - toggle to descending
	dt.Sort("name")
	if dt.SortDirection != SortDesc {
		t.Errorf("Expected SortDirection SortDesc, got %v", dt.SortDirection)
	}

	// Different column - reset to ascending
	dt.Sort("email")
	if dt.SortColumn != "email" {
		t.Errorf("Expected SortColumn 'email', got '%s'", dt.SortColumn)
	}
	if dt.SortDirection != SortAsc {
		t.Errorf("Expected SortDirection SortAsc for new column, got %v", dt.SortDirection)
	}
}

func TestClearSort(t *testing.T) {
	dt := New("test", WithSort("name", SortAsc))
	dt.ClearSort()

	if dt.SortColumn != "" {
		t.Error("Expected SortColumn to be empty")
	}
	if dt.SortDirection != SortNone {
		t.Error("Expected SortDirection to be SortNone")
	}
}

func TestSetFilter(t *testing.T) {
	dt := New("test")
	dt.Page = 5

	dt.SetFilter("search")

	if dt.FilterValue != "search" {
		t.Errorf("Expected FilterValue 'search', got '%s'", dt.FilterValue)
	}
	if dt.Page != 0 {
		t.Error("Expected Page to reset to 0")
	}
}

func TestClearFilter(t *testing.T) {
	dt := New("test", WithFilter("search"))
	dt.Page = 5

	dt.ClearFilter()

	if dt.FilterValue != "" {
		t.Error("Expected FilterValue to be empty")
	}
	if dt.Page != 0 {
		t.Error("Expected Page to reset to 0")
	}
}

func TestPagination(t *testing.T) {
	rows := make([]Row, 25)
	for i := 0; i < 25; i++ {
		rows[i] = Row{ID: string(rune('a' + i))}
	}
	dt := New("test", WithRows(rows), WithPageSize(10))

	// Initial state
	if dt.Page != 0 {
		t.Errorf("Expected initial Page 0, got %d", dt.Page)
	}
	if dt.TotalPages() != 3 {
		t.Errorf("Expected 3 pages, got %d", dt.TotalPages())
	}

	// Next page
	dt.NextPage()
	if dt.Page != 1 {
		t.Errorf("Expected Page 1, got %d", dt.Page)
	}

	// Previous page
	dt.PreviousPage()
	if dt.Page != 0 {
		t.Errorf("Expected Page 0, got %d", dt.Page)
	}

	// Go to page
	dt.GoToPage(2)
	if dt.Page != 2 {
		t.Errorf("Expected Page 2, got %d", dt.Page)
	}

	// First page
	dt.FirstPage()
	if dt.Page != 0 {
		t.Errorf("Expected Page 0, got %d", dt.Page)
	}

	// Last page
	dt.LastPage()
	if dt.Page != 2 {
		t.Errorf("Expected Page 2 (last), got %d", dt.Page)
	}
}

func TestPaginationBounds(t *testing.T) {
	rows := make([]Row, 10)
	dt := New("test", WithRows(rows), WithPageSize(10))

	// Can't go to next when on last page
	dt.NextPage()
	if dt.Page != 0 {
		t.Error("Expected Page to stay at 0 when already on last page")
	}

	// Can't go to previous when on first page
	dt.PreviousPage()
	if dt.Page != 0 {
		t.Error("Expected Page to stay at 0")
	}

	// Invalid page numbers
	dt.GoToPage(-1)
	if dt.Page != 0 {
		t.Error("Expected Page to stay at 0 for negative page")
	}

	dt.GoToPage(100)
	if dt.Page != 0 {
		t.Error("Expected Page to stay at 0 for out of range page")
	}
}

func TestHasNextPreviousPage(t *testing.T) {
	rows := make([]Row, 25)
	dt := New("test", WithRows(rows), WithPageSize(10))

	if !dt.HasNextPage() {
		t.Error("Expected HasNextPage to be true on first page")
	}
	if dt.HasPreviousPage() {
		t.Error("Expected HasPreviousPage to be false on first page")
	}

	dt.Page = 1
	if !dt.HasNextPage() {
		t.Error("Expected HasNextPage to be true on middle page")
	}
	if !dt.HasPreviousPage() {
		t.Error("Expected HasPreviousPage to be true on middle page")
	}

	dt.Page = 2
	if dt.HasNextPage() {
		t.Error("Expected HasNextPage to be false on last page")
	}
	if !dt.HasPreviousPage() {
		t.Error("Expected HasPreviousPage to be true on last page")
	}
}

func TestGetPageRows(t *testing.T) {
	rows := make([]Row, 25)
	for i := 0; i < 25; i++ {
		rows[i] = Row{ID: string(rune('a' + i))}
	}
	dt := New("test", WithRows(rows), WithPageSize(10))

	// First page
	pageRows := dt.GetPageRows()
	if len(pageRows) != 10 {
		t.Errorf("Expected 10 rows on first page, got %d", len(pageRows))
	}

	// Last page (partial)
	dt.Page = 2
	pageRows = dt.GetPageRows()
	if len(pageRows) != 5 {
		t.Errorf("Expected 5 rows on last page, got %d", len(pageRows))
	}
}

func TestRowSelection(t *testing.T) {
	rows := []Row{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}
	dt := New("test", WithRows(rows), WithSelectable(true), WithMultiSelect(true))

	// Select row
	dt.SelectRow("1")
	if !dt.IsRowSelected("1") {
		t.Error("Expected row '1' to be selected")
	}

	// Multi-select
	dt.SelectRow("2")
	if dt.SelectedCount() != 2 {
		t.Errorf("Expected 2 selected rows, got %d", dt.SelectedCount())
	}

	// Deselect
	dt.DeselectRow("1")
	if dt.IsRowSelected("1") {
		t.Error("Expected row '1' to be deselected")
	}

	// Toggle
	dt.ToggleRowSelection("2")
	if dt.IsRowSelected("2") {
		t.Error("Expected row '2' to be toggled off")
	}
}

func TestSingleSelection(t *testing.T) {
	rows := []Row{
		{ID: "1"},
		{ID: "2"},
	}
	dt := New("test", WithRows(rows), WithSelectable(true))
	// MultiSelect is false

	dt.SelectRow("1")
	dt.SelectRow("2")

	if dt.SelectedCount() != 1 {
		t.Errorf("Expected only 1 selected with single select, got %d", dt.SelectedCount())
	}
	if !dt.IsRowSelected("2") {
		t.Error("Expected last selected row '2' to be selected")
	}
}

func TestSelectDeselectAll(t *testing.T) {
	rows := []Row{
		{ID: "1"},
		{ID: "2"},
		{ID: "3", Disabled: true},
	}
	dt := New("test", WithRows(rows), WithSelectable(true), WithMultiSelect(true))

	dt.SelectAll()
	if dt.SelectedCount() != 2 {
		t.Errorf("Expected 2 selected (excluding disabled), got %d", dt.SelectedCount())
	}

	dt.DeselectAll()
	if dt.SelectedCount() != 0 {
		t.Errorf("Expected 0 selected after DeselectAll, got %d", dt.SelectedCount())
	}
}

func TestHasSelectionAndAllSelected(t *testing.T) {
	rows := []Row{
		{ID: "1"},
		{ID: "2"},
	}
	dt := New("test", WithRows(rows), WithSelectable(true), WithMultiSelect(true))

	if dt.HasSelection() {
		t.Error("Expected HasSelection to be false initially")
	}
	if dt.AllSelected() {
		t.Error("Expected AllSelected to be false initially")
	}

	dt.SelectRow("1")
	if !dt.HasSelection() {
		t.Error("Expected HasSelection to be true")
	}
	if dt.AllSelected() {
		t.Error("Expected AllSelected to be false with partial selection")
	}

	dt.SelectAll()
	if !dt.AllSelected() {
		t.Error("Expected AllSelected to be true")
	}
}

func TestGetSelectedRows(t *testing.T) {
	rows := []Row{
		{ID: "1", Data: map[string]any{"name": "Alice"}},
		{ID: "2", Data: map[string]any{"name": "Bob"}},
	}
	dt := New("test", WithRows(rows), WithSelectable(true), WithMultiSelect(true))

	dt.SelectRow("1")
	selected := dt.GetSelectedRows()

	if len(selected) != 1 {
		t.Errorf("Expected 1 selected row, got %d", len(selected))
	}
	if selected[0].ID != "1" {
		t.Error("Expected selected row to be '1'")
	}
}

func TestIsSortedBy(t *testing.T) {
	dt := New("test")

	if dt.IsSortedBy("name") {
		t.Error("Expected IsSortedBy to be false initially")
	}

	dt.Sort("name")
	if !dt.IsSortedBy("name") {
		t.Error("Expected IsSortedBy to be true")
	}
	if !dt.IsSortedAsc("name") {
		t.Error("Expected IsSortedAsc to be true")
	}
	if dt.IsSortedDesc("name") {
		t.Error("Expected IsSortedDesc to be false")
	}

	dt.Sort("name") // Toggle to desc
	if dt.IsSortedAsc("name") {
		t.Error("Expected IsSortedAsc to be false")
	}
	if !dt.IsSortedDesc("name") {
		t.Error("Expected IsSortedDesc to be true")
	}
}

func TestGetColumn(t *testing.T) {
	columns := []Column{
		{ID: "name", Label: "Name"},
		{ID: "email", Label: "Email"},
	}
	dt := New("test", WithColumns(columns))

	col := dt.GetColumn("name")
	if col == nil {
		t.Fatal("Expected to find column 'name'")
	}
	if col.Label != "Name" {
		t.Errorf("Expected Label 'Name', got '%s'", col.Label)
	}

	if dt.GetColumn("notfound") != nil {
		t.Error("Expected nil for not found column")
	}
}

func TestVisibleColumns(t *testing.T) {
	columns := []Column{
		{ID: "name", Label: "Name"},
		{ID: "email", Label: "Email", Hidden: true},
		{ID: "role", Label: "Role"},
	}
	dt := New("test", WithColumns(columns))

	visible := dt.VisibleColumns()
	if len(visible) != 2 {
		t.Errorf("Expected 2 visible columns, got %d", len(visible))
	}
}

func TestShowHideColumn(t *testing.T) {
	columns := []Column{
		{ID: "name", Label: "Name"},
	}
	dt := New("test", WithColumns(columns))

	dt.HideColumn("name")
	if !dt.Columns[0].Hidden {
		t.Error("Expected column to be hidden")
	}

	dt.ShowColumn("name")
	if dt.Columns[0].Hidden {
		t.Error("Expected column to be visible")
	}
}

func TestSetLoading(t *testing.T) {
	dt := New("test")

	dt.SetLoading(true)
	if !dt.Loading {
		t.Error("Expected Loading to be true")
	}

	dt.SetLoading(false)
	if dt.Loading {
		t.Error("Expected Loading to be false")
	}
}

func TestSetData(t *testing.T) {
	rows := []Row{{ID: "1"}}
	dt := New("test", WithRows(rows), WithSelectable(true), WithMultiSelect(true))
	dt.SelectRow("1")

	newRows := []Row{{ID: "a"}, {ID: "b"}}
	dt.SetData(newRows)

	if len(dt.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(dt.Rows))
	}
	if dt.SelectedCount() != 0 {
		t.Error("Expected selection to be cleared")
	}
}

func TestIsEmpty(t *testing.T) {
	dt := New("test")

	if !dt.IsEmpty() {
		t.Error("Expected IsEmpty to be true")
	}

	dt.Rows = []Row{{ID: "1"}}
	if dt.IsEmpty() {
		t.Error("Expected IsEmpty to be false")
	}
}

func TestStartEndIndex(t *testing.T) {
	rows := make([]Row, 25)
	dt := New("test", WithRows(rows), WithPageSize(10))

	if dt.StartIndex() != 1 {
		t.Errorf("Expected StartIndex 1, got %d", dt.StartIndex())
	}
	if dt.EndIndex() != 10 {
		t.Errorf("Expected EndIndex 10, got %d", dt.EndIndex())
	}

	dt.Page = 2 // Last page
	if dt.StartIndex() != 21 {
		t.Errorf("Expected StartIndex 21, got %d", dt.StartIndex())
	}
	if dt.EndIndex() != 25 {
		t.Errorf("Expected EndIndex 25, got %d", dt.EndIndex())
	}
}

func TestRowGetCellValue(t *testing.T) {
	row := Row{
		ID:   "1",
		Data: map[string]any{"name": "Alice", "age": 30},
	}

	name := row.GetCellValue("name")
	if name != "Alice" {
		t.Errorf("Expected 'Alice', got '%v'", name)
	}

	missing := row.GetCellValue("missing")
	if missing != nil {
		t.Error("Expected nil for missing column")
	}
}

func TestRowGetCellString(t *testing.T) {
	row := Row{
		ID:   "1",
		Data: map[string]any{"name": "Alice", "age": 30},
	}

	name := row.GetCellString("name")
	if name != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", name)
	}

	// Non-string returns empty
	age := row.GetCellString("age")
	if age != "" {
		t.Errorf("Expected empty string for non-string, got '%s'", age)
	}
}

func TestTotalRowsAndPages(t *testing.T) {
	dt := New("test", WithPageSize(10))

	if dt.TotalRows() != 0 {
		t.Errorf("Expected TotalRows 0, got %d", dt.TotalRows())
	}
	if dt.TotalPages() != 1 {
		t.Errorf("Expected TotalPages 1 (minimum), got %d", dt.TotalPages())
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
