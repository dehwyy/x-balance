package freezeexpiry

import (
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Worker struct {
	client         *redis.Client
	balanceservice *balanceservice.Service
}

type WorkerOpts struct {
	fx.In
	Client         *redis.Client
	BalanceService *balanceservice.Service
}

func New(opts WorkerOpts) *Worker {
	return &Worker{
		client:         opts.Client,
		balanceservice: opts.BalanceService,
	}
}
