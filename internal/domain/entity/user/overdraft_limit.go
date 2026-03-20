package user

import "github.com/shopspring/decimal"

type OverdraftLimit struct {
	Value decimal.Decimal
}

func NewOverdraftLimit(v decimal.Decimal) OverdraftLimit {
	return OverdraftLimit{Value: v}
}
