package userservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type UpdateUserRequest struct {
	ID             string
	Name           string
	OverdraftLimit decimal.Decimal
}

type UpdateUserResponse struct {
	User *user.User
}

func (s *Service) UpdateUser(
	ctx context.Context,
	req UpdateUserRequest,
) (*UpdateUserResponse, error) {
	var u *user.User

	err := s.tx.Do(ctx, "userservice.UpdateUser", func(ctx context.Context) error {
		var err error
		u, err = s.userRepo.Update(ctx, &user.User{
			ID:             user.ID{Value: req.ID},
			Name:           user.Name{Value: req.Name},
			OverdraftLimit: user.OverdraftLimit{Value: req.OverdraftLimit},
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &UpdateUserResponse{User: u}, nil
}
