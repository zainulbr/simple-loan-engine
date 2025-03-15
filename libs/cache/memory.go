package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	value      interface{}
	expiration *time.Timer
}

type memoryCache struct {
	items sync.Map // map[any]any but is safe for concurrent use by multiple goroutines
}

// NewMemoryCache creates a new instance of MemoryCache
func NewMemoryCache() Cache {
	return &memoryCache{}
}

// Set stores a value in cache with an optional TTL
func (c *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	// If TTL is greater than 0, set expiration
	var timer *time.Timer
	if ttl > 0 {
		timer = time.AfterFunc(ttl, func() {
			c.Delete(key)
		})
	}

	c.items.Store(key, CacheItem{value: value, expiration: timer})
	return nil
}

// Get retrieves a value from cache
func (c *memoryCache) Get(key string) (interface{}, error) {
	item, ok := c.items.Load(key)
	if !ok {
		return nil, fmt.Errorf("Item not found")
	}
	return item.(CacheItem).value, nil
}

// Delete removes an item from cache
func (c *memoryCache) Delete(key string) error {
	if item, ok := c.items.Load(key); ok {
		if item.(CacheItem).expiration != nil {
			item.(CacheItem).expiration.Stop()
		}
		c.items.Delete(key)
	}
	return nil
}

// Flush clears all items in cache
func (c *memoryCache) Flush() error {
	c.items.Range(func(key, value interface{}) bool {
		if value.(CacheItem).expiration != nil {
			value.(CacheItem).expiration.Stop()
		}
		c.items.Delete(key)
		return true
	})
	return nil
}
