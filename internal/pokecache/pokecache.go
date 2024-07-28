package pokecache

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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
	val, err := compress(val)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	value, err := decompress(entry.val)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return value, true
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

// Compress the data
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress the data
func decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
