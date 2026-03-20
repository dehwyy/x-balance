package balanceservice

import (
	"context"

	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

func (s *Service) maybeCreateSnapshot(ctx context.Context, userID user.ID) error {
	if s.config.SnapshotEveryN <= 0 {
		return nil
	}

	snapshotResult, err := s.snapshotRepo.GetLatestByUserID(
		ctx,
		dto.SnapshotGetLatestByUserIDRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return err
	}
	snap := snapshotResult.Snapshot

	countResult, err := s.eventRepo.CountSinceSnapshot(
		ctx,
		dto.EventCountSinceSnapshotRequest{
			UserID:     userID,
			SnapshotID: snap.ID,
		},
	)
	if err != nil {
		return err
	}

	if int(countResult.Count) < s.config.SnapshotEveryN {
		return nil
	}

	sumSinceSnapshot, err := s.eventRepo.SumSinceSnapshot(
		ctx,
		dto.EventSumSinceSnapshotRequest{
			UserID:     userID,
			SnapshotID: snap.ID,
		},
	)
	if err != nil {
		return err
	}

	available, frozen := snap.ComputeBalance(
		sumSinceSnapshot.Available,
		sumSinceSnapshot.Frozen,
	)

	_, err = s.snapshotRepo.Create(
		ctx,
		dto.SnapshotCreateRequest{
			UserID:  userID,
			Balance: snapshot.Balance(available),
			Version: snapshot.Version(0),
		},
	)
	if err != nil {
		return err
	}

	if err := s.balanceCache.Set(
		ctx,
		dto.BalanceCacheSetRequest{
			UserID:    userID,
			Available: available.Sub(frozen),
			Frozen:    decimal.Zero,
		},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to set balance cache after snapshot creation", "err", err)
	}
	return nil
}
