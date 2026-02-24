package telemetry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/livetemplate/lvt/internal/config"
)

// Enabled returns true only when LVT_TELEMETRY is explicitly set to "true" or "1".
// Telemetry is opt-in: it is disabled by default to respect user privacy.
func Enabled() bool {
	v := strings.ToLower(os.Getenv("LVT_TELEMETRY"))
	return v == "true" || v == "1"
}

// Collector tracks generation events. Safe to use even when disabled (noop).
type Collector struct {
	store   Store
	enabled bool
}

// NewCollector opens the telemetry store. Returns a disabled collector if
// telemetry is off or the store cannot be opened (never returns an error to
// callers — telemetry must not break generation).
func NewCollector() (*Collector, error) {
	if !Enabled() {
		return &Collector{enabled: false}, nil
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return &Collector{enabled: false}, nil
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return &Collector{enabled: false}, nil
	}

	dbPath := filepath.Join(configDir, "telemetry.db")
	store, err := OpenSQLite(dbPath)
	if err != nil {
		return &Collector{enabled: false}, nil
	}

	return &Collector{store: store, enabled: true}, nil
}

// NewCollectorWithStore creates a collector with a custom store (for testing).
func NewCollectorWithStore(store Store) *Collector {
	return &Collector{store: store, enabled: true}
}

// Close releases the underlying store.
func (c *Collector) Close() error {
	if c.store != nil {
		return c.store.Close()
	}
	return nil
}

// Store returns the underlying Store (may be nil if disabled).
func (c *Collector) Store() Store {
	return c.store
}

// IsEnabled reports whether this collector is active.
func (c *Collector) IsEnabled() bool {
	return c.enabled
}

// StartCapture begins tracking a generation command. The returned Capture
// accumulates data until Complete() is called. If the collector is disabled
// the Capture is a silent noop.
func (c *Collector) StartCapture(command string, inputs map[string]any) *Capture {
	if !c.enabled {
		return &Capture{noop: true}
	}
	return &Capture{
		store: c.store,
		event: &GenerationEvent{
			ID:         generateID(),
			Timestamp:  time.Now(),
			Command:    command,
			Inputs:     inputs,
			LvtVersion: getLvtVersion(),
		},
	}
}

// RunRetention deletes events older than retentionDays.
func (c *Collector) RunRetention(ctx context.Context, retentionDays int) error {
	if !c.enabled {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	_, err := c.store.DeleteBefore(ctx, cutoff)
	return err
}

// Capture accumulates data for a single generation event.
type Capture struct {
	store Store
	event *GenerationEvent
	noop  bool
}

// NoopCapture returns a capture that silently discards all data.
func NoopCapture() *Capture {
	return &Capture{noop: true}
}

// SetKit records the kit used for generation.
func (cap *Capture) SetKit(kit string) {
	if cap.noop {
		return
	}
	cap.event.Kit = kit
}

// RecordError adds an error to the capture.
func (cap *Capture) RecordError(err GenerationError) {
	if cap.noop {
		return
	}
	cap.event.Errors = append(cap.event.Errors, err)
}

// RecordFileGenerated adds a generated file path.
func (cap *Capture) RecordFileGenerated(path string) {
	if cap.noop {
		return
	}
	cap.event.FilesGenerated = append(cap.event.FilesGenerated, path)
}

// Complete finalises the capture, computes duration, and stores the event.
// It is always safe to call (noop if disabled).
func (cap *Capture) Complete(success bool, validationJSON string) {
	if cap.noop {
		return
	}
	cap.event.Success = success
	cap.event.ValidationJSON = validationJSON
	cap.event.DurationMs = time.Since(cap.event.Timestamp).Milliseconds()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Best-effort: telemetry errors are silently ignored.
	_ = cap.store.Save(ctx, cap.event)
}

// Event returns the underlying event (for testing). Returns nil on noop captures.
func (cap *Capture) Event() *GenerationEvent {
	if cap.noop {
		return nil
	}
	return cap.event
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// getLvtVersion returns the module version from build info, or "dev" if unavailable.
func getLvtVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
