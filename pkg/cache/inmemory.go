package cache

import (
	"runtime"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

const defaultTTL = time.Minute * 30

var bucketCount = uint64(runtime.NumCPU() * 4)

type inMemoryCache struct {
	buckets map[uint64]*cacheBucket
	ttl     *ttl
}

type cacheBucket struct {
	sync.RWMutex
	items map[string]*item
}

func NewInmemoryCache() Cache {
	cache := inMemoryCache{
		buckets: make(map[uint64]*cacheBucket, bucketCount),
	}

	for idx := uint64(0); idx < bucketCount; idx++ {
		cache.buckets[idx] = &cacheBucket{
			items: make(map[string]*item),
		}
	}

	t := newTtl(&cache)
	t.startCleanUpJob()
	cache.ttl = t

	return cache
}

func (c inMemoryCache) Get(key string) (string, error) {
	bucket := c.getBucket(key)
	bucket.RLock()
	defer bucket.RUnlock()

	if bucket.items[key] == nil {
		return "", ErrorKeyNotFound
	}

	if bucket.items[key].isExpired() {
		c.ttl.addExpiredKey(key)
		return "", ErrorKeyNotFound
	}

	return bucket.items[key].getValue(), nil
}

func (c inMemoryCache) Set(key string, value string) error {
	return c.SetWithTtl(key, value, defaultTTL)
}

func (c inMemoryCache) SetWithTtl(key string, value string, ttl time.Duration) error {
	bucket := c.getBucket(key)
	bucket.Lock()
	defer bucket.Unlock()

	bucket.items[key] = newItem(value, time.Now().Add(ttl))

	return nil
}

func (c inMemoryCache) Delete(key string) {
	c.deleteKeys(key)
}

func (c inMemoryCache) Flush() {
	for _, b := range c.buckets {
		b.Lock()
		for k, _ := range b.items {
			delete(b.items, k)
		}
		b.Unlock()
	}
}

func (c inMemoryCache) Teardown() {
	c.Flush()
	c.ttl.shutdown <- struct{}{}
}

func (c inMemoryCache) deleteKeys(keys ...string) {
	for _, k := range keys {
		bucket := c.getBucket(k)

		bucket.Lock()
		delete(bucket.items, k)
		bucket.Unlock()
	}
}

func (c inMemoryCache) getBucket(key string) *cacheBucket {
	hash := xxhash.Sum64([]byte(key))
	bucketKey := hash % bucketCount
	return c.buckets[bucketKey]
}
