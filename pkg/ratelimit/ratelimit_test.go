package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func reqFromIP(ip string) *http.Request {
	r := httptest.NewRequest("GET", "/test", nil)
	r.RemoteAddr = ip + ":12345"
	return r
}

func newTestRL(t *testing.T, opts ...Option) *RateLimiter {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	rl := New(ctx, opts...)
	t.Cleanup(func() {
		cancel()
		rl.Close()
	})
	return rl
}

func wrapHandler(t *testing.T, opts ...Option) http.Handler {
	t.Helper()
	rl := newTestRL(t, opts...)
	return rl.Middleware()(okHandler())
}

func TestDefaults(t *testing.T) {
	cfg := defaultConfig()
	if cfg.RPS != 100 {
		t.Errorf("default RPS = %f, want 100", cfg.RPS)
	}
	if cfg.Burst != 200 {
		t.Errorf("default Burst = %d, want 200", cfg.Burst)
	}
	if cfg.MaxIPs != 10000 {
		t.Errorf("default MaxIPs = %d, want 10000", cfg.MaxIPs)
	}
}

func TestOptionsOverrideDefaults(t *testing.T) {
	rl := newTestRL(t,
		WithRate(50),
		WithBurst(100),
		WithMaxIPs(500),
	)
	_ = rl // just ensure it doesn't panic
}

func TestBasicRateLimiting(t *testing.T) {
	// burst=2, very low rps so tokens don't refill
	wrapped := wrapHandler(t, WithRate(0.001), WithBurst(2), WithMaxIPs(100))

	// First 2 requests succeed (burst)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, w.Code)
		}
	}

	// 3rd request should be rate limited
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request 3: expected 429, got %d", w.Code)
	}
}

func TestRetryAfterHeader(t *testing.T) {
	wrapped := wrapHandler(t, WithRate(0.001), WithBurst(1), WithMaxIPs(100))

	// Use burst
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Should get 429 with Retry-After
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
	if ra := w.Header().Get("Retry-After"); ra != "1" {
		t.Errorf("Retry-After = %q, want %q", ra, "1")
	}
}

func TestCustomDenyHandler(t *testing.T) {
	customCalled := false
	deny := func(w http.ResponseWriter, r *http.Request) {
		customCalled = true
		w.Header().Set("Retry-After", "60")
		http.Redirect(w, r, "/auth?error=rate_limited", http.StatusSeeOther)
	}

	wrapped := wrapHandler(t,
		WithRate(0.001), WithBurst(1), WithMaxIPs(100),
		WithDenyHandler(deny),
	)

	// Use burst
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))

	// Should trigger custom deny handler
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))

	if !customCalled {
		t.Error("custom deny handler was not called")
	}
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
}

func TestLRUEviction(t *testing.T) {
	// maxIPs=3 to force single-mutex path (< defaultNumShards)
	wrapped := wrapHandler(t, WithRate(100), WithBurst(100), WithMaxIPs(3))

	// Fill to capacity with 3 IPs
	for _, ip := range []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"} {
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP(ip))
		if w.Code != http.StatusOK {
			t.Fatalf("IP %s: expected 200, got %d", ip, w.Code)
		}
	}

	// 4th IP should succeed (LRU eviction), not error
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("4.4.4.4"))
	if w.Code != http.StatusOK {
		t.Errorf("4th IP at capacity: expected 200, got %d", w.Code)
	}
}

func TestEvictedIPGetsFreshLimiter(t *testing.T) {
	// burst=1, maxIPs=2 to force single-mutex path
	wrapped := wrapHandler(t, WithRate(100), WithBurst(1), WithMaxIPs(2))

	// IP uses its burst token
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Second request from same IP should be 429 (burst exhausted)
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}

	// Push 1.1.1.1 out by filling with 2 other IPs
	for _, ip := range []string{"2.2.2.2", "3.3.3.3"} {
		w = httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP(ip))
		if w.Code != http.StatusOK {
			t.Fatalf("IP %s: expected 200, got %d", ip, w.Code)
		}
	}

	// 1.1.1.1 returns — should get a fresh limiter
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Errorf("evicted IP returning: expected 200, got %d", w.Code)
	}
}

