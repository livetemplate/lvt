// Package ratelimit provides HTTP rate limiting middleware using token buckets.
//
// It tracks per-IP request rates with LRU eviction when the IP map is full.
// For high-concurrency workloads, a sharded implementation reduces lock
// contention by distributing IPs across multiple mutex-guarded maps.
//
// Usage:
//
//	rl := ratelimit.New(ctx,
//	    ratelimit.WithRate(10),
//	    ratelimit.WithBurst(20),
//	)
//	defer rl.Close()
//	handler := rl.Middleware()(mux)
package ratelimit

import (
	"container/list"
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/time/rate"
)

// defaultNumShards is the number of shards used by the sharded rate limiter.
const defaultNumShards = 16

// Config holds rate limiter configuration.
type Config struct {
	RPS              float64
	Burst            int
	MaxIPs           int
	DenyHandler      http.HandlerFunc
	SweepInterval    time.Duration
	StaleThreshold   time.Duration
	EvictLogInterval time.Duration
}

// Option configures a Config.
type Option func(*Config)

// WithRate sets the steady-state request rate per IP (requests per second).
func WithRate(rps float64) Option {
	return func(c *Config) { c.RPS = rps }
}

// WithBurst sets the maximum burst size per IP.
func WithBurst(burst int) Option {
	return func(c *Config) { c.Burst = burst }
}

// WithMaxIPs sets the maximum number of unique IPs to track.
func WithMaxIPs(n int) Option {
	return func(c *Config) { c.MaxIPs = n }
}

// WithDenyHandler sets a custom handler for rate-limited requests.
func WithDenyHandler(h http.HandlerFunc) Option {
	return func(c *Config) { c.DenyHandler = h }
}

// WithSweepInterval sets how often stale entries are cleaned up.
func WithSweepInterval(d time.Duration) Option {
	return func(c *Config) { c.SweepInterval = d }
}

// WithStaleThreshold sets how long an IP must be idle before cleanup removes it.
func WithStaleThreshold(d time.Duration) Option {
	return func(c *Config) { c.StaleThreshold = d }
}

// WithEvictLogInterval sets the minimum time between eviction log messages.
func WithEvictLogInterval(d time.Duration) Option {
	return func(c *Config) { c.EvictLogInterval = d }
}

func defaultConfig() Config {
	return Config{
		RPS:              100,
		Burst:            200,
		MaxIPs:           10000,
		SweepInterval:    5 * time.Minute,
		StaleThreshold:   10 * time.Minute,
		EvictLogInterval: 30 * time.Second,
	}
}

func defaultDenyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Retry-After", "1")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("rate limit exceeded"))
}

// ipLimiter tracks a per-IP token bucket and its position in the LRU list.
type ipLimiter struct {
	ip       string
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter manages per-IP rate limiting with LRU eviction.
type RateLimiter struct {
	mw     func(http.Handler) http.Handler
	cancel context.CancelFunc
	done   chan struct{}
}

// New creates a RateLimiter that runs a background cleanup goroutine.
// The goroutine is stopped when Close is called or ctx is cancelled.
//
// For MaxIPs >= 16, a sharded implementation is used to reduce lock
// contention. For smaller values, a single-mutex implementation is used.
func New(ctx context.Context, opts ...Option) *RateLimiter {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.DenyHandler == nil {
		cfg.DenyHandler = defaultDenyHandler
	}
	if cfg.MaxIPs <= 0 {
		cfg.MaxIPs = 10000
	}

	ctx, cancel := context.WithCancel(ctx)

	rl := &RateLimiter{cancel: cancel}

	if cfg.MaxIPs < defaultNumShards {
		rl.mw, rl.done = newSingleMutexLimiter(ctx, &cfg)
	} else {
		rl.mw, rl.done = newShardedLimiter(ctx, &cfg, defaultNumShards)
	}

	return rl
}

// Middleware returns the HTTP middleware function.
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return rl.mw
}

// Close stops the cleanup goroutine and waits for it to exit.
func (rl *RateLimiter) Close() {
	rl.cancel()
	<-rl.done
}

