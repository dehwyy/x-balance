package userrepo

import (
	"context"
	"time"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Delete(
	ctx context.Context,
	id user.ID,
) error {
	ctx, span := dspan.Start(ctx, "userrepo.Delete")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	now := time.Now()
	result := db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id.Value).
		Update("deleted_at", now)

	if result.Error != nil {
		return span.Err(result.Error)
	}

	return nil
}
