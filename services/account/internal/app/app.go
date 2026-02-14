package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/Molov30/go-microservices/services/account/internal/config"
	"github.com/Molov30/go-microservices/services/account/internal/grpc"
	"github.com/Molov30/go-microservices/services/account/internal/repository"
	"github.com/Molov30/go-microservices/services/account/internal/service"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	accountRepository *repository.Repository
	accountService    *service.AccountService
	accountHandler    *grpc.Handler
	accountServer     *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	repo, err := a.getAccountRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	service := a.getAccountService(ctx, repo)
	handler := a.getAccountHandler(ctx, service)
	server := a.getAccountServer(ctx, handler)

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

func (a *App) getAccountRepository(_ context.Context) (*repository.Repository, error) {
	if a.accountRepository != nil {
		return a.accountRepository, nil
	}

	repository, err := repository.NewRepository(a.cfg, a.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	a.accountRepository = repository
	return repository, nil
}

func (a *App) getAccountService(_ context.Context, repo *repository.Repository) *service.AccountService {
	if a.accountService != nil {
		return a.accountService
	}

	service := service.NewAccountService(repo, a.logger)
	a.accountService = service
	return service
}

func (a *App) getAccountHandler(_ context.Context, service *service.AccountService) *grpc.Handler {
	if a.accountHandler != nil {
		return a.accountHandler
	}

	handler := grpc.NewHandler(a.logger, service)
	a.accountHandler = handler
	return handler
}

func (a *App) getAccountServer(_ context.Context, handler *grpc.Handler) *grpc.Server {
	if a.accountServer != nil {
		return a.accountServer
	}

	server := grpc.NewServer(a.logger, handler)
	a.accountServer = server
	return server
}
