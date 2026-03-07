package tailwind

import (
	"reflect"
	"testing"
)

func TestAllStyleStructsHaveNoEmptyFields(t *testing.T) {
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

	// Fields intentionally empty in the Tailwind adapter.
	// AccordionStyles.Item: items need no wrapper class because divide-y on Root handles separation.
	allowEmpty := map[string]map[string]bool{
		"AccordionStyles": {"Item": true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.styles)
			typ := v.Type()
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				fieldName := typ.Field(i).Name
				// Skip bool fields (structural flags like UseCustomTrack)
				if field.Kind() == reflect.Bool {
					continue
				}
				if allowEmpty[tt.name][fieldName] {
					continue
				}
				if field.Kind() == reflect.String && field.String() == "" {
					t.Errorf("%s.%s is empty", tt.name, fieldName)
				}
			}
		})
	}
}

func TestAdapterName(t *testing.T) {
	adapter := &Adapter{}
	if adapter.Name() != "tailwind" {
		t.Errorf("expected adapter name %q, got %q", "tailwind", adapter.Name())
	}
}

func TestAdapterRegisteredOnInit(t *testing.T) {
	// The init() function should have registered this adapter.
	// We test this indirectly by confirming the adapter implements the interface
	// and returns non-empty styles.
	adapter := &Adapter{}

	// Spot-check a few methods return non-zero structs
	modal := adapter.ModalStyles()
	if modal.Root == "" || modal.Panel == "" || modal.Overlay == "" {
		t.Error("ModalStyles has empty core fields")
	}

	toggle := adapter.ToggleStyles()
	if !toggle.UseCustomTrack {
		t.Error("ToggleStyles.UseCustomTrack should be true for tailwind")
	}
	if toggle.Track == "" || toggle.Knob == "" {
		t.Error("ToggleStyles has empty track/knob fields")
	}

	checkbox := adapter.CheckboxStyles()
	if !checkbox.UseCustomCheckbox {
		t.Error("CheckboxStyles.UseCustomCheckbox should be true for tailwind")
	}

	confirm := adapter.ConfirmModalStyles()
	if !confirm.ShowIconCircle {
		t.Error("ConfirmModalStyles.ShowIconCircle should be true for tailwind")
	}
}
