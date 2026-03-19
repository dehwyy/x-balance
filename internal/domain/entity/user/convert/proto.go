package userconvert

import (
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	userv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/user/v1"
)

func UserToProto(u *user.User) *userv1.User {
	proto := &userv1.User{
		Id:             u.ID.Value,
		Name:           u.Name.Value,
		OverdraftLimit: u.OverdraftLimit.Value.String(),
		CreatedAt:      timestamppb.New(u.CreatedAt),
		UpdatedAt:      timestamppb.New(u.UpdatedAt),
	}

	if u.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*u.DeletedAt)
	}

	return proto
}

func ProtoToUser(p *userv1.User) *user.User {
	overdraftLimit, _ := decimal.NewFromString(p.OverdraftLimit)

	u := &user.User{
		ID:             user.ID{Value: p.Id},
		Name:           user.Name{Value: p.Name},
		OverdraftLimit: user.OverdraftLimit{Value: overdraftLimit},
		CreatedAt:      p.CreatedAt.AsTime(),
		UpdatedAt:      p.UpdatedAt.AsTime(),
	}

	if p.DeletedAt != nil {
		deletedAt := p.DeletedAt.AsTime()
		u.DeletedAt = &deletedAt
	}

	return u
}
