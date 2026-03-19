package dto

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
)

type SnapshotCreateRequest struct {
	UserID  user.ID
	Balance snapshot.Balance
	Version snapshot.Version
}

type SnapshotCreateResponse struct {
	Snapshot snapshot.Snapshot
}

type SnapshotGetLatestByUserIDRequest struct {
	UserID user.ID
}

type SnapshotGetLatestByUserIDResponse struct {
	Snapshot snapshot.Snapshot
}

type SnapshotUpdateVersionRequest struct {
	Snapshot snapshot.Snapshot
}
