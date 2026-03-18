package balanceservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/shopspring/decimal"
)

type UnfreezeRequest struct {
	UserID        string
	TransactionID string
}

type UnfreezeResponse struct {
	UnfrozenAmount decimal.Decimal
	TransactionID  string
}

func (s *Service) Unfreeze(
	ctx context.Context,
	req UnfreezeRequest,
) (*UnfreezeResponse, error) {
	releaseKey := req.TransactionID + ":release"

	existing, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: releaseKey})
	if err == nil {
		return &UnfreezeResponse{UnfrozenAmount: existing.Amount.Value.Abs(), TransactionID: req.TransactionID}, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	freezeEvent, err := s.eventRepo.GetByTransactionID(ctx, event.TransactionID{Value: req.TransactionID})
	if err != nil {
		return nil, ErrFreezeNotFound
	}

	frozenAmount := freezeEvent.Amount.Value

	err = s.tx.Do(ctx, "balanceservice.Unfreeze", func(ctx context.Context) error {
		releaseEvent := &event.Event{
			UserID:        req.UserID,
			Type:          event.TypeFreezeRelease,
			Amount:        event.Amount{Value: frozenAmount.Neg()},
			TransactionID: event.TransactionID{Value: releaseKey},
		}

		if _, err := s.eventRepo.Create(ctx, releaseEvent); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.freezeScheduler.Cancel(ctx, req.TransactionID)
	_ = s.balanceCache.Invalidate(ctx, req.UserID)

	return &UnfreezeResponse{UnfrozenAmount: frozenAmount, TransactionID: req.TransactionID}, nil
}
