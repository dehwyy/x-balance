package convert

import (
	"github.com/dehwyy/x-balance/internal/application/service/userservice"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetUserRequestToDomain(req *balancepb.GetUserRequest) userservice.GetUserRequest {
	return userservice.GetUserRequest{ID: req.Id}
}

func GetUserResponseToProto(res *userservice.GetUserResponse) *balancepb.GetUserResponse {
	u := res.User
	pb := &balancepb.GetUserResponse{
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
