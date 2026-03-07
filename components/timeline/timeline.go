// Package timeline provides timeline components for the LiveTemplate framework.
//
// Timelines display a list of events in chronological order with optional
// icons, status indicators, and content.
//
// # Available Templates
//
//   - lvt:timeline:default:v1 - Standard vertical timeline
//
// # Basic Usage
//
//	tl := timeline.New("history",
//	    timeline.WithItems(
//	        timeline.NewItem("1", timeline.WithItemTitle("Event 1")),
//	        timeline.NewItem("2", timeline.WithItemTitle("Event 2")),
//	    ),
//	)
//
//	{{template "lvt:timeline:default:v1" .Timeline}}
package timeline

import (
	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
)

// Orientation defines timeline layout direction.
type Orientation string

const (
	// OrientationVertical displays items top to bottom.
	OrientationVertical Orientation = "vertical"
	// OrientationHorizontal displays items left to right.
	OrientationHorizontal Orientation = "horizontal"
)

// Position defines item content placement.
type Position string

const (
	// PositionLeft places content on the left side.
	PositionLeft Position = "left"
	// PositionRight places content on the right side.
	PositionRight Position = "right"
	// PositionAlternate alternates content between left and right.
	PositionAlternate Position = "alternate"
)

// Status defines item status indicators.
type Status string

const (
	// StatusDefault is the default/neutral status.
	StatusDefault Status = "default"
	// StatusPending indicates a pending/waiting state.
	StatusPending Status = "pending"
	// StatusActive indicates the current/active state.
	StatusActive Status = "active"
	// StatusComplete indicates completion.
	StatusComplete Status = "complete"
	// StatusError indicates an error state.
	StatusError Status = "error"
)

// Color defines item indicator colors.
type Color string

const (
	ColorGray   Color = "gray"
	ColorBlue   Color = "blue"
	ColorGreen  Color = "green"
	ColorYellow Color = "yellow"
	ColorRed    Color = "red"
	ColorPurple Color = "purple"
)

// Timeline represents a timeline container.
type Timeline struct {
	base.Base

	// Items in the timeline.
	Items []*TimelineItem

	// Orientation controls layout direction.
	Orientation Orientation

	// Position controls content placement (for vertical orientation).
	Position Position

	// ShowConnectors displays lines between items.
	ShowConnectors bool

	// Reverse displays items in reverse order.
	Reverse bool
}

// Option configures a Timeline.
type Option func(*Timeline)

