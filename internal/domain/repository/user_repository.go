package repository

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
)

//go:generate mockery --name=UserRepository --output=../../../pkg/test/mocks --outpkg=mocks
type UserRepository interface {
	Create(
		ctx context.Context,
		u *user.User,
	) (*user.User, error)

	GetByID(
		ctx context.Context,
		id user.ID,
	) (*user.User, error)

	Update(
		ctx context.Context,
		u *user.User,
	) (*user.User, error)

	Delete(
		ctx context.Context,
		id user.ID,
	) error
}
