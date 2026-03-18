package mocks

import (
	"context"

	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TxManager struct {
	mock.Mock
}

func (_m *TxManager) Do(ctx context.Context, spanName string, fn func(context.Context) error) error {
	// Execute fn directly (no real transaction needed for unit tests)
	return fn(ctx)
}

func (_m *TxManager) GetConnection(ctx context.Context) *gorm.DB {
	ret := _m.Called(ctx)
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(*gorm.DB)
}

func (_m *TxManager) GetRawConnection(ctx context.Context) *gorm.DB {
	ret := _m.Called(ctx)
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(*gorm.DB)
}

var _ txmanager.TxManager = &TxManager{}
