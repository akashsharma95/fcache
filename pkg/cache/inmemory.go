package cache

import (
	"runtime"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

// defaultTTL default expiry time for the cache
const defaultTTL = time.Minute * 30

// bucketCount is number of buckets/shards to create for cache is based on number of CPU * factor 8
// factor 8 was chosen by running benchmarks
var bucketCount = uint64(runtime.NumCPU() * 8)

// inMemoryCache stores key value pair in sharded map
type inMemoryCache struct {
	buckets map[uint64]*cacheBucket
	ttlJob  *ttl
}

// cacheBucket is concurrent safe map
type cacheBucket struct {
	sync.RWMutex
	items map[string]item
}

// NewInmemoryCache creates new in memory cache instance with fixed number of buckets
func NewInmemoryCache() Cache {
	cache := &inMemoryCache{
		buckets: make(map[uint64]*cacheBucket, bucketCount),
	}

	for idx := uint64(0); idx < bucketCount; idx++ {
		cache.buckets[idx] = &cacheBucket{
			items: make(map[string]item),
		}
	}

	// based on different cache eviction policy we can have background workers for cleaning cache
	cache.ttlJob = newTtlJob(cache)
	cache.ttlJob.start()

	return cache
}

// Get gets the key from one of the bucket and returns error if key is not found
func (c *inMemoryCache) Get(key string) (string, error) {
	bucket := c.getBucket(key)
	bucket.RLock()
	defer bucket.RUnlock()

	if value, exists := bucket.items[key]; !exists {
		return "", ErrorKeyNotFound
	} else {
		if value.isExpired() {
			return "", ErrorKeyNotFound
		}

		return value.getValue(), nil
	}
}

// Set stores the key and value in cache with default ttl of 30 mins
func (c *inMemoryCache) Set(key string, value string) error {
	return c.SetWithTtl(key, value, defaultTTL)
}

// SetWithTtl stores the key and value in cache with given ttl value
func (c *inMemoryCache) SetWithTtl(key string, value string, duration time.Duration) error {
	bucket := c.getBucket(key)
	bucket.Lock()
	defer bucket.Unlock()

	bucket.items[key] = newItem(value, time.Now().Add(duration))

	return nil
}

// Delete removes a key from cache
func (c *inMemoryCache) Delete(key string) {
	c.deleteKeys(key)
}

// Flush clears the cache
func (c *inMemoryCache) Flush() {
	for _, bucket := range c.buckets {
		bucket.Lock()
		for k := range bucket.items {
			delete(bucket.items, k)
		}
		bucket.Unlock()
	}
}

// Teardown clears the cache and stops the ttl background job
func (c *inMemoryCache) Teardown() {
	c.Flush()
	c.ttlJob.shutdown <- struct{}{}
}

// deleteKeys helper function to delete the list of keys from cache
func (c *inMemoryCache) deleteKeys(keys ...string) {
	for _, k := range keys {
		bucket := c.getBucket(k)

		bucket.Lock()
		delete(bucket.items, k)
		bucket.Unlock()
	}
}

// getBucket get the bucket where key is stored using consistent hashing
func (c *inMemoryCache) getBucket(key string) *cacheBucket {
	hash := xxhash.Sum64String(key)
	bucketKey := hash % bucketCount
	return c.buckets[bucketKey]
}
