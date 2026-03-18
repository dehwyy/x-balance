package snapshotrepo

import (
	"context"
	"errors"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

var ErrVersionConflict = errors.New("snapshot version conflict")

func (impl *Implementation) UpdateVersion(
	ctx context.Context,
	s *snapshot.Snapshot,
) error {
	ctx, span := dspan.Start(ctx, "snapshotrepo.UpdateVersion")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	result := db.Model(&models.Snapshot{}).
		Where("id = ? AND version = ?", s.ID.Value, s.Version.Value).
		Update("version", s.Version.Value+1)

	if result.Error != nil {
		return span.Err(result.Error)
	}

	if result.RowsAffected == 0 {
		return span.Err(ErrVersionConflict)
	}

	return nil
}
