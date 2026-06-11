package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(k string, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[k] = cacheEntry{
		createdAt: time.Now(),
		val:       v,
	}
}

func (c *Cache) Get(k string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.entries[k]
	if !ok {
		return nil, false
	}

	return v.val, true
}

// loops forever clearing cache entries older than the interval
func (c *Cache) reapLoop(interval time.Duration) {
	for {
		c.mu.Lock()
		for k, entry := range c.entries {
			if time.Now().After(entry.createdAt.Add(interval)) {
				delete(c.entries, k)
			}
			// reap
		}
		c.mu.Unlock()
		time.Sleep(interval)
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{}
	cache.entries = map[string]cacheEntry{}
	go cache.reapLoop(interval)

	return cache
}
