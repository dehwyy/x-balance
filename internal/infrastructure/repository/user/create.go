package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Create(
	ctx context.Context,
	req dto.UserCreateRequest,
) (dto.UserCreateResponse, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Implementation.Create", dspan.Attr("req", req))
	defer span.End()

	m := &models.User{
		Name:           req.Name.Value,
		OverdraftLimit: req.OverdraftLimit.Value,
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return dto.UserCreateResponse{}, span.Err(err)
	}

	response := dto.UserCreateResponse{User: *userconvert.ModelToUser(m)}
	span.WithAttribute("response", response)
	return response, nil
}
