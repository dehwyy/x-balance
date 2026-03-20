package snapshotcron

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

func (w *Worker) createSnapshotForUser(ctx context.Context, userID user.ID) error {
	ctx, span := dspan.Start(
		ctx,
		"snapshotcron.Worker.createSnapshotForUser",
		dspan.Attr("user_id", userID),
	)
	defer span.End()

	if err := w.ensureInitialSnapshot(ctx, userID); err != nil {
		return span.Err(err)
	}

	snapshotResult, err := w.snapshotRepo.GetLatestByUserID(
		ctx,
		dto.SnapshotGetLatestByUserIDRequest{UserID: userID},
	)
	if err != nil {
		return span.Err(err)
	}
	snap := snapshotResult.Snapshot

	countResult, err := w.eventRepo.CountSinceSnapshot(
		ctx,
		dto.EventCountSinceSnapshotRequest{
			UserID:     userID,
			SnapshotID: snap.ID,
		},
	)
	if err != nil {
		return span.Err(err)
	}

	if countResult.Count == 0 {
		return nil
	}

	sumSinceSnapshot, err := w.eventRepo.SumSinceSnapshot(
		ctx,
		dto.EventSumSinceSnapshotRequest{
			UserID:     userID,
			SnapshotID: snap.ID,
		},
	)
	if err != nil {
		return span.Err(err)
	}

	available, _ := snap.ComputeBalance(sumSinceSnapshot.Available, sumSinceSnapshot.Frozen)

	if _, err := w.snapshotRepo.Create(
		ctx,
		dto.SnapshotCreateRequest{
			UserID:  userID,
			Balance: snapshot.Balance(available),
			Version: snapshot.Version(0),
		},
	); err != nil {
		return span.Err(err)
	}

	return nil
}
