package balanceservice

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type FreezeRequest struct {
	UserID               user.ID
	Amount               decimal.Decimal
	TransactionID        event.TransactionID
	FreezeTimeoutSeconds int64
}

type FreezeResponse struct {
	FrozenAmount  decimal.Decimal
	TransactionID event.TransactionID
}

func (s *Service) Freeze(
	ctx context.Context,
	req *FreezeRequest,
) (*FreezeResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.Freeze",
		dspan.Attr("req", req),
	)
	defer span.End()

	existingEvent, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{TransactionID: req.TransactionID},
	)
	if err == nil {
		return dspan.Response(span, &FreezeResponse{FrozenAmount: existingEvent.Event.Amount.Value, TransactionID: req.TransactionID}), nil
	}
	if !isNotFound(err) {
		return nil, span.Err(err)
	}

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(
			ctx,
			"balanceservice.Freeze",
			func(ctx context.Context) error {
				snapshotResult, err := s.snapshotRepo.GetLatestByUserID(
					ctx,
					dto.SnapshotGetLatestByUserIDRequest{UserID: req.UserID},
				)
				if err != nil {
					return err
				}
				snap := snapshotResult.Snapshot

				userDTO, err := s.userRepo.GetByID(
					ctx,
					dto.UserGetByIDRequest{ID: req.UserID},
				)
				if err != nil {
					return err
				}
				u := userDTO.User

				sumSinceSnapshot, err := s.eventRepo.SumSinceSnapshot(
					ctx,
					dto.EventSumSinceSnapshotRequest{UserID: req.UserID, SnapshotID: snap.ID},
				)
				if err != nil {
					return err
				}

				available, _ := snap.ComputeBalance(sumSinceSnapshot.Available, sumSinceSnapshot.Frozen)
				if !u.CanDebit(available, req.Amount) {
					return ErrInsufficientFunds
				}

				if err := s.snapshotRepo.UpdateVersion(
					ctx,
					dto.SnapshotUpdateVersionRequest{Snapshot: snap},
				); err != nil {
					return err
				}

				var expiresAt *time.Time
				if req.FreezeTimeoutSeconds > 0 {
					t := time.Now().Add(time.Duration(req.FreezeTimeoutSeconds) * time.Second)
					expiresAt = &t
				}

				snapID := event.NewSnapshotID(snap.ID.Value)
				if _, err := s.eventRepo.Create(
					ctx,
					dto.EventCreateRequest{
						UserID:          req.UserID,
						Type:            event.TypeFreezeHold,
						Amount:          event.NewAmount(req.Amount),
						TransactionID:   req.TransactionID,
						SnapshotID:      &snapID,
						FreezeExpiresAt: expiresAt,
					},
				); err != nil {
					return err
				}

				return nil
			},
		)
	})
	if err != nil {
		return nil, span.Err(err)
	}

	if req.FreezeTimeoutSeconds > 0 {
		if err := s.freezeScheduler.Schedule(
			ctx,
			dto.FreezeScheduleRequest{
				TransactionID: req.TransactionID,
				TTLSeconds:    req.FreezeTimeoutSeconds,
			},
		); err != nil {
			tlog.FromContext(ctx).Error("failed to schedule freeze expiry", "err", err)
		}
	}

	if err := s.balanceCache.Invalidate(
		ctx,
		dto.BalanceCacheInvalidateRequest{UserID: req.UserID},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to invalidate balance cache", "err", err)
	}

	return dspan.Response(span, &FreezeResponse{FrozenAmount: req.Amount, TransactionID: req.TransactionID}), nil
}
