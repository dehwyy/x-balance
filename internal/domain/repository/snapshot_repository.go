package repository

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
)

//go:generate mockery --name=SnapshotRepository --output=../../../pkg/test/mocks --outpkg=mocks
type SnapshotRepository interface {
	Create(
		ctx context.Context,
		req dto.SnapshotCreateRequest,
	) (dto.SnapshotCreateResponse, error)

	GetLatestByUserID(
		ctx context.Context,
		req dto.SnapshotGetLatestByUserIDRequest,
	) (dto.SnapshotGetLatestByUserIDResponse, error)

	UpdateVersion(
		ctx context.Context,
		req dto.SnapshotUpdateVersionRequest,
	) error
}
