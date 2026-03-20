package main

import (
	"github.com/dehwyy/tracerfx/pkg/tracer"
	"go.uber.org/fx"

	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/runners"
	"github.com/dehwyy/x-balance/internal/runners/modules"
)

func main() {
	fx.New(
		tracer.FxModule(
			tracer.WithHost("92.39.53.47:4317"),
			tracer.WithServiceName("x-balance"),
		),
		fx.Provide(
			config.Load,
			runners.NewGORM,
			runners.NewTxManager,
			runners.NewRedis,
		),

		modules.InfrastructureRepositoryModule,
		modules.InfrastructureGatewayModule,
		modules.ApplicationModule,
		modules.WorkersModule,
		modules.DeliveryModule,

		fx.Invoke(
			runners.RunAPI,
			runners.RunWorkers,
		),
	).Run()
}
