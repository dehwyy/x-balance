package modules

import (
	freezeexpiry "github.com/dehwyy/x-balance/internal/application/worker/freeze_expiry"
	snapshotcron "github.com/dehwyy/x-balance/internal/application/worker/snapshot_cron"
	"go.uber.org/fx"
)

var WorkersModule = fx.Options(
	fx.Provide(
		freezeexpiry.New,
		snapshotcron.New,
	),
)
