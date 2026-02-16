package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Molov30/go-microservices/services/account/internal/mapper"
	"github.com/Molov30/go-microservices/services/account/internal/model"
	repomodel "github.com/Molov30/go-microservices/services/account/internal/repository/model"
)

func (r *Repository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	userRepo := mapper.UserToRepoUser(user)

	res := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(userRepo)

	if res.Error != nil {
		return nil, fmt.Errorf("failed to save user: %w", res.Error)
	}

	user.ID = userRepo.ID
	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	var user repomodel.User
	res := r.db.
		WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", userID).
		First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", res.Error)
	}
	return mapper.RepoUserToUser(&user), nil
}

func (r *Repository) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	var users []*repomodel.User
	res := r.db.
		WithContext(ctx).
		Model(&repomodel.User{}).
		Offset(offset).
		Limit(limit).
		Find(&users)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("failed to get users: %w", res.Error)
	}
	return mapper.RepoUsersToUsers(users), nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint64) error {
	res := r.db.
		WithContext(ctx).
		Where("id = ?", userID).
		Delete(&repomodel.User{})

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrEntryNotFound
		}
		return fmt.Errorf("failed to delete user by id: %w", res.Error)
	}
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID uint64, user *model.UpdateUser) error {
	updateUser := mapper.UpdateUserToRepoUser(user)

	res := r.db.
		WithContext(ctx).
		Where("id = ?", userID).
		Updates(updateUser)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrEntryNotFound
		}
		return fmt.Errorf("failed to update user by id: %w", res.Error)
	}
	return nil
}
