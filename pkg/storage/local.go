package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// LocalStore stores files on the local filesystem.
type LocalStore struct {
	baseDir   string
	urlPrefix string
}

// NewLocalStore creates a new LocalStore that saves files under baseDir
// and generates URLs with the given prefix.
func NewLocalStore(baseDir, urlPrefix string) *LocalStore {
	return &LocalStore{
		baseDir:   baseDir,
		urlPrefix: strings.TrimRight(urlPrefix, "/"),
	}
}

// Save writes the contents of reader to baseDir/key, creating parent directories.
func (s *LocalStore) Save(_ context.Context, key string, reader io.Reader) error {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("storage: create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("storage: create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("storage: write file: %w", err)
	}

	return nil
}

// Open returns a ReadCloser for the file at baseDir/key.
func (s *LocalStore) Open(_ context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("storage: open file: %w", err)
	}
	return f, nil
}

// Delete removes the file at baseDir/key. If key is a URL returned by URL(),
// the prefix is stripped automatically.
func (s *LocalStore) Delete(_ context.Context, key string) error {
	key = s.resolveKey(key)
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: delete file: %w", err)
	}
	return nil
}

// resolveKey strips the URL prefix if the key looks like a URL returned by URL().
func (s *LocalStore) resolveKey(key string) string {
	prefix := s.urlPrefix + "/"
	if strings.HasPrefix(key, prefix) {
		return key[len(prefix):]
	}
	return key
}

// URL returns the public URL for the given key.
func (s *LocalStore) URL(key string) string {
	return s.urlPrefix + "/" + key
}

// FileServer returns an http.Handler that serves files from baseDir.
func (s *LocalStore) FileServer() http.Handler {
	return http.FileServer(http.Dir(s.baseDir))
}
