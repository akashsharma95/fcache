package cache

import (
	"time"
)

const (
	gcDuration = time.Minute * 1
)

// ttl is a background worker which deletes the expired key at some interval
type ttl struct {
	cache    *inMemoryCache
	tick     *time.Ticker
	shutdown chan struct{}
}

// newTtlJob create new ttlJob for a cache
func newTtlJob(cache *inMemoryCache) *ttl {
	return &ttl{
		cache:    cache,
		tick:     time.NewTicker(gcDuration),
		shutdown: make(chan struct{}, 1),
	}
}

// start the ttl job in separate go-routine
func (t *ttl) start() {
	go func() {
		for {
			select {
			case <-t.shutdown:
				t.tick.Stop()
				return

			case <-t.tick.C:
				// iterate over all the buckets and remove the expired keys
				// can be optimized using priority queue with increased code complexity
				var keys []string
				for _, bucket := range t.cache.buckets {
					bucket.RLock()
					for k, v := range bucket.items {
						if v.isExpired() {
							keys = append(keys, k)
						}
					}
					bucket.RUnlock()
				}
				t.cache.deleteKeys(keys...)
			}
		}
	}()
}
