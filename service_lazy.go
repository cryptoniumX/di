package di

import (
	"sync"
)

type Provider[T any] func(*Container) (T, error)
type providerFn func(*Container) (any, error)

func toProviderFn[T any](provider Provider[T]) providerFn {
	return func(Container *Container) (any, error) {
		result, err := provider(Container)
		return result, err
	}
}

type ServiceLazy struct {
	mu       sync.RWMutex
	name     string
	instance any

	// lazy loading
	built    bool
	provider providerFn
}

func newServiceLazy(name string, provider providerFn) Service {
	return &ServiceLazy{
		name: name,

		built:    false,
		provider: provider,
	}
}

//nolint:unused
func (s *ServiceLazy) getName() string {
	return s.name
}

//nolint:unused
func (s *ServiceLazy) getInstance(i *Container) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.built {
		err := s.build(i)
		if err != nil {
			return nil, err
		}
	}

	return s.instance, nil
}

//nolint:unused
func (s *ServiceLazy) build(i *Container) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	instance, err := s.provider(i)
	if err != nil {
		return err
	}

	s.instance = instance
	s.built = true

	return nil
}

func (s *ServiceLazy) healthcheck() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.built {
		return nil
	}

	instance, ok := any(s.instance).(Healthcheckable)
	if ok {
		return instance.HealthCheck()
	}

	return nil
}

func (s *ServiceLazy) shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.built {
		return nil
	}

	instance, ok := any(s.instance).(Shutdownable)
	if ok {
		err := instance.Shutdown()
		if err != nil {
			return err
		}
	}

	s.built = false
	s.instance = nil

	return nil
}

func (s *ServiceLazy) clone() any {
	// reset `build` flag and instance
	return &ServiceLazy{
		name: s.name,

		built:    false,
		provider: s.provider,
	}
}
