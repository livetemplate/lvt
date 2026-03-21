package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStore_Save(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir, "/uploads")

	ctx := context.Background()
	content := []byte("hello world")

	if err := store.Save(ctx, "photos/test.txt", bytes.NewReader(content)); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists on disk
	data, err := os.ReadFile(filepath.Join(dir, "photos", "test.txt"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("got %q, want %q", string(data), "hello world")
	}
}

func TestLocalStore_Open(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir, "/uploads")

	// Write a file first
	os.MkdirAll(filepath.Join(dir, "docs"), 0755)
	os.WriteFile(filepath.Join(dir, "docs", "readme.txt"), []byte("read me"), 0644)

	ctx := context.Background()
	rc, err := store.Open(ctx, "docs/readme.txt")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(data) != "read me" {
		t.Errorf("got %q, want %q", string(data), "read me")
	}
}

func TestLocalStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir, "/uploads")

	// Write a file first
	path := filepath.Join(dir, "temp.txt")
	os.WriteFile(path, []byte("temp"), 0644)

	ctx := context.Background()
	if err := store.Delete(ctx, "temp.txt"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestLocalStore_Delete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir, "/uploads")

	ctx := context.Background()
	if err := store.Delete(ctx, "does-not-exist.txt"); err != nil {
		t.Fatalf("Delete() non-existent should not error, got %v", err)
	}
}

func TestLocalStore_URL(t *testing.T) {
	store := NewLocalStore("/var/uploads", "/uploads")

	url := store.URL("photos/avatar.jpg")
	if url != "/uploads/photos/avatar.jpg" {
		t.Errorf("URL() = %q, want %q", url, "/uploads/photos/avatar.jpg")
	}
}

func TestLocalStore_URL_TrailingSlash(t *testing.T) {
	store := NewLocalStore("/var/uploads", "/uploads/")

	url := store.URL("file.txt")
	if url != "/uploads/file.txt" {
		t.Errorf("URL() = %q, want %q", url, "/uploads/file.txt")
	}
}

func TestLocalStore_SaveAndOpen_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir, "/uploads")
	ctx := context.Background()

	content := []byte("round trip content")
	if err := store.Save(ctx, "rt/file.bin", bytes.NewReader(content)); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	rc, err := store.Open(ctx, "rt/file.bin")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("round trip mismatch: got %q, want %q", data, content)
	}
}

func TestS3Store_URL(t *testing.T) {
	tests := []struct {
		name   string
		config S3StoreConfig
		key    string
		want   string
	}{
		{
			name:   "standard S3 URL",
			config: S3StoreConfig{Bucket: "my-bucket", Region: "us-east-1"},
			key:    "photos/avatar.jpg",
			want:   "https://my-bucket.s3.us-east-1.amazonaws.com/photos/avatar.jpg",
		},
		{
			name:   "CDN prefix",
			config: S3StoreConfig{Bucket: "my-bucket", Region: "us-east-1", CDNPrefix: "https://cdn.example.com"},
			key:    "photos/avatar.jpg",
			want:   "https://cdn.example.com/photos/avatar.jpg",
		},
		{
			name:   "custom endpoint",
			config: S3StoreConfig{Bucket: "my-bucket", Region: "us-east-1", Endpoint: "http://localhost:9000"},
			key:    "photos/avatar.jpg",
			want:   "http://localhost:9000/my-bucket/photos/avatar.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &S3Store{config: tt.config}
			got := store.URL(tt.key)
			if got != tt.want {
				t.Errorf("URL() = %q, want %q", got, tt.want)
			}
		})
	}
}
