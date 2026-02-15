package grpc

import (
	"context"
	"errors"
	"net"

	authpb "github.com/Molov30/go-microservices/generated/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	srv    *grpc.Server
	logger *zerolog.Logger
}

func NewServer(logger *zerolog.Logger, handler authpb.AuthServer) *Server {
	srv := grpc.NewServer()
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthSrv)
	reflection.Register(srv)

	authpb.RegisterAuthServer(srv, handler)

	return &Server{
		srv:    srv,
		logger: logger,
	}
}

func (s *Server) Run(l net.Listener) error {
	err := s.srv.Serve(l)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.srv.Stop()
		return errors.New("grpc server forcing stop")
	case <-stopped:
		return nil
	}
}
