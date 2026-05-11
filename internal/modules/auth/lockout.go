package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"devix-backend/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type LockoutService struct {
	client *redis.Client
	cfg    config.AuthLockoutConfig
	log    zerolog.Logger
}

type LockoutError struct {
	RetryAfter time.Duration
}

func (e *LockoutError) Error() string {
	return fmt.Sprintf("too many login attempts, retry after %s", e.RetryAfter.Round(time.Second))
}

func NewLockoutService(client *redis.Client, cfg config.AuthLockoutConfig, log zerolog.Logger) *LockoutService {
	return &LockoutService{
		client: client,
		cfg:    cfg,
		log:    log.With().Str("component", "auth_lockout").Logger(),
	}
}

func (s *LockoutService) NormalizeKey(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}

func (s *LockoutService) Check(ctx context.Context, identifier string) (time.Duration, bool) {
	if s == nil || s.client == nil {
		return 0, false
	}
	key := "auth:lock:" + identifier
	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return 0, false
	}
	return ttl, true
}

func (s *LockoutService) RegisterFailure(ctx context.Context, identifier string) {
	if s == nil || s.client == nil {
		return
	}
	failKey := "auth:fail:" + identifier
	count, err := s.client.Incr(ctx, failKey).Result()
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to increment auth failure count")
		return
	}
	_ = s.client.Expire(ctx, failKey, s.cfg.Window).Err()

	if int(count) < s.cfg.MaxAttempts {
		return
	}

	over := int(count) - s.cfg.MaxAttempts
	backoff := s.cfg.BaseLock
	for i := 0; i < over; i++ {
		if backoff >= s.cfg.MaxLock {
			backoff = s.cfg.MaxLock
			break
		}
		backoff *= 2
		if backoff > s.cfg.MaxLock {
			backoff = s.cfg.MaxLock
			break
		}
	}

	lockKey := "auth:lock:" + identifier
	if err := s.client.Set(ctx, lockKey, "1", backoff).Err(); err != nil {
		s.log.Warn().Err(err).Msg("failed to set auth lockout")
	}
}

func (s *LockoutService) Clear(ctx context.Context, identifier string) {
	if s == nil || s.client == nil {
		return
	}
	_, _ = s.client.Del(ctx, "auth:fail:"+identifier, "auth:lock:"+identifier).Result()
}
