package userconvert

import (
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToUser(m *models.User) *user.User {
	return &user.User{
		ID:             user.ID(m.ID),
		Name:           user.Name(m.Name),
		OverdraftLimit: user.OverdraftLimit(m.OverdraftLimit),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

func UserToModel(u *user.User) *models.User {
	return &models.User{
		ID:             string(u.ID),
		Name:           string(u.Name),
		OverdraftLimit: decimal.Decimal(u.OverdraftLimit),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt,
	}
}
