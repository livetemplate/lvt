package parser

import (
	"strings"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    int
		wantErr bool
	}{
		{
			name:  "valid fields",
			input: []string{"name:string", "email:string", "age:int"},
			want:  3,
		},
		{
			name:  "single field",
			input: []string{"title:string"},
			want:  1,
		},
		{
			name:    "empty input",
			input:   []string{},
			wantErr: true,
		},
		{
			name:    "invalid format - missing colon",
			input:   []string{"name"},
			wantErr: true,
		},
		{
			name:    "invalid format - empty name",
			input:   []string{":string"},
			wantErr: true,
		},
		{
			name:    "invalid format - empty type",
			input:   []string{"name:"},
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   []string{"id:uuid"},
			wantErr: true,
		},
		{
			name:  "select field with options",
			input: []string{"status:select:active,inactive,pending"},
			want:  1,
		},
		{
			name:  "mixed fields with select",
			input: []string{"name:string", "status:select:open,closed", "priority:int"},
			want:  3,
		},
		{
			name:    "select field without options",
			input:   []string{"status:select"},
			wantErr: true,
		},
		{
			name:    "select field with empty options",
			input:   []string{"status:select:"},
			wantErr: true,
		},
		{
			name:  "select field with adjacent empty options filtered",
			input: []string{"status:select:active,,inactive"},
			want:  1, // empty options are filtered, 2 valid remain
		},
		{
			name:  "select field with whitespace-only option filtered",
			input: []string{"status:select:active, ,inactive"},
			want:  1, // whitespace-only options are filtered, 2 valid remain
		},
		{
			name:    "select field with single option after filtering",
			input:   []string{"status:select:active"},
			wantErr: true, // only 1 valid option, need at least 2
		},
		{
			name:    "select field with all empty options",
			input:   []string{"status:select:,,"},
			wantErr: true, // 0 valid options
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := ParseFields(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(fields) != tt.want {
				t.Errorf("got %d fields, want %d", len(fields), tt.want)
			}
		})
	}
}

func TestParseFieldsSelectProperties(t *testing.T) {
	fields, err := ParseFields([]string{"status:select:active,inactive,pending"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	f := fields[0]
	if f.Name != "status" {
		t.Errorf("expected name 'status', got %q", f.Name)
	}
	if f.GoType != "string" {
		t.Errorf("expected GoType 'string', got %q", f.GoType)
	}
	if f.SQLType != "TEXT" {
		t.Errorf("expected SQLType 'TEXT', got %q", f.SQLType)
	}
	if !f.IsSelect {
		t.Error("expected IsSelect to be true")
	}
	if len(f.SelectOptions) != 3 {
		t.Fatalf("expected 3 options, got %d", len(f.SelectOptions))
	}
	expected := []string{"active", "inactive", "pending"}
	for i, opt := range f.SelectOptions {
		if opt != expected[i] {
			t.Errorf("option %d: expected %q, got %q", i, expected[i], opt)
		}
	}
}

func TestParseFieldsSelectWithOtherFields(t *testing.T) {
	fields, err := ParseFields([]string{"name:string", "status:select:open,closed", "count:int"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	// First field: normal string
	if fields[0].IsSelect {
		t.Error("field 'name' should not be IsSelect")
	}

	// Second field: select
	if !fields[1].IsSelect {
		t.Error("field 'status' should be IsSelect")
	}
	if len(fields[1].SelectOptions) != 2 {
		t.Errorf("expected 2 options, got %d", len(fields[1].SelectOptions))
	}

	// Third field: normal int
	if fields[2].IsSelect {
		t.Error("field 'count' should not be IsSelect")
	}
	if fields[2].GoType != "int64" {
		t.Errorf("expected GoType 'int64', got %q", fields[2].GoType)
	}
}

func TestMapType(t *testing.T) {
	tests := []struct {
		input        string
		wantGo       string
		wantSQL      string
		wantTextarea bool
		wantErr      bool
	}{
		{"string", "string", "TEXT", false, false},
		{"str", "string", "TEXT", false, false},
		{"text", "string", "TEXT", true, false},
		{"textarea", "string", "TEXT", true, false},
		{"longtext", "string", "TEXT", true, false},
		{"int", "int64", "INTEGER", false, false},
		{"integer", "int64", "INTEGER", false, false},
		{"bool", "bool", "BOOLEAN", false, false},
		{"boolean", "bool", "BOOLEAN", false, false},
		{"float", "float64", "REAL", false, false},
		{"float64", "float64", "REAL", false, false},
		{"time", "time.Time", "DATETIME", false, false},
		{"datetime", "time.Time", "DATETIME", false, false},
		{"uuid", "", "", false, true},
		{"unknown", "", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			goType, sqlType, isTextarea, err := MapType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if goType != tt.wantGo {
				t.Errorf("got Go type %s, want %s", goType, tt.wantGo)
			}
			if sqlType != tt.wantSQL {
				t.Errorf("got SQL type %s, want %s", sqlType, tt.wantSQL)
			}
			if isTextarea != tt.wantTextarea {
				t.Errorf("got isTextarea %v, want %v", isTextarea, tt.wantTextarea)
			}
		})
	}
}

func TestFieldsToGoStruct(t *testing.T) {
	fields := []Field{
		{Name: "name", GoType: "string"},
		{Name: "age", GoType: "int64"},
	}

	result := FieldsToGoStruct(fields)

	// Check for expected field declarations
	if !strings.Contains(result, "Name string") {
		t.Error("expected 'Name string' in struct")
	}
	if !strings.Contains(result, "Age int64") {
		t.Error("expected 'Age int64' in struct")
	}
	if !strings.Contains(result, "json:\"name\"") {
		t.Error("expected json tag for name")
	}
	if !strings.Contains(result, "json:\"age\"") {
		t.Error("expected json tag for age")
	}
}

func TestFieldsToSQLColumns(t *testing.T) {
	fields := []Field{
		{Name: "name", SQLType: "TEXT"},
		{Name: "age", SQLType: "INTEGER"},
	}

	result := FieldsToSQLColumns(fields)

	if !strings.Contains(result, "name TEXT NOT NULL") {
		t.Error("expected 'name TEXT NOT NULL' in SQL")
	}
	if !strings.Contains(result, "age INTEGER NOT NULL") {
		t.Error("expected 'age INTEGER NOT NULL' in SQL")
	}
}
