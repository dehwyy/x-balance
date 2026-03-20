package snapshotcron

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

func (w *Worker) ensureInitialSnapshot(ctx context.Context, userID user.ID) error {
	ctx, span := dspan.Start(
		ctx,
		"snapshotcron.Worker.ensureInitialSnapshot",
		dspan.Attr("user_id", userID),
	)
	defer span.End()

	_, err := w.snapshotRepo.GetLatestByUserID(
		ctx,
		dto.SnapshotGetLatestByUserIDRequest{UserID: userID},
	)
	if err == nil {
		return nil
	}
	if !isNotFound(err) {
		return span.Err(err)
	}

	if _, err := w.snapshotRepo.Create(
		ctx,
		dto.SnapshotCreateRequest{
			UserID:  userID,
			Balance: snapshot.NewBalance(decimal.Zero),
			Version: snapshot.NewVersion(0),
		},
	); err != nil {
		return span.Err(err)
	}

	return nil
}
