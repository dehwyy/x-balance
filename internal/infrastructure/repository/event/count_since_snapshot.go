package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) CountSinceSnapshot(
	ctx context.Context,
	req dto.EventCountSinceSnapshotRequest,
) (dto.EventCountSinceSnapshotResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.CountSinceSnapshot", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)

	var snapshotModel models.Snapshot
	if err := db.Where("id = ?", req.SnapshotID.Value).First(&snapshotModel).Error; err != nil {
		return dto.EventCountSinceSnapshotResponse{}, span.Err(err)
	}

	var count int64
	if err := db.Model(&models.Event{}).
		Where("user_id = ? AND created_at > ?", req.UserID.Value, snapshotModel.CreatedAt).
		Count(&count).Error; err != nil {
		return dto.EventCountSinceSnapshotResponse{}, span.Err(err)
	}

	response := dto.EventCountSinceSnapshotResponse{Count: count}
	span.WithAttribute("response", response)
	return response, nil
}
