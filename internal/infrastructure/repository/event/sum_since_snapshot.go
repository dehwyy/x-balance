package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/shopspring/decimal"
)

func (impl *Implementation) SumSinceSnapshot(
	ctx context.Context,
	req dto.EventSumSinceSnapshotRequest,
) (dto.EventSumSinceSnapshotResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.SumSinceSnapshot", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)

	var snapshotModel models.Snapshot
	if err := db.Where("id = ?", string(req.SnapshotID)).First(&snapshotModel).Error; err != nil {
		return dto.EventSumSinceSnapshotResponse{}, span.Err(err)
	}

	type SumResult struct {
		Total decimal.Decimal
	}

	var balanceResult SumResult
	if err := db.Model(&models.Event{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND created_at > ? AND type NOT IN (?, ?)",
			string(req.UserID), snapshotModel.CreatedAt,
			"freeze_hold", "freeze_release",
		).
		Scan(&balanceResult).Error; err != nil {
		return dto.EventSumSinceSnapshotResponse{}, span.Err(err)
	}

	var frozenResult SumResult
	if err := db.Model(&models.Event{}).
		Select("COALESCE(SUM(CASE WHEN type = ? THEN amount WHEN type = ? THEN -amount ELSE 0 END), 0) as total",
			"freeze_hold", "freeze_release").
		Where("user_id = ?", string(req.UserID)).
		Scan(&frozenResult).Error; err != nil {
		return dto.EventSumSinceSnapshotResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventSumSinceSnapshotResponse{Available: balanceResult.Total, Frozen: frozenResult.Total}), nil
}
