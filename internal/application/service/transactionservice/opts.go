package transactionservice

import (
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	EventRepo repository.EventRepository
}

type Service struct {
	eventRepo repository.EventRepository
}

func New(opts Opts) *Service {
	return &Service{
		eventRepo: opts.EventRepo,
	}
}
