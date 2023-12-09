package main

import (
	"log"

	"github.com/cryptoniumX/di"
)

/**
 * Wheel
 */
type Wheel struct {
}

/**
 * Engine
 */
type Engine struct {
}

func (c *Engine) Shutdown() error {
	println("engine stopped")
	return nil
}

/**
 * Car
 */
type Car struct {
	Engine *Engine
	Wheels []*Wheel
}

func (c *Car) Shutdown() error {
	println("car stopped")
	return nil
}

func (c *Car) Start() {
	println("vroooom")
}

/**
 * Run example
 */
func main() {
	Container := di.New()

	// provide wheels
	di.ProvideNamedValue(Container, "wheel-1", &Wheel{})
	di.ProvideNamedValue(Container, "wheel-2", &Wheel{})
	di.ProvideNamedValue(Container, "wheel-3", &Wheel{})
	di.ProvideNamedValue(Container, "wheel-4", &Wheel{})

	// provide car
	di.Provide(Container, func(i *di.Container) (*Car, error) {
		car := Car{
			Engine: di.MustInvoke[*Engine](i),
			Wheels: []*Wheel{
				di.MustInvokeNamed[*Wheel](i, "wheel-1"),
				di.MustInvokeNamed[*Wheel](i, "wheel-2"),
				di.MustInvokeNamed[*Wheel](i, "wheel-3"),
				di.MustInvokeNamed[*Wheel](i, "wheel-4"),
			},
		}

		return &car, nil
	})

	// provide engine
	di.Provide(Container, func(i *di.Container) (*Engine, error) {
		return &Engine{}, nil
	})

	// start car
	car := di.MustInvoke[*Car](Container)
	car.Start()

	err := Container.ShutdownOnSIGTERM()
	if err != nil {
		log.Fatal(err.Error())
	}
}
