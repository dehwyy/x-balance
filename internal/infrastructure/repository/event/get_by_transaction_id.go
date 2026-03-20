package eventrepo

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/dto"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	"github.com/dehwyy/x-balance/internal/infrastructure/repository/models"
)

func (impl *Implementation) GetByTransactionID(
	ctx context.Context,
	req dto.EventGetByTxIDRequest,
) (dto.EventGetByTxIDResponse, error) {
	ctx, span := dspan.Start(ctx, "eventrepo.Implementation.GetByTransactionID", dspan.Attr("req", req))
	defer span.End()

	db := impl.tx.GetConnection(ctx)
	var m models.Event
	if err := db.Where("transaction_id = ?", req.TransactionID.Value).First(&m).Error; err != nil {
		return dto.EventGetByTxIDResponse{}, span.Err(err)
	}

	return dspan.Response(span, dto.EventGetByTxIDResponse{Event: *eventconvert.ModelToEvent(&m)}), nil
}
