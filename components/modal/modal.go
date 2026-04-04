// Package modal provides modal/dialog components.
//
// Available variants:
//   - New() creates a modal dialog (template: "lvt:modal:default:v1")
//   - NewConfirm() creates a confirmation dialog (template: "lvt:modal:confirm:v1")
//   - NewSheet() creates a slide-in sheet (template: "lvt:modal:sheet:v1")
//
// Open/close is handled client-side via onclick handlers and CSS classes.
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
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
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

// Styles returns the resolved ModalStyles for this component.
func (m *Modal) Styles() styles.ModalStyles {
	if s, ok := m.StyleData().(styles.ModalStyles); ok {
		return s
	}
	adapter := styles.ForStyled(m.IsStyled())
	if adapter == nil {
		return styles.ModalStyles{}
	}
	s := adapter.ModalStyles()
	m.SetStyleData(s)
	return s
}

// SizeClass returns CSS class for modal size.
func (m *Modal) SizeClass() string {
	s := m.Styles()
	switch m.Size {
	case SizeSm:
		return s.SizeSm
	case SizeLg:
		return s.SizeLg
	case SizeXl:
		return s.SizeXl
	case SizeFull:
		return s.SizeFull
	default: // md
		return s.SizeMd
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

// Styles returns the resolved ConfirmModalStyles for this component.
func (c *ConfirmModal) Styles() styles.ConfirmModalStyles {
	if s, ok := c.StyleData().(styles.ConfirmModalStyles); ok {
		return s
	}
	adapter := styles.ForStyled(c.IsStyled())
	if adapter == nil {
		return styles.ConfirmModalStyles{}
	}
	s := adapter.ConfirmModalStyles()
	c.SetStyleData(s)
	return s
}

// ConfirmButtonClass returns CSS class for confirm button.
func (c *ConfirmModal) ConfirmButtonClass() string {
	s := c.Styles()
	if c.Destructive {
		return s.ConfirmDestructive
	}
	return s.ConfirmDefault
}

// IconClass returns CSS class for the icon.
func (c *ConfirmModal) IconClass() string {
	s := c.Styles()
	switch c.Icon {
	case "warning":
		if c.Destructive {
			return s.IconWarningDestructive
		}
		return s.IconWarning
	case "info":
		return s.IconInfo
	case "success":
		return s.IconSuccess
	case "error":
		return s.IconError
	default:
		return s.IconDefault
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

// Styles returns the resolved SheetStyles for this component.
func (s *SheetModal) Styles() styles.SheetStyles {
	if st, ok := s.StyleData().(styles.SheetStyles); ok {
		return st
	}
	adapter := styles.ForStyled(s.IsStyled())
	if adapter == nil {
		return styles.SheetStyles{}
	}
	st := adapter.SheetStyles()
	s.SetStyleData(st)
	return st
}

// PositionClass returns CSS classes for position.
func (s *SheetModal) PositionClass() string {
	st := s.Styles()
	switch s.Position {
	case "left":
		return st.PositionLeft
	case "right":
		return st.PositionRight
	case "top":
		return st.PositionTop
	case "bottom":
		return st.PositionBottom
	default:
		return st.PositionRight
	}
}

// SizeClass returns CSS class for sheet size.
func (s *SheetModal) SizeClass() string {
	st := s.Styles()
	if s.IsHorizontal() {
		switch s.Size {
		case SizeSm:
			return st.SizeSmH
		case SizeLg:
			return st.SizeLgH
		case SizeXl:
			return st.SizeXlH
		case SizeFull:
			return st.SizeFullH
		default:
			return st.SizeMdH
		}
	}
	// Vertical
	switch s.Size {
	case SizeSm:
		return st.SizeSmV
	case SizeLg:
		return st.SizeLgV
	case SizeXl:
		return st.SizeXlV
	case SizeFull:
		return st.SizeFullV
	default:
		return st.SizeMdV
	}
}

// TransformClass returns CSS transform for animation.
func (s *SheetModal) TransformClass() string {
	st := s.Styles()
	if s.Open {
		return st.TransformOpen
	}
	switch s.Position {
	case "left":
		return st.TransformLeft
	case "right":
		return st.TransformRight
	case "top":
		return st.TransformTop
	case "bottom":
		return st.TransformBottom
	default:
		return st.TransformRight
	}
}
