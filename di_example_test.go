package di

import (
	"fmt"
)

func ExampleProvide() {
	container := New()

	type test struct {
		foobar string
	}

	Provide(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value, err := Invoke[*test](container)

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleInvoke() {
	container := New()

	type test struct {
		foobar string
	}

	Provide(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value, err := Invoke[*test](container)

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleMustInvoke() {
	container := New()

	type test struct {
		foobar string
	}

	Provide(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value := MustInvoke[*test](container)

	fmt.Println(value)
	// Output:
	// &{foobar}
}

func ExampleProvideNamed() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "my_service", func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value, err := InvokeNamed[*test](container, "my_service")

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleInvokeNamed() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "my_service", func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value, err := InvokeNamed[*test](container, "my_service")

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleMustInvokeNamed() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "my_service", func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})
	value := MustInvokeNamed[*test](container, "my_service")

	fmt.Println(value)
	// Output:
	// &{foobar}
}

func ExampleProvideValue() {
	Container := New()

	type test struct {
		foobar string
	}

	ProvideValue(Container, &test{foobar: "foobar"})
	value, err := Invoke[*test](Container)

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleProvideNamedValue() {
	Container := New()

	type test struct {
		foobar string
	}

	ProvideNamedValue(Container, "my_service", &test{foobar: "foobar"})
	value, err := InvokeNamed[*test](Container, "my_service")

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar}
	// <nil>
}

func ExampleOverride() {
	container := New()

	type test struct {
		foobar string
	}

	Provide(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar1"}, nil
	})
	Override(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar2"}, nil
	})
	value, err := Invoke[*test](container)

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar2}
	// <nil>
}

func ExampleOverrideNamed() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "my_service", func(i *Container) (*test, error) {
		return &test{foobar: "foobar1"}, nil
	})
	OverrideNamed(container, "my_service", func(i *Container) (*test, error) {
		return &test{foobar: "foobar2"}, nil
	})
	value, err := InvokeNamed[*test](container, "my_service")

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar2}
	// <nil>
}

func ExampleOverrideNamedValue() {
	Container := New()

	type test struct {
		foobar string
	}

	ProvideNamedValue(Container, "my_service", &test{foobar: "foobar1"})
	OverrideNamedValue(Container, "my_service", &test{foobar: "foobar2"})
	value, err := InvokeNamed[*test](Container, "my_service")

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// &{foobar2}
	// <nil>
}
