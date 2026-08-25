package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

// Cache — thread-safe TTL in-memory cache with background eviction
type Cache struct {
	mu      sync.RWMutex
	entries map[string]entry
}

func New() *Cache {
	c := &Cache{entries: make(map[string]entry)}
	go c.evict()
	return c
}

func (c *Cache) Set(key string, val any, ttl time.Duration) {
	c.mu.Lock()
	c.entries[key] = entry{value: val, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

func (c *Cache) Flush() {
	c.mu.Lock()
	c.entries = make(map[string]entry)
	c.mu.Unlock()
}

func (c *Cache) Size() int {
	c.mu.RLock()
	n := len(c.entries)
	c.mu.RUnlock()
	return n
}

func (c *Cache) evict() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.entries {
			if now.After(e.expiresAt) {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
