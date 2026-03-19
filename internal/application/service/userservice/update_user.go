package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
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
	req *UpdateUserRequest,
) (*UpdateUserResponse, error) {
	ctx, span := dspan.Start(ctx, "userservice.Service.UpdateUser", dspan.Attr("req", req))
	defer span.End()

	var u user.User

	err := s.tx.Do(ctx, "userservice.UpdateUser", func(ctx context.Context) error {
		updateResp, err := s.userRepo.Update(ctx, dto.UserUpdateRequest{
			User: user.User{
				ID:             user.ID{Value: req.ID},
				Name:           user.Name{Value: req.Name},
				OverdraftLimit: user.OverdraftLimit{Value: req.OverdraftLimit},
			},
		})
		if err != nil {
			return err
		}
		u = updateResp.User
		return nil
	})
	if err != nil {
		return nil, span.Err(err)
	}

	response := &UpdateUserResponse{User: &u}
	span.WithAttribute("response", response)
	return response, nil
}
