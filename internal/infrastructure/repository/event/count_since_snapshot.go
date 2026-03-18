package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) CountSinceSnapshot(
	ctx context.Context,
	userID string,
	snapshotID snapshot.ID,
) (int64, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.CountSinceSnapshot")
	defer span.End()

	db := impl.tx.GetConnection(ctx)

	var snapshotModel models.Snapshot
	if err := db.Where("id = ?", snapshotID.Value).First(&snapshotModel).Error; err != nil {
		return 0, span.Err(err)
	}

	var count int64
	if err := db.Model(&models.Event{}).
		Where("user_id = ? AND created_at > ?", userID, snapshotModel.CreatedAt).
		Count(&count).Error; err != nil {
		return 0, span.Err(err)
	}

	return count, nil
}
