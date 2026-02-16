package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/Molov30/go-microservices/services/account/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, userID uint64) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error
}

type AccountService struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewAccountService(repo Repository, logger *zerolog.Logger) *AccountService {
	return &AccountService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AccountService) CreateUser(ctx context.Context, newUser *model.CreateUser) (*model.User, error) {
	user := &model.User{
		Login:      newUser.Login,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		MiddleName: newUser.MiddleName,
		Age:        newUser.Age,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	user, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *AccountService) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *AccountService) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}

func (s *AccountService) DeleteUser(ctx context.Context, userID uint64) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *AccountService) UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error {
	return s.repo.UpdateUser(ctx, userID, user)
}
