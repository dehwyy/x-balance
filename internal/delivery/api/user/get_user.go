package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func (h *Handler) GetUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.GetUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userDelivery.GetUser", dspan.Attr("req", req))
	defer span.End()

	response, err := h.userservice.GetUser(ctx, userservice.GetUserRequest{
		ID: req.Id,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &userspb.GetUserResponse{User: userconvert.UserToProto(response.User)}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
