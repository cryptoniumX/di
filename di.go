package di

import (
	"fmt"
)

func Provide[T any](i *Container, provider Provider[T]) {
	name := generateServiceName[T]()

	ProvideNamed[T](i, name, provider)
}

func ProvideNamed[T any](i *Container, name string, provider Provider[T]) {
	_i := getContainerOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	providerFn := toProviderFn[T](provider)
	service := newServiceLazy(name, providerFn)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

func ProvideValue[T any](i *Container, value T) {
	name := generateServiceName[T]()

	ProvideNamedValue[T](i, name, value)
}

func ProvideNamedValue[T any](i *Container, name string, value T) {
	_i := getContainerOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

func Override[T any](i *Container, provider Provider[T]) {
	name := generateServiceName[T]()

	OverrideNamed[T](i, name, provider)
}

func OverrideNamed[T any](i *Container, name string, provider Provider[T]) {
	_i := getContainerOrDefault(i)

	providerFn := toProviderFn[T](provider)
	service := newServiceLazy(name, providerFn)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

func OverrideValue[T any](i *Container, value T) {
	name := generateServiceName[T]()

	OverrideNamedValue[T](i, name, value)
}

func OverrideNamedValue[T any](i *Container, name string, value T) {
	_i := getContainerOrDefault(i)

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

func Invoke[T any](i *Container) (T, error) {
	name := generateServiceName[T]()
	return InvokeNamed[T](i, name)
}

func MustInvoke[T any](i *Container) T {
	s, err := Invoke[T](i)
	must(err)
	return s
}

func InvokeNamed[T any](i *Container, name string) (T, error) {
	return invokeImplem[T](i, name, "")
}

func MustInvokeNamed[T any](i *Container, name string) T {
	s, err := InvokeNamed[T](i, name)
	must(err)
	return s
}

func invokeByName(serviceName string, i *Container, dependencyName string, fallbackName string) (interface{}, error) {
	d, err := invokeImplem[any](i, dependencyName, fallbackName)
	if err != nil {
		return nil, fmt.Errorf("Failed to inject dependencies to service %s: %w", serviceName, err)
	}
	return d, nil
}

func invokeImplem[T any](i *Container, name string, fallbackName string) (T, error) {
	_i := getContainerOrDefault(i)

	var serviceAny any
	var ok bool

	serviceAny, ok = _i.get(name)
	if !ok {
		fallbackNames := []string{
			// if name is not found, try to find by pointer name
			fmt.Sprintf("*%s", name),
		}

		if fallbackName != "" {
			fallbackNames = append(fallbackNames, fallbackName)
		}

		for _, fallbackName := range fallbackNames {
			serviceAny, ok = _i.get(fallbackName)
			if ok {
				break
			}
		}

		if !ok {
			return empty[T](), _i.serviceNotFound(
				fmt.Sprintf("name: %s, fallbackName:%s", name, fallbackName),
			)
		}

	}

	service, ok := serviceAny.(Service)
	if !ok {
		return empty[T](), _i.serviceNotFound(name)
	}

	instanceAny, err := service.getInstance(_i)
	if err != nil {
		return empty[T](), err
	}

	_i.onServiceInvoke(name)

	if instance, ok := instanceAny.(T); ok {
		_i.logf("service %s invoked", name)
		return instance, nil
	}

	panic(fmt.Errorf("DI: service `%s` is not of type `%T`", name, empty[T]()))
}

func HealthCheck[T any](i *Container) error {
	name := generateServiceName[T]()
	return getContainerOrDefault(i).healthcheckImplem(name)
}

func HealthCheckNamed(i *Container, name string) error {
	return getContainerOrDefault(i).healthcheckImplem(name)
}

func Shutdown[T any](i *Container) error {
	name := generateServiceName[T]()
	return getContainerOrDefault(i).shutdownImplem(name)
}

func MustShutdown[T any](i *Container) {
	name := generateServiceName[T]()
	must(getContainerOrDefault(i).shutdownImplem(name))
}

func ShutdownNamed(i *Container, name string) error {
	return getContainerOrDefault(i).shutdownImplem(name)
}

func MustShutdownNamed(i *Container, name string) {
	must(getContainerOrDefault(i).shutdownImplem(name))
}
