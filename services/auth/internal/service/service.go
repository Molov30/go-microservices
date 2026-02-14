package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Molov30/go-microservices/services/auth/internal/config"
	"github.com/Molov30/go-microservices/services/auth/internal/model"
	"github.com/Molov30/go-microservices/services/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID uint64) (*model.User, error)
	GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (*model.User, error)
	SaveRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

type AuthService struct {
	repo   Repository
	cfg    *config.Config
	logger *zerolog.Logger
}

func NewAuthService(repo Repository, cfg *config.Config, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *AuthService) Register(ctx context.Context, registerParams *model.UserRegister) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(registerParams.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate from password: %w", err)
	}

	user := &model.User{
		Login:        registerParams.Login,
		Email:        registerParams.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, loginParam *model.UserLogin) (*model.TokenPair, error) {
	user, err := s.repo.GetUserByLoginOrEmail(ctx, loginParam.LoginOrEmail)
	if err != nil {
		if errors.Is(err, repository.ErrEntryNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user by login or email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginParam.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrEntryNotFound) {
			return nil, errors.New("invalid token")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return s.issueTokens(ctx, rt.UserID)
}

func (s *AuthService) Validate(_ context.Context, accessToken string) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(_ *jwt.Token) (any, error) {
		return []byte(s.cfg.JwtSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	uidFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("token invalid subject")
	}

	return uint64(uidFloat), nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint64) (*model.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(s.cfg.AccessTokenTTLMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": accessExp,
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := access.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed sign access token: %w", err)
	}

	refreshRaw := fmt.Sprintf("%d:%d:%s", userID, now.UnixNano(), s.cfg.JwtSecret)
	h := sha256.Sum256([]byte(refreshRaw))
	refreshStr := hex.EncodeToString(h[:])
	refresh := &model.RefreshToken{
		UserID:    userID,
		Token:     refreshStr,
		ExpiresAt: now.Add(time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour),
		CreatedAt: now,
	}

	err = s.repo.SaveRefreshToken(ctx, refresh)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &model.TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}
