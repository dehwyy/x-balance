package snapshotcron

import (
	"context"

	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
)

func (w *Worker) Start(ctx context.Context, cronExpr string) error {
	_, err := w.cron.AddFunc(cronExpr, func() {
		if err := w.createSnapshots(ctx); err != nil {
			tlog.FromContext(ctx).Error("snapshot cron job failed", "err", err)
		}
	})
	if err != nil {
		return err
	}

	w.cron.Start()
	return nil
}
