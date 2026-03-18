package userservice

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type DeleteUserRequest struct {
	ID string
}

func (s *Service) DeleteUser(
	ctx context.Context,
	req DeleteUserRequest,
) error {
	return s.tx.Do(ctx, "userservice.DeleteUser", func(ctx context.Context) error {
		return s.userRepo.Delete(ctx, user.ID{Value: req.ID})
	})
}
