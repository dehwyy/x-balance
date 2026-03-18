package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/domain/entity/snapshot"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/shopspring/decimal"
)

func (impl *Implementation) SumSinceSnapshot(
	ctx context.Context,
	userID string,
	snapshotID snapshot.ID,
) (decimal.Decimal, decimal.Decimal, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.SumSinceSnapshot")
	defer span.End()

	db := impl.tx.GetConnection(ctx)

	var snapshotModel models.Snapshot
	if err := db.Where("id = ?", snapshotID.Value).First(&snapshotModel).Error; err != nil {
		return decimal.Zero, decimal.Zero, span.Err(err)
	}

	type SumResult struct {
		Total decimal.Decimal
	}

	// Sum of non-freeze events (credit/debit) for available balance
	var balanceResult SumResult
	if err := db.Model(&models.Event{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND created_at > ? AND type NOT IN (?, ?)",
			userID, snapshotModel.CreatedAt,
			event.TypeFreezeHold.Value, event.TypeFreezeRelease.Value,
		).
		Scan(&balanceResult).Error; err != nil {
		return decimal.Zero, decimal.Zero, span.Err(err)
	}

	// Sum of active freeze_hold events (minus released)
	var frozenResult SumResult
	if err := db.Model(&models.Event{}).
		Select("COALESCE(SUM(CASE WHEN type = ? THEN amount WHEN type = ? THEN -amount ELSE 0 END), 0) as total",
			event.TypeFreezeHold.Value, event.TypeFreezeRelease.Value).
		Where("user_id = ?", userID).
		Scan(&frozenResult).Error; err != nil {
		return decimal.Zero, decimal.Zero, span.Err(err)
	}

	return balanceResult.Total, frozenResult.Total, nil
}
