package runners

import (
	"github.com/dehwyy/dbfx/pkg/gormfx"
	"github.com/dehwyy/dbfx/pkg/gormfx/postgres"
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/txmanagerfx/pkg/txmanager/gormtx"
	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type NewGORMOpts struct {
	fx.In
	Config *config.Config
}

func NewGORM(opts NewGORMOpts) (*gorm.DB, error) {
	return gormfx.New(
		gormfx.Opts{
			Postgres: &postgres.Opts{
				ConnectionStrings: []string{opts.Config.DatabaseURL},
			},
		},
	)()
}

func NewTxManager(db *gorm.DB) (txmanager.TxManager, error) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Snapshot{},
	); err != nil {
		return nil, err
	}

	return gormtx.New(
		gormtx.Opts{
			DB: db,
		},
	)
}
