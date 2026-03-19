package snapshot

import (
	"time"

	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type Snapshot struct {
	ID        ID
	UserID    user.ID
	Balance   Balance
	Version   Version
	CreatedAt time.Time
}
