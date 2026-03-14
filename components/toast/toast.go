// Package toast provides notification/toast components for displaying messages.
//
// Available variants:
//   - New() creates a toast container (template: "lvt:toast:container:v1")
//   - NewMessage() creates a single toast message (template: "lvt:toast:message:v1")
//
// Required lvt-* attributes: lvt-click
//
// Example usage:
//
//	// In your controller/state
//	Toasts: toast.New("notifications")
//
//	// To show a toast
//	state.Toasts.Add(toast.NewMessage(
//	    toast.WithTitle("Success"),
//	    toast.WithBody("Your changes have been saved."),
//	    toast.WithType(toast.Success),
//	))
//
//	// In your template
//	{{template "lvt:toast:container:v1" .Toasts}}
package toast

import (
	"strconv"

	"github.com/livetemplate/lvt/components/base"
	"github.com/livetemplate/lvt/components/styles"
)

// Type represents the visual style/severity of a toast.
type Type string

const (
	Info    Type = "info"
	Success Type = "success"
	Warning Type = "warning"
	Error   Type = "error"
)

// Position represents where toasts appear on screen.
type Position string

const (
	TopRight     Position = "top-right"
	TopLeft      Position = "top-left"
	TopCenter    Position = "top-center"
	BottomRight  Position = "bottom-right"
	BottomLeft   Position = "bottom-left"
	BottomCenter Position = "bottom-center"
)

// Message represents a single toast notification.
type Message struct {
	ID            string // Unique identifier for this toast
	Title         string // Optional title/header
	Body          string // Main message content
	Type          Type   // Visual style (info, success, warning, error)
	Dismissible   bool   // Whether user can dismiss the toast
	Icon          string // Optional icon (HTML or class name)
	AutoDismissMS int    // Auto-dismiss after this many milliseconds (0 = no auto-dismiss)
}

// Container holds and manages multiple toast notifications.
// Use template "lvt:toast:container:v1" to render.
type Container struct {
	base.Base

	// Messages is the list of active toasts
	Messages []Message

	// Position determines where toasts appear
	Position Position

	// MaxVisible limits how many toasts are shown at once (0 = unlimited)
	MaxVisible int

	// Counter for generating unique IDs
	Counter int `json:"counter"`
}

// New creates a toast container.
//
// Example:
//
//	toasts := toast.New("notifications",
//	    toast.WithPosition(toast.TopRight),
//	    toast.WithMaxVisible(5),
//	)
func New(id string, opts ...ContainerOption) *Container {
	c := &Container{
		Base:       base.NewBase(id, "toast"),
		Messages:   make([]Message, 0),
		Position:   TopRight,
		MaxVisible: 0,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Add adds a new toast message.
func (c *Container) Add(msg Message) {
	// Generate ID if not provided
	if msg.ID == "" {
		c.Counter++
		msg.ID = strconv.Itoa(c.Counter)
	}

	c.Messages = append(c.Messages, msg)

	// Trim to MaxVisible if set
	if c.MaxVisible > 0 && len(c.Messages) > c.MaxVisible {
		c.Messages = c.Messages[len(c.Messages)-c.MaxVisible:]
	}
}

// DefaultAutoDismissMS is the default auto-dismiss duration for success/info toasts.
const DefaultAutoDismissMS = 5000

// AddInfo adds an info toast that auto-dismisses after 5 seconds.
func (c *Container) AddInfo(title, body string) {
	c.Add(Message{Title: title, Body: body, Type: Info, Dismissible: true, AutoDismissMS: DefaultAutoDismissMS})
}

// AddSuccess adds a success toast that auto-dismisses after 5 seconds.
func (c *Container) AddSuccess(title, body string) {
	c.Add(Message{Title: title, Body: body, Type: Success, Dismissible: true, AutoDismissMS: DefaultAutoDismissMS})
}

// AddWarning adds a warning toast (no auto-dismiss).
func (c *Container) AddWarning(title, body string) {
	c.Add(Message{Title: title, Body: body, Type: Warning, Dismissible: true})
}

// AddError adds an error toast (no auto-dismiss).
func (c *Container) AddError(title, body string) {
	c.Add(Message{Title: title, Body: body, Type: Error, Dismissible: true})
}

// Dismiss removes a toast by ID.
func (c *Container) Dismiss(id string) {
	for i, msg := range c.Messages {
		if msg.ID == id {
			c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			return
		}
	}
}

// DismissAll removes all toasts.
func (c *Container) DismissAll() {
	c.Messages = make([]Message, 0)
}

// Count returns the number of active toasts.
func (c *Container) Count() int {
	return len(c.Messages)
}

// HasMessages returns true if there are any toasts.
func (c *Container) HasMessages() bool {
	return len(c.Messages) > 0
}

// VisibleMessages returns the messages to display (respects MaxVisible).
func (c *Container) VisibleMessages() []Message {
	if c.MaxVisible <= 0 || len(c.Messages) <= c.MaxVisible {
		return c.Messages
	}
	return c.Messages[len(c.Messages)-c.MaxVisible:]
}

// Styles returns the resolved ToastStyles for this component.
func (c *Container) Styles() styles.ToastStyles {
	if s, ok := c.StyleData().(styles.ToastStyles); ok {
		return s
	}
	adapter := styles.ForStyled(c.IsStyled())
	if adapter == nil {
		return styles.ToastStyles{}
	}
	s := adapter.ToastStyles()
	c.SetStyleData(s)
	return s
}

// GetPositionClasses returns CSS classes for the container position.
func (c *Container) GetPositionClasses() string {
	s := c.Styles()
	switch c.Position {
	case TopRight:
		return s.PosTopRight
	case TopLeft:
		return s.PosTopLeft
	case TopCenter:
		return s.PosTopCenter
	case BottomRight:
		return s.PosBottomRight
	case BottomLeft:
		return s.PosBottomLeft
	case BottomCenter:
		return s.PosBottomCenter
	default:
		return s.PosTopRight
	}
}

// GetTypeClasses returns CSS classes for a toast type using the style adapter.
// Accepts Type or string to handle values that may lose their named type
// after JSON round-trip (e.g., state cloning via marshal/unmarshal in
// mount.go:cloneStateTyped). The alternative — normalizing at deserialization
// — would require framework-level changes to the generic state cloner.
func (c *Container) GetTypeClasses(t interface{}) string {
	var typ Type
	switch v := t.(type) {
	case Type:
		typ = v
	case string:
		typ = Type(v)
	default:
		typ = Info
	}

	s := c.Styles()
	switch typ {
	case Success:
		return s.TypeSuccess
	case Warning:
		return s.TypeWarning
	case Error:
		return s.TypeError
	default: // Info
		return s.TypeInfo
	}
}

// GetTypeIcon returns a default icon for the toast type.
func GetTypeIcon(t Type) string {
	switch t {
	case Success:
		return `<svg class="w-5 h-5 text-green-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>`
	case Warning:
		return `<svg class="w-5 h-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg>`
	case Error:
		return `<svg class="w-5 h-5 text-red-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/></svg>`
	default: // Info
		return `<svg class="w-5 h-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/></svg>`
	}
}

