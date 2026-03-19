package repository

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
)

//go:generate mockery --name=UserRepository --output=../../../pkg/test/mocks --outpkg=mocks
type UserRepository interface {
	Create(
		ctx context.Context,
		req dto.UserCreateRequest,
	) (dto.UserCreateResponse, error)

	GetByID(
		ctx context.Context,
		req dto.UserGetByIDRequest,
	) (dto.UserGetByIDResponse, error)

	Update(
		ctx context.Context,
		req dto.UserUpdateRequest,
	) (dto.UserUpdateResponse, error)

	Delete(
		ctx context.Context,
		req dto.UserDeleteRequest,
	) error
}
