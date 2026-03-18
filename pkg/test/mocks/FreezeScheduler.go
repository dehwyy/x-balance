package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/stretchr/testify/mock"
)

type FreezeScheduler struct {
	mock.Mock
}

func (_m *FreezeScheduler) Schedule(ctx context.Context, txID string, ttlSeconds int64) error {
	ret := _m.Called(ctx, txID, ttlSeconds)
	return ret.Error(0)
}

func (_m *FreezeScheduler) Cancel(ctx context.Context, txID string) error {
	ret := _m.Called(ctx, txID)
	return ret.Error(0)
}

var _ gateway.FreezeScheduler = &FreezeScheduler{}
