package parser

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FieldMetadata holds validation and HTML rendering metadata derived from the field type.
type FieldMetadata struct {
	ValidateTag   string // e.g. "required,email", "required,min=8"
	HTMLInputType string // e.g. "email", "url", "tel", "password", "text", "number"
	HTMLMinLength int    // 0 = not set
	HTMLMaxLength int    // 0 = not set
	HTMLStep      string // e.g. "0.01" for floats
	IsPassword    bool   // suppress value echo in edit forms
}

type Field struct {
	Name            string
	Type            string
	GoType          string
	SQLType         string
	IsReference     bool
	ReferencedTable string
	OnDelete        string   // CASCADE, SET NULL, RESTRICT, etc.
	IsTextarea      bool     // true if field should render as textarea
	IsSelect        bool     // true if field should render as <select>
	SelectOptions   []string // options for select fields
	IsFile          bool     // true if field is a file upload
	IsImage         bool     // true if field is an image upload (subset of file)
	Metadata        FieldMetadata
}

// ParseFields parses field definitions in the format "name:type name2:type2"
func ParseFields(args []string) ([]Field, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no fields provided")
	}

	var fields []Field
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid field format '%s', expected 'name:type'", arg)
		}

		name := strings.TrimSpace(parts[0])
		typ := strings.TrimSpace(parts[1])

		if name == "" {
			return nil, fmt.Errorf("field name cannot be empty")
		}
		if typ == "" {
			return nil, fmt.Errorf("field type cannot be empty for field '%s'", name)
		}

		// Handle select type: name:select:opt1,opt2,opt3
		if strings.ToLower(typ) == "select" {
			if len(parts) < 3 || strings.TrimSpace(parts[2]) == "" {
				return nil, fmt.Errorf("field '%s': select type requires options, e.g., 'status:select:active,inactive,pending'", name)
			}
			rawOptions := strings.Split(parts[2], ",")
			var options []string
			for _, o := range rawOptions {
				if s := strings.TrimSpace(o); s != "" {
					options = append(options, s)
				}
			}
			if len(options) < 2 {
				return nil, fmt.Errorf("field '%s': select requires at least 2 non-empty options", name)
			}
			fields = append(fields, Field{
				Name:          name,
				Type:          "select",
				GoType:        "string",
				SQLType:       "TEXT",
				IsSelect:      true,
				SelectOptions: options,
				Metadata: FieldMetadata{
					ValidateTag:   "required",
					HTMLInputType: "text",
				},
			})
			continue
		}

		// Handle file/image types: name:file or name:image
		lowerTyp := strings.ToLower(typ)
		if lowerTyp == "file" || lowerTyp == "image" {
			fields = append(fields, Field{
				Name:    name,
				Type:    lowerTyp,
				GoType:  "string",
				SQLType: "TEXT",
				IsFile:  true,
				IsImage: lowerTyp == "image",
				Metadata: FieldMetadata{
					HTMLInputType: "file",
				},
			})
			continue
		}

		// Rejoin remaining parts for types that use colons (e.g., references:table:cascade)
		fullType := strings.Join(parts[1:], ":")

		// Validate type
		goType, sqlType, isTextarea, err := MapType(fullType)
		if err != nil {
			return nil, fmt.Errorf("field '%s': %w", name, err)
		}

		// Parse reference metadata if it's a reference type
		field := Field{
			Name:       name,
			Type:       fullType,
			GoType:     goType,
			SQLType:    sqlType,
			IsTextarea: isTextarea,
			Metadata:   GetFieldMetadata(typ),
		}

		if strings.HasPrefix(strings.ToLower(fullType), "references:") {
			// Parse: references:table_name[:on_delete_action]
			parts := strings.Split(fullType, ":")
			if len(parts) < 2 {
				return nil, fmt.Errorf("field '%s': invalid references syntax, expected 'references:table_name'", name)
			}

			field.IsReference = true
			field.ReferencedTable = parts[1]
			field.Metadata = FieldMetadata{ValidateTag: "required", HTMLInputType: "text"}

			// Default to CASCADE
			field.OnDelete = "CASCADE"

			// Check for custom on_delete action
			if len(parts) > 2 {
				action := strings.ToUpper(parts[2])
				switch action {
				case "CASCADE", "SET NULL", "RESTRICT", "NO ACTION", "SET_NULL":
					if action == "SET_NULL" {
						action = "SET NULL"
					}
					field.OnDelete = action
				default:
					return nil, fmt.Errorf("field '%s': invalid ON DELETE action '%s' (supported: CASCADE, SET_NULL, RESTRICT, NO_ACTION)", name, parts[2])
				}
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// fieldTypeInfo holds the combined type mapping and metadata for a field type.
type fieldTypeInfo struct {
	GoType     string
	SQLType    string
	IsTextarea bool
	Metadata   FieldMetadata
}

// fieldTypeTable is the single source of truth for field type → Go/SQL/metadata mapping.
var fieldTypeTable = map[string]fieldTypeInfo{
	"string":    {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required,min=3", HTMLInputType: "text", HTMLMinLength: 3}},
	"str":       {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required,min=3", HTMLInputType: "text", HTMLMinLength: 3}},
	"text":      {GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: FieldMetadata{ValidateTag: "required,min=3", HTMLInputType: "text", HTMLMinLength: 3}},
	"textarea":  {GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: FieldMetadata{ValidateTag: "required,min=3", HTMLInputType: "text", HTMLMinLength: 3}},
	"longtext":  {GoType: "string", SQLType: "TEXT", IsTextarea: true, Metadata: FieldMetadata{ValidateTag: "required,min=3", HTMLInputType: "text", HTMLMinLength: 3}},
	"int":       {GoType: "int64", SQLType: "INTEGER", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "number"}},
	"integer":   {GoType: "int64", SQLType: "INTEGER", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "number"}},
	"bool":      {GoType: "bool", SQLType: "BOOLEAN", Metadata: FieldMetadata{HTMLInputType: "checkbox"}},
	"boolean":   {GoType: "bool", SQLType: "BOOLEAN", Metadata: FieldMetadata{HTMLInputType: "checkbox"}},
	"float":     {GoType: "float64", SQLType: "REAL", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "number", HTMLStep: "0.01"}},
	"float64":   {GoType: "float64", SQLType: "REAL", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "number", HTMLStep: "0.01"}},
	"decimal":   {GoType: "float64", SQLType: "REAL", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "number", HTMLStep: "0.01"}},
	"time":      {GoType: "time.Time", SQLType: "DATETIME", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "datetime-local"}},
	"datetime":  {GoType: "time.Time", SQLType: "DATETIME", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "datetime-local"}},
	"timestamp": {GoType: "time.Time", SQLType: "DATETIME", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "datetime-local"}},
	"email":     {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required,email", HTMLInputType: "email", HTMLMinLength: 3}},
	"url":       {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required,url", HTMLInputType: "url"}},
	"phone":     {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "tel"}},
	"tel":       {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required", HTMLInputType: "tel"}},
	"password":  {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{ValidateTag: "required,min=8", HTMLInputType: "password", HTMLMinLength: 8, IsPassword: true}},
	"file":      {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{HTMLInputType: "file"}},
	"image":     {GoType: "string", SQLType: "TEXT", Metadata: FieldMetadata{HTMLInputType: "file"}},
}

// supportedTypes returns a comma-separated list of primary supported type names.
func supportedTypes() string {
	// Show primary names (not aliases) in a logical order
	return "string, text, int, bool, float, time, email, url, phone, tel, password, file, image"
}

// MapType maps a user-provided type to Go and SQL types.
// Also handles references syntax: references:table_name[:on_delete_action]
// Returns: goType, sqlType, isTextarea, error
func MapType(typ string) (goType, sqlType string, isTextarea bool, err error) {
	if strings.HasPrefix(strings.ToLower(typ), "references:") {
		return "string", "TEXT", false, nil
	}

	info, ok := fieldTypeTable[strings.ToLower(typ)]
	if !ok {
		return "", "", false, fmt.Errorf("unsupported type '%s' (supported: %s, references:table)", typ, supportedTypes())
	}
	return info.GoType, info.SQLType, info.IsTextarea, nil
}

// GetFieldMetadata returns validation and HTML metadata for a given field type.
func GetFieldMetadata(fieldType string) FieldMetadata {
	if info, ok := fieldTypeTable[strings.ToLower(fieldType)]; ok {
		return info.Metadata
	}
	return FieldMetadata{HTMLInputType: "text"}
}

// FieldsToGoStruct generates Go struct field declarations
func FieldsToGoStruct(fields []Field) string {
	var sb strings.Builder
	for _, f := range fields {
		// Capitalize first letter for export
		fieldName := cases.Title(language.English).String(f.Name)
		sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, f.GoType, f.Name))
	}
	return sb.String()
}

// FieldsToSQLColumns generates SQL column definitions
func FieldsToSQLColumns(fields []Field) string {
	var sb strings.Builder
	for i, f := range fields {
		sb.WriteString(fmt.Sprintf("  %s %s NOT NULL", f.Name, f.SQLType))
		if i < len(fields)-1 {
			sb.WriteString(",\n")
		}
	}
	return sb.String()
}
