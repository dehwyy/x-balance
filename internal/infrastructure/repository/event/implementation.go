package eventrepo

import (
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/domain/repository"
)

var _ repository.EventRepository = &Implementation{}

type Implementation struct {
	tx txmanager.TxManager
}

func New(tx txmanager.TxManager) *Implementation {
	return &Implementation{tx: tx}
}
