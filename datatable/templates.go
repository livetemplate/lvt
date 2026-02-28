package datatable

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/livetemplate/components/base"
)

// templateFS contains all datatable template files embedded at compile time.
//
//go:embed templates/*.tmpl
var templateFS embed.FS

// Templates returns the datatable component's template set for registration
// with the LiveTemplate framework.
//
// Example usage in main.go:
//
//	import "github.com/livetemplate/components/datatable"
//
//	tmpl, err := livetemplate.New("app",
//	    livetemplate.WithComponentTemplates(datatable.Templates()),
//	)
//
// Available templates:
//   - "lvt:datatable:default:v1" - Data table with sorting/pagination
func Templates() *base.TemplateSet {
	return base.NewTemplateSet(templateFS, "templates/*.tmpl", "datatable").
		WithFuncs(template.FuncMap{
			// mod returns the remainder of a/b - used for zebra striping (alternating row colors)
			"mod": func(a, b int) int {
				return a % b
			},
			// isSortedAsc checks if datatable is sorted ascending by column.
			// Works with both *DataTable and map representations.
			"isSortedAsc": func(dt interface{}, columnID string) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.IsSortedAsc(columnID)
				}
				// Handle map representation (from JSON)
				if m, ok := dt.(map[string]interface{}); ok {
					sortColumn, _ := m["SortColumn"].(string)
					sortDir, _ := m["SortDirection"].(float64) // JSON numbers are float64
					return sortColumn == columnID && int(sortDir) == int(SortAsc)
				}
				return false
			},
			// isSortedDesc checks if datatable is sorted descending by column.
			// Works with both *DataTable and map representations.
			"isSortedDesc": func(dt interface{}, columnID string) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.IsSortedDesc(columnID)
				}
				// Handle map representation (from JSON)
				if m, ok := dt.(map[string]interface{}); ok {
					sortColumn, _ := m["SortColumn"].(string)
					sortDir, _ := m["SortDirection"].(float64)
					return sortColumn == columnID && int(sortDir) == int(SortDesc)
				}
				return false
			},
			// getCellValue gets a cell value from a row.
			// Works with both Row and map representations.
			"getCellValue": func(row interface{}, columnID string) interface{} {
				if r, ok := row.(Row); ok {
					return r.GetCellValue(columnID)
				}
				// Handle map representation
				if m, ok := row.(map[string]interface{}); ok {
					if data, ok := m["Data"].(map[string]interface{}); ok {
						return data[columnID]
					}
				}
				return nil
			},
			// getRowID gets the ID from a row.
			// Works with both Row and map representations.
			"getRowID": func(row interface{}) string {
				if r, ok := row.(Row); ok {
					return r.ID
				}
				if m, ok := row.(map[string]interface{}); ok {
					if id, ok := m["ID"].(string); ok {
						return id
					}
				}
				return ""
			},
			// isRowSelected checks if a row is selected.
			// Works with both Row and map representations.
			"isRowSelected": func(row interface{}) bool {
				if r, ok := row.(Row); ok {
					return r.Selected
				}
				if m, ok := row.(map[string]interface{}); ok {
					if selected, ok := m["Selected"].(bool); ok {
						return selected
					}
				}
				return false
			},
			// isRowDisabled checks if a row is disabled.
			// Works with both Row and map representations.
			"isRowDisabled": func(row interface{}) bool {
				if r, ok := row.(Row); ok {
					return r.Disabled
				}
				if m, ok := row.(map[string]interface{}); ok {
					if disabled, ok := m["Disabled"].(bool); ok {
						return disabled
					}
				}
				return false
			},
			// dtID gets the datatable ID.
			// Works with both *DataTable and map representations.
			"dtID": func(dt interface{}) string {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.ID()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if id, ok := m["Id"].(string); ok {
						return id
					}
				}
				return ""
			},
			// dtPageSize gets the page size from datatable.
			"dtPageSize": func(dt interface{}) int {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.PageSize
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if ps, ok := m["PageSize"].(float64); ok {
						return int(ps)
					}
				}
				return 0
			},
			// dtStartIndex gets the start index.
			"dtStartIndex": func(dt interface{}) int {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.StartIndex()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if si, ok := m["StartIndex"].(float64); ok {
						return int(si)
					}
				}
				return 1
			},
			// dtEndIndex gets the end index.
			"dtEndIndex": func(dt interface{}) int {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.EndIndex()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if ei, ok := m["EndIndex"].(float64); ok {
						return int(ei)
					}
				}
				return 0
			},
			// dtTotalRows gets the total row count.
			"dtTotalRows": func(dt interface{}) int {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.TotalRows()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if tr, ok := m["TotalRows"].(float64); ok {
						return int(tr)
					}
				}
				return 0
			},
			// dtHasPrev checks if there's a previous page.
			"dtHasPrev": func(dt interface{}) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.HasPreviousPage()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if hp, ok := m["HasPreviousPage"].(bool); ok {
						return hp
					}
				}
				return false
			},
			// dtHasNext checks if there's a next page.
			"dtHasNext": func(dt interface{}) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.HasNextPage()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if hn, ok := m["HasNextPage"].(bool); ok {
						return hn
					}
				}
				return false
			},
			// dtVisibleColumns gets visible columns.
			"dtVisibleColumns": func(dt interface{}) []Column {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.VisibleColumns()
				}
				// For map representation, try to extract from computed field
				if m, ok := dt.(map[string]interface{}); ok {
					if vc, ok := m["VisibleColumns"].([]interface{}); ok {
						cols := make([]Column, 0, len(vc))
						for _, c := range vc {
							if cm, ok := c.(map[string]interface{}); ok {
								col := Column{
									ID:       getMapString(cm, "ID"),
									Label:    getMapString(cm, "Label"),
									Sortable: getMapBool(cm, "Sortable"),
									Width:    getMapString(cm, "Width"),
									Align:    getMapString(cm, "Align"),
								}
								cols = append(cols, col)
							}
						}
						return cols
					}
				}
				return nil
			},
			// dtPageRows gets current page rows.
			"dtPageRows": func(dt interface{}) interface{} {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.GetPageRows()
				}
				// For map representation, return the PageRows field
				if m, ok := dt.(map[string]interface{}); ok {
					if pr, ok := m["PageRows"]; ok {
						return pr
					}
				}
				return nil
			},
			// dtIsEmpty checks if datatable is empty.
			"dtIsEmpty": func(dt interface{}) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.IsEmpty()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if ie, ok := m["IsEmpty"].(bool); ok {
						return ie
					}
				}
				return true
			},
			// dtAllSelected checks if all rows are selected.
			"dtAllSelected": func(dt interface{}) bool {
				if datatable, ok := dt.(*DataTable); ok {
					return datatable.AllSelected()
				}
				if m, ok := dt.(map[string]interface{}); ok {
					if as, ok := m["AllSelected"].(bool); ok {
						return as
					}
				}
				return false
			},
			// colID extracts column ID from Column or map.
			"colID": func(col interface{}) string {
				if c, ok := col.(Column); ok {
					return c.ID
				}
				if m, ok := col.(map[string]interface{}); ok {
					return getMapString(m, "ID")
				}
				return ""
			},
			// colLabel extracts column label.
			"colLabel": func(col interface{}) string {
				if c, ok := col.(Column); ok {
					return c.Label
				}
				if m, ok := col.(map[string]interface{}); ok {
					return getMapString(m, "Label")
				}
				return ""
			},
			// colSortable checks if column is sortable.
			"colSortable": func(col interface{}) bool {
				if c, ok := col.(Column); ok {
					return c.Sortable
				}
				if m, ok := col.(map[string]interface{}); ok {
					return getMapBool(m, "Sortable")
				}
				return false
			},
			// colWidth gets column width.
			"colWidth": func(col interface{}) string {
				if c, ok := col.(Column); ok {
					return c.Width
				}
				if m, ok := col.(map[string]interface{}); ok {
					return getMapString(m, "Width")
				}
				return ""
			},
			// colAlign gets column alignment.
			"colAlign": func(col interface{}) string {
				if c, ok := col.(Column); ok {
					return c.Align
				}
				if m, ok := col.(map[string]interface{}); ok {
					return getMapString(m, "Align")
				}
				return ""
			},
			// Print for debugging
			"debugType": func(v interface{}) string {
				return fmt.Sprintf("%T", v)
			},
		})
}

// Helper functions for map access
func getMapString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getMapBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}
