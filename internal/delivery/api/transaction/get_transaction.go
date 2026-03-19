package transactionhandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/transactionservice"
	eventconvert "github.com/dehwyy/x-balance/internal/domain/entity/event/convert"
	transactionspb "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
)

func (h *Handler) GetTransaction(
	ctx context.Context,
	req *transactionspb.GetTransactionRequest,
) (*transactionspb.GetTransactionResponse, error) {
	ctx, span := dspan.Start(ctx, "transactionDelivery.GetTransaction", dspan.Attr("req", req))
	defer span.End()

	response, err := h.transactionservice.GetTransaction(ctx, transactionservice.GetTransactionRequest{
		UserID: req.UserId,
		TxID:   req.TxId,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &transactionspb.GetTransactionResponse{
		Transaction: eventconvert.EventToProto(response.Event),
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
