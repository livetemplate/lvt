package storage

import (
	"context"
	"fmt"
	"io"
	"io/fs"
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
// The key is sanitized to prevent path traversal outside baseDir.
func (s *LocalStore) Save(_ context.Context, key string, reader io.Reader) error {
	path, err := s.safePath(key)
	if err != nil {
		return err
	}

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
	path, err := s.safePath(key)
	if err != nil {
		return nil, err
	}
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
	path, err := s.safePath(key)
	if err != nil {
		return err
	}
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

// safePath resolves a key to an absolute path within baseDir,
// rejecting any key that would escape via path traversal.
func (s *LocalStore) safePath(key string) (string, error) {
	absBase, err := filepath.Abs(s.baseDir)
	if err != nil {
		return "", fmt.Errorf("storage: resolve base dir: %w", err)
	}
	joined := filepath.Join(absBase, filepath.FromSlash(key))
	cleaned := filepath.Clean(joined)
	if !strings.HasPrefix(cleaned, absBase+string(filepath.Separator)) && cleaned != absBase {
		return "", fmt.Errorf("storage: path traversal rejected: %q", key)
	}
	return cleaned, nil
}

// URL returns the public URL for the given key.
func (s *LocalStore) URL(key string) string {
	return s.urlPrefix + "/" + key
}

// FileServer returns an http.Handler that serves files from baseDir.
// Directory listings are disabled — only individual files are served.
func (s *LocalStore) FileServer() http.Handler {
	return http.FileServer(noListDir{http.Dir(s.baseDir)})
}

// noListDir wraps an http.FileSystem to disable directory listings.
type noListDir struct{ http.FileSystem }

func (d noListDir) Open(name string) (http.File, error) {
	f, err := d.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	if stat, _ := f.Stat(); stat != nil && stat.IsDir() {
		f.Close()
		return nil, fs.ErrPermission
	}
	return f, nil
}
