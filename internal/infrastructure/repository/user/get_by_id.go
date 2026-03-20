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

func (impl *Implementation) GetByID(
	ctx context.Context,
	req dto.UserGetByIDRequest,
) (dto.UserGetByIDResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"userrepo.Implementation.GetByID",
		dspan.Attr("req", req),
	)
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.User
	if err := db.Where("id = ? AND deleted_at IS NULL", string(req.ID)).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.UserGetByIDResponse{}, repository.ErrNotFound
		}
		return dto.UserGetByIDResponse{}, span.Err(err)
	}

	return dspan.Response(
		span,
		dto.UserGetByIDResponse{
			User: *userconvert.ModelToUser(&m),
		},
	), nil
}
