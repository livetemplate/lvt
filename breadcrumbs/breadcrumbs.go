// Package breadcrumbs provides breadcrumb navigation components for the LiveTemplate framework.
//
// Breadcrumbs display the current page location within a navigational hierarchy.
//
// # Available Templates
//
//   - lvt:breadcrumbs:default:v1 - Standard breadcrumb navigation
//
// # Basic Usage
//
//	bc := breadcrumbs.New("nav",
//	    breadcrumbs.WithItems(
//	        breadcrumbs.NewItem("home", breadcrumbs.WithItemLabel("Home"), breadcrumbs.WithItemHref("/")),
//	        breadcrumbs.NewItem("products", breadcrumbs.WithItemLabel("Products"), breadcrumbs.WithItemHref("/products")),
//	        breadcrumbs.NewItem("current", breadcrumbs.WithItemLabel("Laptop"), breadcrumbs.WithItemCurrent(true)),
//	    ),
//	)
//
//	{{template "lvt:breadcrumbs:default:v1" .Breadcrumbs}}
package breadcrumbs

import (
	"github.com/livetemplate/components/base"
)

// Separator defines the visual separator between items.
type Separator string

const (
	// SeparatorSlash uses "/" as separator.
	SeparatorSlash Separator = "slash"
	// SeparatorChevron uses ">" as separator.
	SeparatorChevron Separator = "chevron"
	// SeparatorArrow uses "→" as separator.
	SeparatorArrow Separator = "arrow"
	// SeparatorDot uses "•" as separator.
	SeparatorDot Separator = "dot"
)

// Size defines breadcrumb text size.
type Size string

const (
	SizeSm   Size = "sm"
	SizeMd   Size = "md"
	SizeLg   Size = "lg"
)

// Breadcrumbs represents a breadcrumb navigation component.
type Breadcrumbs struct {
	base.Base

	// Items in the breadcrumb trail.
	Items []*BreadcrumbItem

	// Separator between items.
	Separator Separator

	// Size of the breadcrumb text.
	Size Size

	// ShowHome displays a home icon as the first item.
	ShowHome bool

	// HomeHref is the URL for the home link.
	HomeHref string

	// Collapsible enables collapsing middle items.
	Collapsible bool

	// MaxVisible is the max items to show when collapsed.
	MaxVisible int
}

// Option configures a Breadcrumbs.
type Option func(*Breadcrumbs)

// New creates a new Breadcrumbs with the given ID and options.
func New(id string, opts ...Option) *Breadcrumbs {
	bc := &Breadcrumbs{
		Base:        base.NewBase(id, "breadcrumbs"),
		Items:       make([]*BreadcrumbItem, 0),
		Separator:   SeparatorChevron,
		Size:        SizeMd,
		ShowHome:    false,
		HomeHref:    "/",
		Collapsible: false,
		MaxVisible:  3,
	}
	for _, opt := range opts {
		opt(bc)
	}
	return bc
}

// AddItem adds an item to the breadcrumbs.
func (bc *Breadcrumbs) AddItem(item *BreadcrumbItem) {
	bc.Items = append(bc.Items, item)
}

// HasItems returns true if there are items.
func (bc *Breadcrumbs) HasItems() bool {
	return len(bc.Items) > 0
}

// ItemCount returns the number of items.
func (bc *Breadcrumbs) ItemCount() int {
	return len(bc.Items)
}

// LastItem returns the last item or nil.
func (bc *Breadcrumbs) LastItem() *BreadcrumbItem {
	if len(bc.Items) == 0 {
		return nil
	}
	return bc.Items[len(bc.Items)-1]
}

// IsCollapsed returns true if items should be collapsed.
func (bc *Breadcrumbs) IsCollapsed() bool {
	return bc.Collapsible && len(bc.Items) > bc.MaxVisible
}

// VisibleItems returns the items to display when collapsed.
func (bc *Breadcrumbs) VisibleItems() []*BreadcrumbItem {
	if !bc.IsCollapsed() {
		return bc.Items
	}
	// Show first item, ellipsis placeholder, and last items
	visible := make([]*BreadcrumbItem, 0)
	visible = append(visible, bc.Items[0])
	// Last MaxVisible-1 items
	start := len(bc.Items) - (bc.MaxVisible - 1)
	visible = append(visible, bc.Items[start:]...)
	return visible
}

// HiddenCount returns the number of hidden items.
func (bc *Breadcrumbs) HiddenCount() int {
	if !bc.IsCollapsed() {
		return 0
	}
	return len(bc.Items) - bc.MaxVisible
}

// SeparatorSymbol returns the separator character/symbol.
func (bc *Breadcrumbs) SeparatorSymbol() string {
	switch bc.Separator {
	case SeparatorSlash:
		return "/"
	case SeparatorArrow:
		return "→"
	case SeparatorDot:
		return "•"
	default:
		return "" // Chevron uses SVG
	}
}

// IsChevronSeparator returns true if using chevron separator.
func (bc *Breadcrumbs) IsChevronSeparator() bool {
	return bc.Separator == SeparatorChevron
}

// SizeClass returns CSS class for size.
func (bc *Breadcrumbs) SizeClass() string {
	switch bc.Size {
	case SizeSm:
		return "text-sm"
	case SizeLg:
		return "text-lg"
	default:
		return "text-base"
	}
}

// BreadcrumbItem represents a single breadcrumb item.
type BreadcrumbItem struct {
	base.Base

	// Label is the display text.
	Label string

	// Href is the link URL.
	Href string

	// Icon is an optional icon before the label.
	Icon string

	// Current marks this as the current page.
	Current bool

	// Disabled prevents interaction.
	Disabled bool
}

// ItemOption configures a BreadcrumbItem.
type ItemOption func(*BreadcrumbItem)

// NewItem creates a new BreadcrumbItem with the given ID and options.
func NewItem(id string, opts ...ItemOption) *BreadcrumbItem {
	item := &BreadcrumbItem{
		Base:    base.NewBase(id, "breadcrumb-item"),
		Current: false,
	}
	for _, opt := range opts {
		opt(item)
	}
	return item
}

// HasLabel returns true if label is set.
func (i *BreadcrumbItem) HasLabel() bool {
	return i.Label != ""
}

// HasHref returns true if href is set.
func (i *BreadcrumbItem) HasHref() bool {
	return i.Href != ""
}

// HasIcon returns true if icon is set.
func (i *BreadcrumbItem) HasIcon() bool {
	return i.Icon != ""
}

// IsClickable returns true if item can be clicked.
func (i *BreadcrumbItem) IsClickable() bool {
	return i.HasHref() && !i.Current && !i.Disabled
}

// LinkClass returns CSS class for link state.
func (i *BreadcrumbItem) LinkClass() string {
	if i.Current {
		return "text-gray-700 font-medium"
	}
	if i.Disabled {
		return "text-gray-400 cursor-not-allowed"
	}
	return "text-gray-500 hover:text-gray-700"
}
