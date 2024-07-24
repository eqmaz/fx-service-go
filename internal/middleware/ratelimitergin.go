package middleware

import (
	"sync"
	"time"
)

type RateLimiterConfig struct {
	MaxRequests int
	Timeframe   int // in seconds
}

// ginRateLimiter is a simple rate limiter for gin
// Probably not the best implementation, or place to keep this
type ginRateLimiter struct {
	visitors map[string]*ginVisitor
	mu       sync.Mutex
}

func newGinRateLimiter() *ginRateLimiter {
	return &ginRateLimiter{
		visitors: make(map[string]*ginVisitor),
	}
}

type ginVisitor struct {
	lastSeen time.Time
	limiter  *time.Ticker
	count    int
}

func (rl *ginRateLimiter) getVisitor(ip string, timeframe int) *ginVisitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := time.NewTicker(time.Duration(timeframe) * time.Second)
		v = &ginVisitor{
			limiter: limiter,
			count:   0,
		}
		rl.visitors[ip] = v
		go func() {
			<-limiter.C
			rl.mu.Lock()
			delete(rl.visitors, ip)
			rl.mu.Unlock()
		}()
	}

	v.lastSeen = time.Now()
	return v
}
