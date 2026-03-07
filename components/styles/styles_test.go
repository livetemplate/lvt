package styles

import (
	"sort"
	"testing"
)

// mockAdapter is a minimal adapter for testing the registry.
type mockAdapter struct {
	name string
}

func (m *mockAdapter) Name() string                           { return m.name }
func (m *mockAdapter) AccordionStyles() AccordionStyles       { return AccordionStyles{} }
func (m *mockAdapter) AutocompleteStyles() AutocompleteStyles { return AutocompleteStyles{} }
func (m *mockAdapter) BreadcrumbsStyles() BreadcrumbsStyles   { return BreadcrumbsStyles{} }
func (m *mockAdapter) BreadcrumbItemStyles() BreadcrumbItemStyles {
	return BreadcrumbItemStyles{}
}
func (m *mockAdapter) DatatableStyles() DatatableStyles   { return DatatableStyles{} }
func (m *mockAdapter) DatepickerStyles() DatepickerStyles { return DatepickerStyles{} }
func (m *mockAdapter) DrawerStyles() DrawerStyles         { return DrawerStyles{} }
func (m *mockAdapter) DropdownStyles() DropdownStyles     { return DropdownStyles{} }
func (m *mockAdapter) MenuStyles() MenuStyles             { return MenuStyles{} }
func (m *mockAdapter) ModalStyles() ModalStyles           { return ModalStyles{} }
func (m *mockAdapter) ConfirmModalStyles() ConfirmModalStyles {
	return ConfirmModalStyles{}
}
func (m *mockAdapter) SheetStyles() SheetStyles   { return SheetStyles{} }
func (m *mockAdapter) PopoverStyles() PopoverStyles { return PopoverStyles{} }
func (m *mockAdapter) ProgressStyles() ProgressStyles {
	return ProgressStyles{}
}
func (m *mockAdapter) CircularProgressStyles() CircularProgressStyles {
	return CircularProgressStyles{}
}
func (m *mockAdapter) SpinnerStyles() SpinnerStyles         { return SpinnerStyles{} }
func (m *mockAdapter) RatingStyles() RatingStyles           { return RatingStyles{} }
func (m *mockAdapter) SkeletonStyles() SkeletonStyles       { return SkeletonStyles{} }
func (m *mockAdapter) AvatarSkeletonStyles() AvatarSkeletonStyles {
	return AvatarSkeletonStyles{}
}
func (m *mockAdapter) CardSkeletonStyles() CardSkeletonStyles {
	return CardSkeletonStyles{}
}
func (m *mockAdapter) TabsStyles() TabsStyles           { return TabsStyles{} }
func (m *mockAdapter) TagsInputStyles() TagsInputStyles { return TagsInputStyles{} }
func (m *mockAdapter) TimelineStyles() TimelineStyles   { return TimelineStyles{} }
func (m *mockAdapter) TimelineItemStyles() TimelineItemStyles {
	return TimelineItemStyles{}
}
func (m *mockAdapter) TimepickerStyles() TimepickerStyles { return TimepickerStyles{} }
func (m *mockAdapter) ToastStyles() ToastStyles           { return ToastStyles{} }
func (m *mockAdapter) ToggleStyles() ToggleStyles         { return ToggleStyles{} }
func (m *mockAdapter) CheckboxStyles() CheckboxStyles     { return CheckboxStyles{} }
func (m *mockAdapter) TooltipStyles() TooltipStyles       { return TooltipStyles{} }

func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	adapters = map[string]StyleAdapter{}
	current = nil
}

func TestRegister(t *testing.T) {
	resetRegistry()

	a := &mockAdapter{name: "test"}
	Register(a)

	if got := Get("test"); got != a {
		t.Errorf("Get(\"test\") = %v, want %v", got, a)
	}
}

func TestRegisterSetsFirstAsDefault(t *testing.T) {
	resetRegistry()

	a1 := &mockAdapter{name: "first"}
	a2 := &mockAdapter{name: "second"}

	Register(a1)
	Register(a2)

	if got := Default(); got != a1 {
		t.Errorf("Default() = %v, want first registered adapter %v", got, a1)
	}
}

func TestSetDefault(t *testing.T) {
	resetRegistry()

	a1 := &mockAdapter{name: "first"}
	a2 := &mockAdapter{name: "second"}

	Register(a1)
	Register(a2)
	SetDefault(a2)

	if got := Default(); got != a2 {
		t.Errorf("Default() = %v, want %v", got, a2)
	}
}

func TestGetUnregistered(t *testing.T) {
	resetRegistry()

	if got := Get("nonexistent"); got != nil {
		t.Errorf("Get(\"nonexistent\") = %v, want nil", got)
	}
}

func TestDefaultWithNoAdapters(t *testing.T) {
	resetRegistry()

	if got := Default(); got != nil {
		t.Errorf("Default() with no adapters = %v, want nil", got)
	}
}

func TestForStyled(t *testing.T) {
	resetRegistry()

	tw := &mockAdapter{name: "tailwind"}
	un := &mockAdapter{name: "unstyled"}

	Register(tw)
	Register(un)
	SetDefault(tw)

	t.Run("styled=true returns default", func(t *testing.T) {
		got := ForStyled(true)
		if got != tw {
			t.Errorf("ForStyled(true) = %v, want tailwind adapter", got)
		}
	})

	t.Run("styled=false returns unstyled", func(t *testing.T) {
		got := ForStyled(false)
		if got != un {
			t.Errorf("ForStyled(false) = %v, want unstyled adapter", got)
		}
	})
}

func TestForStyledFallsBackToDefault(t *testing.T) {
	resetRegistry()

	tw := &mockAdapter{name: "tailwind"}
	Register(tw)

	got := ForStyled(false)
	if got != tw {
		t.Errorf("ForStyled(false) without unstyled adapter = %v, want default %v", got, tw)
	}
}

func TestNames(t *testing.T) {
	resetRegistry()

	Register(&mockAdapter{name: "alpha"})
	Register(&mockAdapter{name: "beta"})

	names := Names()
	sort.Strings(names)

	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("Names() = %v, want [alpha, beta]", names)
	}
}

func TestCount(t *testing.T) {
	resetRegistry()

	if Count() != 0 {
		t.Errorf("Count() with empty registry = %d, want 0", Count())
	}

	Register(&mockAdapter{name: "one"})
	Register(&mockAdapter{name: "two"})

	if Count() != 2 {
		t.Errorf("Count() = %d, want 2", Count())
	}
}

func TestRegisterOverwrite(t *testing.T) {
	resetRegistry()

	a1 := &mockAdapter{name: "same"}
	a2 := &mockAdapter{name: "same"}

	Register(a1)
	Register(a2)

	if got := Get("same"); got != a2 {
		t.Errorf("Re-registering should overwrite: Get(\"same\") = %v, want %v", got, a2)
	}

	if Count() != 1 {
		t.Errorf("Count after overwrite = %d, want 1", Count())
	}
}
