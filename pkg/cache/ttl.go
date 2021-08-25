package cache

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxExpiredBuffer = 64
	gcDuration       = time.Minute * 1
)

type ttl struct {
	cache        *inMemoryCache
	expiredCount uint32
	expiredMap   sync.Map
	tick         *time.Ticker
	shutdown     chan struct{}
}

func newTtl(cache *inMemoryCache) *ttl {
	return &ttl{
		cache:      cache,
		expiredMap: sync.Map{},
		tick:       time.NewTicker(gcDuration),
		shutdown:   make(chan struct{}, 1),
	}
}

func (t *ttl) addExpiredKey(key string) {
	_, exists := t.expiredMap.LoadOrStore(key, nil)
	if !exists {
		atomic.AddUint32(&t.expiredCount, 1)
	}
}

func (t *ttl) startCleanUpJob() {
	log.Println("starting cleanup job")
	go func() {
		for {
			select {
			case <-t.shutdown:
				t.tick.Stop()
				return
			case <-t.tick.C:
				expiredLen := t.expiredCount
				if expiredLen > maxExpiredBuffer/2 {
					var keys []string
					t.expiredMap.Range(func(key, _ interface{}) bool {
						keys = append(keys, key.(string))
						return true
					})
					t.cache.deleteKeys(keys...)
					atomic.AddUint32(&t.expiredCount, -uint32(len(keys)))
				}

				for _, bucket := range t.cache.buckets {
					bucket.Lock()
					for k, v := range bucket.items {
						if v.isExpired() {
							t.addExpiredKey(k)
						}
					}
					bucket.Unlock()
				}
			}
		}
	}()
}
