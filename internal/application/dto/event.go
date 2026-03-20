package dto

import (
	"time"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/pkg/storage"
	"github.com/shopspring/decimal"
)

type EventCreateRequest struct {
	Event event.Event
}

type EventCreateResponse struct {
	Event event.Event
}

type EventGetByIDRequest struct {
	ID event.ID
}

type EventGetByIDResponse struct {
	Event event.Event
}

type EventGetByTxIDRequest struct {
	TransactionID event.TransactionID
}

type EventGetByTxIDResponse struct {
	Event event.Event
}

type EventListRequest struct {
	UserID     user.ID
	Pagination storage.Pagination
	From       *time.Time
	To         *time.Time
}

type EventListResponse struct {
	Events []event.Event
	Total  int64
}

type EventCountSinceSnapshotRequest struct {
	UserID     user.ID
	SnapshotID snapshot.ID
}

type EventCountSinceSnapshotResponse struct {
	Count int64
}

type EventSumSinceSnapshotRequest struct {
	UserID     user.ID
	SnapshotID snapshot.ID
}

type EventSumSinceSnapshotResponse struct {
	Available decimal.Decimal
	Frozen    decimal.Decimal
}
