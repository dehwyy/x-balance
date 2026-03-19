package gateway

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
)

//go:generate mockery --name=BalanceCache --output=../../../pkg/test/mocks --outpkg=mocks
type BalanceCache interface {
	Get(
		ctx context.Context,
		req dto.BalanceCacheGetRequest,
	) (dto.BalanceCacheGetResponse, error)

	Set(
		ctx context.Context,
		req dto.BalanceCacheSetRequest,
	) error

	Invalidate(
		ctx context.Context,
		req dto.BalanceCacheInvalidateRequest,
	) error
}
