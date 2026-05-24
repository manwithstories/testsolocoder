package middleware

import (
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var (
	rateLimiters = make(map[string]*RateLimiter)
	limiterMu    sync.Mutex
)

func GetRateLimiter(name string) *RateLimiter {
	limiterMu.Lock()
	defer limiterMu.Unlock()

	if limiter, exists := rateLimiters[name]; exists {
		return limiter
	}

	limiter := &RateLimiter{
		visitors: make(map[string]*visitor),
	}
	rateLimiters[name] = limiter
	return limiter
}

func (rl *RateLimiter) Allow(key string, maxRequests int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists || time.Since(v.lastSeen) > window {
		rl.visitors[key] = &visitor{
			count:    1,
			lastSeen: time.Now(),
		}
		return true
	}

	v.lastSeen = time.Now()
	v.count++

	return v.count <= maxRequests
}

func (rl *RateLimiter) Cleanup(staleTime time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, v := range rl.visitors {
		if time.Since(v.lastSeen) > staleTime {
			delete(rl.visitors, key)
		}
	}
}
