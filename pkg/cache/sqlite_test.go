package cache

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func newTestSQLiteCache(t *testing.T) *SQLiteCache {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	c, err := NewSQLiteCache(db, time.Minute)
	if err != nil {
		t.Fatalf("NewSQLiteCache() error = %v", err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func TestSQLiteCache_SetGet(t *testing.T) {
	c := newTestSQLiteCache(t)
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

func TestSQLiteCache_Miss(t *testing.T) {
	c := newTestSQLiteCache(t)

	_, found, err := c.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if found {
		t.Error("Get() found = true for nonexistent key")
	}
}

func TestSQLiteCache_TTLExpiration(t *testing.T) {
	c := newTestSQLiteCache(t)
	ctx := context.Background()

	c.Set(ctx, "expires", []byte("data"), 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	_, found, _ := c.Get(ctx, "expires")
	if found {
		t.Error("Get() found expired entry")
	}
}

func TestSQLiteCache_Delete(t *testing.T) {
	c := newTestSQLiteCache(t)
	ctx := context.Background()

	c.Set(ctx, "key", []byte("val"), time.Minute)
	c.Delete(ctx, "key")

	_, found, _ := c.Get(ctx, "key")
	if found {
		t.Error("Get() found deleted key")
	}
}

func TestSQLiteCache_Flush(t *testing.T) {
	c := newTestSQLiteCache(t)
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

func TestSQLiteCache_Overwrite(t *testing.T) {
	c := newTestSQLiteCache(t)
	ctx := context.Background()

	c.Set(ctx, "key", []byte("old"), time.Minute)
	c.Set(ctx, "key", []byte("new"), time.Minute)

	val, _, _ := c.Get(ctx, "key")
	if string(val) != "new" {
		t.Errorf("Get() = %q after overwrite, want %q", val, "new")
	}
}
