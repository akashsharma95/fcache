package cache

import (
	"fmt"
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

func BenchmarkCacheSet(b *testing.B) {
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

func BenchmarkCacheGet(b *testing.B) {
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

func BenchmarkCacheSetGet(b *testing.B) {
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
