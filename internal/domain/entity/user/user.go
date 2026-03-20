package user

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID             ID
	Name           Name
	OverdraftLimit OverdraftLimit
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

func New(id ID, name Name, overdraftLimit OverdraftLimit) User {
	return User{
		ID:             id,
		Name:           name,
		OverdraftLimit: overdraftLimit,
	}
}

// CanDebit проверяет, не превышает ли списание лимит овердрафта.
func (u User) CanDebit(currentAvailable decimal.Decimal, amount decimal.Decimal) bool {
	minAllowed := u.OverdraftLimit.Value.Neg()
	return currentAvailable.Sub(amount).GreaterThanOrEqual(minAllowed)
}
