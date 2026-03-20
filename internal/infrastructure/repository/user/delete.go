package userrepo

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Delete(
	ctx context.Context,
	req dto.UserDeleteRequest,
) error {
	ctx, span := dspan.Start(
		ctx,
		"userrepo.Implementation.Delete",
		dspan.Attr("req", req),
	)
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	now := time.Now()
	result := db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", string(req.ID)).
		Update("deleted_at", now)

	if result.Error != nil {
		return span.Err(result.Error)
	}

	return nil
}
