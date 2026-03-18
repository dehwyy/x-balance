package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/domain/entity/event"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models/convert"
)

func (impl *Implementation) GetByTransactionID(
	ctx context.Context,
	txID event.TransactionID,
) (*event.Event, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.GetByTransactionId")
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Event
	if err := db.Where("transaction_id = ?", txID.Value).First(&m).Error; err != nil {
		return nil, span.Err(err)
	}

	return convert.EventToDomain(&m), nil
}
