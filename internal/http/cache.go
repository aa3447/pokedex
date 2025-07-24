package http
// This file provides a cache implementation for the Pokedex application.

import (
	"time"
	"sync"
)	

type Cache struct {
	CacheData map[string]cacheEntry
}

type cacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

var mux *sync.RWMutex

func NewCache(reapInterval time.Duration) *Cache {
	ticker := time.NewTicker(reapInterval)
	mux = &sync.RWMutex{}
	cache := &Cache{
		CacheData: make(map[string]cacheEntry),
	}
	
	go func() {
		for  t := range ticker.C{
			cache.reapLoop(t, reapInterval)
		}	
	}()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	if c.Exists(key){
		return
	}
	
	mux.Lock()
	defer mux.Unlock()
	c.CacheData[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	mux.RLock()
	defer mux.RUnlock()
	
	entry, exists := c.CacheData[key]
	if !exists {
		return nil, false
	}
	return entry.Val, true
}

func (c *Cache) Exists(key string) bool {
	mux.RLock()
	defer mux.RUnlock()
	
	_, exists := c.CacheData[key]
	return exists
}

// reapLoop removes entries that are older than the reap interval.
func (c *Cache) reapLoop(callTime time.Time, reapInterval time.Duration) {
	mux.Lock()
	defer mux.Unlock()
	
	if c == nil || c.CacheData == nil || len(c.CacheData) == 0 {
		return
	}
	for key, entry := range c.CacheData{
		if entry.CreatedAt.Add(reapInterval).Before(callTime) {
			delete(c.CacheData, key)
		}
	}
}