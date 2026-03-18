package repository

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
)

//go:generate mockery --name=SnapshotRepository --output=../../../pkg/test/mocks --outpkg=mocks
type SnapshotRepository interface {
	Create(
		ctx context.Context,
		s *snapshot.Snapshot,
	) (*snapshot.Snapshot, error)

	GetLatestByUserID(
		ctx context.Context,
		userID string,
	) (*snapshot.Snapshot, error)

	UpdateVersion(
		ctx context.Context,
		s *snapshot.Snapshot,
	) error
}
