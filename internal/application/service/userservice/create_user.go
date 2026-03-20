package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"

	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type CreateUserRequest struct {
	Name           user.Name
	OverdraftLimit user.OverdraftLimit
}

type CreateUserResponse struct {
	User *user.User
}

func (s *Service) CreateUser(
	ctx context.Context,
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"userservice.Service.CreateUser",
		dspan.Attr("req", req),
	)
	defer span.End()

	var u user.User

	err := s.tx.Do(
		ctx,
		"userservice.CreateUser",
		func(ctx context.Context) error {
			createDTO, err := s.userRepo.Create(
				ctx,
				dto.UserCreateRequest{
					Name:           req.Name,
					OverdraftLimit: req.OverdraftLimit,
				},
			)
			if err != nil {
				return err
			}
			u = createDTO.User
			return nil
		},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(
		span,
		&CreateUserResponse{
			User: &u,
		},
	), nil
}
