package eventconvert

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToEvent(m *models.Event) *event.Event {
	return &event.Event{
		ID:              event.ID{Value: m.ID},
		UserID:          m.UserID,
		Type:            event.EventType{Value: m.Type},
		Amount:          event.Amount{Value: m.Amount},
		TransactionID:   event.TransactionID{Value: m.TransactionID},
		SnapshotID:      m.SnapshotID,
		FreezeExpiresAt: m.FreezeExpiresAt,
		CreatedAt:       m.CreatedAt,
	}
}

func EventToModel(e *event.Event) *models.Event {
	return &models.Event{
		ID:              e.ID.Value,
		UserID:          e.UserID,
		Type:            e.Type.Value,
		Amount:          e.Amount.Value,
		TransactionID:   e.TransactionID.Value,
		SnapshotID:      e.SnapshotID,
		FreezeExpiresAt: e.FreezeExpiresAt,
		CreatedAt:       e.CreatedAt,
	}
}
