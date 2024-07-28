package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	m  map[string]cacheEntry
	mu sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		m: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	createdAt := time.Now()
	newEntry := cacheEntry{
		createdAt: createdAt,
		val:       val,
	}
	c.m[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.m[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.mu.Lock()
		currentTime := time.Now()
		for key, entry := range c.m {
			if currentTime.Sub(entry.createdAt) >= interval {
				delete(c.m, key)
			}
		}
		c.mu.Unlock()
	}
}
