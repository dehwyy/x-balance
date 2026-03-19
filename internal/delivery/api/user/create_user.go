package userhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userDelivery.CreateUser", dspan.Attr("req", req))
	defer span.End()

	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	response, err := h.userservice.CreateUser(ctx, userservice.CreateUserRequest{
		Name:           req.Name,
		OverdraftLimit: limit,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &userspb.CreateUserResponse{User: userconvert.UserToProto(response.User)}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
