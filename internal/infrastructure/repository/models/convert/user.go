package convert

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/shopspring/decimal"
)

func UserToModel(u *user.User) *models.User {
	return &models.User{
		ID:             u.ID.Value,
		Name:           u.Name.Value,
		OverdraftLimit: u.OverdraftLimit.Value,
	}
}

func UserToDomain(m *models.User) *user.User {
	return &user.User{
		ID:             user.ID{Value: m.ID},
		Name:           user.Name{Value: m.Name},
		OverdraftLimit: user.OverdraftLimit{Value: m.OverdraftLimit},
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

func NewUser(name string, overdraftLimit decimal.Decimal) *models.User {
	return &models.User{
		Name:           name,
		OverdraftLimit: overdraftLimit,
	}
}
