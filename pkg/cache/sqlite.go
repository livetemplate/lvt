package cache

import (
	"context"
	"database/sql"
	"time"
)

// SQLiteCache is a cache backed by a SQLite database table.
// Data persists across application restarts.
type SQLiteCache struct {
	db   *sql.DB
	stop chan struct{}
}

// NewSQLiteCache creates a persistent cache using the given SQLite database.
// It creates the _cache table if it doesn't exist.
// A background goroutine cleans up expired entries at cleanupInterval.
func NewSQLiteCache(db *sql.DB, cleanupInterval time.Duration) (*SQLiteCache, error) {
	if cleanupInterval <= 0 {
		cleanupInterval = 5 * time.Minute
	}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS _cache (
		key TEXT PRIMARY KEY,
		value BLOB NOT NULL,
		expires_at DATETIME NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	c := &SQLiteCache{db: db, stop: make(chan struct{})}
	go c.cleanup(cleanupInterval)
	return c, nil
}

func (c *SQLiteCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	var value []byte
	err := c.db.QueryRowContext(ctx,
		`SELECT value FROM _cache WHERE key = ? AND expires_at > ?`, key, time.Now(),
	).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return value, true, nil
}

func (c *SQLiteCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	_, err := c.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO _cache (key, value, expires_at) VALUES (?, ?, ?)`,
		key, value, expiresAt,
	)
	return err
}

func (c *SQLiteCache) Delete(ctx context.Context, key string) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM _cache WHERE key = ?`, key)
	return err
}

func (c *SQLiteCache) Flush(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM _cache`)
	return err
}

// Close stops the background cleanup goroutine.
func (c *SQLiteCache) Close() error {
	close(c.stop)
	return nil
}

func (c *SQLiteCache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.db.Exec(`DELETE FROM _cache WHERE expires_at < ?`, time.Now())
		}
	}
}
