package database

import (
	"context"
	"fmt"
	"time"

	"devix-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	if cfg.URL == "" {
		return nil, nil
	}

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if cfg.Password != "" {
		opts.Password = cfg.Password
	}
	opts.DB = cfg.DB
	opts.PoolSize = 20
	opts.MinIdleConns = 5
	opts.MaxRetries = 3
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
