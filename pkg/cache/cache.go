package cache

import (
	"errors"
	"time"
)

var ErrorKeyNotFound = errors.New("key not found")

type Cache interface {
	Get(string) (string, error)
	Set(string, string) error
	SetWithTtl(string, string, time.Duration) error
	Delete(string)
	Flush()
}
