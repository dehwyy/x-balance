package userservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type GetUserRequest struct {
	ID string
}

type GetUserResponse struct {
	User *user.User
}

func (s *Service) GetUser(
	ctx context.Context,
	req GetUserRequest,
) (*GetUserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, user.ID{Value: req.ID})
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{User: u}, nil
}
