package cache

import (
	"time"
)

const (
	defaultGcDuration = time.Minute * 1
)

// ttl is a background worker which deletes the expired key at some interval
type ttl struct {
	cache    *inMemoryCache
	tick     *time.Ticker
	shutdown chan struct{}
}

type TtlOption func(*ttl)

// newTtlJob create new ttlJob for a cache
func newTtlJob(cache *inMemoryCache, opts ...TtlOption) *ttl {
	ttlJob := &ttl{
		cache:    cache,
		tick:     time.NewTicker(defaultGcDuration),
		shutdown: make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(ttlJob)
	}

	return ttlJob
}

func WithGcDuration(duration time.Duration) TtlOption {
	return func(t *ttl) {
		if t.tick != nil {
			t.tick.Reset(duration)
		} else {
			t.tick = time.NewTicker(duration)
		}
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
