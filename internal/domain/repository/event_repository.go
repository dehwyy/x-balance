package repository

import (
	"context"
	"time"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/shopspring/decimal"
)

type ListEventsRequest struct {
	UserID string
	Limit  int
	Offset int
	From   *time.Time
	To     *time.Time
}

//go:generate mockery --name=EventRepository --output=../../../pkg/test/mocks --outpkg=mocks
type EventRepository interface {
	Create(
		ctx context.Context,
		e *event.Event,
	) (*event.Event, error)

	GetByTransactionID(
		ctx context.Context,
		txID event.TransactionID,
	) (*event.Event, error)

	GetByID(
		ctx context.Context,
		id event.ID,
	) (*event.Event, error)

	List(
		ctx context.Context,
		req ListEventsRequest,
	) ([]*event.Event, int64, error)

	CountSinceSnapshot(
		ctx context.Context,
		userID string,
		snapshotID snapshot.ID,
	) (int64, error)

	SumSinceSnapshot(
		ctx context.Context,
		userID string,
		snapshotID snapshot.ID,
	) (decimal.Decimal, decimal.Decimal, error)
}
