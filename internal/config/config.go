package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL    string `env:"DATABASE_URL,required"`
	RedisURL       string `env:"REDIS_URL,required"`
	Port           int    `env:"PORT,required"`
	SnapshotEveryN int    `env:"SNAPSHOT_EVERY_N"`
	SnapshotCron   string `env:"SNAPSHOT_CRON"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
