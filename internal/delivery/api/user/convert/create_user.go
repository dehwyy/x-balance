package convert

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateUserRequestToDomain(req *balancepb.CreateUserRequest) userservice.CreateUserRequest {
	limit, _ := decimal.NewFromString(req.OverdraftLimit)
	return userservice.CreateUserRequest{
		Name:           req.Name,
		OverdraftLimit: limit,
	}
}

func CreateUserResponseToProto(res *userservice.CreateUserResponse) *balancepb.CreateUserResponse {
	return &balancepb.CreateUserResponse{
		User: userToProto(res),
	}
}

func userToProto(res *userservice.CreateUserResponse) *balancepb.User {
	u := res.User
	pb := &balancepb.User{
		Id:             u.ID.Value,
		Name:           u.Name.Value,
		OverdraftLimit: u.OverdraftLimit.Value.String(),
		CreatedAt:      timestamppb.New(u.CreatedAt),
		UpdatedAt:      timestamppb.New(u.UpdatedAt),
	}
	if u.DeletedAt != nil {
		pb.DeletedAt = timestamppb.New(*u.DeletedAt)
	}
	return pb
}
