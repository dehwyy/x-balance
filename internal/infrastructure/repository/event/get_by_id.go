package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetByID(
	ctx context.Context,
	req dto.EventGetByIDRequest,
) (dto.EventGetByIDResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.GetByID", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Event
	if err := db.Where("id = ?", req.ID.Value).First(&m).Error; err != nil {
		return dto.EventGetByIDResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventGetByIDResponse{Event: *eventconvert.ModelToEvent(&m)}), nil
}
