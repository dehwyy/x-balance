package snapshotconvert

import (
	"github.com/shopspring/decimal"

	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	user "github.com/dehwyy/x-balance/internal/domain/entity/user"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToSnapshot(m *models.Snapshot) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        snapshot.ID(m.ID),
		UserID:    user.ID(m.UserID),
		Balance:   snapshot.Balance(m.Balance),
		Version:   snapshot.Version(m.Version),
		CreatedAt: m.CreatedAt,
	}
}

func SnapshotToModel(s *snapshot.Snapshot) *models.Snapshot {
	return &models.Snapshot{
		ID:        string(s.ID),
		UserID:    string(s.UserID),
		Balance:   decimal.Decimal(s.Balance),
		Version:   int64(s.Version),
		CreatedAt: s.CreatedAt,
	}
}
