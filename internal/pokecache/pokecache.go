// Package pokecache implements the caching of the api calls from pokeapi
package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]cacheEntry
	Mu       sync.Mutex
	Interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries:  make(map[string]cacheEntry),
		Interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Find(key string) (val []byte, isCached bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	e, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	return e.val, ok
}

func (c *Cache) reapLoop() {
	t := time.NewTicker(c.Interval)
	defer t.Stop()
	for range t.C {
		c.Mu.Lock()
		for k, v := range c.Entries {
			if c.Interval < time.Since(v.createdAt) {
				delete(c.Entries, k)
			}
		}
		c.Mu.Unlock()
	}
}
