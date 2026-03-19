package gateway

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
)

//go:generate mockery --name=FreezeScheduler --output=../../../pkg/test/mocks --outpkg=mocks
type FreezeScheduler interface {
	Schedule(
		ctx context.Context,
		req dto.FreezeScheduleRequest,
	) error

	Cancel(
		ctx context.Context,
		req dto.FreezeCancelRequest,
	) error
}
