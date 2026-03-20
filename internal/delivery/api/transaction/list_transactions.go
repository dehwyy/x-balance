package transactionhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	transactionconvert "github.com/dehwyy/x-balance/internal/delivery/api/transaction/convert"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

func (h *Handler) ListTransactions(
	ctx context.Context,
	req *transactionspb.ListTransactionsRequest,
) (*transactionspb.ListTransactionsResponse, error) {
	ctx, span := dspan.Start(
		ctx,
		"transactionDelivery.ListTransactions",
		dspan.Attr("req", req),
	)
	defer span.End()

	response, err := h.transactionservice.ListTransactions(
		ctx,
		transactionconvert.ListTransactionsRequestToDomain(req),
	)
	if err != nil {
		return nil, span.Err(err)
	}

	return dspan.Response(
		span,
		transactionconvert.ListTransactionsResponseToProto(response),
	), nil
}
