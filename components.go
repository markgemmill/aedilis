package aedilis

import "fmt"

func generateComponentName(c any) string {
	name := fmt.Sprintf("%T", c)
	if name != "<nil>" {
		return name
	}
	return "noname"
}

// Component can be any type, but the general idea
// is it represents a piece of application functionality
// such as logger, configuration manager, global error handler,
// service managers, etc.
type Component interface{}

// ComponentFunc is a function type returned during component
// initialization that represents an action that can be called
// by the application. An action could start or shutdown the component.
type ComponentFunc func(app *Application) error

type StartFunc func(app *Application) error

type CloseFunc func(app *Application) error

type ComponentOptions struct {
	Component Component
	Starter   ComponentFunc
	Closer    ComponentFunc
	Alias     string
}

func ComponentError(err error) (ComponentOptions, error) {
	return ComponentOptions{}, err
}

func (co *ComponentOptions) Name() string {
	if co.Alias == "" {
		return generateComponentName(co.Component)
	}
	return co.Alias
}

// InitFunc is a function type specific to initializing a component. It
// must return the initialized component object and a start and shutdown function.
// These could all be nil. There is no specific requirement that they exist.
type InitFunc func(app *Application) (ComponentOptions, error)

// ComponentRegistry maintains a named list of application components.
type ComponentRegistry struct {
	components map[string]Component
	order      []string
}

// NewComponentRegistry creates a new registry object.
func NewComponentRegistry() *ComponentRegistry {
	mgr := &ComponentRegistry{}
	mgr.components = make(map[string]Component)
	return mgr
}

// Get returns the named component if it exists, otherwise nil
func (mgr *ComponentRegistry) Get(name string) (Component, bool) {
	cmp, ok := mgr.components[name]
	if ok {
		return cmp, true
	}
	return nil, false
}

// Add adds the named component.
func (mgr *ComponentRegistry) Add(name string, component Component) error {
	_, exists := mgr.Get(name)
	if exists {
		return fmt.Errorf("DefaultComponentManager already contains '%s'", name)
	}
	mgr.components[name] = component
	mgr.order = append(mgr.order, name)
	return nil
}

// Count returns the number of components registered.
func (mgr *ComponentRegistry) Count() int {
	return len(mgr.order)
}
