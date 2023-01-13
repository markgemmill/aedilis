package aedilis

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	Console    *Console
	components *ComponentRegistry
	starters   *ActionRegistry
	closers    *ActionRegistry
}

func New() *Application {
	app := &Application{}
	app.components = NewComponentRegistry()
	app.starters = NewActionRegistry("starter")
	app.closers = NewActionRegistry("shutdown")
	app.Console = NewConsole()
	return app
}

func (app *Application) StartWithInterruptWrapper(startFunc ComponentFunc) ComponentFunc {
	return func(app *Application) error {
		var err error
		go func() {
			err = startFunc(app)
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
		<-quit
		return err
	}
}

func (app *Application) Run(initializers ...InitFunc) error {
	err := Init(app, initializers...)
	if err != nil {
		app.Console.WriteError("%s", err.Error())
		app.Stop()
		return err
	}

	err = app.Start()
	app.Stop()

	return err
}

// RegisterComponent adds the named component to the registry.
func (app *Application) RegisterComponent(name string, component Component) error {
	app.Console.Write("Registering component %s", name)
	return app.components.Add(name, component)
}

func (app *Application) RegisterStarter(name string, starter ComponentFunc) {
	app.Console.Write("Registering start function %s", name)
	app.starters.Add(name, starter)
}

func (app *Application) RegisterCloser(name string, closer ComponentFunc) {
	app.Console.Write("Registering close function %s", name)
	app.closers.Add(name, closer)
}

// Init takes a component initialization function and calls it immediately.
func (app *Application) Init(initializer InitFunc) error {

	options, err := initializer(app)

	if err != nil {
		return err
	}

	name := options.Name()

	if options.Component != nil {
		err = app.RegisterComponent(name, options.Component)
		if err != nil {
			return err
		}
	}

	if options.Starter != nil {
		app.RegisterStarter(name, options.Starter)
	}

	if options.Closer != nil {
		app.RegisterCloser(name, options.Closer)
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
func MustGetComponent[T Component](app *Application, name string) T {
	c, err := GetComponent[T](app, name)
	if err != nil {
		panic(err)
	}
	return c
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
