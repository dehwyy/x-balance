package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
)

func (h *Handler) Unfreeze(ctx context.Context, req *balancev1.UnfreezeRequest) (*balancev1.UnfreezeResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Unfreeze", dspan.Attr("req", req))
	defer span.End()

	response, err := h.balanceservice.Unfreeze(ctx, balanceservice.UnfreezeRequest{
		UserID:        req.UserId,
		TransactionID: req.TransactionId,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &balancev1.UnfreezeResponse{
		UnfrozenAmount: response.UnfrozenAmount.String(),
		TransactionId:  response.TransactionID,
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
