package cache

import (
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrCacheMiss is thrown when there's a cache miss
	ErrCacheMiss = errors.New("cache miss")
)

// Cache is the interface for a cache
type Cache interface {
	Set(key, value string, expiration time.Duration) error
	Get(key string) (string, error)
}
