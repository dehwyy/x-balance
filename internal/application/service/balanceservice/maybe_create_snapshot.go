package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

func (s *Service) maybeCreateSnapshot(ctx context.Context, userID user.ID) error {
	if s.config.SnapshotEveryN <= 0 {
		return nil
	}

	snapResp, err := s.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: userID})
	if err != nil {
		return err
	}
	snap := snapResp.Snapshot

	countResp, err := s.eventRepo.CountSinceSnapshot(ctx, dto.EventCountSinceSnapshotRequest{UserID: userID, SnapshotID: snap.ID})
	if err != nil {
		return err
	}

	if int(countResp.Count) < s.config.SnapshotEveryN {
		return nil
	}

	sumResp, err := s.eventRepo.SumSinceSnapshot(ctx, dto.EventSumSinceSnapshotRequest{UserID: userID, SnapshotID: snap.ID})
	if err != nil {
		return err
	}

	newBalance := snap.Balance.Value.Add(sumResp.Available).Sub(sumResp.Frozen)

	_, err = s.snapshotRepo.Create(ctx, dto.SnapshotCreateRequest{
		UserID:  userID,
		Balance: snapshot.Balance{Value: newBalance},
		Version: snapshot.Version{Value: 0},
	})
	if err != nil {
		return err
	}

	_ = s.balanceCache.Set(ctx, dto.BalanceCacheSetRequest{
		UserID:    userID,
		Available: newBalance.Sub(sumResp.Frozen),
		Frozen:    decimal.Zero,
	})
	return nil
}
