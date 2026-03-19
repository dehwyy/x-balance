package modules

import (
	balancehandler "github.com/dehwyy/x-balance/internal/delivery/api/balance"
	transactionhandler "github.com/dehwyy/x-balance/internal/delivery/api/transaction"
	userhandler "github.com/dehwyy/x-balance/internal/delivery/api/user"
	"go.uber.org/fx"
)

var DeliveryModule = fx.Options(
	fx.Provide(
		userhandler.New,
		balancehandler.New,
		transactionhandler.New,
	),
)
