package gateway

import "context"

//go:generate mockery --name=FreezeScheduler --output=../../../pkg/test/mocks --outpkg=mocks
type FreezeScheduler interface {
	Schedule(
		ctx context.Context,
		txID string,
		ttlSeconds int64,
	) error

	Cancel(
		ctx context.Context,
		txID string,
	) error
}
