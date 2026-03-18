package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type SnapshotRepository struct {
	mock.Mock
}

func (_m *SnapshotRepository) Create(ctx context.Context, s *snapshot.Snapshot) (*snapshot.Snapshot, error) {
	ret := _m.Called(ctx, s)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*snapshot.Snapshot), ret.Error(1)
}

func (_m *SnapshotRepository) GetLatestByUserID(ctx context.Context, userID string) (*snapshot.Snapshot, error) {
	ret := _m.Called(ctx, userID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*snapshot.Snapshot), ret.Error(1)
}

func (_m *SnapshotRepository) UpdateVersion(ctx context.Context, s *snapshot.Snapshot) error {
	ret := _m.Called(ctx, s)
	return ret.Error(0)
}

var _ repository.SnapshotRepository = &SnapshotRepository{}
