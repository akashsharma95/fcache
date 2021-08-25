package cache

import "time"

type item struct {
	value    string
	expireAt time.Time
}

func newItem(value string, expireAt time.Time) *item {
	return &item{
		value:    value,
		expireAt: expireAt,
	}
}

func (i *item) getValue() string {
	return i.value
}

func (i *item) isExpired() bool {
	return i.expireAt.Before(time.Now())
}