// newSingleMutexLimiter creates a rate limiter with a single mutex.
func newSingleMutexLimiter(ctx context.Context, cfg *Config) (func(http.Handler) http.Handler, chan struct{}) {
	var (
		items = make(map[string]*list.Element)
		order = list.New()
		mu    sync.Mutex

		lastEvictLog time.Time
		evictCount   int
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(cfg.SweepInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				now := time.Now()
				for e := order.Back(); e != nil; {
					lim := e.Value.(*ipLimiter)
					prev := e.Prev()
					if now.Sub(lim.lastSeen) > cfg.StaleThreshold {
						order.Remove(e)
						delete(items, lim.ip)
					}
					e = prev
				}
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	deny := cfg.DenyHandler
	rps := cfg.RPS
	burst := cfg.Burst
	maxIPs := cfg.MaxIPs
	evictLogInterval := cfg.EvictLogInterval

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)
			now := time.Now()

			mu.Lock()
			elem, exists := items[ip]
			if exists {
				order.MoveToFront(elem)
				elem.Value.(*ipLimiter).lastSeen = now
			} else {
				if order.Len() >= maxIPs {
					back := order.Back()
					if back != nil {
						evicted := back.Value.(*ipLimiter)
						order.Remove(back)
						delete(items, evicted.ip)
						evictCount++
						if now.Sub(lastEvictLog) >= evictLogInterval {
							slog.Warn("rate limiter evicted least-recent IPs",
								"count", evictCount, "capacity", maxIPs)
							lastEvictLog = now
							evictCount = 0
						}
					}
				}
				lim := &ipLimiter{
					ip:       ip,
					limiter:  rate.NewLimiter(rate.Limit(rps), burst),
					lastSeen: now,
				}
				elem = order.PushFront(lim)
				items[ip] = elem
			}
			lim := elem.Value.(*ipLimiter).limiter
			mu.Unlock()

			if !lim.Allow() {
				deny(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	return middleware, done
}

// shard holds a subset of IP limiters behind its own mutex.
// Padding ensures each shard occupies its own CPU cache line,
// preventing false sharing when adjacent shards are accessed by different cores.
type shard struct {
	mu           sync.Mutex
	items        map[string]*list.Element
	order        *list.List
	maxIPs       int
	lastEvictLog time.Time
	evictCount   int
	_            [64]byte // pad to 128 bytes (2 cache lines)
}

// Compile-time assertion: shard struct must be exactly 128 bytes.
const _ = uint(128 - unsafe.Sizeof(shard{}))

type shardedLimiter struct {
	shards           []shard
	numShards        uint32
	rps              float64
	burst            int
	totalMaxIPs      int
	evictLogInterval time.Duration
}

func (sl *shardedLimiter) shardFor(ip string) *shard {
	h := uint32(2166136261)
	for i := 0; i < len(ip); i++ {
		h ^= uint32(ip[i])
		h *= 16777619
	}
	return &sl.shards[h%sl.numShards]
}

func newShardedLimiter(ctx context.Context, cfg *Config, numShards int) (func(http.Handler) http.Handler, chan struct{}) {
	sl := &shardedLimiter{
		shards:           make([]shard, numShards),
		numShards:        uint32(numShards),
		rps:              cfg.RPS,
		burst:            cfg.Burst,
		totalMaxIPs:      cfg.MaxIPs,
		evictLogInterval: cfg.EvictLogInterval,
	}

	base := cfg.MaxIPs / numShards
	remainder := cfg.MaxIPs % numShards
	for i := range sl.shards {
		shardCap := base
		if i < remainder {
			shardCap++
		}
		sl.shards[i] = shard{
			items:  make(map[string]*list.Element),
			order:  list.New(),
			maxIPs: shardCap,
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(cfg.SweepInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				for i := range sl.shards {
					s := &sl.shards[i]
					s.mu.Lock()
					for e := s.order.Back(); e != nil; {
						lim := e.Value.(*ipLimiter)
						prev := e.Prev()
						if now.Sub(lim.lastSeen) > cfg.StaleThreshold {
							s.order.Remove(e)
							delete(s.items, lim.ip)
						}
						e = prev
					}
					s.mu.Unlock()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	deny := cfg.DenyHandler

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)
			now := time.Now()
			s := sl.shardFor(ip)

			s.mu.Lock()
			elem, exists := s.items[ip]
			if exists {
				s.order.MoveToFront(elem)
				elem.Value.(*ipLimiter).lastSeen = now
			} else {
				if s.order.Len() >= s.maxIPs {
					back := s.order.Back()
					if back != nil {
						evicted := back.Value.(*ipLimiter)
						s.order.Remove(back)
						delete(s.items, evicted.ip)
						s.evictCount++
						if now.Sub(s.lastEvictLog) >= sl.evictLogInterval {
							slog.Warn("rate limiter evicted least-recent IPs",
								"count", s.evictCount,
								"shard_capacity", s.maxIPs,
								"total_capacity", sl.totalMaxIPs)
							s.lastEvictLog = now
							s.evictCount = 0
						}
					}
				}
				lim := &ipLimiter{
					ip:       ip,
					limiter:  rate.NewLimiter(rate.Limit(sl.rps), sl.burst),
					lastSeen: now,
				}
				elem = s.order.PushFront(lim)
				s.items[ip] = elem
			}
			lim := elem.Value.(*ipLimiter).limiter
			s.mu.Unlock()

			if !lim.Allow() {
				deny(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	return middleware, done
}

// GetClientIP extracts the client IP from the request.
// It only trusts X-Forwarded-For / X-Real-IP when the immediate peer is a
// loopback or private address (i.e., behind a reverse proxy).
func GetClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	peerIP := net.ParseIP(host)
	trustedProxy := peerIP != nil && (peerIP.IsLoopback() || peerIP.IsPrivate())

	if trustedProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			clientIP := xff
			if i := strings.IndexByte(xff, ','); i > 0 {
				clientIP = xff[:i]
			}
			if ip := net.ParseIP(strings.TrimSpace(clientIP)); ip != nil {
				return ip.String()
			}
			return strings.TrimSpace(clientIP)
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			if ip := net.ParseIP(strings.TrimSpace(xri)); ip != nil {
				return ip.String()
			}
			return strings.TrimSpace(xri)
		}
	}

	if peerIP != nil {
		return peerIP.String()
	}
	return host
}
