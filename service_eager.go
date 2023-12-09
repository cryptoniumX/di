package di

type serviceEager struct {
	name     string
	instance any
}

func newServiceEager(name string, instance any) Service {
	return &serviceEager{
		name:     name,
		instance: instance,
	}
}

//nolint:unused
func (s *serviceEager) getName() string {
	return s.name
}

//nolint:unused
func (s *serviceEager) getInstance(i *Container) (any, error) {
	return s.instance, nil
}

func (s *serviceEager) healthcheck() error {
	instance, ok := any(s.instance).(Healthcheckable)
	if ok {
		return instance.HealthCheck()
	}

	return nil
}

func (s *serviceEager) shutdown() error {
	instance, ok := any(s.instance).(Shutdownable)
	if ok {
		return instance.Shutdown()
	}

	return nil
}

func (s *serviceEager) clone() any {
	return s
}
