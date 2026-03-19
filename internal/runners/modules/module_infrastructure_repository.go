package modules

import (
	"github.com/dehwyy/x-balance/internal/domain/repository"
	eventrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/event"
	snapshotrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/snapshot"
	userrepo "github.com/dehwyy/x-balance/internal/infrastructure/repository/user"
	"go.uber.org/fx"
)

var InfrastructureRepositoryModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			userrepo.New,
			fx.As(new(repository.UserRepository)),
		),
		fx.Annotate(
			eventrepo.New,
			fx.As(new(repository.EventRepository)),
		),
		fx.Annotate(
			snapshotrepo.New,
			fx.As(new(repository.SnapshotRepository)),
		),
	),
)
