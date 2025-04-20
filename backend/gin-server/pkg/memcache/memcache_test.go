package memcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)

	// Test set and get
	err := cache.Set("key1", "value1", time.Minute)
	assert.NoError(t, err)

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Test non-existent key
	_, exists = cache.Get("nonexistent")
	assert.False(t, exists)

	// Test expiration
	err = cache.Set("key2", "value2", 10*time.Millisecond)
	assert.NoError(t, err)

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	_, exists = cache.Get("key2")
	assert.False(t, exists)

	// Test delete
	err = cache.Set("key3", "value3", time.Minute)
	assert.NoError(t, err)

	err = cache.Delete("key3")
	assert.NoError(t, err)

	_, exists = cache.Get("key3")
	assert.False(t, exists)
}
