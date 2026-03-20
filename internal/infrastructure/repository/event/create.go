package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.EventCreateRequest,
) (dto.EventCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	m := &models.Event{
		UserID:          req.Event.UserID.Value,
		Type:            req.Event.Type.Value,
		Amount:          req.Event.Amount.Value,
		TransactionID:   req.Event.TransactionID.Value,
		FreezeExpiresAt: req.Event.FreezeExpiresAt,
	}
	if req.Event.SnapshotID != nil {
		s := req.Event.SnapshotID.Value
		m.SnapshotID = &s
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.EventCreateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventCreateResponse{Event: *eventconvert.ModelToEvent(m)}), nil
}
