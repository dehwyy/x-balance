package balancecache

import (
	"fmt"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
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
