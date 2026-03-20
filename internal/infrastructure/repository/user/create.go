package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.UserCreateRequest,
) (dto.UserCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	userEntity := &user.User{
		Name:           req.Name,
		OverdraftLimit: req.OverdraftLimit,
	}
	m := userconvert.UserToModel(userEntity)

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.UserCreateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.UserCreateResponse{User: *userconvert.ModelToUser(m)}), nil
}
