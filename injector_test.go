package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Repository interface {
	GetID() string
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
	ProvideValue[Repository](container, repository)
	ProvideValue[RedisClient](container, redisClient)
	ProvideNamedValue[float64](container, "float64Container", 69.69)

	type service struct {
		Repository  Repository  `di:"repository"`
		RedisClient RedisClient `di:"redisClient"`
		Test        *test       `di:"test"`
		Float64     float64     `di:"float64Container"`
	}

	s := service{}
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