func TestMRUNotEvicted(t *testing.T) {
	// maxIPs=3 to force single-mutex path
	wrapped := wrapHandler(t, WithRate(100), WithBurst(100), WithMaxIPs(3))

	// Fill: A, B, C (order: C=front, B, A=back)
	for _, ip := range []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"} {
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP(ip))
		if w.Code != http.StatusOK {
			t.Fatalf("IP %s: expected 200, got %d", ip, w.Code)
		}
	}

	// Touch A → moves to front
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("10.0.0.1"))
	if w.Code != http.StatusOK {
		t.Fatalf("touch A: expected 200, got %d", w.Code)
	}

	// New IP D → evicts B (back), not A
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("10.0.0.4"))
	if w.Code != http.StatusOK {
		t.Fatalf("new IP D: expected 200, got %d", w.Code)
	}

	// A should still be present
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("10.0.0.1"))
	if w.Code != http.StatusOK {
		t.Errorf("A after eviction: expected 200, got %d", w.Code)
	}
}

func TestConcurrentAccess(t *testing.T) {
	wrapped := wrapHandler(t, WithRate(1000), WithBurst(1000), WithMaxIPs(100))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ip := fmt.Sprintf("10.0.%d.%d", id/256, id%256)
			for j := 0; j < 10; j++ {
				w := httptest.NewRecorder()
				wrapped.ServeHTTP(w, reqFromIP(ip))
			}
		}(i)
	}
	wg.Wait()
}

func TestCloseStopsGoroutine(t *testing.T) {
	ctx := context.Background()
	rl := New(ctx, WithRate(100), WithBurst(100))

	done := make(chan struct{})
	go func() {
		rl.Close()
		close(done)
	}()

	select {
	case <-done:
		// OK
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return within 2s")
	}
}

func TestStaleCleanup(t *testing.T) {
	// rps=0.001 so tokens don't refill naturally
	rl := newTestRL(t,
		WithRate(0.001), WithBurst(1), WithMaxIPs(100),
		WithSweepInterval(50*time.Millisecond),
		WithStaleThreshold(100*time.Millisecond),
	)
	wrapped := rl.Middleware()(okHandler())

	// First request uses burst → 200
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("5.5.5.5"))
	if w.Code != http.StatusOK {
		t.Fatalf("first: expected 200, got %d", w.Code)
	}

	// Second request → 429
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("5.5.5.5"))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("second: expected 429, got %d", w.Code)
	}

	// Wait for cleanup
	time.Sleep(400 * time.Millisecond)

	// Entry should be gone — fresh limiter → 200
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("5.5.5.5"))
	if w.Code != http.StatusOK {
		t.Errorf("after cleanup: expected 200, got %d", w.Code)
	}
}

func TestActiveEntrySurvivesCleanup(t *testing.T) {
	rl := newTestRL(t,
		WithRate(0.001), WithBurst(1), WithMaxIPs(100),
		WithSweepInterval(50*time.Millisecond),
		WithStaleThreshold(300*time.Millisecond),
	)
	wrapped := rl.Middleware()(okHandler())

	// Use burst token → 200
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("6.6.6.6"))
	if w.Code != http.StatusOK {
		t.Fatalf("initial: expected 200, got %d", w.Code)
	}

	// Keep alive by touching every 100ms (< 300ms threshold)
	for i := 0; i < 4; i++ {
		time.Sleep(100 * time.Millisecond)
		w = httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP("6.6.6.6"))
		if w.Code != http.StatusTooManyRequests {
			t.Errorf("tick %d: expected 429 (entry alive), got %d", i, w.Code)
		}
	}
}

func TestShardedBoundary(t *testing.T) {
	// MaxIPs=64 hits the sharded path (>= defaultNumShards * 4)
	wrapped := wrapHandler(t, WithRate(100), WithBurst(100), WithMaxIPs(64))

	for i := 0; i < 20; i++ {
		ip := fmt.Sprintf("192.168.%d.%d", i/256, i%256)
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, reqFromIP(ip))
		if w.Code != http.StatusOK {
			t.Errorf("IP %s: expected 200, got %d", ip, w.Code)
		}
	}
}

