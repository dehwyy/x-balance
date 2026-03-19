package balancehandler

import (
	"context"

	"github.com/dehwyy/tracerfx/pkg/tracer/dspan"
	"github.com/dehwyy/x-balance/internal/application/service/balanceservice"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	"github.com/shopspring/decimal"
)

func (h *Handler) Debit(ctx context.Context, req *balancev1.DebitRequest) (*balancev1.DebitResponse, error) {
	ctx, span := dspan.Start(ctx, "balanceDelivery.Debit", dspan.Attr("req", req))
	defer span.End()

	amount, _ := decimal.NewFromString(req.Amount)

	response, err := h.balanceservice.Debit(ctx, balanceservice.DebitRequest{
		UserID:        req.UserId,
		Amount:        amount,
		TransactionID: req.TransactionId,
	})
	if err != nil {
		return nil, span.Err(err)
	}

	protoResponse := &balancev1.DebitResponse{
		NewBalance:    response.NewBalance.String(),
		TransactionId: response.TransactionID,
	}
	span.WithAttribute("response", protoResponse)
	return protoResponse, nil
}
