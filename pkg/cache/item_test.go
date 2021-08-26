package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItem(t *testing.T) {
	item := newItem("value", time.Now().Add(time.Minute))
	assert.NotNil(t, item)
	if !item.expireAt.After(time.Now()) {
		t.Error("expire should be greater than current time")
	}

	assert.Equal(t, "value", item.getValue())
}
