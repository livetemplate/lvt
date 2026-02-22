package telemetry

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func openTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test-telemetry.db")
	store, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSQLiteStore_Schema(t *testing.T) {
	store := openTestStore(t)
	// Verify schema creation by attempting a query
	ctx := context.Background()
	events, err := store.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("list on fresh db: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestCollector_CaptureAndRetrieve(t *testing.T) {
	store := openTestStore(t)
	c := NewCollectorWithStore(store)
	defer c.Close()

	cap := c.StartCapture("gen resource", map[string]any{
		"resource_name": "posts",
		"fields":        []string{"title:string", "body:text"},
	})
	cap.SetKit("multi")
	cap.RecordFileGenerated("app/posts/posts.go")
	cap.RecordFileGenerated("app/posts/posts.tmpl")
	cap.RecordError(GenerationError{
		Phase:   "generation",
		Message: "test warning",
	})
	cap.Complete(true, `{"valid":true}`)

	// Retrieve
	ctx := context.Background()
	event, err := store.Get(ctx, cap.Event().ID)
	if err != nil {
		t.Fatalf("get event: %v", err)
	}
	if event.Command != "gen resource" {
		t.Errorf("expected command 'gen resource', got %q", event.Command)
	}
	if event.Kit != "multi" {
		t.Errorf("expected kit 'multi', got %q", event.Kit)
	}
	if !event.Success {
		t.Error("expected success=true")
	}
	if len(event.FilesGenerated) != 2 {
		t.Errorf("expected 2 files, got %d", len(event.FilesGenerated))
	}
	if len(event.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(event.Errors))
	}
	if event.ValidationJSON != `{"valid":true}` {
		t.Errorf("unexpected validation JSON: %q", event.ValidationJSON)
	}
	if event.DurationMs < 0 {
		t.Errorf("expected non-negative duration, got %d", event.DurationMs)
	}
	if event.Inputs["resource_name"] != "posts" {
		t.Errorf("expected input resource_name=posts, got %v", event.Inputs["resource_name"])
	}
}

func TestCollector_List(t *testing.T) {
	store := openTestStore(t)
	c := NewCollectorWithStore(store)
	defer c.Close()

	// Insert several events
	for i, cmd := range []string{"gen resource", "gen view", "gen resource"} {
		cap := c.StartCapture(cmd, nil)
		cap.Complete(i%2 == 0, "") // alternating success
	}

	ctx := context.Background()

	// List all
	all, err := store.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 events, got %d", len(all))
	}

	// Filter by command
	resources, err := store.List(ctx, ListOptions{Command: "gen resource"})
	if err != nil {
		t.Fatalf("list by command: %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("expected 2 'gen resource' events, got %d", len(resources))
	}

	// Filter by success
	successTrue := true
	successes, err := store.List(ctx, ListOptions{SuccessOnly: &successTrue})
	if err != nil {
		t.Fatalf("list by success: %v", err)
	}
	if len(successes) != 2 {
		t.Errorf("expected 2 successful events, got %d", len(successes))
	}

	// Limit
	limited, err := store.List(ctx, ListOptions{Limit: 1})
	if err != nil {
		t.Fatalf("list with limit: %v", err)
	}
	if len(limited) != 1 {
		t.Errorf("expected 1 event, got %d", len(limited))
	}
}

func TestCollector_Disabled(t *testing.T) {
	c := &Collector{enabled: false}

	cap := c.StartCapture("gen resource", nil)
	if cap.Event() != nil {
		t.Error("expected nil event for noop capture")
	}

	// These should all be safe to call on a noop capture
	cap.SetKit("multi")
	cap.RecordError(GenerationError{Phase: "generation", Message: "test"})
	cap.RecordFileGenerated("file.go")
	cap.Complete(true, "")

	if c.IsEnabled() {
		t.Error("expected collector to be disabled")
	}
}

func TestCollector_Retention(t *testing.T) {
	store := openTestStore(t)
	c := NewCollectorWithStore(store)
	defer c.Close()
	ctx := context.Background()

	// Insert an old event directly
	oldEvent := &GenerationEvent{
		ID:        "old-event",
		Timestamp: time.Now().AddDate(0, 0, -100), // 100 days ago
		Command:   "gen resource",
		Inputs:    map[string]any{},
		Success:   true,
	}
	if err := store.Save(ctx, oldEvent); err != nil {
		t.Fatalf("save old event: %v", err)
	}

	// Insert a recent event
	cap := c.StartCapture("gen view", nil)
	cap.Complete(true, "")

	// Verify both exist
	all, _ := store.List(ctx, ListOptions{})
	if len(all) != 2 {
		t.Fatalf("expected 2 events, got %d", len(all))
	}

	// Run retention with 90-day cutoff
	if err := c.RunRetention(ctx, 90); err != nil {
		t.Fatalf("run retention: %v", err)
	}

	// Only the recent event should remain
	remaining, _ := store.List(ctx, ListOptions{})
	if len(remaining) != 1 {
		t.Fatalf("expected 1 event after retention, got %d", len(remaining))
	}
	if remaining[0].Command != "gen view" {
		t.Errorf("expected remaining event to be 'gen view', got %q", remaining[0].Command)
	}
}

func TestSQLiteStore_CountBySuccess(t *testing.T) {
	store := openTestStore(t)
	c := NewCollectorWithStore(store)
	defer c.Close()
	ctx := context.Background()

	// Insert mix of success/failure
	for _, success := range []bool{true, true, false, true} {
		cap := c.StartCapture("gen resource", nil)
		cap.Complete(success, "")
	}

	since := time.Now().Add(-1 * time.Hour)
	total, successes, err := store.CountBySuccess(ctx, since)
	if err != nil {
		t.Fatalf("count by success: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total=4, got %d", total)
	}
	if successes != 3 {
		t.Errorf("expected successes=3, got %d", successes)
	}
}

func TestNoopCapture(t *testing.T) {
	cap := NoopCapture()
	cap.SetKit("multi")
	cap.RecordError(GenerationError{Phase: "test", Message: "test"})
	cap.RecordFileGenerated("file.go")
	cap.Complete(true, "")

	if cap.Event() != nil {
		t.Error("expected nil event for NoopCapture")
	}
}
