package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func (h *Handler) DeleteUser(ctx context.Context, req *userspb.DeleteUserRequest) (*userspb.DeleteUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userDelivery.DeleteUser", dspan.Attr("req", req))
	defer span.End()

	err := h.userservice.DeleteUser(ctx, userservice.DeleteUserRequest{
		ID: req.Id,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	return &userspb.DeleteUserResponse{}, nil
}
