package cache

import (
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

type LRUExpireCache struct {
	cache *lru.Cache
	lock  sync.RWMutex
}

func NewLRUExpireCache(maxSize int) *LRUExpireCache {
	return &LRUExpireCache{cache: lru.New(maxSize)}
}

type cacheEntry struct {
	value      interface{}
	expireTime time.Time
}

//Add  ttl==0 for forever
func (c *LRUExpireCache) Add(key lru.Key, value interface{}, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Add(key, &cacheEntry{value, time.Now().Add(ttl)})
	// Remove entry from cache after ttl.
	if ttl != 0 {
		time.AfterFunc(ttl, func() { c.Remove(key) })
	}
}

func (c *LRUExpireCache) Get(key lru.Key) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	e, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	if time.Now().After(e.(*cacheEntry).expireTime) {
		go c.Remove(key)
		return nil, false
	}
	return e.(*cacheEntry).value, true
}

func (c *LRUExpireCache) Remove(key lru.Key) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Remove(key)
}
