package snapshotcron

import (
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type Worker struct {
	eventRepo    repository.EventRepository
	snapshotRepo repository.SnapshotRepository
	userRepo     repository.UserRepository
	cron         *cron.Cron
}

type WorkerOpts struct {
	fx.In
	EventRepo    repository.EventRepository
	SnapshotRepo repository.SnapshotRepository
	UserRepo     repository.UserRepository
}

func New(opts WorkerOpts) *Worker {
	return &Worker{
		eventRepo:    opts.EventRepo,
		snapshotRepo: opts.SnapshotRepo,
		userRepo:     opts.UserRepo,
		cron:         cron.New(),
	}
}
