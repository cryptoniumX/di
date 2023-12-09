package main

import (
	"github.com/cryptoniumX/di"
)

func main() {
	Container := di.New()

	// provide wheels
	di.ProvideNamedValue(Container, "wheel-1", NewWheel())
	di.ProvideNamedValue(Container, "wheel-2", NewWheel())
	di.ProvideNamedValue(Container, "wheel-3", NewWheel())
	di.ProvideNamedValue(Container, "wheel-4", NewWheel())

	// provide car
	di.Provide(Container, NewCar)

	// provide engine
	di.Provide(Container, NewEngine)

	// start car
	car := di.MustInvoke[Car](Container)
	car.Start()
}
