package event

import (
	"time"

	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type Event struct {
	ID              ID
	UserID          user.ID
	Type            EventType
	Amount          Amount
	TransactionID   TransactionID
	SnapshotID      *SnapshotID
	FreezeExpiresAt *time.Time
	CreatedAt       time.Time
}
