package repository

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
)

//go:generate mockery --name=EventRepository --output=../../../pkg/test/mocks --outpkg=mocks
type EventRepository interface {
	Create(
		ctx context.Context,
		req dto.EventCreateRequest,
	) (dto.EventCreateResponse, error)

	GetByTransactionID(
		ctx context.Context,
		req dto.EventGetByTxIDRequest,
	) (dto.EventGetByTxIDResponse, error)

	GetByID(
		ctx context.Context,
		req dto.EventGetByIDRequest,
	) (dto.EventGetByIDResponse, error)

	List(
		ctx context.Context,
		req dto.EventListRequest,
	) (dto.EventListResponse, error)

	CountSinceSnapshot(
		ctx context.Context,
		req dto.EventCountSinceSnapshotRequest,
	) (dto.EventCountSinceSnapshotResponse, error)

	SumSinceSnapshot(
		ctx context.Context,
		req dto.EventSumSinceSnapshotRequest,
	) (dto.EventSumSinceSnapshotResponse, error)
}
