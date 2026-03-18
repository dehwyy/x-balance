package mocks

import (
	"context"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (_m *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	ret := _m.Called(ctx, u)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*user.User), ret.Error(1)
}

func (_m *UserRepository) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	ret := _m.Called(ctx, id)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*user.User), ret.Error(1)
}

func (_m *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	ret := _m.Called(ctx, u)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*user.User), ret.Error(1)
}

func (_m *UserRepository) Delete(ctx context.Context, id user.ID) error {
	ret := _m.Called(ctx, id)
	return ret.Error(0)
}

var _ repository.UserRepository = &UserRepository{}
