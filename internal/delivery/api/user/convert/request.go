package userconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
	"github.com/shopspring/decimal"
)

func CreateUserRequestToDomain(req *userspb.CreateUserRequest) *userservice.CreateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return &userservice.CreateUserRequest{
		Name:           req.Name,
		OverdraftLimit: limit,
	}
}

func GetUserRequestToDomain(req *userspb.GetUserRequest) *userservice.GetUserRequest {
	return &userservice.GetUserRequest{ID: req.Id}
}

func UpdateUserRequestToDomain(req *userspb.UpdateUserRequest) *userservice.UpdateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return &userservice.UpdateUserRequest{
		ID:             req.Id,
		Name:           req.Name,
		OverdraftLimit: limit,
	}
}

func DeleteUserRequestToDomain(req *userspb.DeleteUserRequest) *userservice.DeleteUserRequest {
	return &userservice.DeleteUserRequest{ID: req.Id}
}
