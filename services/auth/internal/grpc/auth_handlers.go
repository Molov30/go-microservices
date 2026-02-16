package grpc

import (
	"context"

	authpb "github.com/Molov30/go-microservices/generated/auth"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Molov30/go-microservices/services/auth/internal/mapper"
)

func (h *Handler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	loginParams := mapper.PbLoginToUserLogin(req)
	tokenPair, err := h.accountService.Login(ctx, loginParams)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &authpb.LoginResponse{
		TokenPair: &authpb.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		},
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *authpb.RefreshRequest) (*emptypb.Empty, error) {
	err := h.accountService.Logout(ctx, req.RefreshToken)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	tokenPair, err := h.accountService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &authpb.RefreshResponse{
		TokenPair: &authpb.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		},
	}, nil
}

func (h *Handler) Register(ctx context.Context, req *authpb.RegisterRequest) (*emptypb.Empty, error) {
	registerParams := mapper.PbRegisterToUserRegister(req)
	err := h.accountService.Register(ctx, registerParams)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	userID, err := h.accountService.Validate(ctx, req.AccessToken)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &authpb.ValidateResponse{UserId: userID}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *authpb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := h.accountService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}
