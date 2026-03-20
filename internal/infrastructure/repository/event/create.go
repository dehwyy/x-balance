package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.EventCreateRequest,
) (dto.EventCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	// Import from event model converter
	m := eventconvert.EventToModel(&req.Event)

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.EventCreateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventCreateResponse{Event: *eventconvert.ModelToEvent(m)}), nil
}
