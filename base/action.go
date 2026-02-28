package base

import (
	"strconv"
)

// ActionContext provides context for component action handlers.
// It wraps the data passed from lvt-data-* attributes and provides
// convenient accessor methods.
//
// Components receive ActionContext in their action methods:
//
//	func (d *Dropdown) Select(ctx *ActionContext) error {
//	    value := ctx.Data("value")
//	    d.Selected = d.findOption(value)
//	    return nil
//	}
type ActionContext struct {
	// Action is the action name that was triggered (e.g., "select", "toggle").
	Action string

	// ComponentID is the ID of the component that triggered the action.
	ComponentID string

	// data holds the key-value pairs from lvt-data-* attributes.
	data map[string]string
}

// NewActionContext creates a new ActionContext with the given action name,
// component ID, and data map.
func NewActionContext(action, componentID string, data map[string]string) *ActionContext {
	if data == nil {
		data = make(map[string]string)
	}
	return &ActionContext{
		Action:      action,
		ComponentID: componentID,
		data:        data,
	}
}

// Data returns the value for the given key, or empty string if not present.
// Keys correspond to lvt-data-* attributes without the "lvt-data-" prefix.
//
// Example:
//
//	<button lvt-click="select_myid" lvt-data-value="option1">Click</button>
//
//	// In action handler:
//	value := ctx.Data("value") // returns "option1"
func (ctx *ActionContext) Data(key string) string {
	return ctx.data[key]
}

// DataInt returns the value for the given key as an integer.
// Returns 0 if the key is not present or cannot be parsed as int.
func (ctx *ActionContext) DataInt(key string) int {
	if v, ok := ctx.data[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// DataFloat returns the value for the given key as a float64.
// Returns 0.0 if the key is not present or cannot be parsed as float.
func (ctx *ActionContext) DataFloat(key string) float64 {
	if v, ok := ctx.data[key]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// DataBool returns the value for the given key as a boolean.
// Returns false if the key is not present.
// Values "true", "1", "yes", "on" are considered true (case-insensitive).
func (ctx *ActionContext) DataBool(key string) bool {
	if v, ok := ctx.data[key]; ok {
		switch v {
		case "true", "1", "yes", "on", "True", "TRUE", "Yes", "YES", "On", "ON":
			return true
		}
	}
	return false
}

// HasData returns true if the given key exists in the data map.
func (ctx *ActionContext) HasData(key string) bool {
	_, ok := ctx.data[key]
	return ok
}

// AllData returns a copy of all data key-value pairs.
func (ctx *ActionContext) AllData() map[string]string {
	result := make(map[string]string, len(ctx.data))
	for k, v := range ctx.data {
		result[k] = v
	}
	return result
}

// ActionHandler is the function signature for component action handlers.
// Components implement methods matching this signature to handle user interactions.
//
// Example:
//
//	func (d *Dropdown) Toggle(ctx *ActionContext) error {
//	    d.Open = !d.Open
//	    return nil
//	}
type ActionHandler func(ctx *ActionContext) error

// ActionProvider is an interface for components that provide action handlers.
// Components implement this to expose their handlers to the LiveTemplate framework.
type ActionProvider interface {
	// Actions returns a map of action names to handlers.
	// Action names should be without the component ID suffix.
	// The framework will automatically match "toggle_myid" to the "toggle" handler.
	Actions() map[string]ActionHandler
}
