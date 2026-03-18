package runners

import (
	"context"
	"fmt"

	"github.com/dehwyy/x-balance/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	// Enable keyspace notifications for freeze expiry
	client.ConfigSet(ctx, "notify-keyspace-events", "KEA")

	return client, nil
}
