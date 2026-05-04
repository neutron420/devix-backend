package middleware

import (
	"net/http"
	"sync"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// rateLimitEntry tracks request counts per key.
type rateLimitEntry struct {
	count     int
	expiresAt time.Time
}

// InMemoryRateLimiter provides a simple in-memory rate limiter.
// Use Redis-based rate limiting in production with multiple instances.
type InMemoryRateLimiter struct {
	mu      sync.RWMutex
	entries map[string]*rateLimitEntry
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter.
func NewInMemoryRateLimiter() *InMemoryRateLimiter {
	rl := &InMemoryRateLimiter{
		entries: make(map[string]*rateLimitEntry),
	}

	// Cleanup expired entries every minute
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

// Allow checks if a request is allowed under the rate limit.
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

// RateLimit returns a rate-limiting middleware.
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

// AuthRateLimit returns a stricter rate limiter for auth endpoints.
func AuthRateLimit(limiter *InMemoryRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP + path to rate limit specific auth endpoints
		key := "auth:" + c.ClientIP() + ":" + c.Request.URL.Path

		if !limiter.Allow(key, maxRequests, window) {
			c.Header("Retry-After", window.String())
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()
	}
}