func TestShardedRateLimiting(t *testing.T) {
	// Enough IPs to use sharded path
	wrapped := wrapHandler(t, WithRate(0.001), WithBurst(1), WithMaxIPs(100))

	// First request succeeds
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("8.8.8.8"))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Second from same IP → 429
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("8.8.8.8"))
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestShardedConcurrentAccess(t *testing.T) {
	wrapped := wrapHandler(t, WithRate(1000), WithBurst(1000), WithMaxIPs(1000))

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ip := fmt.Sprintf("10.%d.%d.%d", id/65536, (id/256)%256, id%256)
			for j := 0; j < 10; j++ {
				w := httptest.NewRecorder()
				wrapped.ServeHTTP(w, reqFromIP(ip))
			}
		}(i)
	}
	wg.Wait()
}

func TestGetClientIP_Direct(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.1:12345"
	if ip := GetClientIP(r); ip != "203.0.113.1" {
		t.Errorf("GetClientIP = %q, want %q", ip, "203.0.113.1")
	}
}

func TestGetClientIP_IgnoresXFFFromPublicIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.1:12345"
	r.Header.Set("X-Forwarded-For", "10.0.0.1")
	if ip := GetClientIP(r); ip != "203.0.113.1" {
		t.Errorf("GetClientIP = %q, want %q (should ignore XFF from public IP)", ip, "203.0.113.1")
	}
}

func TestGetClientIP_TrustsXFFFromLoopback(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("X-Forwarded-For", "203.0.113.50, 10.0.0.1")
	if ip := GetClientIP(r); ip != "203.0.113.50" {
		t.Errorf("GetClientIP = %q, want %q", ip, "203.0.113.50")
	}
}

func TestGetClientIP_TrustsXFFFromPrivate(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.0.0.1:12345"
	r.Header.Set("X-Forwarded-For", "198.51.100.5")
	if ip := GetClientIP(r); ip != "198.51.100.5" {
		t.Errorf("GetClientIP = %q, want %q", ip, "198.51.100.5")
	}
}

func TestGetClientIP_TrustsXRealIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("X-Real-IP", "198.51.100.10")
	if ip := GetClientIP(r); ip != "198.51.100.10" {
		t.Errorf("GetClientIP = %q, want %q", ip, "198.51.100.10")
	}
}

func TestGetClientIP_XFFBeforeXRealIP(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("X-Forwarded-For", "203.0.113.1")
	r.Header.Set("X-Real-IP", "198.51.100.10")
	if ip := GetClientIP(r); ip != "203.0.113.1" {
		t.Errorf("GetClientIP = %q, want %q (XFF should take priority)", ip, "203.0.113.1")
	}
}

func TestGetClientIP_NoPort(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.1"
	if ip := GetClientIP(r); ip != "203.0.113.1" {
		t.Errorf("GetClientIP = %q, want %q", ip, "203.0.113.1")
	}
}

func TestMaxIPsDefaultsToTenThousand(t *testing.T) {
	// MaxIPs=0 should default to 10000 (sharded path)
	rl := newTestRL(t, WithRate(100), WithBurst(100), WithMaxIPs(0))
	wrapped := rl.Middleware()(okHandler())

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestNegativeMaxIPsDefaultsToTenThousand(t *testing.T) {
	rl := newTestRL(t, WithRate(100), WithBurst(100), WithMaxIPs(-1))
	wrapped := rl.Middleware()(okHandler())

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDifferentIPsIndependent(t *testing.T) {
	wrapped := wrapHandler(t, WithRate(0.001), WithBurst(1), WithMaxIPs(100))

	// IP A uses burst
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusOK {
		t.Fatalf("IP A: expected 200, got %d", w.Code)
	}

	// IP A is now rate limited
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("1.1.1.1"))
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("IP A again: expected 429, got %d", w.Code)
	}

	// IP B should still succeed
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, reqFromIP("2.2.2.2"))
	if w.Code != http.StatusOK {
		t.Errorf("IP B: expected 200, got %d", w.Code)
	}
}
