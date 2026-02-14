package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Molov30/go-microservices/services/auth/internal/mapper"
	"github.com/Molov30/go-microservices/services/auth/internal/model"
	repomodel "github.com/Molov30/go-microservices/services/auth/internal/repository/model"
	"gorm.io/gorm"
)

func (r *Repository) SaveRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	repoToken := mapper.RefreshTokenToRepoRefresh(token)

	res := r.db.
		WithContext(ctx).
		Create(repoToken)

	if res.Error != nil {
		return fmt.Errorf("failed to save refresh token in repo: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken repomodel.RefreshToken

	res := r.db.
		WithContext(ctx).
		Model(&repomodel.RefreshToken{}).
		Where("token = ?", token).
		First(&refreshToken)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", res.Error)
	}

	return mapper.RepoRefreshTokenToRefresh(&refreshToken), nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	res := r.db.
		WithContext(ctx).
		Model(&repomodel.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", gorm.Expr("NOW()"))

	if res.Error != nil {
		return fmt.Errorf("failed to revoked token: %w", res.Error)
	}
	return nil
}
