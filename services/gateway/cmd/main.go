package main

import (
	"context"
	"log"

	"github.com/Molov30/go-microservices/services/gateway/internal/app"
	"github.com/Molov30/go-microservices/services/gateway/internal/config"
	"github.com/Molov30/go-microservices/services/gateway/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logger.New(cfg)
	app := app.New(&logger, cfg)

	if err := app.Run(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("failed to run app")
	}
}
