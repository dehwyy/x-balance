package userservice

import (
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TX       txmanager.TxManager
	UserRepo repository.UserRepository
}

type Service struct {
	tx       txmanager.TxManager
	userRepo repository.UserRepository
}

func New(opts Opts) *Service {
	return &Service{
		tx:       opts.TX,
		userRepo: opts.UserRepo,
	}
}
