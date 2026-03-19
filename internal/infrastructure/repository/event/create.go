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
		UserID:          req.UserID.Value,
		Type:            req.Type.Value,
		Amount:          req.Amount.Value,
		TransactionID:   req.TransactionID.Value,
		FreezeExpiresAt: req.FreezeExpiresAt,
	}
	if req.SnapshotID != nil {
		s := req.SnapshotID.Value
		m.SnapshotID = &s
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.EventCreateResponse{}, span.Err(err)
	}

	response := dto.EventCreateResponse{Event: *eventconvert.ModelToEvent(m)}
	span.WithAttribute("response", response)
	return response, nil
}
