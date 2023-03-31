package caching_test

import (
	"testing"
	"time"

	"examples/caching"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheErr(t *testing.T) {
	appCache, err := caching.NewBuilder[string, string]().
		Size(-1).
		ExpireAfterWrite(time.Duration(2000) * time.Millisecond).
		Build()

	assert.Error(t, err)
	assert.Nil(t, appCache)
}

func TestCache_GetIfPresent(t *testing.T) {
	appCache, err := caching.NewBuilder[string, string]().
		Size(2).
		ExpireAfterWrite(time.Duration(5) * time.Minute).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	appCache.Put("key", "value")
	value := appCache.GetIfPresent("key")

	assert.NotNil(t, value)
	assert.True(t, value.IsSome())

	actual, err := value.Take()
	assert.NoError(t, err)
	assert.Equal(t, "value", actual)
}

type OrderDTO struct {
	ID int64 `json:"order_id,omitempty"`
}

func TestCache_GetIfPresentForStruct(t *testing.T) {
	appCache, err := caching.NewBuilder[int, OrderDTO]().
		Size(2).
		ExpireAfterWrite(time.Duration(2000) * time.Millisecond).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	orderDTO := new(OrderDTO)
	orderDTO.ID = int64(1)

	appCache.Put(1, *orderDTO)
	value := appCache.GetIfPresent(1)

	assert.NotNil(t, value)
	assert.True(t, value.IsSome())
	actual, err := value.Take()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), actual.ID)
}

func TestCache_GetIfNotPresentString(t *testing.T) {
	appCache, err := caching.NewBuilder[string, string]().
		Size(2).
		ExpireAfterWrite(time.Duration(2000) * time.Millisecond).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	value := appCache.GetIfPresent("key")

	assert.True(t, value.IsNone())
}

func TestCache_GetIfNotPresentInt(t *testing.T) {
	appCache, err := caching.NewBuilder[string, int]().
		Size(2).
		ExpireAfterWrite(time.Duration(2000) * time.Millisecond).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	value := appCache.GetIfPresent("key")

	assert.True(t, value.IsNone())
}

func TestCache_GetIfNotPresentStruct(t *testing.T) {
	appCache, err := caching.NewBuilder[string, OrderDTO]().
		Size(2).
		ExpireAfterWrite(time.Duration(2000) * time.Millisecond).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	value := appCache.GetIfPresent("key")

	assert.True(t, value.IsNone())
}

func TestCacheBuilder_ExpireAfterWrite(t *testing.T) {
	appCache, err := caching.NewBuilder[string, string]().
		Size(2).
		ExpireAfterWrite(time.Duration(100) * time.Millisecond).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	appCache.Put("key", "value")
	time.Sleep(time.Duration(200) * time.Millisecond)

	value := appCache.GetIfPresent("key")

	assert.True(t, value.IsNone())
}

func TestCacheBuilder_Size(t *testing.T) {
	appCache, err := caching.NewBuilder[string, string]().
		Size(1).
		ExpireAfterWrite(time.Duration(100) * time.Minute).
		Build()

	assert.NoError(t, err)
	assert.NotNil(t, appCache)

	appCache.Put("key1", "value1")
	appCache.Put("key2", "value2")

	time.Sleep(time.Duration(1000) * time.Millisecond)

	value1 := appCache.GetIfPresent("key1")
	value2 := appCache.GetIfPresent("key2")

	actual := (value1.IsSome() && !value2.IsSome()) || (!value1.IsSome() && value2.IsSome())
	assert.True(t, actual)
}
