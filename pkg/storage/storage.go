// Package storage provides file storage backends for generated applications.
// It defines the Store interface and provides LocalStore (filesystem) and
// S3Store (AWS S3) implementations.
package storage

import (
	"context"
	"io"
)

// Store is the interface for file storage backends.
type Store interface {
	Save(ctx context.Context, key string, reader io.Reader) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
}
