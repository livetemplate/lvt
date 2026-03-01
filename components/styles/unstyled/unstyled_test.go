package unstyled

import (
	"reflect"
	"strings"
	"testing"
)

func TestAllStyleStructsHaveNoEmptyStringFields(t *testing.T) {
	adapter := &Adapter{}

	tests := []struct {
		name   string
		styles any
	}{
		{"AccordionStyles", adapter.AccordionStyles()},
		{"AutocompleteStyles", adapter.AutocompleteStyles()},
		{"BreadcrumbsStyles", adapter.BreadcrumbsStyles()},
		{"BreadcrumbItemStyles", adapter.BreadcrumbItemStyles()},
		{"DatatableStyles", adapter.DatatableStyles()},
		{"DatepickerStyles", adapter.DatepickerStyles()},
		{"DrawerStyles", adapter.DrawerStyles()},
		{"DropdownStyles", adapter.DropdownStyles()},
		{"MenuStyles", adapter.MenuStyles()},
		{"ModalStyles", adapter.ModalStyles()},
		{"ConfirmModalStyles", adapter.ConfirmModalStyles()},
		{"SheetStyles", adapter.SheetStyles()},
		{"PopoverStyles", adapter.PopoverStyles()},
		{"ProgressStyles", adapter.ProgressStyles()},
		{"CircularProgressStyles", adapter.CircularProgressStyles()},
		{"SpinnerStyles", adapter.SpinnerStyles()},
		{"RatingStyles", adapter.RatingStyles()},
		{"SkeletonStyles", adapter.SkeletonStyles()},
		{"AvatarSkeletonStyles", adapter.AvatarSkeletonStyles()},
		{"CardSkeletonStyles", adapter.CardSkeletonStyles()},
		{"TabsStyles", adapter.TabsStyles()},
		{"TagsInputStyles", adapter.TagsInputStyles()},
		{"TimelineStyles", adapter.TimelineStyles()},
		{"TimelineItemStyles", adapter.TimelineItemStyles()},
		{"TimepickerStyles", adapter.TimepickerStyles()},
		{"ToastStyles", adapter.ToastStyles()},
		{"ToggleStyles", adapter.ToggleStyles()},
		{"CheckboxStyles", adapter.CheckboxStyles()},
		{"TooltipStyles", adapter.TooltipStyles()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.styles)
			typ := v.Type()
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				fieldName := typ.Field(i).Name
				if field.Kind() == reflect.Bool {
					continue
				}
				if field.Kind() == reflect.String && field.String() == "" {
					t.Errorf("%s.%s is empty", tt.name, fieldName)
				}
			}
		})
	}
}

func TestAllClassNamesAreBEM(t *testing.T) {
	invalid := ValidateBEM()
	if len(invalid) > 0 {
		t.Errorf("non-BEM class names found: %v", invalid)
	}
}

func TestAllClassNamesStartWithLvt(t *testing.T) {
	for _, name := range AllClassNames() {
		if !strings.HasPrefix(name, "lvt-") {
			t.Errorf("class %q does not start with lvt-", name)
		}
	}
}

func TestAdapterName(t *testing.T) {
	adapter := &Adapter{}
	if adapter.Name() != "unstyled" {
		t.Errorf("expected adapter name %q, got %q", "unstyled", adapter.Name())
	}
}

func TestStructuralFlagsAreDisabled(t *testing.T) {
	adapter := &Adapter{}

	toggle := adapter.ToggleStyles()
	if toggle.UseCustomTrack {
		t.Error("ToggleStyles.UseCustomTrack should be false for unstyled")
	}

	checkbox := adapter.CheckboxStyles()
	if checkbox.UseCustomCheckbox {
		t.Error("CheckboxStyles.UseCustomCheckbox should be false for unstyled")
	}

	confirm := adapter.ConfirmModalStyles()
	if confirm.ShowIconCircle {
		t.Error("ConfirmModalStyles.ShowIconCircle should be false for unstyled")
	}
}

func TestClassCount(t *testing.T) {
	count := ClassCount()
	if count < 200 {
		t.Errorf("expected at least 200 BEM classes, got %d", count)
	}
}
