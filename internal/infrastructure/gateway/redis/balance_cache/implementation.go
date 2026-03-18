package balancecache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

var _ gateway.BalanceCache = &Implementation{}

type cacheEntry struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

type Implementation struct {
	client *redis.Client
}

func New(client *redis.Client) *Implementation {
	return &Implementation{client: client}
}

func balanceKey(userID string) string {
	return fmt.Sprintf("balance:%s", userID)
}

func (impl *Implementation) Get(
	ctx context.Context,
	userID string,
) (decimal.Decimal, decimal.Decimal, bool, error) {
	val, err := impl.client.Get(ctx, balanceKey(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return decimal.Zero, decimal.Zero, false, nil
		}
		return decimal.Zero, decimal.Zero, false, err
	}

	var entry cacheEntry
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return decimal.Zero, decimal.Zero, false, err
	}

	available, err := decimal.NewFromString(entry.Available)
	if err != nil {
		return decimal.Zero, decimal.Zero, false, err
	}

	frozen, err := decimal.NewFromString(entry.Frozen)
	if err != nil {
		return decimal.Zero, decimal.Zero, false, err
	}

	return available, frozen, true, nil
}

func (impl *Implementation) Set(
	ctx context.Context,
	userID string,
	available decimal.Decimal,
	frozen decimal.Decimal,
) error {
	entry := cacheEntry{
		Available: available.String(),
		Frozen:    frozen.String(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return impl.client.Set(ctx, balanceKey(userID), data, 0).Err()
}

func (impl *Implementation) Invalidate(
	ctx context.Context,
	userID string,
) error {
	return impl.client.Del(ctx, balanceKey(userID)).Err()
}
