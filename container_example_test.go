package di

import (
	"database/sql"
	"fmt"
)

type dbService struct {
	db *sql.DB
}

func (s *dbService) HealthCheck() error {
	return nil
}

func (s *dbService) Shutdown() error {
	return nil
}

func dbServiceProvider(i *Container) (*dbService, error) {
	return &dbService{db: nil}, nil
}

func ExampleNew() {
	Container := New()

	ProvideNamedValue(Container, "PG_URI", "postgres://user:pass@host:5432/db")
	uri, err := InvokeNamed[string](Container, "PG_URI")

	fmt.Println(uri)
	fmt.Println(err)
	// Output:
	// postgres://user:pass@host:5432/db
	// <nil>
}

func ExampleDefaultContainer() {
	ProvideNamedValue(nil, "PG_URI", "postgres://user:pass@host:5432/db")
	uri, err := InvokeNamed[string](nil, "PG_URI")

	fmt.Println(uri)
	fmt.Println(err)
	// Output:
	// postgres://user:pass@host:5432/db
	// <nil>
}

func ExampleNewWithOpts() {
	container := NewWithOpts(&ContainerOpts{
		HookAfterShutdown: func(container *Container, serviceName string) {
			fmt.Printf("service shutdown: %s\n", serviceName)
		},
	})

	ProvideNamed(container, "PG_URI", func(i *Container) (string, error) {
		return "postgres://user:pass@host:5432/db", nil
	})
	MustInvokeNamed[string](container, "PG_URI")
	_ = container.Shutdown()

	// Output:
	// service shutdown: PG_URI
}

func ExampleContainer_ListProvidedServices() {
	container := New()

	ProvideNamedValue(container, "PG_URI", "postgres://user:pass@host:5432/db")
	services := container.ListProvidedServices()

	fmt.Println(services)
	// Output:
	// [PG_URI]
}

func ExampleContainer_ListInvokedServices_invoked() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "SERVICE_NAME", func(i *Container) (test, error) {
		return test{foobar: "foobar"}, nil
	})
	_, _ = InvokeNamed[test](container, "SERVICE_NAME")
	services := container.ListInvokedServices()

	fmt.Println(services)
	// Output:
	// [SERVICE_NAME]
}

func ExampleContainer_ListInvokedServices_notInvoked() {
	container := New()

	type test struct {
		foobar string
	}

	ProvideNamed(container, "SERVICE_NAME", func(i *Container) (test, error) {
		return test{foobar: "foobar"}, nil
	})
	services := container.ListInvokedServices()

	fmt.Println(services)
	// Output:
	// []
}

func ExampleContainer_HealthCheck() {
	container := New()

	Provide(container, dbServiceProvider)
	health := container.HealthCheck()

	fmt.Println(health)
	// Output:
	// map[*di.dbService:<nil>]
}

func ExampleContainer_Shutdown() {
	container := New()

	Provide(container, dbServiceProvider)
	err := container.Shutdown()

	fmt.Println(err)
	// Output:
	// <nil>
}

func ExampleContainer_Clone() {
	container := New()

	ProvideNamedValue(container, "PG_URI", "postgres://user:pass@host:5432/db")
	Container2 := container.Clone()
	services := Container2.ListProvidedServices()

	fmt.Println(services)
	// Output:
	// [PG_URI]
}
