package di

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var DefaultContainer = New()

func getContainerOrDefault(i *Container) *Container {
	if i != nil {
		return i
	}

	return DefaultContainer
}

func New() *Container {
	return NewWithOpts(&ContainerOpts{})
}

type ContainerOpts struct {
	HookAfterRegistration func(injector *Container, serviceName string)
	HookAfterShutdown     func(injector *Container, serviceName string)

	Logf func(format string, args ...any)
}

func NewWithOpts(opts *ContainerOpts) *Container {
	logf := opts.Logf
	if logf == nil {
		logf = func(format string, args ...any) {}
	}

	logf("injector created")

	return &Container{
		mu:       sync.RWMutex{},
		services: make(map[string]any),

		orderedInvocation:      map[string]int{},
		orderedInvocationIndex: 0,

		hookAfterRegistration: opts.HookAfterRegistration,
		hookAfterShutdown:     opts.HookAfterShutdown,

		logf: logf,
	}
}

type Container struct {
	mu       sync.RWMutex
	services map[string]any

	// It should be a graph instead of simple ordered list.
	orderedInvocation      map[string]int // map is faster than slice
	orderedInvocationIndex int

	hookAfterRegistration func(injector *Container, serviceName string)
	hookAfterShutdown     func(injector *Container, serviceName string)

	logf func(format string, args ...any)
}

func (i *Container) ListProvidedServices() []string {
	i.mu.RLock()
	names := keys(i.services)
	i.mu.RUnlock()

	i.logf("exported list of services: %v", names)

	return names
}

func (i *Container) ListInvokedServices() []string {
	i.mu.RLock()
	names := keys(i.orderedInvocation)
	i.mu.RUnlock()

	i.logf("exported list of invoked services: %v", names)

	return names
}

func (i *Container) HealthCheck() map[string]error {
	i.mu.RLock()
	names := keys(i.services)
	i.mu.RUnlock()

	i.logf("requested healthcheck")

	results := map[string]error{}

	for _, name := range names {
		results[name] = i.healthcheckImplem(name)
	}

	i.logf("got healthcheck results: %v", results)

	return results
}

func (i *Container) Shutdown() error {
	i.mu.RLock()
	invocations := invertMap(i.orderedInvocation)
	i.mu.RUnlock()

	i.logf("requested shutdown")

	for index := i.orderedInvocationIndex; index >= 0; index-- {
		name, ok := invocations[index]
		if !ok {
			continue
		}

		err := i.shutdownImplem(name)
		if err != nil {
			return err
		}
	}

	i.logf("shutdowned services")

	return nil
}

// ShutdownOnSIGTERM listens for sigterm signal in order to graceful stop service.
// It will block until receiving a sigterm signal.
func (i *Container) ShutdownOnSIGTERM() error {
	return i.ShutdownOnSignals(syscall.SIGTERM)
}

// ShutdownOnSignals listens for signals defined in signals parameter in order to graceful stop service.
// It will block until receiving any of these signal.
// If no signal is provided in signals parameter, syscall.SIGTERM will be added as default signal.
func (i *Container) ShutdownOnSignals(signals ...os.Signal) error {
	// Make sure there is at least syscall.SIGTERM as a signal
	if len(signals) < 1 {
		signals = append(signals, syscall.SIGTERM)
	}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	<-ch
	signal.Stop(ch)
	close(ch)

	return i.Shutdown()
}

func (i *Container) healthcheckImplem(name string) error {
	i.mu.Lock()

	serviceAny, ok := i.services[name]
	if !ok {
		i.mu.Unlock()
		return fmt.Errorf("DI: could not find service `%s`", name)
	}

	i.mu.Unlock()

	service, ok := serviceAny.(healthcheckableService)
	if ok {
		i.logf("requested healthcheck for service %s", name)

		err := service.healthcheck()
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Container) shutdownImplem(name string) error {
	i.mu.Lock()

	serviceAny, ok := i.services[name]
	if !ok {
		i.mu.Unlock()
		return fmt.Errorf("DI: could not find service `%s`", name)
	}

	i.mu.Unlock()

	service, ok := serviceAny.(shutdownableService)
	if ok {
		i.logf("requested shutdown for service %s", name)

		err := service.shutdown()
		if err != nil {
			return err
		}
	}

	delete(i.services, name)
	delete(i.orderedInvocation, name)

	i.onServiceShutdown(name)

	return nil
}

func (i *Container) exists(name string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	_, ok := i.services[name]
	return ok
}

func (i *Container) get(name string) (any, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	s, ok := i.services[name]
	return s, ok
}

func (i *Container) set(name string, service any) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.services[name] = service

	// defering hook call will unlock mutex
	defer i.onServiceRegistration(name)
}

func (i *Container) remove(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.services, name)
}

func (i *Container) forEach(cb func(name string, service any)) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for name, service := range i.services {
		cb(name, service)
	}
}

func (i *Container) serviceNotFound(name string) error {
	// @TODO: use the Keys+Map functions from `golang.org/x/exp/maps` as
	// soon as it is released in stdlib.
	servicesNames := keys(i.services)
	servicesNames = mAp(servicesNames, func(name string) string {
		return fmt.Sprintf("`%s`", name)
	})

	return fmt.Errorf("DI: could not find service `%s`, available services: %s", name, strings.Join(servicesNames, "\n"))
}

func (i *Container) onServiceInvoke(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.orderedInvocation[name]; !ok {
		i.orderedInvocation[name] = i.orderedInvocationIndex
		i.orderedInvocationIndex++
	}
}

func (i *Container) onServiceRegistration(name string) {
	if i.hookAfterRegistration != nil {
		i.hookAfterRegistration(i, name)
	}
}

func (i *Container) onServiceShutdown(name string) {
	if i.hookAfterShutdown != nil {
		i.hookAfterShutdown(i, name)
	}
}

// Clone clones injector with provided services but not with invoked instances.
func (i *Container) Clone() *Container {
	return i.CloneWithOpts(&ContainerOpts{})
}

// CloneWithOpts clones injector with provided services but not with invoked instances, with options.
func (i *Container) CloneWithOpts(opts *ContainerOpts) *Container {
	clone := NewWithOpts(opts)

	i.mu.RLock()
	defer i.mu.RUnlock()

	for name, serviceAny := range i.services {
		if service, ok := serviceAny.(cloneableService); ok {
			clone.services[name] = service.clone()
		} else {
			clone.services[name] = service
		}
		defer clone.onServiceRegistration(name)
	}

	i.logf("injector cloned")

	return clone
}
