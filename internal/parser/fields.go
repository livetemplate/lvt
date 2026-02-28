package parser

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
		}

		if strings.HasPrefix(strings.ToLower(fullType), "references:") {
			// Parse: references:table_name[:on_delete_action]
			parts := strings.Split(fullType, ":")
			if len(parts) < 2 {
				return nil, fmt.Errorf("field '%s': invalid references syntax, expected 'references:table_name'", name)
			}

			field.IsReference = true
			field.ReferencedTable = parts[1]

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

// MapType maps a user-provided type to Go and SQL types
// Also handles references syntax: references:table_name[:on_delete_action]
// Returns: goType, sqlType, isTextarea, error
func MapType(typ string) (goType, sqlType string, isTextarea bool, err error) {
	// Check if it's a reference type
	if strings.HasPrefix(strings.ToLower(typ), "references:") {
		// Format: references:table_name[:on_delete_action]
		// We return TEXT type to match our primary key type
		// The reference metadata is handled separately
		return "string", "TEXT", false, nil
	}

	switch strings.ToLower(typ) {
	case "string", "str":
		return "string", "TEXT", false, nil
	case "text", "textarea", "longtext":
		return "string", "TEXT", true, nil
	case "int", "integer":
		return "int64", "INTEGER", false, nil
	case "bool", "boolean":
		return "bool", "BOOLEAN", false, nil
	case "float", "float64", "decimal":
		return "float64", "REAL", false, nil
	case "time", "datetime", "timestamp":
		return "time.Time", "DATETIME", false, nil
	default:
		return "", "", false, fmt.Errorf("unsupported type '%s' (supported: string, text, int, bool, float, time, references:table)", typ)
	}
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
