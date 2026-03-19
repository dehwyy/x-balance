package snapshotrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetLatestByUserID(
	ctx context.Context,
	userID string,
) (*snapshot.Snapshot, error) {
	ctx, span := dspan.Start(ctx, "snapshotrepo.GetLatestByUserId")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Snapshot
	if err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&m).Error; err != nil {
		return nil, span.Err(err)
	}

	return snapshotconvert.ModelToSnapshot(&m), nil
}
