package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/Molov30/go-microservices/services/account/internal/config"
	"github.com/Molov30/go-microservices/services/account/internal/logger"
	v1 "github.com/Molov30/go-microservices/services/account/internal/rest/v1"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logger.New(cfg)
	srv := v1.NewServer(cfg, &logger)

	go func() {
		err := srv.Run()
		if err != nil {
			logger.Err(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()

	stopCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	err = srv.Stop(stopCtx)
	if err != nil {
		logger.Err(err).Msg("failed to stop server")
	}
}
