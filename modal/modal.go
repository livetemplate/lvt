// Package modal provides modal/dialog components.
//
// Available variants:
//   - New() creates a modal dialog (template: "lvt:modal:default:v1")
//   - NewConfirm() creates a confirmation dialog (template: "lvt:modal:confirm:v1")
//   - NewSheet() creates a slide-in sheet (template: "lvt:modal:sheet:v1")
//
// Required lvt-* attributes: lvt-click, lvt-click-away, lvt-modal-open, lvt-modal-close
//
// Example usage:
//
//	// In your controller/state
//	DeleteConfirm: modal.NewConfirm("delete",
//	    modal.WithTitle("Delete Item"),
//	    modal.WithMessage("Are you sure?"),
//	)
//
//	// In your template
//	{{template "lvt:modal:confirm:v1" .DeleteConfirm}}
package modal

import (
	"github.com/livetemplate/components/base"
)

// Size defines the modal size.
type Size string

const (
	SizeSm   Size = "sm"
	SizeMd   Size = "md"
	SizeLg   Size = "lg"
	SizeXl   Size = "xl"
	SizeFull Size = "full"
)

// Modal is a dialog component.
// Use template "lvt:modal:default:v1" to render.
type Modal struct {
	base.Base

	// Open indicates whether the modal is visible
	Open bool

	// Title is the modal header
	Title string

	// Size of the modal
	Size Size

	// ShowClose shows the close button
	ShowClose bool

	// CloseOnOverlay closes modal when clicking overlay
	CloseOnOverlay bool

	// CloseOnEscape closes modal on Escape key
	CloseOnEscape bool

	// Centered vertically centers the modal
	Centered bool

	// Scrollable makes the modal body scrollable
	Scrollable bool
}

