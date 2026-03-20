package event

import "github.com/shopspring/decimal"

type Amount struct {
	Value decimal.Decimal
}

func NewAmount(v decimal.Decimal) Amount { return Amount{Value: v} }
