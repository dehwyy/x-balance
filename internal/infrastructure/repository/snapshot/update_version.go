package snapshotrepo

import (
	"context"
	"errors"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

var ErrVersionConflict = errors.New("snapshot version conflict")

func (impl *Implementation) UpdateVersion(
	ctx context.Context,
	req dto.SnapshotUpdateVersionRequest,
) error {
	ctx, span := dspan.Start(ctx, "snapshotrepo.Implementation.UpdateVersion", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	result := db.Model(&models.Snapshot{}).
		Where("id = ? AND version = ?", string(req.Snapshot.ID), int64(req.Snapshot.Version)).
		Update("version", int64(req.Snapshot.Version)+1)

	if result.Error != nil {
		return span.Err(result.Error)
	}

	if result.RowsAffected == 0 {
		return span.Err(ErrVersionConflict)
	}

	return nil
}
