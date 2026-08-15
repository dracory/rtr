package middlewares

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// rateLimiter implements a sliding-window rate limiter using two fixed windows
// (current and previous) with a weighted rate calculation. It is a stdlib-only
// reimplementation of go-chi/httprate's RateLimiter.
type rateLimiter struct {
	requestLimit int
	windowLength time.Duration
	// windowOffset aligns windows to the limiter's start instant rather than
	// the wall clock, so resets spread out instead of all snapping to the same
	// instant.
	windowOffset time.Duration
	keyFn        func(*http.Request) (string, error)
	counter      *limitCounter
	mu           sync.Mutex
}

// limitCounter is an in-memory counter with two buckets (current and previous
// window). It mirrors go-chi/httprate's localCounter.
type limitCounter struct {
	windowLength     time.Duration
	latestWindow     time.Time
	latestCounters   map[string]int
	previousCounters map[string]int
	mu               sync.RWMutex
}

// newLimitCounter creates a new in-memory limit counter.
func newLimitCounter(windowLength time.Duration) *limitCounter {
	return &limitCounter{
		windowLength:     windowLength,
		latestWindow:     time.Now().UTC(),
		latestCounters:   make(map[string]int),
		previousCounters: make(map[string]int),
	}
}

// incrementBy adds amount to the counter for the given key in the current
// window.
func (c *limitCounter) incrementBy(key string, currentWindow time.Time, amount int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evict(currentWindow)

	count := c.latestCounters[key]
	c.latestCounters[key] = count + amount
}

// get returns the counts for the current and previous windows for the given
// key.
func (c *limitCounter) get(key string, currentWindow, previousWindow time.Time) (int, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.latestWindow.Equal(currentWindow) {
		curr := c.latestCounters[key]
		prev := c.previousCounters[key]
		return curr, prev
	}

	if c.latestWindow.Equal(previousWindow) {
		prev := c.latestCounters[key]
		return 0, prev
	}

	return 0, 0
}

// evict shifts or clears the windows when the current window has advanced.
func (c *limitCounter) evict(currentWindow time.Time) {
	if c.latestWindow.Equal(currentWindow) {
		return
	}

	previousWindow := currentWindow.Add(-c.windowLength)
	if c.latestWindow.Equal(previousWindow) {
		c.latestWindow = currentWindow
		// Shift the windows without map re-allocation.
		clear(c.previousCounters)
		c.latestCounters, c.previousCounters = c.previousCounters, c.latestCounters
		return
	}

	c.latestWindow = currentWindow
	clear(c.previousCounters)
	clear(c.latestCounters)
}

// newRateLimiter creates a new rate limiter with the given configuration.
func newRateLimiter(requestLimit int, windowLength time.Duration, keyFn func(*http.Request) (string, error)) *rateLimiter {
	start := time.Now().UTC()
	offset := start.Sub(start.Truncate(windowLength))

	return &rateLimiter{
		requestLimit: requestLimit,
		windowLength: windowLength,
		windowOffset: offset,
		keyFn:        keyFn,
		counter:      newLimitCounter(windowLength),
	}
}

// currentWindow returns the start of the rate-limit window containing t,
// aligned to windowOffset rather than the wall clock.
func (l *rateLimiter) currentWindow(t time.Time) time.Time {
	return t.Add(-l.windowOffset).Truncate(l.windowLength).Add(l.windowOffset)
}

// calculateRate computes the weighted rate across the current and previous
// windows. The previous window's count is weighted by the fraction of time
// remaining in it, giving a smooth sliding-window approximation.
func (l *rateLimiter) calculateRate(key string) float64 {
	now := time.Now().UTC()
	currentWindow := l.currentWindow(now)
	previousWindow := currentWindow.Add(-l.windowLength)

	currCount, prevCount := l.counter.get(key, currentWindow, previousWindow)

	diff := now.Sub(currentWindow)
	rate := float64(prevCount)*(float64(l.windowLength)-float64(diff))/float64(l.windowLength) + float64(currCount)
	return rate
}

// handler returns the middleware handler that enforces the rate limit.
func (l *rateLimiter) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := l.keyFn(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionRequired)
			return
		}

		l.mu.Lock()
		rate := l.calculateRate(key)
		rateInt := int(math.Round(rate))

		if rateInt+1 > l.requestLimit {
			l.mu.Unlock()
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.requestLimit))
			// Clamp remaining to 0; the sliding-window weighted rate can exceed
			// the per-window limit after a burst in the previous window, which
			// would otherwise produce a negative value.
			remaining := l.requestLimit - rateInt
			if remaining < 0 {
				remaining = 0
			}
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("Retry-After", strconv.Itoa(int(l.windowLength.Seconds())))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		currentWindow := l.currentWindow(time.Now().UTC())
		l.counter.incrementBy(key, currentWindow, 1)
		l.mu.Unlock()

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.requestLimit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(l.requestLimit-rateInt-1))

		next.ServeHTTP(w, r)
	})
}

// keyByRemoteAddr uses r.RemoteAddr (including port) as the rate-limit key.
func keyByRemoteAddr(r *http.Request) (string, error) {
	return r.RemoteAddr, nil
}

// keyByIP extracts the client IP from r.RemoteAddr (stripping the port) and
// canonicalizes IPv6 addresses to their /64 prefix. It mirrors
// go-chi/httprate's KeyByIP.
func keyByIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return canonicalizeIP(ip), nil
}

// canonicalizeIP normalizes a client IP string for use as a rate-limit key:
//   - IPv4 addresses are returned unchanged.
//   - IPv6 addresses are reduced to their /64 prefix.
//   - Any other string, including "", is returned unchanged.
func canonicalizeIP(ip string) string {
	isIPv6 := false
	for i := 0; !isIPv6 && i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return ip
		case ':':
			isIPv6 = true
		}
	}
	if !isIPv6 {
		return ip
	}

	ipv6 := net.ParseIP(ip)
	if ipv6 == nil {
		return ip
	}

	return ipv6.Mask(net.CIDRMask(64, 128)).String()
}
