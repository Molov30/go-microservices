package account

import (
	"context"
	"fmt"

	accountpb "github.com/Molov30/go-microservices/generated/account"
	"github.com/Molov30/go-microservices/generated/pagination"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

type Service struct {
	client accountpb.AccountClient
}

func NewAccountService(client accountpb.AccountClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	user, err := s.client.GetUser(ctx, &accountpb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return PbToUser(user.User), nil
}

func (s *Service) GetUsers(ctx context.Context, limit, offset uint32) ([]*model.User, error) {
	pagination := &pagination.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	users, err := s.client.GetUsers(ctx, &accountpb.GetUsersRequest{Pagination: pagination})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return PbsToUsers(users.Users), nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uint64) error {
	_, err := s.client.DeleteUser(ctx, &accountpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *Service) CreateUser(ctx context.Context, user *model.CreateUser) (*model.User, error) {
	res, err := s.client.CreateUser(ctx, &accountpb.CreateUserRequest{
		User: UserCreateToPb(user),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return PbToUser(res.User), nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error {
	_, err := s.client.UpdateUser(ctx, &accountpb.UpdateUserRequest{
		UserId: userID,
		User:   UserUpdateToPb(user),
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func PbToUser(userpb *accountpb.User) *model.User {
	return &model.User{
		ID:         userpb.Id,
		Login:      userpb.Login,
		Email:      userpb.Email,
		Phone:      userpb.Phone,
		FirstName:  userpb.FirstName,
		LastName:   userpb.LastName,
		MiddleName: userpb.MiddleName,
		Age:        userpb.Age,
		CreatedAt:  userpb.CreatedAt.AsTime(),
		UpdatedAt:  userpb.UpdatedAt.AsTime(),
	}
}

func PbsToUsers(pbs []*accountpb.User) []*model.User {
	users := make([]*model.User, 0, len(pbs))
	for _, pb := range pbs {
		users = append(users, PbToUser(pb))
	}
	return users
}

func UserCreateToPb(user *model.CreateUser) *accountpb.CreateUser {
	return &accountpb.CreateUser{
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func UserUpdateToPb(user *model.UpdateUser) *accountpb.UpdateUser {
	return &accountpb.UpdateUser{
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}
