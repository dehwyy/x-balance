package freezescheduler

import (
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var _ gateway.FreezeScheduler = &Implementation{}

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

func freezeKey(txID string) string {
	return "freeze:" + txID
}
