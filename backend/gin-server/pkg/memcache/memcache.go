package memcache

import "time"

// Cache defines an interface for token storage
type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
}

// InMemoryCache implements Cache interface with a local map
type InMemoryCache struct {
	data    map[string]cacheItem
	cleanup time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache with auto cleanup
func NewInMemoryCache(cleanupInterval time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		data:    make(map[string]cacheItem),
		cleanup: cleanupInterval,
	}

	go cache.startCleanupTimer()

	return cache
}

func (c *InMemoryCache) startCleanupTimer() {
	ticker := time.NewTicker(c.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		c.deleteExpired()
	}
}

func (c *InMemoryCache) deleteExpired() {
	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiration) {
			delete(c.data, key)
		}
	}
}

// Set adds a value to the cache with an expiration
func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(key string) error {
	delete(c.data, key)
	return nil
}
