package runners

import (
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	eventrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/event"
	snapshotrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/snapshot"
	userrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/user"
	"go.uber.org/fx"
)

var ReposModule = fx.Module("repos",
	fx.Provide(
		fx.Annotate(
			func(tx txmanager.TxManager) repository.UserRepository {
				return userrepo.New(tx)
			},
		),
		fx.Annotate(
			func(tx txmanager.TxManager) repository.EventRepository {
				return eventrepo.New(tx)
			},
		),
		fx.Annotate(
			func(tx txmanager.TxManager) repository.SnapshotRepository {
				return snapshotrepo.New(tx)
			},
		),
	),
)
