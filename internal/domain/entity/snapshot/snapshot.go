package snapshot

import "time"

type Snapshot struct {
	ID        ID
	UserID    string
	Balance   Balance
	Version   Version
	CreatedAt time.Time
}
