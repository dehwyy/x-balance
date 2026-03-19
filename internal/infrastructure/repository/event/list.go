package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) List(
	ctx context.Context,
	req dto.EventListRequest,
) (dto.EventListResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.List", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx).Where("user_id = ?", req.UserID.Value)

	if req.From != nil {
		db = db.Where("created_at >= ?", req.From)
	}
	if req.To != nil {
		db = db.Where("created_at <= ?", req.To)
	}

	var total int64
	if err := db.Model(&models.Event{}).Count(&total).Error; err != nil {
		return dto.EventListResponse{}, span.Err(err)
	}

	limit := req.Pagination.Limit()
	offset := req.Pagination.Offset()
	if limit > 0 {
		db = db.Limit(limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}

	var ms []models.Event
	if err := db.Order("created_at DESC").Find(&ms).Error; err != nil {
		return dto.EventListResponse{}, span.Err(err)
	}

	events := make([]event.Event, len(ms))
	for i, m := range ms {
		m := m
		events[i] = *eventconvert.ModelToEvent(&m)
	}

	response := dto.EventListResponse{Events: events, Total: total}
	span.WithAttribute("response", response)
	return response, nil
}
