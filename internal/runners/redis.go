package runners

import (
	"context"
	"fmt"

	"github.com/dehwyy/x-balance/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type NewRedisOpts struct {
	fx.In
	Config *config.Config
}

func NewRedis(opts NewRedisOpts) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(opts.Config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(redisOpts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	client.ConfigSet(ctx, "notify-keyspace-events", "KEA")

	return client, nil
}
