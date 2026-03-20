package snapshot

import (
	"time"

	"github.com/shopspring/decimal"

	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type Snapshot struct {
	ID        ID
	UserID    user.ID
	Balance   Balance
	Version   Version
	CreatedAt time.Time
}

// ComputeBalance вычисляет доступный и замороженный баланс относительно снапшота.
// sumAvailable — сумма кредитов и дебетов после снапшота.
// sumFrozen — сумма активных заморозок после снапшота.
func (s Snapshot) ComputeBalance(
	sumAvailable decimal.Decimal,
	sumFrozen decimal.Decimal,
) (available decimal.Decimal, frozen decimal.Decimal) {
	available = decimal.Decimal(s.Balance).Add(sumAvailable).Sub(sumFrozen)
	frozen = sumFrozen
	return
}