// New creates a modal dialog.
//
// Example:
//
//	m := modal.New("settings",
//	    modal.WithTitle("Settings"),
//	    modal.WithSize(modal.SizeLg),
//	)
func New(id string, opts ...Option) *Modal {
	m := &Modal{
		Base:           base.NewBase(id, "modal"),
		Size:           SizeMd,
		ShowClose:      true,
		CloseOnOverlay: true,
		CloseOnEscape:  true,
		Centered:       true,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Show opens the modal.
func (m *Modal) Show() {
	m.Open = true
}

// Hide closes the modal.
func (m *Modal) Hide() {
	m.Open = false
}

// Toggle toggles modal visibility.
func (m *Modal) Toggle() {
	m.Open = !m.Open
}

// HasTitle returns true if modal has a title.
func (m *Modal) HasTitle() bool {
	return m.Title != ""
}

// HasHeader returns true if modal should show header.
func (m *Modal) HasHeader() bool {
	return m.HasTitle() || m.ShowClose
}

// SizeClass returns CSS class for modal size.
func (m *Modal) SizeClass() string {
	switch m.Size {
	case SizeSm:
		return "max-w-sm"
	case SizeLg:
		return "max-w-2xl"
	case SizeXl:
		return "max-w-4xl"
	case SizeFull:
		return "max-w-full mx-4"
	default: // md
		return "max-w-lg"
	}
}

// ConfirmModal is a confirmation dialog.
type ConfirmModal struct {
	base.Base

	// Open indicates whether the modal is visible
	Open bool

	// Title is the dialog header
	Title string

	// Message is the confirmation message
	Message string

	// ConfirmText is the confirm button text
	ConfirmText string

	// CancelText is the cancel button text
	CancelText string

	// Destructive styles the confirm button as destructive
	Destructive bool

	// Icon shows an icon in the dialog
	Icon string
}

// NewConfirm creates a confirmation dialog.
//
// Example:
//
//	c := modal.NewConfirm("delete",
//	    modal.WithConfirmTitle("Delete Item"),
//	    modal.WithConfirmMessage("Are you sure you want to delete this item?"),
//	    modal.WithConfirmDestructive(true),
//	)
func NewConfirm(id string, opts ...ConfirmOption) *ConfirmModal {
	c := &ConfirmModal{
		Base:        base.NewBase(id, "modal"),
		ConfirmText: "Confirm",
		CancelText:  "Cancel",
		Icon:        "warning",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Show opens the confirmation dialog.
func (c *ConfirmModal) Show() {
	c.Open = true
}

// Hide closes the confirmation dialog.
func (c *ConfirmModal) Hide() {
	c.Open = false
}

// HasTitle returns true if dialog has a title.
func (c *ConfirmModal) HasTitle() bool {
	return c.Title != ""
}

// HasMessage returns true if dialog has a message.
func (c *ConfirmModal) HasMessage() bool {
	return c.Message != ""
}

// HasIcon returns true if dialog has an icon.
func (c *ConfirmModal) HasIcon() bool {
	return c.Icon != ""
}

// IsDestructive returns true if action is destructive.
func (c *ConfirmModal) IsDestructive() bool {
	return c.Destructive
}

// ConfirmButtonClass returns CSS class for confirm button.
func (c *ConfirmModal) ConfirmButtonClass() string {
	if c.Destructive {
		return "bg-red-600 hover:bg-red-700 text-white"
	}
	return "bg-blue-600 hover:bg-blue-700 text-white"
}

// IconClass returns CSS class for the icon.
func (c *ConfirmModal) IconClass() string {
	switch c.Icon {
	case "warning":
		if c.Destructive {
			return "text-red-500"
		}
		return "text-yellow-500"
	case "info":
		return "text-blue-500"
	case "success":
		return "text-green-500"
	case "error":
		return "text-red-500"
	default:
		return "text-gray-500"
	}
}

// SheetModal is a slide-in panel/sheet.
type SheetModal struct {
	base.Base

	// Open indicates whether the sheet is visible
	Open bool

	// Title is the sheet header
	Title string

	// Position is where the sheet slides from (left, right, top, bottom)
	Position string

	// Size controls the sheet width/height
	Size Size

	// ShowClose shows the close button
	ShowClose bool

	// CloseOnOverlay closes sheet when clicking overlay
	CloseOnOverlay bool
}

// NewSheet creates a slide-in sheet.
//
// Example:
//
//	s := modal.NewSheet("filters",
//	    modal.WithSheetTitle("Filters"),
//	    modal.WithSheetPosition("right"),
//	)
func NewSheet(id string, opts ...SheetOption) *SheetModal {
	s := &SheetModal{
		Base:           base.NewBase(id, "modal"),
		Position:       "right",
		Size:           SizeMd,
		ShowClose:      true,
		CloseOnOverlay: true,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Show opens the sheet.
func (s *SheetModal) Show() {
	s.Open = true
}

// Hide closes the sheet.
func (s *SheetModal) Hide() {
	s.Open = false
}

// Toggle toggles sheet visibility.
func (s *SheetModal) Toggle() {
	s.Open = !s.Open
}

// HasTitle returns true if sheet has a title.
func (s *SheetModal) HasTitle() bool {
	return s.Title != ""
}

// IsLeft returns true if position is left.
func (s *SheetModal) IsLeft() bool {
	return s.Position == "left"
}

// IsRight returns true if position is right.
func (s *SheetModal) IsRight() bool {
	return s.Position == "right"
}

// IsTop returns true if position is top.
func (s *SheetModal) IsTop() bool {
	return s.Position == "top"
}

// IsBottom returns true if position is bottom.
func (s *SheetModal) IsBottom() bool {
	return s.Position == "bottom"
}

// IsHorizontal returns true if position is left or right.
func (s *SheetModal) IsHorizontal() bool {
	return s.IsLeft() || s.IsRight()
}

// IsVertical returns true if position is top or bottom.
func (s *SheetModal) IsVertical() bool {
	return s.IsTop() || s.IsBottom()
}

// PositionClass returns CSS classes for position.
func (s *SheetModal) PositionClass() string {
	switch s.Position {
	case "left":
		return "left-0 top-0 h-full"
	case "right":
		return "right-0 top-0 h-full"
	case "top":
		return "top-0 left-0 w-full"
	case "bottom":
		return "bottom-0 left-0 w-full"
	default:
		return "right-0 top-0 h-full"
	}
}

// SizeClass returns CSS class for sheet size.
func (s *SheetModal) SizeClass() string {
	if s.IsHorizontal() {
		switch s.Size {
		case SizeSm:
			return "w-64"
		case SizeLg:
			return "w-96"
		case SizeXl:
			return "w-[32rem]"
		case SizeFull:
			return "w-full"
		default:
			return "w-80"
		}
	}
	// Vertical
	switch s.Size {
	case SizeSm:
		return "h-48"
	case SizeLg:
		return "h-96"
	case SizeXl:
		return "h-[32rem]"
	case SizeFull:
		return "h-full"
	default:
		return "h-64"
	}
}

// TransformClass returns CSS transform for animation.
func (s *SheetModal) TransformClass() string {
	if s.Open {
		return "translate-x-0 translate-y-0"
	}
	switch s.Position {
	case "left":
		return "-translate-x-full"
	case "right":
		return "translate-x-full"
	case "top":
		return "-translate-y-full"
	case "bottom":
		return "translate-y-full"
	default:
		return "translate-x-full"
	}
}
