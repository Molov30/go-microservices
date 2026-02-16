package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	accountpb "github.com/Molov30/go-microservices/generated/account"
	authpb "github.com/Molov30/go-microservices/generated/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Molov30/go-microservices/services/gateway/internal/account"
	"github.com/Molov30/go-microservices/services/gateway/internal/auth"
	"github.com/Molov30/go-microservices/services/gateway/internal/config"
	controller "github.com/Molov30/go-microservices/services/gateway/internal/grpc"
	"github.com/Molov30/go-microservices/services/gateway/internal/service"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	gatewayService *service.GatewayService
	authService    *auth.Service
	accountService *account.Service

	gatewayHandler *controller.Handler
	gatewayServer  *controller.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	accountService, err := a.getAccountService(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account service: %w", err)
	}

	authService, err := a.getAuthService(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth service: %w", err)
	}

	service := a.getGatewayService(ctx, authService, accountService)
	handler := a.getgatewayHandler(ctx, service)
	server := a.getGatewayServer(ctx, handler)

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

	return err
}

func (a *App) getAccountService(_ context.Context) (*account.Service, error) {
	if a.accountService != nil {
		return a.accountService, nil
	}

	conn, err := grpc.NewClient(
		a.cfg.AccountGRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	client := accountpb.NewAccountClient(conn)
	a.accountService = account.NewAccountService(client)
	return a.accountService, nil
}

func (a *App) getAuthService(_ context.Context) (*auth.Service, error) {
	if a.authService != nil {
		return a.authService, nil
	}

	conn, err := grpc.NewClient(
		a.cfg.AuthGRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	client := authpb.NewAuthClient(conn)
	a.authService = auth.NewAuthService(client)
	return a.authService, nil
}

func (a *App) getGatewayService(_ context.Context,
	authService service.AuthService,
	accountService service.AccountService,
) *service.GatewayService {
	if a.gatewayService != nil {
		return a.gatewayService
	}

	service := service.NewGatewayService(a.logger, accountService, authService)
	a.gatewayService = service
	return service
}

func (a *App) getgatewayHandler(_ context.Context, service *service.GatewayService) *controller.Handler {
	if a.gatewayHandler != nil {
		return a.gatewayHandler
	}

	handler := controller.NewHandler(a.logger, service)
	a.gatewayHandler = handler
	return handler
}

func (a *App) getGatewayServer(_ context.Context, handler *controller.Handler) *controller.Server {
	if a.gatewayServer != nil {
		return a.gatewayServer
	}

	server := controller.NewServer(a.logger, handler)
	a.gatewayServer = server
	return server
}
