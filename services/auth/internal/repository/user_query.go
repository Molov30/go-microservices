package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Molov30/go-microservices/services/auth/internal/mapper"
	"github.com/Molov30/go-microservices/services/auth/internal/model"
	repomodel "github.com/Molov30/go-microservices/services/auth/internal/repository/model"
)

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	userRepo := mapper.UserToRepoUser(user)

	res := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(userRepo)

	if res.Error != nil {
		return fmt.Errorf("failed to create user: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
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

func (r *Repository) GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (*model.User, error) {
	var user repomodel.User

	res := r.db.
		WithContext(ctx).
		Model(&repomodel.User{}).
		Where("login = ? OR email = ?", loginOrEmail, loginOrEmail).
		First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("failed to get user by email or login: %w", res.Error)
	}

	return mapper.RepoUserToUser(&user), nil
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
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}

	return nil
}
