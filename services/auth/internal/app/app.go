package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/Molov30/go-microservices/services/auth/internal/config"
	"github.com/Molov30/go-microservices/services/auth/internal/grpc"
	"github.com/Molov30/go-microservices/services/auth/internal/repository"
	"github.com/Molov30/go-microservices/services/auth/internal/service"
	"github.com/rs/zerolog"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	authRepository *repository.Repository
	authService    *service.AuthService
	authHandler    *grpc.Handler
	authServer     *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	repo, err := a.getAuthRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	service := a.getAuthService(ctx, repo)
	handler := a.getAuthHandler(ctx, service)
	server := a.getAuthServer(ctx, handler)

	lCfg := &net.ListenConfig{}
	l, err := lCfg.Listen(ctx, "tcp", net.JoinHostPort(a.cfg.Host, a.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen grpc addr: %w", err)
	}

	stopCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)

		a.logger.Info().Msg("run grpc server")
		err := server.Run(l)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-stopCtx.Done():
		a.logger.Info().Msg("stopping with context")
	case err = <-errChan:
	}

	// Shutdown phase uses independent context with fixed timeout, not derived from startup ctx.
	// This is intentional: shutdown must complete even if startup ctx is cancelled.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//nolint:contextcheck // Shutdown context intentionally independent from startup context
	err = errors.Join(err, server.Stop(shutdownCtx))
	//nolint:contextcheck // Shutdown context intentionally independent from startup context
	err = errors.Join(err, repo.Close(shutdownCtx))
	return err
}

func (a *App) getAuthRepository(_ context.Context) (*repository.Repository, error) {
	if a.authRepository != nil {
		return a.authRepository, nil
	}

	repository, err := repository.NewRepository(a.cfg, a.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	a.authRepository = repository
	return repository, nil
}

func (a *App) getAuthService(_ context.Context, repo *repository.Repository) *service.AuthService {
	if a.authService != nil {
		return a.authService
	}

	service := service.NewAuthService(repo, a.cfg, a.logger)
	a.authService = service
	return service
}

func (a *App) getAuthHandler(_ context.Context, service *service.AuthService) *grpc.Handler {
	if a.authHandler != nil {
		return a.authHandler
	}

	handler := grpc.NewHandler(a.logger, service)
	a.authHandler = handler
	return handler
}

func (a *App) getAuthServer(_ context.Context, handler *grpc.Handler) *grpc.Server {
	if a.authServer != nil {
		return a.authServer
	}

	server := grpc.NewServer(a.logger, handler)
	a.authServer = server
	return server
}
