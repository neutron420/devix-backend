package middleware

import (
	"context"
	"fmt"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AbuseDetector struct {
	client    *redis.Client
	threshold int
	window    time.Duration
}

func NewAbuseDetector(client *redis.Client, threshold int, window time.Duration) *AbuseDetector {
	return &AbuseDetector{
		client:    client,
		threshold: threshold,
		window:    window,
	}
}

func (ad *AbuseDetector) RecordSuspicious(ctx context.Context, ip string) {
	if ad.client == nil {
		return
	}
	key := fmt.Sprintf("abuse:%s", ip)
	ad.client.Incr(ctx, key)
	ad.client.Expire(ctx, key, ad.window)
}

func (ad *AbuseDetector) IsBlocked(ctx context.Context, ip string) bool {
	if ad.client == nil {
		return false
	}
	key := fmt.Sprintf("abuse:%s", ip)
	score, err := ad.client.Get(ctx, key).Int()
	if err != nil {
		return false
	}
	return score >= ad.threshold
}

func AbuseProtection(detector *AbuseDetector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if detector == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()

		if detector.IsBlocked(c.Request.Context(), ip) {
			response.Abort(c, apperrors.Forbidden("Your IP has been temporarily blocked due to suspicious activity"))
			return
		}

		ua := c.GetHeader("User-Agent")
		if ua == "" {
			detector.RecordSuspicious(c.Request.Context(), ip)
		}

		c.Next()

		if c.Writer.Status() == 429 || c.Writer.Status() == 401 {
			detector.RecordSuspicious(c.Request.Context(), ip)
		}
	}
}

func LoginProtection(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("login_fail:%s", ip)

		failCount, _ := client.Get(c.Request.Context(), key).Int()
		if failCount >= 10 {
			ttl := client.TTL(c.Request.Context(), key).Val()
			c.Header("Retry-After", fmt.Sprintf("%.0f", ttl.Seconds()))
			response.Abort(c, apperrors.TooManyRequests())
			return
		}

		c.Next()

		if c.Writer.Status() == 401 {
			client.Incr(c.Request.Context(), key)
			client.Expire(c.Request.Context(), key, 15*time.Minute)
		}
	}
}
