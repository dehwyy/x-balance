package event

import (
	"time"

	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
)

type Event struct {
	ID              ID
	UserID          user.ID
	Type            transactionv1.TransactionType
	Amount          Amount
	TransactionID   TransactionID
	SnapshotID      *SnapshotID
	FreezeExpiresAt *time.Time
	CreatedAt       time.Time
}

func New(
	userID user.ID,
	eventType transactionv1.TransactionType,
	amount Amount,
	transactionID TransactionID,
	snapshotID *SnapshotID,
	freezeTimeoutSeconds int64,
) Event {
	var expiresAt *time.Time
	if freezeTimeoutSeconds > 0 {
		t := time.Now().Add(time.Duration(freezeTimeoutSeconds) * time.Second)
		expiresAt = &t
	}
	return Event{
		UserID:          userID,
		Type:            eventType,
		Amount:          amount,
		TransactionID:   transactionID,
		SnapshotID:      snapshotID,
		FreezeExpiresAt: expiresAt,
	}
}
