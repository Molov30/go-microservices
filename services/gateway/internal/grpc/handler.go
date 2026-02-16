package grpc

import (
	"context"

	gatewaypb "github.com/Molov30/go-microservices/generated/gateway"
	"github.com/Molov30/go-microservices/generated/pagination"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

type GatewayService interface {
	Register(ctx context.Context, newUser *model.CreateUser, password string) (*model.User, *model.TokenPair, error)
	Login(ctx context.Context, loginOrEmail, password string) (*model.User, *model.TokenPair, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	ValidateToken(ctx context.Context, accessToken string) (uint64, bool, error)
	CreateUser(ctx context.Context, newUser *model.CreateUser) error
	GetUser(ctx context.Context, userID uint64) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset uint32) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error
}

type Handler struct {
	gatewaypb.UnimplementedGatewayServer

	gatewayService GatewayService
	logger         *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger, gatewayService GatewayService) *Handler {
	return &Handler{
		logger:         logger,
		gatewayService: gatewayService,
	}
}

func (h *Handler) handleError(err error) error {
	h.logger.Error().Err(err).Msg("internal error")
	return status.Error(codes.Internal, "internal server error")
}

func getPagination(pagination *pagination.Pagination) (limit, offset uint32) {
	limit = 10
	offset = 0

	if pagination != nil {
		if pagination.Limit > 0 {
			limit = pagination.Limit
		}
		offset = pagination.Offset
	}

	return limit, offset
}
