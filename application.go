package aedilis

import "fmt"

type Application struct {
	components *ComponentRegistry
	starters   *ActionRegistry
	closers    *ActionRegistry
}

func New() *Application {
	app := &Application{}
	app.components = NewComponentRegistry()
	app.starters = NewActionRegistry()
	app.closers = NewActionRegistry()
	return app
}

func (app *Application) Run(initializers ...InitFunc) error {
	err := Init(app, initializers...)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		app.Stop()
		return err
	}

	err = app.Start()
	app.Stop()

	return err
}

// registerComponent adds the named component to the registry.
func (app *Application) registerComponent(name string, component Component) error {
	fmt.Printf("Registering component %s\n", name)
	return app.components.Add(name, component)
}

// Init takes a component initialization function and calls it immediately.
func (app *Application) Init(initializer InitFunc) error {

	component, starter, closer, err := initializer(app)

	if err != nil {
		return err
	}

	if component != nil {
		name := generateComponentName(component)
		err = app.registerComponent(name, component)
		if err != nil {
			return err
		}
	}

	if starter != nil {
		app.starters.Add(starter)
	}

	if closer != nil {
		app.closers.Add(closer)
	}

	return nil
}

// Start run all starter functions in order they were registered.
func (app *Application) Start() error {
	return app.starters.Execute(app, true)
}

// Stop runs all closer functions in reverse order registered.
func (app *Application) Stop() {
	_ = app.closers.ExecuteReverse(app, false)
}

// GetComponent fetches the named component object from the registry
// as the given generic type.
func GetComponent[T Component](app *Application, name string) (T, error) {
	c, exists := app.components.Get(name)
	if !exists {
		var result T
		return result, fmt.Errorf("component '%s' does not exist", name)
	}
	x, ok := c.(T)
	if ok {
		return x, nil
	}
	return x, fmt.Errorf("component '%s' is not of requested type", name)
}

// MustGetComponent fetches the named component object from the registry. If the component
// does not exist panic is raised.
func MustGetComponent[T Component](app *Application, name string) (T, error) {
	c, err := GetComponent[T](app, name)
	if err != nil {
		panic(err)
	}
	return c, nil
}

func Init(app *Application, initializers ...InitFunc) error {
	for _, initializer := range initializers {
		err := app.Init(initializer)
		if err != nil {
			return err
		}
	}
	return nil
}
