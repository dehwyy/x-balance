package userrepo

import (
	"context"

	"gorm.io/gorm"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	"github.com/dehwyy/x-balance/internal/domain/repository"
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
	if err := db.Where("id = ? AND deleted_at IS NULL", string(req.User.ID)).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.UserUpdateResponse{}, repository.ErrNotFound
		}
		return dto.UserUpdateResponse{}, span.Err(err)
	}

	updated := userconvert.UserToModel(&req.User)
	m.Name = updated.Name
	m.OverdraftLimit = updated.OverdraftLimit

	if err := db.Save(&m).Error; err != nil {
		return dto.UserUpdateResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.UserUpdateResponse{User: *userconvert.ModelToUser(&m)}), nil
}
