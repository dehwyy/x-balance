package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type EventRepository struct {
	mock.Mock
}

func (_m *EventRepository) Create(ctx context.Context, req dto.EventCreateRequest) (dto.EventCreateResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventCreateResponse), ret.Error(1)
}

func (_m *EventRepository) GetByTransactionID(ctx context.Context, req dto.EventGetByTxIDRequest) (dto.EventGetByTxIDResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventGetByTxIDResponse), ret.Error(1)
}

func (_m *EventRepository) GetByID(ctx context.Context, req dto.EventGetByIDRequest) (dto.EventGetByIDResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventGetByIDResponse), ret.Error(1)
}

func (_m *EventRepository) List(ctx context.Context, req dto.EventListRequest) (dto.EventListResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventListResponse), ret.Error(1)
}

func (_m *EventRepository) CountSinceSnapshot(ctx context.Context, req dto.EventCountSinceSnapshotRequest) (dto.EventCountSinceSnapshotResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventCountSinceSnapshotResponse), ret.Error(1)
}

func (_m *EventRepository) SumSinceSnapshot(ctx context.Context, req dto.EventSumSinceSnapshotRequest) (dto.EventSumSinceSnapshotResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.EventSumSinceSnapshotResponse), ret.Error(1)
}

var _ repository.EventRepository = &EventRepository{}
