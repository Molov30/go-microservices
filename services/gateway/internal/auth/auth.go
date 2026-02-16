package auth

import (
	"context"
	"fmt"

	authpb "github.com/Molov30/go-microservices/generated/auth"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

type Service struct {
	client authpb.AuthClient
}

func NewAuthService(client authpb.AuthClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Login(ctx context.Context, loginOrEmail, password string) (*model.TokenPair, error) {
	res, err := s.client.Login(ctx, &authpb.LoginRequest{
		LoginOrEmail: loginOrEmail,
		Password:     password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return PbToTokenPair(res.TokenPair), nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	res, err := s.client.Refresh(ctx, &authpb.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return PbToTokenPair(res.TokenPair), nil
}

func (s *Service) Verify(ctx context.Context, accessToken string) (uint64, error) {
	res, err := s.client.Validate(ctx, &authpb.ValidateRequest{AccessToken: accessToken})
	if err != nil {
		return 0, fmt.Errorf("failed to validate token: %w", err)
	}
	return res.UserId, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.client.Logout(ctx, &authpb.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (s *Service) Register(ctx context.Context, registerParams *model.RegisterUser) error {
	_, err := s.client.Register(ctx, &authpb.RegisterRequest{
		UserId:   registerParams.UserID,
		Login:    registerParams.Login,
		Email:    registerParams.Email,
		Password: registerParams.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uint64) error {
	_, err := s.client.DeleteUser(ctx, &authpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func PbToTokenPair(pb *authpb.TokenPair) *model.TokenPair {
	return &model.TokenPair{
		AccessToken:  pb.AccessToken,
		RefreshToken: pb.RefreshToken,
	}
}
