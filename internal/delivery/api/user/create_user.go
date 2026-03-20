package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	userconvert "github.com/dehwyy/x-balance/internal/delivery/api/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func (h *Handler) CreateUser(
	ctx context.Context,
	req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"userDelivery.CreateUser",
		dspan.Attr("req", req),
	)
	defer span.End()

	response, err := h.userservice.CreateUser(
		ctx,
		userconvert.CreateUserRequestToDomain(req),
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(
		span,
		userconvert.CreateUserResponseToProto(response),
	), nil
}
