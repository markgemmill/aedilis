package aedilis

// ActionRegistry manages a list of functions that perform specific application
// functions. Example component starter or closing functions.
type ActionRegistry struct {
	name    string
	actions map[string]ComponentFunc
	order   []string
}

func NewActionRegistry(name string) *ActionRegistry {
	reg := ActionRegistry{name: name}
	reg.actions = make(map[string]ComponentFunc)
	return &reg
}

// Add a new component function to the registry.
func (reg *ActionRegistry) Add(name string, f ComponentFunc) {
	reg.actions[name] = f
	reg.order = append(reg.order, name)
	//reg.actions = append(reg.actions, f)
}

// Count return the number of registered component functions.
func (reg *ActionRegistry) Count() int {
	return len(reg.order)
}

// Execute all the functions in the registry in the order registered.
func (reg *ActionRegistry) Execute(app *Application, haltOnErr bool) error {
	for _, name := range reg.order {
		app.console.Write("Executing %s function %s\n", reg.name, name)
		err := reg.actions[name](app)
		if haltOnErr == true && err != nil {
			return err
		}
	}
	return nil
}

// ExecuteReverse the functions in the registry in the
// reverse order they were registered.
func (reg *ActionRegistry) ExecuteReverse(app *Application, haltOnErr bool) error {
	for i := len(reg.order) - 1; i >= 0; i-- {
		name := reg.order[i]
		app.console.Write("Executing %s function %s\n", reg.name, name)
		err := reg.actions[name](app)
		if haltOnErr == true && err != nil {
			return err
		}
	}
	return nil
}
