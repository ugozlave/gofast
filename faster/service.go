package faster

import (
	"sync"

	"github.com/ugozlave/gofast"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
	Close()
}

/*
** MemoryCache
 */

type MemoryCache struct {
	store map[string]any
	mu    sync.RWMutex
}

func MemoryCacheBuilder() Builder[*MemoryCache] {
	return func(ctx *gofast.BuilderContext) *MemoryCache {
		return NewMemoryCache(ctx)
	}
}

func NewMemoryCache(ctx *gofast.BuilderContext) *MemoryCache {
	return &MemoryCache{
		store: make(map[string]any),
	}
}

func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, found := c.store[key]
	return value, found
}

func (c *MemoryCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *MemoryCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]any)
}
