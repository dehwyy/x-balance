package eventconvert

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToEvent(m *models.Event) *event.Event {
	e := &event.Event{
		ID:              event.ID{Value: m.ID},
		UserID:          user.ID{Value: m.UserID},
		Type:            event.EventType{Value: m.Type},
		Amount:          event.Amount{Value: m.Amount},
		TransactionID:   event.TransactionID{Value: m.TransactionID},
		FreezeExpiresAt: m.FreezeExpiresAt,
		CreatedAt:       m.CreatedAt,
	}
	if m.SnapshotID != nil {
		snapID := event.SnapshotID{Value: *m.SnapshotID}
		e.SnapshotID = &snapID
	}
	return e
}

func EventToModel(e *event.Event) *models.Event {
	m := &models.Event{
		ID:              e.ID.Value,
		UserID:          e.UserID.Value,
		Type:            e.Type.Value,
		Amount:          e.Amount.Value,
		TransactionID:   e.TransactionID.Value,
		FreezeExpiresAt: e.FreezeExpiresAt,
		CreatedAt:       e.CreatedAt,
	}
	if e.SnapshotID != nil {
		s := e.SnapshotID.Value
		m.SnapshotID = &s
	}
	return m
}
