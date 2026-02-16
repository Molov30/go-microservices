package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

type AccountService interface {
	GetUser(ctx context.Context, userID uint64) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset uint32) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	CreateUser(ctx context.Context, user *model.CreateUser) (*model.User, error)
	UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error
}

type AuthService interface {
	Login(ctx context.Context, loginOrEmail, password string) (*model.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Verify(ctx context.Context, accessToken string) (uint64, error)
	Logout(ctx context.Context, refreshToken string) error
	Register(ctx context.Context, registerParams *model.RegisterUser) error
	DeleteUser(ctx context.Context, userID uint64) error
}

type GatewayService struct {
	logger *zerolog.Logger

	accountService AccountService
	authService    AuthService
}

func NewGatewayService(logger *zerolog.Logger, accountSerivce AccountService, authService AuthService) *GatewayService {
	return &GatewayService{
		logger:         logger,
		accountService: accountSerivce,
		authService:    authService,
	}
}

func (s *GatewayService) Register(ctx context.Context, newUser *model.CreateUser, passowrd string) (*model.User, *model.TokenPair, error) {
	user, err := s.accountService.CreateUser(ctx, newUser)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	registerParams := &model.RegisterUser{
		UserID:   user.ID,
		Login:    user.Login,
		Email:    user.Email,
		Password: passowrd,
	}
	err = s.authService.Register(ctx, registerParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register user: %w", err)
	}

	token, err := s.authService.Login(ctx, user.Login, passowrd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to login user: %w", err)
	}
	return user, token, nil
}

func (s *GatewayService) Login(ctx context.Context, loginOrEmail, password string) (*model.User, *model.TokenPair, error) {
	tokens, err := s.authService.Login(ctx, loginOrEmail, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to login user: %w", err)
	}

	return &model.User{}, tokens, nil
}

func (s *GatewayService) Logout(ctx context.Context, token string) error {
	err := s.authService.Logout(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (s *GatewayService) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	return s.authService.Refresh(ctx, refreshToken)
}

func (s *GatewayService) ValidateToken(ctx context.Context, accessToken string) (uint64, bool, error) {
	userID, err := s.authService.Verify(ctx, accessToken)
	if err != nil {
		return 0, false, fmt.Errorf("failed to verify access token: %w", err)
	}
	return userID, true, nil
}

func (s *GatewayService) CreateUser(ctx context.Context, newUser *model.CreateUser) error {
	_, err := s.accountService.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}
	return nil
}

func (s *GatewayService) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	return s.accountService.GetUser(ctx, userID)
}

func (s *GatewayService) GetUsers(ctx context.Context, limit, offset uint32) ([]*model.User, error) {
	return s.accountService.GetUsers(ctx, limit, offset)
}

func (s GatewayService) DeleteUser(ctx context.Context, userID uint64) error {
	err := s.authService.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return s.accountService.DeleteUser(ctx, userID)
}

func (s GatewayService) UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error {
	return s.accountService.UpdateUser(ctx, userID, user)
}
