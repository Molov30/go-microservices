package mapper

import (
	authpb "github.com/Molov30/go-microservices/generated/auth"

	"github.com/Molov30/go-microservices/services/auth/internal/model"
	repomodel "github.com/Molov30/go-microservices/services/auth/internal/repository/model"
)

func UserToRepoUser(user *model.User) *repomodel.User {
	return &repomodel.User{
		ID:           user.ID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func RepoUserToUser(user *repomodel.User) *model.User {
	return &model.User{
		ID:           user.ID,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func RepoUsersToUsers(repoUsers []*repomodel.User) []*model.User {
	users := make([]*model.User, 0, len(repoUsers))
	for _, repoUser := range repoUsers {
		users = append(users, RepoUserToUser(repoUser))
	}
	return users
}

func RefreshTokenToRepoRefresh(token *model.RefreshToken) *repomodel.RefreshToken {
	return &repomodel.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
	}
}

func RepoRefreshTokenToRefresh(token *repomodel.RefreshToken) *model.RefreshToken {
	return &model.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: token.RevokedAt,
		CreatedAt: token.CreatedAt,
	}
}

func PbLoginToUserLogin(pb *authpb.LoginRequest) *model.UserLogin {
	return &model.UserLogin{
		LoginOrEmail: pb.LoginOrEmail,
		Password:     pb.Password,
	}
}

func PbRegisterToUserRegister(pb *authpb.RegisterRequest) *model.UserRegister {
	return &model.UserRegister{
		UserID:   pb.UserId,
		Login:    pb.Login,
		Email:    pb.Email,
		Password: pb.Password,
	}
}
