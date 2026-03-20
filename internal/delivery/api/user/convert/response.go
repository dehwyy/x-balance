package userconvert

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	domainconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	userspb "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

func CreateUserResponseToProto(resp *userservice.CreateUserResponse) *userspb.CreateUserResponse {
	return &userspb.CreateUserResponse{
		User: domainconvert.UserToProto(resp.User),
	}
}

func GetUserResponseToProto(resp *userservice.GetUserResponse) *userspb.GetUserResponse {
	return &userspb.GetUserResponse{
		User: domainconvert.UserToProto(resp.User),
	}
}

func UpdateUserResponseToProto(resp *userservice.UpdateUserResponse) *userspb.UpdateUserResponse {
	return &userspb.UpdateUserResponse{
		User: domainconvert.UserToProto(resp.User),
	}
}
