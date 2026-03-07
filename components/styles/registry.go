package styles

import (
	"sync"
)

var (
	mu       sync.RWMutex
	adapters = map[string]StyleAdapter{}
	current  StyleAdapter
)

// Register makes a StyleAdapter available by name.
// Typically called from adapter init() functions.
func Register(a StyleAdapter) {
	mu.Lock()
	defer mu.Unlock()
	adapters[a.Name()] = a
	// First registered adapter becomes the default
	if current == nil {
		current = a
	}
}

// Get returns a registered adapter by name, or nil if not found.
func Get(name string) StyleAdapter {
	mu.RLock()
	defer mu.RUnlock()
	return adapters[name]
}

// Default returns the current default adapter.
// Returns nil if no adapters have been registered.
func Default() StyleAdapter {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

// SetDefault sets the default adapter used by all components.
func SetDefault(a StyleAdapter) {
	mu.Lock()
	defer mu.Unlock()
	current = a
}

// ForStyled returns the appropriate adapter based on the styled flag.
// When styled is true, returns the default adapter (typically tailwind).
// When styled is false, returns the "unstyled" adapter if registered,
// otherwise returns the default adapter.
func ForStyled(styled bool) StyleAdapter {
	mu.RLock()
	defer mu.RUnlock()
	if !styled {
		if unstyled, ok := adapters["unstyled"]; ok {
			return unstyled
		}
	}
	return current
}

// Names returns all registered adapter names.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(adapters))
	for name := range adapters {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered adapters.
func Count() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(adapters)
}
