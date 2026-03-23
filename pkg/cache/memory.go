package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type memoryEntry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

// MemoryCache is an in-memory cache with LRU eviction and TTL expiration.
type MemoryCache struct {
	mu       sync.RWMutex
	items    map[string]*list.Element
	order    *list.List // front = most recently used
	maxItems int
	stop     chan struct{}
}

// NewMemoryCache creates an in-memory cache with the given capacity.
// A background goroutine cleans up expired entries at cleanupInterval.
// Call Close() to stop the cleanup goroutine.
func NewMemoryCache(maxEntries int, cleanupInterval time.Duration) *MemoryCache {
	if maxEntries <= 0 {
		maxEntries = 1000
	}
	if cleanupInterval <= 0 {
		cleanupInterval = time.Minute
	}

	c := &MemoryCache{
		items:    make(map[string]*list.Element),
		order:    list.New(),
		maxItems: maxEntries,
		stop:     make(chan struct{}),
	}

	go c.cleanup(cleanupInterval)
	return c
}

func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, bool, error) {
	c.mu.RLock()
	elem, ok := c.items[key]
	if !ok {
		c.mu.RUnlock()
		return nil, false, nil
	}
	entry := elem.Value.(*memoryEntry)
	if time.Now().After(entry.expiresAt) {
		c.mu.RUnlock()
		c.Delete(context.Background(), key)
		return nil, false, nil
	}
	// Copy value to avoid data races
	val := make([]byte, len(entry.value))
	copy(val, entry.value)
	c.mu.RUnlock()

	// Promote to front (requires write lock)
	c.mu.Lock()
	c.order.MoveToFront(elem)
	c.mu.Unlock()

	return val, true, nil
}

func (c *MemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		entry := elem.Value.(*memoryEntry)
		entry.value = make([]byte, len(value))
		copy(entry.value, value)
		entry.expiresAt = time.Now().Add(ttl)
		return nil
	}

	// Evict LRU if at capacity
	for len(c.items) >= c.maxItems {
		back := c.order.Back()
		if back == nil {
			break
		}
		evicted := c.order.Remove(back).(*memoryEntry)
		delete(c.items, evicted.key)
	}

	val := make([]byte, len(value))
	copy(val, value)
	entry := &memoryEntry{
		key:       key,
		value:     val,
		expiresAt: time.Now().Add(ttl),
	}
	elem := c.order.PushFront(entry)
	c.items[key] = elem
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.order.Remove(elem)
		delete(c.items, key)
	}
	return nil
}

func (c *MemoryCache) Flush(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*list.Element)
	c.order.Init()
	return nil
}

// Close stops the background cleanup goroutine.
func (c *MemoryCache) Close() error {
	close(c.stop)
	return nil
}

func (c *MemoryCache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.removeExpired()
		}
	}
}

func (c *MemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, elem := range c.items {
		entry := elem.Value.(*memoryEntry)
		if now.After(entry.expiresAt) {
			c.order.Remove(elem)
			delete(c.items, key)
		}
	}
}
