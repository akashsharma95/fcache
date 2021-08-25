package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInmemoryCache(t *testing.T) {
	c := NewInmemoryCache()
	assert.NotNil(t, c)

	assert.NotNil(t, c.(inMemoryCache).buckets)
	assert.NotNil(t, c.(inMemoryCache).ttl)
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
	bucket := c.(inMemoryCache).getBucket("key")
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

}

func BenchmarkInMemoryCache_Get(b *testing.B) {

}
