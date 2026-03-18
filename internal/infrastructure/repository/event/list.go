package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models/convert"
)

func (impl *Implementation) List(
	ctx context.Context,
	req repository.ListEventsRequest,
) ([]*event.Event, int64, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.List")
	defer span.End()

	db := impl.tx.GetConnection(ctx).Where("user_id = ?", req.UserID)

	if req.From != nil {
		db = db.Where("created_at >= ?", req.From)
	}
	if req.To != nil {
		db = db.Where("created_at <= ?", req.To)
	}

	var total int64
	if err := db.Model(&models.Event{}).Count(&total).Error; err != nil {
		return nil, 0, span.Err(err)
	}

	if req.Limit > 0 {
		db = db.Limit(req.Limit)
	}
	if req.Offset > 0 {
		db = db.Offset(req.Offset)
	}

	var ms []models.Event
	if err := db.Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, span.Err(err)
	}

	events := make([]*event.Event, len(ms))
	for i, m := range ms {
		m := m
		events[i] = convert.EventToDomain(&m)
	}

	return events, total, nil
}
