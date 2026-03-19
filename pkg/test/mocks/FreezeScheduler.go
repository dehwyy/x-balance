package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/stretchr/testify/mock"
)

type FreezeScheduler struct {
	mock.Mock
}

func (_m *FreezeScheduler) Schedule(ctx context.Context, req dto.FreezeScheduleRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

func (_m *FreezeScheduler) Cancel(ctx context.Context, req dto.FreezeCancelRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

var _ gateway.FreezeScheduler = &FreezeScheduler{}
