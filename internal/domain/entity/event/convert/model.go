package eventconvert

import (
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

var typeToDB = map[transactionv1.TransactionType]string{
	transactionv1.TransactionType_TRANSACTION_TYPE_CREDIT:         "credit",
	transactionv1.TransactionType_TRANSACTION_TYPE_DEBIT:          "debit",
	transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_HOLD:    "freeze_hold",
	transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_RELEASE: "freeze_release",
}

var dbToType = map[string]transactionv1.TransactionType{
	"credit":         transactionv1.TransactionType_TRANSACTION_TYPE_CREDIT,
	"debit":          transactionv1.TransactionType_TRANSACTION_TYPE_DEBIT,
	"freeze_hold":    transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_HOLD,
	"freeze_release": transactionv1.TransactionType_TRANSACTION_TYPE_FREEZE_RELEASE,
}

func ModelToEvent(m *models.Event) *event.Event {
	e := &event.Event{
		ID:              event.ID(m.ID),
		UserID:          user.ID(m.UserID),
		Type:            dbToType[m.Type],
		Amount:          event.Amount(m.Amount),
		TransactionID:   event.TransactionID(m.TransactionID),
		FreezeExpiresAt: m.FreezeExpiresAt,
		CreatedAt:       m.CreatedAt,
	}
	if m.SnapshotID != nil {
		snapID := event.SnapshotID(*m.SnapshotID)
		e.SnapshotID = &snapID
	}
	return e
}

func EventToModel(e *event.Event) *models.Event {
	m := &models.Event{
		ID:              string(e.ID),
		UserID:          string(e.UserID),
		Type:            typeToDB[e.Type],
		Amount:          decimal.Decimal(e.Amount),
		TransactionID:   string(e.TransactionID),
		FreezeExpiresAt: e.FreezeExpiresAt,
		CreatedAt:       e.CreatedAt,
	}
	if e.SnapshotID != nil {
		s := string(*e.SnapshotID)
		m.SnapshotID = &s
	}
	return m
}
