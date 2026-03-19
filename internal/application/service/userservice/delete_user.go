package userservice

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type DeleteUserRequest struct {
	ID string
}

func (s *Service) DeleteUser(
	ctx context.Context,
	req *DeleteUserRequest,
) error {
	ctx, span := dspan.Start(ctx, "userservice.Service.DeleteUser", dspan.Attr("req", req))
	defer span.End()

	err := s.tx.Do(ctx, "userservice.DeleteUser", func(ctx context.Context) error {
		return s.userRepo.Delete(ctx, dto.UserDeleteRequest{ID: user.ID{Value: req.ID}})
	})
	if err != nil {
		return span.Err(err)
	}
	return nil
}
