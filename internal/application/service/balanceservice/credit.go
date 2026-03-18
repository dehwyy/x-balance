package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CreditRequest struct {
	UserID        string
	Amount        decimal.Decimal
	TransactionID string
}

type CreditResponse struct {
	NewBalance    decimal.Decimal
	TransactionID string
}

func (s *Service) Credit(
	ctx context.Context,
	req CreditRequest,
) (*CreditResponse, error) {
	existing, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: req.TransactionID})
	if err == nil {
		bal, _, err := s.computeBalance(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		_ = existing
		return &CreditResponse{NewBalance: bal, TransactionID: req.TransactionID}, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	var newBalance decimal.Decimal

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(ctx, "balanceservice.Credit", func(ctx context.Context) error {
			snap, err := s.snapshotRepo.GetLatestByUserID(ctx, req.UserID)
			if err != nil {
				return err
			}

			if err := s.snapshotRepo.UpdateVersion(ctx, snap); err != nil {
				return err
			}

			newEvent := &event.Event{
				UserID:        req.UserID,
				Type:          event.TypeCredit,
				Amount:        event.Amount{Value: req.Amount},
				TransactionID: event.TransactionID{Value: req.TransactionID},
				SnapshotID:    &snap.ID.Value,
			}

			if _, err := s.eventRepo.Create(ctx, newEvent); err != nil {
				return err
			}

			deltaBalance, frozen, err := s.eventRepo.SumSinceSnapshot(ctx, req.UserID, snap.ID)
			if err != nil {
				return err
			}
			newBalance = snap.Balance.Value.Add(deltaBalance).Sub(frozen)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	_ = s.balanceCache.Invalidate(ctx, req.UserID)
	_ = s.maybeCreateSnapshot(ctx, req.UserID)

	return &CreditResponse{NewBalance: newBalance, TransactionID: req.TransactionID}, nil
}

func isNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
