package snapshotcron

func (w *Worker) Stop() { w.cron.Stop() }
