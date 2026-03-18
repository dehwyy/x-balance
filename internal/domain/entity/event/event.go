package event

import "time"

type Event struct {
	ID              ID
	UserID          string
	Type            EventType
	Amount          Amount
	TransactionID   TransactionID
	SnapshotID      *string
	FreezeExpiresAt *time.Time
	CreatedAt       time.Time
}
