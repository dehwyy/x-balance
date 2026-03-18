package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/shopspring/decimal"
)

func (s *Service) maybeCreateSnapshot(ctx context.Context, userID string) error {
	if s.config.SnapshotEveryN <= 0 {
		return nil
	}

	snap, err := s.snapshotRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return err
	}

	count, err := s.eventRepo.CountSinceSnapshot(ctx, userID, snap.ID)
	if err != nil {
		return err
	}

	if int(count) < s.config.SnapshotEveryN {
		return nil
	}

	deltaBalance, frozen, err := s.eventRepo.SumSinceSnapshot(ctx, userID, snap.ID)
	if err != nil {
		return err
	}

	newBalance := snap.Balance.Value.Add(deltaBalance).Sub(frozen)

	_, err = s.snapshotRepo.Create(ctx, &snapshot.Snapshot{
		UserID:  userID,
		Balance: snapshot.Balance{Value: newBalance},
		Version: snapshot.Version{Value: 0},
	})
	if err != nil {
		return err
	}

	_ = s.balanceCache.Set(ctx, userID, newBalance.Sub(frozen), decimal.Zero)
	return nil
}
