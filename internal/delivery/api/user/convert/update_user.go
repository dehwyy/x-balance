package convert

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UpdateUserRequestToDomain(req *balancepb.UpdateUserRequest) userservice.UpdateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return userservice.UpdateUserRequest{
		ID:             req.Id,
		Name:           req.Name,
		OverdraftLimit: limit,
	}
}

func UpdateUserResponseToProto(res *userservice.UpdateUserResponse) *balancepb.UpdateUserResponse {
	u := res.User
	pb := &balancepb.UpdateUserResponse{
		User: &balancepb.User{
			Id:             u.ID.Value,
			Name:           u.Name.Value,
			OverdraftLimit: u.OverdraftLimit.Value.String(),
			CreatedAt:      timestamppb.New(u.CreatedAt),
			UpdatedAt:      timestamppb.New(u.UpdatedAt),
		},
	}
	if u.DeletedAt != nil {
		pb.User.DeletedAt = timestamppb.New(*u.DeletedAt)
	}
	return pb
}

func DeleteUserRequestToDomain(req *balancepb.DeleteUserRequest) userservice.DeleteUserRequest {
	return userservice.DeleteUserRequest{ID: req.Id}
}
