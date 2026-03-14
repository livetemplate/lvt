package generator

import (
	"strings"
	"text/template"

	"github.com/livetemplate/lvt/internal/kits"
	"github.com/livetemplate/lvt/internal/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FieldDataFromFields converts parsed fields to FieldData for template rendering.
func FieldDataFromFields(fields []parser.Field) []FieldData {
	fd := make([]FieldData, len(fields))
	for i, f := range fields {
		fd[i] = FieldData{
			Name:            f.Name,
			GoType:          f.GoType,
			SQLType:         f.SQLType,
			IsReference:     f.IsReference,
			ReferencedTable: f.ReferencedTable,
			OnDelete:        f.OnDelete,
			IsTextarea:      f.IsTextarea,
			IsSelect:        f.IsSelect,
			SelectOptions:   f.SelectOptions,
		}
	}
	return fd
}

type ResourceData struct {
	PackageName          string
	ModuleName           string
	ResourceName         string // Input name, capitalized (e.g., "Users" or "User")
	ResourceNameLower    string // Input name, lowercase (e.g., "users" or "user")
	ResourceNameSingular string // Singular, capitalized (e.g., "User")
	ResourceNamePlural   string // Plural, capitalized (e.g., "Users")
	TableName            string // Plural table name (e.g., "users")
	Fields               []FieldData
	Kit                  *kits.KitInfo  // CSS framework kit (new)
	CSSFramework         string         // CSS framework name: "tailwind", "bulma", "pico", "none" (for backward compatibility)
	DevMode              bool           // Use local client library instead of CDN
	PaginationMode       string         // Pagination mode: "infinite", "load-more", "prev-next", "numbers"
	PageSize             int            // Page size for pagination
	EditMode             string         // Edit mode: "modal", "page"
	Components           ComponentUsage // Which UI components this resource uses
	Styles               string         // Style adapter: "tailwind", "unstyled"
	StylesImportPath     string         // computed import path for style adapter (empty if no components need it)

	// Embedded child resource fields (set when --parent is used)
	ParentResource         string // Parent resource name, lowercase plural (e.g., "posts"). Empty = standalone.
	ParentPackageName      string // Parent package name (e.g., "posts")
	ParentResourceSingular string // Parent resource singular capitalized (e.g., "Post")
	ParentReferenceField   string // FK field referencing parent (e.g., "post_id"), auto-detected
	IsEmbedded             bool   // True when generating as embedded child
}

// NonReferenceFields returns fields excluding the parent reference field.
// Used in embedded templates to omit the parent FK from forms.
func (d ResourceData) NonReferenceFields() []FieldData {
	if d.ParentReferenceField == "" {
		return d.Fields
	}
	var result []FieldData
	for _, f := range d.Fields {
		if f.Name != d.ParentReferenceField {
			result = append(result, f)
		}
	}
	return result
}

type FieldData struct {
	Name            string
	GoType          string
	SQLType         string
	IsReference     bool
	ReferencedTable string
	OnDelete        string
	IsTextarea      bool     // true if field should render as textarea
	IsSelect        bool     // true if field should render as <select>
	SelectOptions   []string // options for select fields
}

type AppData struct {
	AppName      string
	ModuleName   string
	Kit          *kits.KitInfo // CSS framework kit (new)
	DevMode      bool          // Use local client library instead of CDN
	CSSFramework string        // CSS framework name for home page (for backward compatibility)
	Styles       string        // Style adapter: "tailwind", "unstyled"
}

var funcMap = template.FuncMap{
	"title":        cases.Title(language.English).String,
	"lower":        strings.ToLower,
	"upper":        strings.ToUpper,
	"camelCase":    toCamelCase,
	"displayField": getDisplayField,
	"singularize":  singularizeForTemplate,
}

// singularizeForTemplate wraps singularize for use in templates.
func singularizeForTemplate(s string) string {
	return singularize(strings.ToLower(s))
}

// toCamelCase converts snake_case to CamelCase following Go conventions
// Common initialisms like ID, URL, HTTP are kept in all caps
func toCamelCase(s string) string {
	// Common initialisms that should be all caps
	initialisms := map[string]bool{
		"id": true, "url": true, "http": true, "https": true,
		"api": true, "uri": true, "sql": true, "json": true,
		"xml": true, "html": true, "css": true, "js": true,
	}

	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			lower := strings.ToLower(part)
			if initialisms[lower] {
				parts[i] = strings.ToUpper(part)
			} else {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
	}
	return strings.Join(parts, "")
}

// getDisplayField returns the primary display field from a list of fields.
// Priority: title > name > first non-reference string field > first non-reference field > first field
func getDisplayField(fields []FieldData) FieldData {
	if len(fields) == 0 {
		return FieldData{Name: "id", GoType: "string"}
	}

	// Check for "title" field first
	for _, field := range fields {
		if strings.ToLower(field.Name) == "title" {
			return field
		}
	}

	// Check for "name" field second
	for _, field := range fields {
		if strings.ToLower(field.Name) == "name" {
			return field
		}
	}

	// Prefer the first non-reference string field (most likely human-readable)
	for _, field := range fields {
		if !field.IsReference && field.GoType == "string" {
			return field
		}
	}

	// Fall back to first non-reference field
	for _, field := range fields {
		if !field.IsReference {
			return field
		}
	}

	// Last resort: first field
	return fields[0]
}
