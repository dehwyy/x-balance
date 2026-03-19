package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type SnapshotRepository struct {
	mock.Mock
}

func (_m *SnapshotRepository) Create(ctx context.Context, req dto.SnapshotCreateRequest) (dto.SnapshotCreateResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.SnapshotCreateResponse), ret.Error(1)
}

func (_m *SnapshotRepository) GetLatestByUserID(ctx context.Context, req dto.SnapshotGetLatestByUserIDRequest) (dto.SnapshotGetLatestByUserIDResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.SnapshotGetLatestByUserIDResponse), ret.Error(1)
}

func (_m *SnapshotRepository) UpdateVersion(ctx context.Context, req dto.SnapshotUpdateVersionRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

var _ repository.SnapshotRepository = &SnapshotRepository{}
