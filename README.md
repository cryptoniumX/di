# DI - Dependency Injection

> This is a forked of the original [Do](https://github.com/samber/do) package with an extension to allow dynamically injecting depdencies to a struct through reflection.


![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-%23007d9c)
[![GoDoc](https://godoc.org/github.com/cryptoniumX/di?status.svg)](https://pkg.go.dev/github.com/cryptoniumX/di)
![Build Status](https://github.com/cryptoniumX/di/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/cryptoniumX/di)](https://goreportcard.com/report/github.com/cryptoniumX/di)
[![Coverage](https://img.shields.io/codecov/c/github/samber/do)](https://codecov.io/gh/samber/do)
[![License](https://img.shields.io/github/license/samber/do)](./LICENSE)

**⚙️ A dependency injection toolkit based on Go 1.18+ Generics.**

This library implements the Dependency Injection design pattern. It may replace the `uber/dig` fantastic package in simple Go projects. `samber/do` uses Go 1.18+ generics and therefore is typesafe.

## 💡 Features

- Service registration
- Service invocation
- Service health check
- Service shutdown
- Service lifecycle hooks
- Named or anonymous services
- Eagerly or lazily loaded services
- Dependency graph resolution
- Default Container
- Container cloning
- Service override
- Lightweight, no dependencies
- No code generation

🚀 Services are loaded in invocation order.

🕵️ Service health can be checked individually or globally. Services implementing `di.Healthcheckable` interface will be called via `di.HealthCheck[type]()` or `Container.HealthCheck()`.

🛑 Services can be shutdowned properly, in back-initialization order. Services implementing `di.Shutdownable` interface will be called via `di.Shutdown[type]()` or `Container.Shutdown()`.

## Di compared to original Do package
We added a method to allow injecting dependencies dynamically through struct reflection.
Declare a service, and add a `di:"<name>:` tag to the field where you want the container to inject the corresponding dependency.

```go
// main.go
container := New()
repository := newRepository()
redisClient := newRedisClient()
ProvideValue[Repository](container, repository)
ProvideValue[RedisClient](container, redisClient)

// service.go
type service struct {
	Repository  Repository  `di:"repository"`
	RedisClient RedisClient `di:"redisClient"`
}

func newService(
    container *di.Container
) (&service, error) {
    s := service{}
    err := container.Inject(&s)
    return s, err
}
```

## 🚀 Install

```sh
go get github.com/cryptoniumX/di@v1
```

This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0.

This library has no dependencies except the Go std lib.

## 💡 Quick start

You can import `do` using:

```go
import (
    "github.com/cryptoniumX/di"
)
```

Then instanciate services:

```go
func main() {
    container := di.New()

    // provides CarService
    di.Provide(container, NewCarService)

    // provides EngineService
    di.Provide(container, NewEngineService)

    car := di.MustInvoke[*CarService](container)
    car.Start()
    // prints "car starting"

    di.HealthCheck[EngineService](container)
    // returns "engine broken"

    // container.ShutdownOnSIGTERM()    // will block until receiving sigterm signal
    container.Shutdown()
    // prints "car stopped"
}
```

Services:

```go
type EngineService interface{}

func NewEngineService(i *di.Container) (EngineService, error) {
    return &engineServiceImplem{}, nil
}

type engineServiceImplem struct {}

// [Optional] Implements di.Healthcheckable.
func (c *engineServiceImplem) HealthCheck() error {
	return fmt.Errorf("engine broken")
}
```

```go
func NewCarService(i *di.Container) (*CarService, error) {
    engine := di.MustInvoke[EngineService](i)
    car := CarService{Engine: engine}
    return &car, nil
}

type CarService struct {
	Engine EngineService
}

func (c *CarService) Start() {
	println("car starting")
}

// [Optional] Implements di.Shutdownable.
func (c *CarService) Shutdown() error {
	println("car stopped")
	return nil
}
```

## 🤠 Spec

[GoDoc: https://godoc.org/github.com/cryptoniumX/di](https://godoc.org/github.com/cryptoniumX/di)

Container:

- [di.New](https://pkg.go.dev/github.com/cryptoniumX/di#New)
- [di.NewWithOpts](https://pkg.go.dev/github.com/cryptoniumX/di#NewWithOpts)
  - [Container.Clone](https://pkg.go.dev/github.com/cryptoniumX/di#Container.Clone)
  - [Container.CloneWithOpts](https://pkg.go.dev/github.com/cryptoniumX/di#Container.CloneWithOpts)
  - [Container.HealthCheck](https://pkg.go.dev/github.com/cryptoniumX/di#Container.HealthCheck)
  - [Container.Shutdown](https://pkg.go.dev/github.com/cryptoniumX/di#Container.Shutdown)
  - [Container.ShutdownOnSIGTERM](https://pkg.go.dev/github.com/cryptoniumX/di#Container.ShutdownOnSIGTERM)
  - [Container.ShutdownOnSignals](https://pkg.go.dev/github.com/cryptoniumX/di#Container.ShutdownOnSignals)
  - [Container.ListProvidedServices](https://pkg.go.dev/github.com/cryptoniumX/di#Container.ListProvidedServices)
  - [Container.ListInvokedServices](https://pkg.go.dev/github.com/cryptoniumX/di#Container.ListInvokedServices)
- [di.HealthCheck](https://pkg.go.dev/github.com/cryptoniumX/di#HealthCheck)
- [di.HealthCheckNamed](https://pkg.go.dev/github.com/cryptoniumX/di#HealthCheckNamed)
- [di.Shutdown](https://pkg.go.dev/github.com/cryptoniumX/di#Shutdown)
- [di.ShutdownNamed](https://pkg.go.dev/github.com/cryptoniumX/di#ShutdownNamed)
- [di.MustShutdown](https://pkg.go.dev/github.com/cryptoniumX/di#MustShutdown)
- [di.MustShutdownNamed](https://pkg.go.dev/github.com/cryptoniumX/di#MustShutdownNamed)

Service registration:

- [di.Provide](https://pkg.go.dev/github.com/cryptoniumX/di#Provide)
- [di.ProvideNamed](https://pkg.go.dev/github.com/cryptoniumX/di#ProvideNamed)
- [di.ProvideNamedValue](https://pkg.go.dev/github.com/cryptoniumX/di#ProvideNamedValue)
- [di.ProvideValue](https://pkg.go.dev/github.com/cryptoniumX/di#ProvideValue)

Service invocation:

- [di.Invoke](https://pkg.go.dev/github.com/cryptoniumX/di#Invoke)
- [di.MustInvoke](https://pkg.go.dev/github.com/cryptoniumX/di#MustInvoke)
- [di.InvokeNamed](https://pkg.go.dev/github.com/cryptoniumX/di#InvokeNamed)
- [di.MustInvokeNamed](https://pkg.go.dev/github.com/cryptoniumX/di#MustInvokeNamed)

Service override:

- [di.Override](https://pkg.go.dev/github.com/cryptoniumX/di#Override)
- [di.OverrideNamed](https://pkg.go.dev/github.com/cryptoniumX/di#OverrideNamed)
- [di.OverrideNamedValue](https://pkg.go.dev/github.com/cryptoniumX/di#OverrideNamedValue)
- [di.OverrideValue](https://pkg.go.dev/github.com/cryptoniumX/di#OverrideValue)

### Container (DI container)

Build a container for your components. `Container` is responsible for building services in the right order, and managing service lifecycle.

```go
container := di.New()
```

Or use `nil` as the default Container:

```go
di.Provide(nil, func (i *Container) (int, error) {
    return 42, nil
})

service := di.MustInvoke[int](nil)
```

You can check health of services implementing `func HealthCheck() error`.

```go
type DBService struct {
    db *sql.DB
}

func (s *DBService) HealthCheck() error {
    return s.db.Ping()
}

container := di.New()
di.Provide(container, ...)
di.Invoke(container, ...)

statuses := container.HealthCheck()
// map[string]error{
//   "*DBService": nil,
// }
```

De-initialize all compoments properly. Services implementing `func Shutdown() error` will be called synchronously in back-initialization order.

```go
type DBService struct {
    db *sql.DB
}

func (s *DBService) Shutdown() error {
    return s.db.Close()
}

container := di.New()
di.Provide(container, ...)
di.Invoke(container, ...)

// shutdown all services in reverse order
container.Shutdown()
```

List services:

```go
type DBService struct {
    db *sql.DB
}

container := di.New()

di.Provide(container, ...)
println(di.ListProvidedServices())
// output: []string{"*DBService"}

di.Invoke(container, ...)
println(di.ListInvokedServices())
// output: []string{"*DBService"}
```

### Service registration

Services can be registered in multiple way:

- with implicit name (struct or interface name)
- with explicit name
- eagerly
- lazily

Anonymous service, loaded lazily:

```go
type DBService struct {
    db *sql.DB
}

di.Provide[DBService](container, func(i *container) (*DBService, error) {
    db, err := sql.Open(...)
    if err != nil {
        return nil, err
    }

    return &DBService{db: db}, nil
})
```

Named service, loaded lazily:

```go
type DBService struct {
    db *sql.DB
}

di.ProvideNamed(container, "dbconn", func(i *container) (*DBService, error) {
    db, err := sql.Open(...)
    if err != nil {
        return nil, err
    }

    return &DBService{db: db}, nil
})
```

Anonymous service, loaded eagerly:

```go
type Config struct {
    uri string
}

di.ProvideValue[Config](container, Config{uri: "postgres://user:pass@host:5432/db"})
```

Named service, loaded eagerly:

```go
type Config struct {
    uri string
}

di.ProvideNamedValue(container, "configuration", Config{uri: "postgres://user:pass@host:5432/db"})
```

### Service invocation

Loads anonymous service:

```go
type DBService struct {
    db *sql.DB
}

dbService, err := di.Invoke[DBService](container)
```

Loads anonymous service or panics if service was not registered:

```go
type DBService struct {
    db *sql.DB
}

dbService := di.MustInvoke[DBService](container)
```

Loads named service:

```go
config, err := di.InvokeNamed[Config](container, "configuration")
```

Loads named service or panics if service was not registered:

```go
config := di.MustInvokeNamed[Config](container, "configuration")
```

### Individual service healthcheck

Check health of anonymous service:

```go
type DBService struct {
    db *sql.DB
}

dbService, err := di.Invoke[DBService](container)
err = di.HealthCheck[DBService](container)
```

Check health of named service:

```go
config, err := di.InvokeNamed[Config](container, "configuration")
err = di.HealthCheckNamed(container, "configuration")
```

### Individual service shutdown

Unloads anonymous service:

```go
type DBService struct {
    db *sql.DB
}

dbService, err := di.Invoke[DBService](container)
err = di.Shutdown[DBService](container)
```

Unloads anonymous service or panics if service was not registered:

```go
type DBService struct {
    db *sql.DB
}

dbService := di.MustInvoke[DBService](container)
di.MustShutdown[DBService](container)
```

Unloads named service:

```go
config, err := di.InvokeNamed[Config](container, "configuration")
err = di.ShutdownNamed(container, "configuration")
```

Unloads named service or panics if service was not registered:

```go
config := di.MustInvokeNamed[Config](container, "configuration")
di.MustShutdownNamed(container, "configuration")
```

### Service override

By default, providing a service twice will panic. Service can be replaced at runtime using `di.Override` helper.

```go
di.Provide[Vehicle](container, func (i *di.Container) (Vehicle, error) {
    return &CarImplem{}, nil
})

di.Override[Vehicle](container, func (i *di.Container) (Vehicle, error) {
    return &BusImplem{}, nil
})
```

### Hooks

2 lifecycle hooks are available in Containers:

- After registration
- After shutdown

```go
container := di.NewWithOpts(&di.ContainerOpts{
    HookAfterRegistration: func(container *di.Container, serviceName string) {
        fmt.Printf("Service registered: %s\n", serviceName)
    },
    HookAfterShutdown: func(container *di.Container, serviceName string) {
        fmt.Printf("Service stopped: %s\n", serviceName)
    },

    Logf: func(format string, args ...any) {
        log.Printf(format, args...)
    },
})
```

### Cloning Container

Cloned Container have same service registrations as it's parent, but it doesn't share invoked service state.

Clones are useful for unit testing by replacing some services to mocks.

```go
var container *di.Container;

func init() {
    di.Provide[Service](container, func (i *di.Container) (Service, error) {
        return &RealService{}, nil
    })
    di.Provide[*App](container, func (i *di.Container) (*App, error) {
        return &App{i.MustInvoke[Service](i)}, nil
    })
}

func TestService(t *testing.T) {
    i := container.Clone()
    defer i.Shutdown()

    // replace Service to MockService
    di.Override[Service](i, func (i *di.Container) (Service, error) {
        return &MockService{}, nil
    }))

    app := di.Invoke[*App](i)
    // do unit testing with mocked service
}
```

## 🛩 Benchmark

// @TODO

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/cryptoniumX/di)
- Fix [open issues](https://github.com/cryptoniumX/di/issues) or request new features

Don't hesitate ;)

### With Docker

```bash
docker-compose run --rm dev
```

### Without Docker

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/do)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2022 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
