package snapshot

import "github.com/shopspring/decimal"

type Balance struct {
	Value decimal.Decimal
}

func NewBalance(v decimal.Decimal) Balance { return Balance{Value: v} }
