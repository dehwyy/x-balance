package runners

import (
	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/txmanagerfx/pkg/txmanager/gormtx"
	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NewDBOpts struct {
	fx.In
	Config *config.Config
}

func NewDB(opts NewDBOpts) (txmanager.TxManager, error) {
	db, err := gorm.Open(
		postgres.Open(opts.Config.DatabaseURL),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Snapshot{},
	); err != nil {
		return nil, err
	}

	tx, err := gormtx.New(gormtx.Opts{DB: db})
	if err != nil {
		return nil, err
	}

	return tx, nil
}
