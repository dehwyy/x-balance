package balanceservice

import (
	"context"
	"time"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type FreezeRequest struct {
	UserID               string
	Amount               decimal.Decimal
	TransactionID        string
	FreezeTimeoutSeconds int64
}

type FreezeResponse struct {
	FrozenAmount  decimal.Decimal
	TransactionID string
}

func (s *Service) Freeze(
	ctx context.Context,
	req FreezeRequest,
) (*FreezeResponse, error) {
	existing, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: req.TransactionID})
	if err == nil {
		return &FreezeResponse{FrozenAmount: existing.Amount.Value, TransactionID: req.TransactionID}, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(ctx, "balanceservice.Freeze", func(ctx context.Context) error {
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

			var expiresAt *time.Time
			if req.FreezeTimeoutSeconds > 0 {
				t := time.Now().Add(time.Duration(req.FreezeTimeoutSeconds) * time.Second)
				expiresAt = &t
			}

			newEvent := &event.Event{
				UserID:          req.UserID,
				Type:            event.TypeFreezeHold,
				Amount:          event.Amount{Value: req.Amount},
				TransactionID:   event.TransactionID{Value: req.TransactionID},
				SnapshotID:      &snap.ID.Value,
				FreezeExpiresAt: expiresAt,
			}

			if _, err := s.eventRepo.Create(ctx, newEvent); err != nil {
				return err
			}

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	if req.FreezeTimeoutSeconds > 0 {
		_ = s.freezeScheduler.Schedule(ctx, req.TransactionID, req.FreezeTimeoutSeconds)
	}

	_ = s.balanceCache.Invalidate(ctx, req.UserID)

	return &FreezeResponse{FrozenAmount: req.Amount, TransactionID: req.TransactionID}, nil
}
