package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetByID(
	ctx context.Context,
	req dto.UserGetByIDRequest,
) (dto.UserGetByIDResponse, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Implementation.GetByID", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.User
	if err := db.Where("id = ? AND deleted_at IS NULL", req.ID.Value).First(&m).Error; err != nil {
		return dto.UserGetByIDResponse{}, span.Err(err)
	}

	response := dto.UserGetByIDResponse{User: *userconvert.ModelToUser(&m)}
	span.WithAttribute("response", response)
	return response, nil
}