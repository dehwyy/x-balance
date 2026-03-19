package snapshotrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Create(
	ctx context.Context,
	s *snapshot.Snapshot,
) (*snapshot.Snapshot, error) {
	ctx, span := dspan.Start(ctx, "snapshotrepo.Create")
	defer span.End()

	m := &models.Snapshot{
		UserID:  s.UserID,
		Balance: s.Balance.Value,
		Version: s.Version.Value,
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return nil, span.Err(err)
	}

	return snapshotconvert.ModelToSnapshot(m), nil
}
