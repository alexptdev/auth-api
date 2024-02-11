package env

import (
	"errors"
	"github.com/alexptdev/auth-api/internal/config"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPgConfig() (config.PgConfig, error) {

	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{dsn: dsn}, nil
}

func (cfg pgConfig) Dsn() string {
	return cfg.dsn
}
