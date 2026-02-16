package grpc

import (
	"context"

	accountpb "github.com/Molov30/go-microservices/generated/account"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Molov30/go-microservices/services/account/internal/mapper"
)

func (h *Handler) CreateUser(ctx context.Context, req *accountpb.CreateUserRequest) (*accountpb.CreateUserResponse, error) {
	userCreate := mapper.PbToUserCreate(req.User)
	user, err := h.accountService.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &accountpb.CreateUserResponse{
		User: mapper.UserToPb(user),
	}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *accountpb.GetUserRequest) (*accountpb.GetUserResponse, error) {
	user, err := h.accountService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, h.handleError(err)
	}
	return &accountpb.GetUserResponse{
		User: mapper.UserToPb(user),
	}, nil
}

func (h *Handler) GetUsers(ctx context.Context, req *accountpb.GetUsersRequest) (*accountpb.GetUsersResponse, error) {
	limit := int(req.GetPagination().GetLimit())
	offset := int(req.GetPagination().GetOffset())

	users, err := h.accountService.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &accountpb.GetUsersResponse{
		Users:      mapper.UsersToPb(users),
		Pagination: req.GetPagination(),
	}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *accountpb.UpdateUserRequest) (*emptypb.Empty, error) {
	user := mapper.PbToUserUpdate(req.User)
	err := h.accountService.UpdateUser(ctx, req.GetUserId(), user)
	if err != nil {
		return nil, h.handleError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *accountpb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := h.accountService.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, h.handleError(err)
	}

	return &emptypb.Empty{}, nil
}
