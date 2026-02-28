// Package base provides common types and utilities for LiveTemplate components.
//
// All components embed the Base struct which provides ID handling and common functionality.
// Components use functional options for configuration and implement action handlers
// that are automatically registered with the LiveTemplate framework.
package base

// Base provides common functionality for all components.
// Components should embed this struct to gain ID handling and common utilities.
//
// IMPORTANT: Base uses exported fields with JSON tags for proper serialization.
// Do NOT add custom MarshalJSON/UnmarshalJSON to Base - it will break parent
// struct serialization because Go promotes embedded methods to the parent,
// making the parent's json.Marshal only serialize Base fields.
//
// Example:
//
//	type Dropdown struct {
//	    base.Base
//	    Options []Option
//	    Selected *Option
//	}
type Base struct {
	// ComponentID is the unique identifier for this component instance.
	// Used in action names like "toggle_{ID}", "select_{ID}".
	// Exported for JSON serialization. Use ID() method to access.
	ComponentID string `json:"id"`

	// ComponentNamespace is the component type, e.g., "dropdown", "tabs".
	// Used for action naming and template resolution.
	// Exported for JSON serialization. Use Namespace() method to access.
	ComponentNamespace string `json:"namespace"`

	// Styled indicates whether to use Tailwind CSS classes (true) or semantic HTML only (false).
	Styled bool `json:"styled"`
}

// NewBase creates a new Base with the given ID and namespace.
//
// Example:
//
//	func New(id string, opts ...Option) *Dropdown {
//	    d := &Dropdown{
//	        Base: base.NewBase(id, "dropdown"),
//	    }
//	    for _, opt := range opts {
//	        opt(d)
//	    }
//	    return d
//	}
func NewBase(id, namespace string) Base {
	return Base{
		ComponentID:        id,
		ComponentNamespace: namespace,
		Styled:             true, // Default to styled (Tailwind CSS)
	}
}

// ID returns the component's unique identifier.
func (b *Base) ID() string {
	return b.ComponentID
}

// Namespace returns the component's type namespace.
func (b *Base) Namespace() string {
	return b.ComponentNamespace
}

// ActionName generates a namespaced action name for this component.
// For example: ActionName("toggle") returns "toggle_myid" if ID is "myid".
func (b *Base) ActionName(action string) string {
	return action + "_" + b.ComponentID
}

// SetStyled sets whether to use Tailwind CSS classes.
func (b *Base) SetStyled(styled bool) {
	b.Styled = styled
}

// IsStyled returns true if the component should use Tailwind CSS classes.
func (b *Base) IsStyled() bool {
	return b.Styled
}
