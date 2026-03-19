package runners

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/dehwyy/x-balance/internal/application/worker"
	"github.com/dehwyy/x-balance/internal/config"
)

type RunWorkersOpts struct {
	fx.In
	LC             fx.Lifecycle
	Config         *config.Config
	Log            zerolog.Logger
	FreezeWorker   *worker.FreezeExpiryWorker
	SnapshotWorker *worker.SnapshotCronWorker
}

func RunWorkers(opts RunWorkersOpts) {
	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go opts.FreezeWorker.Start(ctx)

			if opts.Config.SnapshotCron != "" {
				if err := opts.SnapshotWorker.Start(ctx, opts.Config.SnapshotCron); err != nil {
					opts.Log.Error().Err(err).Msg("failed to start snapshot cron worker")
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			opts.SnapshotWorker.Stop()
			return nil
		},
	})
}
