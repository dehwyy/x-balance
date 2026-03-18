package userservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/shopspring/decimal"
)

type CreateUserRequest struct {
	Name           string
	OverdraftLimit decimal.Decimal
}

type CreateUserResponse struct {
	User *user.User
}

func (s *Service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (*CreateUserResponse, error) {
	var u *user.User

	err := s.tx.Do(ctx, "userservice.CreateUser", func(ctx context.Context) error {
		var err error
		u, err = s.userRepo.Create(ctx, &user.User{
			Name:           user.Name{Value: req.Name},
			OverdraftLimit: user.OverdraftLimit{Value: req.OverdraftLimit},
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{User: u}, nil
}
