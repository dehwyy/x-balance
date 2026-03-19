package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type GetUserRequest struct {
	ID string
}

type GetUserResponse struct {
	User *user.User
}

func (s *Service) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userservice.Service.GetUser", dspan.Attr("req", req))
	defer span.End()

	userResp, err := s.userRepo.GetByID(ctx, dto.UserGetByIDRequest{ID: user.ID{Value: req.ID}})
	if err != nil {
		return nil, span.Err(err)
	}

	response := &GetUserResponse{User: &userResp.User}
	span.WithAttribute("response", response)
	return response, nil
}
