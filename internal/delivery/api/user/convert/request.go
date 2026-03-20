package userconvert

import (
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func CreateUserRequestToDomain(req *userspb.CreateUserRequest) *userservice.CreateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return &userservice.CreateUserRequest{
		Name:           user.NewName(req.Name),
		OverdraftLimit: user.NewOverdraftLimit(limit),
	}
}

func GetUserRequestToDomain(req *userspb.GetUserRequest) *userservice.GetUserRequest {
	return &userservice.GetUserRequest{
		ID: user.NewID(req.Id),
	}
}

func UpdateUserRequestToDomain(req *userspb.UpdateUserRequest) *userservice.UpdateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return &userservice.UpdateUserRequest{
		ID:             user.NewID(req.Id),
		Name:           user.NewName(req.Name),
		OverdraftLimit: user.NewOverdraftLimit(limit),
	}
}

func DeleteUserRequestToDomain(req *userspb.DeleteUserRequest) *userservice.DeleteUserRequest {
	return &userservice.DeleteUserRequest{
		ID: user.NewID(req.Id),
	}
}
