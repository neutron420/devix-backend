package middleware

import (
	"net/http"
	"sync"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	expiresAt time.Time
}

type InMemoryRateLimiter struct {
	mu      sync.RWMutex
	entries map[string]*rateLimitEntry
}

func NewInMemoryRateLimiter() *InMemoryRateLimiter {
	rl := &InMemoryRateLimiter{
		entries: make(map[string]*rateLimitEntry),
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *InMemoryRateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, entry := range rl.entries {
		if now.After(entry.expiresAt) {
			delete(rl.entries, key)
		}
	}
}

func (rl *InMemoryRateLimiter) Allow(key string, maxRequests int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.entries[key]

	if !exists || now.After(entry.expiresAt) {
		rl.entries[key] = &rateLimitEntry{
			count:     1,
			expiresAt: now.Add(window),
		}
		return true
	}

	if entry.count >= maxRequests {
		return false
	}

	entry.count++
	return true
}

func RateLimit(limiter *InMemoryRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if !limiter.Allow(key, maxRequests, window) {
			c.Header("Retry-After", window.String())
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Header("X-RateLimit-Limit", http.StatusText(maxRequests))
		c.Next()
	}
}

func AuthRateLimit(limiter *InMemoryRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		key := "auth:" + c.ClientIP() + ":" + c.Request.URL.Path

		if !limiter.Allow(key, maxRequests, window) {
			c.Header("Retry-After", window.String())
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()
	}
}
