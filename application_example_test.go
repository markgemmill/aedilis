package aedilis

import (
	"fmt"
)

// ComponentA
type ComponentA struct {
	Name string
}

func (c *ComponentA) DoIt() {
	fmt.Println("ComponentA does it!")
}

func InitA(app *Application) (Component, ComponentFunc, ComponentFunc, error) {
	fmt.Println("Initializing A!")
	c := &ComponentA{"A, eh!"}
	return c, StartA(c), StopA(c), nil
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

func InitB(app *Application) (Component, ComponentFunc, ComponentFunc, error) {
	fmt.Println("Initializing B!")
	a, err := GetComponent[*ComponentA](app, "*aedilis.ComponentA")
	if err != nil {
		return nil, nil, nil, err
	}
	b := &ComponentB{Name: "B, eh!", Other: a}
	fmt.Printf("  - %T uses %T\n", b, a)
	return b, StartB(b), StopB(b), nil
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
