package caching

import (
	"time"

	"github.com/moznion/go-optional"

	lru "github.com/arielsrv/golang-lru"
)

type CacheBuilder[TKey any, TValue any] struct {
	size     int
	duration time.Duration
}

func NewBuilder[TKey any, TValue any]() *CacheBuilder[TKey, TValue] {
	return &CacheBuilder[TKey, TValue]{}
}

func (c *CacheBuilder[TKey, TValue]) Size(size int) *CacheBuilder[TKey, TValue] {
	c.size = size
	return c
}

func (c *CacheBuilder[TKey, TValue]) ExpireAfterWrite(duration time.Duration) *CacheBuilder[TKey, TValue] {
	c.duration = duration
	return c
}

func (c *CacheBuilder[TKey, TValue]) Build() (ICache[TKey, TValue], error) {
	arcCache, err := lru.NewARCWithExpire(c.size, c.duration)
	if err != nil {
		return nil, err
	}

	return NewCache[TKey, TValue](arcCache), nil
}

type ICache[TKey any, TValue any] interface {
	GetIfPresent(key TKey) optional.Option[TValue]
	Put(key TKey, value TValue)
}

type Cache[TKey any, TValue any] struct {
	appCache *lru.ARCCache
}

func NewCache[TKey any, TValue any](appCache *lru.ARCCache) *Cache[TKey, TValue] {
	return &Cache[TKey, TValue]{appCache: appCache}
}

func (c Cache[TKey, TValue]) GetIfPresent(key TKey) optional.Option[TValue] {
	value, hasValue := c.appCache.Get(key)
	if !hasValue {
		var nilValue optional.Option[TValue]
		return nilValue
	}
	return optional.Some[TValue](value.(TValue))
}

func (c Cache[TKey, TValue]) Put(key TKey, value TValue) {
	c.appCache.Add(key, value)
}
