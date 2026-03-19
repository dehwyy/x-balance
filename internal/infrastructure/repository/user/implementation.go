package userrepo

import (
	"go.uber.org/fx"

	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/domain/repository"
)

var _ repository.UserRepository = &Implementation{}

type Opts struct {
	fx.In
	TX txmanager.TxManager
}

type Implementation struct {
	tx txmanager.TxManager
}

func New(opts Opts) *Implementation {
	return &Implementation{tx: opts.TX}
}
