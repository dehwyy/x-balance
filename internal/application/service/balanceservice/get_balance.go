package balanceservice

import (
	"context"

	"github.com/shopspring/decimal"
)

type GetBalanceRequest struct {
	UserID string
}

type GetBalanceResponse struct {
	Available decimal.Decimal
	Frozen    decimal.Decimal
	Total     decimal.Decimal
}

func (s *Service) GetBalance(
	ctx context.Context,
	req GetBalanceRequest,
) (*GetBalanceResponse, error) {
	available, frozen, found, err := s.balanceCache.Get(ctx, req.UserID)
	if err == nil && found {
		return &GetBalanceResponse{
			Available: available,
			Frozen:    frozen,
			Total:     available.Add(frozen),
		}, nil
	}

	available, frozen, err = s.computeBalance(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	_ = s.balanceCache.Set(ctx, req.UserID, available, frozen)

	return &GetBalanceResponse{
		Available: available,
		Frozen:    frozen,
		Total:     available.Add(frozen),
	}, nil
}

func (s *Service) computeBalance(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, error) {
	snap, err := s.snapshotRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	deltaBalance, frozen, err := s.eventRepo.SumSinceSnapshot(ctx, userID, snap.ID)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	available := snap.Balance.Value.Add(deltaBalance).Sub(frozen)
	return available, frozen, nil
}
