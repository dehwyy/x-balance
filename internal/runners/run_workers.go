package runners

import (
	"context"

	"go.uber.org/fx"

	freezeexpiry "github.com/dehwyy/x-balance/internal/application/worker/freeze_expiry"
	snapshotcron "github.com/dehwyy/x-balance/internal/application/worker/snapshot_cron"
	"github.com/dehwyy/x-balance/internal/config"
)

type RunWorkersOpts struct {
	fx.In
	LC             fx.Lifecycle
	Config         *config.Config
	FreezeWorker   *freezeexpiry.Worker
	SnapshotWorker *snapshotcron.Worker
}

func RunWorkers(opts RunWorkersOpts) {
	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go opts.FreezeWorker.Start(ctx)

			if opts.Config.SnapshotCron != "" {
				if err := opts.SnapshotWorker.Start(ctx, opts.Config.SnapshotCron); err != nil {
					return err
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
