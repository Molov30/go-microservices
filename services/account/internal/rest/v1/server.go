package v1

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/Molov30/go-microservices/services/account/internal/config"
	"github.com/Molov30/go-microservices/services/account/internal/rest/v1/generated"
	"github.com/Molov30/go-microservices/services/account/internal/rest/v1/handlers"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config, logger *zerolog.Logger) *Server {
	handler := handlers.NewHandler(logger)
	mux := http.NewServeMux()
	generated.HandlerFromMuxWithBaseURL(handler, mux, "/api/v1")

	srv := &http.Server{
		ReadTimeout: 5 * time.Second,
		Addr:        net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:     mux,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to run server: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
