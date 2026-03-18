package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models/convert"
)

func (impl *Implementation) GetByID(
	ctx context.Context,
	id event.ID,
) (*event.Event, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.GetById")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Event
	if err := db.Where("id = ?", id.Value).First(&m).Error; err != nil {
		return nil, span.Err(err)
	}

	return convert.EventToDomain(&m), nil
}
