package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" json:"service_name" required:"true" default:"account-service"`

	AppEnv string `env:"APP_ENV" json:"app_env" required:"true" default:"development"`

	Host     string `env:"GRPC_HOST" json:"host" required:"true" default:"localhost"`
	Port     string `env:"GRPC_PORT" json:"port" required:"true" default:"9000"`
	LogLevel string `env:"LOG_LEVEL" json:"log_level" required:"true" default:"info"`
	DbDSN    string `env:"DB_DSN" json:"db_dsn" required:"true"`

	JwtSecret             string `env:"JWT_SECRET" json:"jwt_secret" required:"true"`
	AccessTokenTTLMinutes int    `env:"ACCESS_TOKEN_TTL_MINUTES" json:"access_token_ttl_minutes" required:"true" defailt:"50052"`

	RefreshTokenTTLDays int `env:"REFRESH_TOKEN_TTL_DAYS" json:"refresh_token_ttl_days" required:"true" default:"30"`
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
