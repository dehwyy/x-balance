package userconvert

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToUser(m *models.User) *user.User {
	return &user.User{
		ID:             user.ID{Value: m.ID},
		Name:           user.Name{Value: m.Name},
		OverdraftLimit: user.OverdraftLimit{Value: m.OverdraftLimit},
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

func UserToModel(u *user.User) *models.User {
	return &models.User{
		ID:             u.ID.Value,
		Name:           u.Name.Value,
		OverdraftLimit: u.OverdraftLimit.Value,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt,
	}
}
