package modules

import (
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	balancecache "github.com/dehwyy/x-balance/internal/infrastructure/gateway/redis/balance_cache"
	freezescheduler "github.com/dehwyy/x-balance/internal/infrastructure/gateway/redis/freeze_scheduler"
	"go.uber.org/fx"
)

var InfrastructureGatewayModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			balancecache.New,
			fx.As(new(gateway.BalanceCache)),
		),
		fx.Annotate(
			freezescheduler.New,
			fx.As(new(gateway.FreezeScheduler)),
		),
	),
)
