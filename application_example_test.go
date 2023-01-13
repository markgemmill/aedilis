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

	app := New("")
	_ = app.Run(InitA, InitB)

	// output:
	// Initializing A!
	// [aedilis] Registering component *aedilis.ComponentA
	// [aedilis] Registering start function *aedilis.ComponentA
	// [aedilis] Registering close function *aedilis.ComponentA
	// Initializing B!
	//   - *aedilis.ComponentB uses *aedilis.ComponentA
	// [aedilis] Registering component *aedilis.ComponentB
	// [aedilis] Registering start function *aedilis.ComponentB
	// [aedilis] Registering close function *aedilis.ComponentB
	// [aedilis] Executing starter function *aedilis.ComponentA
	// Starting *aedilis.ComponentA...
	// [aedilis] Executing starter function *aedilis.ComponentB
	// Starting *aedilis.ComponentB...
	// [aedilis] Executing shutdown function *aedilis.ComponentB
	// Stopping *aedilis.ComponentB...
	// [aedilis] Executing shutdown function *aedilis.ComponentA
	// Stopping *aedilis.ComponentA...

}

func ExampleRegistryWithName() {

	app := New("testing")
	_ = app.Run(InitA, InitB)

	// output:
	// Initializing A!
	// [testing] Registering component *aedilis.ComponentA
	// [testing] Registering start function *aedilis.ComponentA
	// [testing] Registering close function *aedilis.ComponentA
	// Initializing B!
	//   - *aedilis.ComponentB uses *aedilis.ComponentA
	// [testing] Registering component *aedilis.ComponentB
	// [testing] Registering start function *aedilis.ComponentB
	// [testing] Registering close function *aedilis.ComponentB
	// [testing] Executing starter function *aedilis.ComponentA
	// Starting *aedilis.ComponentA...
	// [testing] Executing starter function *aedilis.ComponentB
	// Starting *aedilis.ComponentB...
	// [testing] Executing shutdown function *aedilis.ComponentB
	// Stopping *aedilis.ComponentB...
	// [testing] Executing shutdown function *aedilis.ComponentA
	// Stopping *aedilis.ComponentA...

}
