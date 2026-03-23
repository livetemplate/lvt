package cache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryCache_SetGet(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	if err := c.Set(ctx, "key1", []byte("value1"), time.Minute); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	val, found, err := c.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !found {
		t.Fatal("Get() found = false, want true")
	}
	if string(val) != "value1" {
		t.Errorf("Get() = %q, want %q", val, "value1")
	}
}

func TestMemoryCache_Miss(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()

	_, found, err := c.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if found {
		t.Error("Get() found = true for nonexistent key")
	}
}

func TestMemoryCache_TTLExpiration(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	c.Set(ctx, "expires", []byte("data"), 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, found, _ := c.Get(ctx, "expires")
	if found {
		t.Error("Get() found expired entry")
	}
}

func TestMemoryCache_LRUEviction(t *testing.T) {
	c := NewMemoryCache(3, time.Minute)
	defer c.Close()
	ctx := context.Background()

	c.Set(ctx, "a", []byte("1"), time.Minute)
	c.Set(ctx, "b", []byte("2"), time.Minute)
	c.Set(ctx, "c", []byte("3"), time.Minute)

	// Access "a" to make it most recently used
	c.Get(ctx, "a")

	// Add "d" — should evict "b" (least recently used)
	c.Set(ctx, "d", []byte("4"), time.Minute)

	_, found, _ := c.Get(ctx, "b")
	if found {
		t.Error("LRU should have evicted 'b'")
	}

	_, found, _ = c.Get(ctx, "a")
	if !found {
		t.Error("'a' should still be in cache (was accessed recently)")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	c.Set(ctx, "key", []byte("val"), time.Minute)
	c.Delete(ctx, "key")

	_, found, _ := c.Get(ctx, "key")
	if found {
		t.Error("Get() found deleted key")
	}
}

func TestMemoryCache_Flush(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	c.Set(ctx, "a", []byte("1"), time.Minute)
	c.Set(ctx, "b", []byte("2"), time.Minute)
	c.Flush(ctx)

	_, foundA, _ := c.Get(ctx, "a")
	_, foundB, _ := c.Get(ctx, "b")
	if foundA || foundB {
		t.Error("Flush() should clear all entries")
	}
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	c.Set(ctx, "key", []byte("old"), time.Minute)
	c.Set(ctx, "key", []byte("new"), time.Minute)

	val, _, _ := c.Get(ctx, "key")
	if string(val) != "new" {
		t.Errorf("Get() = %q after overwrite, want %q", val, "new")
	}
}

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	c := NewMemoryCache(1000, time.Minute)
	defer c.Close()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			c.Set(ctx, key, []byte("value"), time.Minute)
			c.Get(ctx, key)
		}(i)
	}
	wg.Wait()
}
