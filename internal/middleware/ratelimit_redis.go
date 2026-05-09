package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, int, time.Duration) {
	if rl.client == nil {
		return true, maxRequests, 0
	}

	now := time.Now()
	windowStart := now.Add(-window)
	member := fmt.Sprintf("%d:%s", now.UnixNano(), uuid.New().String()[:8])

	pipe := rl.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: member})
	countCmd := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, window)
	_, _ = pipe.Exec(ctx)

	count := countCmd.Val()
	remaining := maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	if int(count) > maxRequests {
		ttl := rl.client.TTL(ctx, key).Val()
		return false, remaining, ttl
	}

	return true, remaining, 0
}

func RedisRateLimit(limiter *RedisRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rl:ip:%s", c.ClientIP())

		allowed, remaining, retryAfter := limiter.Allow(c.Request.Context(), key, maxRequests, window)

		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()
	}
}

func UserRateLimit(limiter *RedisRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if userID, ok := GetUserID(c); ok {
			key = fmt.Sprintf("rl:user:%s:%s", userID.String(), c.Request.URL.Path)
		} else {
			key = fmt.Sprintf("rl:ip:%s:%s", c.ClientIP(), c.Request.URL.Path)
		}

		allowed, remaining, retryAfter := limiter.Allow(c.Request.Context(), key, maxRequests, window)

		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()
	}
}

func AuthRateLimitRedis(limiter *RedisRateLimiter, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rl:auth:%s:%s", c.ClientIP(), c.Request.URL.Path)

		allowed, _, retryAfter := limiter.Allow(c.Request.Context(), key, maxRequests, window)

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()
	}
}
