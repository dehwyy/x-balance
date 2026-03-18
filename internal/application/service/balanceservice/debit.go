package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type DebitRequest struct {
	UserID        string
	Amount        decimal.Decimal
	TransactionID string
}

type DebitResponse struct {
	NewBalance    decimal.Decimal
	TransactionID string
}

func (s *Service) Debit(
	ctx context.Context,
	req DebitRequest,
) (*DebitResponse, error) {
	existing, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: req.TransactionID})
	if err == nil {
		bal, _, err := s.computeBalance(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		_ = existing
		return &DebitResponse{NewBalance: bal, TransactionID: req.TransactionID}, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	var newBalance decimal.Decimal

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(ctx, "balanceservice.Debit", func(ctx context.Context) error {
			snap, err := s.snapshotRepo.GetLatestByUserID(ctx, req.UserID)
			if err != nil {
				return err
			}

			u, err := s.userRepo.GetByID(ctx, user.ID{Value: req.UserID})
			if err != nil {
				return err
			}

			deltaBalance, frozen, err := s.eventRepo.SumSinceSnapshot(ctx, req.UserID, snap.ID)
			if err != nil {
				return err
			}

			available := snap.Balance.Value.Add(deltaBalance).Sub(frozen)
			minAllowed := u.OverdraftLimit.Value.Neg()
			if available.Sub(req.Amount).LessThan(minAllowed) {
				return ErrInsufficientFunds
			}

			if err := s.snapshotRepo.UpdateVersion(ctx, snap); err != nil {
				return err
			}

			newEvent := &event.Event{
				UserID:        req.UserID,
				Type:          event.TypeDebit,
				Amount:        event.Amount{Value: req.Amount.Neg()},
				TransactionID: event.TransactionID{Value: req.TransactionID},
				SnapshotID:    &snap.ID.Value,
			}

			if _, err := s.eventRepo.Create(ctx, newEvent); err != nil {
				return err
			}

			newBalance = available.Sub(req.Amount)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	_ = s.balanceCache.Invalidate(ctx, req.UserID)
	_ = s.maybeCreateSnapshot(ctx, req.UserID)

	return &DebitResponse{NewBalance: newBalance, TransactionID: req.TransactionID}, nil
}
