package cache

import "time"

// Cache provides caching interface
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Clear() error
}
