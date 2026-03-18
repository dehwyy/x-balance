package runners

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	"github.com/dehwyy/x-balance/internal/config"
	"go.uber.org/fx"
)

var ServicesModule = fx.Module("services",
	fx.Provide(
		userservice.New,
		transactionservice.New,
		func(cfg *config.Config) balanceservice.BalanceConfig {
			return balanceservice.BalanceConfig{
				SnapshotEveryN: cfg.SnapshotEveryN,
			}
		},
		balanceservice.New,
	),
)
