package telemetry

import (
	"context"
	"time"
)

// Store persists and retrieves generation events.
type Store interface {
	Save(ctx context.Context, event *GenerationEvent) error
	Get(ctx context.Context, id string) (*GenerationEvent, error)
	List(ctx context.Context, opts ListOptions) ([]*GenerationEvent, error)
	CountBySuccess(ctx context.Context, since time.Time) (total, successes int, err error)
	DeleteBefore(ctx context.Context, before time.Time) (int64, error)
	Close() error
}

// ListOptions controls filtering when listing events.
type ListOptions struct {
	Since       time.Time
	SuccessOnly *bool
	Command     string
	Limit       int
}
