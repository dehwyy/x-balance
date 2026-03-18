package worker

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type SnapshotCronWorker struct {
	eventRepo    repository.EventRepository
	snapshotRepo repository.SnapshotRepository
	userRepo     repository.UserRepository
	log          zerolog.Logger
	cron         *cron.Cron
}

func NewSnapshotCronWorker(
	eventRepo repository.EventRepository,
	snapshotRepo repository.SnapshotRepository,
	userRepo repository.UserRepository,
	log zerolog.Logger,
) *SnapshotCronWorker {
	return &SnapshotCronWorker{
		eventRepo:    eventRepo,
		snapshotRepo: snapshotRepo,
		userRepo:     userRepo,
		log:          log,
		cron:         cron.New(),
	}
}

func (w *SnapshotCronWorker) Start(ctx context.Context, cronExpr string) error {
	_, err := w.cron.AddFunc(cronExpr, func() {
		if err := w.createSnapshots(ctx); err != nil {
			w.log.Error().Err(err).Msg("snapshot cron job failed")
		}
	})
	if err != nil {
		return err
	}

	w.cron.Start()
	return nil
}

func (w *SnapshotCronWorker) Stop() {
	w.cron.Stop()
}

func (w *SnapshotCronWorker) createSnapshots(ctx context.Context) error {
	w.log.Info().Msg("running snapshot cron job")
	return nil
}

func (w *SnapshotCronWorker) CreateSnapshotForUser(ctx context.Context, userID string) error {
	snap, err := w.snapshotRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			snap = &snapshot.Snapshot{
				UserID:  userID,
				Balance: snapshot.Balance{Value: decimal.Zero},
				Version: snapshot.Version{Value: 0},
			}
			snap, err = w.snapshotRepo.Create(ctx, snap)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	count, err := w.eventRepo.CountSinceSnapshot(ctx, userID, snap.ID)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	deltaBalance, frozen, err := w.eventRepo.SumSinceSnapshot(ctx, userID, snap.ID)
	if err != nil {
		return err
	}

	newBalance := snap.Balance.Value.Add(deltaBalance).Sub(frozen)

	_, err = w.snapshotRepo.Create(ctx, &snapshot.Snapshot{
		UserID:  userID,
		Balance: snapshot.Balance{Value: newBalance},
		Version: snapshot.Version{Value: 0},
	})
	return err
}

func (w *SnapshotCronWorker) EnsureInitialSnapshot(ctx context.Context, userID string) error {
	_, err := w.snapshotRepo.GetLatestByUserID(ctx, userID)
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	_, err = w.snapshotRepo.Create(ctx, &snapshot.Snapshot{
		UserID:  userID,
		Balance: snapshot.Balance{Value: decimal.Zero},
		Version: snapshot.Version{Value: 0},
	})
	return err
}
