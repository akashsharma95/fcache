package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInmemoryCache(t *testing.T) {
	c := NewInmemoryCache()
	assert.NotNil(t, c)

	assert.NotNil(t, c.(*inMemoryCache).buckets)
	assert.NotNil(t, c.(*inMemoryCache).ttlJob)
}

func TestInMemoryCache_GetSet(t *testing.T) {
	c := NewInmemoryCache()

	err := c.Set("key", "value")
	assert.NoError(t, err)

	v, err := c.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)
}

func TestInMemoryCache_Delete(t *testing.T) {
	c := NewInmemoryCache()

	err := c.Set("key", "value")
	assert.NoError(t, err)

	c.Delete("key")

	v, err := c.Get("key")
	assert.ErrorIs(t, err, ErrorKeyNotFound)
	assert.Equal(t, "", v)
}

func TestInMemoryCache_SetWithTtl(t *testing.T) {
	c := NewInmemoryCache()

	err := c.SetWithTtl("key", "value", time.Minute*10)
	assert.NoError(t, err)
	bucket := c.(*inMemoryCache).getBucket("key")
	if !bucket.items["key"].expireAt.After(time.Now()) {
		t.Error("expireAt should be greater than current time")
	}
}

func TestInMemoryCache_Flush(t *testing.T) {
	c := NewInmemoryCache()

	err := c.Set("key", "value")
	assert.NoError(t, err)

	v, err := c.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	c.Flush()

	v, err = c.Get("key")
	assert.ErrorIs(t, err, ErrorKeyNotFound)
	assert.Equal(t, "", v)
}

func BenchmarkInMemoryCache_Set(b *testing.B) {
	const items = 1 << 16
	c := NewInmemoryCache()
	defer c.Flush()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				err := c.Set(string(k), v)
				if err != nil {
					b.Error(err)
				}
			}
		}
	})
}

func BenchmarkInMemoryCache_Get(b *testing.B) {
	const items = 1 << 16
	c := NewInmemoryCache()
	defer c.Flush()
	k := []byte("\x00\x00\x00\x00")
	v := "value"
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		err := c.Set(string(k), v)
		if err != nil {
			b.Error(err)
		}
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				got, err := c.Get(string(k))
				if err != nil {
					b.Error(err)
				}

				if got != v {
					panic(fmt.Errorf("BUG: invalid value obtained; got %s; want %q", got, v))
				}
			}
		}
	})
}

func BenchmarkInMemoryCache_SetGet(b *testing.B) {
	const items = 1 << 16
	c := NewInmemoryCache()
	defer c.Flush()
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				err := c.Set(string(k), v)
				if err != nil {
					b.Error(err)
				}
			}
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				got, err := c.Get(string(k))
				if err != nil {
					b.Error(err)
				}

				if got != v {
					panic(fmt.Errorf("BUG: invalid value obtained; got %s; want %q", got, v))
				}
			}
		}
	})
}

func BenchmarkStdMap_Set(b *testing.B) {
	const items = 1 << 16
	m := make(map[string]string)
	var mu sync.Mutex
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				mu.Lock()
				m[string(k)] = v
				mu.Unlock()
			}
		}
	})
}

func BenchmarkStdMap_Get(b *testing.B) {
	const items = 1 << 16
	m := make(map[string]string)
	k := []byte("\x00\x00\x00\x00")
	v := "value"
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		m[string(k)] = v
	}

	var mu sync.RWMutex
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				mu.RLock()
				vv := m[string(k)]
				mu.RUnlock()
				if vv != v {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkStdMap_SetGet(b *testing.B) {
	const items = 1 << 16
	m := make(map[string]string)
	var mu sync.RWMutex
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				mu.Lock()
				m[string(k)] = v
				mu.Unlock()
			}
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				mu.RLock()
				vv := m[string(k)]
				mu.RUnlock()
				if vv != v {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkSyncMap_Set(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				m.Store(string(k), v)
			}
		}
	})
}

func BenchmarkSyncMap_Get(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	k := []byte("\x00\x00\x00\x00")
	v := "value"
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		m.Store(string(k), v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				vv, ok := m.Load(string(k))
				if !ok || vv.(string) != v {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}

func BenchmarkSyncMap_SetGet(b *testing.B) {
	const items = 1 << 16
	m := sync.Map{}
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := "value"
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				m.Store(string(k), v)
			}
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				vv, ok := m.Load(string(k))
				if !ok || vv.(string) != v {
					panic(fmt.Errorf("BUG: unexpected value; got %q; want %q", vv, v))
				}
			}
		}
	})
}
