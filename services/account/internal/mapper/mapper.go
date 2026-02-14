package mapper

import (
	accountpb "github.com/Molov30/go-microservices/generated/account"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Molov30/go-microservices/services/account/internal/model"
	repomodel "github.com/Molov30/go-microservices/services/account/internal/repository/model"
)

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

func PbsToUsers(pbs []*accountpb.User) []*model.User {
	users := make([]*model.User, 0, len(pbs))
	for _, pb := range pbs {
		users = append(users, PbToUser(pb))
	}
	return users
}

func PbToUserCreate(pb *accountpb.CreateUser) *model.CreateUser {
	return &model.CreateUser{
		Login:      pb.Login,
		Email:      pb.Email,
		FirstName:  pb.FirstName,
		LastName:   pb.LastName,
		MiddleName: pb.MiddleName,
		Password:   pb.Password,
		Age:        pb.Age,
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

func UserToRepoUser(user *model.User) *repomodel.User {
	return &repomodel.User{
		ID:         user.ID,
		Login:      user.Login,
		Phone:      user.Phone,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func RepoUserToUser(user *repomodel.User) *model.User {
	return &model.User{
		ID:         user.ID,
		Login:      user.Login,
		Phone:      user.Phone,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func RepoUsersToUsers(repoUsers []*repomodel.User) []*model.User {
	users := make([]*model.User, 0, len(repoUsers))
	for _, repoUser := range repoUsers {
		users = append(users, RepoUserToUser(repoUser))
	}
	return users
}

func UpdateUserToRepoUser(updateUser *model.UpdateUser) *repomodel.User {
	return &repomodel.User{
		Email:      updateUser.Email,
		Phone:      updateUser.Phone,
		FirstName:  updateUser.FirstName,
		LastName:   updateUser.LastName,
		MiddleName: updateUser.MiddleName,
		Age:        updateUser.Age,
	}
}
