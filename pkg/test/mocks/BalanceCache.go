package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/stretchr/testify/mock"
)

type BalanceCache struct {
	mock.Mock
}

func (_m *BalanceCache) Get(ctx context.Context, req dto.BalanceCacheGetRequest) (dto.BalanceCacheGetResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.BalanceCacheGetResponse), ret.Error(1)
}

func (_m *BalanceCache) Set(ctx context.Context, req dto.BalanceCacheSetRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

func (_m *BalanceCache) Invalidate(ctx context.Context, req dto.BalanceCacheInvalidateRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

var _ gateway.BalanceCache = &BalanceCache{}
