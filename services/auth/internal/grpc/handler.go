package grpc

import (
	"context"
	"errors"

	authpb "github.com/Molov30/go-microservices/generated/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Molov30/go-microservices/services/auth/internal/model"
	"github.com/Molov30/go-microservices/services/auth/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, registerParams *model.UserRegister) error
	Login(ctx context.Context, loginParam *model.UserLogin) (*model.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Validate(ctx context.Context, accessToken string) (uint64, error)
	Logout(ctx context.Context, refreshToken string) error
	DeleteUser(ctx context.Context, userID uint64) error
}

type Handler struct {
	authpb.UnimplementedAuthServer

	accountService AuthService
	logger         *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger, accountService AuthService) *Handler {
	return &Handler{
		logger:         logger,
		accountService: accountService,
	}
}

func (h *Handler) handleError(err error) error {
	if errors.Is(err, repository.ErrEntryNotFound) {
		return status.Error(codes.NotFound, "entry not found")
	}

	h.logger.Error().Err(err).Msg("internal error")
	return status.Error(codes.Internal, "internal server error")
}
