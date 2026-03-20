package snapshotrepo

import (
	"context"

	"gorm.io/gorm"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetLatestByUserID(
	ctx context.Context,
	req dto.SnapshotGetLatestByUserIDRequest,
) (dto.SnapshotGetLatestByUserIDResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"snapshotrepo.Implementation.GetLatestByUserID",
		dspan.Attr("req", req),
	)
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Snapshot
	if err := db.Where("user_id = ?", string(req.UserID)).
		Order("created_at DESC").
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.SnapshotGetLatestByUserIDResponse{}, repository.ErrNotFound
		}
		return dto.SnapshotGetLatestByUserIDResponse{}, span.Err(err)
	}

	return dspan.Response(
		span,
		dto.SnapshotGetLatestByUserIDResponse{
			Snapshot: *snapshotconvert.ModelToSnapshot(&m),
		},
	), nil
}
