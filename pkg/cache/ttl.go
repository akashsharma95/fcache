package cache

import (
	"log"
	"time"
)

const (
	gcDuration       = time.Minute * 1
)

type ttl struct {
	cache        *inMemoryCache
	tick         *time.Ticker
	shutdown     chan struct{}
}

func newTtl(cache *inMemoryCache) *ttl {
	return &ttl{
		cache:      cache,
		tick:       time.NewTicker(gcDuration),
		shutdown:   make(chan struct{}, 1),
	}
}

func (t *ttl) startCleanUpJob() {
	log.Println("starting cleanup job")
	go func() {
		for {
			select {
			case <-t.shutdown:
				t.tick.Stop()
				log.Println("cleanup job shutdown")
				return
			case <-t.tick.C:
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
