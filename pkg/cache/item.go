package cache

import "time"

// item contains value and expire-at time for value
type item struct {
	value    string
	expireAt time.Time
}

// newItem create new item with value and expire-at
func newItem(value string, expireAt time.Time) item {
	return item{
		value:    value,
		expireAt: expireAt,
	}
}

// getValue get the value
func (i item) getValue() string {
	return i.value
}

// isExpired returns true if the item is expired based on expireAt value
func (i item) isExpired() bool {
	return i.expireAt.Before(time.Now())
}
