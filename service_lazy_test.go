package di

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type lazyTest struct {
	idx int
	err error
}

func (t *lazyTest) Shutdown() error {
	return t.err
}

func TestServiceLazyName(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	provider1 := func(i *Container) (int, error) {
		return 42, nil
	}
	provider2 := func(i *Container) (test, error) {
		return _test, nil
	}

	service1 := newServiceLazy("foobar", toProviderFn[int](provider1))
	is.Equal("foobar", service1.getName())

	service2 := newServiceLazy("foobar", toProviderFn[test](provider2))
	is.Equal("foobar", service2.getName())
}

func TestServiceLazyInstance(t *testing.T) {
	is := assert.New(t)

	type test struct {
		foobar string
	}
	_test := test{foobar: "foobar"}

	provider1 := func(i *Container) (int, error) {
		return 42, nil
	}
	provider2 := func(i *Container) (test, error) {
		return _test, nil
	}
	provider3 := func(i *Container) (int, error) {
		panic("error")
	}
	provider4 := func(i *Container) (int, error) {
		panic(fmt.Errorf("error"))
	}
	provider5 := func(i *Container) (int, error) {
		return 42, fmt.Errorf("error")
	}

	i := New()

	service1 := newServiceLazy("foobar", toProviderFn[int](provider1))
	instance1, err1 := service1.getInstance(i)
	is.Nil(err1)
	is.Equal(42, instance1)

	service2 := newServiceLazy("hello", toProviderFn[test](provider2))
	instance2, err2 := service2.getInstance(i)
	is.Nil(err2)
	is.Equal(_test, instance2)

	is.Panics(func() {
		service3 := newServiceLazy("baz", toProviderFn[int](provider3))
		_, _ = service3.getInstance(i)
	})

	is.NotPanics(func() {
		service4 := newServiceLazy("plop", toProviderFn[int](provider4))
		instance4, err4 := service4.getInstance(i)
		is.NotNil(err4)
		is.Empty(instance4)
		expected := fmt.Errorf("error")
		is.Equal(expected, err4)
	})

	is.NotPanics(func() {
		service5 := newServiceLazy("plop", toProviderFn[int](provider5))
		instance5, err5 := service5.getInstance(i)
		is.NotNil(err5)
		is.Empty(instance5)
		expected := fmt.Errorf("error")
		is.Equal(expected, err5)
	})
}

func TestServiceLazyInstanceShutDown(t *testing.T) {
	is := assert.New(t)

	index := 1
	provider1 := func(i *Container) (*lazyTest, error) {
		index++
		return &lazyTest{index, nil}, nil
	}
	provider2 := func(i *Container) (*lazyTest, error) {
		index++
		return &lazyTest{index, assert.AnError}, nil
	}

	i := New()

	service1 := newServiceLazy("foobar", toProviderFn[*lazyTest](provider1))
	instance1, err := service1.getInstance(i)
	assert.NotNil(t, instance1)
	is.Nil(err)
	is.True(service1.(*ServiceLazy).built)
	err = service1.shutdown()
	is.False(service1.(*ServiceLazy).built)
	is.Nil(err)
	instance2, err := service1.getInstance(i)
	assert.NotNil(t, instance2)
	is.Nil(err)

	service2 := newServiceLazy("foobar", toProviderFn[*lazyTest](provider2)).(*ServiceLazy)
	is.False(service2.built)
	is.Nil(err)
	err = service2.build(i)
	is.Nil(err)
	is.True(service2.built)
	err = service2.shutdown()
	is.Error(assert.AnError, err)
	is.True(service2.built)
}
