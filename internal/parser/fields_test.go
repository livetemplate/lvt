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
		{
			name:  "email field type",
			input: []string{"email:email"},
			want:  1,
		},
		{
			name:  "url field type",
			input: []string{"website:url"},
			want:  1,
		},
		{
			name:  "phone field type",
			input: []string{"phone:phone"},
			want:  1,
		},
		{
			name:  "tel field type",
			input: []string{"phone:tel"},
			want:  1,
		},
		{
			name:  "password field type",
			input: []string{"secret:password"},
			want:  1,
		},
		{
			name:  "mixed fields with new types",
			input: []string{"name:string", "email:email", "phone:phone", "website:url", "secret:password"},
			want:  5,
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
		{"email", "string", "TEXT", false, false},
		{"url", "string", "TEXT", false, false},
		{"phone", "string", "TEXT", false, false},
		{"tel", "string", "TEXT", false, false},
		{"password", "string", "TEXT", false, false},
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

func TestGetFieldMetadata(t *testing.T) {
	tests := []struct {
		fieldType     string
		wantValidate  string
		wantInputType string
		wantMinLen    int
		wantPassword  bool
		wantStep      string
	}{
		{"email", "required,email", "email", 3, false, ""},
		{"url", "required,url", "url", 0, false, ""},
		{"phone", "required", "tel", 0, false, ""},
		{"tel", "required", "tel", 0, false, ""},
		{"password", "required,min=8", "password", 8, true, ""},
		{"string", "required,min=3", "text", 3, false, ""},
		{"str", "required,min=3", "text", 3, false, ""},
		{"int", "required", "number", 0, false, ""},
		{"bool", "", "checkbox", 0, false, ""},
		{"float", "required", "number", 0, false, "0.01"},
		{"text", "required,min=3", "text", 3, false, ""},
		{"unknown_type", "", "text", 0, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.fieldType, func(t *testing.T) {
			m := GetFieldMetadata(tt.fieldType)
			if m.ValidateTag != tt.wantValidate {
				t.Errorf("ValidateTag = %q, want %q", m.ValidateTag, tt.wantValidate)
			}
			if m.HTMLInputType != tt.wantInputType {
				t.Errorf("HTMLInputType = %q, want %q", m.HTMLInputType, tt.wantInputType)
			}
			if m.HTMLMinLength != tt.wantMinLen {
				t.Errorf("HTMLMinLength = %d, want %d", m.HTMLMinLength, tt.wantMinLen)
			}
			if m.IsPassword != tt.wantPassword {
				t.Errorf("IsPassword = %v, want %v", m.IsPassword, tt.wantPassword)
			}
			if m.HTMLStep != tt.wantStep {
				t.Errorf("HTMLStep = %q, want %q", m.HTMLStep, tt.wantStep)
			}
		})
	}
}

func TestParseFieldsMetadata(t *testing.T) {
	fields, err := ParseFields([]string{"email:email", "secret:password", "website:url", "phone:tel"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}

	// email field
	if fields[0].Metadata.ValidateTag != "required,email" {
		t.Errorf("email ValidateTag = %q, want %q", fields[0].Metadata.ValidateTag, "required,email")
	}
	if fields[0].Metadata.HTMLInputType != "email" {
		t.Errorf("email HTMLInputType = %q, want %q", fields[0].Metadata.HTMLInputType, "email")
	}

	// password field
	if fields[1].Metadata.ValidateTag != "required,min=8" {
		t.Errorf("password ValidateTag = %q, want %q", fields[1].Metadata.ValidateTag, "required,min=8")
	}
	if !fields[1].Metadata.IsPassword {
		t.Error("password field should have IsPassword=true")
	}
	if fields[1].Metadata.HTMLMinLength != 8 {
		t.Errorf("password HTMLMinLength = %d, want 8", fields[1].Metadata.HTMLMinLength)
	}

	// url field
	if fields[2].Metadata.ValidateTag != "required,url" {
		t.Errorf("url ValidateTag = %q, want %q", fields[2].Metadata.ValidateTag, "required,url")
	}

	// tel field
	if fields[3].Metadata.HTMLInputType != "tel" {
		t.Errorf("tel HTMLInputType = %q, want %q", fields[3].Metadata.HTMLInputType, "tel")
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
