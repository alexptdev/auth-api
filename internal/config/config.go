package config

import (
	"github.com/joho/godotenv"
)

type GrpcConfig interface {
	Address() string
}

type PgConfig interface {
	Dsn() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