// New creates a new Timeline with the given ID and options.
func New(id string, opts ...Option) *Timeline {
	t := &Timeline{
		Base:           base.NewBase(id, "timeline"),
		Items:          make([]*TimelineItem, 0),
		Orientation:    OrientationVertical,
		Position:       PositionLeft,
		ShowConnectors: true,
		Reverse:        false,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// AddItem adds an item to the timeline.
func (t *Timeline) AddItem(item *TimelineItem) {
	t.Items = append(t.Items, item)
}

// RemoveItem removes an item by ID.
func (t *Timeline) RemoveItem(id string) {
	items := make([]*TimelineItem, 0, len(t.Items))
	for _, item := range t.Items {
		if item.ID() != id {
			items = append(items, item)
		}
	}
	t.Items = items
}

// GetItem returns an item by ID or nil if not found.
func (t *Timeline) GetItem(id string) *TimelineItem {
	for _, item := range t.Items {
		if item.ID() == id {
			return item
		}
	}
	return nil
}

// HasItems returns true if timeline has items.
func (t *Timeline) HasItems() bool {
	return len(t.Items) > 0
}

// ItemCount returns the number of items.
func (t *Timeline) ItemCount() int {
	return len(t.Items)
}

// IsVertical returns true if orientation is vertical.
func (t *Timeline) IsVertical() bool {
	return t.Orientation == OrientationVertical
}

// IsHorizontal returns true if orientation is horizontal.
func (t *Timeline) IsHorizontal() bool {
	return t.Orientation == OrientationHorizontal
}

// IsAlternate returns true if position is alternate.
func (t *Timeline) IsAlternate() bool {
	return t.Position == PositionAlternate
}

// Styles returns the resolved TimelineStyles for this component.
// It uses lazy resolution: the result is cached after the first call.
func (t *Timeline) Styles() styles.TimelineStyles {
	if s, ok := t.StyleData().(styles.TimelineStyles); ok {
		return s
	}
	adapter := styles.ForStyled(t.IsStyled())
	if adapter == nil {
		return styles.TimelineStyles{}
	}
	s := adapter.TimelineStyles()
	t.SetStyleData(s)
	return s
}

// OrientationClass returns CSS class for orientation.
func (t *Timeline) OrientationClass() string {
	st := t.Styles()
	if t.IsHorizontal() {
		return st.HorizontalRoot
	}
	return st.VerticalRoot
}

// TimelineItem represents a single timeline entry.
type TimelineItem struct {
	base.Base

	// Title is the item heading.
	Title string

	// Description is additional content.
	Description string

	// Time displays a timestamp or date.
	Time string

	// Icon is the icon name or SVG.
	Icon string

	// Status indicates the item state.
	Status Status

	// Color is the indicator color.
	Color Color

	// Active highlights this item.
	Active bool

	// Completed marks item as done.
	Completed bool
}

// ItemOption configures a TimelineItem.
type ItemOption func(*TimelineItem)

// NewItem creates a new TimelineItem with the given ID and options.
func NewItem(id string, opts ...ItemOption) *TimelineItem {
	item := &TimelineItem{
		Base:   base.NewBase(id, "timeline-item"),
		Status: StatusDefault,
		Color:  ColorGray,
	}
	for _, opt := range opts {
		opt(item)
	}
	return item
}

// HasTitle returns true if title is set.
func (i *TimelineItem) HasTitle() bool {
	return i.Title != ""
}

// HasDescription returns true if description is set.
func (i *TimelineItem) HasDescription() bool {
	return i.Description != ""
}

// HasTime returns true if time is set.
func (i *TimelineItem) HasTime() bool {
	return i.Time != ""
}

// HasIcon returns true if icon is set.
func (i *TimelineItem) HasIcon() bool {
	return i.Icon != ""
}

// IsPending returns true if status is pending.
func (i *TimelineItem) IsPending() bool {
	return i.Status == StatusPending
}

// IsActive returns true if status is active or Active is true.
func (i *TimelineItem) IsActive() bool {
	return i.Status == StatusActive || i.Active
}

// IsComplete returns true if status is complete or Completed is true.
func (i *TimelineItem) IsComplete() bool {
	return i.Status == StatusComplete || i.Completed
}

// IsError returns true if status is error.
func (i *TimelineItem) IsError() bool {
	return i.Status == StatusError
}

// Styles returns the resolved TimelineItemStyles for this component.
// It uses lazy resolution: the result is cached after the first call.
func (i *TimelineItem) Styles() styles.TimelineItemStyles {
	if s, ok := i.StyleData().(styles.TimelineItemStyles); ok {
		return s
	}
	adapter := styles.ForStyled(i.IsStyled())
	if adapter == nil {
		return styles.TimelineItemStyles{}
	}
	s := adapter.TimelineItemStyles()
	i.SetStyleData(s)
	return s
}

// IndicatorClass returns CSS class for the indicator dot.
func (i *TimelineItem) IndicatorClass() string {
	st := i.Styles()
	switch i.Color {
	case ColorBlue:
		return st.ColorBlue
	case ColorGreen:
		return st.ColorGreen
	case ColorYellow:
		return st.ColorYellow
	case ColorRed:
		return st.ColorRed
	case ColorPurple:
		return st.ColorPurple
	default:
		return st.ColorGray
	}
}

// StatusClass returns CSS class based on status.
func (i *TimelineItem) StatusClass() string {
	st := i.Styles()
	switch i.Status {
	case StatusPending:
		return st.StatusPending
	case StatusActive:
		return st.StatusActive
	case StatusComplete:
		return st.StatusComplete
	case StatusError:
		return st.StatusError
	default:
		return st.StatusDefault
	}
}

// RingClass returns ring CSS class for active items.
func (i *TimelineItem) RingClass() string {
	if i.IsActive() {
		st := i.Styles()
		switch i.Color {
		case ColorBlue:
			return st.RingBlue
		case ColorGreen:
			return st.RingGreen
		case ColorYellow:
			return st.RingYellow
		case ColorRed:
			return st.RingRed
		case ColorPurple:
			return st.RingPurple
		default:
			return st.RingGray
		}
	}
	return ""
}
