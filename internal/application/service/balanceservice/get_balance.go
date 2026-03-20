package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
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
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.GetBalance",
		dspan.Attr("req", req),
	)
	defer span.End()

	cacheResult, err := s.balanceCache.Get(
		ctx,
		dto.BalanceCacheGetRequest{
			UserID: req.UserID,
		},
	)
	if err == nil && cacheResult.Found {
		return dspan.Response(
			span,
			&GetBalanceResponse{
				Available: cacheResult.Available,
				Frozen:    cacheResult.Frozen,
				Total:     cacheResult.Available.Add(cacheResult.Frozen),
			},
		), nil
	}

	available, frozen, err := s.computeBalance(
		ctx,
		req.UserID,
	)
	if err != nil {
		return nil, span.Err(err)
	}

	if err := s.balanceCache.Set(
		ctx,
		dto.BalanceCacheSetRequest{
			UserID:    req.UserID,
			Available: available,
			Frozen:    frozen,
		},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to set balance cache", "err", err)
	}

	return dspan.Response(
		span,
		&GetBalanceResponse{
			Available: available,
			Frozen:    frozen,
			Total:     available.Add(frozen),
		},
	), nil
}

func (s *Service) computeBalance(ctx context.Context, userID user.ID) (decimal.Decimal, decimal.Decimal, error) {
	snapshotResult, err := s.snapshotRepo.GetLatestByUserID(
		ctx,
		dto.SnapshotGetLatestByUserIDRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	snap := snapshotResult.Snapshot

	sumSinceSnapshot, err := s.eventRepo.SumSinceSnapshot(
		ctx,
		dto.EventSumSinceSnapshotRequest{
			UserID:     userID,
			SnapshotID: snap.ID,
		},
	)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	available, frozen := snap.ComputeBalance(
		sumSinceSnapshot.Available,
		sumSinceSnapshot.Frozen,
	)
	return available, frozen, nil
}
