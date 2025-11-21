package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

const (
	address  string = "localhost:8082"
	logLevel string = "info"
	kafka    string = "0.0.0.0:9092"
)

type Config struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
	Kafka    string `env:"KAFKA_BROKERS"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{
		Address:  address,
		LogLevel: logLevel,
		Kafka: kafka,
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
