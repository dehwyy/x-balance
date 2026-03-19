package snapshotconvert

import (
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func ModelToSnapshot(m *models.Snapshot) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        snapshot.ID{Value: m.ID},
		UserID:    m.UserID,
		Balance:   snapshot.Balance{Value: m.Balance},
		Version:   snapshot.Version{Value: m.Version},
		CreatedAt: m.CreatedAt,
	}
}

func SnapshotToModel(s *snapshot.Snapshot) *models.Snapshot {
	return &models.Snapshot{
		ID:        s.ID.Value,
		UserID:    s.UserID,
		Balance:   s.Balance.Value,
		Version:   s.Version.Value,
		CreatedAt: s.CreatedAt,
	}
}
