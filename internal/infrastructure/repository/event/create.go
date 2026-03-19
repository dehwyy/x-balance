package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Create(
	ctx context.Context,
	e *event.Event,
) (*event.Event, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Create")
	defer span.End()

	m := &models.Event{
		UserID:          e.UserID,
		Type:            e.Type.Value,
		Amount:          e.Amount.Value,
		TransactionID:   e.TransactionID.Value,
		SnapshotID:      e.SnapshotID,
		FreezeExpiresAt: e.FreezeExpiresAt,
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return nil, span.Err(err)
	}

	return eventconvert.ModelToEvent(m), nil
}
