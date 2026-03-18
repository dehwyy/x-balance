package userservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/pkg/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const testUserID = "user-123"

func newUserService(userRepo *mocks.UserRepository) *userservice.Service {
	tx := &mocks.TxManager{}
	return userservice.New(userservice.Opts{
		TX:       tx,
		UserRepo: userRepo,
	})
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}
	overdraft, _ := decimal.NewFromString("100")

	expectedUser := &user.User{
		ID:             user.ID{Value: "new-user-id"},
		Name:           user.Name{Value: "John Doe"},
		OverdraftLimit: user.OverdraftLimit{Value: overdraft},
	}

	userRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(expectedUser, nil)

	svc := newUserService(userRepo)
	resp, err := svc.CreateUser(ctx, userservice.CreateUserRequest{
		Name:           "John Doe",
		OverdraftLimit: overdraft,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedUser, resp.User)
	userRepo.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}
	overdraft, _ := decimal.NewFromString("100")

	expectedErr := errors.New("repository error")
	userRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil, expectedErr)

	svc := newUserService(userRepo)
	_, err := svc.CreateUser(ctx, userservice.CreateUserRequest{
		Name:           "John Doe",
		OverdraftLimit: overdraft,
	})

	assert.ErrorIs(t, err, expectedErr)
	userRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}
	overdraft, _ := decimal.NewFromString("100")

	expectedUser := &user.User{
		ID:             user.ID{Value: testUserID},
		Name:           user.Name{Value: "John Doe"},
		OverdraftLimit: user.OverdraftLimit{Value: overdraft},
	}

	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(expectedUser, nil)

	svc := newUserService(userRepo)
	resp, err := svc.GetUser(ctx, userservice.GetUserRequest{
		ID: testUserID,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedUser, resp.User)
	userRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}

	userRepo.On("GetByID", ctx, user.ID{Value: testUserID}).Return(nil, gorm.ErrRecordNotFound)

	svc := newUserService(userRepo)
	_, err := svc.GetUser(ctx, userservice.GetUserRequest{
		ID: testUserID,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	userRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}
	overdraft, _ := decimal.NewFromString("200")

	updatedUser := &user.User{
		ID:             user.ID{Value: testUserID},
		Name:           user.Name{Value: "John Smith"},
		OverdraftLimit: user.OverdraftLimit{Value: overdraft},
	}

	userRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(updatedUser, nil)

	svc := newUserService(userRepo)
	resp, err := svc.UpdateUser(ctx, userservice.UpdateUserRequest{
		ID:             testUserID,
		Name:           "John Smith",
		OverdraftLimit: overdraft,
	})

	require.NoError(t, err)
	assert.Equal(t, updatedUser, resp.User)
	userRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}
	overdraft, _ := decimal.NewFromString("200")

	userRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil, gorm.ErrRecordNotFound)

	svc := newUserService(userRepo)
	_, err := svc.UpdateUser(ctx, userservice.UpdateUserRequest{
		ID:             testUserID,
		Name:           "John Smith",
		OverdraftLimit: overdraft,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	userRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}

	userRepo.On("Delete", ctx, user.ID{Value: testUserID}).Return(nil)

	svc := newUserService(userRepo)
	err := svc.DeleteUser(ctx, userservice.DeleteUserRequest{
		ID: testUserID,
	})

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.UserRepository{}

	userRepo.On("Delete", ctx, user.ID{Value: testUserID}).Return(gorm.ErrRecordNotFound)

	svc := newUserService(userRepo)
	err := svc.DeleteUser(ctx, userservice.DeleteUserRequest{
		ID: testUserID,
	})

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	userRepo.AssertExpectations(t)
}
