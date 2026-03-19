package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	userconvert "github.com/dehwyy/x-balance/internal/delivery/api/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func (h *Handler) GetUser(
	ctx context.Context,
	req *userspb.GetUserRequest,
) (*userspb.GetUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userDelivery.GetUser", dspan.Attr("req", req))
	defer span.End()

	response, err := h.userservice.GetUser(ctx, userconvert.GetUserRequestToDomain(req))
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := userconvert.GetUserResponseToProto(response)
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
