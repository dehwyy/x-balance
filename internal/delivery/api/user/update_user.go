package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	userconvert "github.com/dehwyy/x-balance/internal/delivery/api/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func (h *Handler) UpdateUser(
	ctx context.Context,
	req *userspb.UpdateUserRequest,
) (*userspb.UpdateUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userDelivery.UpdateUser", dspan.Attr("req", req))
	defer span.End()

	response, err := h.userservice.UpdateUser(ctx, userconvert.UpdateUserRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := userconvert.UpdateUserResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
