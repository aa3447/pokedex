package http

import (
	"testing"
	"time"
)
 
func TestNewCache(t *testing.T) {
	cache := NewCache(24 * time.Hour)
	if cache == nil {
		t.Error("Expected cache to be created, got nil")
		return
	}
	if len(cache.CacheData) != 0{
		t.Error("Expected empty cache data on initialization")
		return
	}
}

func TestCacheAddAndGetAndExists(t *testing.T) {
	cache := NewCache(24 * time.Hour)
	key := "testKey"
	value := []byte("testValue")

	cache.Add(key, value)

	if !cache.Exists(key) {
		t.Errorf("Expected cache to contain key %s", key)
		return
	}

	retrievedValue, exists := cache.Get(key)
	if !exists {
		t.Errorf("Expected to retrieve value for key %s", key)
		return
	}
	if string(retrievedValue) != string(value) {
		t.Errorf("Expected value %s, got %s", value, retrievedValue)
		return
	}
}

func TestCacheReapLoop(t *testing.T) {
	cache := NewCache(1 * time.Second) // Short interval for testing
	key := "testKey"
	value := []byte("testValue")

	cache.Add(key, value)

	time.Sleep(2 * time.Second) // Wait for the reap loop to run

	if cache.Exists(key) {
		t.Errorf("Expected cache to have removed key %s after reap interval", key)
		return
	}
}