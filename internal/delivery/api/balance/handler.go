package balancehandler

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	BalanceService *balanceservice.Service
}

type Handler struct {
	balancev1.UnimplementedBalanceServiceServer
	balanceservice *balanceservice.Service
}

func New(opts Opts) *Handler {
	return &Handler{
		balanceservice: opts.BalanceService,
	}
}
