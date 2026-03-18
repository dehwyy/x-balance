package runners

import (
	"github.com/dehwyy/x-balance/internal/config"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
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

	return db, nil
}
