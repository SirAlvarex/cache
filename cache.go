package cache

import (
	"errors"
	"sync"
)

// Cache is an object that can be used for caching key:value pairs
type Cache struct {
	cache       map[string]interface{}
	cacheRWLock sync.RWMutex
	ignoreCache bool
}

// New creates and returns a new Cache
func New() *Cache {
	c := Cache{
		cache:       make(map[string]interface{}),
		cacheRWLock: sync.RWMutex{},
	}
	return &c
}

// Disable will temporarily disable cache lookups.  Useful for threading.
func (v *Cache) Disable() {
	v.ignoreCache = true
}

// Enable will  enable cache lookups if previously disabled
func (v *Cache) Enable() {
	v.ignoreCache = false
}

// Check checks the cache for the cacheKey value and returns it.  If not found, then error is not nil
func (v *Cache) Check(cacheKey string) (cacheValue interface{}, err error) {
	if v.ignoreCache {
		return cacheValue, errors.New("cache is disabled")
	}
	v.cacheRWLock.RLock()
	defer v.cacheRWLock.RUnlock()
	if value, ok := v.cache[cacheKey]; ok {
		return value, nil
	}
	return cacheValue, errors.New("not found")
}

// Update updates the VSTS Cache
func (v *Cache) Update(cacheKey string, cacheValue interface{}) {
	v.cacheRWLock.Lock()
	defer v.cacheRWLock.Unlock()
	v.cache[cacheKey] = cacheValue

	return
}

// Invalidate invalidates a key in the VSTS Cache
func (v *Cache) Invalidate(cacheKey string) {
	v.cacheRWLock.Lock()
	defer v.cacheRWLock.Unlock()
	delete(v.cache, cacheKey)
}

// Purge invalidates a key in the VSTS Cache
func (v *Cache) Purge() {
	v.cacheRWLock.Lock()
	defer v.cacheRWLock.Unlock()
	v.cache = make(map[string]interface{})
}
