package internal

import (
	"time"
	"sync"
	"fmt"
)

type cache struct {
	cacheEntries 	map[string]cacheEntry
	mu 				sync.Mutex
	interval 		time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func (c cache) Add(key string, val []byte) {
	fmt.Printf("cache.Add() key: %s\n", key)
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cacheEntries[key]
	if ok {
		fmt.Printf("cache.Add(%s) overwriting existing key's value: %w\n", key, val)
	} 
	entry = cacheEntry{val: val, createdAt: time.Now()}
	c.cacheEntries[key] = entry
}

func (c cache) Get(key string) ([]byte, bool) {
	fmt.Printf("cache.Get(%s)\n", key)
	entry, ok := c.cacheEntries[key]
	if ok {
		fmt.Printf("cache hit for key %s\n", key)
		return entry.val, ok
	} 
	fmt.Printf("Cache miss on key %s\n", key)
	return []byte{}, ok
}


func NewCache(interval time.Duration) *cache {
	eMap := make(map[string]cacheEntry)
	mutex := sync.Mutex{}
	c := cache{cacheEntries: eMap, mu: mutex, interval: interval}
	go c.reapLoop()
	return &c
}

func (c cache) reapLoop() {
	ticker := time.Tick(c.interval)
	for range ticker {
		for key, val := range c.cacheEntries {
			if time.Since(val.createdAt) > c.interval {
				fmt.Printf("Deleting key %s\n", key)
				delete(c.cacheEntries, key)
			}
		}
	}
}
