package snapshotrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.SnapshotCreateRequest,
) (dto.SnapshotCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "snapshotrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	snap := &snapshot.Snapshot{
		UserID:  req.UserID,
		Balance: req.Balance,
		Version: req.Version,
	}
	m := snapshotconvert.SnapshotToModel(snap)

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.SnapshotCreateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.SnapshotCreateResponse{Snapshot: *snapshotconvert.ModelToSnapshot(m)}), nil
}
