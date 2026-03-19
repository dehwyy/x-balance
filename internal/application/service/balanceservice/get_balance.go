package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type GetBalanceRequest struct {
	UserID user.ID
}

type GetBalanceResponse struct {
	Available decimal.Decimal
	Frozen    decimal.Decimal
	Total     decimal.Decimal
}

func (s *Service) GetBalance(
	ctx context.Context,
	req *GetBalanceRequest,
) (*GetBalanceResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceservice.Service.GetBalance", dspan.Attr("req", req))
	defer span.End()

	cacheResp, err := s.balanceCache.Get(ctx, dto.BalanceCacheGetRequest{UserID: req.UserID})
	if err == nil && cacheResp.Found {
		response := &GetBalanceResponse{
			Available: cacheResp.Available,
			Frozen:    cacheResp.Frozen,
			Total:     cacheResp.Available.Add(cacheResp.Frozen),
		}
		span.WithAttribute("response", response)
		return response, nil
	}

	available, frozen, err := s.computeBalance(ctx, req.UserID)
	if err != nil {
		return nil, span.Err(err)
	}

	_ = s.balanceCache.Set(ctx, dto.BalanceCacheSetRequest{
		UserID:    req.UserID,
		Available: available,
		Frozen:    frozen,
	})

	response := &GetBalanceResponse{
		Available: available,
		Frozen:    frozen,
		Total:     available.Add(frozen),
	}
	span.WithAttribute("response", response)
	return response, nil
}

func (s *Service) computeBalance(ctx context.Context, userID user.ID) (decimal.Decimal, decimal.Decimal, error) {
	snapResp, err := s.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: userID})
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	snap := snapResp.Snapshot

	sumResp, err := s.eventRepo.SumSinceSnapshot(ctx, dto.EventSumSinceSnapshotRequest{UserID: userID, SnapshotID: snap.ID})
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	available := snap.Balance.Value.Add(sumResp.Available).Sub(sumResp.Frozen)
	return available, sumResp.Frozen, nil
}
