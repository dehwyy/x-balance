package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	"github.com/shopspring/decimal"
)

func (h *Handler) Freeze(ctx context.Context, req *balancev1.FreezeRequest) (*balancev1.FreezeResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Freeze", dspan.Attr("req", req))
	defer span.End()

	amount, _ := decimal.NewFromString(req.Amount)

	response, err := h.balanceservice.Freeze(ctx, balanceservice.FreezeRequest{
		UserID:               req.UserId,
		Amount:               amount,
		TransactionID:        req.TransactionId,
		FreezeTimeoutSeconds: req.FreezeTimeoutSeconds,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &balancev1.FreezeResponse{
		FrozenAmount:  response.FrozenAmount.String(),
		TransactionId: response.TransactionID,
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
