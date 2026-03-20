package balanceservice

import (
	"context"
	"errors"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
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
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.Unfreeze",
		dspan.Attr("req", req),
	)
	defer span.End()

	releaseKey := req.TransactionID.ReleaseKey()

	existingEvent, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{
			TransactionID: releaseKey,
		},
	)
	if err == nil {
		return dspan.Response(
			span,
			&UnfreezeResponse{
				UnfrozenAmount: decimal.Decimal(existingEvent.Event.Amount).Abs(),
				TransactionID:  req.TransactionID,
			},
		), nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, span.Err(err)
	}

	frozenEvent, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{
			TransactionID: req.TransactionID,
		},
	)
	if err != nil {
		return nil, span.Err(ErrFreezeNotFound)
	}

	frozenAmount := decimal.Decimal(frozenEvent.Event.Amount)

	err = s.tx.Do(
		ctx,
		"balanceservice.Unfreeze",
		func(ctx context.Context) error {
			newEvent := event.New(
				req.UserID,
				transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_RELEASE,
				event.Amount(frozenAmount.Neg()),
				releaseKey,
				nil,
				0,
			)
			if _, err := s.eventRepo.Create(
				ctx,
				dto.EventCreateRequest{
					Event: newEvent,
				},
			); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	if err := s.freezeScheduler.Cancel(
		ctx,
		dto.FreezeCancelRequest{
			TransactionID: req.TransactionID,
		},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to cancel freeze scheduler", "err", err)
	}
	if err := s.balanceCache.Invalidate(
		ctx,
		dto.BalanceCacheInvalidateRequest{UserID: req.UserID},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to invalidate balance cache", "err", err)
	}

	return dspan.Response(
		span,
		&UnfreezeResponse{
			UnfrozenAmount: frozenAmount,
			TransactionID:  req.TransactionID,
		},
	), nil
}
