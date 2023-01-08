package aedilis

import "fmt"

// ComponentA
type ComponentA struct {
	Name string
}

func (c *ComponentA) DoIt() {
	fmt.Println("ComponentA does it!")
}

func InitA(app *Application) (ComponentOptions, error) {
	fmt.Println("Initializing A!")
	c := &ComponentA{"A, eh!"}
	return ComponentOptions{
		Component: c,
		Starter:   StartA(c),
		Closer:    StopA(c),
		Alias:     "",
	}, nil
}

func StartA(cmp *ComponentA) ComponentFunc {
	return func(app *Application) error {
		fmt.Printf("Starting %T...\n", cmp)
		return nil
	}
}

func StopA(cmp *ComponentA) ComponentFunc {
	return func(app *Application) error {
		fmt.Printf("Stopping %T...\n", cmp)
		return nil
	}
}

// ComponentB
type ComponentB struct {
	Name  string
	Other *ComponentA
}

func (c *ComponentB) DoIt() {
	fmt.Println("ComponentB does it!")
}

func InitB(app *Application) (ComponentOptions, error) {
	fmt.Println("Initializing B!")
	a, err := GetComponent[*ComponentA](app, "*aedilis.ComponentA")
	if err != nil {
		return ComponentError(err)
	}
	b := &ComponentB{Name: "B, eh!", Other: a}
	fmt.Printf("  - %T uses %T\n", b, a)
	return ComponentOptions{Component: b, Starter: StartB(b), Closer: StopB(b)}, nil
}

func StartB(c *ComponentB) ComponentFunc {
	return func(app *Application) error {
		fmt.Printf("Starting %T...\n", c)
		return nil
	}
}

func StopB(c *ComponentB) ComponentFunc {
	return func(app *Application) error {
		fmt.Printf("Stopping %T...\n", c)
		return nil
	}
}

func ExampleRegistry() {

	app := New()
	_ = app.Run(InitA, InitB)

	// output:
	// Initializing A!
	// Registering component *aedilis.ComponentA
	// Initializing B!
	//   - *aedilis.ComponentB uses *aedilis.ComponentA
	// Registering component *aedilis.ComponentB
	// Starting *aedilis.ComponentA...
	// Starting *aedilis.ComponentB...
	// Stopping *aedilis.ComponentB...
	// Stopping *aedilis.ComponentA...

}
