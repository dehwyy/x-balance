package transactionhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	transactionv1 "github.com/dehwyy/x-balance/internal/generated/pb/common/transaction/v1"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

func (h *Handler) ListTransactions(
	ctx context.Context,
	req *transactionspb.ListTransactionsRequest,
) (*transactionspb.ListTransactionsResponse, error) {
	ctx, span := dspan.Start(ctx, "transactionDelivery.ListTransactions", dspan.Attr("req", req))
	defer span.End()

	domainReq := transactionservice.ListTransactionsRequest{
		UserID: req.UserId,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}
	if req.From != nil {
		t := req.From.AsTime()
		domainReq.From = &t
	}
	if req.To != nil {
		t := req.To.AsTime()
		domainReq.To = &t
	}

	response, err := h.transactionservice.ListTransactions(ctx, domainReq)
	if err != nil {
		return nil, span.Err(err)
	}

	txs := make([]*transactionv1.Transaction, len(response.Events))
	for i, e := range response.Events {
		txs[i] = eventconvert.EventToProto(e)
	}

	protoResponse := &transactionspb.ListTransactionsResponse{
		Transactions: txs,
		Total:        response.Total,
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
