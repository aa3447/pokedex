package http
// This file provides a cache implementation for the Pokedex application.

import (
	"time"
)	

type Cache struct {
	CacheData map[string]cacheEntry
}

type cacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

func NewCache(reapInterval time.Duration) *Cache {
	ticker := time.NewTicker(reapInterval)
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

func (c *Cache) Add(key string, val []byte){
	c.CacheData[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, exists := c.CacheData[key]
	if !exists {
		return nil, false
	}
	return entry.Val, true
}

func (c *Cache) reapLoop(time time.Time, interval time.Duration){
	for key, entry := range c.CacheData{
		if entry.CreatedAt.Add(interval).Before(time) {
			delete(c.CacheData, key)
		}
	}
}