package gateway

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockery --name=BalanceCache --output=../../../pkg/test/mocks --outpkg=mocks
type BalanceCache interface {
	Get(
		ctx context.Context,
		userID string,
	) (available decimal.Decimal, frozen decimal.Decimal, found bool, err error)

	Set(
		ctx context.Context,
		userID string,
		available decimal.Decimal,
		frozen decimal.Decimal,
	) error

	Invalidate(
		ctx context.Context,
		userID string,
	) error
}
