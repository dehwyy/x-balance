package runners

import (
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	balancecache "github.com/dehwyy/x-balance/internal/infrastructure/gateway/redis/balance_cache"
	freezescheduler "github.com/dehwyy/x-balance/internal/infrastructure/gateway/redis/freeze_scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var GatewayModule = fx.Module("gateway",
	fx.Provide(
		fx.Annotate(
			func(client *redis.Client) gateway.BalanceCache {
				return balancecache.New(client)
			},
		),
		fx.Annotate(
			func(client *redis.Client) gateway.FreezeScheduler {
				return freezescheduler.New(client)
			},
		),
	),
)
