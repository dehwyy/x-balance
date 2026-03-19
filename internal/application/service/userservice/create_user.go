package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
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
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userservice.Service.CreateUser", dspan.Attr("req", req))
	defer span.End()

	var u user.User

	err := s.tx.Do(ctx, "userservice.CreateUser", func(ctx context.Context) error {
		createResp, err := s.userRepo.Create(ctx, dto.UserCreateRequest{
			Name:           user.Name{Value: req.Name},
			OverdraftLimit: user.OverdraftLimit{Value: req.OverdraftLimit},
		})
		if err != nil {
			return err
		}
		u = createResp.User
		return nil
	})
	if err != nil {
		return nil, span.Err(err)
	}

	response := &CreateUserResponse{User: &u}
	span.WithAttribute("response", response)
	return response, nil
}
