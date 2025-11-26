package e2e

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"testing"

	"github.com/chromedp/chromedp"
	e2etest "github.com/livetemplate/lvt/testing"
)

// Chrome pool singleton variables for lazy initialization
var (
	chromePool     *ChromePool
	chromePoolOnce sync.Once
	chromePoolMu   sync.Mutex
)

// ChromePool manages a pool of reusable Chrome containers
type ChromePool struct {
	containers []*ChromeContainer
	available  chan *ChromeContainer
	mu         sync.Mutex
	t          *testing.T
}

// ChromeContainer wraps a Chrome container instance
type ChromeContainer struct {
	port int
	url  string
}

// NewChromePool creates a pool with n Chrome containers
func NewChromePool(t *testing.T, size int) *ChromePool {
	if t != nil {
		t.Helper()
	}

	pool := &ChromePool{
		containers: make([]*ChromeContainer, 0, size),
		available:  make(chan *ChromeContainer, size),
		t:          t,
	}

	// Start Chrome containers
	log.Printf("Starting Chrome pool with %d containers...", size)
	for i := 0; i < size; i++ {
		container := pool.startChrome(i)
		pool.containers = append(pool.containers, container)
		pool.available <- container
	}

	log.Printf("✅ Chrome pool ready with %d containers", size)
	return pool
}

func (p *ChromePool) startChrome(index int) *ChromeContainer {
	if p.t != nil {
		p.t.Helper()
	}

	// Allocate unique port for this container
	port, err := e2etest.GetFreePort()
	if err != nil {
		if p.t != nil {
			p.t.Fatalf("Failed to allocate port for Chrome pool container %d: %v", index, err)
		}
		panic(fmt.Sprintf("Failed to allocate port for Chrome pool container %d: %v", index, err))
	}

	// Start Chrome container
	if err := e2etest.StartDockerChrome(p.t, port); err != nil {
		if p.t != nil {
			p.t.Fatalf("Failed to start Chrome pool container %d: %v", index, err)
		}
		panic(fmt.Sprintf("Failed to start Chrome pool container %d: %v", index, err))
	}

	// Get WebSocket URL
	url := fmt.Sprintf("ws://localhost:%d", port)

	log.Printf("Chrome pool container %d started on port %d", index, port)

	return &ChromeContainer{
		port: port,
		url:  url,
	}
}

// Get retrieves an available Chrome container from the pool
// Blocks if all containers are in use
func (p *ChromePool) Get() *ChromeContainer {
	log.Printf("Waiting for available Chrome container (pool size: %d)...", len(p.containers))
	container := <-p.available
	log.Printf("Got Chrome container on port %d", container.port)
	return container
}

// Release returns a Chrome container to the pool after cleanup
func (p *ChromePool) Release(container *ChromeContainer) {
	// Reset Chrome state
	p.resetChrome(container)

	// Return to pool
	p.available <- container
}

func (p *ChromePool) resetChrome(container *ChromeContainer) {
	// Create allocator context from WS URL
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), container.url)
	defer allocCancel()

	// Create fresh context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Clear cookies, storage, navigate to blank
	chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Clear storage
			chromedp.Evaluate(`
				localStorage.clear();
				sessionStorage.clear();
			`, nil).Do(ctx)
			return nil
		}),
		chromedp.Navigate("about:blank"),
	)
}

// Cleanup stops all Chrome containers
// Called from TestMain after all tests complete, so we pass nil to avoid logging to completed test
func (p *ChromePool) Cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, container := range p.containers {
		e2etest.StopDockerChrome(nil, container.port)
	}

	log.Println("✅ Chrome pool cleaned up")
}

// GetPooledChrome returns a Chrome context from the pool
// Lazily initializes the Chrome pool on first call
func GetPooledChrome(t *testing.T) (context.Context, context.CancelFunc, func()) {
	t.Helper()

	// Check if Docker is available
	if _, err := exec.Command("docker", "version").CombinedOutput(); err != nil {
		t.Skip("Docker not available, skipping E2E test")
	}

	// Lazily initialize the Chrome pool
	chromePoolOnce.Do(func() {
		chromePoolMu.Lock()
		defer chromePoolMu.Unlock()

		t.Log("Initializing Chrome pool (first use)...")
		// Pool size of 8 to handle up to 8 concurrent tests (we have 7 tests using Chrome)
		chromePool = NewChromePool(t, 8)
	})

	container := chromePool.Get()

	// Create allocator from WS URL
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), container.url)

	// Create context for this test
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Return context and cleanup function
	cleanup := func() {
		cancel()
		allocCancel()
		chromePool.Release(container)
	}

	return ctx, cancel, cleanup
}
