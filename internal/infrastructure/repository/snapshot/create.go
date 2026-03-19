package snapshotrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	snapshotconvert "github.com/dehwyy/x-balance/internal/domain/entity/snapshot/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.SnapshotCreateRequest,
) (dto.SnapshotCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "snapshotrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	m := &models.Snapshot{
		UserID:  req.UserID.Value,
		Balance: req.Balance.Value,
		Version: req.Version.Value,
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.SnapshotCreateResponse{}, span.Err(err)
	}

	response := dto.SnapshotCreateResponse{Snapshot: *snapshotconvert.ModelToSnapshot(m)}
	span.WithAttribute("response", response)
	return response, nil
}