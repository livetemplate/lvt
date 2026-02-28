package generator

import "testing"

func TestDetectUsedComponents_AlwaysUsesModalAndToast(t *testing.T) {
	data := ResourceData{
		Fields: []FieldData{
			{Name: "name", GoType: "string"},
		},
	}
	usage := DetectUsedComponents(data)

	if !usage.UseModal {
		t.Error("expected UseModal to be true (delete confirmation is always needed)")
	}
	if !usage.UseToast {
		t.Error("expected UseToast to be true (CRUD feedback is always needed)")
	}
	if usage.UseDropdown {
		t.Error("expected UseDropdown to be false when no select fields present")
	}
}

func TestDetectUsedComponents_DropdownWithSelectField(t *testing.T) {
	data := ResourceData{
		Fields: []FieldData{
			{Name: "name", GoType: "string"},
			{Name: "status", GoType: "string", IsSelect: true, SelectOptions: []string{"active", "inactive"}},
		},
	}
	usage := DetectUsedComponents(data)

	if !usage.UseDropdown {
		t.Error("expected UseDropdown to be true when select fields present")
	}
}

func TestDetectUsedComponents_NoFields(t *testing.T) {
	data := ResourceData{}
	usage := DetectUsedComponents(data)

	if !usage.UseModal {
		t.Error("expected UseModal to be true even with no fields")
	}
	if !usage.UseToast {
		t.Error("expected UseToast to be true even with no fields")
	}
	if usage.UseDropdown {
		t.Error("expected UseDropdown to be false with no fields")
	}
}
