package transactionhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	transactionconvert "github.com/dehwyy/x-balance/internal/delivery/api/transaction/convert"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

func (h *Handler) GetTransaction(
	ctx context.Context,
	req *transactionspb.GetTransactionRequest,
) (*transactionspb.GetTransactionResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"transactionDelivery.GetTransaction",
		dspan.Attr("req", req),
	)
	defer span.End()

	response, err := h.transactionservice.GetTransaction(
		ctx,
		transactionconvert.GetTransactionRequestToDomain(req),
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(
		span,
		transactionconvert.GetTransactionResponseToProto(response),
	), nil
}
