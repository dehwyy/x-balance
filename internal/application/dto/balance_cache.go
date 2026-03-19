package dto

import (
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type BalanceCacheGetRequest struct {
	UserID user.ID
}

type BalanceCacheGetResponse struct {
	Available decimal.Decimal
	Frozen    decimal.Decimal
	Found     bool
}

type BalanceCacheSetRequest struct {
	UserID    user.ID
	Available decimal.Decimal
	Frozen    decimal.Decimal
}

type BalanceCacheInvalidateRequest struct {
	UserID user.ID
}
