package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type EventRepository struct {
	mock.Mock
}

func (_m *EventRepository) Create(ctx context.Context, e *event.Event) (*event.Event, error) {
	ret := _m.Called(ctx, e)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*event.Event), ret.Error(1)
}

func (_m *EventRepository) GetByTransactionID(ctx context.Context, txID event.TransactionID) (*event.Event, error) {
	ret := _m.Called(ctx, txID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*event.Event), ret.Error(1)
}

func (_m *EventRepository) GetByID(ctx context.Context, id event.ID) (*event.Event, error) {
	ret := _m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*event.Event), ret.Error(1)
}

func (_m *EventRepository) List(ctx context.Context, req repository.ListEventsRequest) ([]*event.Event, int64, error) {
	ret := _m.Called(ctx, req)
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(int64), ret.Error(2)
	}
	return ret.Get(0).([]*event.Event), ret.Get(1).(int64), ret.Error(2)
}

func (_m *EventRepository) CountSinceSnapshot(ctx context.Context, userID string, snapshotID snapshot.ID) (int64, error) {
	ret := _m.Called(ctx, userID, snapshotID)
	return ret.Get(0).(int64), ret.Error(1)
}

func (_m *EventRepository) SumSinceSnapshot(ctx context.Context, userID string, snapshotID snapshot.ID) (decimal.Decimal, decimal.Decimal, error) {
	ret := _m.Called(ctx, userID, snapshotID)
	return ret.Get(0).(decimal.Decimal), ret.Get(1).(decimal.Decimal), ret.Error(2)
}

var _ repository.EventRepository = &EventRepository{}
