package mapper

import (
	accountpb "github.com/Molov30/go-microservices/generated/account"
	authpb "github.com/Molov30/go-microservices/generated/auth"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

func CreateUserToUser(cUser *model.CreateUser) *model.User {
	return &model.User{
		Login:      cUser.Login,
		Email:      cUser.Email,
		Phone:      cUser.Phone,
		FirstName:  cUser.FirstName,
		LastName:   cUser.LastName,
		MiddleName: cUser.MiddleName,
		Age:        cUser.Age,
	}
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

func PbToUserUpdate(pb *accountpb.UpdateUser) *model.UpdateUser {
	return &model.UpdateUser{
		Email:      pb.Email,
		Phone:      pb.Phone,
		FirstName:  pb.FirstName,
		LastName:   pb.LastName,
		MiddleName: pb.MiddleName,
		Age:        pb.Age,
	}
}

func UserToPb(user *model.User) *accountpb.User {
	return &accountpb.User{
		Id:         user.ID,
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}
}

func UsersToPb(users []*model.User) []*accountpb.User {
	pbs := make([]*accountpb.User, 0, len(users))
	for _, user := range users {
		pbs = append(pbs, UserToPb(user))
	}
	return pbs
}

func TokenPairToPb(tokens *model.TokenPair) *authpb.TokenPair {
	return &authpb.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func PbToUserCreate(pb *accountpb.CreateUser) *model.CreateUser {
	return &model.CreateUser{
		Login:      pb.Login,
		Email:      pb.Email,
		FirstName:  pb.FirstName,
		LastName:   pb.LastName,
		MiddleName: pb.MiddleName,
		Age:        pb.Age,
	}
}
