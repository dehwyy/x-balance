package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type BalanceCache struct {
	mock.Mock
}

func (_m *BalanceCache) Get(ctx context.Context, userID string) (decimal.Decimal, decimal.Decimal, bool, error) {
	ret := _m.Called(ctx, userID)
	return ret.Get(0).(decimal.Decimal), ret.Get(1).(decimal.Decimal), ret.Bool(2), ret.Error(3)
}

func (_m *BalanceCache) Set(ctx context.Context, userID string, available decimal.Decimal, frozen decimal.Decimal) error {
	ret := _m.Called(ctx, userID, available, frozen)
	return ret.Error(0)
}

func (_m *BalanceCache) Invalidate(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)
	return ret.Error(0)
}

var _ gateway.BalanceCache = &BalanceCache{}
