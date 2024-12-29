package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data      map[string][]byte
	timestamp time.Time
	duration  time.Duration
	mu        sync.RWMutex
}

func NewCache(duration time.Duration) *Cache {
	return &Cache{
		data:     make(map[string][]byte),
		duration: duration,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Since(c.timestamp) > c.duration {
		return nil, false
	}

	return data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = data
	c.timestamp = time.Now()
}
