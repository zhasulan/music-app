package cache

import (
	"sync"
	"time"
)

type entry struct {
	value      any
	expiration time.Time
}

type InMemory struct {
	data map[string]entry
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewInMemory(ttl time.Duration) *InMemory {
	return &InMemory{
		data: make(map[string]entry),
		ttl:  ttl,
	}
}

func (c *InMemory) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[key]
	if !ok || time.Now().After(e.expiration) {
		return nil, false
	}
	return e.value, true
}

func (c *InMemory) Set(key string, val any, ttlOverride ...time.Duration) {
	ttl := c.ttl
	if len(ttlOverride) > 0 {
		ttl = ttlOverride[0]
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entry{value: val, expiration: time.Now().Add(ttl)}
}
