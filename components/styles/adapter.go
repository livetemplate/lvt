// Package styles provides the style adapter system for swappable CSS frameworks.
//
// Components resolve their CSS classes through a StyleAdapter instead of
// hardcoding Tailwind classes. Register adapters (e.g. "tailwind", "unstyled")
// at init time and switch between them globally or per-component.
//
// Usage:
//
//	// Set the default adapter (typically in main or init)
//	styles.SetDefault(styles.Get("tailwind"))
//
//	// Components automatically resolve styles
//	modal := modal.New("my-modal")
//	fmt.Println(modal.Styles().Panel) // Tailwind classes
package styles

// StyleAdapter provides CSS classes for all component types.
// Each method returns a typed struct with class names for that component.
type StyleAdapter interface {
	// Name returns the adapter identifier (e.g. "tailwind", "unstyled").
	Name() string

	// Component style methods — one per component variant.

	AccordionStyles() AccordionStyles
	AutocompleteStyles() AutocompleteStyles
	BreadcrumbsStyles() BreadcrumbsStyles
	BreadcrumbItemStyles() BreadcrumbItemStyles
	DatatableStyles() DatatableStyles
	DatepickerStyles() DatepickerStyles
	DrawerStyles() DrawerStyles
	DropdownStyles() DropdownStyles
	MenuStyles() MenuStyles
	ModalStyles() ModalStyles
	ConfirmModalStyles() ConfirmModalStyles
	SheetStyles() SheetStyles
	PopoverStyles() PopoverStyles
	ProgressStyles() ProgressStyles
	CircularProgressStyles() CircularProgressStyles
	SpinnerStyles() SpinnerStyles
	RatingStyles() RatingStyles
	SkeletonStyles() SkeletonStyles
	AvatarSkeletonStyles() AvatarSkeletonStyles
	CardSkeletonStyles() CardSkeletonStyles
	TabsStyles() TabsStyles
	TagsInputStyles() TagsInputStyles
	TimelineStyles() TimelineStyles
	TimelineItemStyles() TimelineItemStyles
	TimepickerStyles() TimepickerStyles
	ToastStyles() ToastStyles
	ToggleStyles() ToggleStyles
	CheckboxStyles() CheckboxStyles
	TooltipStyles() TooltipStyles
}
