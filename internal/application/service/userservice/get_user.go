package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"

	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type GetUserRequest struct {
	ID user.ID
}

type GetUserResponse struct {
	User *user.User
}

func (s *Service) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"userservice.Service.GetUser",
		dspan.Attr("req", req),
	)
	defer span.End()

	userDTO, err := s.userRepo.GetByID(
		ctx,
		dto.UserGetByIDRequest{
			ID: req.ID,
		},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(
		span,
		&GetUserResponse{
			User: &userDTO.User,
		},
	), nil
}
