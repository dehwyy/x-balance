package balanceservice

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
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
	ctx, span := dspan.Start(ctx, "balanceservice.Service.Freeze", dspan.Attr("req", req))
	defer span.End()

	existingResp, err := s.eventRepo.GetByTransactionID(ctx, dto.EventGetByTxIDRequest{TransactionID: req.TransactionID})
	if err == nil {
		response := &FreezeResponse{FrozenAmount: existingResp.Event.Amount.Value, TransactionID: req.TransactionID}
		span.WithAttribute("response", response)
		return response, nil
	}
	if !isNotFound(err) {
		return nil, span.Err(err)
	}

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(ctx, "balanceservice.Freeze", func(ctx context.Context) error {
			snapResp, err := s.snapshotRepo.GetLatestByUserID(ctx, dto.SnapshotGetLatestByUserIDRequest{UserID: req.UserID})
			if err != nil {
				return err
			}
			snap := snapResp.Snapshot

			userResp, err := s.userRepo.GetByID(ctx, dto.UserGetByIDRequest{ID: req.UserID})
			if err != nil {
				return err
			}
			u := userResp.User

			sumResp, err := s.eventRepo.SumSinceSnapshot(ctx, dto.EventSumSinceSnapshotRequest{UserID: req.UserID, SnapshotID: snap.ID})
			if err != nil {
				return err
			}

			available := snap.Balance.Value.Add(sumResp.Available).Sub(sumResp.Frozen)
			minAllowed := u.OverdraftLimit.Value.Neg()
			if available.Sub(req.Amount).LessThan(minAllowed) {
				return ErrInsufficientFunds
			}

			if err := s.snapshotRepo.UpdateVersion(ctx, dto.SnapshotUpdateVersionRequest{Snapshot: snap}); err != nil {
				return err
			}

			var expiresAt *time.Time
			if req.FreezeTimeoutSeconds > 0 {
				t := time.Now().Add(time.Duration(req.FreezeTimeoutSeconds) * time.Second)
				expiresAt = &t
			}

			snapID := event.SnapshotID{Value: snap.ID.Value}
			if _, err := s.eventRepo.Create(ctx, dto.EventCreateRequest{
				UserID:          req.UserID,
				Type:            event.TypeFreezeHold,
				Amount:          event.Amount{Value: req.Amount},
				TransactionID:   req.TransactionID,
				SnapshotID:      &snapID,
				FreezeExpiresAt: expiresAt,
			}); err != nil {
				return err
			}

			return nil
		})
	})
	if err != nil {
		return nil, span.Err(err)
	}

	if req.FreezeTimeoutSeconds > 0 {
		_ = s.freezeScheduler.Schedule(ctx, dto.FreezeScheduleRequest{
			TransactionID: req.TransactionID,
			TTLSeconds:    req.FreezeTimeoutSeconds,
		})
	}

	_ = s.balanceCache.Invalidate(ctx, dto.BalanceCacheInvalidateRequest{UserID: req.UserID})

	response := &FreezeResponse{FrozenAmount: req.Amount, TransactionID: req.TransactionID}
	span.WithAttribute("response", response)
	return response, nil
}
