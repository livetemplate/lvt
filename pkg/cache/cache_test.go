package cache

import (
	"context"
	"testing"
	"time"
)

type testPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func TestGetSetJSON(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()
	ctx := context.Background()

	post := testPost{Title: "Hello", Content: "World"}
	if err := SetJSON(c, ctx, "post:1", post, time.Minute); err != nil {
		t.Fatalf("SetJSON() error = %v", err)
	}

	result, found, err := GetJSON[testPost](c, ctx, "post:1")
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if !found {
		t.Fatal("GetJSON() found = false")
	}
	if result.Title != "Hello" || result.Content != "World" {
		t.Errorf("GetJSON() = %+v, want {Title:Hello Content:World}", result)
	}
}

func TestGetJSON_Miss(t *testing.T) {
	c := NewMemoryCache(100, time.Minute)
	defer c.Close()

	_, found, err := GetJSON[testPost](c, context.Background(), "missing")
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if found {
		t.Error("GetJSON() found = true for missing key")
	}
}
