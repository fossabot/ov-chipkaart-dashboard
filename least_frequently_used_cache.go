package main

import "github.com/pkg/errors"

// MapLFUCache is the data structure for a least frequently used cache implemented using a map
type MapLFUCache struct {
	hashMap map[interface{}]interface{}
}

var (
	// ErrorCacheMiss is the error when a cache miss occurs
	ErrorCacheMiss = errors.New("cache miss")
)

// NewLFUStringCache creates a new cache with capacity
func NewLFUStringCache(capacity int) LFUCache {
	return &MapLFUCache{hashMap: make(map[interface{}]interface{}, capacity)}
}

// Get returns an item from the cache using the key
// TODO: Implement the frequency counter
func (cache MapLFUCache) Get(key interface{}) (val interface{}, err error) {
	val, ok := cache.hashMap[key]
	if ok {
		return val, nil
	}

	return val, ErrorCacheMiss
}

// Put gets an item from the cache
// TODO: Implement the cache invalidation logic
func (cache *MapLFUCache) Put(key interface{}, value interface{}) {
	cache.hashMap[key] = value
}
