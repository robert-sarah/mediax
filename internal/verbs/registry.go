package verbs

import (
	"fmt"
)

// Verb represents a mediax command verb
type Verb interface {
	Name() string
	Description() string
	Usage() string
	Execute(args []string, flags map[string]string) error
}

// Registry holds all registered verbs
var registry = make(map[string]Verb)

// Register adds a verb to the registry
func Register(v Verb) {
	registry[v.Name()] = v
}

// Get retrieves a verb by name
func Get(name string) (Verb, bool) {
	v, ok := registry[name]
	return v, ok
}

// GetAll returns all registered verbs
func GetAll() map[string]Verb {
	return registry
}

// Execute runs a verb with the given arguments
func Execute(name string, args []string, flags map[string]string) error {
	v, ok := Get(name)
	if !ok {
		return fmt.Errorf("unknown verb: %s. Run 'mediax verbs' to see all available verbs", name)
	}
	return v.Execute(args, flags)
}
