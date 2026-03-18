package freezescheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
)

var _ gateway.FreezeScheduler = &Implementation{}

type Implementation struct {
	client *redis.Client
}

func New(client *redis.Client) *Implementation {
	return &Implementation{client: client}
}

func freezeKey(txID string) string {
	return fmt.Sprintf("freeze:%s", txID)
}

func (impl *Implementation) Schedule(
	ctx context.Context,
	txID string,
	ttlSeconds int64,
) error {
	if ttlSeconds <= 0 {
		return nil
	}

	return impl.client.Set(
		ctx,
		freezeKey(txID),
		"",
		time.Duration(ttlSeconds)*time.Second,
	).Err()
}

func (impl *Implementation) Cancel(
	ctx context.Context,
	txID string,
) error {
	return impl.client.Del(ctx, freezeKey(txID)).Err()
}
