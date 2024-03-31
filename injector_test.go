package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Repository[T any] interface {
	GetID() T
}

type repository struct {
	id string
}

func newRepository() *repository {
	return &repository{id: "42"}
}

func (r *repository) GetID() string {
	return r.id
}

type RedisClient interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type redisClient struct {
}

func newRedisClient() *redisClient {
	return &redisClient{}
}

func (r *redisClient) Get(key string) (string, error) {
	return "value", nil
}

func (r *redisClient) Set(key string, value string) error {
	return nil
}

func TestInject(t *testing.T) {
	type test struct {
		foobar string
	}

	container := New()
	Provide(container, func(i *Container) (*test, error) {
		return &test{foobar: "foobar"}, nil
	})

	repository := newRepository()
	redisClient := newRedisClient()
	ProvideValue[Repository[string]](container, repository)
	ProvideValue[RedisClient](container, redisClient)
	ProvideNamedValue[float64](container, "float64Container", 69.69)

	type TestService struct {
		Repository  Repository[string] `di:""`
		RedisClient RedisClient        `di:""`
		Test        *test              `di:""`
		Float64     float64            `di:"float64Container"`
	}

	s := TestService{}
	err := container.Inject(&s)
	assert.NoError(t, err)

	id := s.Repository.GetID()
	assert.Equal(t, "42", id)

	value, err := s.RedisClient.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, value, "value")
	assert.Equal(t, s.Test.foobar, "foobar")
	assert.Equal(t, s.Float64, 69.69)
}

func TestInjectFailPassByValue(t *testing.T) {
	type service struct {
		Dependency1 int    `di:"depdency1"`
		Dependency2 string `di:"depdency2"`
	}

	container := New()
	Provide(container, func(i *Container) (int, error) {
		return 3, nil
	})
	Provide(container, func(i *Container) (string, error) {
		return "string", nil
	})

	s := service{}

	// try to pass by value
	expectedMsg := "Inject: Must pass a pointer to a struct"
	err := container.Inject(s)
	assert.Containsf(t, err.Error(), expectedMsg, "expected error containing %q, got %s", expectedMsg, err)

	// try to pass a map
	m := make(map[string]interface{})
	err = container.Inject(m)
	assert.Containsf(t, err.Error(), expectedMsg, "expected error containing %q, got %s", expectedMsg, err)
	err = container.Inject(&s)
	assert.NoError(t, err)

}
