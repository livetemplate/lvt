// Package cache provides caching backends for generated applications.
// It defines the Cache interface and provides MemoryCache (LRU+TTL)
// and SQLiteCache (persistent) implementations.
package cache

import (
	"context"
	"encoding/json"
	"time"
)

// Cache is the interface for cache backends.
// Get returns (value, found, error). A nil error with found=false means cache miss.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}

// Closeable is optionally implemented by caches that need cleanup (e.g., stop background goroutines).
type Closeable interface {
	Close() error
}

// GetJSON retrieves a value from the cache and unmarshals it into the target type.
func GetJSON[T any](c Cache, ctx context.Context, key string) (T, bool, error) {
	var zero T
	data, found, err := c.Get(ctx, key)
	if err != nil || !found {
		return zero, found, err
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, false, err
	}
	return result, true, nil
}

// SetJSON marshals the value to JSON and stores it in the cache.
func SetJSON(c Cache, ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, data, ttl)
}
