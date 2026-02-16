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

	Host            string `env:"GRPC_HOST" json:"host" required:"true" default:"localhost"`
	Port            string `env:"GRPC_PORT" json:"port" required:"true" default:"50053"`
	LogLevel        string `env:"LOG_LEVEL" json:"log_level" required:"true" default:"info"`
	AccountGRPCHost string `env:"ACCOUNT_GRPC_HOST" json:"account_grpc_host" required:"true"`
	AuthGRPCHost    string `env:"AUTH_GRPC_HOST" json:"auth_grpc_host" required:"true"`
	JwtSecret       string `env:"JWT_SECRET" json:"jwt_secret" required:"true"`
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
