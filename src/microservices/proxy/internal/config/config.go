package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

const (
	address                string = "localhost:8000"
	logLevel               string = "info"
	monolithUrl            string = "http://monolith:8080"
	moviesServiceUrl       string = "http://movies-service:8081"
	eventsServiceUrl       string = "http://events-service:8082"
	gradualMigration       bool   = true
	moviesMigrationPercent int    = 10
)

type Config struct {
	Address                string `env:"SERVER_ADDRESS"`
	LogLevel               string `env:"LOG_LEVEL"`
	MonolithUrl            string `env:"MONOLITH_URL"`
	MoviesServiceUrl       string `env:"MOVIES_SERVICE_URL"`
	EventsServiceUrl       string `env:"EVENTS_SERVICE_URL"`
	GradualMigration       bool   `env:"GRADUAL_MIGRATION"`
	MoviesMigrationPercent int    `env:"MOVIES_MIGRATION_PERCENT"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{
		Address:                address,
		LogLevel:               logLevel,
		MonolithUrl:            monolithUrl,
		MoviesServiceUrl:       moviesServiceUrl,
		EventsServiceUrl:       eventsServiceUrl,
		GradualMigration:       gradualMigration,
		MoviesMigrationPercent: moviesMigrationPercent,
	}

	if err := cnf.initEnvs(); err != nil {
		return Config{}, err
	}

	return cnf, nil
}

func (cnf *Config) initEnvs() error {
	if err := env.Parse(cnf); err != nil {
		return fmt.Errorf("InitEnvs: parse envs fail: %w", err)
	}

	return nil
}
