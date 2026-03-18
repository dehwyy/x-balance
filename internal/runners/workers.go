package runners

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/dehwyy/x-balance/internal/application/worker"
	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var WorkersModule = fx.Module("workers",
	fx.Provide(
		worker.NewFreezeExpiryWorker,
		worker.NewSnapshotCronWorker,
	),
	fx.Invoke(func(
		lc fx.Lifecycle,
		freezeWorker *worker.FreezeExpiryWorker,
		snapshotWorker *worker.SnapshotCronWorker,
		cfg *config.Config,
		log zerolog.Logger,
		// These are needed to satisfy FreezeExpiryWorker dependencies
		_ *redis.Client,
		_ *balanceservice.Service,
		// These are needed to satisfy SnapshotCronWorker dependencies
		_ repository.EventRepository,
		_ repository.SnapshotRepository,
		_ repository.UserRepository,
	) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go freezeWorker.Start(ctx)

				if cfg.SnapshotCron != "" {
					if err := snapshotWorker.Start(ctx, cfg.SnapshotCron); err != nil {
						log.Error().Err(err).Msg("failed to start snapshot cron worker")
					}
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				snapshotWorker.Stop()
				return nil
			},
		})
	}),
)
