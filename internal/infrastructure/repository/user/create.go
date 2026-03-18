package userrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models/convert"
)

func (impl *Implementation) Create(
	ctx context.Context,
	u *user.User,
) (*user.User, error) {
	ctx, span := dspan.Start(ctx, "userrepo.Create")
	defer span.End()

	m := &models.User{
		Name:           u.Name.Value,
		OverdraftLimit: u.OverdraftLimit.Value,
	}

	db := impl.tx.GetConnection(ctx)
	if err := db.Create(m).Error; err != nil {
		return nil, span.Err(err)
	}

	return convert.UserToDomain(m), nil
}
