// Package storage provides file storage backends for generated applications.
// It defines the Store interface and provides LocalStore (filesystem) and
// S3Store (AWS S3) implementations.
package storage

import (
	"context"
	"io"
)

// Store is the interface for file storage backends.
// Keys are opaque storage paths (e.g., "galleries/id/photo.jpg").
// Delete also accepts URLs returned by URL() — implementations strip
// their known prefix to recover the underlying key.
type Store interface {
	Save(ctx context.Context, key string, reader io.Reader) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
}
