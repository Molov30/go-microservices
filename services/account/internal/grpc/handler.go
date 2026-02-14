package grpc

import (
	"context"
	"errors"

	accountpb "github.com/Molov30/go-microservices/generated/account"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Molov30/go-microservices/services/account/internal/model"
	"github.com/Molov30/go-microservices/services/account/internal/repository"
)

type AccountService interface {
	CreateUser(ctx context.Context, user *model.CreateUser) error
	GetUser(ctx context.Context, userID uint64) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error
}

type Handler struct {
	accountpb.UnimplementedAccountServer

	accountService AccountService
	logger         *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger, accountService AccountService) *Handler {
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
