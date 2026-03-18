package userhandler

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	"github.com/dehwyy/x-balance/internal/delivery/api/user/convert"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
)

type Handler struct {
	balancepb.UnimplementedUserServiceServer
	svc *userservice.Service
}

func New(svc *userservice.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateUser(
	ctx context.Context,
	req *balancepb.CreateUserRequest,
) (*balancepb.CreateUserResponse, error) {
	res, err := h.svc.CreateUser(ctx, convert.CreateUserRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.CreateUserResponseToProto(res), nil
}

func (h *Handler) GetUser(
	ctx context.Context,
	req *balancepb.GetUserRequest,
) (*balancepb.GetUserResponse, error) {
	res, err := h.svc.GetUser(ctx, convert.GetUserRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.GetUserResponseToProto(res), nil
}

func (h *Handler) UpdateUser(
	ctx context.Context,
	req *balancepb.UpdateUserRequest,
) (*balancepb.UpdateUserResponse, error) {
	res, err := h.svc.UpdateUser(ctx, convert.UpdateUserRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return convert.UpdateUserResponseToProto(res), nil
}

func (h *Handler) DeleteUser(
	ctx context.Context,
	req *balancepb.DeleteUserRequest,
) (*balancepb.DeleteUserResponse, error) {
	if err := h.svc.DeleteUser(ctx, convert.DeleteUserRequestToDomain(req)); err != nil {
		return nil, err
	}
	return &balancepb.DeleteUserResponse{}, nil
}
