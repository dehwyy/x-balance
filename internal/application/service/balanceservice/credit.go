package balanceservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
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
	ctx, span := dspan.Start(
		ctx,
		"balanceservice.Service.Credit",
		dspan.Attr("req", req),
	)
	defer span.End()

	_, err := s.eventRepo.GetByTransactionID(
		ctx,
		dto.EventGetByTxIDRequest{TransactionID: req.TransactionID},
	)
	if err == nil {
		bal, _, err := s.computeBalance(ctx, req.UserID)
		if err != nil {
			return nil, span.Err(err)
		}
		return dspan.Response(span, &CreditResponse{NewBalance: bal, TransactionID: req.TransactionID}), nil
	}
	if !isNotFound(err) {
		return nil, span.Err(err)
	}

	var newBalance decimal.Decimal

	err = s.withRetry(ctx, func(ctx context.Context) error {
		return s.tx.Do(
			ctx,
			"balanceservice.Credit",
			func(ctx context.Context) error {
				snapshotResult, err := s.snapshotRepo.GetLatestByUserID(
					ctx,
					dto.SnapshotGetLatestByUserIDRequest{UserID: req.UserID},
				)
				if err != nil {
					return err
				}
				snap := snapshotResult.Snapshot

				if err := s.snapshotRepo.UpdateVersion(
					ctx,
					dto.SnapshotUpdateVersionRequest{Snapshot: snap},
				); err != nil {
					return err
				}

				snapID := event.NewSnapshotID(snap.ID.Value)
				if _, err := s.eventRepo.Create(
					ctx,
					dto.EventCreateRequest{
						UserID:        req.UserID,
						Type:          event.TypeCredit,
						Amount:        event.NewAmount(req.Amount),
						TransactionID: req.TransactionID,
						SnapshotID:    &snapID,
					},
				); err != nil {
					return err
				}

				sumSinceSnapshot, err := s.eventRepo.SumSinceSnapshot(
					ctx,
					dto.EventSumSinceSnapshotRequest{UserID: req.UserID, SnapshotID: snap.ID},
				)
				if err != nil {
					return err
				}
				newBalance, _ = snap.ComputeBalance(sumSinceSnapshot.Available, sumSinceSnapshot.Frozen)
				return nil
			},
		)
	})
	if err != nil {
		return nil, span.Err(err)
	}

	if err := s.balanceCache.Invalidate(
		ctx,
		dto.BalanceCacheInvalidateRequest{UserID: req.UserID},
	); err != nil {
		tlog.FromContext(ctx).Error("failed to invalidate balance cache", "err", err)
	}
	if err := s.maybeCreateSnapshot(ctx, req.UserID); err != nil {
		tlog.FromContext(ctx).Error("failed to maybe create snapshot", "err", err)
	}

	return dspan.Response(span, &CreditResponse{NewBalance: newBalance, TransactionID: req.TransactionID}), nil
}

func isNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
