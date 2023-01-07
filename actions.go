package aedilis

// ActionRegistry manages a list of functions that perform specific application
// functions. Example component starter or closing functions.
type ActionRegistry struct {
	actions []ComponentFunc
}

func NewActionRegistry() *ActionRegistry {
	reg := ActionRegistry{}
	return &reg
}

// Add a new component function to the registry.
func (reg *ActionRegistry) Add(f ComponentFunc) {
	reg.actions = append(reg.actions, f)
}

// Count return the number of registered component functions.
func (reg *ActionRegistry) Count() int {
	return len(reg.actions)
}

// Execute all the functions in the registry in the order registered.
func (reg *ActionRegistry) Execute(app *Application, haltOnErr bool) error {
	for _, f := range reg.actions {
		err := f(app)
		if haltOnErr == true && err != nil {
			return err
		}
	}
	return nil
}

// ExecuteReverse the functions in the registry in the
// reverse order they were registered.
func (reg *ActionRegistry) ExecuteReverse(app *Application, haltOnErr bool) error {
	for i := len(reg.actions) - 1; i >= 0; i-- {
		err := reg.actions[i](app)
		if haltOnErr == true && err != nil {
			return err
		}
	}
	return nil
}
