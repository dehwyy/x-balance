package balancecache

import (
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var _ gateway.BalanceCache = &Implementation{}

type cacheEntry struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

type Implementation struct {
	client *redis.Client
}

type Opts struct {
	fx.In
	Client *redis.Client
}

func New(opts Opts) *Implementation {
	return &Implementation{client: opts.Client}
}

func balanceKey(userID string) string {
	return "balance:" + userID
}
