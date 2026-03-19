package generator

import (
	"testing"

	"github.com/livetemplate/lvt/internal/parser"
)

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

func TestFieldDataFromFieldsCopiesMetadata(t *testing.T) {
	fields := []parser.Field{
		{
			Name:   "email",
			GoType: "string",
			Metadata: parser.FieldMetadata{
				ValidateTag:   "required,email",
				HTMLInputType: "email",
				HTMLMinLength: 3,
			},
		},
		{
			Name:   "secret",
			GoType: "string",
			Metadata: parser.FieldMetadata{
				ValidateTag:   "required,min=8",
				HTMLInputType: "password",
				HTMLMinLength: 8,
				IsPassword:    true,
			},
		},
		{
			Name:   "price",
			GoType: "float64",
			Metadata: parser.FieldMetadata{
				ValidateTag:   "required",
				HTMLInputType: "number",
				HTMLStep:      "0.01",
			},
		},
	}

	fd := FieldDataFromFields(fields)

	if len(fd) != 3 {
		t.Fatalf("expected 3 FieldData, got %d", len(fd))
	}

	// email
	if fd[0].ValidateTag != "required,email" {
		t.Errorf("email ValidateTag = %q, want %q", fd[0].ValidateTag, "required,email")
	}
	if fd[0].HTMLInputType != "email" {
		t.Errorf("email HTMLInputType = %q, want %q", fd[0].HTMLInputType, "email")
	}
	if fd[0].HTMLMinLength != 3 {
		t.Errorf("email HTMLMinLength = %d, want 3", fd[0].HTMLMinLength)
	}

	// password
	if fd[1].ValidateTag != "required,min=8" {
		t.Errorf("password ValidateTag = %q, want %q", fd[1].ValidateTag, "required,min=8")
	}
	if !fd[1].IsPassword {
		t.Error("password IsPassword should be true")
	}
	if fd[1].HTMLMinLength != 8 {
		t.Errorf("password HTMLMinLength = %d, want 8", fd[1].HTMLMinLength)
	}

	// float
	if fd[2].HTMLStep != "0.01" {
		t.Errorf("float HTMLStep = %q, want %q", fd[2].HTMLStep, "0.01")
	}
}
