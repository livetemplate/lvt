// Package tailwind provides a Tailwind CSS style adapter.
//
// This adapter returns Tailwind CSS class names for all component style
// structs, matching the classes previously hardcoded in templates and
// CSS class methods.
package tailwind

import "github.com/livetemplate/lvt/components/styles"

func init() {
	styles.Register(&Adapter{})
}

// Adapter implements styles.StyleAdapter using Tailwind CSS classes.
type Adapter struct{}

func (a *Adapter) Name() string { return "tailwind" }

func (a *Adapter) AccordionStyles() styles.AccordionStyles       { return accordionStyles() }
func (a *Adapter) AutocompleteStyles() styles.AutocompleteStyles { return autocompleteStyles() }
func (a *Adapter) BreadcrumbsStyles() styles.BreadcrumbsStyles   { return breadcrumbsStyles() }
func (a *Adapter) BreadcrumbItemStyles() styles.BreadcrumbItemStyles {
	return breadcrumbItemStyles()
}
func (a *Adapter) DatatableStyles() styles.DatatableStyles   { return datatableStyles() }
func (a *Adapter) DatepickerStyles() styles.DatepickerStyles { return datepickerStyles() }
func (a *Adapter) DrawerStyles() styles.DrawerStyles         { return drawerStyles() }
func (a *Adapter) DropdownStyles() styles.DropdownStyles     { return dropdownStyles() }
func (a *Adapter) MenuStyles() styles.MenuStyles             { return menuStyles() }
func (a *Adapter) ModalStyles() styles.ModalStyles           { return modalStyles() }
func (a *Adapter) ConfirmModalStyles() styles.ConfirmModalStyles {
	return confirmModalStyles()
}
func (a *Adapter) SheetStyles() styles.SheetStyles     { return sheetStyles() }
func (a *Adapter) PopoverStyles() styles.PopoverStyles { return popoverStyles() }
func (a *Adapter) ProgressStyles() styles.ProgressStyles {
	return progressStyles()
}
func (a *Adapter) CircularProgressStyles() styles.CircularProgressStyles {
	return circularProgressStyles()
}
func (a *Adapter) SpinnerStyles() styles.SpinnerStyles         { return spinnerStyles() }
func (a *Adapter) RatingStyles() styles.RatingStyles           { return ratingStyles() }
func (a *Adapter) SkeletonStyles() styles.SkeletonStyles       { return skeletonStyles() }
func (a *Adapter) AvatarSkeletonStyles() styles.AvatarSkeletonStyles {
	return avatarSkeletonStyles()
}
func (a *Adapter) CardSkeletonStyles() styles.CardSkeletonStyles {
	return cardSkeletonStyles()
}
func (a *Adapter) TabsStyles() styles.TabsStyles           { return tabsStyles() }
func (a *Adapter) TagsInputStyles() styles.TagsInputStyles { return tagsInputStyles() }
func (a *Adapter) TimelineStyles() styles.TimelineStyles   { return timelineStyles() }
func (a *Adapter) TimelineItemStyles() styles.TimelineItemStyles {
	return timelineItemStyles()
}
func (a *Adapter) TimepickerStyles() styles.TimepickerStyles { return timepickerStyles() }
func (a *Adapter) ToastStyles() styles.ToastStyles           { return toastStyles() }
func (a *Adapter) ToggleStyles() styles.ToggleStyles         { return toggleStyles() }
func (a *Adapter) CheckboxStyles() styles.CheckboxStyles     { return checkboxStyles() }
func (a *Adapter) TooltipStyles() styles.TooltipStyles       { return tooltipStyles() }

// Verify interface compliance at compile time.
var _ styles.StyleAdapter = (*Adapter)(nil)
