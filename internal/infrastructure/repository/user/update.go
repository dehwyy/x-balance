package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Update(
	ctx context.Context,
	req dto.UserUpdateRequest,
) (dto.UserUpdateResponse, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Implementation.Update", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.User
	if err := db.Where("id = ? AND deleted_at IS NULL", req.User.ID.Value).First(&m).Error; err != nil {
		return dto.UserUpdateResponse{}, span.Err(err)
	}

	m.Name = req.User.Name.Value
	m.OverdraftLimit = req.User.OverdraftLimit.Value

	if err := db.Save(&m).Error; err != nil {
		return dto.UserUpdateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.UserUpdateResponse{User: *userconvert.ModelToUser(&m)}), nil
}
