package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (_m *UserRepository) Create(ctx context.Context, req dto.UserCreateRequest) (dto.UserCreateResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.UserCreateResponse), ret.Error(1)
}

func (_m *UserRepository) GetByID(ctx context.Context, req dto.UserGetByIDRequest) (dto.UserGetByIDResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.UserGetByIDResponse), ret.Error(1)
}

func (_m *UserRepository) Update(ctx context.Context, req dto.UserUpdateRequest) (dto.UserUpdateResponse, error) {
	ret := _m.Called(ctx, req)
	return ret.Get(0).(dto.UserUpdateResponse), ret.Error(1)
}

func (_m *UserRepository) Delete(ctx context.Context, req dto.UserDeleteRequest) error {
	ret := _m.Called(ctx, req)
	return ret.Error(0)
}

var _ repository.UserRepository = &UserRepository{}
