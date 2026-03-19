package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type UnfreezeRequest struct {
	UserID        user.ID
	TransactionID event.TransactionID
}

type UnfreezeResponse struct {
	UnfrozenAmount decimal.Decimal
	TransactionID  event.TransactionID
}

func (s *Service) Unfreeze(
	ctx context.Context,
	req *UnfreezeRequest,
) (*UnfreezeResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceservice.Service.Unfreeze", dspan.Attr("req", req))
	defer span.End()

	releaseKey := event.TransactionID{Value: req.TransactionID.Value + ":release"}

	existingResp, err := s.eventRepo.GetByTransactionID(ctx, dto.EventGetByTxIDRequest{TransactionID: releaseKey})
	if err == nil {
		response := &UnfreezeResponse{UnfrozenAmount: existingResp.Event.Amount.Value.Abs(), TransactionID: req.TransactionID}
		span.WithAttribute("response", response)
		return response, nil
	}
	if !isNotFound(err) {
		return nil, span.Err(err)
	}

	freezeResp, err := s.eventRepo.GetByTransactionID(ctx, dto.EventGetByTxIDRequest{TransactionID: req.TransactionID})
	if err != nil {
		return nil, span.Err(ErrFreezeNotFound)
	}

	frozenAmount := freezeResp.Event.Amount.Value

	err = s.tx.Do(ctx, "balanceservice.Unfreeze", func(ctx context.Context) error {
		if _, err := s.eventRepo.Create(ctx, dto.EventCreateRequest{
			UserID:        req.UserID,
			Type:          event.TypeFreezeRelease,
			Amount:        event.Amount{Value: frozenAmount.Neg()},
			TransactionID: releaseKey,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, span.Err(err)
	}

	_ = s.freezeScheduler.Cancel(ctx, dto.FreezeCancelRequest{TransactionID: req.TransactionID})
	_ = s.balanceCache.Invalidate(ctx, dto.BalanceCacheInvalidateRequest{UserID: req.UserID})

	response := &UnfreezeResponse{UnfrozenAmount: frozenAmount, TransactionID: req.TransactionID}
	span.WithAttribute("response", response)
	return response, nil
}
