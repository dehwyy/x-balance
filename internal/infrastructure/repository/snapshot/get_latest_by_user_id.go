package snapshotrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetLatestByUserID(
	ctx context.Context,
	req dto.SnapshotGetLatestByUserIDRequest,
) (dto.SnapshotGetLatestByUserIDResponse, error) {
	ctx, span := dspan.Start(ctx, "snapshotrepo.Implementation.GetLatestByUserID", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Snapshot
	if err := db.Where("user_id = ?", req.UserID.Value).
		Order("created_at DESC").
		First(&m).Error; err != nil {
		return dto.SnapshotGetLatestByUserIDResponse{}, span.Err(err)
	}

	response := dto.SnapshotGetLatestByUserIDResponse{Snapshot: *snapshotconvert.ModelToSnapshot(&m)}
	span.WithAttribute("response", response)
	return response, nil
}