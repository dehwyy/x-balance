package main

import (
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"github.com/dehwyy/x-balance/internal/application/worker"
	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/runners"
	"github.com/dehwyy/x-balance/internal/runners/modules"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			runners.NewDB,
			runners.NewRedis,
			func() zerolog.Logger { return log.Logger },
			worker.NewFreezeExpiryWorker,
			worker.NewSnapshotCronWorker,
		),
		txmanager.NewGorm(),

		modules.InfrastructureRepositoryModule,
		modules.InfrastructureGatewayModule,
		modules.ApplicationModule,
		modules.DeliveryModule,

		fx.Invoke(
			runners.RunAPI,
			runners.RunWorkers,
		),
	).Run()
}
