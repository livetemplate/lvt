package generator

import "testing"

func TestGetDisplayField(t *testing.T) {
	tests := []struct {
		name     string
		fields   []FieldData
		wantName string
	}{
		{
			name:     "empty fields returns id fallback",
			fields:   nil,
			wantName: "id",
		},
		{
			name: "prefers title field",
			fields: []FieldData{
				{Name: "post_id", GoType: "string", IsReference: true},
				{Name: "title", GoType: "string"},
				{Name: "content", GoType: "string"},
			},
			wantName: "title",
		},
		{
			name: "prefers name field over others",
			fields: []FieldData{
				{Name: "post_id", GoType: "string", IsReference: true},
				{Name: "name", GoType: "string"},
				{Name: "age", GoType: "int"},
			},
			wantName: "name",
		},
		{
			name: "skips reference fields in fallback",
			fields: []FieldData{
				{Name: "post_id", GoType: "string", IsReference: true},
				{Name: "author", GoType: "string"},
				{Name: "text", GoType: "string"},
			},
			wantName: "author",
		},
		{
			name: "prefers string field over non-string",
			fields: []FieldData{
				{Name: "count", GoType: "int"},
				{Name: "description", GoType: "string"},
			},
			wantName: "description",
		},
		{
			name: "falls back to non-reference non-string field",
			fields: []FieldData{
				{Name: "post_id", GoType: "string", IsReference: true},
				{Name: "active", GoType: "bool"},
			},
			wantName: "active",
		},
		{
			name: "last resort uses reference field",
			fields: []FieldData{
				{Name: "post_id", GoType: "string", IsReference: true},
			},
			wantName: "post_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDisplayField(tt.fields)
			if got.Name != tt.wantName {
				t.Errorf("getDisplayField() = %q, want %q", got.Name, tt.wantName)
			}
		})
	}
}
