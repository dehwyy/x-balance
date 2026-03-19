package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CreditRequest struct {
	UserID        user.ID
	Amount        decimal.Decimal
	TransactionID event.TransactionID
}

type CreditResponse struct {
	NewBalance    decimal.Decimal
	TransactionID event.TransactionID
}

func (s *Service) Credit(
	ctx context.Context,
	req *CreditRequest,
) (*CreditResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceservice.Service.Credit", dspan.Attr("req", req))
	defer span.End()

	_, err := s.eventRepo.GetByTransactionID(ctx, dto.EventGetByTxIDRequest{TransactionID: req.TransactionID})
	if err == nil {
		bal, _, err := s.computeBalance(ctx, req.UserID)
		if err != nil {
			return nil, span.Err(err)
		}
		response := &CreditResponse{NewBalance: bal, TransactionID: req.TransactionID}
		span.WithAttribute("response", response)
		return response, nil
	}
	if !isNotFound(err) {
		return nil, span.Err(err)
	}

	var newBalance decimal.Decimal

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(ctx, "balanceservice.Credit", func(ctx context.Context) error {
			snapResp, err := s.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: req.UserID})
			if err != nil {
				return err
			}
			snap := snapResp.Snapshot

			if err := s.snapshotRepo.UpdateVersion(ctx, dto.SnapshotUpdateVersionRequest{Snapshot: snap}); err != nil {
				return err
			}

			snapID := event.SnapshotID{Value: snap.ID.Value}
			if _, err := s.eventRepo.Create(ctx, dto.EventCreateRequest{
				UserID:        req.UserID,
				Type:          event.TypeCredit,
				Amount:        event.Amount{Value: req.Amount},
				TransactionID: req.TransactionID,
				SnapshotID:    &snapID,
			}); err != nil {
				return err
			}

			sumResp, err := s.eventRepo.SumSinceSnapshot(ctx, dto.EventSumSinceSnapshotRequest{UserID: req.UserID, SnapshotID: snap.ID})
			if err != nil {
				return err
			}
			newBalance = snap.Balance.Value.Add(sumResp.Available).Sub(sumResp.Frozen)
			return nil
		})
	})
	if err != nil {
		return nil, span.Err(err)
	}

	_ = s.balanceCache.Invalidate(ctx, dto.BalanceCacheInvalidateRequest{UserID: req.UserID})
	_ = s.maybeCreateSnapshot(ctx, req.UserID)

	response := &CreditResponse{NewBalance: newBalance, TransactionID: req.TransactionID}
	span.WithAttribute("response", response)
	return response, nil
}

func isNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
