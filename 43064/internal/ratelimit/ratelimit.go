package ratelimit

import (
	"sync"
	"time"

	"github.com/notification-center/internal/config"
	"github.com/notification-center/internal/logger"
	"golang.org/x/time/rate"
	"go.uber.org/zap"
)

type RateLimiter struct {
	globalLimiter   *rate.Limiter
	channelLimiters map[uint]*rate.Limiter
	mu              sync.RWMutex
	cfg             *config.RateLimitConfig
}

var instance *RateLimiter
var once sync.Once

func GetInstance(cfg *config.RateLimitConfig) *RateLimiter {
	once.Do(func() {
		instance = &RateLimiter{
			channelLimiters: make(map[uint]*rate.Limiter),
			cfg:             cfg,
		}
		if cfg.Global.Enabled {
			instance.globalLimiter = rate.NewLimiter(
				rate.Limit(cfg.Global.RequestsPerSecond),
				cfg.Global.Burst,
			)
			logger.Info("global rate limiter initialized",
				zap.Float64("rps", cfg.Global.RequestsPerSecond),
				zap.Int("burst", cfg.Global.Burst),
			)
		}
	})
	return instance
}

func (rl *RateLimiter) AllowGlobal() bool {
	if rl.globalLimiter == nil {
		return true
	}
	return rl.globalLimiter.Allow()
}

func (rl *RateLimiter) AllowChannel(channelID uint, channelRPS float64, burst int) bool {
	rl.mu.RLock()
	limiter, exists := rl.channelLimiters[channelID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter, exists = rl.channelLimiters[channelID]
		if !exists {
			rps := channelRPS
			if rps <= 0 {
				rps = rl.cfg.DefaultChannel.RequestsPerSecond
			}
			b := burst
			if b <= 0 {
				b = rl.cfg.DefaultChannel.Burst
			}
			if rl.cfg.DefaultChannel.Enabled {
				limiter = rate.NewLimiter(rate.Limit(rps), b)
			}
			rl.channelLimiters[channelID] = limiter
			logger.Info("channel rate limiter initialized",
				zap.Uint("channel_id", channelID),
				zap.Float64("rps", rps),
				zap.Int("burst", b),
			)
		}
		rl.mu.Unlock()
	}

	if limiter == nil {
		return true
	}
	return limiter.Allow()
}

func (rl *RateLimiter) Allow(channelID uint, channelRPS float64, burst int) bool {
	if !rl.AllowGlobal() {
		logger.Warn("global rate limit exceeded")
		return false
	}
	if !rl.AllowChannel(channelID, channelRPS, burst) {
		logger.Warn("channel rate limit exceeded", zap.Uint("channel_id", channelID))
		return false
	}
	return true
}

func (rl *RateLimiter) WaitGlobal(ctx interface{}) error {
	if rl.globalLimiter == nil {
		return nil
	}
	return rl.globalLimiter.Wait(nil)
}

func (rl *RateLimiter) ResetChannel(channelID uint) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.channelLimiters, channelID)
}

func (rl *RateLimiter) ResetAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.channelLimiters = make(map[uint]*rate.Limiter)
}

func (rl *RateLimiter) GlobalStats() (float64, int) {
	if rl.globalLimiter == nil {
		return 0, 0
	}
	return float64(rl.globalLimiter.Limit()), rl.globalLimiter.Burst()
}

func (rl *RateLimiter) ChannelStats(channelID uint) (float64, int, bool) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	limiter, exists := rl.channelLimiters[channelID]
	if !exists || limiter == nil {
		return 0, 0, false
	}
	return float64(limiter.Limit()), limiter.Burst(), true
}

type TokenBucket struct {
	capacity     int
	tokens       int
	fillRate     float64
	lastRefilled time.Time
	mu           sync.Mutex
}

func NewTokenBucket(capacity int, fillRatePerSecond float64) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		fillRate:     fillRatePerSecond,
		lastRefilled: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefilled).Seconds()
	newTokens := int(elapsed * tb.fillRate)

	if newTokens > 0 {
		tb.tokens = min(tb.tokens+newTokens, tb.capacity)
		tb.lastRefilled = now
	}
}

func (tb *TokenBucket) Tokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
