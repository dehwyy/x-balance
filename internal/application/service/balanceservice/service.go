package balanceservice

import (
	"context"
	"errors"

	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/domain/gateway"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	snapshotrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/snapshot"
	"go.uber.org/fx"
)

var ErrMaxRetriesExceeded = errors.New("max retries exceeded due to version conflict")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrFreezeNotFound = errors.New("freeze not found")

const maxRetries = 3

type Opts struct {
	fx.In

	TX              txmanager.TxManager
	EventRepo       repository.EventRepository
	SnapshotRepo    repository.SnapshotRepository
	UserRepo        repository.UserRepository
	BalanceCache    gateway.BalanceCache
	FreezeScheduler gateway.FreezeScheduler
	Config          BalanceConfig
}

type BalanceConfig struct {
	SnapshotEveryN int
}

type Service struct {
	tx              txmanager.TxManager
	eventRepo       repository.EventRepository
	snapshotRepo    repository.SnapshotRepository
	userRepo        repository.UserRepository
	balanceCache    gateway.BalanceCache
	freezeScheduler gateway.FreezeScheduler
	config          BalanceConfig
}

func New(opts Opts) *Service {
	return &Service{
		tx:              opts.TX,
		eventRepo:       opts.EventRepo,
		snapshotRepo:    opts.SnapshotRepo,
		userRepo:        opts.UserRepo,
		balanceCache:    opts.BalanceCache,
		freezeScheduler: opts.FreezeScheduler,
		config:          opts.Config,
	}
}

func (s *Service) withRetry(ctx context.Context, fn func(context.Context) error) error {
	for i := 0; i < maxRetries; i++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}
		if !errors.Is(err, snapshotrepo.ErrVersionConflict) {
			return err
		}
	}
	return ErrMaxRetriesExceeded
}
