package grpc

import (
	"context"

	gatewaypb "github.com/Molov30/go-microservices/generated/gateway"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Molov30/go-microservices/services/gateway/internal/grpc/interceptor"
	"github.com/Molov30/go-microservices/services/gateway/internal/mapper"
)

func (h *Handler) Register(ctx context.Context, req *gatewaypb.RegisterRequest) (*gatewaypb.RegisterResponse, error) {
	user := mapper.PbToUserCreate(req.User)
	createdUser, tokens, err := h.gatewayService.Register(ctx, user, req.Password)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.RegisterResponse{
		User:   mapper.UserToPb(createdUser),
		Tokens: mapper.TokenPairToPb(tokens),
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *gatewaypb.LoginRequest) (*gatewaypb.LoginResponse, error) {
	user, tokens, err := h.gatewayService.Login(ctx, req.LoginOrEmail, req.Password)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.LoginResponse{
		User:   mapper.UserToPb(user),
		Tokens: mapper.TokenPairToPb(tokens),
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *gatewaypb.LogoutRequest) (*empty.Empty, error) {
	err := h.gatewayService.Logout(ctx, req.RefreshToken)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) Refresh(ctx context.Context, req *gatewaypb.RefreshRequest) (*gatewaypb.RefreshResponse, error) {
	tokens, err := h.gatewayService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.RefreshResponse{
		TokenPair: mapper.TokenPairToPb(tokens),
	}, nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *gatewaypb.ValidateTokenRequest) (*gatewaypb.ValidateTokenResponse, error) {
	userID, isValid, err := h.gatewayService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &gatewaypb.ValidateTokenResponse{
		UserId:  userID,
		IsValid: isValid,
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *gatewaypb.CreateUserRequest) (*empty.Empty, error) {
	user := mapper.PbToUserCreate(req.User)
	err := h.gatewayService.CreateUser(ctx, user)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *gatewaypb.GetUserRequest) (*gatewaypb.GetUserResponse, error) {
	user, err := h.gatewayService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.GetUserResponse{User: mapper.UserToPb(user)}, nil
}

func (h *Handler) GetCurrentUser(ctx context.Context, _ *empty.Empty) (*gatewaypb.GetCurrentUserResponse, error) {
	userID, err := interceptor.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, h.handleError(err)
	}

	user, err := h.gatewayService.GetUser(ctx, userID)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.GetCurrentUserResponse{User: mapper.UserToPb(user)}, nil
}

func (h *Handler) GetUsers(ctx context.Context, req *gatewaypb.GetUsersRequest) (*gatewaypb.GetUsersResponse, error) {
	limit, offset := getPagination(req.GetPagination())

	users, err := h.gatewayService.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &gatewaypb.GetUsersResponse{Users: mapper.UsersToPb(users), Pagination: req.Pagination}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *gatewaypb.UpdateUserRequest) (*empty.Empty, error) {
	user := mapper.PbToUserUpdate(req.User)
	err := h.gatewayService.UpdateUser(ctx, req.UserId, user)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) UpdateCurrentUser(ctx context.Context, req *gatewaypb.UpdateCurrentUserRequest) (*empty.Empty, error) {
	userID, err := interceptor.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, h.handleError(err)
	}

	user := mapper.PbToUserUpdate(req.User)
	err = h.gatewayService.UpdateUser(ctx, userID, user)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *gatewaypb.DeleteUserRequest) (*empty.Empty, error) {
	err := h.gatewayService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteCurrentUser(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	userID, err := interceptor.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, h.handleError(err)
	}

	err = h.gatewayService.DeleteUser(ctx, userID)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}
