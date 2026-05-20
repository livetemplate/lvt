//go:build browser

package e2e

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// TestMain sets up shared resources before running tests and cleans up after
func TestMain(m *testing.M) {
	// Setup signal handling for cleanup on interrupt (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("🛑 Interrupted - cleaning up test containers...")
		chromePoolMu.Lock()
		if chromePool != nil {
			chromePool.Cleanup()
		}
		chromePoolMu.Unlock()
		cleanupAllTestContainers()
		log.Println("✅ Cleanup complete")
		os.Exit(1)
	}()

	// Cleanup any leftover containers from previous runs (Chrome + app containers)
	// This is safe to run even if Docker is not available - it will just log a warning
	cleanupAllTestContainers()

	// Run tests
	code := m.Run()

	// Cleanup Chrome pool if it was initialized
	chromePoolMu.Lock()
	if chromePool != nil {
		chromePool.Cleanup()
	}
	chromePoolMu.Unlock()

	// Final cleanup of any remaining containers
	cleanupAllTestContainers()

	os.Exit(code)
}
