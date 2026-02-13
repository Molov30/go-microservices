package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" required:"true" default:"account-service"`
	AppEnv      string `env:"APP_ENV" required:"true" default:"development"`
	Host        string `env:"HTTP_HOST" required:"true" default:"localhost"`
	Port        string `env:"HTTP_PORT" required:"true" default:"9000"`
	LogLevel    string `env:"LOG_LEVEL" required:"true" default:"info"`
	DbDSN       string `env:"DB_DSN" required:"true"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env config: %w", err)
	}

	return cfg, nil
}
