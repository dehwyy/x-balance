package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"

	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type UpdateUserRequest struct {
	ID             user.ID
	Name           user.Name
	OverdraftLimit user.OverdraftLimit
}

type UpdateUserResponse struct {
	User *user.User
}

func (s *Service) UpdateUser(
	ctx context.Context,
	req *UpdateUserRequest,
) (*UpdateUserResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"userservice.Service.UpdateUser",
		dspan.Attr("req", req),
	)
	defer span.End()

	var u user.User

	err := s.tx.Do(
		ctx,
		"userservice.UpdateUser",
		func(ctx context.Context) error {
			updateDTO, err := s.userRepo.Update(
				ctx,
				dto.UserUpdateRequest{
					User: user.New(req.ID, req.Name, req.OverdraftLimit),
				},
			)
			if err != nil {
				return err
			}
			u = updateDTO.User
			return nil
		},
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(span, &UpdateUserResponse{User: &u}), nil
}
