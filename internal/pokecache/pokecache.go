// Package pokecache implements the caching of the api calls from pokeapi
package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry
	Mutex   *sync.Mutex
}
type cacheEntry struct {
	createAt time.Time
	val      []byte
}

func (c *Cache) NewCache(interval time.Duration) (cache *Cache) {
	ticker := time.NewTicker(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.Entries[key] = cacheEntry{createAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) (val []byte, isCached bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	entry, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, ok
}

func (c *Cache) reapLoop() {
}
