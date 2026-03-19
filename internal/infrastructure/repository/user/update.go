package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	userconvert "github.com/dehwyy/x-balance/internal/domain/entity/user/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) Update(
	ctx context.Context,
	u *user.User,
) (*user.User, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Update")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.User
	if err := db.Where("id = ? AND deleted_at IS NULL", u.ID.Value).First(&m).Error; err != nil {
		return nil, span.Err(err)
	}

	m.Name = u.Name.Value
	m.OverdraftLimit = u.OverdraftLimit.Value

	if err := db.Save(&m).Error; err != nil {
		return nil, span.Err(err)
	}

	return userconvert.ModelToUser(&m), nil
}
