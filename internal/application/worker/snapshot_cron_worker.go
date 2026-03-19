package worker

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
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
	uid := user.ID{Value: userID}

	snapResp, err := w.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: uid})
	var snap snapshot.Snapshot
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			createResp, err := w.snapshotRepo.Create(ctx, dto.SnapshotCreateRequest{
				UserID:  uid,
				Balance: snapshot.Balance{Value: decimal.Zero},
				Version: snapshot.Version{Value: 0},
			})
			if err != nil {
				return err
			}
			snap = createResp.Snapshot
		} else {
			return err
		}
	} else {
		snap = snapResp.Snapshot
	}

	countResp, err := w.eventRepo.CountSinceSnapshot(ctx, dto.EventCountSinceSnapshotRequest{
		UserID:     uid,
		SnapshotID: snap.ID,
	})
	if err != nil {
		return err
	}

	if countResp.Count == 0 {
		return nil
	}

	sumResp, err := w.eventRepo.SumSinceSnapshot(ctx, dto.EventSumSinceSnapshotRequest{
		UserID:     uid,
		SnapshotID: snap.ID,
	})
	if err != nil {
		return err
	}

	newBalance := snap.Balance.Value.Add(sumResp.Available).Sub(sumResp.Frozen)

	_, err = w.snapshotRepo.Create(ctx, dto.SnapshotCreateRequest{
		UserID:  uid,
		Balance: snapshot.Balance{Value: newBalance},
		Version: snapshot.Version{Value: 0},
	})
	return err
}

func (w *SnapshotCronWorker) EnsureInitialSnapshot(ctx context.Context, userID string) error {
	uid := user.ID{Value: userID}

	_, err := w.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: uid})
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	_, err = w.snapshotRepo.Create(ctx, dto.SnapshotCreateRequest{
		UserID:  uid,
		Balance: snapshot.Balance{Value: decimal.Zero},
		Version: snapshot.Version{Value: 0},
	})
	return err
}
